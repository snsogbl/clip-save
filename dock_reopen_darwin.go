//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -ObjC -framework Cocoa
// Implemented in dock_reopen_darwin.m
void RegisterReopenObserver(void);
// Forward declaration to Go callback
extern void goOnAppReopen(void);
*/
import "C"

//export goOnAppReopen
func goOnAppReopen() {
	if dockReopenCallback != nil {
		dockReopenCallback()
	}
}

var dockReopenCallback func()

func initDockReopen(callback func()) {
	dockReopenCallback = callback
	C.RegisterReopenObserver()
}
