package common

import (
	"encoding/json"
	"fmt"
)

const AppVersion = "3"

// 支持的语言
const (
	LangChinese = "zh-CN"
	LangEnglish = "en-US"
	LangFrench  = "fr-FR"
	LangArabic  = "ar-SA"
)

// 默认语言
const DefaultLanguage = LangChinese

// 翻译键值对
type Translations map[string]string

// 语言包
type LanguagePack struct {
	Language     string       `json:"language"`
	Translations Translations `json:"translations"`
}

// 全局语言包存储
var languagePacks = make(map[string]LanguagePack)
var currentLanguage = DefaultLanguage

// 初始化国际化
func InitI18n() error {
	// 加载中文语言包
	zhPack := LanguagePack{
		Language: LangChinese,
		Translations: Translations{
			// 应用信息
			"app.title":       "剪存 - 剪贴板历史",
			"app.name":        "剪存",
			"app.description": "剪贴板历史管理工具",
			"app.version":     AppVersion,
		},
	}

	// 加载英文语言包
	enPack := LanguagePack{
		Language: LangEnglish,
		Translations: Translations{
			// 应用信息
			"app.title":       "ClipSave - Clipboard History",
			"app.name":        "ClipSave",
			"app.description": "Clipboard History Management Tool",
			"app.version":     AppVersion,
		},
	}

	// 加载法文语言包
	frPack := LanguagePack{
		Language: LangFrench,
		Translations: Translations{
			// 应用信息
			"app.title":       "ClipSave - Historique du Presse-papiers",
			"app.name":        "ClipSave",
			"app.description": "Outil de Gestion de l'Historique du Presse-papiers",
			"app.version":     AppVersion,
		},
	}

	// 加载阿拉伯语语言包
	arPack := LanguagePack{
		Language: LangArabic,
		Translations: Translations{
			// 应用信息
			"app.title":       "ClipSave - سجل الحافظة",
			"app.name":        "ClipSave",
			"app.description": "أداة إدارة سجل الحافظة",
			"app.version":     AppVersion,
		},
	}

	// 存储语言包
	languagePacks[LangChinese] = zhPack
	languagePacks[LangEnglish] = enPack
	languagePacks[LangFrench] = frPack
	languagePacks[LangArabic] = arPack

	// 尝试从数据库加载用户的语言设置
	settingsJSON, err := GetSetting("app_settings")
	if err == nil && settingsJSON != "" {
		var settings map[string]interface{}
		if json.Unmarshal([]byte(settingsJSON), &settings) == nil {
			if lang, ok := settings["language"].(string); ok && (lang == LangChinese || lang == LangEnglish || lang == LangFrench || lang == LangArabic) {
				currentLanguage = lang
			}
		}
	}

	return nil
}

// 设置当前语言
func SetLanguage(lang string) error {
	if _, exists := languagePacks[lang]; !exists {
		return fmt.Errorf("unsupported language: %s", lang)
	}
	currentLanguage = lang

	// 保存语言设置到数据库
	settingsJSON, err := GetSetting("app_settings")
	var settings map[string]interface{}
	if err != nil || settingsJSON == "" {
		settings = make(map[string]interface{})
	} else {
		json.Unmarshal([]byte(settingsJSON), &settings)
	}
	settings["language"] = lang
	newSettingsJSON, _ := json.Marshal(settings)
	SaveSetting("app_settings", string(newSettingsJSON))

	return nil
}

// 获取当前语言
func GetCurrentLanguage() string {
	return currentLanguage
}

// 获取翻译文本
func T(key string, args ...interface{}) string {
	pack, exists := languagePacks[currentLanguage]
	if !exists {
		pack = languagePacks[LangChinese] // 回退到中文
	}

	text, exists := pack.Translations[key]
	if !exists {
		return key // 如果找不到翻译，返回键名
	}

	if len(args) > 0 {
		return fmt.Sprintf(text, args...)
	}
	return text
}

// 获取支持的语言列表
func GetSupportedLanguages() []string {
	languages := make([]string, 0, len(languagePacks))
	for lang := range languagePacks {
		languages = append(languages, lang)
	}
	return languages
}
