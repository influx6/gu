package gu

import "sync"

// ViewUpdateSubscriber defines a interface that which is used to subscribe specifically for
// events  ViewUpdate type.
type ViewUpdateSubscriber interface {
	Receive(ViewUpdate)
}

//=========================================================================================================

// ViewUpdateHandler defines a structure type which implements the
// ViewUpdateSubscriber interface and the EventDistributor interface.
type ViewUpdateHandler struct {
	handle func(ViewUpdate)
}

// NewViewUpdateHandler returns a new instance of a ViewUpdateHandler.
func NewViewUpdateHandler(fn func(ViewUpdate)) *ViewUpdateHandler {
	return &ViewUpdateHandler{
		handle: fn,
	}
}

// Receive takes the giving value and execute it against the underline handler.
func (sn *ViewUpdateHandler) Receive(elem ViewUpdate) {
	sn.handle(elem)
}

// Handle takes the giving value and asserts the expected value to match the
// ViewUpdate type then passes it to the Receive method.
func (sn *ViewUpdateHandler) Handle(receive interface{}) {
	if elem, ok := receive.(ViewUpdate); ok {
		sn.Receive(elem)
	}
}

//=========================================================================================================

// ViewUpdateNotification defines a structure type which must be used to
// receive ViewUpdate type has a event.
type ViewUpdateNotification struct {
	sml        sync.Mutex
	subs       []ViewUpdateSubscriber
	validation func(ViewUpdate) bool
	register   map[ViewUpdateSubscriber]int
}

// NewViewUpdateNotificationWith returns a new instance of ViewUpdateNotification.
func NewViewUpdateNotificationWith(validation func(ViewUpdate) bool) *ViewUpdateNotification {
	var elem ViewUpdateNotification

	elem.validation = validation
	elem.register = make(map[ViewUpdateSubscriber]int, 0)

	return &elem
}

// NewViewUpdateNotification returns a new instance of NewViewUpdateNotification.
func NewViewUpdateNotification() *ViewUpdateNotification {
	var elem ViewUpdateNotification
	elem.register = make(map[ViewUpdateSubscriber]int, 0)

	return &elem
}

// UnNotify removes the given subscriber from the notification's list if found from future events.
func (sn *ViewUpdateNotification) UnNotify(sub ViewUpdateSubscriber) {
	sn.do(func() {
		index, ok := sn.register[sub]
		if !ok {
			return
		}

		sn.subs = append(sn.subs[:index], sn.subs[index+1:]...)
	})
}

// Notify adds the given subscriber into the notification list and will await an update of
// a new event of the given ViewUpdate type.
func (sn *ViewUpdateNotification) Notify(sub ViewUpdateSubscriber) {
	sn.do(func() {
		sn.register[sub] = len(sn.subs)
		sn.subs = append(sn.subs, sub)
	})
}

// Handle takes the giving value and asserts the expected value to be of
// the type and pass on to it's underline subscribers else ignoring the event.
func (sn *ViewUpdateNotification) Handle(elem interface{}) {
	if elemEvent, ok := elem.(ViewUpdate); ok {
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
func (sn *ViewUpdateNotification) do(fn func()) {
	if fn == nil {
		return
	}

	sn.sml.Lock()
	defer sn.sml.Unlock()

	fn()
}
