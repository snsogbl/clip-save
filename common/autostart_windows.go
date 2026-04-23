//go:build windows

package common

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// Windows 开机自启 - 使用 HKCU Run 注册表项（无需管理员权限）
const (
	autoStartRegPath = `Software\Microsoft\Windows\CurrentVersion\Run`
	autoStartAppKey  = "ClipSave" // 注册表值名，必须全局唯一
)

// SetAutoStart 设置 Windows 开机自启
// enable=true  写入注册表；enable=false 删除注册表项
func SetAutoStart(enable bool) error {
	key, _, err := registry.CreateKey(
		registry.CURRENT_USER,
		autoStartRegPath,
		registry.SET_VALUE|registry.QUERY_VALUE,
	)
	if err != nil {
		return fmt.Errorf("打开注册表失败: %w", err)
	}
	defer key.Close()

	if enable {
		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("获取可执行文件路径失败: %w", err)
		}
		// 路径加引号，避免含空格时解析错误
		cmd := fmt.Sprintf(`"%s"`, exe)
		if err := key.SetStringValue(autoStartAppKey, cmd); err != nil {
			return fmt.Errorf("写入注册表失败: %w", err)
		}
		return nil
	}

	// 关闭：删除值；不存在时视为成功
	if err := key.DeleteValue(autoStartAppKey); err != nil {
		if err == registry.ErrNotExist || strings.Contains(err.Error(), "cannot find") {
			return nil
		}
		return fmt.Errorf("删除注册表项失败: %w", err)
	}
	return nil
}

// IsAutoStartEnabled 查询是否已开启开机自启
// 返回 true 的条件：注册表存在对应键值，且值中指向的路径与当前 exe 相同
func IsAutoStartEnabled() (bool, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, autoStartRegPath, registry.QUERY_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, err
	}
	defer key.Close()

	val, _, err := key.GetStringValue(autoStartAppKey)
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, err
	}

	exe, err := os.Executable()
	if err != nil {
		// 拿不到当前 exe 时，只要值存在就算开启
		return val != "", nil
	}
	return strings.Contains(val, exe), nil
}
