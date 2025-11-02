//go:build windows

package main

import (
	"context"
	_ "embed"
	"os"

	"github.com/getlantern/systray"

	"goWeb3/common"
)

//go:embed appicon_64.png
var trayIconData []byte

var (
	trayAppInstance *App
	trayCtx         context.Context
	showMenuItem    *systray.MenuItem
	quitMenuItem    *systray.MenuItem
)

// initTray 初始化系统托盘（Windows）
func initTray(app *App, ctx context.Context) {
	trayAppInstance = app
	trayCtx = ctx

	// 加载托盘图标
	iconData := loadTrayIcon()

	// 定义 onReady 回调
	onReady := func() {
		systray.SetIcon(iconData)
		systray.SetTooltip(common.T("app.name"))

		showMenuItem = systray.AddMenuItem(common.T("menu.showWindow"), "显示主窗口")
		systray.AddSeparator()
		quitMenuItem = systray.AddMenuItem("退出", "退出应用")

		// 监听菜单点击
		go func() {
			for {
				select {
				case <-showMenuItem.ClickedCh:
					if trayAppInstance != nil {
						trayAppInstance.ShowWindow()
					}
				case <-quitMenuItem.ClickedCh:
					if trayAppInstance != nil {
						trayAppInstance.ForceQuit()
					}
					systray.Quit()
					return
				}
			}
		}()
	}

	// 在单独的 goroutine 中运行托盘（因为 systray.Run 会阻塞）
	go systray.Run(onReady, onExit)
}

func onExit() {
	// 清理资源
}

// loadTrayIcon 加载托盘图标（从 embed 的资源文件）
func loadTrayIcon() []byte {
	if len(trayIconData) > 0 {
		return trayIconData
	}

	// 如果 embed 失败，尝试从文件系统读取
	data, err := os.ReadFile("appicon_64.png")
	if err == nil {
		return data
	}

	// 如果都失败，返回空（systray 会使用默认图标）
	return nil
}
