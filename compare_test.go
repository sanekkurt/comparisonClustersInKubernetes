package main

import (
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

var (

	clusterClientSet1 *fake.Clientset
	clusterClientSet2 *fake.Clientset

	labelForDeployment1           map[string]string
	replicasForDeployment1        int32
	pointerReplicasForDeployment1 *int32
	selectorForDeployment1        metav1.LabelSelector
	pointerSelectorForDeployment1 *metav1.LabelSelector


	labelForDeployment1Fake map[string]string
	selectorForDeployment1Fake        metav1.LabelSelector
	pointerSelectorForDeployment1Fake *metav1.LabelSelector

	objectInformation1            InformationAboutObject
	objectInformation2            InformationAboutObject
)

func initEnvironmentForFirstTest() {
	labelForDeployment1 = make(map[string]string)
	labelForDeployment1["run"] = "test-depl1"
	replicasForDeployment1 = 2
	pointerReplicasForDeployment1 = &replicasForDeployment1
	selectorForDeployment1 = metav1.LabelSelector{
		MatchLabels: labelForDeployment1,
	}
	pointerSelectorForDeployment1 = &selectorForDeployment1
	clusterClientSet1 = fake.NewSimpleClientset(&v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1",
			Namespace: "default",
			Labels:    labelForDeployment1,
		},
		Spec: v12.DeploymentSpec{
			Replicas: pointerReplicasForDeployment1,
			Selector: pointerSelectorForDeployment1,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "container-test-1",
							Image: "image-test-1",
						},
						{
							Name:  "container-test-2",
							Image: "image-test-2",
						},
					},
				},
			},
		},
	})

	clusterClientSet2 = fake.NewSimpleClientset(&v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1",
			Namespace: "default",
			Labels:    labelForDeployment1,
		},
		Spec: v12.DeploymentSpec{
			Replicas: pointerReplicasForDeployment1,
			Selector: pointerSelectorForDeployment1,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "container-test-1",
							Image: "image-test-1",
						},
					},
				},
			},
		},
	})
}

func initEnvironmentForSecondTest() {
	labelForDeployment1Fake = make(map[string]string)
	labelForDeployment1Fake["start"] = "test-depl1"
	selectorForDeployment1Fake = metav1.LabelSelector{
		MatchLabels: labelForDeployment1Fake,
	}
	pointerSelectorForDeployment1Fake = &selectorForDeployment1Fake
	clusterClientSet2 = fake.NewSimpleClientset(&v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1",
			Namespace: "default",
			Labels:    labelForDeployment1Fake,
		},
		Spec: v12.DeploymentSpec{
			Replicas: pointerReplicasForDeployment1,
			Selector: pointerSelectorForDeployment1Fake,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "container-test-1",
							Image: "image-test-1",
						},
						{
							Name:  "container-test-2",
							Image: "image-test-2",
						},
					},
				},
			},
		},
	})
}

func initEnvironmentForThirdTest() {
	clusterClientSet1 = fake.NewSimpleClientset(&v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1",
			Namespace: "default",
			Labels:    labelForDeployment1,
		},
		Spec: v12.DeploymentSpec{
			Replicas: pointerReplicasForDeployment1,
			Selector: pointerSelectorForDeployment1,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "container-test-1",
							Image: "image-test-1",
						},
						{
							Name:  "container-test-2-fake",
							Image: "image-test-2",
						},
					},
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-86c99d8f49-dpzrz",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "container-test-1",
						Image: "image-test-1",
					},
					{
						Name:  "container-test-2-fake",
						Image: "image-test-2",
					},
				},
		},
		Status: v1.PodStatus{
				ContainerStatuses: []v1.ContainerStatus{
					{
						Name:  "container-test-1",
						Image: "image-test-1",
						ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
					},
					{
						Name:  "container-test-2-fake",
						Image: "image-test-2",
						ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
					},
				},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-86c99d8f49-shc7j",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2-fake",
					Image: "image-test-2",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2-fake",
					Image: "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	})

	clusterClientSet2 = fake.NewSimpleClientset(&v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1",
			Namespace: "default",
			Labels:    labelForDeployment1,
		},
		Spec: v12.DeploymentSpec{
			Replicas: pointerReplicasForDeployment1,
			Selector: pointerSelectorForDeployment1,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "container-test-1",
							Image: "image-test-1",
						},
						{
							Name:  "container-test-2",
							Image: "image-test-2",
						},
					},
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-987646d67-kflz7",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-987646d67-x4sgn",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	})
}

func initEnvironmentForFourthTest() {
	clusterClientSet1 = fake.NewSimpleClientset(&v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1",
			Namespace: "default",
			Labels:    labelForDeployment1,
		},
		Spec: v12.DeploymentSpec{
			Replicas: pointerReplicasForDeployment1,
			Selector: pointerSelectorForDeployment1,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "container-test-1",
							Image: "image-test-1",
						},
						{
							Name:  "container-test-2",
							Image: "image-test-2",
						},
					},
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-86c99d8f49-dpzrz",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-86c99d8f49-shc7j",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	})

	clusterClientSet2 = fake.NewSimpleClientset(&v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1",
			Namespace: "default",
			Labels:    labelForDeployment1,
		},
		Spec: v12.DeploymentSpec{
			Replicas: pointerReplicasForDeployment1,
			Selector: pointerSelectorForDeployment1,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "container-test-1",
							Image: "image-test-1",
						},
						{
							Name:  "container-test-2",
							Image: "image-test-2-fake",
						},
					},
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-987646d67-kflz7",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2-fake",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2-fake",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-987646d67-x4sgn",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2-fake",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2-fake",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	})
}

func initEnvironmentForFifthTest() {
	clusterClientSet1 = fake.NewSimpleClientset(&v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1",
			Namespace: "default",
			Labels:    labelForDeployment1,
		},
		Spec: v12.DeploymentSpec{
			Replicas: pointerReplicasForDeployment1,
			Selector: pointerSelectorForDeployment1,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "container-test-1",
							Image: "image-test-1",
						},
						{
							Name:  "container-test-2",
							Image: "image-test-2",
						},
					},
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-86c99d8f49-dpzrz",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-86c99d8f49-shc7j",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	})

	clusterClientSet2 = fake.NewSimpleClientset(&v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1",
			Namespace: "default",
			Labels:    labelForDeployment1,
		},
		Spec: v12.DeploymentSpec{
			Replicas: pointerReplicasForDeployment1,
			Selector: pointerSelectorForDeployment1,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "container-test-1",
							Image: "image-test-1",
						},
						{
							Name:  "container-test-2",
							Image: "image-test-2",
						},
					},
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-987646d67-kflz7",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	})
}

func initEnvironmentForSixthTest() {
	clusterClientSet1 = fake.NewSimpleClientset(&v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1",
			Namespace: "default",
			Labels:    labelForDeployment1,
		},
		Spec: v12.DeploymentSpec{
			Replicas: pointerReplicasForDeployment1,
			Selector: pointerSelectorForDeployment1,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "container-test-1",
							Image: "image-test-1",
						},
						{
							Name:  "container-test-2",
							Image: "image-test-2",
						},
					},
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-86c99d8f49-dpzrz",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-86c99d8f49-shc7j",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
			},
		},
	})

	clusterClientSet2 = fake.NewSimpleClientset(&v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1",
			Namespace: "default",
			Labels:    labelForDeployment1,
		},
		Spec: v12.DeploymentSpec{
			Replicas: pointerReplicasForDeployment1,
			Selector: pointerSelectorForDeployment1,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "container-test-1",
							Image: "image-test-1",
						},
						{
							Name:  "container-test-2",
							Image: "image-test-2",
						},
					},
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-987646d67-kflz7",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-987646d67-x4sgn",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	})
}

func initEnvironmentForSeventhTest() {
	clusterClientSet1 = fake.NewSimpleClientset(&v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1",
			Namespace: "default",
			Labels:    labelForDeployment1,
		},
		Spec: v12.DeploymentSpec{
			Replicas: pointerReplicasForDeployment1,
			Selector: pointerSelectorForDeployment1,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "container-test-1",
							Image: "image-test-1",
						},
						{
							Name:  "container-test-2",
							Image: "image-test-2",
						},
					},
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-86c99d8f49-dpzrz",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	})

	clusterClientSet2 = fake.NewSimpleClientset(&v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1",
			Namespace: "default",
			Labels:    labelForDeployment1,
		},
		Spec: v12.DeploymentSpec{
			Replicas: pointerReplicasForDeployment1,
			Selector: pointerSelectorForDeployment1,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "container-test-1",
							Image: "image-test-1",
						},
						{
							Name:  "container-test-2",
							Image: "image-test-2",
						},
					},
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-987646d67-kflz7",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2-fake",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2-fake",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	})
}

func initEnvironmentForEighthTest() {
	clusterClientSet1 = fake.NewSimpleClientset(&v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1",
			Namespace: "default",
			Labels:    labelForDeployment1,
		},
		Spec: v12.DeploymentSpec{
			Replicas: pointerReplicasForDeployment1,
			Selector: pointerSelectorForDeployment1,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "container-test-1",
							Image: "image-test-1",
						},
						{
							Name:  "container-test-2",
							Image: "image-test-2",
						},
					},
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-86c99d8f49-dpzrz",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-86c99d8f49-shc7j",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	})

	clusterClientSet2 = fake.NewSimpleClientset(&v12.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1",
			Namespace: "default",
			Labels:    labelForDeployment1,
		},
		Spec: v12.DeploymentSpec{
			Replicas: pointerReplicasForDeployment1,
			Selector: pointerSelectorForDeployment1,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "container-test-1",
							Image: "image-test-1",
						},
						{
							Name:  "container-test-2",
							Image: "image-test-2",
						},
					},
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-987646d67-kflz7",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployment-1-987646d67-x4sgn",
			Namespace: "default",
			Labels: labelForDeployment1,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
				},
			},
		},
		Status: v1.PodStatus{
			ContainerStatuses: []v1.ContainerStatus{
				{
					Name:  "container-test-1",
					Image: "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:  "container-test-2",
					Image: "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f914704f3a1c6065fb56cfa34e",
				},
			},
		},
	})
}

func TestCompareContainers(t *testing.T) {
	//проверка на разное количество контейнеров в шаблонах
	initEnvironmentForFirstTest()
	deployments1, _ := clusterClientSet1.AppsV1().Deployments("default").List(metav1.ListOptions{})
	deployments2, _ := clusterClientSet2.AppsV1().Deployments("default").List(metav1.ListOptions{})
	objectInformation1.Selector = deployments1.Items[0].Spec.Selector
	objectInformation1.Template = deployments1.Items[0].Spec.Template
	objectInformation2.Selector = deployments2.Items[0].Spec.Selector
	objectInformation2.Template = deployments2.Items[0].Spec.Template
	v := CompareContainers(objectInformation1, objectInformation2, "default", clusterClientSet1, clusterClientSet2)
	if v != ErrorDiffersTemplatesNumber   {
		t.Error("Error expected: 'The number templates of containers differs'. But it was returned: ", v)
	}

	//Проверка на несовпадение MatchLabels
	initEnvironmentForSecondTest()
	deployments2, _ = clusterClientSet2.AppsV1().Deployments("default").List(metav1.ListOptions{})
	objectInformation2.Selector = deployments2.Items[0].Spec.Selector
	objectInformation2.Template = deployments2.Items[0].Spec.Template
	v = CompareContainers(objectInformation1, objectInformation2, "default", clusterClientSet1, clusterClientSet2)
	if v != ErrorMatchlabelsNotEqual   {
		t.Error("Error expected: 'MatchLabels are not equal'. But it was returned: ", v)
	}

	//Проверка на несовпадение имен контейнеров в шаблоне
	initEnvironmentForThirdTest()
	deployments1, _ = clusterClientSet1.AppsV1().Deployments("default").List(metav1.ListOptions{})
	deployments2, _ = clusterClientSet2.AppsV1().Deployments("default").List(metav1.ListOptions{})
	objectInformation1.Selector = deployments1.Items[0].Spec.Selector
	objectInformation1.Template = deployments1.Items[0].Spec.Template
	objectInformation2.Selector = deployments2.Items[0].Spec.Selector
	objectInformation2.Template = deployments2.Items[0].Spec.Template
	v = CompareContainers(objectInformation1, objectInformation2, "default", clusterClientSet1, clusterClientSet2)
	if v != ErrorContainerNamesTemplate   {
		t.Error("Error expected: 'Container names in template are not equal'. But it was returned: ", v)
	}

	//Првоерка на несовпадение имен образов контейнеров в шаблоне
	initEnvironmentForFourthTest()
	deployments1, _ = clusterClientSet1.AppsV1().Deployments("default").List(metav1.ListOptions{})
	deployments2, _ = clusterClientSet2.AppsV1().Deployments("default").List(metav1.ListOptions{})
	objectInformation1.Selector = deployments1.Items[0].Spec.Selector
	objectInformation1.Template = deployments1.Items[0].Spec.Template
	objectInformation2.Selector = deployments2.Items[0].Spec.Selector
	objectInformation2.Template = deployments2.Items[0].Spec.Template
	v = CompareContainers(objectInformation1, objectInformation2, "default", clusterClientSet1, clusterClientSet2)
	if v != ErrorContainerImagesTemplate   {
		t.Error("Error expected: 'Container name images in template are not equal'. But it was returned: ", v)
	}

	//Проверка на разное количество Pod'ов
	initEnvironmentForFifthTest()
	deployments1, _ = clusterClientSet1.AppsV1().Deployments("default").List(metav1.ListOptions{})
	deployments2, _ = clusterClientSet2.AppsV1().Deployments("default").List(metav1.ListOptions{})
	objectInformation1.Selector = deployments1.Items[0].Spec.Selector
	objectInformation1.Template = deployments1.Items[0].Spec.Template
	objectInformation2.Selector = deployments2.Items[0].Spec.Selector
	objectInformation2.Template = deployments2.Items[0].Spec.Template
	v = CompareContainers(objectInformation1, objectInformation2, "default", clusterClientSet1, clusterClientSet2)
	if v != ErrorPodsCount   {
		t.Error("Error expected: 'The pods count are different'. But it was returned: ", v)
	}

	//Проверка на разное количество контейнеров в Pod'ах
	initEnvironmentForSixthTest()
	deployments1, _ = clusterClientSet1.AppsV1().Deployments("default").List(metav1.ListOptions{})
	deployments2, _ = clusterClientSet2.AppsV1().Deployments("default").List(metav1.ListOptions{})
	objectInformation1.Selector = deployments1.Items[0].Spec.Selector
	objectInformation1.Template = deployments1.Items[0].Spec.Template
	objectInformation2.Selector = deployments2.Items[0].Spec.Selector
	objectInformation2.Template = deployments2.Items[0].Spec.Template
	v = CompareContainers(objectInformation1, objectInformation2, "default", clusterClientSet1, clusterClientSet2)
	if v != ErrorContainersCountInPod   {
		t.Error("Error expected: 'The containers count in pod are different'. But it was returned: ", v)
	}

	//Проверка на отличающиеся имена образов в Pod'e и Template
	initEnvironmentForSeventhTest()
	deployments1, _ = clusterClientSet1.AppsV1().Deployments("default").List(metav1.ListOptions{})
	deployments2, _ = clusterClientSet2.AppsV1().Deployments("default").List(metav1.ListOptions{})
	objectInformation1.Selector = deployments1.Items[0].Spec.Selector
	objectInformation1.Template = deployments1.Items[0].Spec.Template
	objectInformation2.Selector = deployments2.Items[0].Spec.Selector
	objectInformation2.Template = deployments2.Items[0].Spec.Template
	v = CompareContainers(objectInformation1, objectInformation2, "default", clusterClientSet1, clusterClientSet2)
	if v != ErrorContainerImageTemplatePod   {
		t.Error("Error expected: 'The container image in the template does not match the actual image in the Pod'. But it was returned: ", v)
	}

	//Проверка на разные ImageID в Pod'ах
	initEnvironmentForEighthTest()
	deployments1, _ = clusterClientSet1.AppsV1().Deployments("default").List(metav1.ListOptions{})
	deployments2, _ = clusterClientSet2.AppsV1().Deployments("default").List(metav1.ListOptions{})
	objectInformation1.Selector = deployments1.Items[0].Spec.Selector
	objectInformation1.Template = deployments1.Items[0].Spec.Template
	objectInformation2.Selector = deployments2.Items[0].Spec.Selector
	objectInformation2.Template = deployments2.Items[0].Spec.Template
	v = CompareContainers(objectInformation1, objectInformation2, "default", clusterClientSet1, clusterClientSet2)
	if v != ErrorDifferentImageIdInPods   {
		t.Error("Error expected: 'The ImageID in Pods is different'. But it was returned: ", v)
	}

	//Проверка на случай отсутсвия контейнера в другом Pod'е не нужна она и так в стоке работает, на нее реагирует проверка количества контейнеров
}
