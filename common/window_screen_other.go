//go:build !darwin

package common

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// MoveWindowToCurrentScreen 将窗口移动到当前聚焦的屏幕中心（非 macOS 平台使用 WindowCenter）
func MoveWindowToCurrentScreen(ctx context.Context) {
	runtime.WindowCenter(ctx)
}
