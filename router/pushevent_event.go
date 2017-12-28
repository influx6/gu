package router

import "sync"

// PushEventSubscriber defines a interface that which is used to subscribe specifically for
// events  PushEvent type.
type PushEventSubscriber interface {
	Receive(PushEvent)
}

//=========================================================================================================

// PushEventHandler defines a structure type which implements the
// PushEventSubscriber interface and the EventDistributor interface.
type PushEventHandler struct {
	handle func(PushEvent)
}

// NewPushEventHandler returns a new instance of a PushEventHandler.
func NewPushEventHandler(fn func(PushEvent)) *PushEventHandler {
	return &PushEventHandler{
		handle: fn,
	}
}

// Receive takes the giving value and execute it against the underline handler.
func (sn *PushEventHandler) Receive(elem PushEvent) {
	sn.handle(elem)
}

// Handle takes the giving value and asserts the expected value to match the
// PushEvent type then passes it to the Receive method.
func (sn *PushEventHandler) Handle(receive interface{}) {
	if elem, ok := receive.(PushEvent); ok {
		sn.Receive(elem)
	}
}

//=========================================================================================================

// PushEventNotification defines a structure type which must be used to
// receive PushEvent type has a event.
type PushEventNotification struct {
	sml        sync.Mutex
	subs       []PushEventSubscriber
	validation func(PushEvent) bool
	register   map[PushEventSubscriber]int
}

// NewPushEventNotificationWith returns a new instance of PushEventNotification.
func NewPushEventNotificationWith(validation func(PushEvent) bool) *PushEventNotification {
	var elem PushEventNotification

	elem.validation = validation
	elem.register = make(map[PushEventSubscriber]int, 0)

	return &elem
}

// NewPushEventNotification returns a new instance of NewPushEventNotification.
func NewPushEventNotification() *PushEventNotification {
	var elem PushEventNotification
	elem.register = make(map[PushEventSubscriber]int, 0)

	return &elem
}

// UnNotify removes the given subscriber from the notification's list if found from future events.
func (sn *PushEventNotification) UnNotify(sub PushEventSubscriber) {
	sn.do(func() {
		index, ok := sn.register[sub]
		if !ok {
			return
		}

		sn.subs = append(sn.subs[:index], sn.subs[index+1:]...)
	})
}

// Notify adds the given subscriber into the notification list and will await an update of
// a new event of the given PushEvent type.
func (sn *PushEventNotification) Notify(sub PushEventSubscriber) {
	sn.do(func() {
		sn.register[sub] = len(sn.subs)
		sn.subs = append(sn.subs, sub)
	})
}

// Handle takes the giving value and asserts the expected value to be of
// the type and pass on to it's underline subscribers else ignoring the event.
func (sn *PushEventNotification) Handle(elem interface{}) {
	if elemEvent, ok := elem.(PushEvent); ok {
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
func (sn *PushEventNotification) do(fn func()) {
	if fn == nil {
		return
	}

	sn.sml.Lock()
	defer sn.sml.Unlock()

	fn()
}
