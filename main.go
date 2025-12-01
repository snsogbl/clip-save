package main

import (
	"context"
	"embed"
	"log"
	"runtime"

	"goWeb3/common"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 判断是否是 macOS
	isMac := runtime.GOOS == "darwin"

	// 初始化国际化
	if err := common.InitI18n(); err != nil {
		log.Fatal("初始化国际化失败:", err)
	}

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

	appMenu := menu.NewMenu()

	appSubMenu := appMenu.AddSubmenu(common.T("app.name"))
	appSubMenu.AddText("About "+common.T("app.name"), keys.CmdOrCtrl("b"), func(_ *menu.CallbackData) {
		app.ShowAbout()
	})
	if isMac {
		appSubMenu.AddText("Hide "+common.T("app.name"), keys.CmdOrCtrl("h"), func(_ *menu.CallbackData) {
			app.HideWindow()
		})
		appSubMenu.AddText("Show "+common.T("app.name"), keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
			app.ShowWindow()
		})
	}
	appSubMenu.AddSeparator()
	appSubMenu.AddText("Setting "+common.T("app.name"), keys.CmdOrCtrl(","), func(_ *menu.CallbackData) {
		app.ShowSetting()
	})
	appSubMenu.AddSeparator()
	appSubMenu.AddText("Quit "+common.T("app.name"), keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		app.ForceQuit()
	})

	// 添加默认编辑菜单以支持标准复制粘贴快捷键
	appMenu.Append(menu.EditMenu())

	displaySubMenu := appMenu.AddSubmenu(common.T("menu.display"))
	if isMac {
		displaySubMenu.AddText(common.T("menu.showWindow"), keys.CmdOrCtrl("0"), func(_ *menu.CallbackData) {
			app.ShowWindow()
		})
		displaySubMenu.AddSeparator()
	}
	displaySubMenu.AddText(common.T("menu.list"), keys.CmdOrCtrl("left"), func(_ *menu.CallbackData) {
		app.SwitchLeftTab("all")
	})
	displaySubMenu.AddText(common.T("menu.favorite"), keys.CmdOrCtrl("right"), func(_ *menu.CallbackData) {
		app.SwitchLeftTab("fav")
	})
	displaySubMenu.AddSeparator()
	displaySubMenu.AddText(common.T("menu.prev"), keys.CmdOrCtrl("up"), func(_ *menu.CallbackData) {
		app.PrevItem()
	})
	displaySubMenu.AddText(common.T("menu.next"), keys.CmdOrCtrl("down"), func(_ *menu.CallbackData) {
		app.NextItem()
	})
	displaySubMenu.AddSeparator()
	displaySubMenu.AddText(common.T("menu.search"), keys.CmdOrCtrl("f"), func(_ *menu.CallbackData) {
		app.SearchItem()
	})
	displaySubMenu.AddSeparator()
	displaySubMenu.AddText(common.T("menu.copyCurrent"), keys.CmdOrCtrl("enter"), func(_ *menu.CallbackData) {
		app.CopyCurrentItem()
		app.HideWindowAndQuit()
		app.AutoPasteCurrentItem()
	})
	displaySubMenu.AddText(common.T("menu.deleteCurrent"), keys.CmdOrCtrl("backspace"), func(_ *menu.CallbackData) {
		app.DeleteCurrentItem()
	})
	displaySubMenu.AddText(common.T("menu.favoriteCurrent"), keys.CmdOrCtrl("d"), func(_ *menu.CallbackData) {
		app.CollectCurrentItem()
	})
	displaySubMenu.AddText(common.T("menu.translateCurrent"), keys.CmdOrCtrl("t"), func(_ *menu.CallbackData) {
		app.TranslateCurrentItem()
	})

	// 注册剪贴板（后台持续运行）
	clipboardListener := common.RegisterClipboardListener()
	go func() {
		for range clipboardListener {
			// 向前端发送剪贴板更新事件，触发前端刷新
			if app.ctx != nil {
				wailsRuntime.EventsEmit(app.ctx, "clipboard.updated")
			}
		}
	}()

	// 注册全局快捷键
	go app.RestartRegisterHotkey()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     common.T("app.name"),
		Width:     1280,
		Height:    800,
		Frameless: false,
		OnBeforeClose: func(ctx context.Context) bool {
			if isMac && !common.IsForceQuit() {
				app.HideWindow()
				return true
			}
			return false
		},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		Menu:             appMenu,
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			WebviewIsTransparent: false,
			About: &mac.AboutInfo{
				Title:   common.T("app.name"),
				Message: common.T("app.description") + "\n版本 " + common.T("app.version"),
			},
		},
		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			BackdropType:                      windows.Mica,
			DisablePinchZoom:                  false,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme:                             windows.SystemDefault,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar:   windows.RGB(20, 20, 20),
				DarkModeTitleText:  windows.RGB(200, 200, 200),
				DarkModeBorder:     windows.RGB(20, 0, 20),
				LightModeTitleBar:  windows.RGB(200, 200, 200),
				LightModeTitleText: windows.RGB(20, 20, 20),
				LightModeBorder:    windows.RGB(200, 200, 200),
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
