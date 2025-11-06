//go:build darwin
// +build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AppKit
// Implemented in statusbar_darwin.m
void setActivationPolicy(int policy);
*/
import "C"

func SetDockIconVisibility(visible int) {
	// 调用 Cgo 函数
	C.setActivationPolicy(C.int(visible))
}
