package core

import (
	"encoding/json"

	events "github.com/gu-io/gu/eventx"
	"github.com/gu-io/gu/notifications/mque"
)

// GetEvent returns the giving event structure suited to the provided type.
func GetEvent(eventName string, eventJSON []byte, handle mque.End) (*events.BaseEvent, error) {
	switch eventName {
	case "AnimationEvent":
		var eventObject events.AnimationEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "AudioProcessingEvent":
		var eventObject events.AudioProcessingEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "BeforeInputEvent":
		var eventObject events.BeforeInputEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "BeforeUnloadEvent":
		var eventObject events.BeforeUnloadEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "BlobEvent":
		var eventObject events.BlobEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "ChangeEvent":
		var eventObject events.ChangeEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "ClipboardEvent":
		var eventObject events.ClipboardEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "CloseEvent":
		var eventObject events.CloseEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "CompositionEvent":
		var eventObject events.CompositionEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "CSSFontFaceLoadEvent":
		var eventObject events.CSSFontFaceLoadEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "CustomEvent":
		var eventObject events.CustomEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "DeviceLightEvent":
		var eventObject events.DeviceLightEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "DeviceMotionEvent":
		var eventObject events.DeviceMotionEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "DeviceOrientationEvent":
		var eventObject events.DeviceOrientationEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "DeviceProximityEvent":
		var eventObject events.DeviceProximityEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "DOMTransactionEvent":
		var eventObject events.DOMTransactionEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "DragStartEvent":
		var eventObject events.DragStartEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "DragEndEvent":
		var eventObject events.DragEndEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "DragEnterEvent":
		var eventObject events.DragEnterEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "DragLeaveEvent":
		var eventObject events.DragLeaveEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "DragOverEvent":
		var eventObject events.DragOverEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "DropEvent":
		var eventObject events.DropEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "DragEvent":
		var eventObject events.DragEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "EditingBeforeInputEvent":
		var eventObject events.EditingBeforeInputEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "ErrorEvent":
		var eventObject events.ErrorEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "FocusEvent":
		var eventObject events.FocusEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "GamepadEvent":
		var eventObject events.GamepadEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "HashChangeEvent":
		var eventObject events.HashChangeEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "IDBVersionChangeEvent":
		var eventObject events.IDBVersionChangeEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "KeyboardEvent":
		var eventObject events.KeyboardEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "MediaStreamEvent":
		var eventObject events.MediaStreamEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "MessageEvent":
		var eventObject events.MessageEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "MouseEvent":
		var eventObject events.MouseEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "MutationEvent":
		var eventObject events.MutationEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "OfflineAudioCompletionEvent":
		var eventObject events.OfflineAudioCompletionEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "PageTransitionEvent":
		var eventObject events.PageTransitionEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "PointerEvent":
		var eventObject events.PointerEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "PopStateEvent":
		var eventObject events.PopStateEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "ProgressEvent":
		var eventObject events.ProgressEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "RelatedEvent":
		var eventObject events.RelatedEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "RTCPeerConnectionIceEvent":
		var eventObject events.RTCPeerConnectionIceEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "SensorEvent":
		var eventObject events.SensorEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "StorageEvent":
		var eventObject events.StorageEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "SVGEvent":
		var eventObject events.SVGEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "SVGZoomEvent":
		var eventObject events.SVGZoomEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "TimeEvent":
		var eventObject events.TimeEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "TouchEvent":
		var eventObject events.TouchEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "TrackEvent":
		var eventObject events.TrackEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "TransitionEvent":
		var eventObject events.TransitionEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "UIEvent":
		var eventObject events.UIEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "UserProximityEvent":
		var eventObject events.UserProximityEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil
	case "WheelEvent":
		var eventObject events.WheelEvent

		if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
			return nil, err
		}

		return events.NewBaseEvent(&eventObject, handle), nil

	}

	var eventObject events.BasicEventMap

	if err := json.Unmarshal(eventJSON, &eventObject); err != nil {
		return nil, err
	}

	return events.NewBaseEvent(&eventObject, handle), nil
}
