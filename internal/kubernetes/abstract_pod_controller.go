package kubernetes

import (
	"k8s-cluster-comparator/internal/logging"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"sync"
)

// AbstractPodController is an generalized abstraction above deployment/statefulset/daemonset/etc kubernetes pod controllers
type AbstractPodController struct {
	Name string

	Metadata AbstractObjectMetadata

	Labels      map[string]string
	Annotations map[string]string

	Replicas *int32

	PodLabelSelector *v12.LabelSelector
	PodTemplateSpec  v1.PodTemplateSpec
}

// ObjectKindWrapper is a wrapper function that transforms object kind name to a canonical form
func ObjectKindWrapper(kind string) string {
	return strings.ToLower(kind)
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
		wg.Add(1)

		go func(wg *sync.WaitGroup, channel chan bool, name string, index1 IsAlreadyComparedFlag, map1, map2 map[string]IsAlreadyComparedFlag) {
			defer func() {
				wg.Done()
			}()

			if index2, ok := map2[name]; ok {
				index1.Check = true
				map1[name] = index1

				index2.Check = true
				map2[name] = index2

				apc1 := apc1List[index1.Index]
				apc2 := apc2List[index2.Index]
				kind := ObjectKindWrapper(apc1.Metadata.Type.Kind)

				if _, err := CompareAbstractObjectMetadata(apc1.Metadata, apc2.Metadata); err != nil {
					logging.Log.Infof("metadata compare error: %s", err.Error())

					channel <- true

					return
				}

				logging.Log.Debugf("----- Start checking '%s:%s' pod controller spec -----", kind, apc1.Name)

				if apc1.Replicas != nil || apc2.Replicas != nil {
					if *apc1.Replicas != *apc2.Replicas {
						logging.Log.Infof("%s:%s: number of replicas is different: %d and %d", kind, apc1.Replicas, apc2.Replicas)
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
			} else {
				logging.Log.Infof("%s %s presents in 1st cluster but absents in 2nd one", name)
				flag = true
			}
			channel <- flag
		}(wg, channel, name, index1, map1, map2)

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
			logging.Log.Infof("%s %s presents in 2nd cluster but absents in 1st one", name)
			flag = true
		}
	}

	return flag
}
