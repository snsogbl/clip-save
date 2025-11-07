package common

import "sync/atomic"

var forceQuitFlag atomic.Bool

// SetForceQuit 设置强制退出标志
func SetForceQuit() {
	forceQuitFlag.Store(true)
}

// IsForceQuit 检查是否强制退出
func IsForceQuit() bool {
	return forceQuitFlag.Load()
}

// ClearForceQuit 清除强制退出标志
func ClearForceQuit() {
	forceQuitFlag.Store(false)
}
