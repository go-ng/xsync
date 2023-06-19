package xsync

import "sync"

// Pool is a type-asserted variant of sync.Pool.
type Pool[T any] struct {
	sync.Pool
	Reset func(*T)
}

// Get analogous to (*sync.Pool).Get, but already type asserted.
func (p *Pool[T]) Get() T {
	return p.Pool.Get().(T)
}

// Put analogous to (*sync.Pool).Put, but already type asserted,
// and with Reset function called.
func (p *Pool[T]) Put(x T) {
	if p.Reset != nil {
		p.Reset(&x)
	}
	p.Pool.Put(x)
}

// PoolR is a wrapper for a Pool to wrap/unwrap automatically
// the Releasable structure.
type PoolR[T any] struct {
	Pool Pool[T]
}

// Get analogous to (*sync.Pool).Get, but already type asserted.
func (p *PoolR[T]) Get() Releasable[T] {
	return Releasable[T]{
		Value: p.Pool.Get(),
		Pool:  &p.Pool,
	}
}

// NewPoolR returns a new Pool which returns value with method
// Release to put then back to the Pool (after use).
func NewPoolR[T any](
	initFunc func(*T),
	resetFunc func(*T),
) *PoolR[T] {
	p := &PoolR[T]{}
	p.Pool.New = func() any {
		var r T
		if initFunc != nil {
			initFunc(&r)
		}
		return r
	}
	if resetFunc != nil {
		p.Pool.Reset = func(in *T) {
			resetFunc(in)
		}
	}
	return p
}

// Releasable is a wrapper for a value to make it releasable
// pack to its Pool.
type Releasable[T any] struct {
	// TODO: make the field embedded when
	// `type Value[T any] = T` will be permitted (Go1.22?).
	//
	// Partially related ticket: https://github.com/golang/go/issues/46477
	Value T
	Pool  *Pool[T]
}

// Release puts the value back to the Pool.
func (v Releasable[T]) Release() {
	v.Pool.Put(v.Value)
}
