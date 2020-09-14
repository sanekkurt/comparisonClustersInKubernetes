package main

import (
	"errors"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)
 //var ()
var labelForDeployment1 map[string]string
var replicasForDeployment1 int32
var pointerReplicasForDeployment1 *int32
var selectorForDeployment1 metav1.LabelSelector
var pointerSelectorForDeployment1 *metav1.LabelSelector

var objectInformation1 InformationAboutObject
var objectInformation2 InformationAboutObject

/*var labelForDeployment2 map[string]string
var replicasForDeployment2 int32
var pointerReplicasForDeployment2 *int32
var selectorForDeployment2 metav1.LabelSelector
var pointerSelectorForDeployment2 *metav1.LabelSelector

var labelForDeployment3 map[string]string
var replicasForDeployment3 int32
var pointerReplicasForDeployment3 *int32
var selectorForDeployment3 metav1.LabelSelector
var pointerSelectorForDeployment3 *metav1.LabelSelector*/

func initEnvironment() {
	labelForDeployment1 = make(map[string]string)
	labelForDeployment1["run"] = "test-depl1"
	replicasForDeployment1 = 2
	pointerReplicasForDeployment1 = &replicasForDeployment1
	selectorForDeployment1 = metav1.LabelSelector{
		MatchLabels: labelForDeployment1,
	}
	pointerSelectorForDeployment1 = &selectorForDeployment1

	/*labelForDeployment2["run"] = "test-depl2"
	replicasForDeployment2 = 1
	pointerReplicasForDeployment2 = &replicasForDeployment2
	selectorForDeployment2 = metav1.LabelSelector{
		MatchLabels: labelForDeployment2,
	}
	pointerSelectorForDeployment2 = &selectorForDeployment2

	labelForDeployment3["run"] = "test-depl3"
	replicasForDeployment3 = 3
	pointerReplicasForDeployment3 = &replicasForDeployment3
	selectorForDeployment3 = metav1.LabelSelector{
		MatchLabels: labelForDeployment3,
	}
	pointerSelectorForDeployment3 = &selectorForDeployment3*/

}

func TestCompareContainers(t *testing.T) {
	initEnvironment()
	clusterClientSet1 := fake.NewSimpleClientset(&v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v12.DeploymentSpec{
			Replicas: pointerReplicasForDeployment1,
			Selector: pointerSelectorForDeployment1,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "container-test-1",
							Image: "image-test-1",
						},
						{
							Name: "container-test-2",
							Image: "image-test-2",
						},
					},
				},
			},
		},
	})

	clusterClientSet2 := fake.NewSimpleClientset(&v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v12.DeploymentSpec{
			Replicas: pointerReplicasForDeployment1,
			Selector: pointerSelectorForDeployment1,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "container-test-1",
							Image: "image-test-1",
						},
					},
				},
			},
		},
	})
	deployments1, _ := clusterClientSet1.AppsV1().Deployments("default").List(metav1.ListOptions{})
	deployments2, _ := clusterClientSet2.AppsV1().Deployments("default").List(metav1.ListOptions{})
	objectInformation1.Selector = deployments1.Items[0].Spec.Selector
	objectInformation1.Template = deployments1.Items[0].Spec.Template
	objectInformation2.Selector = deployments2.Items[0].Spec.Selector
	objectInformation2.Template = deployments2.Items[0].Spec.Template
	v := CompareContainers(objectInformation1, objectInformation2, "default", clusterClientSet1, clusterClientSet2)
	fdsfsfs :=errors.New("The number templates of containers differs")
	if v != fdsfsfs {
		t.Error("Error expected: 'The number templates of containers differs'. But it was returned: ", v)
	}
}
