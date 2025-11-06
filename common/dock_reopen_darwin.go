//go:build darwin

package common

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -ObjC -framework Cocoa
// Implemented in dock_reopen_darwin.m
void RegisterReopenObserver(void);
// Forward declaration to Go callback
extern void goOnAppReopen(void);
extern void goSetForceQuit(void);
*/
import "C"

var dockReopenCallback func()
var forceQuitCallback func()

//export goOnAppReopen
func goOnAppReopen() {
	if dockReopenCallback != nil {
		dockReopenCallback()
	}
}

//export goSetForceQuit
func goSetForceQuit() {
	if forceQuitCallback != nil {
		forceQuitCallback()
	}
}

// InitDockReopen 初始化 Dock 重新打开监听
func InitDockReopen(callback func()) {
	dockReopenCallback = callback
	C.RegisterReopenObserver()
}

// SetForceQuitCallback 设置强制退出回调
func SetForceQuitCallback(callback func()) {
	forceQuitCallback = callback
}
