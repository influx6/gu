// Package core.js provides javascript functions which provide similar functionalities
// to allow patching provided virtual DOM and query events and dom nodes as needed.


// Annonymouse function that initializes all properties on the root.
function GuClient(onMessages, SendChannel) {

    // GuJS provides the central object which is used to holds the handle for
    // all gujs properties.
    var GuJS = {};


    // unwanted defines specific object functions we dont want passed in during
    // collection of object property names.
    var unwanted = { "constructor": true, "toString": true };

    GuJS.eventsCore = {};
    GuJS.currentAppID = null;

    // GuJS.Dispatch defines a function which dispatches event object to the external
    // API.
    GuJS.Dispatch = function(name, model, meta) {
        SendChannel({ "type": name, "meta": meta, "data": model });
    };

    // GuJS.MakeEventCallback defines a function to generate a callback for an event
    // meta provided
    GuJS.MakeEventCallback = function(target, eventMeta) {
        return function(eventObj) {

            // Do we match the event and possible targets for the event
            // selector.
            GuJS.each(target.querySelectorAll(eventMeta.EventSelector), function(possible) {
                if (eventObj.target !== possible) {
                    return
                }

                if (eventMeta.PreventDefault) {
                    eventObj.preventDefault()
                }

                if (eventMeta.StopImmediatePropagation) {
                    eventObj.stopImmediatePropagation()
                }

                if (eventMeta.StopPropagation) {
                    eventObj.stopPropagation()
                }

                GuJS.Dispatch(GuJS.GetEvent(eventObj), eventMeta)
            })
        }
    };


    // GuJS.ExecuteCommand executes the provided command received.
    GuJS.ExecuteCommand = function(co) {
        if (co == null || co === undefined) {
            return
        }

        var command

        // if we are dealing with a string then parse with json.
        switch (co.constructor) {
            case String:
                command = JSON.parse(co)
            case Object:
                command = co
        }


        var head = document.querySelector("head")
        var body = document.querySelector("body")

        switch (command.Command) {
            case "RenderApp":
                // Rendering the app response is to clear what is currently in the view.
                // We want specific replicate the way the gopherjs driver updates apps
                // in swapping out the current content with the new content received.

                var app = command.App
                GuJS.currentAppID = app.AppID

                // Retrieve events map related to the giving app.
                var appEvents = GuJS.eventsCore[app.AppID] || { views: {}, base: { headEvents: [], bodyEvents: [] } }
                GuJS.eventsCore[app.AppID] = appEvents

                var nonGuHead = head.querySelectorAll("*:not([data-gen='gu'])")
                var nonGuBody = head.querySelectorAll("*:not([data-gen='gu'])")

                // Deregister all head base events.
                GuJS.each(appEvents.base.headEvents, function(cb) {
                    head.removeEventListener(cb.Event.Event, cb.Callback)
                })

                // Deregister all body base events.
                GuJS.each(appEvents.base.bodyEvents, function(cb) {
                    body.removeEventListener(cb.Event.Event, cb.Callback)
                })

                // Deregister all view events.
                GuJS.each(appEvents.views, function(view) {
                    GuJS.each(view, function(cb) {
                        body.removeEventListener(cb.Event.Event, cb.Callback)
                    })
                })

                appEvents.base.headEvents = [];
                appEvents.base.bodyEvents = [];
                appEvents.base.views = {};

                var headHTML = []
                var bodyHTML = []

                // Add the resource markup for the header.
                GuJS.each(app.HeadResources, function(item) {

                    // Generate the fragment for the giving markup.
                    var fragment = GuJS.createDOMFragment(item.Markup)

                    // Register all events for this markup.
                    GuJS.each(item.Events, function(event) {
                        var newEvent = {}
                        newEvent.Event = event
                        newEvent.Callback = GuJS.MakeEventCallback(head, event)

                        head.addEventListener(event.Event, newEvent.Callback, event.UseCapture);
                        appEvents.base.headEvents.push(newEvent);
                    })


                    headHTML.push(fragment);
                })

                GuJS.each(app.Head, function(item) {
                    var viewEvents = appEvents.views[item.ViewID] || []
                    appEvents.views[item.ViewID] = viewEvents

                    // Generate the fragment for the giving markup.
                    var fragment = GuJS.createDOMFragment(item.Tree.Markup)

                    // Register all events for this markup.
                    GuJS.each(item.Tree.Events, function(event) {
                        var newEvent = {}
                        newEvent.Event = event
                        newEvent.Callback = GuJS.MakeEventCallback(head, event)

                        head.addEventListener(event.Event, newEvent.Callback, event.UseCapture);
                        viewEvents.puhs(newEvent)
                    })

                    headHTML.push(fragment)
                })

                // Add the resource markup for the header.
                GuJS.each(app.Body, function(item) {
                    var viewEvents = appEvents.views[item.ViewID] || []
                    appEvents.views[item.ViewID] = viewEvents

                    // Generate the fragment for the giving markup.
                    var fragment = GuJS.createDOMFragment(item.Tree.Markup)

                    // Register all events for this markup.
                    GuJS.each(item.Tree.Events, function(event) {
                        var newEvent = {}
                        newEvent.Event = event
                        newEvent.Callback = GuJS.MakeEventCallback(body, event)

                        body.addEventListener(event.Event, newEvent.Callback, event.UseCapture);
                        viewEvents.push(newEvent)
                    })

                    bodyHTML.push(fragment)
                })

                GuJS.each(app.BodyResources, function(item) {

                    // Generate the fragment for the giving markup.
                    var fragment = GuJS.createDOMFragment(item.Markup)

                    // Register all events for this markup.
                    GuJS.each(item.Events, function(event) {
                        var newEvent = {}
                        newEvent.Event = event
                        newEvent.Callback = GuJS.MakeEventCallback(body, event)

                        body.addEventListener(event.Event, newEvent.Callback, event.UseCapture);
                        appEvents.base.bodyEvents.push(newEvent);
                    })


                    bodyHTML.push(fragment)
                })

                head.innerHTML = ""

                if (nonGuHead.length) {
                    GuJS.each(nonGuHead, function(item) {
                        head.appendChild(item)
                    })
                }

                if (headHTML.length) {
                    GuJS.each(headHTML, function(item) {
                        head.appendChild(item)
                    })
                }

                body.innerHTML = ""

                if (nonGuBody.length) {
                    GuJS.each(nonGuBody, function(item) {
                        body.appendChild(item)
                    })
                }

                if (bodyHTML.length) {
                    GuJS.each(bodyHTML, function(item) {
                        body.appendChild(item)
                    })
                }

                return

            case "RenderView":
                // Rendering the app response is to clear what is currently in the view.
                // We want specific replicate the way the gopherjs driver updates apps views
                // in swapping out the current content of the given view with the new content.

                var view = command.View

                // If the view is from a different app then don't service.
                // An App must be rendered before a view can be updated independently.
                if (GuJS.currentAppID && view.AppID !== GuJS.currentAppID) {
                    return
                }

                // Retrieve events map related to the giving app.
                var appEvents = GuJS.eventsCore[view.AppID] || { views: {}, base: { headEvents: [], bodyEvents: [] } }
                GuJS.eventsCore[view.AppID] = appEvents

                var viewEvents = appEvents.views[view.ViewID] || []
                appEvents.views[view.ViewID] = viewEvents


                // Deregister all view events.
                GuJS.each(viewEvents, function(view) {
                    GuJS.each(view, function(cb) {
                        body.removeEventListener(cb.Event.Event, cb.Callback)
                    })
                })

                var fragmentDOM = GuJS.createDOMFragment(view.Tree.Markup)
                GuJS.PatchDOM(fragmentDOM, body, false)

                // Register all events for this markup.
                GuJS.each(view.Tree.Events, function(event) {
                    var newEvent = {}
                    newEvent.Event = event
                    newEvent.Callback = GuJS.MakeEventCallback(body, event)

                    body.addEventListener(event.Event, newEvent.Callback, event.UseCapture);
                    viewEvents.push(newEvent)
                })

                return

            default:
                console.log("Command not support: ", command);
        }
    };


    // GuJS.PatchDOM patches the provided elements into the target from the current DOM.
    // It crawls a liveDOM version of the DOM, removing, replacing and adding node
    // changes as needed, until the dom resembles it's shadow/fragmentDOM.
    GuJS.PatchDOM = function(fragmentDOM, liveDOM, replace) {
        if (!liveDOM.hasChildNodes()) {
            liveDOM.appendChild(fragmentDOM)
            return
        }

        var shadowNodes = fragmentDOM.childNodes || []
        var liveNodes = liveDOM.childNodes || []

        for (var index = 0; index < shadowNodes.length; index++) {
            var node = shadowNodes[index]

            if (node.constructor === Text) {
                if (GuJS.isEmptyTextNode(node)) {
                    liveDOM.appendChild(node)
                    continue
                }


                if (index < liveNodes.length) {
                    liveNode = liveNodes[index]
                    liveDOM.insertBefore(liveNode, node)
                    continue
                }

                liveDOM.appendChild(node)
                continue
            }

            var nodeTagName = node.tagName
            var nodeId = node.getAttribute("id")
                // var nodeClass = node.getAttribute("class")
            var nodeUID = node.getAttribute("uid")
            var nodeAttr = node.attributes
            var nodeKids = node.childNodes
            var nodeSel = nodeTagName + "[uid=" + nodeUID + "]"
            var nodeRemoved = node.hasAttribute("NodeRemoved")
            var nodeHash = node.getAttribute("hash")

            if (!nodeId && !nodeUID && !nodeHash) {
                GuJS.addIfNoEqual(liveDOM, node)
                continue
            }

            if (!nodeUID && !nodeHash) {
                if (nodeId) {
                    var found = liveDOM.querySelectorAll("#" + nodeId)
                    if (!found.length) {
                        liveDOM.appendChild(node)
                        continue
                    }

                    liveDOM.replaceNode(found, node)
                    continue
                }
            }

            var allTargets = liveDOM.querySelectorAll(nodeSel)
            if (!allTargets.length) {
                liveDOM.appendChild(node)
                continue
            }

            for (var jindex = 0; jindex < allTargets.length; jindex++) {
                var curTarget = allTargets[jindex]

                if (nodeRemoved) {
                    liveDOM.remove(curTarget)
                    continue
                }

                if (replace) {
                    liveDOM.replaceNode(curTarget, node)
                    continue
                }

                var liveHash = curTarget.getAttribute("hash")

                if (liveHash === nodeHash) {
                    continue
                }


                if (!curTarget.childNodes.length) {
                    liveDOM.replaceNode(curTarget, node)
                    continue
                }

                GuJS.removeAllTextNodes(curTarget)

                for (var key in nodeAttr) {
                    var attr = nodeAttr[key]
                    curTarget.setAttribute(attr.nodeName, attr.nodeValue)
                }

                curTargetChilds = curTarget.childNodes
                if (!curTargetChilds.length) {
                    curTarget.innerHTML = ""
                    curTarget.appendChild.apply(curTarget, nodeKids)
                    continue
                }

                GuJS.PatchDOM(node, curTarget, replace)
            }
        }
    }

    // addIfNoEqual adds a giving node into the target if its not found to match any
    // child nodes of the target and if one is found then that is replaced with the
    // provided new node.
    GuJS.addIfNoEqual = function(target, node) {
        var list = target.childNodes
        for (var i = 0; i < list.length; i++) {
            var against = list[i]
            if (against.isEqualNode(node)) {
                target.replaceNode(against, node)
                return
            }
        }

        target.appendChild(node)
    }

    // isEmptyTextNode returns true/false if the node is an empty text node.
    GuJS.isEmptyTextNode = function(node) {
        if (node.nodeType !== 3) {
            return false
        }

        return node.textContent === ""
    }

    // removeAllTextNodes removes all residing textnodes in the provided node.
    GuJS.removeAllTextNodes = function(parent) {
        var list = parent.childNodes

        for (var i = 0; i < list.length; i++) {
            var node = list[i]
            if (node.nodeType === 3) {
                parent.removeChild(node)
            }
        }
    }

    // createDOMFragment creates a DocumentFragment from the provided HTML.
    GuJS.createDOMFragment = function(elemString) {
        var div = document.createElement("div")
        div.innerHTML = elemString

        var fragment = document.createDocumentFragment()

        GuJS.each(div.childNodes, function(node) {
            nodeType = GuJS.Type(node)
            if (nodeType.match(/HTML|Node|Element|Document|Text/)) {
                fragment.appendChild(node)
            }
        })

        div = null

        return fragment
    }

    // GetEvent returns the event as a object which can be jsonified and
    // sent over the pipeline.
    GuJS.GetEvent = function(ev) {
        var eventObj

        var c = ev.constructor
        switch (c) {
            case MutationRecord:
                eventObj = GuJS.DeepClone(ev, {
                    Functions: false,
                })

                var added = GuJS.map(eventObj.addedNodes, function(elem) {
                    return GuJS.StringifyHTML(elem)
                })

                var removed = GuJS.map(eventObj.removedNodes, function(elem) {
                    return GuJS.StringifyHTML(elem)
                })

                var presib = GuJS.StringifyHTML(eventObj.preSibling)
                var nextsib = GuJS.StringifyHTML(eventObj.nextSibling)

                eventObj.AddedNodes = added
                eventObj.RemovedNodes = removed
                eventObj.PreSibling = presib
                eventObj.NextSibling = nextsib

            case MediaStream:
                eventObj = GuJS.toMediaStream(ev)

            default:
                eventObj = GuJS.DeepClone(ev, {
                    Functions: false,
                })
        }

        return eventObj
    }

    // GuJS.Type returns the type of the native constructor of the passed in object.
    GuJS.Type = function(item) {
        if (item !== undefined && item != null) {
            return (item.constructor.toString().match(/function (.*)\(/)[1])
        }
    }

    // filters out the giving items not matching the provided function.
    GuJS.filter = function(item, fn) {
        var filtered = []

        for (key in item) {
            if (fn(item[key], key, item)) {
                filtered.push(item[key])
            }
        }

        return filtered
    }

    // GuJS.map maps new values throug the provided GuJS.skipping null returns.
    GuJS.map = function(item, fn) {
        var mapped = []

        for (key in item) {
            var res = fn(item[key], key, item)
            if (res) {
                mapped.push(item[key])
            }
        }

        return mapped
    }

    // each runs through all items in the provided list.
    GuJS.each = function(list, fn) {
        if ('length' in list) {
            for (var i = 0; i < list.length; i++) {
                fn(list[i], i, list)
            }

            return
        }

        for (key in list) {
            fn(list[key], key, list)
        }
    }

    // GuJS.reverse returns the list reversed in order.
    GuJS.reverse = function(list) {
        var reversed = []

        for (var i = list.length - 1; i > 0; i--) {
            reversed.push(list[i])
        }

        return reversed
    }

    // GuJS.mapFlattern maps new values throug the provided GuJS.skipping null returns.
    GuJS.mapFlattern = function(item, fn) {
        var mapped = []

        for (key in item) {
            var res = fn(item[key], key, item)
            if (!res) {
                continue
            }

            switch (GuJS.Type(res)) {
                case "Array":
                    Array.prototype.push.apply(mapped, res)
                    break
                default:
                    mapped.push(item[key])
            }
        }

        return mapped
    }

    // GuJS.capitalize returns a capitalized string.
    GuJS.capitalize = function(val) {
        if (val !== "") {
            var newVal = [val[0].toUpperCase()]
            newVal.push(val.substring(1))
            return newVal.join('')
        }

        return val
    }

    // GuJS.isUpperCase returns true/false if the string is in uppercase.
    GuJS.isUpperCase = function(val) {
        return val.toUpperCase() === val
    }

    // GuJS.Keys returns the constructor keys for the giving object.
    GuJS.Keys = function(item) {
        // If we can use the getOwnPropertyNames GuJS.in ES5 then use this has
        // inherited properties are desired as well.
        if ("getOwnPropertyNames" in Object) {
            return Object.getOwnPropertyNames(item)
        }

        // If we can use the Object.keys GuJS.in ES5 then use this has
        // we can manage with the provided set.
        if ("keys" in Object) {
            return Object.keys(item)
        }

        var keys = []

        for (var key in item) {
            keys.push(key)
        }

        // Check if keys are empty and if not, probably declared object
        // returned.
        if (keys.length) {
            return keys
        }

        // Attempt using the __proto__ object if we can copy. We are probably back in
        // Old JS land.
        if (item.__proto__) {
            for (var key in item.__proto__) {
                keys.push(key)
            }

            return keys
        }

        // Attempt using the protototype object if we can copy. We are probably back in
        // Old JS land.
        if (item.prototype) {
            for (var key in item.prototype) {
                keys.push(key)
            }

            return keys
        }

        // Digress to access prototype from constructor and
        // using the protototype object if we can copy. We are probably back in
        // Old JS land.
        if (item.constructor.prototype) {
            for (var key in item.constructor.prototype) {
                keys.push(key)
            }

            return keys
        }

        return keys
    }

    // exceptObjects are objects which we dont want uncloned but kept intact.
    // Also, these elements will lead to massive cyclic issues, keep them intact and deal
    // in another approach.
    var exceptObjects = { HTMLElement: true, NodeList: true, HTMLDocument: true, Node: true, Document: true }

    // defaultOptions defines a set of optional values allowed when cloning objects.
    var defaultOptions = { Functions: true, LastTree: [] }

    // GuJS.DeepClone clones all internal properties of the provided object, re-creating
    // internal key-value pairs accessible to the object even in prototype inheritance.
    // Functions are not runned except for custom types which are checked accordingly.
    GuJS.DeepClone = function(item, options) {
        if (item === undefined || item == null) {
            return item
        }

        c = item.constructor

        if (!options) {
            options = defaultOptions
        }

        switch (c) {
            case Function:
                if (options.AllowFunctions) {
                    return item
                }

                return

            case Number:
                return item

            case HTMLElement:
                return item

            case Node:
                return item

            case Element:
                return item

            case Boolean:
                return item

            case String:
                return item

            case Blob:
                return GuJS.fromBlob(item)

            case File:
                return GuJS.fromFile(item)

            case Uint8Array:
                var newArray = new Uint8Array
                for (var index in item) {
                    newArray.push(item[index])
                }

                return newArray

            case Float64Array:
                var newArray = new Float32Array
                for (var index in item) {
                    newArray.push(item[index])
                }

                return newArray

            case Float32Array:
                var newArray = new Float32Array
                for (var index in item) {
                    newArray.push(item[index])
                }

                return newArray

            case TouchList:
                return GuJS.toTouches(item)

            case MediaStream:
                return GuJS.toMediaStream(item)

            case Gamepad:
                return GuJS.toGamepad(item)

            case DataTransfer:
                return GuJS.toDataTransfer(item)

            case Array:
                var newArray = []
                for (var index in item) {
                    indexElem = item[index]
                    newArray[index] = GuJS.DeepClone(indexElem, options)
                }

                return newArray

            default:
                var newObj = {}
                var roots = GuJS.GetRoots(item)

                // If the element is a child of this givng root in the exceptObjects
                // then return it has is, because we need it intact and unchanged.
                for (var root in roots) {
                    var base = roots[root]
                    if (exceptObjects[GuJS.Type(base)]) {
                        return item
                    }
                }

                // If we are passed the previous tree, then check if we have
                // someone in that root as well.
                // if(options.LastTree){
                //   for(var root in options.LastTree){
                //     var base = roots[root]
                //     if(exceptObjects[GuJS.Type(base)]){
                //       return item
                //     }
                //   }
                // }

                var rootProtos = GuJS.reverse(GuJS.filter(roots, function(val) {
                    return GuJS.Type(val) != "Object"
                }))

                // Are we dealing with a empty object without parent, then we are probably
                // dealing with a declared GuJS.map/hash that points directly to the Object constructor.
                if (rootProtos.length === 0) {
                    rootProtos.push(item)
                }

                // Run through all parent constructs and pull keys, we want
                // to have all inherited properties as well.
                var keys = GuJS.mapFlattern(rootProtos, function(root) {
                    return GuJS.filter(GuJS.Keys(root), function(val) {

                        // If functions are not allowed and we have on here, then skip.
                        if (!options.AllowFunctions && GuJS.Type(item[val]) === "Function") {
                            return false
                        }

                        var allowed = !unwanted[val]
                        var isNotConstant = !(GuJS.isUpperCase(val))

                        return allowed && isNotConstant
                    })
                })

                for (var index in keys) {
                    var key = keys[index]
                    newObj[GuJS.capitalize(key)] = GuJS.DeepClone(item[key], {
                        Functions: options.Functions,
                        // LastTree: roots,
                    })
                }

                return newObj
        }
    }

    // GuJS.StringifyHTML returns the html of the element and it's content.
    GuJS.StringifyHTML = function(elem, deep) {
        var div = document.createElement("div")
        div.appendChild(elem.cloneNode(deep))
        return div.innerHTML
    }

    // GuJS.GetRoots retrieves all root properties for which the element runs down.
    GuJS.GetRoots = function(o) {
        var roots = []
        var found = {}

        var proto = o.constructor.prototype

        while (true) {
            if (proto == undefined || proto == null) {
                break
            }

            if (found[proto]) {
                break
            }

            roots.push(proto)
            found[proto] = true

            if ("__proto__" in proto) {
                proto = proto.__proto__
            }
        }

        return roots
    }

    // GuJS.fromBlob transform the providded Object blob into a byte slice.
    GuJS.fromBlob = function(o) {
        if (o == null || o == undefined) {
            return
        }

        var data = null

        fileReader = new FileReader()
        fileReader.onloadend = function() {
            data = new Uint8Array(fileReader.result)
        }

        fileReader.readAsArrayBuffer(o)

        return data
    }

    // GuJS.fromFile transform the providded Object blob into a byte slice.
    GuJS.fromFile = function(o) {
        if (o == null || o == undefined) {
            return
        }

        var data = null

        fileReader = new FileReader()
        fileReader.onloadend = function() {
            data = new Uint8Array(fileReader.result)
        }

        fileReader.readAsArrayBuffer(o)

        return data
    }

    // toInputSourceCapability returns the events.InputDeviceCapabilities from the object.
    GuJS.toInputSourceCapability = function(o) {
        if (o == null || o == undefined) {
            return
        }

        return {
            FiresTouchEvent: o.firesTouchEvent,
        }
    }

    // GuJS.toMotionData returns a motionData object from the Object.
    GuJS.toMotionData = function(o) {
        var md = { X: 0.0, Y: 0.0, Z: 0.0 }

        if (o == null || o == undefined) {
            return md
        }

        md.X = o.x
        md.Y = o.y
        md.Z = o.z
        return md
    }

    // GuJS.toRotationData returns a RotationData object from the Object.
    GuJS.toRotationData = function(o) {
        if (o == null || o == undefined) {
            return
        }

        md.Alpha = o.alpha
        md.Beta = o.beta
        md.Gamma = o.gamma
        return md
    }

    // GuJS.toMediaStream returns a events.MediaStream object.
    GuJS.toMediaStream = function(o) {
        if (o == null || o == undefined) {
            return
        }

        stream.Active = o.active
        stream.Ended = o.ended
        stream.ID = o.id
        stream.Audios = []
        stream.Videos = []

        var audioTracks = o.getAudioTracks()
        if (audioTracks != null && audioTracks != undefined) {
            for (i = 0; i < audioTracks.length; i++) {
                var track = audioTracks[i]
                var settings = track.getSettings()

                stream.Audios.push({
                    Enabled: track.enabled,
                    ID: track.id,
                    Kind: track.kind,
                    Label: track.label,
                    Muted: track.muted,
                    ReadyState: track.readyState,
                    Remote: track.remote,
                    AudioSettings: {
                        ChannelCount: settings.channelCount.Int(),
                        EchoCancellation: settings.echoCancellation,
                        Latency: settings.latency,
                        SampleRate: settings.sampleRate.Int64(),
                        SampleSize: settings.sampleSize.Int64(),
                        Volume: settings.volume,
                        MediaTrackSettings: {
                            DeviceID: settings.deviceId,
                            GroupID: settings.groupId,
                        },
                    },
                })
            }
        }

        var videosTracks = o.getVideoTracks()
        if (videosTracks != null && videosTracks != undefined) {
            for (i = 0; i < videosTracks.length; i++) {
                var track = videosTracks[i]
                var settings = track.getSettings()

                stream.Videos.push({
                    Enabled: track.enabled,
                    ID: track.id,
                    Kind: track.kind,
                    Label: track.label,
                    Muted: track.muted,
                    ReadyState: track.readyState,
                    Remote: track.remote,
                    VideoSettings: {
                        AspectRatio: settings.aspectRation,
                        FrameRate: settings.frameRate,
                        Height: settings.height.Int64(),
                        Width: settings.width.Int64(),
                        FacingMode: settings.facingMode,
                        MediaTrackSettings: {
                            DeviceID: settings.deviceId,
                            GroupID: settings.groupId,
                        },
                    },
                })
            }
        }

        return stream
    }

    GuJS.toTouches = function(o) {
        if (o == null || o == undefined) {
            return
        }

        var touches = []

        for (i = 0; i < o.length; i++) {
            var ev = o.item(i)
            touches.push({
                ClientX: ev.clientX,
                ClientY: ev.clientY,
                OffsetX: ev.offsetX,
                OffsetY: ev.offsetY,
                PageX: ev.pageX,
                PageY: ev.pageY,
                ScreenX: ev.screenX,
                ScreenY: ev.screenY,
                Identifier: ev.identifier,
            })

        }

        return touches
    }

    // toGamepad returns a Gamepad struct from the js object.
    GuJS.toGamepad = function(o) {
        var pad = {}

        if (o == null || o == undefined) {
            return pad
        }

        pad.DisplayID = o.displayId
        pad.ID = o.id
        pad.Index = o.index.Int()
        pad.Mapping = o.mapping
        pad.Connected = o.connected
        pad.Timestamp = o.timestamp
        pad.Axes = []
        pad.Buttons = []

        var axes = o.axes
        if (axes != null && axes != undefined) {
            for (i = 0; i < axes.length; i++) {
                pad.Axes.push(axes[i])
            }
        }

        var buttons = o.buttons
        if (buttons != null && buttons != undefined) {
            for (i = 0; i < buttons.length; i++) {
                button = buttons[i]
                pad.Buttons.push({
                    Value: button.value,
                    Pressed: button.pressed,
                })
            }
        }

        return pad
    }

    // toDataTransfer returns a transfer object from the Object.
    GuJS.toDataTransfer = function(o) {
        if (o == null || o == undefined) {
            return
        }

        var dt = {}
        dt.DropEffect = o.dropEffect
        dt.EffectAllowed = o.effectAllowed
        df.Types = o.types
        df.Items = []

        var items = o.items
        if (items != null && items != undefined) {
            for (i = 0; i < items.length; i++) {
                item = items.DataTransferItem(i)
                dItems.push({
                    Name: item.name,
                    Size: item.size.Int(),
                    Data: GuJS.fromFile(item),
                })
            }
        }

        var dFiles = []

        files = o.files
        if (files != null && files != undefined) {
            for (i = 0; i < files.length; i++) {
                item = files[i]
                dFiles.push({
                    Name: item.name,
                    Size: item.size.Int(),
                    Data: GuJS.fromFile(item),
                })
            }
        }

        dt.Items = { Items: dItems }
        dt.Files = dFiles
        return dt
    }


    onMessages(GuJS.ExecuteCommand)
}