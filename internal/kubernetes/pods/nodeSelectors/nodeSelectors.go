package nodeSelectors

import (
	"context"
	"k8s-cluster-comparator/internal/logging"
)

func CompareNodeSelectors(ctx context.Context, nodeSelector1, nodeSelector2 map[string]string) {

	var (
		log = logging.FromContext(ctx)
		//diffs = make([]types.KubeObjectsDifference, 0)
	)

	for key1, value1 := range nodeSelector1 {

		flag := false
		for key2, value2 := range nodeSelector2 {

			if value1 == value2 && key1 == key2 {

				flag = true
				delete(nodeSelector1, key1)
				delete(nodeSelector2, key1)
				break

			} else if key1 == key2 {

				log.Warnf("%s. %s-%s vs %s-%s", ErrorDiffersNodeSelectorsInTemplates.Error(), key1, value1, key2, value2)
				flag = true
				delete(nodeSelector1, key1)
				delete(nodeSelector2, key1)
				break
			}
		}
		if !flag {
			log.Warnf("node selector %s-%s does not exist in second cluster pod template", key1, value1)
			delete(nodeSelector1, key1)
		}
	}

	if len(nodeSelector2) != 0 {
		for key, value := range nodeSelector2 {
			log.Warnf("node selector %s-%s does not exist in first cluster pod template", key, value)
			delete(nodeSelector2, key)
		}
	}
}