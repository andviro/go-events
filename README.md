# go-events

Simple callback-based events modeled after C# events. Implementation is based on
brilliant [imkira/go-observer](https://github.com/imkira/go-observer) package. See its README for explanation why it's
the most efficient way to implement pubsub/observer/event system.

## Zero initialization

The main point of Event[T] object is to be usable without initialization and to
not require calling of closer method. Disposing of resources is handled by
client which must call closer function. For example usage, also see
[advanced example](https://github.com/andviro/go-events/blob/master/_example/advanced.go):

```go
// instantiate Event object
var evt events.Event[int]
cancel := evt.Handle(func(i int) {
	fmt.Println("received event:", i)
})
defer cancel() // Ensure that closer function is called

// ...

go func() {
	for i := 0; i < 10; i++ {
		evt.Invoke(i) // Send event
		time.Sleep(1 * time.Second)
	}
}()

// Handle events concurrently
time.Sleep(5 * time.Second)
```
