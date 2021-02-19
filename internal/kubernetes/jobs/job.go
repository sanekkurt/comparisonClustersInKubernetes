package jobs

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/metadata"
	"k8s-cluster-comparator/internal/kubernetes/pods"
	"k8s-cluster-comparator/internal/logging"

	"sync"

	v12 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/kubernetes/types"
)

const (
	jobKind          = "job"
	objectBatchLimit = 25
)

func addItemsToJobList(ctx context.Context, clientSet kubernetes.Interface, namespace string, limit int64) (*v12.JobList, error) {
	log := logging.FromContext(ctx)

	log.Debugf("addItemsToJobList started")
	defer log.Debugf("addItemsToJobList completed")

	var (
		batch *v12.JobList
		jobs  = &v12.JobList{
			Items: make([]v12.Job, 0),
		}

		continueToken string

		err error
	)

	for {
		batch, err = clientSet.BatchV1().Jobs(namespace).List(metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		})
		if err != nil {
			return nil, err
		}

		log.Debugf("addItemsToJobList: %d objects received", len(batch.Items))

		jobs.Items = append(jobs.Items, batch.Items...)

		jobs.TypeMeta = batch.TypeMeta
		jobs.ListMeta = batch.ListMeta

		if batch.Continue == "" {
			break
		}

		continueToken = batch.Continue
	}

	jobs.Continue = ""

	return jobs, err
}

type JobsComparator struct {
}

func NewJobsComparator(ctx context.Context, namespace string) JobsComparator {
	return JobsComparator{}
}

// Compare compare Jobs in different clusters
func (cmp JobsComparator) Compare(ctx context.Context, namespace string) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("kind", jobKind))
		cfg = config.FromContext(ctx)
	)
	ctx = logging.WithLogger(ctx, log)

	jobs1, err := addItemsToJobList(ctx, cfg.Connections.Cluster1.ClientSet, namespace, objectBatchLimit)
	if err != nil {
		return nil, fmt.Errorf("cannot obtain jobs list from 1st cluster: %w", err)
	}

	jobs2, err := addItemsToJobList(ctx, cfg.Connections.Cluster2.ClientSet, namespace, objectBatchLimit)
	if err != nil {
		return nil, fmt.Errorf("cannot obtain jobs list from 2st cluster: %w", err)
	}

	mapJobs1, mapJobs2 := prepareJobsMaps(ctx, jobs1, jobs2)

	_ = setInformationAboutJobs(ctx, mapJobs1, mapJobs2, jobs1, jobs2, namespace)

	return nil, nil
}

// prepareJobsMaps add value secrets in map
func prepareJobsMaps(ctx context.Context, jobs1, jobs2 *v12.JobList) (map[string]types.IsAlreadyComparedFlag, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)

		mapJobs1 = make(map[string]types.IsAlreadyComparedFlag)
		mapJobs2 = make(map[string]types.IsAlreadyComparedFlag)

		indexCheck types.IsAlreadyComparedFlag
	)

OUTER1:
	for index, value := range jobs1.Items {
		if cfg.Skips.IsSkippedEntity(jobKind, value.Name) {
			log.With(zap.String("name", value.Name)).Debugf("job/%s is skipped from comparison", value.Name)
			continue
		}

		if value.OwnerReferences != nil {
			for _, owner := range value.OwnerReferences {
				if owner.Kind == "CronJob" {
					log.Debugf("job/%s is skipped from comparison because it is owned by CronJob", value.Name)
					continue OUTER1
				}
			}
		}

		indexCheck.Index = index
		mapJobs1[value.Name] = indexCheck

	}

OUTER2:
	for index, value := range jobs2.Items {
		if cfg.Skips.IsSkippedEntity(jobKind, value.Name) {
			log.With(zap.String("name", value.Name)).Debugf("job/%s is skipped from comparison", value.Name)
			continue
		}

		if value.OwnerReferences != nil {
			for _, owner := range value.OwnerReferences {
				if owner.Kind == "CronJob" {
					log.Debugf("job/%s is skipped from comparison because it is owned by CronJob", value.Name)
					continue OUTER2
				}
			}

		}

		indexCheck.Index = index
		mapJobs2[value.Name] = indexCheck

	}

	return mapJobs1, mapJobs2
}

// setInformationAboutJobs set information about jobs
func setInformationAboutJobs(ctx context.Context, map1, map2 map[string]types.IsAlreadyComparedFlag, jobs1, jobs2 *v12.JobList, namespace string) bool {
	var (
		log  = logging.FromContext(ctx)
		flag bool
	)

	if len(map1) != len(map2) {
		log.Warnw("object counts are different", zap.Int("objectsCount1st", len(map1)), zap.Int("objectsCount2nd", len(map2)))
		flag = true
	}

	wg := &sync.WaitGroup{}
	channel := make(chan bool, len(map1))

	for name, index1 := range map1 {
		ctx = logging.WithLogger(ctx, log.With(zap.String("objectName", name)))

		if index2, ok := map2[name]; ok {
			wg.Add(1)

			index1.Check = true
			map1[name] = index1
			index2.Check = true
			map2[name] = index2

			compareJobSpecs(ctx, wg, channel, name, namespace, &jobs1.Items[index1.Index], &jobs2.Items[index2.Index])
		} else {
			log.With(zap.String("objectName", name)).Warn("job does not exist in 2nd cluster")
			flag = true
			channel <- flag
		}
	}

	wg.Wait()

	close(channel)

	for ch := range channel {
		if ch {
			flag = true
		}
	}

	for name, index := range map2 {
		if !index.Check {
			log.With(zap.String("objectName", name)).Warn("job does not exist in 1st cluster")
			flag = true
		}
	}

	return flag
}

func compareJobSpecs(ctx context.Context, wg *sync.WaitGroup, channel chan bool, name, namespace string, obj1, obj2 *v12.Job) {
	var (
		log  = logging.FromContext(ctx)
		flag bool
	)
	defer func() {
		wg.Done()
	}()

	log.Debugf("----- Start checking job/%s -----", name)

	if !metadata.IsMetadataDiffers(ctx, obj1.ObjectMeta, obj2.ObjectMeta) {
		channel <- true
		return
	}

	bDiff, err := compareJobSpecInternals(ctx, obj1.Spec, obj2.Spec)
	if err != nil || bDiff {
		log.Warnw(err.Error())
		flag = true
	}

	log.Debugf("----- End checking job/%s -----", name)
	channel <- flag
}

func compareJobSpecInternals(ctx context.Context, obj1, obj2 v12.JobSpec) (bool, error) {
	log := logging.FromContext(ctx)

	if obj1.BackoffLimit != nil && obj2.BackoffLimit != nil {
		if *obj1.BackoffLimit != *obj2.BackoffLimit {
			log.Warnw("Job backoff limit is different", zap.Int32("backoffLimit1", *obj1.BackoffLimit), zap.Int32("backoffLimit2", *obj2.BackoffLimit))
			return true, ErrorBackoffLimitDifferent
		}
	} else if obj1.BackoffLimit != nil || obj2.BackoffLimit != nil {
		return true, ErrorBackoffLimitDifferent
	}

	if obj1.Template.Spec.RestartPolicy != obj2.Template.Spec.RestartPolicy {
		log.Warnw("Job restartPolicy limit is different", zap.String("restartPolicy1", string(obj1.Template.Spec.RestartPolicy)), zap.String("restartPolicy2", string(obj1.Template.Spec.RestartPolicy)))
		return true, ErrorRestartPolicyDifferent
	}

	castJob1ForCompareContainers := types.InformationAboutObject{
		Template: obj1.Template,
		Selector: nil,
	}
	castJob2ForCompareContainers := types.InformationAboutObject{
		Template: obj2.Template,
		Selector: nil,
	}

	bDiff, err := pods.ComparePodSpecs(ctx, castJob1ForCompareContainers, castJob2ForCompareContainers)
	if err != nil || bDiff {
		return bDiff, err
	}

	return false, nil
}
