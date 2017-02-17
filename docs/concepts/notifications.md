Notifications
=============

In Gu there exists a central notification backbone package `notifications`, which exposes a system that allows registering specific functions of specific types of structures to be called when such structures are dispatched into the system to allow a decoupled form of communication.

*This provide loose coupling as is needed.*

Using the `notifications` package is simple. By simply registering a function expecting a type, this sets up this function to be called once such type is seen.

```go

import "github.com/gu-io/gu/notifications"

type event struct{
  EventName string
  EventType string
}

func main(){

  notifications.Subscribe(func(eventName string){
    fmt.Printf("EventName[%q] occured.\n", eventName)
  })


  notifications.Subscribe(func(e event){
    fmt.Printf("EventName[%q] and EventType[%q] occured.\n", e.EventName, e.EventType)
  })

  notifications.Dispatch("Click") => `EventName["Click"] occured.`
  notifications.Dispatch(event{EventName:"ScrollDown", EventType:"scroll"}) => `EventName["ScrollDown"] and EventType["scroll"] occured.`
}
```
