package worker

import (
	"context"
	"onepkg/util/log"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

var (
	taskPool   = sync.Pool{New: newTask}
	workerPool = sync.Pool{New: newWorker}
)

func newTask() interface{}   { return &task{} }
func newWorker() interface{} { return &worker{} }

type task struct {
	ctx context.Context
	f   func()

	next *task
}

func (t *task) recycle() {
	t.next = nil

	taskPool.Put(t)
}

type worker struct {
	pool *pool
}

func (w *worker) recycle() {
	w.pool = nil

	workerPool.Put(w)
}

func (w *worker) runTask(t *task) {
	defer func() {
		rec := recover()
		if rec != nil {
			log.Error(t.ctx, "GOPOOL: worker panic in pool: %s: %v: %s", w.pool.name, rec, debug.Stack())
			if w.pool.panicHandler != nil {
				w.pool.panicHandler(t.ctx, rec)
			}
		}
	}()

	t.f()
}

func (w *worker) run() {
	go func() {
		for {
			w.pool.taskLock.Lock()

			// if there's no task to do, exit
			if w.pool.taskHead == nil {
				// worker exit, count inc
				atomic.AddInt32(&w.pool.work, -1)
				w.pool.taskLock.Unlock()
				w.recycle()
				return
			}

			t := w.pool.taskHead
			w.pool.taskHead = w.pool.taskHead.next

			// task count inc
			atomic.AddInt32(&w.pool.task, -1)
			w.pool.taskLock.Unlock()

			w.runTask(t)

			// 归还资源
			t.recycle()
		}
	}()
}

// pool 业务池子, 目前没有支持丢弃策略, 默认无限续杯. 更好用处是对于需要限流, 限制 max qps 场景
type pool struct {
	// The name of the pool
	name string

	// capacity of the pool, the maximum number of goroutines that are actually working
	cap int32

	// Task returns the number of running tasks
	task int32

	// Record the number of running workers
	work int32

	// linked list of tasks
	taskHead *task
	taskTail *task
	taskLock sync.Mutex

	// This method will be called when the worker panic
	panicHandler func(ctx context.Context, rec interface{})
}

// Name returns the corresponding pool name.
func (p *pool) Name() string {
	return p.name
}

// Task returns the number of running tasks, Used to monitor operation status
func (p *pool) Task() int32 {
	return atomic.LoadInt32(&p.task)
}

// Work returns the number of running workers
func (p *pool) Work() int32 {
	return atomic.LoadInt32(&p.work)
}

// CtxGo executes f and accepts the context.
func (p *pool) Go(ctx context.Context, f func()) {
	// task pool get and init
	t := taskPool.Get().(*task)
	t.ctx = ctx
	t.f = f

	p.taskLock.Lock()
	if p.taskHead == nil {
		p.taskHead = t
		p.taskTail = t
	} else {
		p.taskTail.next = t
		p.taskTail = t
	}
	p.taskLock.Unlock()

	atomic.AddInt32(&p.task, 1)

	// 直到 go worker == pool cap 最大 worker 容量
	if p.Work() < p.cap {
		// worker add 1
		atomic.AddInt32(&p.work, 1)

		w := workerPool.Get().(*worker)
		w.pool = p
		w.run()
	}
}

// SetPanicHandler sets the panic handler rec = recover().
func (p *pool) SetPanicHandler(f func(ctx context.Context, rec interface{})) {
	p.panicHandler = f
}
