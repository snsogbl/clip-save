package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	gRuntime "runtime"
	"sync"
	"time"

	"goWeb3/common"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	qrcodegen "github.com/skip2/go-qrcode"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/clipboard"
)

// App struct
type App struct {
	ctx                  context.Context
	isWindowHidden       bool
	isUserSetAlwaysOnTop bool        // 用户是否设置了置顶
	sayProcess           *os.Process // 保存 say 进程对象
	sayProcessMutex      sync.Mutex  // 保护并发访问
}

// ShowAbout 显示关于对话框
func (a *App) ShowAbout() {
	if a.ctx == nil {
		return
	}
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   common.T("app.name"),
		Message: common.T("app.description") + "\n" + common.T("app.version"),
	})
}

// ShowSetting 显示设置对话框
func (a *App) ShowSetting() {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "nav.setting")
	}
}

// RunScript 显示脚本选择器
func (a *App) RunScript() {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "nav.runScript")
	}
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	// 添加 panic 恢复机制
	defer func() {
		if r := recover(); r != nil {
			log.Printf("startup 函数崩溃恢复: %v", r)
		}
	}()

	a.ctx = ctx
	log.Println("Wails 应用启动成功")

	// 延迟初始化应用切换监听器，确保 NSApplication 已完全初始化
	go func() {
		time.Sleep(500 * time.Millisecond)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("初始化应用切换监听器失败: %v", r)
			}
		}()
		common.InitAppSwitchListener()
	}()

	// 初始化统计模块 - 添加错误处理
	func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("初始化统计模块崩溃: %v", r)
			}
		}()
		if err := common.InitAnalytics(); err != nil {
			log.Printf("初始化统计模块失败: %v", err)
		}
	}()

	// 延迟注册 Dock 点击激活时的自动恢复与强退标记（仅 macOS 生效，其他平台为 no-op）
	// 确保 NSApplication 已完全初始化后再注册
	go func() {
		time.Sleep(500 * time.Millisecond)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("注册 Dock 重新打开监听失败: %v", r)
			}
		}()
		common.InitDockReopen(func() {
			a.ShowWindow()
		})
	}()
	common.SetForceQuitCallback(func() { common.SetForceQuit() })

	// 根据设置调整 Dock 图标可见性（仅 macOS 生效）
	if gRuntime.GOOS == "darwin" {
		go func() {
			// 延迟执行，确保应用已完全启动
			time.Sleep(300 * time.Millisecond)
			settingsJSON, err := common.GetSetting("app_settings")
			if err == nil && settingsJSON != "" {
				var settings map[string]interface{}
				if err := json.Unmarshal([]byte(settingsJSON), &settings); err == nil {
					if backgroundMode, ok := settings["backgroundMode"].(bool); ok && backgroundMode {
						// 开启后台模式：隐藏 Dock 图标
						common.SetDockIconVisibility(2)
						log.Println("已根据设置启用后台模式（隐藏 Dock 图标）")
					}
				}
			}
		}()
	}

	// 延迟调整窗口控制按钮位置，确保窗口已创建（仅 macOS 生效）
	go func() {
		time.Sleep(200 * time.Millisecond)
		common.AdjustWindowButtons()
		log.Println("已调整窗口控制按钮位置")
	}()

	// 设置脚本事件回调函数，用于发送事件到前端
	common.SetScriptEventCallback(func(eventName string, data interface{}) {
		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, eventName, data)
		}
	})
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	log.Println("Wails 应用关闭")
	// 停止脚本 HTTP 服务器
	if err := common.StopScriptHTTPServer(); err != nil {
		log.Printf("停止脚本 HTTP 服务器失败: %v", err)
	}
	if err := common.CloseDB(); err != nil {
		log.Printf("关闭数据库失败: %v", err)
	}
	// 清理窗口按钮观察者（仅 macOS 生效）
	common.CleanupWindowButtonsObserver()
}

// SearchClipboardItems 搜索剪贴板项目（供前端调用）
// loadImageData: 是否加载图片数据（极简模式下需要显示图片缩略图）
func (a *App) SearchClipboardItems(isFavorite bool, keyword string, filterType string, limit int, loadImageData bool) ([]common.ClipboardItem, error) {
	items, err := common.SearchClipboardItems(isFavorite, keyword, filterType, limit, loadImageData)
	if err != nil {
		log.Printf("搜索剪贴板项目失败: %v", err)
		return []common.ClipboardItem{}, err
	}
	return items, nil
}

// ToggleFavorite 切换收藏状态（供前端调用）
func (a *App) ToggleFavorite(id string) (int, error) {
	newVal, err := common.ToggleFavorite(id)
	if err != nil {
		log.Printf("切换收藏失败: %v", err)
		return 0, err
	}
	return newVal, nil
}

// GetClipboardItems 获取剪贴板项目列表（供前端调用）
func (a *App) GetClipboardItems(limit int) ([]common.ClipboardItem, error) {
	items, err := common.GetClipboardItems(limit)
	if err != nil {
		log.Printf("获取剪贴板项目失败: %v", err)
		return []common.ClipboardItem{}, err
	}
	return items, nil
}

// GetClipboardItemByID 根据ID获取剪贴板项目（供前端调用）
func (a *App) GetClipboardItemByID(id string) (*common.ClipboardItem, error) {
	item, err := common.GetClipboardItemByID(id)
	if err != nil {
		log.Printf("获取剪贴板项目失败: %v", err)
		return nil, err
	}
	return item, nil
}

// DeleteClipboardItem 删除剪贴板项目（供前端调用）
func (a *App) DeleteClipboardItem(id string) error {
	err := common.DeleteClipboardItem(id)
	if err != nil {
		log.Printf("删除剪贴板项目失败: %v", err)
		return err
	}
	return nil
}

// CopyTextToClipboard 复制文本到剪贴板（供前端调用）
func (a *App) CopyTextToClipboard(text string) error {
	clipboard.Write(clipboard.FmtText, []byte(text))
	log.Printf("已复制文本到剪贴板: %s", text)
	return nil
}

// CopyToClipboard 复制项目到剪贴板（供前端调用）
func (a *App) CopyToClipboard(id string) error {
	item, err := common.GetClipboardItemByID(id)
	if err != nil {
		return fmt.Errorf("获取项目失败: %v", err)
	}

	// 根据类型复制到剪贴板
	if item.ContentType == "Image" && len(item.ImageData) > 0 {
		// 复制图片
		clipboard.Write(clipboard.FmtImage, []byte(item.ImageData))
		log.Printf("已复制图片到剪贴板: %s", id)
	} else if item.ContentType == "File" && item.FilePaths != "" {
		// 复制文件（不是文本，而是真实的文件 URL）
		err := common.WriteFileURLs(item.FilePaths)
		if err != nil {
			log.Printf("复制文件失败: %v", err)
			return fmt.Errorf("复制文件失败: %v", err)
		}
		log.Printf("已复制文件到剪贴板: %s", id)
	} else {
		// 复制文本
		clipboard.Write(clipboard.FmtText, []byte(item.Content))
		log.Printf("已复制文本到剪贴板: %s", id)
	}

	return nil
}

// GetStatistics 获取统计信息（供前端调用）
func (a *App) GetStatistics() (map[string]interface{}, error) {
	stats, err := common.GetStatistics()
	if err != nil {
		log.Printf("获取统计信息失败: %v", err)
		return nil, err
	}
	return stats, nil
}

// ClearItemsOlderThanDays 清除超过指定天数的项目（供前端调用）
func (a *App) ClearItemsOlderThanDays(days int) error {
	err := common.ClearItemsOlderThanDays(days)
	if err != nil {
		log.Printf("清除超过 %d 天的项目失败: %v", days, err)
		return err
	}
	return nil
}

// ClearAllItems 清除所有剪贴板项目（供前端调用）
func (a *App) ClearAllItems() error {
	err := common.ClearAllItems()
	if err != nil {
		log.Printf("清除所有项目失败: %v", err)
		return err
	}
	return nil
}

// SaveAppSettings 保存应用设置（供前端调用）
func (a *App) SaveAppSettings(settingsJSON string) error {
	err := common.SaveSetting("app_settings", settingsJSON)
	if err != nil {
		log.Printf("保存应用设置失败: %v", err)
		return err
	}
	log.Printf("已保存应用设置")
	return nil
}

// GetAppSettings 获取应用设置（供前端调用）
func (a *App) GetAppSettings() (string, error) {
	settings, err := common.GetSetting("app_settings")
	if err != nil {
		log.Printf("获取应用设置失败: %v", err)
		return "", err
	}
	return settings, nil
}

// GetCurrentLanguage 获取当前语言（供前端调用）
func (a *App) GetCurrentLanguage() (string, error) {
	return common.GetCurrentLanguage(), nil
}

// SetAutoStart 设置开机自启（仅 Windows 有效，其他平台 no-op）
func (a *App) SetAutoStart(enable bool) error {
	return common.SetAutoStart(enable)
}

// IsAutoStartEnabled 查询开机自启状态（供前端调用）
func (a *App) IsAutoStartEnabled() (bool, error) {
	return common.IsAutoStartEnabled()
}

// SetLanguage 设置语言（供前端调用）
func (a *App) SetLanguage(lang string) error {
	err := common.SetLanguage(lang)
	if err != nil {
		log.Printf("设置语言失败: %v", err)
		return err
	}
	log.Printf("语言已设置为: %s", lang)
	return nil
}

// SetDockIconVisibility 设置 Dock 图标可见性（供前端调用，仅 macOS 生效）
func (a *App) SetDockIconVisibility(visible int) error {
	common.SetDockIconVisibility(visible)
	log.Printf("Dock 图标可见性已设置为: %d", visible)

	// 设置后台模式后，确保窗口仍然显示（不自动隐藏）
	if visible == 2 && a.ctx != nil {
		// 延迟确保 Activation Policy 设置完成后再显示窗口
		go func() {
			time.Sleep(10 * time.Millisecond)
			runtime.WindowShow(a.ctx)
			runtime.WindowUnminimise(a.ctx)
			log.Println("✅ 后台模式设置后，窗口已保持显示")
		}()
	}

	return nil
}

// GetSupportedLanguages 获取支持的语言列表（供前端调用）
func (a *App) GetSupportedLanguages() ([]string, error) {
	return common.GetSupportedLanguages(), nil
}

// VerifyPassword 验证密码（供前端调用）
func (a *App) VerifyPassword(password string) (bool, error) {
	// 获取设置
	settingsJSON, err := common.GetSetting("app_settings")
	if err != nil {
		log.Printf("获取设置失败: %v", err)
		return false, err
	}

	if settingsJSON == "" {
		// 没有设置，密码验证失败
		return false, nil
	}

	// 解析设置
	var settings map[string]interface{}
	if err := json.Unmarshal([]byte(settingsJSON), &settings); err != nil {
		log.Printf("解析设置失败: %v", err)
		return false, err
	}

	// 获取存储的密码hash
	storedPassword, ok := settings["password"].(string)
	if !ok || storedPassword == "" {
		// 没有设置密码
		return false, nil
	}

	// 计算输入密码的hash
	inputHash := hashPassword(password)

	// 比较hash
	isValid := inputHash == storedPassword
	if isValid {
		log.Println("✅ 密码验证成功")
	} else {
		log.Println("❌ 密码验证失败")
	}

	return isValid, nil
}

// hashPassword 计算密码的SHA-256哈希
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// OpenFileInFinder 在系统文件管理器中显示/打开文件（供前端调用）
func (a *App) OpenFileInFinder(filePath string) error {
	switch gRuntime.GOOS {
	case "darwin":
		// macOS: Finder
		cmd := exec.Command("open", "-R", filePath)
		if err := cmd.Run(); err != nil {
			log.Printf("在 Finder 中打开文件失败: %v", err)
			return fmt.Errorf("打开文件失败: %v", err)
		}
		log.Printf("已在 Finder 中打开文件: %s", filePath)
		return nil
	case "windows":
		// Windows: Explorer，/select, 展示并选中文件
		// 如果是目录，则直接打开目录
		if fi, err := os.Stat(filePath); err == nil && fi.IsDir() {
			cmd := exec.Command("explorer", filePath)
			if err := cmd.Start(); err != nil {
				log.Printf("在资源管理器中打开目录失败: %v", err)
				return fmt.Errorf("打开目录失败: %v", err)
			}
			log.Printf("已在资源管理器中打开目录: %s", filePath)
			return nil
		}
		cmd := exec.Command("explorer", "/select,", filePath)
		// 使用 Start 避免捕获 explorer 的非零退出码导致误报
		if err := cmd.Start(); err != nil {
			log.Printf("在资源管理器中显示文件失败: %v", err)
			return fmt.Errorf("打开文件失败: %v", err)
		}
		log.Printf("已在资源管理器中显示文件: %s", filePath)
		return nil
	default:
		// Linux: xdg-open 直接打开路径
		cmd := exec.Command("xdg-open", filePath)
		if err := cmd.Run(); err != nil {
			log.Printf("在文件管理器中打开失败: %v", err)
			return fmt.Errorf("打开文件失败: %v", err)
		}
		log.Printf("已在文件管理器中打开: %s", filePath)
		return nil
	}
}

// GetFileInfo 获取文件详细信息（供前端调用）
func (a *App) GetFileInfo(id string) ([]common.FileInfo, error) {
	item, err := common.GetClipboardItemByID(id)
	if err != nil {
		return nil, fmt.Errorf("获取项目失败: %v", err)
	}

	if item.ContentType != "File" || item.FileInfo == "" {
		return nil, fmt.Errorf("不是文件类型")
	}

	var fileInfos []common.FileInfo
	if err := json.Unmarshal([]byte(item.FileInfo), &fileInfos); err != nil {
		return nil, fmt.Errorf("解析文件信息失败: %v", err)
	}

	return fileInfos, nil
}

// OpenURL 在默认浏览器中打开 URL（供前端调用）
func (a *App) OpenURL(urlStr string) error {
	// 尝试解码 URL（如果已经被编码）
	decodedURL, err := url.QueryUnescape(urlStr)
	if err != nil {
		// 如果解码失败，使用原始 URL
		log.Printf("URL 解码失败，使用原始 URL: %v", err)
		decodedURL = urlStr
	}

	switch gRuntime.GOOS {
	case "darwin":
		cmd := exec.Command("open", decodedURL)
		if err := cmd.Run(); err != nil {
			log.Printf("打开 URL 失败: %v (原始: %s, 解码后: %s)", err, urlStr, decodedURL)
			return fmt.Errorf("打开 URL 失败: %v", err)
		}
		log.Printf("已在浏览器中打开 URL: %s (原始: %s)", decodedURL, urlStr)
		return nil
	case "windows":
		// 使用 rundll32 调起默认浏览器；用 Start 避免非零退出码误报
		cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", decodedURL)
		if err := cmd.Start(); err != nil {
			log.Printf("在 Windows 打开 URL 失败: %v (原始: %s, 解码后: %s)", err, urlStr, decodedURL)
			return fmt.Errorf("打开 URL 失败: %v", err)
		}
		log.Printf("已在浏览器中打开 URL: %s (原始: %s)", decodedURL, urlStr)
		return nil
	default:
		// Linux: xdg-open
		cmd := exec.Command("xdg-open", decodedURL)
		if err := cmd.Run(); err != nil {
			log.Printf("在 Linux 打开 URL 失败: %v (原始: %s, 解码后: %s)", err, urlStr, decodedURL)
			return fmt.Errorf("打开 URL 失败: %v", err)
		}
		log.Printf("已在浏览器中打开 URL: %s (原始: %s)", decodedURL, urlStr)
		return nil
	}
}

// SayText 使用 macOS 的 say 命令播放文字（仅 macOS）
func (a *App) SayText(text string) error {
	if gRuntime.GOOS != "darwin" {
		return fmt.Errorf("say 命令仅在 macOS 系统上可用")
	}

	if text == "" {
		return fmt.Errorf("文本不能为空")
	}

	// 先停止之前的播放（如果有）
	a.StopSay()

	// 启动 say 命令
	cmd := exec.Command("say", text)
	if err := cmd.Start(); err != nil {
		log.Printf("启动播放失败: %v", err)
		return fmt.Errorf("启动播放失败: %v", err)
	}
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "say.started")
	}

	// 保存进程对象
	a.sayProcessMutex.Lock()
	a.sayProcess = cmd.Process
	a.sayProcessMutex.Unlock()

	// 在 goroutine 中等待进程结束，清除进程对象并发送事件
	go func() {
		cmd.Wait()
		a.sayProcessMutex.Lock()
		a.sayProcess = nil
		a.sayProcessMutex.Unlock()

		// 播放完成，发送事件通知前端
		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, "say.finished")
		}
	}()

	log.Printf("已启动播放，进程 ID: %d", cmd.Process.Pid)
	return nil
}

// StopSay 停止当前正在播放的 say 命令（仅 macOS）
func (a *App) StopSay() error {
	if gRuntime.GOOS != "darwin" {
		return fmt.Errorf("stop say 命令仅在 macOS 系统上可用")
	}

	a.sayProcessMutex.Lock()
	process := a.sayProcess
	a.sayProcess = nil // 先清除，避免重复调用
	a.sayProcessMutex.Unlock()

	if process == nil {
		// 没有正在播放的进程
		return nil
	}

	// 使用 Process.Kill() 终止进程（沙盒环境下可以终止自己启动的进程）
	if err := process.Kill(); err != nil {
		log.Printf("停止播放失败（进程可能已结束）: %v", err)
		return nil // 不返回错误，因为进程可能已经结束
	}

	// 停止成功，发送事件通知前端
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "say.stopped")
	}

	log.Printf("已停止播放（进程 ID: %d）", process.Pid)
	return nil
}

func (a *App) IsSayPlaying() bool {
	a.sayProcessMutex.Lock()
	defer a.sayProcessMutex.Unlock()
	return a.sayProcess != nil
}

// ShowWindow 显示并聚焦窗口（供快捷键调用）
func (a *App) ShowWindow() {
	if a.ctx != nil {
		// 在激活本应用之前，记录当前前台应用
		common.RecordPreviousAppPID()

		// 如果窗口之前是隐藏状态，需要移动到当前活动的桌面空间
		runtime.WindowShow(a.ctx)
		common.EnsureWindowOnCurrentScreen(a.ctx)

		runtime.WindowUnminimise(a.ctx)

		// 通知前端选中第一个列表项
		// 使用 goroutine 异步发送事件，避免在 CGO 回调中直接调用导致信号错误
		if a.isWindowHidden {
			go func() {
				// 短暂延迟确保窗口操作已完成
				time.Sleep(50 * time.Millisecond)
				if a.ctx != nil {
					runtime.EventsEmit(a.ctx, "window.show")
				}
			}()
		}

		// 清除隐藏标记
		a.isWindowHidden = false

		// 临时设置置顶，确保窗口获得焦点
		runtime.WindowSetAlwaysOnTop(a.ctx, true)

		// 如果用户没有设置置顶，延迟取消置顶
		if !a.isUserSetAlwaysOnTop {
			go func() {
				time.Sleep(100 * time.Millisecond)
				if a.ctx != nil {
					runtime.WindowSetAlwaysOnTop(a.ctx, false)
				}
			}()
		}
	}
}

// SetWindowAlwaysOnTop 设置窗口置顶状态（供前端调用，同时更新全局变量）
func (a *App) SetWindowAlwaysOnTop(alwaysOnTop bool) error {
	if a.ctx != nil {
		a.isUserSetAlwaysOnTop = alwaysOnTop
		runtime.WindowSetAlwaysOnTop(a.ctx, alwaysOnTop)
		return nil
	}
	return fmt.Errorf("窗口上下文不存在")
}

// PrevItem 菜单：上一条
func (a *App) PrevItem() {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "nav.prev")
	}
}

// NextItem 菜单：下一条
func (a *App) NextItem() {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "nav.next")
	}
}

// ForceQuit 标记强制退出并真正退出应用
func (a *App) ForceQuit() {
	if a.ctx == nil {
		return
	}
	common.SetForceQuit()
	runtime.Quit(a.ctx)
}

// SwitchLeftTab 菜单：切换列表
func (a *App) SwitchLeftTab(tab string) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "nav.switch", tab)
	}
}

// CopyCurrentItem 菜单：复制当前项
func (a *App) CopyCurrentItem() {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "copy.current")
	}
}

// DeleteCurrentItem 菜单：删除当前项
func (a *App) DeleteCurrentItem() {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "delete.current")
	}
}

// CollectCurrentItem 菜单：收藏当前项
func (a *App) CollectCurrentItem() {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "collect.current")
	}
}

// SearchItem 菜单：查找
func (a *App) SearchItem() {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "search.item")
	}
}

// EnterItem
func (a *App) EnterItem() {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "enter.item")
	}
}

// HideWindow 隐藏窗口
func (a *App) HideWindow() {
	if a.ctx != nil {
		// Windows: 最小化而不是隐藏，确保任务栏图标可见
		if gRuntime.GOOS == "windows" {
			// runtime.WindowMinimise(a.ctx)
		} else {
			a.isWindowHidden = true
			// 其他平台：保持原有隐藏行为
			runtime.WindowHide(a.ctx)
		}
	}
}

func (a *App) HideWindowAndQuit() {
	if a.ctx != nil {
		// Windows: 最小化而不是隐藏，确保任务栏图标可见
		if gRuntime.GOOS == "windows" {
			// runtime.WindowMinimise(a.ctx)
		} else {
			a.isWindowHidden = true
			// 其他平台：保持原有隐藏行为
			runtime.Hide(a.ctx)
		}
	}
}

func (a *App) AutoPasteCurrentItem() {
	if a.ctx != nil {
		go common.PasteCmdV()
	}
}

// 激活应用
func (a *App) ActivatePreviousApp() {
	if a.ctx != nil {
		go common.ActivatePreviousApp()
	}
}

// AutoPasteCurrentItemToPreviousApp 自动粘贴到之前的前台应用（直接发送到进程）
func (a *App) AutoPasteCurrentItemToPreviousApp() {
	if a.ctx != nil {
		go common.PasteCmdVToPreviousApp()
	}
}

// SaveImagePNG 通过系统对话框将 Base64 PNG 保存到本地（供前端调用）
func (a *App) SaveImagePNG(base64Data string, suggestedName string) (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("应用上下文未初始化")
	}

	// 生成默认文件名
	now := time.Now()
	pad := func(n int) string { return fmt.Sprintf("%02d", n) }
	defaultName := fmt.Sprintf(
		"clipboard-%d%s%s-%s%s%s.png",
		now.Year(), pad(int(now.Month())), pad(now.Day()),
		pad(now.Hour()), pad(now.Minute()), pad(now.Second()),
	)
	if suggestedName != "" {
		defaultName = suggestedName
	}

	// 弹出保存对话框
	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		DefaultFilename: defaultName,
		Filters: []runtime.FileFilter{
			{DisplayName: "PNG 图片", Pattern: "*.png"},
		},
	})
	if err != nil {
		return "", fmt.Errorf("选择保存路径失败: %v", err)
	}
	if path == "" {
		// 用户取消
		return "", nil
	}

	// 解码 Base64 数据
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("解码图片失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("写入文件失败: %v", err)
	}

	log.Printf("图片已保存到: %s", path)
	return path, nil
}

// DetectQRCode 检测图片中是否包含二维码（供前端调用）
func (a *App) DetectQRCode(base64Data string) (bool, error) {
	// 解码 Base64 数据
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return false, fmt.Errorf("解码图片失败: %v", err)
	}

	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return false, fmt.Errorf("解码图片失败: %v", err)
	}

	// 将图片转换为灰度图
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return false, fmt.Errorf("转换图片失败: %v", err)
	}

	// 创建二维码读取器
	reader := qrcode.NewQRCodeReader()

	// 尝试识别二维码
	_, err = reader.Decode(bmp, nil)
	if err != nil {
		// 如果没有找到二维码，返回false
		return false, nil
	}

	// 找到二维码
	return true, nil
}

// RecognizeQRCode 识别图片中的二维码内容（供前端调用）
func (a *App) RecognizeQRCode(base64Data string) (string, error) {
	// 解码 Base64 数据
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("解码图片失败: %v", err)
	}

	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("解码图片失败: %v", err)
	}

	// 将图片转换为灰度图
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", fmt.Errorf("转换图片失败: %v", err)
	}

	// 创建二维码读取器
	reader := qrcode.NewQRCodeReader()

	// 尝试识别二维码
	result, err := reader.Decode(bmp, nil)
	if err != nil {
		return "", fmt.Errorf("识别二维码失败: %v", err)
	}

	// 返回二维码内容
	return result.GetText(), nil
}

// GenerateQRCode 生成二维码并返回Base64编码的PNG图片（供前端调用）
func (a *App) GenerateQRCode(content string, size int) (string, error) {
	if content == "" {
		return "", fmt.Errorf("内容不能为空")
	}

	// 设置默认尺寸
	if size <= 0 {
		size = 256
	}

	// 生成二维码
	qr, err := qrcodegen.New(content, qrcodegen.Medium)
	if err != nil {
		return "", fmt.Errorf("生成二维码失败: %v", err)
	}

	// 转换为PNG
	img := qr.Image(size)

	// 编码为PNG
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return "", fmt.Errorf("编码PNG失败: %v", err)
	}

	// 转换为Base64
	base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
	return base64Str, nil
}

// CopyImageToClipboard 将Base64编码的图片复制到剪贴板（供前端调用）
func (a *App) CopyImageToClipboard(base64Data string) error {
	if base64Data == "" {
		return fmt.Errorf("图片数据不能为空")
	}

	// 解码Base64数据
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return fmt.Errorf("解码图片数据失败: %v", err)
	}

	// 写入剪贴板
	done := clipboard.Write(clipboard.FmtImage, data)
	<-done // 等待写入完成

	log.Printf("图片已复制到剪贴板，大小: %d bytes", len(data))
	return nil
}

// TranslateCurrentItem 翻译当前项（供前端调用）
func (a *App) TranslateCurrentItem() {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "translate.current")
	}
}

// PlayCurrentItem 播放当前项（供前端调用）
func (a *App) PlayCurrentItem() {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "play.current")
	}
}

// 添加互斥锁防止重复调用
var hotkeyRestartMutex sync.Mutex

func (a *App) RestartRegisterHotkey() error {
	// 使用互斥锁防止重复调用
	hotkeyRestartMutex.Lock()
	defer hotkeyRestartMutex.Unlock()

	log.Println("🔄 重启注册快捷键")

	// 先取消当前注册的快捷键
	common.UnregisterHotkey()

	// 等待一小段时间确保旧快捷键完全清理
	time.Sleep(100 * time.Millisecond)

	// 获取设置
	settingsJSON, err := common.GetSetting("app_settings")
	if err != nil {
		log.Printf("获取应用设置失败: %v", err)
	}

	// 解析设置
	var settings map[string]interface{}
	if err := json.Unmarshal([]byte(settingsJSON), &settings); err != nil {
		log.Printf("解析应用设置失败: %v", err)
	}

	// 获取快捷键设置
	hotkey := "Command+Option+c" // 默认快捷键
	if hotkeyVal, ok := settings["hotkey"].(string); ok && hotkeyVal != "" {
		hotkey = hotkeyVal
	}

	// 注册快捷键
	if err := common.RegisterHotkey(hotkey, func() {
		a.ShowWindow()
	}); err != nil {
		log.Printf("⚠️ 注册快捷键失败: %v", err)
		return fmt.Errorf("注册快捷键失败: %v", err)
	}

	log.Printf("✅ 快捷键注册成功: %s", hotkey)
	return nil
}

// GetAllUserScripts 获取所有用户脚本
func (a *App) GetAllUserScripts() ([]common.UserScript, error) {
	return common.GetAllUserScripts()
}

// GetEnabledUserScriptsByTrigger 根据触发类型获取启用的脚本
func (a *App) GetEnabledUserScriptsByTrigger(trigger string) ([]common.UserScript, error) {
	return common.GetEnabledUserScripts(trigger)
}

// UpdateUserScriptOrder 更新单个脚本顺序
func (a *App) UpdateUserScriptOrder(scriptID string, sortOrder int) error {
	return common.UpdateUserScriptOrder(scriptID, sortOrder)
}

// GetUserScriptByID 根据 ID 获取脚本
func (a *App) GetUserScriptByID(id string) (*common.UserScript, error) {
	return common.GetUserScriptByID(id)
}

// GetUserScriptsByIDs 根据 ID 列表批量获取脚本
func (a *App) GetUserScriptsByIDs(ids []string) ([]common.UserScript, error) {
	return common.GetUserScriptsByIDs(ids)
}

// SaveUserScript 保存用户脚本
func (a *App) SaveUserScript(scriptJSON string) error {
	var script common.UserScript
	if err := json.Unmarshal([]byte(scriptJSON), &script); err != nil {
		return fmt.Errorf("解析脚本数据失败: %v", err)
	}
	return common.SaveUserScript(&script)
}

// DeleteUserScript 删除用户脚本
func (a *App) DeleteUserScript(id string) error {
	// 如果脚本的 HTTP 服务已启用，先禁用它
	if common.IsScriptHTTPServiceEnabled(id) {
		if err := common.DisableScriptHTTPService(id); err != nil {
			log.Printf("禁用脚本 HTTP 服务失败: %v", err)
		}
	}
	return common.DeleteUserScript(id)
}

// EnableScriptHTTPService 启用脚本的 HTTP 服务（供前端调用）
func (a *App) EnableScriptHTTPService(scriptID string) error {
	return common.EnableScriptHTTPService(scriptID)
}

// DisableScriptHTTPService 禁用脚本的 HTTP 服务（供前端调用）
func (a *App) DisableScriptHTTPService(scriptID string) error {
	return common.DisableScriptHTTPService(scriptID)
}

// IsScriptHTTPServiceEnabled 检查脚本的 HTTP 服务是否已启用（供前端调用）
func (a *App) IsScriptHTTPServiceEnabled(scriptID string) (bool, error) {
	return common.IsScriptHTTPServiceEnabled(scriptID), nil
}

// GetScriptHTTPURL 获取脚本的 HTTP 服务 URL（供前端调用）
func (a *App) GetScriptHTTPURL(scriptID string) (string, error) {
	return common.GetScriptHTTPURL(scriptID)
}

// SetScriptHTTPResult 设置脚本执行结果（供前端调用）
func (a *App) SetScriptHTTPResult(requestID string, resultJSON string) error {
	var result common.ScriptHTTPResult
	if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
		return fmt.Errorf("解析结果失败: %v", err)
	}
	common.SetScriptHTTPResult(requestID, result)
	return nil
}

// HttpRequest 通用的 HTTP 请求代理函数（用于绕过 CORS 限制）
// method: HTTP 方法（GET, POST, PUT, DELETE 等）
// requestUrl: 请求 URL
// headersJson: 请求头 JSON 字符串，格式如 {"Content-Type": "application/json", "Authorization": "Bearer token"}
// bodyJson: 请求体 JSON 字符串（GET 请求可为空字符串）
// 返回响应体的 JSON 字符串和错误信息
func (a *App) HttpRequest(method string, requestUrl string, headersJson string, bodyJson string) (string, error) {
	// 解析请求头
	var headers map[string]string
	if headersJson != "" {
		if err := json.Unmarshal([]byte(headersJson), &headers); err != nil {
			return "", fmt.Errorf("解析请求头失败: %v", err)
		}
	} else {
		headers = make(map[string]string)
	}

	// 创建请求体
	var bodyReader io.Reader
	if bodyJson != "" && (method == "POST" || method == "PUT" || method == "PATCH") {
		bodyReader = bytes.NewBufferString(bodyJson)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest(method, requestUrl, bodyReader)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 如果没有设置 Content-Type 且有请求体，默认设置为 application/json
	if bodyReader != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// 发送请求
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 构建响应对象
	response := map[string]interface{}{
		"status":     resp.StatusCode,
		"statusText": resp.Status,
		"headers":    resp.Header,
		"body":       string(body),
	}

	// 如果响应是 JSON，尝试解析
	var bodyObj interface{}
	if err := json.Unmarshal(body, &bodyObj); err == nil {
		response["body"] = bodyObj
	}

	// 返回 JSON 格式的响应
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("序列化响应失败: %v", err)
	}

	return string(responseJSON), nil
}
