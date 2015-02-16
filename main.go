package main

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
	"unsafe"
)

// converting pointer to slice
// extensionsSlice := (*[1 << 30]*C.char)(unsafe.Pointer(extensions))[:numExtensions:numExtensions]

func queryExtension(display *C.Display, name string) bool {
	var major C.int
	var firstEvent C.int
	var firstError C.int
	strRecord := C.CString(name)
	defer C.free(unsafe.Pointer(strRecord))
	res := C.XQueryExtension(display, strRecord, &major, &firstEvent, &firstError)
	return 1 == int(res)
}

//export eventCallbackGo
func eventCallbackGo(eventType C.int, code C.int) {
	switch eventType {
	case C.KeyPress:
		fmt.Println("KeyPress", code)
	case C.KeyRelease:
		fmt.Println("KeyRelease", code)
	case C.ButtonPress:
		fmt.Println("ButtonPress", code)
	case C.ButtonRelease:
		fmt.Println("ButtonRelease", code)
	}
}

func main() {
	dataDisplay := C.XOpenDisplay(nil)
	controlDisplay := C.XOpenDisplay(nil)

	if dataDisplay == nil {
		log.Fatal("Failed to open display")
	}

	defer C.XCloseDisplay(dataDisplay)
	defer C.XCloseDisplay(controlDisplay)

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
	fmt.Println("got to bottom..")
}
