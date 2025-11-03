//go:build darwin
// +build darwin

package common

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework ApplicationServices
#include <ApplicationServices/ApplicationServices.h>
#include <unistd.h>

static void sendCmdV() {
    const CGKeyCode keyV = (CGKeyCode)9;

    CGEventRef vdown = CGEventCreateKeyboardEvent(NULL, keyV, true);
    if (vdown != NULL) {
        CGEventSetFlags(vdown, kCGEventFlagMaskCommand);
        CGEventPost(kCGHIDEventTap, vdown);
        CFRelease(vdown);
    }

    // 稍作延迟，确保按下事件被处理
    usleep(15000); // 15ms

    CGEventRef vup = CGEventCreateKeyboardEvent(NULL, keyV, false);
    if (vup != NULL) {
        CGEventSetFlags(vup, kCGEventFlagMaskCommand);
        CGEventPost(kCGHIDEventTap, vup);
        CFRelease(vup);
    }
}
*/
import "C"

import "time"

// PasteCmdV
func PasteCmdV() {
	// 等待前台焦点稳定（例如先隐藏本应用后）
	time.Sleep(120 * time.Millisecond)
	C.sendCmdV()
}
