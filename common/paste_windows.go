//go:build windows
// +build windows

package common

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"
)

var (
	// 新增的 Windows API 函数（避免与现有文件冲突）
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
	pasteLogger            *log.Logger
)

// InitAppSwitchListener 初始化应用切换监听器（Windows版本）
func InitAppSwitchListener() {
	// 初始化日志文件
	initPasteLogger()

	// 获取当前进程ID
	pid, _, _ := procGetCurrentProcessId_paste.Call()
	currentProcessId_paste = uint32(pid)

	logPaste("初始化应用切换监听器，当前进程ID: %d", currentProcessId_paste)

	// 记录当前前台窗口（如果不是我们的应用）
	recordCurrentForegroundWindow()
}

// initPasteLogger 初始化粘贴功能的日志记录器
func initPasteLogger() {
	// 获取用户的临时目录
	tempDir := os.TempDir()
	logDir := filepath.Join(tempDir, "ClipSave")

	// 创建日志目录
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		log.Printf("创建日志目录失败: %v", err)
		return
	}

	// 创建日志文件（按日期命名）
	logFileName := fmt.Sprintf("paste_%s.log", time.Now().Format("2006-01-02"))
	logFilePath := filepath.Join(logDir, logFileName)

	// 打开或创建日志文件
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("创建日志文件失败: %v", err)
		return
	}

	// 创建多重输出：同时输出到控制台和文件
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	pasteLogger = log.New(multiWriter, "[PASTE] ", log.LstdFlags|log.Lshortfile)

	logPaste("=== 粘贴功能日志初始化完成 ===")
	logPaste("日志文件位置: %s", logFilePath)
	logPaste("提示: 可以在文件资源管理器中输入 %%TEMP%%\\ClipSave 查看日志文件")
}

// logPaste 记录粘贴相关的日志
func logPaste(format string, args ...any) {
	if pasteLogger != nil {
		pasteLogger.Printf(format, args...)
	} else {
		// 如果日志器未初始化，回退到标准日志
		log.Printf("[PASTE] "+format, args...)
	}
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

		logPaste("当前前台窗口: %x, 进程ID: %d, 我们的进程ID: %d", hwnd, processId, currentProcessId_paste)

		// 如果不是我们的进程，记录这个窗口
		if processId != currentProcessId_paste {
			previousWindow_paste = hwnd
			logPaste("记录前台窗口: %x", hwnd)
		}
	}
}

// ActivateAppByPID 通过 PID 激活指定应用（Windows版本）
func ActivateAppByPID(pid int) bool {
	logPaste("尝试激活进程ID为 %d 的应用", pid)

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
		logPaste("找到目标窗口: %x", targetWindow)
		return activateWindow_paste(targetWindow)
	}

	logPaste("未找到进程ID为 %d 的可见窗口", pid)
	return false
}

// activateWindow_paste 激活指定窗口
func activateWindow_paste(hwnd uintptr) bool {
	logPaste("尝试激活窗口: %x", hwnd)

	// 检查窗口是否有效
	isValid, _, _ := procIsWindow_paste.Call(hwnd)
	if isValid == 0 {
		logPaste("窗口无效: %x", hwnd)
		return false
	}

	// 先显示窗口
	procShowWindow_paste.Call(hwnd, SW_RESTORE_PASTE)
	procShowWindow_paste.Call(hwnd, SW_SHOW_PASTE)

	// 将窗口置于前台
	procBringWindowToTop_paste.Call(hwnd)

	// 设置为前台窗口
	result, _, _ := procSetForegroundWindow_paste.Call(hwnd)

	success := result != 0
	logPaste("窗口激活结果: %x, 成功: %v", hwnd, success)
	return success
}

// PasteCmdV 发送 Ctrl+V 粘贴命令（Windows版本）
func PasteCmdV() {
	logPaste("执行 PasteCmdV")
	time.Sleep(120 * time.Millisecond)
	sendCtrlV_paste()
}

// sendCtrlV_paste 发送 Ctrl+V 组合键
func sendCtrlV_paste() {
	logPaste("发送 Ctrl+V 组合键")

	// 方法1：使用 SendInput
	if sendCtrlVWithSendInput_paste() {
		return
	}

	// 方法2：回退到 keybd_event
	logPaste("SendInput 失败，尝试 keybd_event")
	sendCtrlVWithKeybdEvent_paste()
}

// sendCtrlVWithSendInput_paste 使用 SendInput API 发送 Ctrl+V
func sendCtrlVWithSendInput_paste() bool {
	inputs := make([]INPUT_PASTE, 4)

	// Ctrl 按下
	inputs[0] = INPUT_PASTE{
		Type: INPUT_KEYBOARD_PASTE,
		Ki: KEYBDINPUT_PASTE{
			Wvk:     VK_CONTROL_PASTE,
			Wscan:   0,
			Dwflags: 0,
		},
	}

	// V 按下
	inputs[1] = INPUT_PASTE{
		Type: INPUT_KEYBOARD_PASTE,
		Ki: KEYBDINPUT_PASTE{
			Wvk:     VK_V_PASTE,
			Wscan:   0,
			Dwflags: 0,
		},
	}

	// V 释放
	inputs[2] = INPUT_PASTE{
		Type: INPUT_KEYBOARD_PASTE,
		Ki: KEYBDINPUT_PASTE{
			Wvk:     VK_V_PASTE,
			Wscan:   0,
			Dwflags: KEYEVENTF_KEYUP_PASTE,
		},
	}

	// Ctrl 释放
	inputs[3] = INPUT_PASTE{
		Type: INPUT_KEYBOARD_PASTE,
		Ki: KEYBDINPUT_PASTE{
			Wvk:     VK_CONTROL_PASTE,
			Wscan:   0,
			Dwflags: KEYEVENTF_KEYUP_PASTE,
		},
	}

	// 发送输入事件
	result, _, err := procSendInput_paste.Call(4, uintptr(unsafe.Pointer(&inputs[0])), unsafe.Sizeof(inputs[0]))
	logPaste("SendInput 结果: %d, 错误: %v", result, err)

	success := result == 4
	if !success {
		logPaste("SendInput 失败: 只处理了 %d/4 个输入事件", result)
	}
	return success
}

// sendCtrlVWithKeybdEvent_paste 使用 keybd_event API 发送 Ctrl+V
func sendCtrlVWithKeybdEvent_paste() {
	procKeybdEvent := modUser32.NewProc("keybd_event")

	// Ctrl 按下
	procKeybdEvent.Call(VK_CONTROL_PASTE, 0, 0, 0)
	time.Sleep(10 * time.Millisecond)

	// V 按下
	procKeybdEvent.Call(VK_V_PASTE, 0, 0, 0)
	time.Sleep(10 * time.Millisecond)

	// V 释放
	procKeybdEvent.Call(VK_V_PASTE, 0, KEYEVENTF_KEYUP_PASTE, 0)
	time.Sleep(10 * time.Millisecond)

	// Ctrl 释放
	procKeybdEvent.Call(VK_CONTROL_PASTE, 0, KEYEVENTF_KEYUP_PASTE, 0)

	logPaste("已使用 keybd_event 发送 Ctrl+V")
}

// ActivatePreviousApp 激活之前的应用
func ActivatePreviousApp() {
	if previousWindow_paste != 0 {
		isValid, _, _ := procIsWindow_paste.Call(previousWindow_paste)
		if isValid != 0 {
			activateWindow_paste(previousWindow_paste)
		} else {
			previousWindow_paste = 0
		}
	}
}

// PasteCmdVToPreviousApp 快速激活应用并发送 Ctrl+V（Windows版本）
func PasteCmdVToPreviousApp() {
	logPaste("开始执行 PasteCmdVToPreviousApp, 目标窗口: %x", previousWindow_paste)

	// 方法1：使用记录的窗口句柄
	if previousWindow_paste != 0 {
		success := pasteToWindow_paste(previousWindow_paste)
		if success {
			logPaste("成功使用窗口句柄完成粘贴")
			return
		}
	}

	// 方法2：尝试查找当前前台窗口
	logPaste("尝试查找当前前台窗口")
	currentForeground, _, _ := procGetForegroundWindow.Call()
	if currentForeground != 0 && currentForeground != previousWindow_paste {
		var processId uint32
		procGetWindowThreadProcessId.Call(currentForeground, uintptr(unsafe.Pointer(&processId)))

		if processId != currentProcessId_paste {
			success := pasteToWindow_paste(currentForeground)
			if success {
				logPaste("成功向当前前台窗口完成粘贴")
				return
			}
		}
	}

	// 方法3：全局发送
	logPaste("使用全局发送作为最后的回退方案")
	time.Sleep(150 * time.Millisecond)
	sendCtrlV_paste()
}

// pasteToWindow_paste 向指定窗口发送粘贴命令
func pasteToWindow_paste(hwnd uintptr) bool {
	isValid, _, _ := procIsWindow_paste.Call(hwnd)
	if isValid == 0 {
		logPaste("窗口无效: %x", hwnd)
		return false
	}

	logPaste("尝试向窗口 %x 发送粘贴命令", hwnd)

	success := activateWindow_paste(hwnd)
	if !success {
		logPaste("激活窗口失败: %x", hwnd)
		return false
	}

	time.Sleep(100 * time.Millisecond)
	sendCtrlV_paste()
	logPaste("已向窗口 %x 发送 Ctrl+V", hwnd)

	return true
}
