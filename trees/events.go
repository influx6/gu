package trees

import (
	"fmt"
	"strings"
)

// EventBroadcast defines a struct which gets published for the events.
type EventBroadcast struct {
	EventName string      `json:"event"`
	EventID   string      `json:"event_id"`
	Event     EventObject `json:"event_object"`
}

// EventObject defines a interface for the basic methods which events needs
// expose.
type EventObject interface {
	RemoveEvent()
	Underlying() interface{}
}

// Remover defines an interface that exposes a method to remove a giving event.
type Remover interface {
	Remove()
}

// Event provide a meta registry for helps in registering events for dom markups
// which is translated to the nodes themselves
type Event struct {
	Type                     string
	PreventDefault           bool
	StopPropagation          bool
	UseCapture               bool
	StopImmediatePropagation bool
	Tree                     *Markup
	Remove                   Remover
	secTarget                string
}

// NewEvent returns a event object that allows registering events to eventlisteners.
func NewEvent(evtype string, evtarget string, preventdef bool, stopPropagation bool, stopImmediate bool, useCapture bool) *Event {
	return &Event{
		Type:                     evtype,
		PreventDefault:           preventdef,
		StopPropagation:          stopPropagation,
		StopImmediatePropagation: stopImmediate,
		UseCapture:               useCapture,
		secTarget:                evtarget,
	}
}

// Target returns the target of the giving event.
func (e *Event) Target() string {
	if e.Tree != nil {
		return e.Tree.EventID()
	}

	return e.secTarget
}

// EventJSON defines a struct which contains the giving events and
// and tree of the giving tree.
type EventJSON struct {
	ParentSelector           string `json:"ParentSelector"`
	EventSelector            string `json:"EventSelector"`
	EventName                string `json:"EventName"`
	Event                    string `json:"Event"`
	PreventDefault           bool   `json:"PreventDefault"`
	StopPropagation          bool   `json:"StopPropagation"`
	UseCapture               bool   `json:"UseCapture"`
	StopImmediatePropagation bool   `json:"StopImmediatePropagation"`
}

// EventJSON returns the event json structure which represent the giving event.
func (e *Event) EventJSON() EventJSON {
	return EventJSON{
		Event:                    e.Type,
		UseCapture:               e.UseCapture,
		EventName:                e.EventName(),
		EventSelector:            e.EventSelector(),
		ParentSelector:           e.ParentEventSelector(),
		PreventDefault:           e.PreventDefault,
		StopPropagation:          e.StopPropagation,
		StopImmediatePropagation: e.StopImmediatePropagation,
	}
}

// ParentEventSelector returns the parent selector for this events markup tree.
func (e *Event) ParentEventSelector() string {
	if e.Tree != nil {
		return e.Tree.IDSelector(true)
	}

	return ""
}

// EventSelector returns the selector for this events tree.
func (e *Event) EventSelector() string {
	if e.Tree != nil && e.secTarget == "" {
		return e.Tree.IDSelector(false)
	}

	return e.secTarget
}

// EventName returns the giving name of the event.
func (e *Event) EventName() string {
	eventName := strings.ToUpper(e.Type[:1]) + e.Type[1:]
	if strings.HasSuffix(eventName, "Event") {
		return eventName
	}

	return eventName + "Event"
}

// ID returns the uique event id string for the event.
func (e *Event) ID() string {
	return fmt.Sprintf("%s#%s", e.Target(), e.Type)
}

// Clone  returns a new Event object from this.
func (e *Event) Clone() *Event {
	return &Event{
		Type:                     e.Type,
		secTarget:                e.secTarget,
		PreventDefault:           e.PreventDefault,
		UseCapture:               e.UseCapture,
		StopPropagation:          e.StopPropagation,
		StopImmediatePropagation: e.StopImmediatePropagation,
	}
}

// Apply adds the event into the elements events lists
func (e *Event) Apply(ex *Markup) {
	if !ex.allowEvents {
		return
	}

	e.Tree = ex

	ex.AddEvent(*e)
}

// String returns the string representation of the giving event.
func (e *Event) String() string {
	return fmt.Sprintf("%#v", e.EventJSON())
}

//==============================================================================
