//go:build !darwin

package common

// InitDockReopen 非 macOS 平台为空实现
func InitDockReopen(callback func()) {}

// SetForceQuitCallback 非 macOS 平台为空实现
func SetForceQuitCallback(callback func()) {}
