package gu

import "sync"

// AppUpdateSubscriber defines a interface that which is used to subscribe specifically for
// events  AppUpdate type.
type AppUpdateSubscriber interface {
	Receive(AppUpdate)
}

//=========================================================================================================

// AppUpdateHandler defines a structure type which implements the
// AppUpdateSubscriber interface and the EventDistributor interface.
type AppUpdateHandler struct {
	handle func(AppUpdate)
}

// NewAppUpdateHandler returns a new instance of a AppUpdateHandler.
func NewAppUpdateHandler(fn func(AppUpdate)) *AppUpdateHandler {
	return &AppUpdateHandler{
		handle: fn,
	}
}

// Receive takes the giving value and execute it against the underline handler.
func (sn *AppUpdateHandler) Receive(elem AppUpdate) {
	sn.handle(elem)
}

// Handle takes the giving value and asserts the expected value to match the
// AppUpdate type then passes it to the Receive method.
func (sn *AppUpdateHandler) Handle(receive interface{}) {
	if elem, ok := receive.(AppUpdate); ok {
		sn.Receive(elem)
	}
}

//=========================================================================================================

// AppUpdateNotification defines a structure type which must be used to
// receive AppUpdate type has a event.
type AppUpdateNotification struct {
	sml        sync.Mutex
	subs       []AppUpdateSubscriber
	validation func(AppUpdate) bool
	register   map[AppUpdateSubscriber]int
}

// NewAppUpdateNotificationWith returns a new instance of AppUpdateNotification.
func NewAppUpdateNotificationWith(validation func(AppUpdate) bool) *AppUpdateNotification {
	var elem AppUpdateNotification

	elem.validation = validation
	elem.register = make(map[AppUpdateSubscriber]int, 0)

	return &elem
}

// NewAppUpdateNotification returns a new instance of NewAppUpdateNotification.
func NewAppUpdateNotification() *AppUpdateNotification {
	var elem AppUpdateNotification
	elem.register = make(map[AppUpdateSubscriber]int, 0)

	return &elem
}

// UnNotify removes the given subscriber from the notification's list if found from future events.
func (sn *AppUpdateNotification) UnNotify(sub AppUpdateSubscriber) {
	sn.do(func() {
		index, ok := sn.register[sub]
		if !ok {
			return
		}

		sn.subs = append(sn.subs[:index], sn.subs[index+1:]...)
	})
}

// Notify adds the given subscriber into the notification list and will await an update of
// a new event of the given AppUpdate type.
func (sn *AppUpdateNotification) Notify(sub AppUpdateSubscriber) {
	sn.do(func() {
		sn.register[sub] = len(sn.subs)
		sn.subs = append(sn.subs, sub)
	})
}

// Handle takes the giving value and asserts the expected value to be of
// the type and pass on to it's underline subscribers else ignoring the event.
func (sn *AppUpdateNotification) Handle(elem interface{}) {
	if elemEvent, ok := elem.(AppUpdate); ok {
		if sn.validation != nil && sn.validation(elemEvent) {
			sn.do(func() {
				for _, sub := range sn.subs {
					sub.Receive(elemEvent)
				}
			})

			return
		}

		sn.do(func() {
			for _, sub := range sn.subs {
				sub.Receive(elemEvent)
			}
		})
	}
}

// do performs action with the mutex locked and unlocked appropriately, ensuring safe
// concurrent access.
func (sn *AppUpdateNotification) do(fn func()) {
	if fn == nil {
		return
	}

	sn.sml.Lock()
	defer sn.sml.Unlock()

	fn()
}
