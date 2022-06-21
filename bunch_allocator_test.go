package xsync

import (
	"fmt"
	"runtime"
	"testing"
)

type testType1 [1]byte

type testType8 [8]byte

type testType32 [32]byte

type testType256 [256]byte

type testType1024 [1024]byte

type testType65536 [65536]byte

var (
	v1     *testType1
	v8     *testType8
	v32    *testType32
	v256   *testType256
	v1024  *testType1024
	v65536 *testType65536
)

func BenchmarkBunchAllocator(b *testing.B) {
	for _, structSize := range []uint{1, 8, 32, 256, 1024, 65536} {
		b.Run(fmt.Sprintf("structSize%d", structSize), func(b *testing.B) {

			var (
				noBunchGetter func()
				bunchGetter   func()
			)
			switch structSize {
			case 1:
				noBunchGetter = func() {
					v1 = &testType1{}
				}
				allocator := NewBunchAllocator[testType1]()
				bunchGetter = func() {
					v1 = allocator.Get()
				}
			case 8:
				noBunchGetter = func() {
					v8 = &testType8{}
				}
				allocator := NewBunchAllocator[testType8]()
				bunchGetter = func() {
					v8 = allocator.Get()
				}
			case 32:
				noBunchGetter = func() {
					v32 = &testType32{}
				}
				allocator := NewBunchAllocator[testType32]()
				bunchGetter = func() {
					v32 = allocator.Get()
				}
			case 256:
				noBunchGetter = func() {
					v256 = &testType256{}
				}
				allocator := NewBunchAllocator[testType256]()
				bunchGetter = func() {
					v256 = allocator.Get()
				}
			case 1024:
				noBunchGetter = func() {
					v1024 = &testType1024{}
				}
				allocator := NewBunchAllocator[testType1024]()
				bunchGetter = func() {
					v1024 = allocator.Get()
				}
			case 65536:
				noBunchGetter = func() {
					v65536 = &testType65536{}
				}
				allocator := NewBunchAllocator[testType65536]()
				bunchGetter = func() {
					v65536 = allocator.Get()
				}
			}
			b.Run("no_bunch", func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					noBunchGetter()
				}
				runtime.GC()
				runtime.GC()
			})

			b.Run("bunch", func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					bunchGetter()
				}
				runtime.GC()
				runtime.GC()
			})
		})
	}
}
