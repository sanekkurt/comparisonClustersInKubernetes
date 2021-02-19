package types

import (
	"sync"
)

type OnceSettableFlag struct {
	m    sync.Mutex
	flag bool
}

func (o *OnceSettableFlag) SetFlag(v bool) {
	//o.m.Lock()
	//defer func() {
	//	o.m.Unlock()
	//}()

	if o.flag {
		return
	}

	o.flag = true
}

func (o *OnceSettableFlag) GetFlag() bool {
	//o.m.Lock()
	//defer func() {
	//	o.m.Unlock()
	//}()

	return o.flag
}
