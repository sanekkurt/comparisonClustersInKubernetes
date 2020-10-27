# k8s-cluster-comparator :: a tiny Go tool to compare two kube-clusters

> k8s-cluster-comparator allow comparing two kube-clusters in terms of resource presents, 
> pod specifications (container images, environment, ...) and data equality in ConfigMap / Secret objects

## What is compared

* configuration storage KV-maps
    * ConfigMaps
    * Secrets (special secrets (service account tokens, Helm release metadata, etc) are excluded)


* pod controllers
    * Deployments
    * StatefulSets
    * DaemonSets
    

* one-hop pod-controllers
    * Jobs
    * CronJobs
    
* network-related resources
    * Services
    * Ingresses
    
## How to use

*Coming Soon*