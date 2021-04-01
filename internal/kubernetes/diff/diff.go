package diff

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap/zapcore"
	"k8s-cluster-comparator/internal/logging"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Object struct {
	Type metav1.TypeMeta
	Meta metav1.ObjectMeta
}

type ObjectsDiff struct {
	Msg   string
	Final bool
}

type DiffsBatch struct {
	m sync.RWMutex

	Object Object
	Diffs  []ObjectsDiff
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

func (s *DiffsStorage) NewBatch(objType metav1.TypeMeta, objMeta metav1.ObjectMeta) *DiffsBatch {
	//batch := make(DiffsBatch, 0)
	batch := DiffsBatch{
		Object: Object{
			Type: objType,
			Meta: objMeta,
		},
		Diffs: make([]ObjectsDiff, 0),
	}

	s.m.Lock()
	defer s.m.Unlock()

	s.Batches = append(s.Batches, batch)

	return &s.Batches[len(s.Batches)-1]
}

//func (s *DiffsStorage) Add(ctx context.Context, objType *metav1.TypeMeta, objMeta *metav1.ObjectMeta, final bool, logLevel zapcore.Level, msg string, variables ...interface{}) bool {
//
//	var (
//		log = logging.FromContext(ctx)
//	)
//
//	switch logLevel {
//	case zapcore.WarnLevel:
//		log.Warnf(msg, variables...)
//	case zapcore.ErrorLevel:
//		log.Errorf(msg, variables...)
//	case zapcore.FatalLevel:
//		log.Fatalf(msg, variables...)
//	case zapcore.PanicLevel:
//		log.Panicf(msg, variables...)
//	}
//
//	diff := ObjectsDiff{
//		Msg:   msg,
//		Final: final,
//	}
//
//	diff.Msg = fmt.Sprintf(msg, variables...)
//	diff.Final = final
//
//	s.m.Lock()
//	//	s.Batches[0] = append(s.Batches[0], diff)
//	defer s.m.Unlock()
//
//	return final == true
//}

func (s *DiffsBatch) Add(ctx context.Context, final bool, logLevel zapcore.Level, msg string, variables ...interface{}) bool {
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
		Msg:   msg,
		Final: final,
	}

	diff.Msg = fmt.Sprintf(msg, variables...)
	diff.Final = final

	s.m.Lock()
	defer s.m.Unlock()

	s.Diffs = append(s.Diffs, diff)

	return final == true
}
