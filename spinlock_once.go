package xsync

import (
	"sync/atomic"

	"github.com/xaionaro-go/spinlock"
)

type LazyInitSpinlock struct {
	init   func()
	inited int32
	locker spinlock.Locker
}

func (once *LazyInitSpinlock) InitedRun(fn func()) {
	if atomic.LoadInt32(&once.inited) != 0 {
		return
	}

	defer atomic.StoreInt32(&once.inited, 1)

	once.locker.Lock()
	defer once.locker.Unlock()

	fn()
}

func (once *LazyInitSpinlock) ImmutableRun(fn func()) {
	if atomic.LoadInt32(&once.inited) != 0 {
		fn()
		return
	}

	once.locker.Lock()
	defer once.locker.Unlock()

	fn()
}
