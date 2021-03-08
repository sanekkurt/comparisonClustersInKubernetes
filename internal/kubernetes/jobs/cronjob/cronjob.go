package cronjob

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"k8s-cluster-comparator/internal/kubernetes/jobs/common"

	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/kubernetes/metadata"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"

	"sync"

	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	cronJobKind = "cronjob"
)

func cronJobsRetrieveBatchLimit(ctx context.Context) int64 {
	cfg := config.FromContext(ctx)

	if limit := cfg.Tasks.CronJobs.BatchSize; limit != 0 {
		return limit
	}

	if limit := cfg.Common.DefaultBatchSize; limit != 0 {
		return limit
	}

	return 25
}

func addItemsToCronJobList(ctx context.Context, clientSet kubernetes.Interface, namespace string, limit int64) (*v1beta1.CronJobList, error) {
	log := logging.FromContext(ctx)

	log.Debugf("addItemsToCronJobList started")
	defer log.Debugf("addItemsToCronJobList completed")

	var (
		batch    *v1beta1.CronJobList
		cronJobs = &v1beta1.CronJobList{
			Items: make([]v1beta1.CronJob, 0),
		}

		continueToken string

		err error
	)

	for {
		batch, err = clientSet.BatchV1beta1().CronJobs(namespace).List(metav1.ListOptions{
			Limit:    limit,
			Continue: continueToken,
		})
		if err != nil {
			return nil, err
		}

		log.Debugf("addItemsToCronJobList: %d objects received", len(batch.Items))

		cronJobs.Items = append(cronJobs.Items, batch.Items...)

		cronJobs.TypeMeta = batch.TypeMeta
		cronJobs.ListMeta = batch.ListMeta

		if batch.Continue == "" {
			break
		}

		continueToken = batch.Continue
	}

	cronJobs.Continue = ""

	return cronJobs, err
}

type CronJobsComparator struct {
}

func NewCronJobsComparator(ctx context.Context, namespace string) CronJobsComparator {
	return CronJobsComparator{}
}

func (cmp CronJobsComparator) Compare(ctx context.Context, namespace string) ([]types.KubeObjectsDifference, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("kind", cronJobKind))
		cfg = config.FromContext(ctx)
	)
	ctx = logging.WithLogger(ctx, log)

	if !cfg.Workloads.Enabled ||
		!cfg.Tasks.Enabled ||
		!cfg.Tasks.CronJobs.Enabled {
		log.Infof("'%s' kind skipped from comparison due to configuration", cronJobKind)
		return nil, nil
	}

	cronJobs1, err := addItemsToCronJobList(ctx, cfg.Connections.Cluster1.ClientSet, namespace, cronJobsRetrieveBatchLimit(ctx))
	if err != nil {
		return nil, fmt.Errorf("cannot obtain cronJobs list from 1st cluster: %w", err)
	}

	cronJobs2, err := addItemsToCronJobList(ctx, cfg.Connections.Cluster2.ClientSet, namespace, cronJobsRetrieveBatchLimit(ctx))
	if err != nil {
		return nil, fmt.Errorf("cannot obtain cronJobs list from 2st cluster: %w", err)
	}

	mapJobs1, mapJobs2 := prepareCronJobsMaps(ctx, cronJobs1, cronJobs2)

	_ = setInformationAboutCronJobs(ctx, mapJobs1, mapJobs2, cronJobs1, cronJobs2, namespace)

	return nil, nil
}

func prepareCronJobsMaps(ctx context.Context, cronJobs1, cronJobs2 *v1beta1.CronJobList) (map[string]types.IsAlreadyComparedFlag, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)

		mapCronJobs1 = make(map[string]types.IsAlreadyComparedFlag)
		mapCronJobs2 = make(map[string]types.IsAlreadyComparedFlag)

		indexCheck types.IsAlreadyComparedFlag
	)

	for index, value := range cronJobs1.Items {
		if cfg.ExcludesIncludes.IsSkippedEntity(cronJobKind, value.Name) {
			log.With(zap.String("name", value.Name)).Debugf("cronjob/%s is skipped from comparison", value.Name)
			continue
		}

		indexCheck.Index = index
		mapCronJobs1[value.Name] = indexCheck

	}

	for index, value := range cronJobs2.Items {
		if cfg.ExcludesIncludes.IsSkippedEntity(cronJobKind, value.Name) {
			log.With(zap.String("name", value.Name)).Debugf("cronjob/%s is skipped from comparison", value.Name)
			continue
		}

		indexCheck.Index = index
		mapCronJobs2[value.Name] = indexCheck

	}
	return mapCronJobs1, mapCronJobs2
}

// setInformationAboutCronJobs set information about jobs
func setInformationAboutCronJobs(ctx context.Context, map1, map2 map[string]types.IsAlreadyComparedFlag, cronJobs1, cronJobs2 *v1beta1.CronJobList, namespace string) bool {
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

			compareCronJobSpecs(ctx, wg, channel, name, namespace, &cronJobs1.Items[index1.Index], &cronJobs2.Items[index2.Index])
		} else {
			log.With(zap.String("objectName", name)).Warn("cronjob does not exist in 2nd cluster")
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
			log.With(zap.String("objectName", name)).Warn("cronjob does not exist in 1st cluster")

			flag = true
		}
	}

	return flag
}

func compareCronJobSpecs(ctx context.Context, wg *sync.WaitGroup, channel chan bool, name, namespace string, obj1, obj2 *v1beta1.CronJob) {
	var (
		log = logging.FromContext(ctx)

		flag bool
	)
	defer func() {
		wg.Done()
	}()

	log.Debugf("----- Start checking cronjob/%s -----", name)

	if metadata.IsMetadataDiffers(ctx, obj1.ObjectMeta, obj2.ObjectMeta) {
		channel <- true
		return
	}

	bDiff, err := compareCronJobSpecInternals(ctx, *obj1, *obj2)
	if err != nil || bDiff {
		log.Warnw(err.Error())
		flag = true
	}

	log.Debugf("----- End checking cronjob/%s -----", name)
	channel <- flag
}

func compareCronJobSpecInternals(ctx context.Context, obj1, obj2 v1beta1.CronJob) (bool, error) {
	log := logging.FromContext(ctx)

	if obj1.Spec.Schedule != obj2.Spec.Schedule {
		log.Warnw("CronJob schedule is different", zap.String("schedule1", obj1.Spec.Schedule), zap.String("schedule2", obj2.Spec.Schedule))
		return true, common.ErrorScheduleDifferent
	}

	bDiff, err := common.CompareJobSpecInternals(ctx, obj1.Spec.JobTemplate.Spec, obj1.Spec.JobTemplate.Spec)

	return bDiff, err
}
