package jobs

import (
	"fmt"
	"k8s-cluster-comparator/internal/kubernetes/common"
	"k8s-cluster-comparator/internal/kubernetes/kv_maps"
	"k8s-cluster-comparator/internal/kubernetes/skipper"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sync"
)

func CompareCronJobs(clientSet1, clientSet2 kubernetes.Interface, namespace string, skipEntityList skipper.SkipEntitiesList) (bool, error) {
	var (
		isClustersDiffer bool
	)

	cronJobs1, err := clientSet1.BatchV1beta1().CronJobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("cannot obtain jobs list from 1st cluster: %w", err)
	}

	cronJobs2, err := clientSet2.BatchV1beta1().CronJobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("cannot obtain jobs list from 2nd cluster: %w", err)
	}

	mapJobs1, mapJobs2 := prepareCronJobsMaps(cronJobs1, cronJobs2, skipEntityList.GetByKind("cronJobs"))

	isClustersDiffer = setInformationAboutCronJobs(mapJobs1, mapJobs2, cronJobs1, cronJobs2)

	return isClustersDiffer, nil
}

func prepareCronJobsMaps(cronJobs1, cronJobs2 *v1beta1.CronJobList, skipEntities skipper.SkipComponentNames) (map[string]types.IsAlreadyComparedFlag, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	mapCronJobs1 := make(map[string]types.IsAlreadyComparedFlag)
	mapCronJobs2 := make(map[string]types.IsAlreadyComparedFlag)
	var indexCheck types.IsAlreadyComparedFlag

	for index, value := range cronJobs1.Items {
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("cronJob %s is skipped from comparison due to its name", value.Name)
			continue
		}
		indexCheck.Index = index
		mapCronJobs1[value.Name] = indexCheck

	}
	for index, value := range cronJobs2.Items {
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("cronJob %s is skipped from comparison due to its name", value.Name)
			continue
		}
		indexCheck.Index = index
		mapCronJobs2[value.Name] = indexCheck

	}
	return mapCronJobs1, mapCronJobs2
}

// setInformationAboutCronJobs set information about jobs
func setInformationAboutCronJobs(map1, map2 map[string]types.IsAlreadyComparedFlag, cronJobs1, cronJobs2 *v1beta1.CronJobList) bool {
	var (
		flag bool
	)

	if len(map1) != len(map2) {
		log.Infof("cronJobs counts are different")
		flag = true
	}

	wg := &sync.WaitGroup{}
	channel := make(chan bool, len(map1))

	for name, index1 := range map1 {
		if index2, ok := map2[name]; ok {
			wg.Add(1)

			index1.Check = true
			map1[name] = index1
			index2.Check = true
			map2[name] = index2

			compareCronJobSpecInternals(wg, channel, name, &cronJobs1.Items[index1.Index], &cronJobs2.Items[index2.Index])
		} else {
			log.Infof("cronJob '%s' does not exist in 2nd cluster", name)
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

			log.Infof("cronJob '%s' does not exist in 1st cluster", name)
			flag = true

		}
	}

	return flag
}

func compareCronJobSpecInternals(wg *sync.WaitGroup, channel chan bool, name string, cronJob1, cronJob2 *v1beta1.CronJob) {
	var (
		flag bool
	)
	defer func() {
		wg.Done()
	}()

	log.Debugf("----- Start checking cronJob: '%s' -----", name)

	if !kv_maps.AreKVMapsEqual(cronJob1.ObjectMeta.Labels, cronJob2.ObjectMeta.Labels, common.SkippedKubeLabels) {
		log.Infof("metadata of cronJob '%s' differs: different labels", cronJob1.Name)
		channel <- true
		return
	}

	if !kv_maps.AreKVMapsEqual(cronJob1.ObjectMeta.Labels, cronJob2.ObjectMeta.Labels, nil) {
		log.Infof("metadata of cronJob '%s' differs: different annotations", cronJob2.Name)
		channel <- true
		return
	}

	err := compareSpecInCronJobs(*cronJob1, *cronJob2)
	if err != nil {
		log.Infof("CronJob %s: %s", name, err.Error())
		flag = true
	}

	log.Debugf("----- End checking cronJob: '%s' -----", name)
	channel <- flag
}

func compareSpecInCronJobs(cronJob1, cronJob2 v1beta1.CronJob) error {

	if cronJob1.Spec.Schedule != cronJob2.Spec.Schedule {
		return fmt.Errorf("%w. CronJob name: %s. CronJob 1 - %s, cronJob2 - %s ", ErrorScheduleDifferent, cronJob1.Name, cronJob1.Spec.Schedule, cronJob2.Spec.Schedule)
	}

	err := compareSpecInJobs(cronJob1.Spec.JobTemplate.Spec, cronJob1.Spec.JobTemplate.Spec)
	if err != nil {
		return err
	}
	//castJob1ForCompareContainers := types.InformationAboutObject{
	//	Template: job1.Spec.Template,
	//	Selector: nil,
	//}
	//castJob2ForCompareContainers := types.InformationAboutObject{
	//	Template: job2.Spec.Template,
	//	Selector: nil,
	//}
	//
	//err := pod_controllers.CompareContainers(castJob1ForCompareContainers, castJob2ForCompareContainers, "default",  true, true, nil, nil)
	//if err != nil {
	//	return err
	//}

	return nil
}