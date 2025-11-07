//go:build !darwin
// +build !darwin

package common

// SetDockIconVisibility 其他平台的空实现
func SetDockIconVisibility(visible int) {
	// 其他平台不需要实现
}

