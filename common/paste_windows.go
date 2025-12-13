//go:build windows
// +build windows

package common

import (
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var (
	// Windows API 函数
	procSetForegroundWindow_paste = modUser32.NewProc("SetForegroundWindow")
	procSendInput_paste           = modUser32.NewProc("SendInput")
	procIsWindow_paste            = modUser32.NewProc("IsWindow")
	procEnumWindows_paste         = modUser32.NewProc("EnumWindows")
	procShowWindow_paste          = modUser32.NewProc("ShowWindow")
	procBringWindowToTop_paste    = modUser32.NewProc("BringWindowToTop")
	procAttachThreadInput_paste   = modUser32.NewProc("AttachThreadInput")
	procGetCurrentThreadId_paste  = modKernel32.NewProc("GetCurrentThreadId")
	procSetActiveWindow_paste     = modUser32.NewProc("SetActiveWindow")
	procGetCurrentProcessId_paste = modKernel32.NewProc("GetCurrentProcessId")
)

const (
	// 键盘输入相关常量
	INPUT_KEYBOARD_PASTE  = 1
	KEYEVENTF_KEYUP_PASTE = 0x0002
	VK_CONTROL_PASTE      = 0x11
	VK_V_PASTE            = 0x56
	SW_RESTORE_PASTE      = 9
	SW_SHOW_PASTE         = 5
)

type INPUT_PASTE struct {
	Type uint32
	Ki   KEYBDINPUT_PASTE
	_    [8]byte // padding for union
}

type KEYBDINPUT_PASTE struct {
	Wvk         uint16
	Wscan       uint16
	Dwflags     uint32
	Time        uint32
	DwextraInfo uintptr
}

// 全局变量
var (
	previousWindow_paste   uintptr
	currentProcessId_paste uint32
	windowMutex            sync.Mutex
)

// InitAppSwitchListener 初始化应用切换监听器（Windows版本）
func InitAppSwitchListener() {
	// 获取当前进程ID
	pid, _, _ := procGetCurrentProcessId_paste.Call()
	currentProcessId_paste = uint32(pid)

	// 记录当前前台窗口（如果不是我们的应用）
	recordCurrentForegroundWindow()

	// 启动后台窗口监听器
	startWindowMonitor()
}

// RecordPreviousAppPID 记录当前前台应用（在激活本应用之前调用）
func RecordPreviousAppPID() {
	recordCurrentForegroundWindow()
}

// recordCurrentForegroundWindow 记录当前前台窗口
func recordCurrentForegroundWindow() {
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd != 0 {
		var processId uint32
		procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&processId)))

		// 如果不是我们的进程，记录这个窗口
		if processId != currentProcessId_paste {
			previousWindow_paste = hwnd
		}
	}
}

// startWindowMonitor 启动窗口监听器
func startWindowMonitor() {
	go func() {
		var lastWindow uintptr
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			current, _, _ := procGetForegroundWindow.Call()
			if current != 0 && current != lastWindow {
				var processId uint32
				procGetWindowThreadProcessId.Call(current, uintptr(unsafe.Pointer(&processId)))

				// 如果不是我们的进程，记录这个窗口
				if processId != currentProcessId_paste {
					windowMutex.Lock()
					previousWindow_paste = current
					windowMutex.Unlock()
				}
				lastWindow = current
			}
		}
	}()
}

// ActivateAppByPID 通过 PID 激活指定应用（Windows版本）
func ActivateAppByPID(pid int) bool {
	targetPid := uint32(pid)
	var targetWindow uintptr

	// 枚举所有窗口，找到属于目标进程的窗口
	enumProc := syscall.NewCallback(func(hwnd uintptr, lParam uintptr) uintptr {
		var processId uint32
		procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&processId)))

		if processId == targetPid {
			isValid, _, _ := procIsWindow_paste.Call(hwnd)
			if isValid != 0 {
				targetWindow = hwnd
				return 0 // 停止枚举
			}
		}
		return 1 // 继续枚举
	})

	procEnumWindows_paste.Call(enumProc, 0)

	if targetWindow != 0 {
		return activateWindow_paste(targetWindow)
	}

	return false
}

// activateWindow_paste 激活指定窗口
func activateWindow_paste(hwnd uintptr) bool {
	// 检查窗口是否有效
	isValid, _, _ := procIsWindow_paste.Call(hwnd)
	if isValid == 0 {
		return false
	}

	// 标准激活流程
	procShowWindow_paste.Call(hwnd, SW_RESTORE_PASTE)
	time.Sleep(50 * time.Millisecond)

	procShowWindow_paste.Call(hwnd, SW_SHOW_PASTE)
	time.Sleep(50 * time.Millisecond)

	procBringWindowToTop_paste.Call(hwnd)
	time.Sleep(50 * time.Millisecond)

	result, _, _ := procSetForegroundWindow_paste.Call(hwnd)

	if result != 0 {
		return true
	}

	// 强制激活
	return forceActivateWindow_paste(hwnd)
}

// forceActivateWindow_paste 强制激活窗口
func forceActivateWindow_paste(hwnd uintptr) bool {
	// 获取目标窗口的线程ID
	var targetProcessId uint32
	targetThreadId, _, _ := procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&targetProcessId)))

	// 获取当前线程ID
	currentThreadId, _, _ := procGetCurrentThreadId_paste.Call()

	if targetThreadId != currentThreadId {
		// 附加到目标线程的输入队列
		procAttachThreadInput_paste.Call(currentThreadId, targetThreadId, 1)
	}

	// 尝试激活
	procSetActiveWindow_paste.Call(hwnd)
	result, _, _ := procSetForegroundWindow_paste.Call(hwnd)

	// 分离线程输入队列
	if targetThreadId != currentThreadId {
		procAttachThreadInput_paste.Call(currentThreadId, targetThreadId, 0)
	}

	return result != 0
}

// PasteCmdV 发送 Ctrl+V 粘贴命令（Windows版本）
func PasteCmdV() {
	// 线程安全地读取目标窗口
	windowMutex.Lock()
	targetWindow := previousWindow_paste
	windowMutex.Unlock()

	// 如果有记录的目标窗口，确保激活它
	if targetWindow != 0 {
		isValid, _, _ := procIsWindow_paste.Call(targetWindow)
		if isValid != 0 {
			activateWindow_paste(targetWindow)
			time.Sleep(200 * time.Millisecond)
		}
	}

	// 检查当前前台窗口
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd != 0 {
		var processId uint32
		procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&processId)))
		if processId == currentProcessId_paste {
			// 如果前台窗口是我们自己，尝试再次激活目标窗口
			if targetWindow != 0 {
				activateWindow_paste(targetWindow)
				time.Sleep(200 * time.Millisecond)
			}
		}
	}

	sendCtrlV_paste()
}

// sendCtrlV_paste 发送 Ctrl+V 组合键
func sendCtrlV_paste() {
	// 优先使用 SendInput
	if sendCtrlVWithSendInput_paste() {
		return
	}

	// SendInput 失败，使用 keybd_event 备用
	sendCtrlVWithKeybdEvent_paste()
}

// sendCtrlVWithSendInput_paste 使用 SendInput API 发送 Ctrl+V
func sendCtrlVWithSendInput_paste() bool {
	inputs := make([]INPUT_PASTE, 4)

	// Ctrl 按下
	inputs[0] = INPUT_PASTE{
		Type: INPUT_KEYBOARD_PASTE,
		Ki: KEYBDINPUT_PASTE{
			Wvk:         VK_CONTROL_PASTE,
			Wscan:       0,
			Dwflags:     0,
			Time:        0,
			DwextraInfo: 0,
		},
	}

	// V 按下
	inputs[1] = INPUT_PASTE{
		Type: INPUT_KEYBOARD_PASTE,
		Ki: KEYBDINPUT_PASTE{
			Wvk:         VK_V_PASTE,
			Wscan:       0,
			Dwflags:     0,
			Time:        0,
			DwextraInfo: 0,
		},
	}

	// V 释放
	inputs[2] = INPUT_PASTE{
		Type: INPUT_KEYBOARD_PASTE,
		Ki: KEYBDINPUT_PASTE{
			Wvk:         VK_V_PASTE,
			Wscan:       0,
			Dwflags:     KEYEVENTF_KEYUP_PASTE,
			Time:        0,
			DwextraInfo: 0,
		},
	}

	// Ctrl 释放
	inputs[3] = INPUT_PASTE{
		Type: INPUT_KEYBOARD_PASTE,
		Ki: KEYBDINPUT_PASTE{
			Wvk:         VK_CONTROL_PASTE,
			Wscan:       0,
			Dwflags:     KEYEVENTF_KEYUP_PASTE,
			Time:        0,
			DwextraInfo: 0,
		},
	}

	// 发送输入事件
	result, _, _ := procSendInput_paste.Call(
		uintptr(4),
		uintptr(unsafe.Pointer(&inputs[0])),
		uintptr(unsafe.Sizeof(inputs[0])),
	)

	return result == 4
}

// sendCtrlVWithKeybdEvent_paste 使用 keybd_event API 发送 Ctrl+V
func sendCtrlVWithKeybdEvent_paste() {
	procKeybdEvent := modUser32.NewProc("keybd_event")

	// Ctrl 按下
	procKeybdEvent.Call(VK_CONTROL_PASTE, 0, 0, 0)
	time.Sleep(50 * time.Millisecond)

	// V 按下
	procKeybdEvent.Call(VK_V_PASTE, 0, 0, 0)
	time.Sleep(50 * time.Millisecond)

	// V 释放
	procKeybdEvent.Call(VK_V_PASTE, 0, KEYEVENTF_KEYUP_PASTE, 0)
	time.Sleep(50 * time.Millisecond)

	// Ctrl 释放
	procKeybdEvent.Call(VK_CONTROL_PASTE, 0, KEYEVENTF_KEYUP_PASTE, 0)
}

// ActivatePreviousApp 激活之前的应用
func ActivatePreviousApp() {
	windowMutex.Lock()
	targetWindow := previousWindow_paste
	windowMutex.Unlock()

	if targetWindow != 0 {
		isValid, _, _ := procIsWindow_paste.Call(targetWindow)
		if isValid != 0 {
			// 多次尝试激活
			for i := 0; i < 3; i++ {
				success := activateWindow_paste(targetWindow)
				if success {
					break
				}
				time.Sleep(200 * time.Millisecond)
			}
		} else {
			windowMutex.Lock()
			previousWindow_paste = 0
			windowMutex.Unlock()
		}
	}
}

// PasteCmdVToPreviousApp 快速激活应用并发送 Ctrl+V（Windows版本）
func PasteCmdVToPreviousApp() {
	windowMutex.Lock()
	targetWindow := previousWindow_paste
	windowMutex.Unlock()

	// 方法1：使用记录的窗口句柄
	if targetWindow != 0 {
		isValid, _, _ := procIsWindow_paste.Call(targetWindow)
		if isValid != 0 {
			success := pasteToWindow_paste(targetWindow)
			if success {
				return
			}
		} else {
			windowMutex.Lock()
			previousWindow_paste = 0
			windowMutex.Unlock()
		}
	}

	// 方法2：查找其他窗口
	otherWindow := findAnyOtherWindow()
	if otherWindow != 0 {
		success := pasteToWindow_paste(otherWindow)
		if success {
			return
		}
	}

	// 方法3：当前前台窗口
	currentForeground, _, _ := procGetForegroundWindow.Call()
	if currentForeground != 0 && currentForeground != targetWindow {
		var processId uint32
		procGetWindowThreadProcessId.Call(currentForeground, uintptr(unsafe.Pointer(&processId)))

		if processId != currentProcessId_paste {
			success := pasteToWindow_paste(currentForeground)
			if success {
				return
			}
		}
	}

	// 方法4：全局发送
	time.Sleep(150 * time.Millisecond)
	sendCtrlV_paste()
}

// pasteToWindow_paste 向指定窗口发送粘贴命令
func pasteToWindow_paste(hwnd uintptr) bool {
	isValid, _, _ := procIsWindow_paste.Call(hwnd)
	if isValid == 0 {
		return false
	}

	success := activateWindow_paste(hwnd)
	if !success {
		return false
	}

	time.Sleep(100 * time.Millisecond)
	sendCtrlV_paste()

	return true
}

// findAnyOtherWindow 查找任何其他可见的窗口
func findAnyOtherWindow() uintptr {
	var foundWindow uintptr

	enumProc := syscall.NewCallback(func(hwnd uintptr, lParam uintptr) uintptr {
		// 检查窗口是否可见
		procIsWindowVisible := modUser32.NewProc("IsWindowVisible")
		visible, _, _ := procIsWindowVisible.Call(hwnd)
		if visible == 0 {
			return 1
		}

		// 获取窗口进程ID
		var processId uint32
		procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&processId)))

		// 跳过自己的进程
		if processId == currentProcessId_paste {
			return 1
		}

		// 检查窗口是否有标题
		windowTitle := getWindowTitle(hwnd)
		if len(windowTitle) > 0 {
			foundWindow = hwnd
			return 0 // 停止枚举
		}

		return 1
	})

	procEnumWindows_paste.Call(enumProc, 0)
	return foundWindow
}

// getWindowTitle 获取窗口标题
func getWindowTitle(hwnd uintptr) string {
	procGetWindowText := modUser32.NewProc("GetWindowTextW")

	buf := make([]uint16, 256)
	length, _, _ := procGetWindowText.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))

	if length > 0 {
		return syscall.UTF16ToString(buf[:length])
	}

	return ""
}
