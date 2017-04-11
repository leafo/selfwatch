package selfwatch

/*
#cgo LDFLAGS: -lXext -lX11 -lXtst
#include <stdlib.h>
#include <stdio.h>
#include <X11/Xlib.h>
#include <X11/extensions/record.h>
#include <X11/extensions/XTest.h>

void event_callback_cgo(XPointer priv, XRecordInterceptData *hook);
*/
import "C"

import (
	"fmt"
	"log"
	"time"
	"unsafe"
)

var instance *Recorder

type Recorder struct {
	KeyPress      func(code int32)
	KeyRelease    func(code int32)
	ButtonPress   func(code int32)
	ButtonRelease func(code int32)
	display       *C.Display
}

func NewRecorder() *Recorder {
	if instance != nil {
		log.Fatal("recorder already exists")
	}

	instance = &Recorder{}
	return instance
}

func (recorder *Recorder) Bind() error {
	dataDisplay := C.XOpenDisplay(nil)
	controlDisplay := C.XOpenDisplay(nil)

	for dataDisplay == nil {
		log.Print("Failed to open display, trying again in 10s")
		time.Sleep(time.Second * 10)
		dataDisplay = C.XOpenDisplay(nil)
		controlDisplay = C.XOpenDisplay(nil)
	}

	defer C.XCloseDisplay(dataDisplay)
	defer C.XCloseDisplay(controlDisplay)

	recorder.display = controlDisplay

	C.XSynchronize(controlDisplay, 1)

	if !queryExtension(dataDisplay, "RECORD") {
		log.Fatal("RECORD extension not present")
	}

	rr := C.XRecordAllocRange()
	if rr == nil {
		log.Fatal("XRecordAllocRange failed")
	}

	rr.device_events.first = C.KeyPress
	rr.device_events.last = C.MotionNotify
	rcs := C.XRecordAllClients

	rc := C.XRecordCreateContext(controlDisplay, 0, (*C.XRecordClientSpec)(unsafe.Pointer(&rcs)), 1, &rr, 1)

	if int(rc) == 0 {
		log.Fatal("XRecordCreateContext failed")
	}

	C.XRecordEnableContext(dataDisplay, rc, (C.XRecordInterceptProc)(unsafe.Pointer(C.event_callback_cgo)), nil)
	return nil
}

//export eventCallbackGo
func eventCallbackGo(eventType C.int, code C.int) {
	if instance == nil {
		return
	}

	switch eventType {
	case C.KeyPress:
		fmt.Println("KeyPress", code)
		if instance.KeyPress != nil {
			instance.KeyPress(int32(code))
		}
	case C.KeyRelease:
		fmt.Println("KeyRelease", code)
		if instance.KeyRelease != nil {
			instance.KeyRelease(int32(code))
		}

		window := instance.GetInputFocus()
		instance.ListProperties(window)

	case C.ButtonPress:
		fmt.Println("ButtonPress", code)
		if instance.ButtonPress != nil {
			instance.ButtonPress(int32(code))
		}
	case C.ButtonRelease:
		fmt.Println("ButtonRelease", code)
		if instance.ButtonRelease != nil {
			instance.ButtonRelease(int32(code))
		}
	}
}

func queryExtension(display *C.Display, name string) bool {
	var major C.int
	var firstEvent C.int
	var firstError C.int
	strRecord := C.CString(name)
	defer C.free(unsafe.Pointer(strRecord))
	res := C.XQueryExtension(display, strRecord, &major, &firstEvent, &firstError)
	return 1 == int(res)
}

func (r *Recorder) GetInputFocus() C.Window {
	var window C.Window
	var revert C.int
	C.XGetInputFocus(r.display, &window, &revert)
	return window
}

func (r *Recorder) GetWindowAttributes(window C.Window) {
	var attributes C.XWindowAttributes
	C.XGetWindowAttributes(r.display, window, &attributes)
	fmt.Println(attributes)
}

func (r *Recorder) ListProperties(window C.Window) {
	var numProperties C.int
	properties := C.XListProperties(r.display, window, &numProperties)
	fmt.Println(numProperties, properties)
}
