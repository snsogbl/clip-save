//go:build !darwin
// +build !darwin

package common

// AdjustWindowButtons 调整窗口控制按钮位置（非 macOS 平台为空操作）
func AdjustWindowButtons() {
	// 非 macOS 平台不需要调整
}

func CleanupWindowButtonsObserver() {
	// 非 macOS 平台不需要清理
}
