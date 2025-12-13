package selfwatch

/*
#include <stdio.h>

#include <X11/Xlib.h>
#include <X11/extensions/record.h>
#include <X11/extensions/XTest.h>

#include <X11/Xlibint.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/cursorfont.h>
#include <X11/keysymdef.h>
#include <X11/keysym.h>


void eventCallbackGo(int type, int code);

void event_callback_cgo(XPointer priv, XRecordInterceptData *hook) {
	if (hook->category != XRecordFromServer) {
		XRecordFreeData(hook);
		return;
	}

	xEvent *event = (xEvent*)hook->data;
	int type = event->u.u.type;
	int code = event->u.u.detail;
	XRecordFreeData(hook);
	eventCallbackGo(type, code);
}
*/
import "C"
