package common

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"golang.design/x/hotkey"
)

var (
	hk           *hotkey.Hotkey
	hotkeyCancel context.CancelFunc
	hotkeyMutex  sync.RWMutex
)

// HotkeyCallback 快捷键回调函数类型
type HotkeyCallback func()

// parseHotkeyString 解析快捷键字符串，例如 "ctrl+shift+k" 或 "cmd+alt+c" 或 "Control+v"
func parseHotkeyString(hotkeyStr string) ([]hotkey.Modifier, hotkey.Key, error) {
	parts := strings.Split(strings.ToLower(hotkeyStr), "+")
	if len(parts) < 2 {
		return nil, 0, fmt.Errorf("快捷键格式错误，至少需要两个键，例如: ctrl+shift+k")
	}

	var mods []hotkey.Modifier
	var key hotkey.Key

	// 解析修饰键
	for i := 0; i < len(parts)-1; i++ {
		part := strings.TrimSpace(parts[i])
		switch part {
		case "ctrl", "control":
			mods = append(mods, hotkey.ModCtrl)
		case "shift":
			mods = append(mods, hotkey.ModShift)
		case "alt", "option":
			mods = append(mods, hotkey.ModOption)
		case "cmd", "command", "meta":
			mods = append(mods, hotkey.ModCmd)
		default:
			return nil, 0, fmt.Errorf("不支持的修饰键: %s", part)
		}
	}

	// 解析主键
	keyStr := strings.TrimSpace(parts[len(parts)-1])
	key = parseKey(keyStr)
	if key == 0 {
		return nil, 0, fmt.Errorf("不支持的按键: %s", keyStr)
	}

	return mods, key, nil
}

// parseKey 解析按键字符串
func parseKey(keyStr string) hotkey.Key {
	switch keyStr {
	case "a":
		return hotkey.KeyA
	case "b":
		return hotkey.KeyB
	case "c":
		return hotkey.KeyC
	case "d":
		return hotkey.KeyD
	case "e":
		return hotkey.KeyE
	case "f":
		return hotkey.KeyF
	case "g":
		return hotkey.KeyG
	case "h":
		return hotkey.KeyH
	case "i":
		return hotkey.KeyI
	case "j":
		return hotkey.KeyJ
	case "k":
		return hotkey.KeyK
	case "l":
		return hotkey.KeyL
	case "m":
		return hotkey.KeyM
	case "n":
		return hotkey.KeyN
	case "o":
		return hotkey.KeyO
	case "p":
		return hotkey.KeyP
	case "q":
		return hotkey.KeyQ
	case "r":
		return hotkey.KeyR
	case "s":
		return hotkey.KeyS
	case "t":
		return hotkey.KeyT
	case "u":
		return hotkey.KeyU
	case "v":
		return hotkey.KeyV
	case "w":
		return hotkey.KeyW
	case "x":
		return hotkey.KeyX
	case "y":
		return hotkey.KeyY
	case "z":
		return hotkey.KeyZ
	case "0":
		return hotkey.Key0
	case "1":
		return hotkey.Key1
	case "2":
		return hotkey.Key2
	case "3":
		return hotkey.Key3
	case "4":
		return hotkey.Key4
	case "5":
		return hotkey.Key5
	case "6":
		return hotkey.Key6
	case "7":
		return hotkey.Key7
	case "8":
		return hotkey.Key8
	case "9":
		return hotkey.Key9
	case "space":
		return hotkey.KeySpace
	case "enter", "return":
		return hotkey.KeyReturn
	case "tab":
		return hotkey.KeyTab
	case "escape", "esc":
		return hotkey.KeyEscape
	case "delete", "del":
		return hotkey.KeyDelete
	case "up":
		return hotkey.KeyUp
	case "down":
		return hotkey.KeyDown
	case "left":
		return hotkey.KeyLeft
	case "right":
		return hotkey.KeyRight
	case "f1":
		return hotkey.KeyF1
	case "f2":
		return hotkey.KeyF2
	case "f3":
		return hotkey.KeyF3
	case "f4":
		return hotkey.KeyF4
	case "f5":
		return hotkey.KeyF5
	case "f6":
		return hotkey.KeyF6
	case "f7":
		return hotkey.KeyF7
	case "f8":
		return hotkey.KeyF8
	case "f9":
		return hotkey.KeyF9
	case "f10":
		return hotkey.KeyF10
	case "f11":
		return hotkey.KeyF11
	case "f12":
		return hotkey.KeyF12
	default:
		return 0
	}
}

// RegisterHotkey 注册全局快捷键
func RegisterHotkey(hotkeyStr string, callback HotkeyCallback) error {
	// 先取消之前注册的热键
	UnregisterHotkey()

	// 解析快捷键字符串
	mods, key, err := parseHotkeyString(hotkeyStr)
	if err != nil {
		return fmt.Errorf("解析快捷键失败: %v", err)
	}

	// 创建新的上下文用于取消
	ctx, cancel := context.WithCancel(context.Background())

	hotkeyMutex.Lock()
	hotkeyCancel = cancel
	hotkeyMutex.Unlock()

	go func() {
		err := listenHotkeyWithContext(ctx, key, mods, callback)
		if err != nil {
			log.Printf("热键监听失败: %v", err)
		}
	}()

	log.Printf("成功注册快捷键: %s", hotkeyStr)
	return nil
}

func listenHotkeyWithContext(ctx context.Context, key hotkey.Key, mods []hotkey.Modifier, callback HotkeyCallback) (err error) {
	hotkeyMutex.Lock()
	hk = hotkey.New(mods, key)
	hotkeyMutex.Unlock()

	err = hk.Register()
	if err != nil {
		return
	}
	defer func() {
		hotkeyMutex.Lock()
		if hk != nil {
			hk.Unregister()
			hk = nil
		}
		hotkeyMutex.Unlock()
	}()

	// 持续监听热键事件
	for {
		select {
		case <-ctx.Done():
			log.Println("热键监听已取消")
			return nil
		default:
			select {
			case <-ctx.Done():
				log.Println("热键监听已取消")
				return nil
			case <-hk.Keydown():
				log.Printf("热键被按下")
				select {
				case <-ctx.Done():
					log.Println("热键监听已取消")
					return nil
				case <-hk.Keyup():
					log.Printf("热键被松开")
					if callback != nil {
						callback()
					}
				}
			}
		}
	}
}

// UnregisterHotkey 取消注册快捷键
func UnregisterHotkey() {
	hotkeyMutex.Lock()
	defer hotkeyMutex.Unlock()

	if hotkeyCancel != nil {
		hotkeyCancel()
		hotkeyCancel = nil
		log.Println("✅ 快捷键已取消注册")
	}

	// 确保热键对象也被清理
	if hk != nil {
		hk.Unregister()
		hk = nil
	}
}
