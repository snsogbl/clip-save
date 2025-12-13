//go:build windows
// +build windows

package common

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
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
	windowMutex            sync.Mutex // 保护 previousWindow_paste 的并发访问
)

// InitAppSwitchListener 初始化应用切换监听器（Windows版本）
func InitAppSwitchListener() {
	// 强制初始化日志文件（即使失败也要继续）
	initPasteLogger()

	// 立即写入启动日志，确认函数被调用
	logPaste("=== Windows 粘贴功能启动 ===")
	logPaste("InitAppSwitchListener 函数被调用")

	// 获取当前进程ID
	pid, _, _ := procGetCurrentProcessId_paste.Call()
	currentProcessId_paste = uint32(pid)

	logPaste("初始化应用切换监听器，当前进程ID: %d", currentProcessId_paste)

	// 检查管理员权限
	CheckAdminRights()

	// 记录当前前台窗口（如果不是我们的应用）
	recordCurrentForegroundWindow()

	// 启动后台窗口监听器（使用轮询方式，更可靠）
	startWindowMonitor()

	logPaste("=== Windows 粘贴功能初始化完成 ===")
	logPaste("提示：可以使用 TestWindowsWindowMonitoring() 测试窗口监听功能")
	logPaste("提示：可以使用 TestWindowsPasteFunction() 测试完整粘贴功能")
}

// 全局日志文件句柄
var pasteLogFile *os.File

// initPasteLogger 初始化粘贴功能的日志记录器
func initPasteLogger() {
	// 获取用户的临时目录
	tempDir := os.TempDir()
	logDir := filepath.Join(tempDir, "ClipSave")

	log.Printf("[PASTE] 尝试创建日志目录: %s", logDir)

	// 创建日志目录
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		log.Printf("[PASTE] 创建日志目录失败: %v", err)
		log.Printf("[PASTE] 将只使用控制台输出")
		return
	}

	log.Printf("[PASTE] 日志目录创建成功: %s", logDir)

	// 创建日志文件（按日期命名）
	logFileName := fmt.Sprintf("paste_%s.log", time.Now().Format("2006-01-02"))
	logFilePath := filepath.Join(logDir, logFileName)

	log.Printf("[PASTE] 尝试创建日志文件: %s", logFilePath)

	// 打开或创建日志文件
	pasteLogFile, err = os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("[PASTE] 创建日志文件失败: %v", err)
		log.Printf("[PASTE] 将只使用控制台输出")
		return
	}

	// 写入初始化日志
	logPaste("=== 粘贴功能日志初始化完成 ===")
	logPaste("日志文件位置: %s", logFilePath)
	logPaste("提示: 可以在文件资源管理器中输入 %%TEMP%%\\ClipSave 查看日志文件")

	log.Printf("[PASTE] ✅ 日志文件初始化成功")
}

// logPaste 记录粘贴相关的日志
func logPaste(format string, args ...any) {
	message := fmt.Sprintf("[PASTE] "+format, args...)

	// 总是输出到控制台
	log.Print(message)

	// 如果日志文件可用，直接写入文件
	if pasteLogFile != nil {
		timestamp := time.Now().Format("2006/01/02 15:04:05")
		fileMessage := fmt.Sprintf("%s %s\n", timestamp, message)

		pasteLogFile.WriteString(fileMessage)
		pasteLogFile.Sync() // 立即刷新
	}
}

// RecordPreviousAppPID 记录当前前台应用（在激活本应用之前调用）
func RecordPreviousAppPID() {
	logPaste("=== RecordPreviousAppPID 被调用 ===")
	recordCurrentForegroundWindow()
	logPaste("=== RecordPreviousAppPID 完成 ===")
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
			// 获取窗口标题用于调试
			windowTitle := getWindowTitle(hwnd)
			previousWindow_paste = hwnd
			logPaste("✅ 记录前台窗口: %x, 进程ID: %d, 标题: %s", hwnd, processId, windowTitle)
		} else {
			logPaste("⚠️ 当前前台窗口是我们自己的应用，不记录")
		}
	} else {
		logPaste("❌ 无法获取前台窗口")
	}
}

// getWindowTitle 获取窗口标题（用于调试）
func getWindowTitle(hwnd uintptr) string {
	procGetWindowText := modUser32.NewProc("GetWindowTextW")

	// 创建缓冲区
	buf := make([]uint16, 256)

	// 获取窗口标题
	length, _, _ := procGetWindowText.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))

	if length > 0 {
		return syscall.UTF16ToString(buf[:length])
	}

	return "Unknown"
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
		logPaste("❌ 窗口无效: %x", hwnd)
		return false
	}

	// 方法1：标准激活流程
	logPaste("方法1: 标准窗口激活")

	// 先显示窗口
	procShowWindow_paste.Call(hwnd, SW_RESTORE_PASTE)
	time.Sleep(50 * time.Millisecond)

	procShowWindow_paste.Call(hwnd, SW_SHOW_PASTE)
	time.Sleep(50 * time.Millisecond)

	// 将窗口置于前台
	procBringWindowToTop_paste.Call(hwnd)
	time.Sleep(50 * time.Millisecond)

	// 设置为前台窗口
	result, _, _ := procSetForegroundWindow_paste.Call(hwnd)

	if result != 0 {
		logPaste("✅ 标准方法激活成功: %x", hwnd)
		return true
	}

	// 方法2：强制激活 (使用线程附加)
	logPaste("⚠️ 标准方法失败，尝试方法2: 强制激活")
	return forceActivateWindow_paste(hwnd)
}

// forceActivateWindow_paste 强制激活窗口（使用线程附加技术）
func forceActivateWindow_paste(hwnd uintptr) bool {
	logPaste("尝试强制激活窗口: %x", hwnd)

	// 获取目标窗口的线程ID
	var targetProcessId uint32
	targetThreadId, _, _ := procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&targetProcessId)))

	// 获取当前线程ID
	currentThreadId, _, _ := procGetCurrentThreadId_paste.Call()

	logPaste("目标线程ID: %d, 当前线程ID: %d", targetThreadId, currentThreadId)

	if targetThreadId != currentThreadId {
		// 附加到目标线程的输入队列
		procAttachThreadInput_paste.Call(currentThreadId, targetThreadId, 1)
		logPaste("已附加到目标线程输入队列")
	}

	// 尝试激活
	procSetActiveWindow_paste.Call(hwnd)
	result, _, _ := procSetForegroundWindow_paste.Call(hwnd)

	// 分离线程输入队列
	if targetThreadId != currentThreadId {
		procAttachThreadInput_paste.Call(currentThreadId, targetThreadId, 0)
		logPaste("已分离线程输入队列")
	}

	success := result != 0
	logPaste("强制激活结果: %x, 成功: %v", hwnd, success)
	return success
}

// PasteCmdV 发送 Ctrl+V 粘贴命令（Windows版本）
func PasteCmdV() {
	log.Printf("[PASTE-FORCE] PasteCmdV 被调用！")
	logPaste("=== 开始执行 PasteCmdV ===")

	// Windows 需要更长的延迟确保窗口激活
	logPaste("等待 500ms 确保窗口激活完成")
	time.Sleep(500 * time.Millisecond)

	// 检查当前前台窗口
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd != 0 {
		windowTitle := getWindowTitle(hwnd)
		logPaste("当前前台窗口: %s", windowTitle)
	}

	logPaste("发送 Ctrl+V 到当前前台窗口")
	sendCtrlV_paste()
	logPaste("=== PasteCmdV 执行完成 ===")
}

// sendCtrlV_paste 发送 Ctrl+V 组合键
func sendCtrlV_paste() {
	logPaste("开始发送 Ctrl+V 组合键")

	// 方法1：使用 keybd_event (最兼容)
	logPaste("尝试方法1: keybd_event (最兼容)")
	sendCtrlVWithKeybdEvent_paste()

	// 短暂延迟后尝试其他方法作为备份
	time.Sleep(100 * time.Millisecond)

	// 方法2：使用 SendInput
	logPaste("尝试方法2: SendInput")
	if sendCtrlVWithSendInput_paste() {
		logPaste("✅ SendInput 方法成功")
	}

	// 方法3：使用 PostMessage
	logPaste("尝试方法3: PostMessage")
	sendCtrlVWithPostMessage_paste()
}

// sendCtrlVWithSendInput_paste 使用 SendInput API 发送 Ctrl+V
func sendCtrlVWithSendInput_paste() bool {
	logPaste("尝试使用 SendInput 发送 Ctrl+V")

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
	result, _, err := procSendInput_paste.Call(
		uintptr(4),
		uintptr(unsafe.Pointer(&inputs[0])),
		uintptr(unsafe.Sizeof(inputs[0])),
	)

	logPaste("SendInput 调用完成 - 结果: %d, 错误: %v", result, err)
	logPaste("SendInput 参数 - 事件数: 4, 结构体大小: %d", unsafe.Sizeof(inputs[0]))

	success := result == 4
	if !success {
		logPaste("❌ SendInput 失败: 只处理了 %d/4 个输入事件", result)

		// 获取更详细的错误信息
		if err != nil {
			logPaste("❌ SendInput 系统错误: %v", err)
		}
	} else {
		logPaste("✅ SendInput 成功发送了所有 4 个输入事件")
	}
	return success
}

// sendCtrlVWithKeybdEvent_paste 使用 keybd_event API 发送 Ctrl+V
func sendCtrlVWithKeybdEvent_paste() {
	procKeybdEvent := modUser32.NewProc("keybd_event")

	logPaste("使用 keybd_event 发送 Ctrl+V")

	// Ctrl 按下
	result1, _, _ := procKeybdEvent.Call(VK_CONTROL_PASTE, 0, 0, 0)
	logPaste("Ctrl 按下结果: %d", result1)
	time.Sleep(50 * time.Millisecond)

	// V 按下
	result2, _, _ := procKeybdEvent.Call(VK_V_PASTE, 0, 0, 0)
	logPaste("V 按下结果: %d", result2)
	time.Sleep(50 * time.Millisecond)

	// V 释放
	result3, _, _ := procKeybdEvent.Call(VK_V_PASTE, 0, KEYEVENTF_KEYUP_PASTE, 0)
	logPaste("V 释放结果: %d", result3)
	time.Sleep(50 * time.Millisecond)

	// Ctrl 释放
	result4, _, _ := procKeybdEvent.Call(VK_CONTROL_PASTE, 0, KEYEVENTF_KEYUP_PASTE, 0)
	logPaste("Ctrl 释放结果: %d", result4)

	logPaste("✅ keybd_event 序列完成")
}

// sendCtrlVWithPostMessage_paste 使用 PostMessage 发送 Ctrl+V (新增方法)
func sendCtrlVWithPostMessage_paste() {
	logPaste("尝试使用 PostMessage 发送 Ctrl+V")

	// 获取当前前台窗口
	hwnd, _, _ := procGetForegroundWindow.Call()
	if hwnd == 0 {
		logPaste("❌ PostMessage: 无法获取前台窗口")
		return
	}

	procPostMessage := modUser32.NewProc("PostMessageW")

	const (
		WM_KEYDOWN = 0x0100
		WM_KEYUP   = 0x0101
		WM_CHAR    = 0x0102
	)

	// 发送 Ctrl+V 消息
	// 先发送 Ctrl 按下
	procPostMessage.Call(hwnd, WM_KEYDOWN, VK_CONTROL_PASTE, 0)
	time.Sleep(10 * time.Millisecond)

	// 发送 V 按下
	procPostMessage.Call(hwnd, WM_KEYDOWN, VK_V_PASTE, 0)
	time.Sleep(10 * time.Millisecond)

	// 发送 V 释放
	procPostMessage.Call(hwnd, WM_KEYUP, VK_V_PASTE, 0)
	time.Sleep(10 * time.Millisecond)

	// 发送 Ctrl 释放
	procPostMessage.Call(hwnd, WM_KEYUP, VK_CONTROL_PASTE, 0)

	logPaste("✅ PostMessage 方法执行完成")
}

// ActivatePreviousApp 激活之前的应用
func ActivatePreviousApp() {
	log.Printf("[PASTE-FORCE] ActivatePreviousApp 被调用！")
	logPaste("=== 开始激活之前的应用 ===")

	// 线程安全地读取窗口句柄
	windowMutex.Lock()
	targetWindow := previousWindow_paste
	windowMutex.Unlock()

	if targetWindow != 0 {
		logPaste("检查记录的窗口: %x", targetWindow)
		isValid, _, _ := procIsWindow_paste.Call(targetWindow)
		if isValid != 0 {
			windowTitle := getWindowTitle(targetWindow)
			logPaste("窗口有效，标题: %s，尝试激活", windowTitle)

			// 多次尝试激活
			success := false
			for i := 0; i < 3; i++ {
				success = activateWindow_paste(targetWindow)
				if success {
					logPaste("✅ 第 %d 次尝试激活成功", i+1)
					break
				}
				logPaste("⚠️ 第 %d 次激活失败，重试", i+1)
				time.Sleep(200 * time.Millisecond)
			}

			if !success {
				logPaste("❌ 多次尝试后仍然激活失败")
			}
		} else {
			logPaste("❌ 记录的窗口已无效，清除记录")
			windowMutex.Lock()
			previousWindow_paste = 0
			windowMutex.Unlock()
		}
	} else {
		logPaste("⚠️ 没有记录的前台窗口")
	}

	logPaste("=== ActivatePreviousApp 完成 ===")
}

// PasteCmdVToPreviousApp 快速激活应用并发送 Ctrl+V（Windows版本）
func PasteCmdVToPreviousApp() {
	logPaste("=== 开始执行 PasteCmdVToPreviousApp ===")

	// 线程安全地读取目标窗口
	windowMutex.Lock()
	targetWindow := previousWindow_paste
	windowMutex.Unlock()

	logPaste("目标窗口: %x", targetWindow)
	logPaste("当前进程ID: %d", currentProcessId_paste)

	// 方法1：使用记录的窗口句柄
	if targetWindow != 0 {
		logPaste("方法1: 使用记录的窗口句柄")

		// 验证窗口是否仍然有效
		isValid, _, _ := procIsWindow_paste.Call(targetWindow)
		if isValid != 0 {
			windowTitle := getWindowTitle(targetWindow)
			logPaste("目标窗口有效: %s (句柄: %x)", windowTitle, targetWindow)

			success := pasteToWindow_paste(targetWindow)
			if success {
				logPaste("✅ 成功使用窗口句柄完成粘贴")
				return
			}
			logPaste("❌ 窗口句柄方法失败，尝试其他方法")
		} else {
			logPaste("❌ 记录的窗口句柄已无效，清除记录")
			windowMutex.Lock()
			previousWindow_paste = 0
			windowMutex.Unlock()
		}
	} else {
		logPaste("⚠️ 没有记录的窗口句柄")
		logPaste("提示：请确保在使用粘贴功能前切换过其他应用")
	}

	// 方法2：查找最近使用的任何其他窗口
	logPaste("方法2: 查找最近使用的其他窗口")
	otherWindow := findAnyOtherWindow()
	if otherWindow != 0 {
		success := pasteToWindow_paste(otherWindow)
		if success {
			logPaste("✅ 成功向其他窗口完成粘贴")
			return
		}
		logPaste("❌ 其他窗口方法失败")
	} else {
		logPaste("❌ 没有找到其他可用窗口")
	}

	// 方法3：尝试查找当前前台窗口
	logPaste("方法3: 尝试查找当前前台窗口")
	currentForeground, _, _ := procGetForegroundWindow.Call()
	if currentForeground != 0 && currentForeground != targetWindow {
		var processId uint32
		procGetWindowThreadProcessId.Call(currentForeground, uintptr(unsafe.Pointer(&processId)))
		windowTitle := getWindowTitle(currentForeground)
		logPaste("当前前台窗口: %s (句柄: %x, 进程ID: %d)", windowTitle, currentForeground, processId)

		if processId != currentProcessId_paste {
			success := pasteToWindow_paste(currentForeground)
			if success {
				logPaste("✅ 成功向当前前台窗口完成粘贴")
				return
			}
			logPaste("❌ 当前前台窗口方法失败")
		} else {
			logPaste("⚠️ 当前前台窗口是我们自己的应用")
		}
	} else {
		logPaste("⚠️ 没有找到有效的前台窗口")
	}

	// 方法4：全局发送（类似 Mac 的回退机制）
	logPaste("方法4: 使用全局发送作为最后的回退方案")
	logPaste("⚠️ 所有窗口定位方法都失败，使用全局键盘输入")
	time.Sleep(150 * time.Millisecond)
	sendCtrlV_paste()
	logPaste("=== PasteCmdVToPreviousApp 执行完成 ===")
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

// TestPasteFunction 测试粘贴功能（供调试使用）
func TestPasteFunction() {
	logPaste("=== 开始测试粘贴功能 ===")

	// 测试日志功能
	logPaste("✅ 日志功能正常")

	// 测试进程ID获取
	logPaste("当前进程ID: %d", currentProcessId_paste)

	// 测试前台窗口获取
	currentForeground, _, _ := procGetForegroundWindow.Call()
	if currentForeground != 0 {
		var processId uint32
		procGetWindowThreadProcessId.Call(currentForeground, uintptr(unsafe.Pointer(&processId)))
		logPaste("当前前台窗口: %x, 进程ID: %d", currentForeground, processId)
	} else {
		logPaste("❌ 无法获取前台窗口")
	}

	// 测试记录的窗口
	if previousWindow_paste != 0 {
		logPaste("记录的前台窗口: %x", previousWindow_paste)
		isValid, _, _ := procIsWindow_paste.Call(previousWindow_paste)
		if isValid != 0 {
			logPaste("✅ 记录的窗口仍然有效")
		} else {
			logPaste("❌ 记录的窗口已无效")
		}
	} else {
		logPaste("⚠️ 没有记录的前台窗口")
	}

	// 测试键盘输入
	logPaste("测试发送 Ctrl+V...")
	sendCtrlV_paste()

	logPaste("=== 粘贴功能测试完成 ===")
}

// SimpleTestPaste 简单测试粘贴功能（直接发送到当前窗口）
func SimpleTestPaste() {
	logPaste("=== 开始简单粘贴测试 ===")
	logPaste("注意：请确保有文本应用（如记事本）处于前台")

	// 等待3秒让用户切换到目标应用
	for i := 3; i > 0; i-- {
		logPaste("倒计时: %d 秒", i)
		time.Sleep(1 * time.Second)
	}

	logPaste("开始发送 Ctrl+V...")
	sendCtrlV_paste()

	logPaste("=== 简单粘贴测试完成 ===")
}

// TestLogWriting 测试日志写入功能
func TestLogWriting() {
	log.Printf("[PASTE-FORCE] TestLogWriting 开始")

	logPaste("=== 测试日志写入功能 ===")
	logPaste("当前时间: %s", time.Now().Format("2006-01-02 15:04:05"))
	logPaste("测试消息1: 这是一条测试日志")
	logPaste("测试消息2: 日志文件应该能看到这些内容")
	logPaste("测试消息3: 如果看到这些，说明日志写入正常")

	// 检查日志文件状态
	if pasteLogFile != nil {
		// 获取文件信息
		fileInfo, err := pasteLogFile.Stat()
		if err == nil {
			logPaste("✅ 日志文件状态正常，大小: %d 字节", fileInfo.Size())
		} else {
			logPaste("❌ 获取日志文件状态失败: %v", err)
		}

		// 强制刷新
		pasteLogFile.Sync()
		logPaste("✅ 已强制刷新日志文件")

		// 再次检查文件大小
		fileInfo2, err := pasteLogFile.Stat()
		if err == nil {
			logPaste("✅ 刷新后文件大小: %d 字节", fileInfo2.Size())
		}
	} else {
		logPaste("❌ 日志文件句柄为空")
	}

	logPaste("=== 日志写入测试完成 ===")

	log.Printf("[PASTE-FORCE] TestLogWriting 完成")
}

// TestKeyboardOnly 只测试键盘输入，不涉及窗口操作
func TestKeyboardOnly() {
	logPaste("=== 纯键盘输入测试 ===")
	logPaste("请确保有文本编辑器在前台（如记事本）")
	logPaste("3秒后将发送 Ctrl+V")

	// 倒计时
	for i := 3; i > 0; i-- {
		logPaste("倒计时: %d", i)
		time.Sleep(1 * time.Second)
	}

	logPaste("开始发送键盘事件...")

	// 直接使用 keybd_event，最简单的方式
	procKeybdEvent := modUser32.NewProc("keybd_event")

	// Ctrl 按下
	procKeybdEvent.Call(VK_CONTROL_PASTE, 0, 0, 0)
	time.Sleep(100 * time.Millisecond)

	// V 按下
	procKeybdEvent.Call(VK_V_PASTE, 0, 0, 0)
	time.Sleep(100 * time.Millisecond)

	// V 释放
	procKeybdEvent.Call(VK_V_PASTE, 0, KEYEVENTF_KEYUP_PASTE, 0)
	time.Sleep(100 * time.Millisecond)

	// Ctrl 释放
	procKeybdEvent.Call(VK_CONTROL_PASTE, 0, KEYEVENTF_KEYUP_PASTE, 0)

	logPaste("✅ 键盘事件发送完成")
	logPaste("=== 测试结束 ===")
}

// TestWindowMonitoring 测试窗口监听功能
func TestWindowMonitoring() {
	logPaste("=== 测试窗口监听功能 ===")
	logPaste("请在接下来的10秒内切换到不同的应用程序")
	logPaste("观察日志中是否有窗口切换记录")

	// 记录当前前台窗口
	recordCurrentForegroundWindow()

	// 等待10秒，观察窗口切换
	for i := 10; i > 0; i-- {
		logPaste("剩余时间: %d 秒", i)
		time.Sleep(1 * time.Second)
	}

	// 显示最终记录的窗口
	windowMutex.Lock()
	finalWindow := previousWindow_paste
	windowMutex.Unlock()

	if finalWindow != 0 {
		windowTitle := getWindowTitle(finalWindow)
		logPaste("✅ 最终记录的窗口: %s (句柄: %x)", windowTitle, finalWindow)
	} else {
		logPaste("❌ 没有记录到任何窗口")
	}

	logPaste("=== 窗口监听测试完成 ===")
}

// CheckAdminRights 检查是否以管理员身份运行
func CheckAdminRights() bool {
	procIsUserAnAdmin := modShell32.NewProc("IsUserAnAdmin")
	if procIsUserAnAdmin.Find() != nil {
		logPaste("⚠️ 无法检查管理员权限")
		return false
	}

	result, _, _ := procIsUserAnAdmin.Call()
	isAdmin := result != 0

	if isAdmin {
		logPaste("✅ 当前以管理员身份运行")
	} else {
		logPaste("⚠️ 当前未以管理员身份运行，可能影响粘贴功能")
		logPaste("建议：右键点击应用 → 以管理员身份运行")
	}

	return isAdmin
}

// findAnyOtherWindow 查找任何其他可见的窗口（类似 Mac 的简单回退）
func findAnyOtherWindow() uintptr {
	logPaste("查找任何其他可见窗口")

	var foundWindow uintptr

	// 枚举所有窗口，找到第一个不是我们进程的可见窗口
	enumProc := syscall.NewCallback(func(hwnd uintptr, lParam uintptr) uintptr {
		// 检查窗口是否可见
		procIsWindowVisible := modUser32.NewProc("IsWindowVisible")
		visible, _, _ := procIsWindowVisible.Call(hwnd)
		if visible == 0 {
			return 1 // 继续枚举
		}

		// 获取窗口进程ID
		var processId uint32
		procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&processId)))

		// 跳过自己的进程
		if processId == currentProcessId_paste {
			return 1
		}

		// 检查窗口是否有标题（过滤掉系统窗口）
		windowTitle := getWindowTitle(hwnd)
		if len(windowTitle) > 0 {
			logPaste("找到其他窗口: %s (句柄: %x)", windowTitle, hwnd)
			foundWindow = hwnd
			return 0 // 停止枚举，使用第一个找到的
		}

		return 1 // 继续枚举
	})

	procEnumWindows_paste.Call(enumProc, 0)

	if foundWindow == 0 {
		logPaste("未找到其他可见窗口")
	}

	return foundWindow
}

// startWindowMonitor 启动窗口监听器（使用轮询方式，更可靠）
func startWindowMonitor() {
	logPaste("启动窗口切换监听器（轮询方式）")

	// 使用轻量级轮询方式监听窗口变化（更可靠）
	go func() {
		var lastWindow uintptr
		ticker := time.NewTicker(500 * time.Millisecond) // 每500ms检查一次
		defer ticker.Stop()

		for range ticker.C {
			current, _, _ := procGetForegroundWindow.Call()
			if current != 0 && current != lastWindow {
				var processId uint32
				procGetWindowThreadProcessId.Call(current, uintptr(unsafe.Pointer(&processId)))

				// 如果不是我们的进程，记录这个窗口
				if processId != currentProcessId_paste {
					windowTitle := getWindowTitle(current)
					logPaste("轮询检测到窗口切换: %s (句柄: %x, 进程ID: %d)", windowTitle, current, processId)

					// 线程安全地更新记录的前台窗口
					windowMutex.Lock()
					previousWindow_paste = current
					windowMutex.Unlock()
				}
				lastWindow = current
			}
		}
	}()

	logPaste("✅ 窗口轮询监听器启动成功")
}
