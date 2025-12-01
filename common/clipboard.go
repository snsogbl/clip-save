package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"  // GIF 格式支持
	_ "image/jpeg" // JPEG 格式支持
	"image/png"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"

	_ "golang.org/x/image/bmp"  // BMP 格式支持
	_ "golang.org/x/image/tiff" // TIFF 格式支持
	_ "golang.org/x/image/webp" // WebP 格式支持

	"golang.design/x/clipboard"
)

// ClipboardItem 剪贴板项目结构
type ClipboardItem struct {
	ID          string
	Content     string
	ContentType string
	ContentHash string // 内容哈希值，用于去重
	ImageData   []byte // 图片数据（PNG格式）
	FilePaths   string // 文件路径（JSON 数组格式）
	FileInfo    string // 文件信息（JSON 格式）
	Timestamp   time.Time
	Source      string
	CharCount   int
	WordCount   int
	IsFavorite  int // 0/1
}

// 剪贴板更新通知监听器（只发送信号，不传递数据）
var clipboardListener chan struct{}

func init() {
	// 初始化剪贴板
	err := clipboard.Init()
	if err != nil {
		log.Printf("初始化剪贴板失败: %v", err)
		return
	}

	// 初始化数据库
	if err := InitDB(); err != nil {
		log.Printf("数据库初始化失败: %v", err)
	}

	// 启动剪贴板
	go run()
}

// RegisterClipboardListener 注册剪贴板更新监听器
func RegisterClipboardListener() chan struct{} {
	clipboardListener = make(chan struct{}, 10)
	return clipboardListener
}

// notifyListeners 通知监听器（只发送信号）
func notifyListeners() {
	if clipboardListener != nil {
		select {
		case clipboardListener <- struct{}{}:
		default:
			// 如果 channel 已满，跳过这次通知
			log.Printf("channel 已满，跳过通知")
		}
	}
}

func run() {
	var lastTextContent string
	var lastImageHash string
	var lastFileHash string
	var lastPasteboardChangeCount int

	// 用于追踪应用切换历史
	// var currentAppName string
	// var previousAppName string
	// var appSwitchTime time.Time

	// 缩短轮询间隔到 50ms，以便更及时地捕获剪贴板变化
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		// 使用 changeCount 精确检测剪贴板是否变化
		currentChangeCount := GetPasteboardChangeCount()
		// log.Printf("🔄 剪贴板变化计数: %d", currentChangeCount)
		if currentChangeCount == lastPasteboardChangeCount {
			// 剪贴板没有变化，继续下一次循环
			continue
		}
		// 获取当前活动应用
		lastPasteboardChangeCount = currentChangeCount

		// 获取当前活动应用
		sourceAppName := GetFrontmostAppName()

		// 优先级1: 先检查是否有文件
		fileJSON, fileCount := ReadFileURLs()
		if fileCount > 0 && fileJSON != "" {
			// 使用完整路径集合的稳定哈希，避免前缀相同导致的误判
			fileHash := calculateFilePathsHash(fileJSON)
			if fileHash != lastFileHash {
				lastFileHash = fileHash
				lastTextContent = ""
				lastImageHash = ""
				handleFileClipboard(fileJSON, fileCount, sourceAppName, fileHash)
			}
			continue
		}

		// 优先级2: 检查是否有图片 - 尝试多种格式
		imgData := tryReadImage()
		if len(imgData) > 0 {
			// 统一转换为PNG格式后计算哈希，确保相同图片内容产生相同哈希值
			pngData, err := convertToPNG(imgData)
			if err != nil {
				log.Printf("❌ 转换图片为PNG失败: %v", err)
				continue
			}
			// 对PNG数据计算哈希值来判断是否是新图片
			h := sha256.Sum256(pngData)
			imageHash := hex.EncodeToString(h[:])
			if imageHash != lastImageHash {
				lastImageHash = imageHash
				lastTextContent = ""
				lastFileHash = ""

				handleImageClipboard(imgData, sourceAppName, imageHash)
			}
		} else {
			// 优先级3: 没有图片，检查文本
			textData := clipboard.Read(clipboard.FmtText)
			if len(textData) > 0 {
				content := string(textData)
				if content != lastTextContent && content != "" {
					lastTextContent = content
					lastImageHash = ""
					lastFileHash = ""
					handleTextClipboard(content, sourceAppName)
				}
			}
		}
	}
}

// tryReadImage 尝试从剪贴板读取图片，支持多种格式
func tryReadImage() []byte {
	// 常见的图片 UTI 类型
	imageTypes := []string{
		"public.tiff",        // TIFF 格式（macOS 常用）
		"public.png",         // PNG 格式
		"public.jpeg",        // JPEG 格式
		"com.compuserve.gif", // GIF 格式
		"public.image",       // 通用图片类型
		"com.apple.pict",     // PICT 格式（旧 macOS）
		"com.microsoft.bmp",  // BMP 格式
	}

	// 尝试读取各种图片类型
	for _, imageType := range imageTypes {
		imgData := ReadPasteboardData(imageType)
		if len(imgData) > 0 {
			return imgData
		}
	}

	// 如果没有找到图片，尝试标准的 clipboard.FmtImage
	imgData := clipboard.Read(clipboard.FmtImage)
	if len(imgData) > 0 {
		log.Printf("🎨 从 clipboard.FmtImage 读取到图片，大小: %d bytes", len(imgData))
		return imgData
	}

	return nil
}

// convertToPNG 将任意格式的图片数据转换为PNG格式
func convertToPNG(imgData []byte) ([]byte, error) {
	// 解码图片
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return nil, fmt.Errorf("解码图片失败: %v", err)
	}

	// 转换为PNG格式
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, fmt.Errorf("编码PNG失败: %v", err)
	}

	return buf.Bytes(), nil
}

// handleTextClipboard 处理文本剪贴板
func handleTextClipboard(content string, appName string) {
	timestamp := time.Now()
	item := ClipboardItem{
		ID:          fmt.Sprintf("%d", timestamp.UnixNano()),
		Content:     content,
		ContentType: detectContentType(content),
		Timestamp:   timestamp,
		Source:      appName,
		CharCount:   len([]rune(content)),
		WordCount:   countWords(content),
	}

	// 计算内容哈希
	item.ContentHash = calculateContentHash(&item)

	// log.Printf("📝 新文本剪贴板: %s, 类型: %s", truncateString(item.Content, 50), item.ContentType)

	// 保存到数据库
	if err := SaveClipboardItem(&item); err != nil {
		log.Printf("保存剪贴板内容失败: %v", err)
	} else {
		// 通知监听器
		notifyListeners()
	}
}

// handleImageClipboard 处理图片剪贴板
func handleImageClipboard(imgData []byte, appName string, precomputedHash string) {
	// 解码图片
	img, format, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Printf("❌ 解码图片失败: %v (数据头: %X)", err, imgData[:min(16, len(imgData))])
		return
	}

	// 转换为PNG格式存储
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		log.Printf("❌ 编码PNG失败: %v", err)
		return
	}

	timestamp := time.Now()
	pngData := buf.Bytes()

	// 生成缩略图描述
	bounds := img.Bounds()
	imageDesc := fmt.Sprintf("图片 %dx%d (%s)", bounds.Dx(), bounds.Dy(), format)

	item := ClipboardItem{
		ID:          fmt.Sprintf("%d", timestamp.UnixNano()),
		Content:     imageDesc,
		ContentType: "Image",
		ImageData:   pngData,
		Timestamp:   timestamp,
		Source:      appName,
		CharCount:   len(pngData),
		WordCount:   0,
	}

	// 计算内容哈希（优先使用外部预计算避免重复开销）
	if precomputedHash != "" {
		item.ContentHash = precomputedHash
	} else {
		item.ContentHash = calculateContentHash(&item)
	}

	// 保存到数据库
	if err := SaveClipboardItem(&item); err != nil {
		log.Printf("❌ 保存图片剪贴板失败: %v", err)
	} else {
		// 通知监听器
		notifyListeners()
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

// min 返回较小的整数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// detectContentType 检测内容类型
func detectContentType(content string) string {
	if len(content) == 0 {
		return "Text"
	}

	// 去除首尾空白
	content = strings.TrimSpace(content)

	// 检测是否为 JSON（对象或数组）
	if (strings.HasPrefix(content, "{") && strings.HasSuffix(content, "}")) ||
		(strings.HasPrefix(content, "[") && strings.HasSuffix(content, "]")) {
		var js interface{}
		if json.Unmarshal([]byte(content), &js) == nil {
			return "JSON"
		}
	}

	// 检测是否为URL
	if len(content) > 4 && (content[:4] == "http" || content[:3] == "www") {
		return "URL"
	}

	// 检测是否为颜色代码
	if isColorFormat(content) {
		return "Color"
	}

	// 默认为文本
	return "Text"
}

// isColorFormat 检测是否为颜色格式
func isColorFormat(content string) bool {
	// HEX 格式: #RGB 或 #RRGGBB 或 #RRGGBBAA
	hexPattern := regexp.MustCompile(`^#([0-9A-Fa-f]{3}|[0-9A-Fa-f]{6}|[0-9A-Fa-f]{8})$`)
	if hexPattern.MatchString(content) {
		return true
	}

	// RGB 格式: rgb(r, g, b)
	rgbPattern := regexp.MustCompile(`^rgb\s*\(\s*\d{1,3}\s*,\s*\d{1,3}\s*,\s*\d{1,3}\s*\)$`)
	if rgbPattern.MatchString(content) {
		return true
	}

	// RGBA 格式: rgba(r, g, b, a)
	rgbaPattern := regexp.MustCompile(`^rgba\s*\(\s*\d{1,3}\s*,\s*\d{1,3}\s*,\s*\d{1,3}\s*,\s*[0-9.]+\s*\)$`)
	if rgbaPattern.MatchString(content) {
		return true
	}

	// HSL 格式: hsl(h, s%, l%)
	hslPattern := regexp.MustCompile(`^hsl\s*\(\s*\d{1,3}\s*,\s*\d{1,3}%\s*,\s*\d{1,3}%\s*\)$`)
	if hslPattern.MatchString(content) {
		return true
	}

	// HSLA 格式: hsla(h, s%, l%, a)
	hslaPattern := regexp.MustCompile(`^hsla\s*\(\s*\d{1,3}\s*,\s*\d{1,3}%\s*,\s*\d{1,3}%\s*,\s*[0-9.]+\s*\)$`)
	return hslaPattern.MatchString(content)
}

// countWords 统计单词数（智能识别中英文）
// 中文/日文/韩文等字符按字数统计，英文按单词统计
func countWords(content string) int {
	if len(content) == 0 {
		return 0
	}

	count := 0
	inWord := false

	for _, r := range content {
		// 判断是否为 CJK 字符（中文、日文、韩文）
		if isCJK(r) {
			// 如果之前在处理英文单词，先结算这个单词
			if inWord {
				count++
				inWord = false
			}
			// CJK 字符每个都算一个"单词"
			count++
		} else if isWordCharacter(r) {
			// 英文字母、数字等
			if !inWord {
				inWord = true
			}
		} else {
			// 空格、标点符号等分隔符
			if inWord {
				count++
				inWord = false
			}
		}
	}

	// 处理最后一个单词
	if inWord {
		count++
	}

	return count
}

// isCJK 判断是否为中日韩字符
func isCJK(r rune) bool {
	return (r >= 0x4E00 && r <= 0x9FFF) || // CJK 统一表意文字
		(r >= 0x3400 && r <= 0x4DBF) || // CJK 扩展 A
		(r >= 0x20000 && r <= 0x2A6DF) || // CJK 扩展 B
		(r >= 0x2A700 && r <= 0x2B73F) || // CJK 扩展 C
		(r >= 0x2B740 && r <= 0x2B81F) || // CJK 扩展 D
		(r >= 0x2B820 && r <= 0x2CEAF) || // CJK 扩展 E
		(r >= 0xF900 && r <= 0xFAFF) || // CJK 兼容表意文字
		(r >= 0x2F800 && r <= 0x2FA1F) || // CJK 兼容表意文字补充
		(r >= 0x3040 && r <= 0x309F) || // 日文平假名
		(r >= 0x30A0 && r <= 0x30FF) || // 日文片假名
		(r >= 0xAC00 && r <= 0xD7AF) // 韩文音节
}

// isWordCharacter 判断是否为单词字符（字母、数字、下划线等）
func isWordCharacter(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-'
}

// FileInfo 文件信息结构
type FileInfo struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	IsDir     bool   `json:"is_dir"`
	Exists    bool   `json:"exists"`
	Extension string `json:"extension"`
}

// handleFileClipboard 处理文件剪贴板
func handleFileClipboard(fileJSON string, fileCount int, appName string, precomputedHash string) {
	// 解析文件路径列表
	var filePaths []string
	if err := json.Unmarshal([]byte(fileJSON), &filePaths); err != nil {
		log.Printf("❌ 解析文件路径失败: %v", err)
		return
	}

	if len(filePaths) == 0 {
		return
	}

	// 收集文件信息
	fileInfos := make([]FileInfo, 0, len(filePaths))
	var totalSize int64

	for _, path := range filePaths {
		info := getFileInfo(path)
		fileInfos = append(fileInfos, info)
		totalSize += info.Size
	}

	// 生成内容描述
	var content string
	if len(filePaths) == 1 {
		info := fileInfos[0]
		if info.IsDir {
			content = fmt.Sprintf("📁 %s", info.Name)
		} else {
			content = fmt.Sprintf("📄 %s (%s)", info.Name, formatFileSize(info.Size))
		}
	} else {
		content = fmt.Sprintf("📦 %d 个文件/文件夹 (%s)", len(filePaths), formatFileSize(totalSize))
	}

	// 序列化文件信息为 JSON
	fileInfoJSON, err := json.Marshal(fileInfos)
	if err != nil {
		log.Printf("❌ 序列化文件信息失败: %v", err)
		return
	}

	timestamp := time.Now()
	item := ClipboardItem{
		ID:          fmt.Sprintf("%d", timestamp.UnixNano()),
		Content:     content,
		ContentType: "File",
		FilePaths:   fileJSON,
		FileInfo:    string(fileInfoJSON),
		Timestamp:   timestamp,
		Source:      appName,
		CharCount:   len(content),
		WordCount:   len(filePaths),
	}

	// 计算内容哈希（优先使用外部预计算避免重复开销）
	if precomputedHash != "" {
		item.ContentHash = precomputedHash
	} else {
		item.ContentHash = calculateContentHash(&item)
	}

	log.Printf("📁 新文件剪贴板: %s", content)

	// 保存到数据库
	if err := SaveClipboardItem(&item); err != nil {
		log.Printf("❌ 保存文件剪贴板失败: %v", err)
	} else {
		// 通知监听器
		notifyListeners()
	}
}

// getFileInfo 获取文件信息
func getFileInfo(path string) FileInfo {
	info := FileInfo{
		Name: filepath.Base(path),
		Path: path,
	}

	stat, err := os.Stat(path)
	if err != nil {
		// 文件不存在或无法访问
		info.Exists = false
		info.Extension = filepath.Ext(path)
		return info
	}

	info.Exists = true
	info.IsDir = stat.IsDir()
	info.Size = stat.Size()

	if !info.IsDir {
		info.Extension = strings.ToLower(filepath.Ext(path))
	}

	return info
}

// formatFileSize 格式化文件大小
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// calculateContentHash 计算剪贴板项目的内容哈希值
func calculateContentHash(item *ClipboardItem) string {
	switch item.ContentType {
	case "Text", "URL", "Color":
		// 文本类型直接对内容计算哈希
		hash := sha256.Sum256([]byte(item.Content))
		return hex.EncodeToString(hash[:])
	case "Image":
		// 图片类型对图片数据计算哈希
		if len(item.ImageData) == 0 {
			return ""
		}
		hash := sha256.Sum256(item.ImageData)
		return hex.EncodeToString(hash[:])
	case "File":
		// 文件类型对排序后的文件路径计算哈希
		return calculateFilePathsHash(item.FilePaths)
	default:
		// 其他类型对内容计算哈希
		hash := sha256.Sum256([]byte(item.Content))
		return hex.EncodeToString(hash[:])
	}
}

// calculateFilePathsHash 计算文件路径的哈希值
func calculateFilePathsHash(filePathsJSON string) string {
	if filePathsJSON == "" {
		return ""
	}

	// 解析文件路径列表
	var filePaths []string
	if err := json.Unmarshal([]byte(filePathsJSON), &filePaths); err != nil {
		// 如果解析失败，直接对原始字符串计算哈希
		hash := sha256.Sum256([]byte(filePathsJSON))
		return hex.EncodeToString(hash[:])
	}

	// 对路径列表排序以确保相同的文件集合产生相同的哈希
	sortedPaths := make([]string, len(filePaths))
	copy(sortedPaths, filePaths)
	sort.Strings(sortedPaths)

	// 将排序后的路径重新序列化为JSON
	sortedJSON, err := json.Marshal(sortedPaths)
	if err != nil {
		// 如果序列化失败，直接对原始字符串计算哈希
		hash := sha256.Sum256([]byte(filePathsJSON))
		return hex.EncodeToString(hash[:])
	}

	// 对排序后的JSON计算哈希
	hash := sha256.Sum256(sortedJSON)
	return hex.EncodeToString(hash[:])
}
