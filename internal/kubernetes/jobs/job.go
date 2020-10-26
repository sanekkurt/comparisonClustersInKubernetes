package jobs

import (
	"fmt"
	"k8s-cluster-comparator/internal/kubernetes/common"
	"k8s-cluster-comparator/internal/kubernetes/kv_maps"
	"k8s-cluster-comparator/internal/kubernetes/pod_controllers"
	"sync"

	"k8s-cluster-comparator/internal/kubernetes/skipper"
	"k8s-cluster-comparator/internal/kubernetes/types"
	v12 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CompareJobs compare jobs in different clusters
func CompareJobs(clientSet1, clientSet2 kubernetes.Interface, namespace string, skipEntityList skipper.SkipEntitiesList) (bool, error) {
	var (
		isClustersDiffer bool
	)

	jobs1, err := clientSet1.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("cannot obtain jobs list from 1st cluster: %w", err)
	}

	jobs2, err := clientSet2.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("cannot obtain jobs list from 2nd cluster: %w", err)
	}

	mapJobs1, mapJobs2 := prepareJobsMaps(jobs1, jobs2, skipEntityList.GetByKind("jobs"))

	isClustersDiffer = setInformationAboutJobs(mapJobs1, mapJobs2, jobs1, jobs2, namespace)

	return isClustersDiffer, nil
}

// prepareJobsMaps add value secrets in map
func prepareJobsMaps(jobs1, jobs2 *v12.JobList, skipEntities skipper.SkipComponentNames) (map[string]types.IsAlreadyComparedFlag, map[string]types.IsAlreadyComparedFlag) { //nolint:gocritic,unused
	mapJobs1 := make(map[string]types.IsAlreadyComparedFlag)
	mapJobs2 := make(map[string]types.IsAlreadyComparedFlag)
	var indexCheck types.IsAlreadyComparedFlag

	for index, value := range jobs1.Items {
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("job %s is skipped from comparison due to its name", value.Name)
			continue
		}
		indexCheck.Index = index
		mapJobs1[value.Name] = indexCheck

	}
	for index, value := range jobs2.Items {
		if skipEntities.IsSkippedEntity(value.Name) {
			log.Debugf("job %s is skipped from comparison due to its name", value.Name)
			continue
		}
		indexCheck.Index = index
		mapJobs2[value.Name] = indexCheck

	}
	return mapJobs1, mapJobs2
}

// setInformationAboutJobs set information about jobs
func setInformationAboutJobs(map1, map2 map[string]types.IsAlreadyComparedFlag, jobs1, jobs2 *v12.JobList, namespace string) bool {
	var (
		flag bool
	)

	if len(map1) != len(map2) {
		log.Infof("jobs counts are different")
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

			compareJobSpecInternals(wg, channel, name, namespace, &jobs1.Items[index1.Index], &jobs2.Items[index2.Index])
		} else {
			log.Infof("job '%s' does not exist in 2nd cluster", name)
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

			log.Infof("job '%s' does not exist in 1st cluster", name)
			flag = true

		}
	}

	return flag
}

func compareJobSpecInternals(wg *sync.WaitGroup, channel chan bool, name, namespace string, job1, job2 *v12.Job) {
	var (
		flag bool
	)
	defer func() {
		wg.Done()
	}()

	log.Debugf("----- Start checking job: '%s' -----", name)

	if !kv_maps.AreKVMapsEqual(job1.ObjectMeta.Labels, job2.ObjectMeta.Labels, common.SkippedKubeLabels) {
		log.Infof("metadata of job '%s' differs: different labels", job1.Name)
		channel <- true
		return
	}

	if !kv_maps.AreKVMapsEqual(job1.ObjectMeta.Labels, job2.ObjectMeta.Labels, nil) {
		log.Infof("metadata of job '%s' differs: different annotations", job2.Name)
		channel <- true
		return
	}

	err := compareSpecInJobs(job1.Spec, job2.Spec, namespace)
	if err != nil {
		log.Infof("Job %s: %s", name, err.Error())
		flag = true
	}

	log.Debugf("----- End checking job: '%s' -----", name)
	channel <- flag
}

func compareSpecInJobs(job1, job2 v12.JobSpec, namespace string) error {

	if job1.BackoffLimit != nil && job2.BackoffLimit != nil {
		if *job1.BackoffLimit != *job2.BackoffLimit {
			return fmt.Errorf("%w. Job 1 - %d, Job 2 - %d", ErrorBackoffLimitDifferent, &job1.BackoffLimit, &job2.BackoffLimit )
		}
	} else if job1.BackoffLimit != nil || job2.BackoffLimit != nil {
		return ErrorBackoffLimitDifferent
	}

	if job1.Template.Spec.RestartPolicy != job2.Template.Spec.RestartPolicy {
		return fmt.Errorf("%w. Job 1 - %s, Job 2 - %s", ErrorRestartPolicyDifferent, job1.Template.Spec.RestartPolicy, job2.Template.Spec.RestartPolicy)
	}

	castJob1ForCompareContainers := types.InformationAboutObject{
		Template: job1.Template,
		Selector: nil,
	}
	castJob2ForCompareContainers := types.InformationAboutObject{
		Template: job2.Template,
		Selector: nil,
	}

	err := pod_controllers.CompareContainers(castJob1ForCompareContainers, castJob2ForCompareContainers, namespace,  true, true, nil, nil)
	if err != nil {
		return err
	}

	return nil
}