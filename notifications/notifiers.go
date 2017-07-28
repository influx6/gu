package notifications

// AppEvent defines a struct to contain a event which occurs to be delivered to
// a giving AppNotification instance.
//
//@notification:event
type AppEvent struct {
	UUID  string
	Event interface{}
}

// AppNotification defines a structure which provides a local notification
// framework for the pubsub.
func AppNotification(uid string) {
	// var app AppNotification
	// app.uid = uid

	// Subscribe(func(event AppEvent) {
	// 	if event.UUID != app.uid {
	// 		return
	// 	}

	// 	app.Dispatch(event.Event)
	// })

	// return &app
}
