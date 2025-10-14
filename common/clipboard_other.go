//go:build !darwin
// +build !darwin

package common

import "fmt"

// GetFrontmostAppName 获取当前活动应用程序的名称（非 macOS 平台返回 System）
func GetFrontmostAppName() string {
	return "System"
}

// ReadPasteboardData 读取指定类型的剪贴板数据（非 macOS 平台不支持）
func ReadPasteboardData(typeName string) []byte {
	return nil
}

// ReadFileURLs 读取剪贴板中的文件 URL 列表（非 macOS 平台不支持）
func ReadFileURLs() (string, int) {
	return "", 0
}

// WriteFileURLs 将文件 URL 写入剪贴板（非 macOS 平台不支持）
func WriteFileURLs(filePaths string) error {
	return fmt.Errorf("不支持的平台")
}

// GetPasteboardChangeCount 获取剪贴板的变化计数（非 macOS 平台不支持）
func GetPasteboardChangeCount() int {
	return 0
}
