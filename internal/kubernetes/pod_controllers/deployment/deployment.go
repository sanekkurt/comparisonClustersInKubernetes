package deployment

import (
	"context"
	"fmt"
	"k8s-cluster-comparator/internal/config"
	"k8s-cluster-comparator/internal/consts"
	kubectx "k8s-cluster-comparator/internal/kubernetes/context"
	"k8s-cluster-comparator/internal/kubernetes/diff"
	pccommon "k8s-cluster-comparator/internal/kubernetes/pod_controllers/common"
	"k8s-cluster-comparator/internal/kubernetes/types"
	"k8s-cluster-comparator/internal/logging"
	"sort"
	"sync"
	"time"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	objectKind = "deployment"
)

type Comparator struct {
	Kind      string
	Namespace string
	BatchSize int64
}

func NewComparator(ctx context.Context, namespace string) *Comparator {
	return &Comparator{
		Kind:      objectKind,
		Namespace: namespace,
		BatchSize: getBatchLimit(ctx),
	}
}

func (cmp *Comparator) FieldSelectorProvider(ctx context.Context) string {
	return ""
}

func (cmp *Comparator) LabelSelectorProvider(ctx context.Context) string {
	return ""
}

func (cmp *Comparator) collectIncludedFromCluster(ctx context.Context) (map[string]appsv1.Deployment, error) {
	var (
		log       = logging.FromContext(ctx)
		cfg       = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		objects = make(map[string]appsv1.Deployment)
	)

	log.Debugf("%T: collectIncludedFromCluster started", cmp)
	defer log.Debugf("%T: collectIncludedFromCluster completed", cmp)

	for name := range cfg.ExcludesIncludes.NameBasedSkip {
		obj, err := clientSet.AppsV1().Deployments(cmp.Namespace).Get(ctx, string(name), metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				log.With(zap.String("objectName", string(name))).Warnf("%s/%s not found in cluster", cmp.Kind, name)
				continue
			}
			return nil, err
		}
		objects[obj.Name] = *obj
	}

	for name := range cfg.ExcludesIncludes.FullResourceNamesSkip[types.ObjectKind(cmp.Kind)] {
		obj, err := clientSet.AppsV1().Deployments(cmp.Namespace).Get(ctx, string(name), metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				log.With(zap.String("objectName", string(name))).Warnf("%s/%s not found in cluster", cmp.Kind, name)
				continue
			}
			return nil, err
		}
		objects[obj.Name] = *obj
	}

	return objects, nil
}

func (cmp *Comparator) collectFromClusterWithoutExcludes(ctx context.Context) (map[string]appsv1.Deployment, error) {
	var (
		log       = logging.FromContext(ctx)
		cfg       = config.FromContext(ctx)
		clientSet = kubectx.ClientSetFromContext(ctx)

		batch   *appsv1.DeploymentList
		objects = make(map[string]appsv1.Deployment)

		continueToken string

		err error
	)

	log.Debugf("%T: collectFromClusterWithoutExcludes started", cmp)
	defer log.Debugf("%T: collectFromClusterWithoutExcludes completed", cmp)

forOuterLoop:
	for {
		select {
		case <-ctx.Done():
			return nil, context.Canceled
		default:
			batch, err = clientSet.AppsV1().Deployments(cmp.Namespace).List(ctx, metav1.ListOptions{
				Limit:         cmp.BatchSize,
				FieldSelector: cmp.FieldSelectorProvider(ctx),
				LabelSelector: cmp.LabelSelectorProvider(ctx),
				Continue:      continueToken,
			})
			if err != nil {
				return nil, err
			}

			log.Debugf("%d %s retrieved", len(batch.Items), cmp.Kind)

		forInnerLoop:
			for _, obj := range batch.Items {
				if _, ok := objects[obj.Name]; ok {
					log.With("objectName", obj.Name).Warnf("%s/%s already present in comparison list", cmp.Kind, obj.Name)
				}

				if cfg.ExcludesIncludes.IsSkippedEntity(cmp.Kind, obj.Name) {
					log.With(zap.String("objectName", obj.Name)).Debugf("%s/%s is skipped from comparison", cmp.Kind, obj.Name)
					continue forInnerLoop
				}

				if *obj.Spec.Replicas != obj.Status.ReadyReplicas {
					log.With(zap.String("objectName", obj.Name)).Warnf("%s/%s is progressing now, comparison might be inaccurate", cmp.Kind, obj.Name)
				}

				objects[obj.Name] = obj
			}

			if batch.Continue == "" {
				break forOuterLoop
			}

			continueToken = batch.Continue
		}
	}

	return objects, nil
}

func (cmp *Comparator) collectFromCluster(ctx context.Context) (map[string]appsv1.Deployment, error) {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)
	)

	log.Debugf("%T: collectFromCluster started", cmp)
	defer log.Debugf("%T: collectFromCluster completed", cmp)

	if cfg.Common.WorkMode == consts.EverythingButNotExcludesWorkMode {
		return cmp.collectFromClusterWithoutExcludes(ctx)
	} else {
		return cmp.collectIncludedFromCluster(ctx)
	}
}

// Compare compares list of Deployment objects in two given k8s-clusters
func (cmp *Comparator) Compare(ctx context.Context) (*diff.DiffsStorage, error) {
	var (
		log = logging.FromContext(ctx).With(zap.String("kind", cmp.Kind))
		cfg = config.FromContext(ctx)

		err error
	)

	ctx = logging.WithLogger(ctx, log)

	if !cfg.Workloads.Enabled ||
		!cfg.Workloads.PodControllers.Enabled ||
		!cfg.Workloads.PodControllers.Deployments.Enabled {
		log.Debugf("'%s' kind skipped from comparison due to configuration", cmp.Kind)
		return nil, nil
	}

	objects, err := cmp.collect(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve objects for comparision: %w", err)
	}

	err = cmp.compare(ctx, objects[0], objects[1])
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (cmp *Comparator) collect(ctx context.Context) ([]map[string]appsv1.Deployment, error) {
	var (
		log = logging.FromContext(ctx)
		cfg = config.FromContext(ctx)

		objects = make([]map[string]appsv1.Deployment, 2, 2)
		wg      = &sync.WaitGroup{}

		err error
	)

	wg.Add(2)

	for idx, clientSet := range []kubernetes.Interface{
		cfg.Connections.Cluster1.ClientSet,
		cfg.Connections.Cluster2.ClientSet,
	} {
		go func(idx int, clientSet kubernetes.Interface) {
			defer wg.Done()

			objects[idx], err = cmp.collectFromCluster(kubectx.WithClientSet(ctx, clientSet))
			if err != nil {
				log.Fatalf("cannot obtain %ss from cluster #%d: %s", cmp.Kind, idx+1, err.Error())
			}
		}(idx, clientSet)
	}

	wg.Wait()

	return objects, nil
}

func (cmp *Comparator) compare(ctx context.Context, map1, map2 map[string]appsv1.Deployment) error {
	var (
		apcs = make([]map[string]*pccommon.AbstractPodController, 2, 2)
		cfg  = config.FromContext(ctx)
	)

	for idx, objs := range []map[string]appsv1.Deployment{map1, map2} {
		apcs[idx] = cmp.prepareAPCMap(ctx, objs)
	}

	if cfg.Common.CheckingCreationTimestampDeploymentsLimit {
		clearApcs, err := cmp.prepareReplicaSets(ctx, apcs)
		if err != nil {
			return err
		}

		err = pccommon.CompareAbstractPodControllerMaps(kubectx.WithNamespace(ctx, cmp.Namespace), cmp.Kind, clearApcs[0], clearApcs[1])
		if err != nil {
			return err
		}

		return nil

	} else {
		err := pccommon.CompareAbstractPodControllerMaps(kubectx.WithNamespace(ctx, cmp.Namespace), cmp.Kind, apcs[0], apcs[1])
		if err != nil {
			return err
		}

		return nil
	}

}

func (cmp *Comparator) prepareAPCMap(ctx context.Context, objs map[string]appsv1.Deployment) map[string]*pccommon.AbstractPodController {
	var (
		apcs = make(map[string]*pccommon.AbstractPodController)
	)

	for name, obj := range objs {
		apcs[name] = &pccommon.AbstractPodController{
			Metadata: types.AbstractObjectMetadata{
				Type: metav1.TypeMeta{
					Kind:       cmp.Kind,
					APIVersion: "apps/v1",
				},
				Meta: obj.ObjectMeta,
			},
			Name:             obj.Name,
			Labels:           obj.Labels,
			Annotations:      obj.Annotations,
			Replicas:         obj.Spec.Replicas,
			PodLabelSelector: obj.Spec.Selector,
			PodTemplateSpec:  obj.Spec.Template,
		}
	}

	return apcs
}

func (cmp *Comparator) prepareReplicaSets(ctx context.Context, apcs []map[string]*pccommon.AbstractPodController) ([]map[string]*pccommon.AbstractPodController, error) {
	var (
		cfg           = config.FromContext(ctx)
		log           = logging.FromContext(ctx)
		continueToken string
		limitTime     = -time.Duration(cfg.Workloads.PodControllers.Deployments.DiscardDeploymentsUpdatedLaterTime) * time.Minute
	)

	select {
	case <-ctx.Done():
		return nil, context.Canceled

	default:
		for _, apc := range apcs {
			for key, value := range apc {

				matchLabels, _ := metav1.LabelSelectorAsSelector(value.PodLabelSelector)
				var replicaSetList []appsv1.ReplicaSet

			Loop:
				for {
					batch, err := cfg.Connections.Cluster1.ClientSet.AppsV1().ReplicaSets(cmp.Namespace).List(ctx, metav1.ListOptions{
						Limit:         cmp.BatchSize,
						LabelSelector: matchLabels.String(),
						Continue:      continueToken,
					})
					if err != nil {
						return nil, err
					}

					log.Debugf("%d %s for %s deployment retrieved", len(batch.Items), "ReplicaSets", key)

					for _, replicaSet := range batch.Items {
						replicaSetList = append(replicaSetList, replicaSet)
					}

					if batch.Continue == "" {
						break Loop
					}

					continueToken = batch.Continue
				}

				sort.SliceStable(replicaSetList, func(i, j int) bool {
					return replicaSetList[i].CreationTimestamp.Time.After(replicaSetList[j].CreationTimestamp.Time)
				})

				if replicaSetList[0].CreationTimestamp.Time.After(time.Now().Add(limitTime)) {
					log.Warnf("latest replicaSet have creationTimestamp '%s' after limit '%s'. The deployment %s will be excluded", replicaSetList[0].CreationTimestamp.Time, limitTime, key)
					delete(apc, key)
				}
			}
		}

		return apcs, nil

	}

}
