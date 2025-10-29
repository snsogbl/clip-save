//go:build !darwin && !windows
// +build !darwin,!windows

package common

// HotkeyCallback 快捷键回调函数类型
type HotkeyCallback func()

// RegisterHotkey 非 macOS/Windows 平台暂不支持，返回错误
func RegisterHotkey(hotkeyStr string, callback HotkeyCallback) error {
	return nil
}

// UnregisterHotkey 非 macOS/Windows 平台无操作
func UnregisterHotkey() {}
