package pod_controllers

import (
	"sync"

	"k8s.io/client-go/kubernetes"

	"k8s-cluster-comparator/internal/kubernetes/types"
)

type clusterCompareTask struct {
	Client                   kubernetes.Interface
	APCList                  []AbstractPodController
	IsAlreadyCheckedFlagsMap map[string]types.IsAlreadyComparedFlag
}

// comparePodControllerSpecs compares abstracted pod controller specifications in two k8s clusters
func comparePodControllerSpecs(c1, c2 *clusterCompareTask, namespace string) bool {
	var (
		flag bool

		wg      = &sync.WaitGroup{}
		channel = make(chan bool, len(c1.IsAlreadyCheckedFlagsMap))
	)

	if len(c1.IsAlreadyCheckedFlagsMap) != len(c2.IsAlreadyCheckedFlagsMap) {
		log.Infof("controllers count are different")
		flag = true
	}

	for name, index1 := range c1.IsAlreadyCheckedFlagsMap {
		if index2, ok := c2.IsAlreadyCheckedFlagsMap[name]; ok {
			wg.Add(1)

			index1.Check = true
			c1.IsAlreadyCheckedFlagsMap[name] = index1

			index2.Check = true
			c2.IsAlreadyCheckedFlagsMap[name] = index2

			apc1 := c1.APCList[index1.Index]
			apc2 := c2.APCList[index2.Index]

			go comparePodControllerSpecInternals(wg, channel, name, namespace, c1.Client, c2.Client, &apc1, &apc2)
		} else {
			log.Infof("%s %s presents in 1st cluster but absents in 2nd one", c1.APCList[index1.Index].Metadata.Type.Kind, name)
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

	for name, index := range c2.IsAlreadyCheckedFlagsMap {
		if !index.Check {
			log.Infof("%s %s presents in 2nd cluster but absents in 1st one", c2.APCList[index.Index].Metadata.Type.Kind, name)
			flag = true
		}
	}

	return flag
}

func comparePodControllerSpecInternals(wg *sync.WaitGroup, channel chan bool, name, namespace string, c1, c2 kubernetes.Interface, apc1, apc2 *AbstractPodController) {
	var (
		flag bool
	)

	defer func() {
		wg.Done()
	}()

	kind := types.ObjectKindWrapper(apc1.Metadata.Type.Kind)

	log.Debugf("----- Start checking '%s:%s' pod controller spec -----", kind, apc1.Name)

	if apc1.Replicas != nil || apc2.Replicas != nil {
		if *apc1.Replicas != *apc2.Replicas {
			log.Infof("%s:%s: number of replicas is different: %d and %d", kind, name, *apc1.Replicas, *apc2.Replicas)
			flag = true
		}
	}
	if (apc1.Replicas != nil && apc2.Replicas == nil) || (apc2.Replicas != nil && apc1.Replicas == nil) {
		log.Infof("%s:%s: strange replicas specification difference: %#v and %#v", kind, apc1.Replicas, apc2.Replicas)
		flag = true
	}

	// fill in the information that will be used for comparison
	object1 := types.InformationAboutObject{
		Template: apc1.PodTemplateSpec,
		Selector: apc1.PodLabelSelector,
	}
	object2 := types.InformationAboutObject{
		Template: apc2.PodTemplateSpec,
		Selector: apc2.PodLabelSelector,
	}

	err := CompareContainers(object1, object2, namespace, c1, c2)
	if err != nil {
		log.Infof("%s %s: %s", kind, name, err.Error())
		flag = true
	}

	log.Debugf("----- End checking %s: '%s' -----", kind, name)

	channel <- flag
}
