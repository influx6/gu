Notifications
=============

In Gu there exists a central notification backbone package `notifications`, which exposes a system that allows registering specific functions of specific types of structures to be called when such structures are dispatched into the system to allow a decoupled form of communication.

*This provide loose coupling between components as is needed.*

Using the `notifications` package is simple. By simply registering a function expecting a type, this sets up this function to be called once such type is seen.

```go

import "github.com/gu-io/gu/notifications"

type event struct{
  EventName string
  EventType string
}


func main(){

  notifications.Subscribe(func(eventName interface{}){
    fmt.Printf("EventName[%+q] occured.\n", eventName)
  })

  notifications.Dispatch("Click") => `EventName["Click"] occured.`
}
```

## Custom Notification

Include in the Gu library is a code generation system which allows you to annotate
a given struct type to be an event, which sets of functions and structures should be
generated for.

We equally understand of the importance of lazy developers, as we are one ourselves, hence
this provides us a quick and seamless way to plug into the central notification system, whilst
ensuring to keep type safety by generating the needed code to convert the interface to the
expected type, before notifying the provided function or subscribers.

By annotating structures with `@notification:event` and with a call to `gu generate`,
any structures which has such annotations will have the event handling and assertion
strucutures generated for it.

```go

//@notification:event
type EventForward struct{
  X int
  Angle float64
}

```

See example usage in core:

- AppEvent
    Annotation: https://github.com/gu-io/gu/blob/master/notifications/notifiers.go#L6
    Generated: https://github.com/gu-io/gu/blob/master/notifications/appevent_event.go

- ViewUpdate
    Annotation: https://github.com/gu-io/gu/blob/master/gu.go#L97
    Generated: https://github.com/gu-io/gu/blob/master/viewupdate_event.go
