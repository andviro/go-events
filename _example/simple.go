package main

import (
	"fmt"
	"time"

	"github.com/andviro/go-events"
)

func main() {
	var evt events.Event[int]
	cancel := evt.Handle(func(i int) {
		fmt.Println("received event:", i)
	})
	defer cancel()
	go func() {
		for i := 0; i < 10; i++ {
			evt.Invoke(i)
			time.Sleep(1 * time.Second)
		}
	}()
	time.Sleep(5 * time.Second)
}
