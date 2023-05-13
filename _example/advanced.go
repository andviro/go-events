package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/andviro/go-events"
)

type EventSource struct {
	onInt    events.Event[int]
	onString events.Event[string]
	closers  []func()
}

type Handler[T any] func(ctx context.Context, arg T) error

// _subscribe decorates simple event with some error handling, context etc.
func _subscribe[T any](dest *EventSource, event *events.Event[T], handler Handler[T]) {
	ctx := context.TODO()
	cancel := event.Handle(func(arg T) {
		if err := handler(ctx, arg); err != nil {
			log.Printf("WARN: error in '%T' handler: %+v", arg, err)
		}
	})
	dest.closers = append(dest.closers, cancel)
}

func (e *EventSource) Close() {
	for _, c := range e.closers {
		c()
	}
}

func (e *EventSource) OnInt(h Handler[int]) {
	_subscribe(e, &e.onInt, h)
}

func (e *EventSource) OnString(h Handler[string]) {
	_subscribe(e, &e.onString, h)
}

func (e *EventSource) Run() {
	// generate some events
	for i := 0; i < 5; i++ {
		time.Sleep(500 * time.Millisecond)
		e.onInt.Invoke(i)
		time.Sleep(500 * time.Millisecond)
		e.onString.Invoke(strconv.Itoa(i))
	}
}

func main() {
	var e EventSource
	defer e.Close()
	e.OnInt(func(ctx context.Context, arg int) error {
		if arg == 3 {
			return fmt.Errorf("invalid value: %d", arg)
		}
		log.Printf("INFO: got %T event: %+v", arg, arg)
		return nil
	})
	e.OnString(func(ctx context.Context, arg string) error {
		if arg == "2" {
			return fmt.Errorf("bad string %q", arg)
		}
		log.Printf("INFO: got %T event: %+v", arg, arg)
		return nil
	})
	e.Run()
}
