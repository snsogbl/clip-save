//go:build windows
// +build windows

package common

import "golang.design/x/hotkey"

// mapModifier 将字符串修饰键映射为 Windows 下的热键修饰常量
func mapModifier(part string) (hotkey.Modifier, bool) {
	switch part {
	case "ctrl", "control":
		return hotkey.ModCtrl, true
	case "shift":
		return hotkey.ModShift, true
	case "alt", "option":
		return hotkey.ModAlt, true
	case "cmd", "command", "meta":
		return hotkey.ModWin, true
	default:
		return 0, false
	}
}
