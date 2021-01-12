package common

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ErrorObjectCacheNoSuchKindCache = errors.New("there is no such object kind cache yet")
	ErrorObjectCacheNoSuchObject    = errors.New("object not cached yet")
)

type objectsCache map[string]interface{}
type kindsCache map[string]objectsCache

type ObjectCache struct {
	m sync.RWMutex
	// [objectKind][objectName]Object
	cache kindsCache
}

func (c *ObjectCache) getKey(meta v1.TypeMeta) string {
	return strings.ToLower(fmt.Sprintf("%s/%s", meta.APIVersion, meta.Kind))
}

func (c *ObjectCache) Get(meta v1.TypeMeta, name string) (interface{}, error) {
	c.m.RLock()
	defer func() {
		c.m.RUnlock()
	}()

	k := c.getKey(meta)

	_, ok := c.cache[k]
	if !ok {
		return nil, ErrorObjectCacheNoSuchKindCache
	}

	obj, ok := c.cache[k][name]
	if !ok {
		return nil, ErrorObjectCacheNoSuchObject
	}

	return obj, nil
}

func (c *ObjectCache) Put(meta v1.TypeMeta, name string, obj interface{}) {
	c.m.Lock()
	defer func() {
		c.m.Unlock()
	}()

	k := c.getKey(meta)

	_, ok := c.cache[k]
	if !ok {
		c.cache[k] = make(objectsCache)
	}

	c.cache[k][name] = obj
}
