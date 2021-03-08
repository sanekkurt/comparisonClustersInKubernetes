package job

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/jobs/common"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/metadata"
	"k8s-cluster-comparator/internal/logging"

	"sync"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/kubernetes/types"
)

const (
	jobKind = "job"
)

func jobsRetrieveBatchLimit(ctx context.Context) int64 {
	cfg := config.FromContext(ctx)

	if limit := cfg.Tasks.Jobs.BatchSize; limit != 0 {
		return limit
	}

	if limit := cfg.Common.DefaultBatchSize; limit != 0 {
		return limit
	}

	return 25
}

func addItemsToJobList(ctx context.Context, clientSet kubernetes.Interface, namespace string, limit int64) (*batchv1.JobList, error) {
	log := logging.FromContext(ctx)

	log.Debugf("addItemsToJobList started")
	defer log.Debugf("addItemsToJobList completed")

	var (
		batch *batchv1.JobList
		jobs  = &batchv1.JobList{
			Items: make([]batchv1.Job, 0),
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

	if !cfg.Workloads.Enabled ||
		!cfg.Tasks.Enabled ||
		!cfg.Tasks.Jobs.Enabled {
		log.Infof("'%s' kind skipped from comparison due to configuration", jobKind)
		return nil, nil
	}

	jobs1, err := addItemsToJobList(ctx, cfg.Connections.Cluster1.ClientSet, namespace, jobsRetrieveBatchLimit(ctx))
	if err != nil {
		return nil, fmt.Errorf("cannot obtain jobs list from 1st cluster: %w", err)
	}

	jobs2, err := addItemsToJobList(ctx, cfg.Connections.Cluster2.ClientSet, namespace, jobsRetrieveBatchLimit(ctx))
	if err != nil {
		return nil, fmt.Errorf("cannot obtain jobs list from 2st cluster: %w", err)
	}

	mapJobs1, mapJobs2 := prepareJobsMaps(ctx, jobs1, jobs2)

	_ = setInformationAboutJobs(ctx, mapJobs1, mapJobs2, jobs1, jobs2, namespace)

	return nil, nil
}

// prepareJobsMaps add value secrets in map
func prepareJobsMaps(ctx context.Context, jobs1, jobs2 *batchv1.JobList) (map[string]types.IsAlreadyComparedFlag, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)

		mapJobs1 = make(map[string]types.IsAlreadyComparedFlag)
		mapJobs2 = make(map[string]types.IsAlreadyComparedFlag)

		indexCheck types.IsAlreadyComparedFlag
	)

OUTER1:
	for index, value := range jobs1.Items {
		if cfg.ExcludesIncludes.IsSkippedEntity(jobKind, value.Name) {
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
		if cfg.ExcludesIncludes.IsSkippedEntity(jobKind, value.Name) {
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
func setInformationAboutJobs(ctx context.Context, map1, map2 map[string]types.IsAlreadyComparedFlag, jobs1, jobs2 *batchv1.JobList, namespace string) bool {
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

func compareJobSpecs(ctx context.Context, wg *sync.WaitGroup, channel chan bool, name, namespace string, obj1, obj2 *batchv1.Job) {
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

	bDiff, err := common.CompareJobSpecInternals(ctx, obj1.Spec, obj2.Spec)
	if err != nil || bDiff {
		log.Warnw(err.Error())
		flag = true
	}

	log.Debugf("----- End checking job/%s -----", name)
	channel <- flag
}
