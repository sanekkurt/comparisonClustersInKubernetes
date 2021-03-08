package service

import (
	"errors"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

var (
	selector1 map[string]string
	selector2 map[string]string
)

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
	err := compareSpecInServices(*service1, *service2)
	if !errors.Is(errors.Unwrap(err), ErrorPortsCountDifferent) {
		t.Error("Error expected: 'the ports count are different'. But it was returned: ", err)
	}

	initEnvironmentForSecondTest3()
	service1, _ = clusterClientSet1.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	service2, _ = clusterClientSet2.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	err = compareSpecInServices(*service1, *service2)
	if !errors.Is(errors.Unwrap(err), ErrorPortInServicesDifferent) {
		t.Error("Error expected: 'the port in the services is different'. But it was returned: ", err)
	}

	initEnvironmemtForThirdTest3()
	service1, _ = clusterClientSet1.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	service2, _ = clusterClientSet2.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	err = compareSpecInServices(*service1, *service2)
	if !errors.Is(errors.Unwrap(err), ErrorSelectorsCountDifferent) {
		t.Error("Error expected: 'the selectors count are different'. But it was returned: ", err)
	}

	initEnvironmemtForFourthTest3()
	service1, _ = clusterClientSet1.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	service2, _ = clusterClientSet2.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	err = compareSpecInServices(*service1, *service2)
	if !errors.Is(errors.Unwrap(err), ErrorSelectorInServicesDifferent) {
		t.Error("Error expected: 'the selector in the services is different'. But it was returned: ", err)
	}

	initEnvironmemtForFifthTest3()
	service1, _ = clusterClientSet1.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	service2, _ = clusterClientSet2.CoreV1().Services("default").Get("testService", metav1.GetOptions{})
	err = compareSpecInServices(*service1, *service2)
	if !errors.Is(errors.Unwrap(err), ErrorTypeInServicesDifferent) {
		t.Error("the type in the services is different'. But it was returned: ", err)
	}
}
