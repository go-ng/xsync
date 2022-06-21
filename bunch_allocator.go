package xsync

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	bunchSize = 256
)

type bunchAllocatorPool[T any] struct {
	storage [bunchSize]T
	ptr     uintptr
	end     uintptr
}

// BunchAllocator allocates a bunch of objects and then returns them one by one.
// This supposed to reduce overhead on memory management processes, but currently
// it works well only for objects of small size (around 8-256 bytes). If the structure
// is large, then it only adds overhead.
//
// Latest benchmark results:
//
// name                               old time/op    new time/op    delta
// BunchAllocator/structSize1-16        7.62ns ± 5%    7.34ns ± 4%      ~     (p=0.310 n=5+5)
// BunchAllocator/structSize8-16        10.5ns ± 2%     8.1ns ±15%   -22.84%  (p=0.008 n=5+5)
// BunchAllocator/structSize32-16       17.0ns ± 7%    11.9ns ± 5%   -29.60%  (p=0.008 n=5+5)
// BunchAllocator/structSize256-16      44.7ns ± 3%    33.1ns ±14%   -26.06%  (p=0.008 n=5+5)
// BunchAllocator/structSize1024-16      156ns ±11%     190ns ± 5%   +21.78%  (p=0.008 n=5+5)
// BunchAllocator/structSize65536-16    4.35µs ±27%    5.41µs ± 5%      ~     (p=0.095 n=5+5)
//
// name                               old alloc/op   new alloc/op   delta
// BunchAllocator/structSize8-16         8.00B ± 0%     9.00B ± 0%   +12.50%  (p=0.008 n=5+5)
// BunchAllocator/structSize32-16        32.0B ± 0%     37.0B ± 0%   +15.62%  (p=0.008 n=5+5)
// BunchAllocator/structSize256-16        256B ± 0%      288B ± 0%   +12.50%  (p=0.008 n=5+5)
// BunchAllocator/structSize1024-16     1.02kB ± 0%    1.06kB ± 0%    +3.09%  (p=0.008 n=5+5)
// BunchAllocator/structSize65536-16    65.5kB ± 0%    65.6kB ± 0%    +0.05%  (p=0.008 n=5+5)
type BunchAllocator[T any] struct {
	*bunchAllocatorPool[T]
	ptrDelta   uintptr
	poolLocker sync.Mutex
}

// NewBunchAllocator returns a new instance of BunchAllocator
func NewBunchAllocator[T any]() *BunchAllocator[T] {
	a := &BunchAllocator[T]{}
	a.ptrDelta = unsafe.Sizeof(a.storage[0])
	a.setNewPool()
	return a
}

// Get returns a new zero-valued instance of type T
func (a *BunchAllocator[T]) Get() *T {
	for {
		if obj := a.pool().get(a.ptrDelta); obj != nil {
			return obj
		}
		a.updatePool()
	}
}

func (a *BunchAllocator[T]) updatePool() {
	a.poolLocker.Lock()
	defer a.poolLocker.Unlock()

	if atomic.LoadUintptr((*uintptr)((unsafe.Pointer)(&a.pool().ptr))) < a.end {
		return
	}

	a.setNewPool()
}

func (a *BunchAllocator[T]) setNewPool() {
	newPool := &bunchAllocatorPool[T]{}
	newPool.ptr = (uintptr)(unsafe.Pointer(&newPool.storage[0]))
	newPool.end = newPool.ptr + a.ptrDelta*bunchSize
	atomic.StorePointer((*unsafe.Pointer)((unsafe.Pointer)(&a.bunchAllocatorPool)), (unsafe.Pointer)(newPool))
}

func (a *BunchAllocator[T]) pool() *bunchAllocatorPool[T] {
	return (*bunchAllocatorPool[T])(atomic.LoadPointer((*unsafe.Pointer)((unsafe.Pointer)(&a.bunchAllocatorPool))))
}

func (p *bunchAllocatorPool[T]) get(ptrDelta uintptr) *T {
	ptr := (unsafe.Pointer)(atomic.AddUintptr(&p.ptr, ptrDelta))
	obj := unsafe.Add(ptr, -ptrDelta)
	if uintptr(obj) < p.end {
		return (*T)(obj)
	}
	return nil
}
