//go:build darwin
// +build darwin

package common

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework ApplicationServices -framework Cocoa
#include <ApplicationServices/ApplicationServices.h>
#include <Cocoa/Cocoa.h>
#include <unistd.h>

// 全局变量记录之前的前台应用 PID
static pid_t previousAppPID = 0;

// 记录当前前台应用（在激活本应用之前调用）
static void recordPreviousAppPID() {
    NSWorkspace *workspace = [NSWorkspace sharedWorkspace];
    NSRunningApplication *currentApp = [NSRunningApplication currentApplication];
    NSRunningApplication *frontmostApp = [workspace frontmostApplication];

    if (frontmostApp != nil &&
        [frontmostApp processIdentifier] != [currentApp processIdentifier]) {
        previousAppPID = [frontmostApp processIdentifier];
    }
}

// 初始化应用切换监听器
static void setupAppSwitchListener() {
    NSWorkspace *workspace = [NSWorkspace sharedWorkspace];
    NSNotificationCenter *center = [workspace notificationCenter];
    NSRunningApplication *currentApp = [NSRunningApplication currentApplication];

    // 初始化：记录当前前台应用（如果不是我们的应用）
    NSRunningApplication *frontmostApp = [workspace frontmostApplication];
    if (frontmostApp != nil &&
        [frontmostApp processIdentifier] != [currentApp processIdentifier]) {
        previousAppPID = [frontmostApp processIdentifier];
    }

    // 监听应用激活通知
    [center addObserverForName:NSWorkspaceDidActivateApplicationNotification
                        object:nil
                         queue:[NSOperationQueue mainQueue]
                    usingBlock:^(NSNotification *note) {
        NSRunningApplication *app = [[note userInfo] objectForKey:NSWorkspaceApplicationKey];
        if (app != nil) {
            NSRunningApplication *currentApp = [NSRunningApplication currentApplication];
            // 如果不是我们的应用，记录为之前的前台应用
            if ([app processIdentifier] != [currentApp processIdentifier]) {
                previousAppPID = [app processIdentifier];
            }
        }
    }];
}

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

// 获取之前的前台应用的进程 ID
static pid_t getPreviousAppPID() {
    // 验证记录的 PID 是否有效
    if (previousAppPID != 0) {
        NSRunningApplication *app = [NSRunningApplication runningApplicationWithProcessIdentifier:previousAppPID];
        if (app != nil && ![app isTerminated]) {
            return previousAppPID;
        }
        // 如果进程已终止，清除记录
        previousAppPID = 0;
    }
    return 0;
}

// 发送 Cmd+V 到指定进程
static int sendCmdVToPID(pid_t pid) {
    if (pid == 0) {
        return 0;
    }

    const CGKeyCode keyV = (CGKeyCode)9;

    CGEventRef vdown = CGEventCreateKeyboardEvent(NULL, keyV, true);
    if (vdown != NULL) {
        CGEventSetFlags(vdown, kCGEventFlagMaskCommand);
        CGEventPostToPid(pid, vdown);
        CFRelease(vdown);
    }

    // 稍作延迟，确保按下事件被处理
    usleep(15000); // 15ms

    CGEventRef vup = CGEventCreateKeyboardEvent(NULL, keyV, false);
    if (vup != NULL) {
        CGEventSetFlags(vup, kCGEventFlagMaskCommand);
        CGEventPostToPid(pid, vup);
        CFRelease(vup);
    }

    return 1;
}

// 通过 PID 激活指定应用
static int activateAppByPID(pid_t pid) {
    if (pid == 0) {
        return 0;
    }

    NSRunningApplication *app = [NSRunningApplication runningApplicationWithProcessIdentifier:pid];
    if (app == nil || [app isTerminated]) {
        return 0;
    }

    BOOL success = [app activateWithOptions:NSApplicationActivateIgnoringOtherApps];
    return success ? 1 : 0;
}

// 激活自己的应用
static void activateSelf() {
    [[NSApplication sharedApplication] activateIgnoringOtherApps:YES];
}
*/
import "C"

import (
	"time"
)

// InitAppSwitchListener 初始化应用切换监听器（在应用启动时调用）
func InitAppSwitchListener() {
	C.setupAppSwitchListener()
}

// RecordPreviousAppPID 记录当前前台应用（在激活本应用之前调用）
func RecordPreviousAppPID() {
	C.recordPreviousAppPID()
}

// ActivateAppByPID 通过 PID 激活指定应用
func ActivateAppByPID(pid int) bool {
	result := C.activateAppByPID(C.pid_t(pid))
	return result == 1
}

// PasteCmdV 发送 Cmd+V 粘贴命令
func PasteCmdV() {
	// 等待前台焦点稳定（例如先隐藏本应用后）
	time.Sleep(120 * time.Millisecond)
	C.sendCmdV()
}

func ActivatePreviousApp() {
	// 获取之前的前台应用的进程 ID（优先使用记录的 PID）
	pid := C.getPreviousAppPID()
	if pid == 0 {
		return
	}

	// 先激活应用（必需，否则某些应用无法接收键盘事件）
	ActivateAppByPID(int(pid))
}

// PasteCmdVToPreviousApp 快速激活应用并发送 Cmd+V（优化等待时间，减少窗口切换感）
func PasteCmdVToPreviousApp() {
	// 获取之前的前台应用的进程 ID（优先使用记录的 PID）
	pid := C.getPreviousAppPID()
	if pid == 0 {
		// 如果找不到目标应用，回退到全局发送
		time.Sleep(50 * time.Millisecond)
		C.sendCmdV()
		return
	}

	// 先激活应用（必需，否则某些应用无法接收键盘事件）
	success := ActivateAppByPID(int(pid))
	if !success {
		// 如果激活失败，尝试直接发送到进程
		C.sendCmdVToPID(pid)
		return
	}

	// 发送 Cmd+V
	C.sendCmdV()

	// 粘贴完成后，激活本应用（置顶模式下需要）
	go func() {
		C.activateSelf()
	}()
}

// TestPasteFunction macOS 平台空实现（仅用于 Windows 调试）
func TestPasteFunction() {
	// macOS 不需要特殊的调试功能
}

// SimpleTestPaste macOS 平台空实现
func SimpleTestPaste() {
	// macOS 不需要特殊的调试功能
}

// TestLogWriting macOS 平台空实现
func TestLogWriting() {
	// macOS 不需要特殊的调试功能
}
