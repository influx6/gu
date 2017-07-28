package notifications

// AppEvent defines a struct to contain a event which occurs to be delivered to
// a giving AppNotification instance.
type AppEvent struct {
	UUID  string
	Event interface{}
}

// // AppNotification defines a structure which provides a local notification
// // framework for the pubsub.
// type AppNotification struct {
// 	core *mque.MQue
// 	uid  string
// }

// // New returns a new instance of the AppNotification type.
// func New(uid string) *AppNotification {
// 	var app AppNotification
// 	app.core = mque.New()
// 	app.uid = uid

// 	Subscribe(func(event AppEvent) {
// 		if event.UUID != app.uid {
// 			return
// 		}

// 		app.Dispatch(event.Event)
// 	}, func() {
// 		//TODO: Do we want to do anything here?
// 	})

// 	return &app
// }

// // Subscribe adds a new listener to the dispatcher.
// func (app *AppNotification) Subscribe(q interface{}, end ...func()) mque.End {
// 	return app.core.Q(q, end...)
// }

// // Dispatch emits a event into the notification  callback listeners.
// func (app *AppNotification) Dispatch(q interface{}) {
// 	app.core.Run(q)
// }

// // GlobalDispatch emits a event to the global dispatch callback listeners.
// func (app *AppNotification) GlobalDispatch(q interface{}) {
// 	dispatch.Run(q)
// }
