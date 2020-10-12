package networking

import (
	"errors"
	"testing"

	"k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/fake"
)

var (
	// variables for the container comparison function test
	clusterClientSet1 *fake.Clientset
	clusterClientSet2 *fake.Clientset
)

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
		Spec: v1beta1.IngressSpec{},
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
		Host: "host",
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
	err := compareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorTLSInIngressesDifferent) {
		t.Error("the TLS in the ingresses are different'. But it was returned: ", err)
	}

	initEnvironmentForSecondTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = compareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorTLSCountDifferent) {
		t.Error("the TLS count in the ingresses are different'. But it was returned: ", err)
	}

	initEnvironmentForThirdTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = compareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorSecretNameInTLSDifferent) {
		t.Error("the secret name in the TLS are different'. But it was returned: ", err)
	}

	initEnvironmentForFifthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = compareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorHostsCountDifferent) {
		t.Error("the hosts count in the TLS are different'. But it was returned: ", err)
	}

	initEnvironmentForSixthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = compareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorNameHostDifferent) {
		t.Error("the name host in the TLS are different'. But it was returned: ", err)
	}

	initEnvironmentForSeventhTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = compareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorHostsInIngressesDifferent) {
		t.Error("the hosts in the ingresses are different'. But it was returned: ", err)
	}

	initEnvironmentForEighthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = compareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorBackendInIngressesDifferent) {
		t.Error("the backend in the ingresses are different'. But it was returned: ", err)
	}

	initEnvironmentForNinthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = compareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorServiceNameInBackendDifferent) {
		t.Error("the service name in the backend are different'. But it was returned: ", err)
	}

	initEnvironmentForTenthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = compareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorBackendServicePortDifferent) {
		t.Error("the service port in the backend are different'. But it was returned: ", err)
	}

	initEnvironmentForEleventhTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = compareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorRulesInIngressesDifferent) {
		t.Error("the rules in the ingresses are different'. But it was returned: ", err)
	}

	initEnvironmentForTwelvesTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = compareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorRulesCountDifferent) {
		t.Error("the rules count in the ingresses is different'. But it was returned: ", err)
	}

	initEnvironmentForThirteenthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = compareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorHostNameInRuleDifferent) {
		t.Error("the hosts name in the rule are different'. But it was returned: ", err)
	}

	initEnvironmentForFourteenthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = compareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorHTTPInIngressesDifferent) {
		t.Error("the HTTP in the ingresses is different'. But it was returned: ", err)
	}

	initEnvironmentForFifteenthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = compareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorPathsCountDifferent) {
		t.Error("the paths count in the ingresses is different'. But it was returned: ", err)
	}

	initEnvironmentForSixteenthTest4()
	ingress1, _ = clusterClientSet1.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	ingress2, _ = clusterClientSet2.NetworkingV1beta1().Ingresses("default").Get("testIngress", metav1.GetOptions{})
	err = compareSpecInIngresses(*ingress1, *ingress2)
	if !errors.Is(errors.Unwrap(err), ErrorPathValueDifferent) {
		t.Error("the path value in the ingresses is different'. But it was returned: ", err)
	}
}
