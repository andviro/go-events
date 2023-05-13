package events

import (
	"sync"

	"github.com/imkira/go-observer/v2"
)

type Event[T any] struct {
	prop observer.Property[T]
	once sync.Once
	mu   sync.RWMutex
}

func (e *Event[T]) Handle(h func(T)) (unhandle func()) {
	e.once.Do(func() {
		e.mu.Lock()
		defer e.mu.Unlock()
		var t T
		e.prop = observer.NewProperty(t)
	})
	stop := make(chan struct{})
	stream := e.prop.Observe()
	go func() {
		for {
			select {
			case <-stop:
				return
			case <-stream.Changes():
				h(stream.Next())
			}
		}
	}()
	return func() {
		close(stop)
	}
}

func (e *Event[T]) Invoke(val T) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.prop == nil {
		return
	}
	e.prop.Update(val)
}
