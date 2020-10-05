package pod_controllers

import (
	"strings"
	"sync"

	"k8s-cluster-comparator/internal/logging"
)

// ObjectKindWrapper is a wrapper function that transforms object kind name to a canonical form
func ObjectKindWrapper(kind string) string {
	return strings.ToLower(kind)
}

func comparePodControllerSpecInternals(wg *sync.WaitGroup, channel chan bool, name, namespace string, apc1, apc2 *AbstractPodController) {
	var (
		flag bool
	)

	defer func() {
		wg.Done()
	}()

	kind := ObjectKindWrapper(apc1.Metadata.Type.Kind)

	logging.Log.Debugf("----- Start checking '%s:%s' pod controller spec -----", kind, apc1.Name)

	if apc1.Replicas != nil || apc2.Replicas != nil {
		if *apc1.Replicas != *apc2.Replicas {
			logging.Log.Infof("%s:%s: number of replicas is different: %d and %d", kind, name, *apc1.Replicas, *apc2.Replicas)
			flag = true
		}
	}
	if (apc1.Replicas != nil && apc2.Replicas == nil) || (apc2.Replicas != nil && apc1.Replicas == nil) {
		logging.Log.Infof("%s:%s: strange replicas specification difference: %#v and %#v", kind, apc1.Replicas, apc2.Replicas)
		flag = true
	}

	// fill in the information that will be used for comparison
	object1 := InformationAboutObject{
		Template: apc1.PodTemplateSpec,
		Selector: apc1.PodLabelSelector,
	}
	object2 := InformationAboutObject{
		Template: apc2.PodTemplateSpec,
		Selector: apc2.PodLabelSelector,
	}

	err := CompareContainers(object1, object2, namespace, Client1, Client2)
	if err != nil {
		logging.Log.Infof("%s %s: %s", kind, name, err.Error())
		flag = true
	}

	logging.Log.Debugf("----- End checking %s: '%s' -----", kind, name)

	channel <- flag
}

// ComparePodControllerSpecs compares abstracted pod controller specifications in two k8s clusters
func ComparePodControllerSpecs(map1, map2 map[string]IsAlreadyComparedFlag, apc1List, apc2List []AbstractPodController, namespace string) bool {
	var (
		flag bool

		wg      = &sync.WaitGroup{}
		channel = make(chan bool, len(map1))
	)

	if len(map1) != len(map2) {
		logging.Log.Infof("controllers count are different")
		flag = true
	}

	for name, index1 := range map1 {
		if index2, ok := map2[name]; ok {
			wg.Add(1)

			index1.Check = true
			map1[name] = index1

			index2.Check = true
			map2[name] = index2

			apc1 := apc1List[index1.Index]
			apc2 := apc2List[index2.Index]

			go comparePodControllerSpecInternals(wg, channel, name, namespace, &apc1, &apc2)
		} else {
			logging.Log.Infof("%s %s presents in 1st cluster but absents in 2nd one", apc1List[index1.Index].Metadata.Type.Kind, name)
			flag = true
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
			logging.Log.Infof("%s %s presents in 2nd cluster but absents in 1st one", apc2List[index.Index].Metadata.Type.Kind, name)
			flag = true
		}
	}

	return flag
}
