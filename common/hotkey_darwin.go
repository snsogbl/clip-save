package common

import (
	"context"
	"fmt"
	"strings"

	hook "github.com/robotn/gohook"
)

var (
	currentHotkey string
	hotkeyCancel  context.CancelFunc

	// 修饰键状态跟踪
	modifierKeysPressed = map[string]bool{
		"ctrl":  false,
		"shift": false,
		"alt":   false,
		"cmd":   false,
	}
)

// HotkeyCallback 快捷键回调函数类型
type HotkeyCallback func()

// ParseHotkey 解析快捷键字符串，例如 "Control+V"
func ParseHotkey(hotkeyStr string) ([]string, error) {
	parts := strings.Split(hotkeyStr, "+")
	if len(parts) < 1 {
		return nil, fmt.Errorf("无效的快捷键格式")
	}

	// 转换修饰键格式，按照 gohook 的要求：主键在前，修饰键在后
	keys := make([]string, 0, len(parts))

	// 先添加主键（最后一个）
	mainKey := strings.ToLower(strings.TrimSpace(parts[len(parts)-1]))
	keys = append(keys, mainKey)

	// 再添加修饰键
	for i := 0; i < len(parts)-1; i++ {
		mod := strings.ToLower(strings.TrimSpace(parts[i]))

		// 转换修饰键名称
		switch mod {
		case "control", "ctrl":
			keys = append(keys, "ctrl")
		case "shift":
			keys = append(keys, "shift")
		case "alt", "option":
			keys = append(keys, "alt")
		case "command", "cmd":
			keys = append(keys, "cmd")
		default:
			// 其他修饰键直接使用
			keys = append(keys, mod)
		}
	}

	return keys, nil
}

// isExactHotkeyMatch 检查事件是否精确匹配注册的快捷键组合
func isExactHotkeyMatch(e hook.Event, registeredKeys []string) bool {
	// 获取当前按下的修饰键状态
	currentModifiers := getCurrentModifierKeys(e)

	// 分析注册的快捷键组合
	mainKey, modifiers := parseRegisteredKeys(registeredKeys)

	// 检查主键是否匹配
	if !isMainKeyMatch(e, mainKey) {
		return false
	}

	// 检查修饰键是否精确匹配
	return isModifiersExactMatch(currentModifiers, modifiers)
}

// setupModifierKeyTracking 设置修饰键状态跟踪
func setupModifierKeyTracking() {
	// 注册修饰键的按下和释放事件
	modifierKeys := []string{"ctrl", "shift", "alt", "cmd"}

	for _, key := range modifierKeys {
		// 创建局部变量来避免闭包问题
		keyCopy := key

		// 注册按键按下事件
		hook.Register(hook.KeyDown, []string{key}, func(e hook.Event) {
			modifierKeysPressed[keyCopy] = true
			// log.Printf("🔑 修饰键按下: %s (键码: %d)", keyCopy, e.Keycode)
		})

		// 注册按键释放事件
		hook.Register(hook.KeyUp, []string{key}, func(e hook.Event) {
			modifierKeysPressed[keyCopy] = false
			// log.Printf("🔓 修饰键释放: %s (键码: %d)", keyCopy, e.Keycode)
		})
	}

	// log.Printf("✅ 修饰键状态跟踪已设置: %v", modifierKeys)
}

// getCurrentModifierKeys 从全局状态获取当前按下的修饰键
func getCurrentModifierKeys(e hook.Event) map[string]bool {
	// 返回当前修饰键状态的副本
	modifiers := make(map[string]bool)
	for key, pressed := range modifierKeysPressed {
		modifiers[key] = pressed
	}

	// 添加调试信息
	// log.Printf("🔍 当前修饰键状态: %v", modifiers)

	// 过滤掉未按下的修饰键，只返回实际按下的修饰键
	pressedModifiers := make(map[string]bool)
	for key, pressed := range modifiers {
		if pressed {
			pressedModifiers[key] = true
		}
	}

	return pressedModifiers
}

// parseRegisteredKeys 解析注册的快捷键，分离主键和修饰键
func parseRegisteredKeys(keys []string) (string, []string) {
	if len(keys) == 0 {
		return "", nil
	}

	// 主键是第一个（根据ParseHotkey的实现）
	mainKey := keys[0]

	// 修饰键是其余的
	modifiers := keys[1:]

	return mainKey, modifiers
}

// isMainKeyMatch 检查主键是否匹配
func isMainKeyMatch(e hook.Event, expectedMainKey string) bool {
	// 使用gohook库的Keycode映射来检查主键是否匹配
	expectedKeycode, exists := hook.Keycode[expectedMainKey]
	if !exists {
		// log.Printf("⚠️ 未找到主键的键码映射: %s", expectedMainKey)
		return false
	}

	// 检查事件的Keycode是否匹配期望的主键
	if e.Keycode == expectedKeycode {
		// log.Printf("✅ 主键匹配: %s (键码: %d)", expectedMainKey, e.Keycode)
		return true
	}

	// 也检查Keychar是否匹配（作为备用检查）
	if e.Keychar == rune(expectedMainKey[0]) {
		// log.Printf("✅ 主键通过Keychar匹配: %s (字符: %c)", expectedMainKey, e.Keychar)
		return true
	}

	// log.Printf("❌ 主键不匹配: 期望=%s(键码:%d), 实际键码=%d, 实际字符=%c",
	// 	expectedMainKey, expectedKeycode, e.Keycode, e.Keychar)
	return false
}

// isModifiersExactMatch 检查修饰键是否精确匹配
func isModifiersExactMatch(currentModifiers map[string]bool, expectedModifiers []string) bool {
	// 创建期望的修饰键映射
	expected := make(map[string]bool)
	for _, mod := range expectedModifiers {
		expected[mod] = true
	}

	// log.Printf("🔍 修饰键匹配检查:")
	// log.Printf("   期望的修饰键: %v", expectedModifiers)
	// log.Printf("   当前按下的修饰键: %v", getPressedModifiers(currentModifiers))

	// 检查是否所有期望的修饰键都被按下
	for mod := range expected {
		if !currentModifiers[mod] {
			// log.Printf("❌ 缺少期望的修饰键: %s", mod)
			return false
		}
	}

	// 检查是否有额外的修饰键被按下（除了期望的）
	// 注意：currentModifiers 现在只包含实际按下的修饰键
	// 所以这里只需要检查是否有不在期望列表中的修饰键被按下
	for mod := range currentModifiers {
		if !expected[mod] {
			// 有额外的修饰键被按下，不匹配
			// log.Printf("❌ 检测到额外的修饰键: %s", mod)
			return false
		}
	}

	// log.Printf("✅ 修饰键精确匹配")
	return true
}

// resetModifierKeys 重置所有修饰键状态为未按下
func resetModifierKeys() {
	for key := range modifierKeysPressed {
		modifierKeysPressed[key] = false
	}
	// log.Printf("🔄 修饰键状态已重置")
}

// RegisterHotkey 注册全局快捷键
func RegisterHotkey(hotkeyStr string, callback HotkeyCallback) error {
	// 如果已有快捷键在运行，先停止
	UnregisterHotkey()

	// 重置修饰键状态
	resetModifierKeys()

	keys, err := ParseHotkey(hotkeyStr)
	if err != nil {
		return err
	}

	// log.Printf("🔥 注册快捷键: %s -> keys: %v", hotkeyStr, keys)

	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	hotkeyCancel = cancel
	currentHotkey = hotkeyStr

	// 在新的 goroutine 中启动快捷键
	go func() {
		// 注册修饰键状态跟踪
		setupModifierKeyTracking()

		// 注册快捷键组合（按照参考案例的方式）
		hook.Register(hook.KeyDown, keys, func(e hook.Event) {
			// log.Printf("🎯 快捷键事件触发: %s, 事件详情: Keycode=%d, Keychar=%c, Mask=0x%04X",
			// 	hotkeyStr, e.Keycode, e.Keychar, e.Mask)

			// 检查是否精确匹配注册的快捷键组合
			if isExactHotkeyMatch(e, keys) {
				// log.Printf("🔥 快捷键精确匹配触发: %s", hotkeyStr)
				if callback != nil {
					callback()
				}
			} else {
				// log.Printf("🔍 快捷键部分匹配但不符合精确条件: %s", hotkeyStr)
			}
		})

		// 启动钩子
		s := hook.Start()
		// log.Printf("✅ 快捷键注册成功: %s", hotkeyStr)

		// 取消信号
		go func() {
			<-ctx.Done()
			// log.Println("⏹️ 停止快捷键")
			hook.End()
		}()

		// 等待事件（这会阻塞直到程序结束）
		<-hook.Process(s)
	}()

	return nil
}

// UnregisterHotkey 取消注册快捷键
func UnregisterHotkey() {
	if hotkeyCancel != nil {
		// log.Printf("🔥 取消快捷键: %s", currentHotkey)
		hotkeyCancel()
		hotkeyCancel = nil
		currentHotkey = ""
	}
}

// GetCurrentHotkey 获取当前注册的快捷键
func GetCurrentHotkey() string {
	return currentHotkey
}
