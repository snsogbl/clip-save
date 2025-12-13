//go:build windows
// +build windows

package common

import (
	"fmt"
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

	logPaste("=== Windows 粘贴功能初始化完成 ===")
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
	// 强制输出到控制台
	log.Printf("[PASTE-FORCE] RecordPreviousAppPID 被调用！")
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

	if previousWindow_paste != 0 {
		logPaste("检查记录的窗口: %x", previousWindow_paste)
		isValid, _, _ := procIsWindow_paste.Call(previousWindow_paste)
		if isValid != 0 {
			windowTitle := getWindowTitle(previousWindow_paste)
			logPaste("窗口有效，标题: %s，尝试激活", windowTitle)

			// 多次尝试激活
			success := false
			for i := 0; i < 3; i++ {
				success = activateWindow_paste(previousWindow_paste)
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
			previousWindow_paste = 0
		}
	} else {
		logPaste("⚠️ 没有记录的前台窗口")
	}

	logPaste("=== ActivatePreviousApp 完成 ===")
}

// PasteCmdVToPreviousApp 快速激活应用并发送 Ctrl+V（Windows版本）
func PasteCmdVToPreviousApp() {
	// 强制输出到控制台，确保能看到函数被调用
	log.Printf("[PASTE-FORCE] PasteCmdVToPreviousApp 被调用！")

	logPaste("=== 开始执行 PasteCmdVToPreviousApp ===")
	logPaste("目标窗口: %x", previousWindow_paste)
	logPaste("当前进程ID: %d", currentProcessId_paste)

	// 方法1：使用记录的窗口句柄
	if previousWindow_paste != 0 {
		logPaste("方法1: 使用记录的窗口句柄")
		success := pasteToWindow_paste(previousWindow_paste)
		if success {
			logPaste("✅ 成功使用窗口句柄完成粘贴")
			return
		}
		logPaste("❌ 窗口句柄方法失败，尝试其他方法")
	} else {
		logPaste("⚠️ 没有记录的窗口句柄")
	}

	// 方法2：尝试查找当前前台窗口
	logPaste("方法2: 尝试查找当前前台窗口")
	currentForeground, _, _ := procGetForegroundWindow.Call()
	if currentForeground != 0 && currentForeground != previousWindow_paste {
		var processId uint32
		procGetWindowThreadProcessId.Call(currentForeground, uintptr(unsafe.Pointer(&processId)))
		logPaste("当前前台窗口: %x, 进程ID: %d", currentForeground, processId)

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

	// 方法3：全局发送
	logPaste("方法3: 使用全局发送作为最后的回退方案")
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
	log.Printf("[PASTE-FORCE] 开始测试纯键盘输入")

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

	log.Printf("[PASTE-FORCE] 纯键盘输入测试完成")
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
