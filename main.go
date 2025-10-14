package main

import (
	"embed"
	"encoding/json"
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

	// 注册剪贴板监听器（后台持续运行）
	clipboardListener := common.RegisterClipboardListener()
	go func() {
		for newItem := range clipboardListener {
			log.Printf("📋 收到剪贴板更新通知: %s", truncateString(newItem.Content, 50))
		}
	}()

	// 注册全局快捷键
	go func() {
		// 等待应用完全启动后再注册快捷键
		// 从数据库读取快捷键设置
		settingsJSON, err := common.GetSetting("app_settings")
		if err == nil && settingsJSON != "" {
			var settings map[string]interface{}
			if err := json.Unmarshal([]byte(settingsJSON), &settings); err == nil {
				hotkey := "Control+V" // 默认快捷键
				if hotkeyVal, ok := settings["hotkey"].(string); ok && hotkeyVal != "" {
					hotkey = hotkeyVal
				}

				// 注册快捷键
				if err := common.RegisterHotkey(hotkey, func() {
					app.ShowWindow()
				}); err != nil {
					log.Printf("⚠️ 注册快捷键失败: %v", err)
				}
			}
		} else {
			// 使用默认快捷键
			if err := common.RegisterHotkey("Control+V", func() {
				app.ShowWindow()
			}); err != nil {
				log.Printf("⚠️ 注册默认快捷键失败: %v", err)
			}
		}
	}()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "剪存 - 剪贴板历史",
		Width:             1024,
		Height:            800,
		HideWindowOnClose: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 245, G: 245, B: 247, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: false,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "剪存",
				Message: "剪贴板历史管理工具\n版本 1.0.0",
			},
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
