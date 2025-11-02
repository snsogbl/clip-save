//go:build !windows

package main

import (
	"context"
)

// initTray 其他平台的空实现
func initTray(app *App, ctx context.Context) {
	// macOS 和 Linux 暂不支持托盘，或使用平台特定的实现
}
