// Package eventx attempts to provide a concrete go type for different events
// related to the DOM. Not all the events and their functionality can be supported.
// It exists has a package to allow access to this events without any tie into the
// corresponding js primitive. If you prefer a fuller support. Use GopherJS (https://github.com/gopherjs/gopherjs).
package eventx

import (
	"time"

	"github.com/gu-io/gu/common"
	"github.com/gu-io/gu/router/cache"
)

// Element defines a string type which contains the markup of the giving element.
type Element string

// DeltaMode is an indication of the units of measurement for a delta value.
type DeltaMode uint64

// KeyCode represents a system and implementation dependent numerical
// code identifying the unmodified value of the pressed key.
type KeyCode uint8

// Keyboard keys.
const (
	KeyBackspace          KeyCode = 8
	KeyTab                KeyCode = 9
	KeyEnter              KeyCode = 13
	KeyShift              KeyCode = 16
	KeyCtrl               KeyCode = 17
	KeyAlt                KeyCode = 18
	KeyPauseBreak         KeyCode = 19
	KeyCapsLock           KeyCode = 20
	KeyEsc                KeyCode = 27
	KeySpace              KeyCode = 32
	KeyPageUp             KeyCode = 33
	KeyPageDown           KeyCode = 34
	KeyEnd                KeyCode = 35
	KeyHome               KeyCode = 36
	KeyLeft               KeyCode = 37
	KeyUp                 KeyCode = 38
	KeyRight              KeyCode = 39
	KeyDown               KeyCode = 40
	KeyPrintScreen        KeyCode = 44
	KeyInsert             KeyCode = 45
	KeyDelete             KeyCode = 46
	Key0                  KeyCode = 48
	Key1                  KeyCode = 49
	Key2                  KeyCode = 50
	Key3                  KeyCode = 51
	Key4                  KeyCode = 52
	Key5                  KeyCode = 53
	Key6                  KeyCode = 54
	Key7                  KeyCode = 55
	Key8                  KeyCode = 56
	Key9                  KeyCode = 57
	KeyA                  KeyCode = 65
	KeyB                  KeyCode = 66
	KeyC                  KeyCode = 67
	KeyD                  KeyCode = 68
	KeyE                  KeyCode = 69
	KeyF                  KeyCode = 70
	KeyG                  KeyCode = 71
	KeyH                  KeyCode = 72
	KeyI                  KeyCode = 73
	KeyJ                  KeyCode = 74
	KeyK                  KeyCode = 75
	KeyL                  KeyCode = 76
	KeyM                  KeyCode = 77
	KeyN                  KeyCode = 78
	KeyO                  KeyCode = 79
	KeyP                  KeyCode = 80
	KeyQ                  KeyCode = 81
	KeyR                  KeyCode = 82
	KeyS                  KeyCode = 83
	KeyT                  KeyCode = 84
	KeyU                  KeyCode = 85
	KeyV                  KeyCode = 86
	KeyW                  KeyCode = 87
	KeyX                  KeyCode = 88
	KeyY                  KeyCode = 88
	KeyZ                  KeyCode = 90
	KeyMeta               KeyCode = 91
	KeyMenu               KeyCode = 93
	KeyNumPad0            KeyCode = 96
	KeyNumPad1            KeyCode = 97
	KeyNumPad2            KeyCode = 98
	KeyNumPad3            KeyCode = 99
	KeyNumPad4            KeyCode = 100
	KeyNumPad5            KeyCode = 101
	KeyNumPad6            KeyCode = 102
	KeyNumPad7            KeyCode = 103
	KeyNumPad8            KeyCode = 104
	KeyNumPad9            KeyCode = 105
	KeyNumPadMult         KeyCode = 106
	KeyNumPadPlus         KeyCode = 107
	KeyNumPadMin          KeyCode = 109
	KeyNumPadDot          KeyCode = 110
	KeyNumPadDecimal      KeyCode = 111
	KeyF1                 KeyCode = 112
	KeyF2                 KeyCode = 113
	KeyF3                 KeyCode = 114
	KeyF4                 KeyCode = 115
	KeyF5                 KeyCode = 116
	KeyF6                 KeyCode = 117
	KeyF7                 KeyCode = 118
	KeyF8                 KeyCode = 119
	KeyF9                 KeyCode = 120
	KeyF10                KeyCode = 121
	KeyF11                KeyCode = 122
	KeyF12                KeyCode = 123
	KeyNumLock            KeyCode = 144
	KeyMute               KeyCode = 173
	KeyVolumeDown         KeyCode = 174
	KeyVolumeUp           KeyCode = 175
	KeySemicolon          KeyCode = 186
	KeyEqual              KeyCode = 187
	KeyComa               KeyCode = 188
	KeyDash               KeyCode = 189
	KeyDot                KeyCode = 190
	KeySlash              KeyCode = 191
	KeyBackquote          KeyCode = 192
	KeySquareBracketLeft  KeyCode = 219
	KeyBackslash          KeyCode = 220
	KeySquareBracketRight KeyCode = 221
	KeyQuote              KeyCode = 222
)

// KeyLocation represents the location of the key on the keyboard or
// other input device.
type KeyLocation uint8

// Keyboard locations constants.
const (
	KeyLocationStandard KeyLocation = iota
	KeyLocationLeft
	KeyLocationRight
	KeyLocationNumpad
)

const (
	// DeltaPixel defines the mouse movement state.
	DeltaPixel = 0

	// DeltaLine defines the mouse movement state.
	DeltaLine = 1

	// DeltaPage defines the mouse movement state.
	DeltaPage = 2
)

// ChangeEvent represents the data passed in a onchange event.
type ChangeEvent struct {
	Core  interface{} `json:"-"`
	Value string
}

// InputDeviceCapabilities defines a struct to contain input capbilities for a
// inputtype.
type InputDeviceCapabilities struct {
	FiresTouchEvent bool
}

// IDBVersionChangeEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type IDBVersionChangeEvent struct {
	Core       interface{} `json:"-"`
	OldVersion int64
	NewVersion int64
}

// HashChangeEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type HashChangeEvent struct {
	Core interface{} `json:"-"`
	Old  string
	New  string
}

// Button defines a struct which holds button information as
// related with Gamepads.
type Button struct {
	Value   float64
	Pressed bool
}

// Gamepad defines a struct which holds the gamepad object porperties.
type Gamepad struct {
	DisplayID string
	ID        string
	Index     int
	Mapping   string
	Connected bool
	Buttons   []Button
	Axes      []float64
	Timestamp float64
}

// GamepadEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type GamepadEvent struct {
	Core    interface{} `json:"-"`
	Gamepad Gamepad
}

// AnimationEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type AnimationEvent struct {
	Core          interface{} `json:"-"`
	AnimationName string
	PseudoElement string
	ElapsedTime   float64
}

// AudioBuffer defines a struct to represent the buffer associated with
// with a giving AudioProcessingEvent.
// When copying channel data from javascript ensure to follow this:
// myArrayBuffer.copyFromChannel(destination,channelNumber,startInChannel);
// Where:
// destination => the array buffer
// channelNumber => channel number (starts from 0....totalChannels)
// startInChannel => index of internal channel array
// var myArrayBuffer = audioCtx.createBuffer(2, frameCount, audioCtx.sampleRate);
// var anotherArray = new Float32Array;
// myArrayBuffer.copyFromChannel(anotherArray,1,0);
type AudioBuffer struct {
	SampleRate         float64
	Duration           time.Duration
	Channels           int
	SampleFramesLength int
	ChannelData        [][]byte
}

// AudioProcessingEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
// Deprecated Event.
type AudioProcessingEvent struct {
	Core         interface{} `json:"-"`
	PlaybackTime float64
	InputBuffer  AudioBuffer
	OutputBuffer AudioBuffer
}

// BeforeUnloadEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type BeforeUnloadEvent struct {
	Core interface{} `json:"-"`
}

// BeforeInputEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type BeforeInputEvent struct {
	Core interface{} `json:"-"`
}

// BlobEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type BlobEvent struct {
	Core interface{} `json:"-"`
	Data []byte
}

// DataTransferItem defines a DataTransferItem file item.
type DataTransferItem struct {
	Name string
	Data []byte
	Size int
}

// DataTransferItemList defines a struct which contains a list of DataTransferItems.
type DataTransferItemList struct {
	Items []DataTransferItem
}

// DataTransfer defines a struct to represent the data retrieved from the data
// transfer object.
type DataTransfer struct {
	DropEffect    string
	EffectAllowed string
	Files         []DataTransferItem
	Items         DataTransferItemList
	Types         []string
}

// ClipboardEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type ClipboardEvent struct {
	Core interface{} `json:"-"`
	Data DataTransfer
}

// CloseEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type CloseEvent struct {
	Core     interface{} `json:"-"`
	Code     int
	Reason   string
	WasClean bool
}

// CompositionEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type CompositionEvent struct {
	Core   interface{} `json:"-"`
	Text   string
	Data   string
	Locale string
}

// CSSFontFaceLoadEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type CSSFontFaceLoadEvent struct {
	Core interface{} `json:"-"`
}

// CustomEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type CustomEvent struct {
	Core   interface{} `json:"-"`
	Detail interface{}
}

// DropEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type DropEvent struct {
	*MouseEvent
	Core         interface{} `json:"-"`
	DataTransfer DataTransfer
}

// DragLeaveEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type DragLeaveEvent struct {
	*MouseEvent
	Core         interface{} `json:"-"`
	DataTransfer DataTransfer
}

// DragStartEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type DragStartEvent struct {
	*MouseEvent
	Core         interface{} `json:"-"`
	DataTransfer DataTransfer
}

// DragEndEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type DragEndEvent struct {
	*MouseEvent
	Core         interface{} `json:"-"`
	DataTransfer DataTransfer
}

// DragOverEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type DragOverEvent struct {
	*MouseEvent
	Core         interface{} `json:"-"`
	DataTransfer DataTransfer
}

// DragExitEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type DragExitEvent struct {
	*MouseEvent
	Core         interface{} `json:"-"`
	DataTransfer DataTransfer
}

// DragEnterEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type DragEnterEvent struct {
	*MouseEvent
	Core         interface{} `json:"-"`
	DataTransfer DataTransfer
}

// DragEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type DragEvent struct {
	*MouseEvent
	Core         interface{} `json:"-"`
	DataTransfer DataTransfer
}

// DeviceLightEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type DeviceLightEvent struct {
	Core  interface{} `json:"-"`
	Value float64
}

// MotionData defines a struct contain motion data.
type MotionData struct {
	X float64
	Y float64
	Z float64
}

// RotationData defines a struct contain motion data.
type RotationData struct {
	Alpha float64
	Beta  float64
	Gamma float64
}

// DeviceMotionEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type DeviceMotionEvent struct {
	Core                         interface{} `json:"-"`
	Interval                     float64
	Acceleration                 MotionData
	AccelerationIncludingGravity MotionData
	RotationRate                 RotationData
}

// DeviceOrientationEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type DeviceOrientationEvent struct {
	Core     interface{} `json:"-"`
	Absolute bool
	Alpha    float64
	Beta     float64
	Gamma    float64
}

// DeviceProximityEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type DeviceProximityEvent struct {
	Core  interface{} `json:"-"`
	Max   float64
	Min   float64
	Value float64
}

// FetchEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type FetchEvent struct {
	Core     interface{} `json:"-"`
	IsReload bool
	Request  cache.Request
	ClientID string
}

// DOMTransactionEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type DOMTransactionEvent struct {
	Core interface{} `json:"-"`
}

// EditingBeforeInputEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type EditingBeforeInputEvent struct {
	Core interface{} `json:"-"`
}

// ErrorEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type ErrorEvent struct {
	Core       interface{} `json:"-"`
	Message    string
	Filename   string
	LineNumber int
	ColNumber  int
	Error      error
}

// InputEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type InputEvent struct {
	Core        interface{} `json:"-"`
	Data        string
	IsComposing bool
}

// FocusEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type FocusEvent struct {
	Core interface{} `json:"-"`
}

// MutationRecord defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type MutationRecord struct {
	Type             string
	AddedNodes       []Element
	RemovedNodes     []Element
	PreSibling       Element
	NextSibling      Element
	AttributeName    string
	AttributreNameNS string
}

// MutationEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type MutationEvent struct {
	Core interface{} `json:"-"`
}

// MouseEvent represents data fired when interacting
// with a pointing device (such as a mouse).
type MouseEvent struct {
	*UIEvent
	Core     interface{} `json:"-"`
	ClientX  float64
	ClientY  float64
	PageX    float64
	PageY    float64
	OffsetX  float64
	OffsetY  float64
	ScreenX  float64
	ScreenY  float64
	MovemenX float64
	MovemenY float64
	Region   int
	Button   int
	Detail   int
	AltKey   bool
	CtrlKey  bool
	MetaKey  bool
	ShiftKey bool
}

// WebGLContextEvent represents data fired when a wheel button of a
// pointing device (usually a mouse) is rotated.
type WebGLContextEvent struct {
	StatusMessage string
}

// WheelEvent represents data fired when a wheel button of a
// pointing device (usually a mouse) is rotated.
type WheelEvent struct {
	Core      interface{} `json:"-"`
	DeltaX    float64
	DeltaY    float64
	DeltaZ    float64
	DeltaMode DeltaMode
}

// KeyboardEvent represents data fired when the keyboard is used.
type KeyboardEvent struct {
	Core          interface{} `json:"-"`
	CharCode      int
	KeyCode       KeyCode
	KeyLocation   KeyLocation
	Location      int
	Key           string
	KeyIdentifier string
	Locale        string
	AltKey        bool
	CtrlKey       bool
	MetaKey       bool
	ShiftKey      bool
	Repeat        bool
	ModifiedState bool
}

// OfflineAudioCompletionEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type OfflineAudioCompletionEvent struct {
	Core           interface{} `json:"-"`
	RenderedBuffer AudioBuffer
}

// PageTransitionEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type PageTransitionEvent struct {
	Core      interface{} `json:"-"`
	Persisted bool
}

// PointerEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type PointerEvent struct {
	Core        interface{} `json:"-"`
	PointerID   int
	Width       int
	Height      int
	Pressure    float64
	TiltX       float64
	TiltY       float64
	IsPrimary   bool
	PointerType string
}

// PopStateEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type PopStateEvent struct {
	Core interface{} `json:"-"`
}

// ProgressEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type ProgressEvent struct {
	Core             interface{} `json:"-"`
	LengthComputable bool
	Loaded           uint64
	Total            int
}

// RelatedEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type RelatedEvent struct {
	Core          interface{} `json:"-"`
	RelatedTarget Element
}

// RTCPeerConnectionIceEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type RTCPeerConnectionIceEvent struct {
	Core      interface{} `json:"-"`
	Candidate string
}

// RTCIdentityEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type RTCIdentityEvent struct {
	Core      interface{} `json:"-"`
	Assertion string
}

// SensorEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type SensorEvent struct {
	Core interface{} `json:"-"`
}

// StorageEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type StorageEvent struct {
	Core        interface{} `json:"-"`
	Key         string
	NewValue    string
	OldValue    string
	URL         string
	StorageArea interface{}
}

// TimeEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type TimeEvent struct {
	Core         interface{} `json:"-"`
	Long         interface{}
	AbstractView interface{}
}

// TransitionEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type TransitionEvent struct {
	Core          interface{} `json:"-"`
	PropertyName  string
	ElapsedTime   float64
	PseudoElement string
}

// UserProximityEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type UserProximityEvent struct {
	Core interface{} `json:"-"`
	Near bool
}

// UIEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type UIEvent struct {
	Core               interface{} `json:"-"`
	IsChar             bool
	LayerX             float64
	LayerY             float64
	PageX              float64
	PageY              float64
	Detail             int
	SourceCapabilities *InputDeviceCapabilities
}

// TrackEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type TrackEvent struct {
	Core interface{} `json:"-"`
}

// Touch defines a struct which holds touch list data related to touch events.
type Touch struct {
	Identifier float64
	ClientX    float64
	ClientY    float64
	PageX      float64
	PageY      float64
	OffsetX    float64
	OffsetY    float64
	ScreenX    float64
	ScreenY    float64
	Target     Element
}

// TouchList defines a list which holds touch data related to touch events.
type TouchList struct {
	Touches []Touch
	Length  int
}

// TouchEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type TouchEvent struct {
	Core          interface{} `json:"-"`
	AltKey        bool
	CtrlKey       bool
	MetaKey       bool
	ShiftKey      bool
	TargetTouches TouchList
	Touches       TouchList
}

// SVGZoomEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type SVGZoomEvent struct {
	Core interface{} `json:"-"`
}

// SVGEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type SVGEvent struct {
	Core interface{} `json:"-"`
}

// MessageEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
type MessageEvent struct {
	Core   interface{} `json:"-"`
	Data   []byte
	Origin string
	Source string
	Port   int
}

// MediaStream defines a struct for the media stream event.
// API: MediaStream.onaddtrack, MediaStream.onremovetrack.
// API: MediaStream.getTracks, MediaStream.getAudioTracks, MediaStream.getVideoTracks.
type MediaStream struct {
	Active bool
	Ended  bool
	ID     string
	Audios []MediaStreamTrack
	Videos []MediaStreamTrack
}

// MediaTrackSettings defines the struct which contains settiings for the MediaTrack
// API.
type MediaTrackSettings struct {
	DeviceID string
	GroupID  string
}

// MediaAudioTrackSettings defines the struct which contains settiings for the MediaTrack
// API.
type MediaAudioTrackSettings struct {
	MediaTrackSettings
	ChannelCount     int
	EchoCancellation bool
	Latency          float64
	SampleRate       int64
	SampleSize       int64
	Volume           float64
}

// MediaVideoTrackSettings defines the struct which contains settiings for the MediaTrack
// API.
type MediaVideoTrackSettings struct {
	MediaTrackSettings
	AspectRatio float64
	FacingMode  string
	FrameRate   float64
	Height      int64
	Width       int64
}

// MediaStreamTrack defines a track of a MediaStream.
type MediaStreamTrack struct {
	Core          interface{} `json:"-"`
	Enabled       bool
	ID            string
	Kind          string
	Label         string
	Muted         bool
	ReadyState    bool
	Remote        bool
	AudioSettings *MediaAudioTrackSettings
	VideoSettings *MediaVideoTrackSettings
}

// MediaStreamEvent defines a struct to contain the values of a promiximity
// event fired from a giving DOM.
// When using this combine them with the MediaStreamTrack.
type MediaStreamEvent struct {
	Core   interface{} `json:"-"`
	Stream MediaStream
}

// BasicEventMap defines a event type which defines a event type which is not
// supported by this package.
type BasicEventMap map[string]string

// BaseEvent defines the standard event returned to wrap a core event.
type BaseEvent struct {
	Core   interface{} `json:"-"`
	event  interface{}
	handle common.Remover
}

// NewBaseEvent returns a new instance of the base evet.
func NewBaseEvent(event interface{}, handle common.Remover) *BaseEvent {
	return &BaseEvent{
		event:  event,
		handle: handle,
	}
}

// RemoveEvent calls the finalizer of the giving event.
func (b *BaseEvent) RemoveEvent() {
	if b.handle != nil {
		b.handle.Remove()
	}
}

// Underlying returns the core event object returned by the
// base.
func (b *BaseEvent) Underlying() interface{} {
	return b.event
}

// GetEvent returns the base event type expected from the provided value.
func GetEvent() *BaseEvent {
	return nil
}
