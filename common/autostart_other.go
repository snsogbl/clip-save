//go:build !windows

package common

// 非 Windows 平台：no-op 实现，保证跨平台编译通过
// 后续可扩展：macOS -> LaunchAgent plist；Linux -> ~/.config/autostart/*.desktop

// SetAutoStart 非 Windows 平台暂不支持开机自启
func SetAutoStart(enable bool) error {
	return nil
}

// IsAutoStartEnabled 非 Windows 平台始终返回 false
func IsAutoStartEnabled() (bool, error) {
	return false, nil
}
