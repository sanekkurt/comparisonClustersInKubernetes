package diff

import (
	"context"
	"fmt"
	"go.uber.org/zap/zapcore"
	"k8s-cluster-comparator/internal/logging"
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

type DiffsBatch struct {
	//	m sync.RWMutex
	Diffs []ObjectsDiff
}

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
	//batch := make(DiffsBatch, 0)
	batch := DiffsBatch{
		Diffs: make([]ObjectsDiff, 0),
	}
	s.m.Lock()
	s.Batches = append(s.Batches, batch)
	defer s.m.Unlock()

	return &s.Batches[len(s.Batches)-1]
}

func (s *DiffsStorage) Add(ctx context.Context, objType *metav1.TypeMeta, objMeta *metav1.ObjectMeta, final bool, logLevel zapcore.Level, msg string, variables ...interface{}) bool {

	var (
		log = logging.FromContext(ctx)
	)

	switch logLevel {
	case zapcore.WarnLevel:
		log.Warnf(msg, variables...)
	case zapcore.ErrorLevel:
		log.Errorf(msg, variables...)
	case zapcore.FatalLevel:
		log.Fatalf(msg, variables...)
	case zapcore.PanicLevel:
		log.Panicf(msg, variables...)
	}

	diff := ObjectsDiff{
		Object: Object{},

		Msg:   msg,
		Final: final,
	}

	if objType != nil && objMeta != nil {
		diff.Object.Type = *objType
		diff.Object.Meta = *objMeta
	}

	diff.Msg = fmt.Sprintf(msg, variables...)
	diff.Final = final

	s.m.Lock()
	//	s.Batches[0] = append(s.Batches[0], diff)
	defer s.m.Unlock()

	return final == true
}

func (s *DiffsBatch) Add(ctx context.Context, objType *metav1.TypeMeta, objMeta *metav1.ObjectMeta, final bool, logLevel zapcore.Level, msg string, variables ...interface{}) bool {

	var (
		log = logging.FromContext(ctx)
	)

	switch logLevel {
	case zapcore.WarnLevel:
		log.Warnf(msg, variables...)
	case zapcore.ErrorLevel:
		log.Errorf(msg, variables...)
	case zapcore.FatalLevel:
		log.Fatalf(msg, variables...)
	case zapcore.PanicLevel:
		log.Panicf(msg, variables...)
	}

	diff := ObjectsDiff{
		Object: Object{},

		Msg:   msg,
		Final: final,
	}

	if objType != nil && objMeta != nil {
		diff.Object.Type = *objType
		diff.Object.Meta = *objMeta
	}

	diff.Msg = fmt.Sprintf(msg, variables...)
	diff.Final = final

	s.Diffs = append(s.Diffs, diff)

	return final == true
}
