//go:build windows
// +build windows

package common

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	// 重用已有的 DLL 引用，只声明新需要的 proc
	procSetForegroundWindow = modUser32.NewProc("SetForegroundWindow")
	procFindWindowW         = modUser32.NewProc("FindWindowW")
	procSendInput           = modUser32.NewProc("SendInput")
	procGetCurrentProcessId = modKernel32.NewProc("GetCurrentProcessId")
	procIsWindow            = modUser32.NewProc("IsWindow")
	procGetWindowTextW      = modUser32.NewProc("GetWindowTextW")
	procEnumWindows         = modUser32.NewProc("EnumWindows")
	procShowWindow          = modUser32.NewProc("ShowWindow")
	procBringWindowToTop    = modUser32.NewProc("BringWindowToTop")
)

const (
	INPUT_KEYBOARD            = 1
	KEYEVENTF_KEYUP           = 0x0002
	VK_CONTROL                = 0x11
	VK_V                      = 0x56
	PROCESS_QUERY_INFORMATION = 0x0400
	SW_RESTORE                = 9
	SW_SHOW                   = 5
)

type INPUT struct {
	Type uint32
	Ki   KEYBDINPUT
	_    [8]byte // padding for union
}

type KEYBDINPUT struct {
	Wvk         uint16
	Wscan       uint16
	Dwflags     uint32
	Time        uint32
	DwextraInfo uintptr
}

// 全局变量记录之前的前台窗口
var previousWindow uintptr
var currentProcessId uint32

// InitAppSwitchListener 初始化应用切换监听器（Windows版本）
func InitAppSwitchListener() {
	// 获取当前进程ID
	pid, _, _ := procGetCurrentProcessId.Call()
	currentProcessId = uint32(pid)

	// 记录当前前台窗口（如果不是我们的应用）
	recordCurrentForegroundWindow()
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
		if processId != currentProcessId {
			previousWindow = hwnd
		}
	}
}

// ActivateAppByPID 通过 PID 激活指定应用（Windows版本）
func ActivateAppByPID(pid int) bool {
	// 在Windows中，我们需要通过枚举窗口来找到对应进程的窗口
	targetPid := uint32(pid)
	var targetWindow uintptr

	// 枚举所有窗口，找到属于目标进程的窗口
	enumProc := syscall.NewCallback(func(hwnd uintptr, lParam uintptr) uintptr {
		var processId uint32
		procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&processId)))

		if processId == targetPid {
			// 检查窗口是否有效
			isValid, _, _ := procIsWindow.Call(hwnd)
			if isValid != 0 {
				targetWindow = hwnd
				return 0 // 停止枚举
			}
		}
		return 1 // 继续枚举
	})

	procEnumWindows.Call(enumProc, 0)

	if targetWindow != 0 {
		return activateWindow(targetWindow)
	}

	return false
}

// activateWindow 激活指定窗口
func activateWindow(hwnd uintptr) bool {
	// 先显示窗口
	procShowWindow.Call(hwnd, SW_RESTORE)
	procShowWindow.Call(hwnd, SW_SHOW)

	// 将窗口置于前台
	procBringWindowToTop.Call(hwnd)

	// 设置为前台窗口
	result, _, _ := procSetForegroundWindow.Call(hwnd)
	return result != 0
}

// PasteCmdV 发送 Ctrl+V 粘贴命令（Windows版本）
func PasteCmdV() {
	// 等待前台焦点稳定
	time.Sleep(120 * time.Millisecond)
	sendCtrlV()
}

// sendCtrlV 发送 Ctrl+V 组合键
func sendCtrlV() {
	inputs := make([]INPUT, 4)

	// Ctrl 按下
	inputs[0] = INPUT{
		Type: INPUT_KEYBOARD,
		Ki: KEYBDINPUT{
			Wvk: VK_CONTROL,
		},
	}

	// V 按下
	inputs[1] = INPUT{
		Type: INPUT_KEYBOARD,
		Ki: KEYBDINPUT{
			Wvk: VK_V,
		},
	}

	// V 释放
	inputs[2] = INPUT{
		Type: INPUT_KEYBOARD,
		Ki: KEYBDINPUT{
			Wvk:     VK_V,
			Dwflags: KEYEVENTF_KEYUP,
		},
	}

	// Ctrl 释放
	inputs[3] = INPUT{
		Type: INPUT_KEYBOARD,
		Ki: KEYBDINPUT{
			Wvk:     VK_CONTROL,
			Dwflags: KEYEVENTF_KEYUP,
		},
	}

	procSendInput.Call(4, uintptr(unsafe.Pointer(&inputs[0])), unsafe.Sizeof(inputs[0]))
}

// ActivatePreviousApp 激活之前的应用
func ActivatePreviousApp() {
	if previousWindow != 0 {
		// 检查窗口是否仍然有效
		isValid, _, _ := procIsWindow.Call(previousWindow)
		if isValid != 0 {
			activateWindow(previousWindow)
		} else {
			// 窗口已无效，清除记录
			previousWindow = 0
		}
	}
}

// PasteCmdVToPreviousApp 快速激活应用并发送 Ctrl+V（Windows版本）
func PasteCmdVToPreviousApp() {
	if previousWindow != 0 {
		// 检查窗口是否仍然有效
		isValid, _, _ := procIsWindow.Call(previousWindow)
		if isValid != 0 {
			// 激活之前的窗口
			success := activateWindow(previousWindow)
			if success {
				// 等待窗口激活
				time.Sleep(50 * time.Millisecond)
				// 发送 Ctrl+V
				sendCtrlV()
				return
			}
		} else {
			// 窗口已无效，清除记录
			previousWindow = 0
		}
	}

	// 如果找不到目标窗口，回退到全局发送
	time.Sleep(50 * time.Millisecond)
	sendCtrlV()
}
