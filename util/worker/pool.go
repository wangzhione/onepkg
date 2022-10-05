package worker

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// Pool go pool
type Pool interface {
	// Go executes f and accepts the context.
	Go(ctx context.Context, f func())

	// Name returns the corresponding pool name.
	Name() string
	// Task returns the number of running tasks, Used to monitor operation status
	Task() int32
	// Work returns the number of running workers
	Work() int32

	// SetPanicHandler sets the panic handler rec = recover().
	SetPanicHandler(f func(ctx context.Context, rec interface{}))
}

// NewPool creates a new pool with the given name, cap and config.
func NewPool(threshold int32) Pool {
	return &pool{
		name: uuid.New().String(),
		cap:  threshold,
	}
}

var poolMap sync.Map

// RegisterPool registers a new pool to the global map.
// GetPool can be used to get the registered pool by name.
// returns error if the same name is registered.
func RegisterPool(p Pool) error {
	_, loaded := poolMap.LoadOrStore(p.Name(), p)
	if loaded {
		return fmt.Errorf("name: %s already registered", p.Name())
	}
	return nil
}

// GetPool gets the registered pool by name.
// Returns nil if not registered.
func GetPool(name string) Pool {
	p, ok := poolMap.Load(name)
	if !ok {
		return nil
	}
	return p.(Pool)
}
