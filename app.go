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
	"log"
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

// GetCurrentLanguage 获取当前语言（供前端调用）
func (a *App) GetCurrentLanguage() (string, error) {
	return common.GetCurrentLanguage(), nil
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
			if err := cmd.Run(); err != nil {
				log.Printf("在资源管理器中打开目录失败: %v", err)
				return fmt.Errorf("打开目录失败: %v", err)
			}
			log.Printf("已在资源管理器中打开目录: %s", filePath)
			return nil
		}
		cmd := exec.Command("explorer", "/select,", filePath)
		if err := cmd.Run(); err != nil {
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
		// 使用 rundll32 调起默认浏览器，避免 cmd/start 的转义问题
		cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", decodedURL)
		if err := cmd.Run(); err != nil {
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

// HideWindow 隐藏窗口
func (a *App) HideWindow() {
	if a.ctx != nil {
		// Windows: 最小化而不是隐藏，确保任务栏图标可见
		if gRuntime.GOOS == "windows" {
			// runtime.WindowMinimise(a.ctx)
		} else {
			// 其他平台：保持原有隐藏行为
			runtime.Hide(a.ctx)
		}
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
