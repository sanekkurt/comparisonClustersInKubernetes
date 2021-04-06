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

func NewDiffsStorage(ctx context.Context) *DiffsStorage {
	ds := &DiffsStorage{
		wg: &sync.WaitGroup{},
	}

	return ds
}

func (s *DiffsStorage) Finalize(ctx context.Context) {
	var (
		log = logging.FromContext(ctx)
	)

	log.Debugf("Closing %d collecting goroutines...", len(s.batches))

	for _, batch := range s.batches {
		close(batch.diffsCh)
	}

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
