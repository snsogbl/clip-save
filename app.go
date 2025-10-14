package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os/exec"

	"goWeb3/common"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/clipboard"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	log.Println("Wails 应用启动成功")
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	log.Println("Wails 应用关闭")
	if err := common.CloseDB(); err != nil {
		log.Printf("关闭数据库失败: %v", err)
	}
}

// SearchClipboardItems 搜索剪贴板项目（供前端调用）
func (a *App) SearchClipboardItems(keyword string, filterType string, limit int) ([]common.ClipboardItem, error) {
	items, err := common.SearchClipboardItems(keyword, filterType, limit)
	if err != nil {
		log.Printf("搜索剪贴板项目失败: %v", err)
		return []common.ClipboardItem{}, err
	}
	return items, nil
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

// ClearOldItems 清除旧项目（供前端调用）
func (a *App) ClearOldItems(keepCount int) error {
	err := common.ClearOldItems(keepCount)
	if err != nil {
		log.Printf("清除旧项目失败: %v", err)
		return err
	}
	return nil
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

// OpenFileInFinder 在 Finder 中打开文件（供前端调用）
func (a *App) OpenFileInFinder(filePath string) error {
	// 使用 open -R 命令在 Finder 中显示文件
	cmd := exec.Command("open", "-R", filePath)
	err := cmd.Run()
	if err != nil {
		log.Printf("在 Finder 中打开文件失败: %v", err)
		return fmt.Errorf("打开文件失败: %v", err)
	}
	log.Printf("已在 Finder 中打开文件: %s", filePath)
	return nil
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

	// 使用 open 命令在默认浏览器中打开 URL
	cmd := exec.Command("open", decodedURL)
	err = cmd.Run()
	if err != nil {
		log.Printf("打开 URL 失败: %v (原始: %s, 解码后: %s)", err, urlStr, decodedURL)
		return fmt.Errorf("打开 URL 失败: %v", err)
	}
	log.Printf("已在浏览器中打开 URL: %s (原始: %s)", decodedURL, urlStr)
	return nil
}

// ShowWindow 显示并聚焦窗口（供快捷键调用）
func (a *App) ShowWindow() {
	if a.ctx != nil {
		runtime.WindowShow(a.ctx)
		runtime.WindowUnminimise(a.ctx)
		runtime.WindowSetAlwaysOnTop(a.ctx, true)
		// 延迟取消置顶，确保窗口获得焦点
		go func() {
			// 短暂延迟后取消置顶
			runtime.WindowSetAlwaysOnTop(a.ctx, false)
		}()
		log.Println("🪟 窗口已显示并聚焦")
	}
}
