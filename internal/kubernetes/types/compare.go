package types

import (
	"context"
	"sync"
)

type KubeObjectsDifference struct {
	Msg string
	Final bool
}

type KubeObjectsDiffsStorage struct {
	m sync.RWMutex
	Diffs []KubeObjectsDifference
}

func NewKubeObjectsDiffsStorage () *KubeObjectsDiffsStorage {
	return &KubeObjectsDiffsStorage{
		Diffs: make([]KubeObjectsDifference, 0),
	}
}

func (s *KubeObjectsDiffsStorage) Add(msg string, final bool) bool {
	s.m.Lock()
	defer s.m.Unlock()

	s.Diffs = append(s.Diffs, KubeObjectsDifference{
		Msg:   msg,
		Final: final,
	})

	return final == true
}

type KubeResourceComparator interface {
	Compare(context.Context) ([]KubeObjectsDifference, error)

	//fieldSelectorProvider (context.Context) string
	//labelSelectorProvider (context.Context) string
	//
	//collect(ctx context.Context) (interface{}, error)
}
