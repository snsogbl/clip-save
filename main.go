package main

import (
	"embed"
	"log"

	"goWeb3/common"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 确保程序退出时关闭数据库
	defer func() {
		if err := common.CloseDB(); err != nil {
			log.Printf("关闭数据库失败: %v", err)
		}
		// 取消注册快捷键
		common.UnregisterHotkey()
	}()

	// Create an instance of the app structure
	app := NewApp()

	// 注册剪贴板（后台持续运行）
	clipboardListener := common.RegisterClipboardListener()
	go func() {
		for newItem := range clipboardListener {
			log.Printf("📋 收到剪贴板更新通知: %s", truncateString(newItem.Content, 50))
		}
	}()

	// 注册全局快捷键
	go app.RestartRegisterHotkey()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "剪存 - 剪贴板历史",
		Width:             1280,
		Height:            800,
		Frameless:         false,
		HideWindowOnClose: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Mac: &mac.Options{
			WebviewIsTransparent: false,
			About: &mac.AboutInfo{
				Title:   "剪存",
				Message: "剪贴板历史管理工具\n版本 1.0.5",
			},
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatal("启动 Wails 应用失败:", err)
	}
}

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}
