package diff

import (
	"context"
	"fmt"
	"sync"

	"k8s-cluster-comparator/internal/logging"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type object struct {
	Type metav1.TypeMeta
	Meta metav1.ObjectMeta
}

type DiffsStorage struct {
	wg *sync.WaitGroup

	batches []*DiffsBatch
}

type objectsDiff struct {
	Msg   string
	Final bool
}

type DiffsBatch struct {
	object object
	diffs  []objectsDiff

	diffsCh chan objectsDiff

	s    *DiffsStorage
	once sync.Once
}

//type DiffsBatch struct {
//	object object
//
//	diffs []objectsDiff
//}

func NewDiffsStorage(ctx context.Context) *DiffsStorage {
	var (
		log    = logging.FromContext(ctx)
		syncCh = make(chan struct{})
	)

	ds := &DiffsStorage{
		wg: &sync.WaitGroup{},
		//batches: make([]DiffsBatch, 0),
	}

	go func(ds *DiffsStorage) {
		log.Infof("[SYNC BEGIN]")

		syncCh <- struct{}{}

		log.Infof("[WORKER COMPLETED]")
	}(ds)

	<-syncCh

	log.Infof("[SYNC DONE]")

	return ds
}

func (s *DiffsStorage) Finalize() {
	s.wg.Wait()
}

func (s *DiffsStorage) NewLazyBatch(objType metav1.TypeMeta, objMeta metav1.ObjectMeta) *DiffsBatch {
	b := &DiffsBatch{
		object: object{
			Type: objType,
			Meta: objMeta,
		},

		diffsCh: make(chan objectsDiff),

		s: s,
	}

	return b
}

//
//type Diff struct {
//	Ctx       context.Context
//	Final     bool
//	LogLevel  zapcore.Level
//	Msg       string
//	Variables []interface{}
//}

//type ChanForDiff chan Diff
//
//func (s *DiffsStorage) NewChannel(objType metav1.TypeMeta, objMeta metav1.ObjectMeta) *ChanForDiff {
//
//	ch := make(ChanForDiff)
//	go cyclicReadingFromChannel(ch, s, objType, objMeta)
//
//	return &ch
//}

//func cyclicReadingFromChannel(c ChanForDiff, difStorage *DiffsStorage, objType metav1.TypeMeta, objMeta metav1.ObjectMeta) {
//	var batchCreated bool
//	var dif *DiffsBatch
//	for {
//		val, ok := <-c
//		if ok {
//			if !batchCreated {
//				dif = difStorage.NewLazyBatch(objType, objMeta)
//				batchCreated = true
//			}
//
//			dif.Add(val.Ctx, val.Final, val.LogLevel, val.Msg, val.Variables...)
//		} else {
//			break
//		}
//	}
//}

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
//	diff := objectsDiff{
//		Msg:   msg,
//		Final: final,
//	}
//
//	diff.Msg = fmt.Sprintf(msg, variables...)
//	diff.Final = final
//
//	s.m.Lock()
//	//	s.batches[0] = append(s.batches[0], diff)
//	defer s.m.Unlock()
//
//	return final == true
//}

func (b *DiffsBatch) lazyInit(ctx context.Context) {
	var (
		log = logging.FromContext(ctx)
	)

	b.once.Do(func() {
		log.Debug("New diffsBatch created")

		b.s.wg.Add(1)
		b.s.batches = append(b.s.batches, b)

		// since we are here, there is one diff at least
		b.diffs = make([]objectsDiff, 0, 1)

		go func(b *DiffsBatch) {
			log.Debug("New diffsBatch channel reader created")

			for diff := range b.diffsCh {
				b.diffs = append(b.diffs, diff)
			}

			b.s.wg.Done()
		}(b)
	})
}

func (b *DiffsBatch) Add(ctx context.Context, final bool, msg string, fields ...interface{}) bool {
	b.lazyInit(ctx)

	diffLog(ctx, b.object, msg, fields...)

	diff := objectsDiff{
		Msg:   fmt.Sprintf(msg, fields...),
		Final: final,
	}

	b.diffsCh <- diff

	return final
}
