//go:build !darwin

package common

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// EnsureWindowOnCurrentScreen 确保窗口在当前聚焦的屏幕上（非 macOS 平台使用 WindowCenter）
// 如果窗口已经在当前屏幕上，保持原位置不变
func EnsureWindowOnCurrentScreen(ctx context.Context) {
	// 非 macOS 平台：使用 WindowCenter（会居中到窗口当前所在的屏幕）
	runtime.WindowCenter(ctx)
}
