package job

import (
	"context"
	"fmt"
	"k8s-cluster-comparator/internal/consts"
	kubectx "k8s-cluster-comparator/internal/kubernetes/context"
	"k8s.io/apimachinery/pkg/api/errors"

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
	objectKind = "job"
)

type Comparator struct {
	Kind      string
	Namespace string
	BatchSize int64
}

func NewJobsComparator(ctx context.Context, namespace string) *Comparator {
	return &Comparator{
		Kind:      objectKind,
		Namespace: namespace,
		BatchSize: getBatchLimit(ctx),
	}
}

func (cmp *Comparator) fieldSelectorProvider(ctx context.Context) string {
	return ""
}

func (cmp *Comparator) labelSelectorProvider(ctx context.Context) string {
	return ""
}

func (cmp *Comparator) collectIncludedFromCluster(ctx context.Context) (map[string]batchv1.Job, error) {
	var (
		log       = logging.FromContext(ctx)
		cfg       = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		objects = make(map[string]batchv1.Job)
	)

	log.Debugf("%T: collectIncludedFromCluster started", cmp)
	defer log.Debugf("%T: collectIncludedFromCluster completed", cmp)

	for name := range cfg.ExcludesIncludes.NameBasedSkip {
		obj, err := clientSet.BatchV1().Jobs(cmp.Namespace).Get(string(name), metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				log.With(zap.String("objectName", string(name))).Warnf("%s/%s not found in cluster", cmp.Kind, name)
				continue
			}
			return nil, err
		}
		objects[obj.Name] = *obj
	}

	for name := range cfg.ExcludesIncludes.FullResourceNamesSkip[types.ObjectKind(cmp.Kind)] {
		obj, err := clientSet.BatchV1().Jobs(cmp.Namespace).Get(string(name), metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				log.With(zap.String("objectName", string(name))).Warnf("%s/%s not found in cluster", cmp.Kind, name)
				continue
			}
			return nil, err
		}
		objects[obj.Name] = *obj
	}

	return objects, nil
}

func (cmp *Comparator) collectFromClusterWithoutExcludes(ctx context.Context) (map[string]batchv1.Job, error) {
	var (
		log       = logging.FromContext(ctx)
		cfg       = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		batch   *batchv1.JobList
		objects = make(map[string]batchv1.Job)

		continueToken string

		err error
	)

	log.Debugf("%T: collectFromClusterWithoutExcludes started", cmp)
	defer log.Debugf("%T: collectFromClusterWithoutExcludes completed", cmp)

forOuterLoop:
	for {
		select {
		case <-ctx.Done():
			return nil, context.Canceled
		default:
			batch, err = clientSet.BatchV1().Jobs(cmp.Namespace).List(metav1.ListOptions{
				Limit:         cmp.BatchSize,
				FieldSelector: cmp.fieldSelectorProvider(ctx),
				LabelSelector: cmp.labelSelectorProvider(ctx),
				Continue:      continueToken,
			})
			if err != nil {
				return nil, err
			}

			log.Debugf("%d %s retrieved", len(batch.Items), cmp.Kind)

		forInnerLoop:
			for _, obj := range batch.Items {
				if _, ok := objects[obj.Name]; ok {
					log.With("objectName", obj.Name).Warnf("%s/%s already present in comparison list", cmp.Kind, obj.Name)
				}

				if cfg.ExcludesIncludes.IsSkippedEntity(cmp.Kind, obj.Name) {
					log.With(zap.String("objectName", obj.Name)).Debugf("%s/%s is skipped from comparison", cmp.Kind, obj.Name)
					continue forInnerLoop
				}

				objects[obj.Name] = obj
			}

			if batch.Continue == "" {
				break forOuterLoop
			}

			continueToken = batch.Continue
		}
	}

	return objects, nil
}

func (cmp *Comparator) collectFromCluster(ctx context.Context) (map[string]batchv1.Job, error) {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)
	)

	log.Debugf("%T: collectFromCluster started", cmp)
	defer log.Debugf("%T: collectFromCluster completed", cmp)

	if cfg.Common.WorkMode == consts.EverythingButNotExcludesWorkMode {
		return cmp.collectFromClusterWithoutExcludes(ctx)
	} else {
		return cmp.collectIncludedFromCluster(ctx)
	}
}

func (cmp *Comparator) Compare(ctx context.Context) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("kind", cmp.Kind))
		cfg = config.FromContext(ctx)
	)
	ctx = logging.WithLogger(ctx, log)

	if !cfg.Tasks.Enabled ||
		!cfg.Tasks.Jobs.Enabled {
		log.Infof("'%s' kind skipped from comparison due to configuration", cmp.Kind)
		return nil, nil
	}

	objects, err := cmp.collect(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve objects for comparision: %w", err)
	}

	diff := cmp.compare(ctx, objects[0], objects[1])

	return diff, nil
}

func (cmp *Comparator) collect(ctx context.Context) ([]map[string]batchv1.Job, error) {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)

		objects = make([]map[string]batchv1.Job, 2, 2)
		wg      = &sync.WaitGroup{}

		err error
	)

	wg.Add(2)

	for idx, clientSet := range []kubernetes.Interface{
		cfg.Connections.Cluster1.ClientSet,
		cfg.Connections.Cluster2.ClientSet,
	} {
		go func(idx int, clientSet kubernetes.Interface) {
			defer wg.Done()

			objects[idx], err = cmp.collectFromCluster(kubectx.WithClientSet(ctx, clientSet))
			if err != nil {
				log.Fatalf("cannot obtain %ss from cluster #%d: %s", cmp.Kind, idx+1, err.Error())
			}
		}(idx, clientSet)
	}

	wg.Wait()

	return objects, nil
}

func (cmp *Comparator) compare(ctx context.Context, map1, map2 map[string]batchv1.Job) []types.KubeObjectsDifference {
	var (
		log = logging.FromContext(ctx)

		diffs = make([]types.KubeObjectsDifference, 0)
	)

	if len(map1) != len(map2) {
		log.Warnw("object counts are different", zap.Int("objectsCount1st", len(map1)), zap.Int("objectsCount2nd", len(map2)))
	}

	for name, obj1 := range map1 {
		ctx = logging.WithLogger(ctx, log.With(zap.String("objectName", name)))

		select {
		case <-ctx.Done():
			log.Warnw(context.Canceled.Error())
			return nil
		default:
			if obj2, ok := map2[name]; ok {
				diff := compareJobSpecs(ctx, name, &obj1, &obj2)

				diffs = append(diffs, diff...)

				delete(map1, name)
				delete(map2, name)
			} else {
				log.With(zap.String("objectName", name)).Warnf("%s does not exist in 2nd cluster", cmp.Kind)
			}
		}
	}

	for name, _ := range map2 {
		log.With(zap.String("objectName", name)).Warnf("%s does not exist in 1st cluster", cmp.Kind)
	}

	return diffs
}

func compareJobSpecs(ctx context.Context, name string, obj1, obj2 *batchv1.Job) []types.KubeObjectsDifference {
	var (
		log = logging.FromContext(ctx)
	)

	ctx = logging.WithLogger(ctx, log)

	log.Debugf("job/%s compare started", name)
	defer func() {
		log.Debugf("job/%s compare completed", name)
	}()

	metadata.IsMetadataDiffers(ctx, obj1.ObjectMeta, obj2.ObjectMeta)

	err := common.CompareJobSpecInternals(ctx, obj1.Spec, obj2.Spec)
	if err != nil {
		return nil //log.Warnw(err.Error())
	}
	return nil

}

//const (
//	jobKind = "job"
//)
//
//func jobsRetrieveBatchLimit(ctx context.Context) int64 {
//	cfg := config.FromContext(ctx)
//
//	if limit := cfg.Tasks.Jobs.BatchSize; limit != 0 {
//		return limit
//	}
//
//	if limit := cfg.Common.DefaultBatchSize; limit != 0 {
//		return limit
//	}
//
//	return 25
//}
//
//func addItemsToJobList(ctx context.Context, clientSet kubernetes.Interface, namespace string, limit int64) (*batchv1.JobList, error) {
//	log := logging.FromContext(ctx)
//
//	log.Debugf("addItemsToJobList started")
//	defer log.Debugf("addItemsToJobList completed")
//
//	var (
//		batch *batchv1.JobList
//		jobs  = &batchv1.JobList{
//			Items: make([]batchv1.Job, 0),
//		}
//
//		continueToken string
//
//		err error
//	)
//
//	for {
//		batch, err = clientSet.BatchV1().Jobs(namespace).List(metav1.ListOptions{
//			Limit:    limit,
//			Continue: continueToken,
//		})
//		if err != nil {
//			return nil, err
//		}
//
//		log.Debugf("addItemsToJobList: %d objects received", len(batch.Items))
//
//		jobs.Items = append(jobs.Items, batch.Items...)
//
//		jobs.TypeMeta = batch.TypeMeta
//		jobs.ListMeta = batch.ListMeta
//
//		if batch.Continue == "" {
//			break
//		}
//
//		continueToken = batch.Continue
//	}
//
//	jobs.Continue = ""
//
//	return jobs, err
//}
//
//type JobsComparator struct {
//}
//
//func NewJobsComparator(ctx context.Context, namespace string) JobsComparator {
//	return JobsComparator{}
//}
//
//// Compare compare Jobs in different clusters
//func (cmp JobsComparator) Compare(ctx context.Context, namespace string) ([]types.KubeObjectsDifference, error) {
//	var (
//		log = logging.FromContext(ctx).With(zap.String("kind", jobKind))
//		cfg = config.FromContext(ctx)
//	)
//	ctx = logging.WithLogger(ctx, log)
//
//	if !cfg.Workloads.Enabled ||
//		!cfg.Tasks.Enabled ||
//		!cfg.Tasks.Jobs.Enabled {
//		log.Infof("'%s' kind skipped from comparison due to configuration", jobKind)
//		return nil, nil
//	}
//
//	jobs1, err := addItemsToJobList(ctx, cfg.Connections.Cluster1.ClientSet, namespace, jobsRetrieveBatchLimit(ctx))
//	if err != nil {
//		return nil, fmt.Errorf("cannot obtain jobs list from 1st cluster: %w", err)
//	}
//
//	jobs2, err := addItemsToJobList(ctx, cfg.Connections.Cluster2.ClientSet, namespace, jobsRetrieveBatchLimit(ctx))
//	if err != nil {
//		return nil, fmt.Errorf("cannot obtain jobs list from 2st cluster: %w", err)
//	}
//
//	mapJobs1, mapJobs2 := prepareJobsMaps(ctx, jobs1, jobs2)
//
//	_ = setInformationAboutJobs(ctx, mapJobs1, mapJobs2, jobs1, jobs2, namespace)
//
//	return nil, nil
//}
//
//// prepareJobsMaps add value secrets in map
//func prepareJobsMaps(ctx context.Context, jobs1, jobs2 *batchv1.JobList) (map[string]types.IsAlreadyComparedFlag, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
//	var (
//		log = logging.FromContext(ctx)
//		cfg = config.FromContext(ctx)
//
//		mapJobs1 = make(map[string]types.IsAlreadyComparedFlag)
//		mapJobs2 = make(map[string]types.IsAlreadyComparedFlag)
//
//		indexCheck types.IsAlreadyComparedFlag
//	)
//
//OUTER1:
//	for index, value := range jobs1.Items {
//		if cfg.ExcludesIncludes.IsSkippedEntity(jobKind, value.Name) {
//			log.With(zap.String("name", value.Name)).Debugf("job/%s is skipped from comparison", value.Name)
//			continue
//		}
//
//		if value.OwnerReferences != nil {
//			for _, owner := range value.OwnerReferences {
//				if owner.Kind == "CronJob" {
//					log.Debugf("job/%s is skipped from comparison because it is owned by CronJob", value.Name)
//					continue OUTER1
//				}
//			}
//		}
//
//		indexCheck.Index = index
//		mapJobs1[value.Name] = indexCheck
//
//	}
//
//OUTER2:
//	for index, value := range jobs2.Items {
//		if cfg.ExcludesIncludes.IsSkippedEntity(jobKind, value.Name) {
//			log.With(zap.String("name", value.Name)).Debugf("job/%s is skipped from comparison", value.Name)
//			continue
//		}
//
//		if value.OwnerReferences != nil {
//			for _, owner := range value.OwnerReferences {
//				if owner.Kind == "CronJob" {
//					log.Debugf("job/%s is skipped from comparison because it is owned by CronJob", value.Name)
//					continue OUTER2
//				}
//			}
//
//		}
//
//		indexCheck.Index = index
//		mapJobs2[value.Name] = indexCheck
//
//	}
//
//	return mapJobs1, mapJobs2
//}
//
//// setInformationAboutJobs set information about jobs
//func setInformationAboutJobs(ctx context.Context, map1, map2 map[string]types.IsAlreadyComparedFlag, jobs1, jobs2 *batchv1.JobList, namespace string) bool {
//	var (
//		log  = logging.FromContext(ctx)
//		flag bool
//	)
//
//	if len(map1) != len(map2) {
//		log.Warnw("object counts are different", zap.Int("objectsCount1st", len(map1)), zap.Int("objectsCount2nd", len(map2)))
//		flag = true
//	}
//
//	wg := &sync.WaitGroup{}
//	channel := make(chan bool, len(map1))
//
//	for name, index1 := range map1 {
//		ctx = logging.WithLogger(ctx, log.With(zap.String("objectName", name)))
//
//		if index2, ok := map2[name]; ok {
//			wg.Add(1)
//
//			index1.Check = true
//			map1[name] = index1
//			index2.Check = true
//			map2[name] = index2
//
//			compareJobSpecs(ctx, wg, channel, name, namespace, &jobs1.Items[index1.Index], &jobs2.Items[index2.Index])
//		} else {
//			log.With(zap.String("objectName", name)).Warn("job does not exist in 2nd cluster")
//			flag = true
//			channel <- flag
//		}
//	}
//
//	wg.Wait()
//
//	close(channel)
//
//	for ch := range channel {
//		if ch {
//			flag = true
//		}
//	}
//
//	for name, index := range map2 {
//		if !index.Check {
//			log.With(zap.String("objectName", name)).Warn("job does not exist in 1st cluster")
//			flag = true
//		}
//	}
//
//	return flag
//}
//
//func compareJobSpecs(ctx context.Context, wg *sync.WaitGroup, channel chan bool, name, namespace string, obj1, obj2 *batchv1.Job) {
//	var (
//		log  = logging.FromContext(ctx)
//		flag bool
//	)
//	defer func() {
//		wg.Done()
//	}()
//
//	log.Debugf("----- Start checking job/%s -----", name)
//
//	if !metadata.IsMetadataDiffers(ctx, obj1.ObjectMeta, obj2.ObjectMeta) {
//		channel <- true
//		return
//	}
//
//	bDiff, err := common.CompareJobSpecInternals(ctx, obj1.Spec, obj2.Spec)
//	if err != nil || bDiff {
//		log.Warnw(err.Error())
//		flag = true
//	}
//
//	log.Debugf("----- End checking job/%s -----", name)
//	channel <- flag
//}
