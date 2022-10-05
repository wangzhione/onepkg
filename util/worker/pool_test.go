package worker

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

const benchmarkTimes = 10000

func DoCopyStack(a, b int) int {
	if b < 100 {
		return DoCopyStack(0, b+1)
	}
	return 0
}

func testFunc() {
	_ = DoCopyStack(0, 0)
}

func testPanicFunc() {
	panic("test")
}

func TestPool(t *testing.T) {
	p := NewPool(100)

	var n int32
	var wg sync.WaitGroup
	for i := 0; i < 2000; i++ {
		wg.Add(1)
		p.Go(context.Background(), func() {
			defer wg.Done()
			atomic.AddInt32(&n, 1)
		})
	}
	wg.Wait()

	if n != 2000 {
		t.Error(n)
	}
}

func TestPoolPanic(t *testing.T) {
	p := NewPool(100)
	p.Go(context.Background(), testPanicFunc)
}

func BenchmarkPool(b *testing.B) {
	// runtime.GOMAXPROCS ( 逻辑 CPU 数量 )
	// < 1: 不修改任何数值
	// = 1: 单核心执行
	// > 1: 多核并发执行
	p := NewPool(int32(runtime.GOMAXPROCS(0)))

	var wg sync.WaitGroup

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(benchmarkTimes)
		for j := 0; j < benchmarkTimes; j++ {
			p.Go(context.Background(), func() {
				testFunc()
				wg.Done()
			})
		}
		wg.Wait()
	}

	/*
		Running tool: /opt/homebrew/bin/go test -benchmem -run=^$ -bench ^BenchmarkPool$ onepkg/util/worker -v

		goos: darwin
		goarch: arm64
		pkg: onepkg/util/worker
		BenchmarkPool
		BenchmarkPool-10    	     267	   4386748 ns/op	  202664 B/op	   10718 allocs/op
		PASS
		ok  	onepkg/util/worker	2.031s
	*/
}

func BenchmarkGo(b *testing.B) {
	var wg sync.WaitGroup

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.Add(benchmarkTimes)
		for j := 0; j < benchmarkTimes; j++ {
			go func() {
				testFunc()
				wg.Done()
			}()
		}
		wg.Wait()
	}

	/*
		Running tool: /opt/homebrew/bin/go test -benchmem -run=^$ -bench ^BenchmarkGo$ onepkg/util/worker -v

		goos: darwin
		goarch: arm64
		pkg: onepkg/util/worker
		BenchmarkGo
		BenchmarkGo-10    	     146	   9122579 ns/op	  167999 B/op	   10018 allocs/op
		PASS
		ok  	onepkg/util/worker	2.391s
	*/
}
