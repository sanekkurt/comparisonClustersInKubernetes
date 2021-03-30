package diff

import (
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Object struct {
	Type metav1.TypeMeta
	Meta metav1.ObjectMeta
}

type ObjectsDiff struct {
	Object Object

	Msg   string
	Final bool
}

type DiffsBatch []ObjectsDiff

type DiffsStorage struct {
	m       sync.RWMutex
	Batches []DiffsBatch
}

func NewDiffsStorage() *DiffsStorage {
	return &DiffsStorage{
		Batches: make([]DiffsBatch, 0),
	}
}

func (s *DiffsStorage) NewBatch() *DiffsBatch {
	batch := make(DiffsBatch, 0)
	s.Batches = append(s.Batches, batch)

	return &s.Batches[len(s.Batches)-1]
}

func (s *DiffsStorage) Add(objType *metav1.TypeMeta, objMeta *metav1.ObjectMeta, msg string, final bool) bool {
	diff := ObjectsDiff{
		Object: Object{},

		Msg:   msg,
		Final: final,
	}

	if objType != nil && objMeta != nil {
		diff.Object.Type = *objType
		diff.Object.Meta = *objMeta
	}

	s.m.Lock()
	defer s.m.Unlock()

	return final == true
}
