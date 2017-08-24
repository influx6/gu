// +build ignore

// The generation of this package was inspired by Neelance work on DOM (https://github.com/neelance/dom)

package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type event struct {
	Name string
	Link string
	Desc string
}

func main() {
	nameMap := map[string]string{
		"CssRuleViewCSSLinkClicked": "CSSRuleViewCSSLinkClicked",
		"CssRuleViewRefreshed":      "CSSRuleViewRefreshed",
		"CssRuleViewChanged":        "CSSRuleViewChange",
		"afterscriptexecute":        "AfterScriptExecute",
		"afterprint":                "AfterPrint",
		"animationend":              "AnimationEnd",
		"animationiteration":        "AnimationIteration",
		"animationstart":            "AnimationStart",
		"audioprocess":              "AudioProcess",
		"beforeprint":               "BeforePrint",
		"beforeunload":              "BeforeUnload",
		"canplay":                   "CanPlay",
		"canplaythrough":            "CanPlayThrough",
		"chargingchange":            "ChargingChange",
		"chargingtimechange":        "ChargingTimeChange",
		"compassneedscalibration":   "CompassNeedsCalibration",
		"compositionend":            "CompositionEnd",
		"compositionstart":          "CompositionStart",
		"compositionupdate":         "CompositionUpdate",
		"contextmenu":               "ContextMenu",
		"dblclick":                  "DblClick",
		"devicelight":               "DeviceLight",
		"devicemotion":              "DeviceMotion",
		"deviceorientation":         "DeviceOrientation",
		"deviceproximity":           "DeviceProximity",
		"dischargingtimechange":     "DischargingTimeChange",
		"dragend":                   "DragEnd",
		"dragenter":                 "DragEnter",
		"dragleave":                 "DragLeave",
		"dragover":                  "DragOver",
		"dragstart":                 "DragStart",
		"durationchange":            "DurationChange",
		"focusin":                   "FocusIn",
		"focusout":                  "FocusOut",
		"fullscreenchange":          "FullScreenChange",
		"fullscreenerror":           "FullScreenError",
		"gamepadconnected":          "GamepadConnected",
		"gamepaddisconnected":       "GamepadDisconnected",
		"hashchange":                "HashChange",
		"beforescriptexecute":       "BeforeScriptExecute",
		"beforeinstallprompt":       "BeforeInstallPrompt",
		"keydown":                   "KeyDown",
		"keypress":                  "KeyPress",
		"keyup":                     "KeyUp",
		"languagechange":            "LanguageChange",
		"levelchange":               "LevelChange",
		"loadeddata":                "LoadedData",
		"loadedmetadata":            "LoadedMetadata",
		"loadend":                   "LoadEnd",
		"loadstart":                 "LoadStart",
		"mousedown":                 "MouseDown",
		"mouseenter":                "MouseEnter",
		"mouseleave":                "MouseLeave",
		"mousemove":                 "MouseMove",
		"mouseout":                  "MouseOut",
		"mouseover":                 "MouseOver",
		"mouseup":                   "MouseUp",
		"noupdate":                  "NoUpdate",
		"orientationchange":         "OrientationChange",
		"pagehide":                  "PageHide",
		"pageshow":                  "PageShow",
		"smartcard-remove":          "SmartCardRemove",
		"smartcard-insert":          "SmartCardInsert",
		"pointerlockchange":         "PointerLockChange",
		"pointerlockerror":          "PointerLockError",
		"popstate":                  "PopState",
		"ratechange":                "RateChange",
		"readystatechange":          "ReadystateChange",
		"timeupdate":                "TimeUpdate",
		"touchcancel":               "TouchCancel",
		"touchend":                  "TouchEnd",
		"touchenter":                "TouchEnter",
		"touchleave":                "TouchLeave",
		"touchmove":                 "TouchMove",
		"touchstart":                "TouchStart",
		"transitionend":             "TransitionEnd",
		"updateready":               "UpdateReady",
		"upgradeneeded":             "UpgradeNeeded",
		"userproximity":             "UserProximity",
		"versionchange":             "VersionChange",
		"visibilitychange":          "VisibilityChange",
		"volumechange":              "VolumeChange",
	}

	ignore := map[string]bool{
		"error": true,
	}

	doc, err := goquery.NewDocument("https://developer.mozilla.org/en-US/docs/Web/Events")
	if err != nil {
		panic(err)
	}

	events := make(map[string]*event)

	tables := doc.Find(".standard-table")

	tables.Each(func(ind int, item *goquery.Selection) {
		item.Eq(0).Find("tr").Each(func(i int, s *goquery.Selection) {
			cols := s.Find("td")
			if cols.Length() == 0 || cols.Find(".icon-thumbs-down-alt").Length() != 0 {
				return
			}
			link := cols.Eq(0).Find("a").Eq(0)
			var e event

			// fmt.Printf("Name: %q -> %q\n", nameMap[link.Text()], link.Text())

			if newName, ok := nameMap[link.Text()]; ok {
				e.Name = newName
			} else {
				e.Name = link.Text()
			}

			e.Link, _ = link.Attr("href")
			e.Desc = strings.TrimSpace(cols.Eq(3).Text())
			if e.Desc == "" {
				e.Desc = "(no documentation)"
			}

			funName := e.Name
			if strings.Contains(funName, "-") {
				parts := strings.Split(funName, "-")
				for ind, sm := range parts {
					parts[ind] = capitalize(sm)
				}

				funName = strings.Join(parts, "")
			}

			funName = capitalize(funName)

			if e.Name == "" || ignore[e.Name] {
				return
			}

			events[funName] = &e
		})
	})

	var names []string
	for name := range events {
		names = append(names, name)
	}

	sort.Strings(names)

	file, err := os.Create("event.gen.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprint(file, `// The generation of this package was inspired by Neelance work on DOM (https://github.com/neelance/dom)

//go:generate go run generate.go

// Documentation source: "Event reference" by Mozilla Contributors, https://developer.mozilla.org/en-US/docs/Web/Events, licensed under CC-BY-SA 2.5.

//Package events defines the event binding system that combines different libraries to create a interesting event system.
package events

import (
	"github.com/gu-io/gu/trees"
	"github.com/gu-io/gu/common"
	"github.com/gu-io/gu/notifications"
)


// EventHandler defines a function type for event callbacks.
type EventHandler func(common.EventObject, *trees.Markup)

// WrapHandler wraps the function returning a EventHandler to call the provided
// function to be called when the event occurs without need for the arguments.
func WrapHandler(callback func()) EventHandler {
	return func(ev common.EventObject, root *trees.Markup){
		callback()
	}
}

// WrapEventOnlyHandler wraps the function returning a EventHandler to call the provided
// function to be called when the event occurs without need for the arguments.
func WrapEventOnlyHandler(callback func(common.EventObject)) EventHandler {
	return func(ev common.EventObject, root *trees.Markup) {
		callback(ev)
	}
}
`)

	for _, name := range names {
		e := events[name]
		fmt.Fprintf(file, `
// %sEvent Documentation is as below: %q
// https://developer.mozilla.org%s
// This event provides options() to be called when the events is triggered and an optional selector which will override the internal selector
// mechanism of the domtrees.Element i.e if the selectorOverride argument is an empty string then domtrees.Element will create an
// appropriate selector matching its type and uid value in this format  (ElementType[uid='UID_VALUE']) but if
// the selector value is not empty then that becomes the default selector used match the event with.
func %sEvent(callback interface{}, options ...trees.EventOptions) *trees.Event {
	var handler EventHandler

	switch cb := callback.(type){
	case func():
		handler = WrapHandler(cb)
	case func(common.EventObject):
		handler = WrapEventOnlyHandler(cb)
	case func(common.EventObject, *trees.Markup):
		handler = cb
	default:
		panic("Unacceptable type for event callback")
	}

	ops := append([]trees.EventOptions{trees.EventType(%q)}, options...)

	ev := trees.NewEvent(ops...)

	eventHandler := common.NewEventBroadcastHandler(func (evm common.EventBroadcast){
		if ev.ID() != evm.EventID{
			return
		}

		handler(evm.Event, ev.Tree)
	})

	ev.Remove = notifications.SubscribeWithRemover(eventHandler)

	return ev
}
`, name, e.Desc, e.Link[6:], name, e.Name)
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}
