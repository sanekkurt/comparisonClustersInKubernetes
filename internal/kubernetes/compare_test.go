package kubernetes

import (
	"errors"
	v12 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

var (
	// variables for the container comparison function test
	clusterClientSet1 *fake.Clientset
	clusterClientSet2 *fake.Clientset

	labelForDeployment1           map[string]string
	replicasForDeployment1        int32
	pointerReplicasForDeployment1 *int32
	selectorForDeployment1        metav1.LabelSelector
	pointerSelectorForDeployment1 *metav1.LabelSelector

	labelForDeployment1Fake           map[string]string
	selectorForDeployment1Fake        metav1.LabelSelector
	pointerSelectorForDeployment1Fake *metav1.LabelSelector

	objectInformation1 InformationAboutObject
	objectInformation2 InformationAboutObject

	// variables for testing the variable comparison function in containers
	env1                   []v1.EnvVar
	env2                   []v1.EnvVar
	temp                   v1.EnvVar
	dataInFirstConfigMap   map[string]string
	dataInSecondConfigMap  map[string]string
	dataInFirstSecret      map[string][]byte
	dataInSecondSecret     map[string][]byte
	envVarSource           v1.EnvVarSource
	pointerEnvVarSource    *v1.EnvVarSource
	secretKeyRef           v1.SecretKeySelector
	pointerSecretKeyRef    *v1.SecretKeySelector
	configMapKeyRef        v1.ConfigMapKeySelector
	pointerConfigMapKeyRef *v1.ConfigMapKeySelector

	envVarSource2        v1.EnvVarSource
	pointerEnvVarSource2 *v1.EnvVarSource
	secretKeyRef2        v1.SecretKeySelector
	pointerSecretKeyRef2 *v1.SecretKeySelector

	// variables for testing the spec comparison function in services
	selector1 map[string]string
	selector2 map[string]string
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
			Name:      "deployment-1-86c99d8f49-dpzrz",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2-fake",
					Image:   "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1-86c99d8f49-shc7j",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2-fake",
					Image:   "image-test-2",
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
			Name:      "deployment-1-987646d67-kflz7",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1-987646d67-x4sgn",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2",
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
			Name:      "deployment-1-86c99d8f49-dpzrz",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1-86c99d8f49-shc7j",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2",
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
			Name:      "deployment-1-987646d67-kflz7",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2-fake",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1-987646d67-x4sgn",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2-fake",
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
			Name:      "deployment-1-86c99d8f49-dpzrz",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1-86c99d8f49-shc7j",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2",
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
			Name:      "deployment-1-987646d67-kflz7",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2",
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
			Name:      "deployment-1-86c99d8f49-dpzrz",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1-86c99d8f49-shc7j",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
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
			Name:      "deployment-1-987646d67-kflz7",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1-987646d67-x4sgn",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2",
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
			Name:      "deployment-1-86c99d8f49-dpzrz",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2",
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
			Name:      "deployment-1-987646d67-kflz7",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2-fake",
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
			Name:      "deployment-1-86c99d8f49-dpzrz",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1-86c99d8f49-shc7j",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2",
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
			Name:      "deployment-1-987646d67-kflz7",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f91470a23ef380320609ee6e69ac68d20bc804f3a1c6065fb56cfa34e",
				},
			},
		},
	}, &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1-987646d67-x4sgn",
			Namespace: "default",
			Labels:    labelForDeployment1,
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
					Name:    "container-test-1",
					Image:   "image-test-1",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-1@sha256:e8fc56926ac3d5705772f13befbaee3aa2fc6e9c52faee3d96b26612cd77556c",
				},
				{
					Name:    "container-test-2",
					Image:   "image-test-2",
					ImageID: "docker-hub.binary.alfabank.ru/image-test-2@sha256:7d6a3c8f914704f3a1c6065fb56cfa34e",
				},
			},
		},
	})
}

// TestCompareContainers check CompareContainers function
func TestCompareContainers(t *testing.T) {
	// Checking for different number of containers in templates
	initEnvironmentForFirstTest()
	deployments1, _ := clusterClientSet1.AppsV1().Deployments("default").List(metav1.ListOptions{})
	deployments2, _ := clusterClientSet2.AppsV1().Deployments("default").List(metav1.ListOptions{})
	objectInformation1.Selector = deployments1.Items[0].Spec.Selector
	objectInformation1.Template = deployments1.Items[0].Spec.Template
	objectInformation2.Selector = deployments2.Items[0].Spec.Selector
	objectInformation2.Template = deployments2.Items[0].Spec.Template
	err := CompareContainers(objectInformation1, objectInformation2, "default", clusterClientSet1, clusterClientSet2)
	if !errors.Is(err, ErrorDiffersTemplatesNumber) {
		t.Error("Error expected: 'The number templates of containers differs'. But it was returned: ", err)
	}

	// Checking for MatchLabels mismatch
	initEnvironmentForSecondTest()
	deployments2, _ = clusterClientSet2.AppsV1().Deployments("default").List(metav1.ListOptions{})
	objectInformation2.Selector = deployments2.Items[0].Spec.Selector
	objectInformation2.Template = deployments2.Items[0].Spec.Template
	err = CompareContainers(objectInformation1, objectInformation2, "default", clusterClientSet1, clusterClientSet2)
	if !errors.Is(err, ErrorMatchlabelsNotEqual) {
		t.Error("Error expected: 'MatchLabels are not equal'. But it was returned: ", err)
	}

	// Check for mismatch of names of the containers in the template
	initEnvironmentForThirdTest()
	deployments1, _ = clusterClientSet1.AppsV1().Deployments("default").List(metav1.ListOptions{})
	deployments2, _ = clusterClientSet2.AppsV1().Deployments("default").List(metav1.ListOptions{})
	objectInformation1.Selector = deployments1.Items[0].Spec.Selector
	objectInformation1.Template = deployments1.Items[0].Spec.Template
	objectInformation2.Selector = deployments2.Items[0].Spec.Selector
	objectInformation2.Template = deployments2.Items[0].Spec.Template
	err = CompareContainers(objectInformation1, objectInformation2, "default", clusterClientSet1, clusterClientSet2)
	if !errors.Is(err, ErrorContainerNamesTemplate) {
		t.Error("Error expected: 'Container names in template are not equal'. But it was returned: ", err)
	}

	// Checking for mismatched container image names in the template
	initEnvironmentForFourthTest()
	deployments1, _ = clusterClientSet1.AppsV1().Deployments("default").List(metav1.ListOptions{})
	deployments2, _ = clusterClientSet2.AppsV1().Deployments("default").List(metav1.ListOptions{})
	objectInformation1.Selector = deployments1.Items[0].Spec.Selector
	objectInformation1.Template = deployments1.Items[0].Spec.Template
	objectInformation2.Selector = deployments2.Items[0].Spec.Selector
	objectInformation2.Template = deployments2.Items[0].Spec.Template

	err = CompareContainers(objectInformation1, objectInformation2, "default", clusterClientSet1, clusterClientSet2)
	if !errors.Is(err, ErrorContainerImagesTemplate) {
		t.Error("Error expected: 'Container name images in template are not equal'. But it was returned: ", err)
	}

	// Checking for different counts Pods
	initEnvironmentForFifthTest()
	deployments1, _ = clusterClientSet1.AppsV1().Deployments("default").List(metav1.ListOptions{})
	deployments2, _ = clusterClientSet2.AppsV1().Deployments("default").List(metav1.ListOptions{})
	objectInformation1.Selector = deployments1.Items[0].Spec.Selector
	objectInformation1.Template = deployments1.Items[0].Spec.Template
	objectInformation2.Selector = deployments2.Items[0].Spec.Selector
	objectInformation2.Template = deployments2.Items[0].Spec.Template
	err = CompareContainers(objectInformation1, objectInformation2, "default", clusterClientSet1, clusterClientSet2)
	if !errors.Is(err, ErrorPodsCount) {
		t.Error("Error expected: 'The pods count are different'. But it was returned: ", err)
	}

	// Checking for different number of containers in Pods
	initEnvironmentForSixthTest()
	deployments1, _ = clusterClientSet1.AppsV1().Deployments("default").List(metav1.ListOptions{})
	deployments2, _ = clusterClientSet2.AppsV1().Deployments("default").List(metav1.ListOptions{})
	objectInformation1.Selector = deployments1.Items[0].Spec.Selector
	objectInformation1.Template = deployments1.Items[0].Spec.Template
	objectInformation2.Selector = deployments2.Items[0].Spec.Selector
	objectInformation2.Template = deployments2.Items[0].Spec.Template
	err = CompareContainers(objectInformation1, objectInformation2, "default", clusterClientSet1, clusterClientSet2)
	if !errors.Is(err, ErrorContainersCountInPod) {
		t.Error("Error expected: 'The containers count in pod are different'. But it was returned: ", err)
	}

	// Checking for different image names in Pod and Template
	initEnvironmentForSeventhTest()
	deployments1, _ = clusterClientSet1.AppsV1().Deployments("default").List(metav1.ListOptions{})
	deployments2, _ = clusterClientSet2.AppsV1().Deployments("default").List(metav1.ListOptions{})
	objectInformation1.Selector = deployments1.Items[0].Spec.Selector
	objectInformation1.Template = deployments1.Items[0].Spec.Template
	objectInformation2.Selector = deployments2.Items[0].Spec.Selector
	objectInformation2.Template = deployments2.Items[0].Spec.Template
	err = CompareContainers(objectInformation1, objectInformation2, "default", clusterClientSet1, clusterClientSet2)
	if !errors.Is(err, ErrorContainerImageTemplatePod) {
		t.Error("Error expected: 'The container image in the template does not match the actual image in the Pod'. But it was returned: ", err)
	}

	// Checking for different ImageID's in Pods
	initEnvironmentForEighthTest()
	deployments1, _ = clusterClientSet1.AppsV1().Deployments("default").List(metav1.ListOptions{})
	deployments2, _ = clusterClientSet2.AppsV1().Deployments("default").List(metav1.ListOptions{})
	objectInformation1.Selector = deployments1.Items[0].Spec.Selector
	objectInformation1.Template = deployments1.Items[0].Spec.Template
	objectInformation2.Selector = deployments2.Items[0].Spec.Selector
	objectInformation2.Template = deployments2.Items[0].Spec.Template
	err = CompareContainers(objectInformation1, objectInformation2, "default", clusterClientSet1, clusterClientSet2)
	if !errors.Is(errors.Unwrap(err), ErrorDifferentImageIDInPods) {
		t.Error("Error expected: 'The ImageID in Pods is different'. But it was returned: ", err)
	}
}

func initEnvironmentForFirstTest2() {

	temp.Name = "Hello"
	temp.Value = "World"
	env1 = append(env1, temp)
	env2 = append(env2, temp)

	temp.Name = "I"
	temp.Value = "Am"
	env1 = append(env1, temp)

	temp.Name = "From"
	temp.Value = "Russia"
	env2 = append(env2, temp)

	temp.Name = "The"
	temp.Value = "End"
	env2 = append(env2, temp)

	dataInFirstConfigMap = make(map[string]string)
	dataInSecondConfigMap = make(map[string]string)

	clusterClientSet1 = fake.NewSimpleClientset(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "configMapInCluster",
			Namespace: "default",
		},
		Data: dataInFirstConfigMap,
	},
	)
	clusterClientSet2 = fake.NewSimpleClientset(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "configMapInCluster",
			Namespace: "default",
		},
		Data: dataInSecondConfigMap,
	},
	)
}

func initEnvironmentForSecondTest2() {
	temp.Name = "Not"
	temp.Value = "End"
	env1 = append(env1, temp)
}

func initEnvironmentForThirdTest2() {
	env1 = append(env1[:1], env1[1+2:]...)
	env2 = append(env2[:1], env2[1+2:]...)
	secretKeyRef.Name = "test"
	secretKeyRef.Key = "test"
	pointerSecretKeyRef = &secretKeyRef
	envVarSource.SecretKeyRef = pointerSecretKeyRef
	pointerEnvVarSource = &envVarSource
	temp.Name = "testValueFrom"
	temp.ValueFrom = pointerEnvVarSource
	env1 = append(env1, temp)
	secretKeyRef2.Name = "test"
	secretKeyRef2.Key = "test2222"
	pointerSecretKeyRef2 = &secretKeyRef2
	envVarSource2.SecretKeyRef = pointerSecretKeyRef2
	pointerEnvVarSource2 = &envVarSource2
	temp.Name = "testValueFrom"
	temp.ValueFrom = pointerEnvVarSource2
	env2 = append(env2, temp)
}

func initEnvironmentForFourthTest2() {
	env1 = append(env1[:1], env1[1+1:]...)
	env2 = append(env2[:1], env2[1+1:]...)

	secretKeyRef.Name = "secretInCluster"
	secretKeyRef.Key = "test"
	pointerSecretKeyRef = &secretKeyRef
	envVarSource.SecretKeyRef = pointerSecretKeyRef
	pointerEnvVarSource = &envVarSource
	temp.Name = "testValueFrom"
	temp.ValueFrom = pointerEnvVarSource
	env1 = append(env1, temp)
	env2 = append(env2, temp)

	dataInFirstSecret = make(map[string][]byte)
	dataInSecondSecret = make(map[string][]byte)

	dataInFirstSecret["test"] = []byte("test")
	dataInSecondSecret["test"] = []byte("fakeTest")

	clusterClientSet1 = fake.NewSimpleClientset(&v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secretInCluster",
			Namespace: "default",
		},
		Data: dataInFirstSecret,
	},
	)
	clusterClientSet2 = fake.NewSimpleClientset(&v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secretInCluster",
			Namespace: "default",
		},
		Data: dataInSecondSecret,
	},
	)
}

func initEnvironmentForFifthTest2() {
	env1 = append(env1[:1], env1[1+1:]...)
	env2 = append(env2[:1], env2[1+1:]...)

	configMapKeyRef.Name = "configMapInCluster"
	configMapKeyRef.Key = "test"
	pointerConfigMapKeyRef = &configMapKeyRef
	envVarSource.SecretKeyRef = nil
	envVarSource.ConfigMapKeyRef = pointerConfigMapKeyRef
	pointerEnvVarSource = &envVarSource

	env1 = append(env1, temp)
	env2 = append(env2, temp)

	dataInFirstConfigMap["test"] = "test"
	dataInSecondConfigMap["test"] = "fakeTest"

	clusterClientSet1 = fake.NewSimpleClientset(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "configMapInCluster",
			Namespace: "default",
		},
		Data: dataInFirstConfigMap,
	},
	)
	clusterClientSet2 = fake.NewSimpleClientset(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "configMapInCluster",
			Namespace: "default",
		},
		Data: dataInSecondConfigMap,
	},
	)
}

// TestCompareEnvInContainers check CompareEnvInContainers function
func TestCompareEnvInContainers(t *testing.T) {
	initEnvironmentForFirstTest2()
	err := CompareEnvInContainers(env1, env2, "default", clusterClientSet1, clusterClientSet2)
	if !errors.Is(err, ErrorNumberVariables) {
		t.Error("Error expected: 'The number of variables in containers differs'. But it was returned: ", err)
	}

	initEnvironmentForSecondTest2()
	err = CompareEnvInContainers(env1, env2, "default", clusterClientSet1, clusterClientSet2)
	if !errors.Is(errors.Unwrap(err), ErrorEnvironmentNotEqual) {
		t.Error("Error expected: 'The environment in containers not equal'. But it was returned: ", err)
	}

	initEnvironmentForThirdTest2()
	err = CompareEnvInContainers(env1, env2, "default", clusterClientSet1, clusterClientSet2)
	if !errors.Is(errors.Unwrap(err), ErrorEnvironmentNotEqual) {
		t.Error("Error expected: 'The environment in containers not equal'. But it was returned: ", err)
	}

	initEnvironmentForFourthTest2()
	err = CompareEnvInContainers(env1, env2, "default", clusterClientSet1, clusterClientSet2)
	if !errors.Is(errors.Unwrap(err), ErrorDifferentValueSecretKey) {
		t.Error("Error expected: 'The value for the SecretKey is different'. But it was returned: ", err)
	}

	initEnvironmentForFifthTest2()
	err = CompareEnvInContainers(env1, env2, "default", clusterClientSet1, clusterClientSet2)
	if !errors.Is(errors.Unwrap(err), ErrorDifferentValueConfigMapKey) {
		t.Error("Error expected: 'The value for the ConfigMapKey is different'. But it was returned: ", err)
	}
}

func initEnvironmentForFirstTest3() {
	clusterClientSet1 = fake.NewSimpleClientset(&v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testService",
			Namespace: "default",
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "port1",
					NodePort: 80,
				},
			},
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testService",
			Namespace: "default",
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "port1",
					NodePort: 80,
				},
				{
					Name:     "port2",
					NodePort: 81,
				},
			},
		},
	})
}

func initEnvironmentForSecondTest3() {
	clusterClientSet1 = fake.NewSimpleClientset(&v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testService",
			Namespace: "default",
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "port1",
					NodePort: 80,
				},
				{
					Name:     "port2",
					NodePort: 81,
				},
			},
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testService",
			Namespace: "default",
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "port1",
					NodePort: 80,
				},
				{
					Name:     "port2",
					NodePort: 88,
				},
			},
		},
	})
}

func initEnvironmemtForThirdTest3() {
	selector1 = make(map[string]string)
	selector2 = make(map[string]string)
	selector1["one"] = "1"
	selector1["two"] = "2"
	selector2["one"] = "1"
	selector2["two"] = "2"
	selector2["three"] = "3"

	clusterClientSet1 = fake.NewSimpleClientset(&v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testService",
			Namespace: "default",
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "port1",
					NodePort: 80,
				},
			},
			Selector: selector1,
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testService",
			Namespace: "default",
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "port1",
					NodePort: 80,
				},
			},
			Selector: selector2,
		},
	})
}

func initEnvironmemtForFourthTest3() {
	selector1["three"] = "1"

	clusterClientSet1 = fake.NewSimpleClientset(&v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testService",
			Namespace: "default",
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "port1",
					NodePort: 80,
				},
			},
			Selector: selector1,
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testService",
			Namespace: "default",
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "port1",
					NodePort: 80,
				},
			},
			Selector: selector2,
		},
	})
}

func initEnvironmemtForFifthTest3() {
	selector1["three"] = "3"
	clusterClientSet1 = fake.NewSimpleClientset(&v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testService",
			Namespace: "default",
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "port1",
					NodePort: 80,
				},
			},
			Selector: selector1,
			Type:     "string",
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testService",
			Namespace: "default",
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "port1",
					NodePort: 80,
				},
			},
			Selector: selector2,
			Type:     "int",
		},
	})
}

func TestCompareSpecInServices(t *testing.T) {
	initEnvironmentForFirstTest3()
	service1, _ := clusterClientSet1.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	service2, _ := clusterClientSet2.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	err := CompareSpecInServices(*service1, *service2)
	if !errors.Is(errors.Unwrap(err), ErrorPortsCountDifferent) {
		t.Error("Error expected: 'the ports count are different'. But it was returned: ", err)
	}

	initEnvironmentForSecondTest3()
	service1, _ = clusterClientSet1.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	service2, _ = clusterClientSet2.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	err = CompareSpecInServices(*service1, *service2)
	if !errors.Is(errors.Unwrap(err), ErrorPortInServicesDifferent) {
		t.Error("Error expected: 'the port in the services is different'. But it was returned: ", err)
	}

	initEnvironmemtForThirdTest3()
	service1, _ = clusterClientSet1.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	service2, _ = clusterClientSet2.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	err = CompareSpecInServices(*service1, *service2)
	if !errors.Is(errors.Unwrap(err), ErrorSelectorsCountDifferent) {
		t.Error("Error expected: 'the selectors count are different'. But it was returned: ", err)
	}

	initEnvironmemtForFourthTest3()
	service1, _ = clusterClientSet1.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	service2, _ = clusterClientSet2.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	err = CompareSpecInServices(*service1, *service2)
	if !errors.Is(errors.Unwrap(err), ErrorSelectorInServicesDifferent) {
		t.Error("Error expected: 'the selector in the services is different'. But it was returned: ", err)
	}

	initEnvironmemtForFifthTest3()
	service1, _ = clusterClientSet1.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	service2, _ = clusterClientSet2.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	err = CompareSpecInServices(*service1, *service2)
	if !errors.Is(errors.Unwrap(err), ErrorTypeInServicesDifferent) {
		t.Error("the type in the services is different'. But it was returned: ", err)
	}
}

func initEnvironmentForFirthTest4() {
	clusterClientSet1 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
				},
			},
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
		},
	})
}

func initEnvironmentForSecondTest4() {
	clusterClientSet1 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName1",
				},
				{
					SecretName: "secretName2",
				},
			},
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
				},
			},
		},
	})
}

func initEnvironmentForThirdTest4() {
	clusterClientSet1 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName1",
				},
			},
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
				},
			},
		},
	})
}

func initEnvironmentForFifthTest4() {
	hosts1 := []string{"host1", "host2"}
	hosts2 := []string{"host2"}
	clusterClientSet1 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts2,
				},
			},
		},
	})
}

func initEnvironmentForSixthTest4() {
	hosts1 := []string{"host1", "host2"}
	hosts2 := []string{"host2", "host4"}
	clusterClientSet1 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts2,
				},
			},
		},
	})
}

func initEnvironmentForSeventhTest4() {
	hosts1 := []string{"host1", "host2"}
	clusterClientSet1 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
				},
			},
		},
	})
}

func initEnvironmentForEighthTest4() {
	backend1 := v1beta1.IngressBackend{
		ServiceName: "tempName",
		ServicePort: intstr.IntOrString{
			Type:   58,
			IntVal: 32,
			StrVal: "strVal",
		},
	}
	hosts1 := []string{"host1", "host2"}
	clusterClientSet1 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
			Backend: &backend1,
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
		},
	})
}

func initEnvironmentForNinthTest4() {
	backend1 := v1beta1.IngressBackend{
		ServiceName: "Name",
		ServicePort: intstr.IntOrString{
			Type:   58,
			IntVal: 32,
			StrVal: "strVal",
		},
	}
	backend2 := v1beta1.IngressBackend{
		ServiceName: "FakeName",
		ServicePort: intstr.IntOrString{
			Type:   58,
			IntVal: 32,
			StrVal: "strVal",
		},
	}
	hosts1 := []string{"host1", "host2"}
	clusterClientSet1 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
			Backend: &backend1,
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
			Backend: &backend2,
		},
	})
}

func initEnvironmentForTenthTest4() {
	backend1 := v1beta1.IngressBackend{
		ServiceName: "Name",
		ServicePort: intstr.IntOrString{
			Type:   58,
			IntVal: 32,
			StrVal: "fakeVal",
		},
	}
	backend2 := v1beta1.IngressBackend{
		ServiceName: "Name",
		ServicePort: intstr.IntOrString{
			Type:   58,
			IntVal: 32,
			StrVal: "strVal",
		},
	}
	hosts1 := []string{"host1", "host2"}
	clusterClientSet1 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
			Backend: &backend1,
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
			Backend: &backend2,
		},
	})
}

func initEnvironmentForEleventhTest4() {
	backend := v1beta1.IngressBackend{
		ServiceName: "Name",
		ServicePort: intstr.IntOrString{
			Type:   58,
			IntVal: 32,
			StrVal: "fakeVal",
		},
	}
	http := v1beta1.IngressRuleValue{
		HTTP: &v1beta1.HTTPIngressRuleValue{
			Paths: []v1beta1.HTTPIngressPath{
				{
					Path:    "dfdfdf",
					Backend: backend,
				},
			},
		},
	}
	rule1 := v1beta1.IngressRule{
		Host:             "host",
		IngressRuleValue: http,
	}
	var rules1 = []v1beta1.IngressRule{rule1}

	hosts1 := []string{"host1", "host2"}
	clusterClientSet1 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
			Rules: rules1,
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
		},
	})
}

func initEnvironmentForTwelvesTest4() {
	backend := v1beta1.IngressBackend{
		ServiceName: "Name",
		ServicePort: intstr.IntOrString{
			Type:   58,
			IntVal: 32,
			StrVal: "fakeVal",
		},
	}
	http := v1beta1.IngressRuleValue{
		HTTP: &v1beta1.HTTPIngressRuleValue{
			Paths: []v1beta1.HTTPIngressPath{
				{
					Path:    "dfdfdf",
					Backend: backend,
				},
			},
		},
	}
	rule1 := v1beta1.IngressRule{
		Host:             "host",
		IngressRuleValue: http,
	}
	var rules1 = []v1beta1.IngressRule{rule1, rule1}
	var rules2 = []v1beta1.IngressRule{rule1}

	hosts1 := []string{"host1", "host2"}
	clusterClientSet1 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
			Rules: rules1,
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
			Rules: rules2,
		},
	})
}

func initEnvironmentForThirteenthTest4() {
	backend := v1beta1.IngressBackend{
		ServiceName: "Name",
		ServicePort: intstr.IntOrString{
			Type:   58,
			IntVal: 32,
			StrVal: "fakeVal",
		},
	}
	http := v1beta1.IngressRuleValue{
		HTTP: &v1beta1.HTTPIngressRuleValue{
			Paths: []v1beta1.HTTPIngressPath{
				{
					Path:    "dfdfdf",
					Backend: backend,
				},
			},
		},
	}
	rule1 := v1beta1.IngressRule{
		Host:             "host",
		IngressRuleValue: http,
	}
	rule2 := v1beta1.IngressRule{
		Host:             "fakeHost",
		IngressRuleValue: http,
	}
	var rules1 = []v1beta1.IngressRule{rule1, rule1}
	var rules2 = []v1beta1.IngressRule{rule1, rule2}

	hosts1 := []string{"host1", "host2"}
	clusterClientSet1 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
			Rules: rules1,
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
			Rules: rules2,
		},
	})
}

func initEnvironmentForFourteenthTest4() {
	backend := v1beta1.IngressBackend{
		ServiceName: "Name",
		ServicePort: intstr.IntOrString{
			Type:   58,
			IntVal: 32,
			StrVal: "fakeVal",
		},
	}
	http := v1beta1.IngressRuleValue{
		HTTP: &v1beta1.HTTPIngressRuleValue{
			Paths: []v1beta1.HTTPIngressPath{
				{
					Path:    "dfdfdf",
					Backend: backend,
				},
			},
		},
	}
	rule1 := v1beta1.IngressRule{
		Host:             "host",
	}
	rule2 := v1beta1.IngressRule{
		Host:             "host",
		IngressRuleValue: http,
	}
	var rules1 = []v1beta1.IngressRule{rule1, rule1}
	var rules2 = []v1beta1.IngressRule{rule1, rule2}

	hosts1 := []string{"host1", "host2"}
	clusterClientSet1 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
			Rules: rules1,
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
			Rules: rules2,
		},
	})
}

func initEnvironmentForFifteenthTest4() {
	backend := v1beta1.IngressBackend{
		ServiceName: "Name",
		ServicePort: intstr.IntOrString{
			Type:   58,
			IntVal: 32,
			StrVal: "fakeVal",
		},
	}
	http1 := v1beta1.IngressRuleValue{
		HTTP: &v1beta1.HTTPIngressRuleValue{
			Paths: []v1beta1.HTTPIngressPath{
				{
					Path:    "path1",
					Backend: backend,
				},
				{
					Path:    "path2",
					Backend: backend,
				},
			},
		},
	}
	http2 := v1beta1.IngressRuleValue{
		HTTP: &v1beta1.HTTPIngressRuleValue{
			Paths: []v1beta1.HTTPIngressPath{
				{
					Path:    "path1",
					Backend: backend,
				},
			},
		},
	}
	rule1 := v1beta1.IngressRule{
		Host:             "host",
		IngressRuleValue: http1,
	}
	rule2 := v1beta1.IngressRule{
		Host:             "host",
		IngressRuleValue: http2,
	}
	var rules1 = []v1beta1.IngressRule{rule1, rule1}
	var rules2 = []v1beta1.IngressRule{rule1, rule2}

	hosts1 := []string{"host1", "host2"}
	clusterClientSet1 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
			Rules: rules1,
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
			Rules: rules2,
		},
	})
}

func initEnvironmentForSixteenthTest4() {
	backend := v1beta1.IngressBackend{
		ServiceName: "Name",
		ServicePort: intstr.IntOrString{
			Type:   58,
			IntVal: 32,
			StrVal: "fakeVal",
		},
	}
	http1 := v1beta1.IngressRuleValue{
		HTTP: &v1beta1.HTTPIngressRuleValue{
			Paths: []v1beta1.HTTPIngressPath{
				{
					Path:    "path1",
					Backend: backend,
				},
				{
					Path:    "path2",
					Backend: backend,
				},
			},
		},
	}
	http2 := v1beta1.IngressRuleValue{
		HTTP: &v1beta1.HTTPIngressRuleValue{
			Paths: []v1beta1.HTTPIngressPath{
				{
					Path:    "path1",
					Backend: backend,
				},
				{
					Path:    "path1",
					Backend: backend,
				},
			},
		},
	}
	rule1 := v1beta1.IngressRule{
		Host:             "host",
		IngressRuleValue: http1,
	}
	rule2 := v1beta1.IngressRule{
		Host:             "host",
		IngressRuleValue: http2,
	}
	var rules1 = []v1beta1.IngressRule{rule1, rule1}
	var rules2 = []v1beta1.IngressRule{rule1, rule2}

	hosts1 := []string{"host1", "host2"}
	clusterClientSet1 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
			Rules: rules1,
		},
	})
	clusterClientSet2 = fake.NewSimpleClientset(&v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testIngress",
			Namespace: "default",
		},
		Spec: v1beta1.IngressSpec{
			TLS: []v1beta1.IngressTLS{
				{
					SecretName: "secretName",
					Hosts:      hosts1,
				},
			},
			Rules: rules2,
		},
	})
}

func TestCompareSpecInIngresses(t *testing.T) {
	initEnvironmentForFirthTest4()
	ingress1, _ := clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ := clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err := CompareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorTLSInIngressesDifferent) {
		t.Error("the TLS in the ingresses are different'. But it was returned: ", err)
	}

	initEnvironmentForSecondTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = CompareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorTLSCountDifferent) {
		t.Error("the TLS count in the ingresses are different'. But it was returned: ", err)
	}

	initEnvironmentForThirdTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = CompareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorSecretNameInTLSDifferent) {
		t.Error("the secret name in the TLS are different'. But it was returned: ", err)
	}

	initEnvironmentForFifthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = CompareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorHostsCountDifferent) {
		t.Error("the hosts count in the TLS are different'. But it was returned: ", err)
	}

	initEnvironmentForSixthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = CompareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorNameHostDifferent) {
		t.Error("the name host in the TLS are different'. But it was returned: ", err)
	}

	initEnvironmentForSeventhTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = CompareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorHostsInIngressesDifferent) {
		t.Error("the hosts in the ingresses are different'. But it was returned: ", err)
	}

	initEnvironmentForEighthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = CompareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorBackendInIngressesDifferent) {
		t.Error("the backend in the ingresses are different'. But it was returned: ", err)
	}

	initEnvironmentForNinthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = CompareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorServiceNameInBackendDifferent) {
		t.Error("the service name in the backend are different'. But it was returned: ", err)
	}

	initEnvironmentForTenthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = CompareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorBackendServicePortDifferent) {
		t.Error("the service port in the backend are different'. But it was returned: ", err)
	}

	initEnvironmentForEleventhTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = CompareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorRulesInIngressesDifferent) {
		t.Error("the rules in the ingresses are different'. But it was returned: ", err)
	}

	initEnvironmentForTwelvesTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = CompareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorRulesCountDifferent) {
		t.Error("the rules count in the ingresses is different'. But it was returned: ", err)
	}

	initEnvironmentForThirteenthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = CompareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorHostNameInRuleDifferent) {
		t.Error("the hosts name in the rule are different'. But it was returned: ", err)
	}

	initEnvironmentForFourteenthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = CompareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorHTTPInIngressesDifferent) {
		t.Error("the HTTP in the ingresses is different'. But it was returned: ", err)
	}

	initEnvironmentForFifteenthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = CompareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorPathsCountDifferent) {
		t.Error("the paths count in the ingresses is different'. But it was returned: ", err)
	}

	initEnvironmentForSixteenthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = CompareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorPathValueDifferent) {
		t.Error("the path value in the ingresses is different'. But it was returned: ", err)
	}
}
