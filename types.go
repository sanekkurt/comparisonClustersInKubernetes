package main

import (
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CheckerFlag struct {
	index int
	check bool
}

//структура для описания основной информации в контейнере для сравнения
type Container struct {
	name string
	image string
	imageID string
}

//структура для универсализации сравнительной функции, позволяет в нее передать информацию как deployment'ов, так и statefulset'ов
type InformationAboutObject struct {
	Template v12.PodTemplateSpec
	Selector *v1.LabelSelector
	//Pods *v12.PodList
}

type KubeconfigYaml struct {
	APIVersion string `json:"apiVersion"`
	Clusters   []struct {
		Cluster struct {
			CertificateAuthorityData string `json:"certificate-authority-data"`
			Server                   string `json:"server"`
		} `json:"cluster"`
		Name string `json:"name"`
	} `json:"clusters"`
	Contexts []struct {
		Context struct {
			Cluster string `json:"cluster"`
			User    string `json:"user"`
		} `json:"context"`
		Name string `json:"name"`
	} `json:"contexts"`
	CurrentContext string `json:"current-context"`
	Kind           string `json:"kind"`
	Preferences    struct {
	} `json:"preferences"`
	Users []struct {
		Name string `json:"name"`
		User struct {
			ClientCertificateData string `json:"client-certificate-data"`
			ClientKeyData         string `json:"client-key-data"`
		} `json:"user"`
	} `json:"users"`
}