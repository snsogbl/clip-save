//go:build darwin
// +build darwin

package common

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AppKit
// Implemented in statusbar_darwin.m
void setActivationPolicy(int policy);
*/
import "C"

// SetDockIconVisibility 设置 Dock 图标可见性（仅 macOS）
func SetDockIconVisibility(visible int) {
	// 调用 Cgo 函数
	C.setActivationPolicy(C.int(visible))
}
