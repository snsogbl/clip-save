package common

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB 初始化数据库
func InitDB() error {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户目录失败: %v", err)
	}

	// 创建应用数据目录
	appDir := filepath.Join(homeDir, ".clipsave")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("创建应用目录失败: %v", err)
	}

	// 数据库文件路径
	dbPath := filepath.Join(appDir, "clipboard.db")
	log.Printf("数据库路径: %s", dbPath)

	// 打开数据库连接
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("打开数据库失败: %v", err)
	}

	DB = db

	// 创建表
	if err := createTables(); err != nil {
		return fmt.Errorf("创建表失败: %v", err)
	}

	// 检查并添加新字段（兼容老用户）
	if err := checkAndAddNewFields(); err != nil {
		return fmt.Errorf("添加新字段失败: %v", err)
	}

	// 初始化默认设置
	if err := initDefaultSettings(); err != nil {
		log.Printf("警告: 初始化默认设置失败: %v", err)
		// 不返回错误，允许应用继续运行
	}

	// 检查是否是第一次创建表（表为空时）
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM clipboard_items").Scan(&count)
	if err != nil {
		log.Printf("检查剪贴板记录失败: %v", err)
	} else if count == 0 {
		// 第一次创建表，添加默认文本记录
		if err := initDefaultTextRecord(); err != nil {
			log.Printf("警告: 初始化默认文本记录失败: %v", err)
			// 不返回错误，允许应用继续运行
		}
	}

	log.Println("数据库初始化成功")
	return nil
}

// createTables 创建数据表
func createTables() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS clipboard_items (
		id TEXT PRIMARY KEY,
		content TEXT NOT NULL,
		content_type TEXT NOT NULL,
		image_data BLOB,
		file_paths TEXT,
		file_info TEXT,
		timestamp DATETIME NOT NULL,
		source TEXT,
		char_count INTEGER,
		word_count INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_timestamp ON clipboard_items(timestamp DESC);
	CREATE INDEX IF NOT EXISTS idx_content_type ON clipboard_items(content_type);

	CREATE TABLE IF NOT EXISTS app_settings (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := DB.Exec(createTableSQL)
	return err
}

// SaveClipboardItem 保存剪贴板项目（支持去重）
func SaveClipboardItem(item *ClipboardItem) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 检查是否存在相同内容的项目
	if item.ContentHash != "" {
		var existingID string
		checkSQL := `SELECT id FROM clipboard_items WHERE content_hash = ? AND content_type = ? LIMIT 1`
		err := DB.QueryRow(checkSQL, item.ContentHash, item.ContentType).Scan(&existingID)

		if err == nil {
			// 找到重复项，先删除旧记录
			deleteSQL := `DELETE FROM clipboard_items WHERE id = ?`
			_, deleteErr := DB.Exec(deleteSQL, existingID)
			if deleteErr != nil {
				log.Printf("⚠️ 删除重复项目失败: %v", deleteErr)
			} else {
				log.Printf("🔄 删除重复项目: ID=%s", existingID)
			}
		}
	}

	// 插入新记录
	insertSQL := `
	INSERT INTO clipboard_items (id, content, content_type, content_hash, image_data, file_paths, file_info, timestamp, source, char_count, word_count)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := DB.Exec(insertSQL,
		item.ID,
		item.Content,
		item.ContentType,
		item.ContentHash,
		item.ImageData,
		item.FilePaths,
		item.FileInfo,
		item.Timestamp,
		item.Source,
		item.CharCount,
		item.WordCount,
	)

	if err != nil {
		return fmt.Errorf("保存剪贴板项目失败: %v", err)
	}

	log.Printf("已保存剪贴板项目: ID=%s, 类型=%s, 哈希=%s", item.ID, item.ContentType, item.ContentHash[:8])
	return nil
}

// GetClipboardItems 获取剪贴板项目列表
func GetClipboardItems(limit int) ([]ClipboardItem, error) {
	if DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	query := `
	SELECT id, content, content_type, content_hash, image_data, file_paths, file_info, timestamp, source, char_count, word_count
	FROM clipboard_items
	ORDER BY timestamp DESC
	LIMIT ?
	`

	rows, err := DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("查询剪贴板项目失败: %v", err)
	}
	defer rows.Close()

	var items []ClipboardItem
	for rows.Next() {
		var item ClipboardItem
		err := rows.Scan(
			&item.ID,
			&item.Content,
			&item.ContentType,
			&item.ContentHash,
			&item.ImageData,
			&item.FilePaths,
			&item.FileInfo,
			&item.Timestamp,
			&item.Source,
			&item.CharCount,
			&item.WordCount,
		)
		if err != nil {
			log.Printf("扫描行失败: %v", err)
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

// GetClipboardItemByID 根据ID获取剪贴板项目
func GetClipboardItemByID(id string) (*ClipboardItem, error) {
	if DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	query := `
	SELECT id, content, content_type, content_hash, image_data, file_paths, file_info, timestamp, source, char_count, word_count
	FROM clipboard_items
	WHERE id = ?
	`

	var item ClipboardItem
	err := DB.QueryRow(query, id).Scan(
		&item.ID,
		&item.Content,
		&item.ContentType,
		&item.ContentHash,
		&item.ImageData,
		&item.FilePaths,
		&item.FileInfo,
		&item.Timestamp,
		&item.Source,
		&item.CharCount,
		&item.WordCount,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("未找到剪贴板项目")
	}
	if err != nil {
		return nil, fmt.Errorf("查询剪贴板项目失败: %v", err)
	}

	return &item, nil
}

// DeleteClipboardItem 删除剪贴板项目
func DeleteClipboardItem(id string) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	deleteSQL := `DELETE FROM clipboard_items WHERE id = ?`
	result, err := DB.Exec(deleteSQL, id)
	if err != nil {
		return fmt.Errorf("删除剪贴板项目失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("未找到要删除的项目")
	}

	log.Printf("已删除剪贴板项目: ID=%s", id)
	return nil
}

// ClearOldItems 清除旧的剪贴板项目（保留最近N条）
func ClearOldItems(keepCount int) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	deleteSQL := `
	DELETE FROM clipboard_items
	WHERE id NOT IN (
		SELECT id FROM clipboard_items
		ORDER BY timestamp DESC
		LIMIT ?
	)
	`

	result, err := DB.Exec(deleteSQL, keepCount)
	if err != nil {
		return fmt.Errorf("清除旧项目失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("已清除 %d 条旧的剪贴板项目", rowsAffected)
	return nil
}

// ClearItemsOlderThanDays 清除超过指定天数的剪贴板项目
func ClearItemsOlderThanDays(days int) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 计算截止日期
	cutoffDate := time.Now().AddDate(0, 0, -days)

	deleteSQL := `
	DELETE FROM clipboard_items
	WHERE timestamp < ?
	`

	result, err := DB.Exec(deleteSQL, cutoffDate.Format("2006-01-02 15:04:05"))
	if err != nil {
		return fmt.Errorf("清除超过 %d 天的项目失败: %v", days, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		log.Printf("已清除 %d 条超过 %d 天的剪贴板项目", rowsAffected, days)
	}
	return nil
}

// ClearAllItems 清除所有剪贴板项目
func ClearAllItems() error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	deleteSQL := `DELETE FROM clipboard_items`

	result, err := DB.Exec(deleteSQL)
	if err != nil {
		return fmt.Errorf("清除所有项目失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("已清除所有剪贴板项目，共 %d 条", rowsAffected)
	return nil
}

// SearchClipboardItems 搜索剪贴板项目
func SearchClipboardItems(keyword string, filterType string, limit int) ([]ClipboardItem, error) {
	if DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	query := `
	SELECT id, content, content_type, content_hash, image_data, file_paths, file_info, timestamp, source, char_count, word_count
	FROM clipboard_items
	WHERE 1=1
	`
	args := []interface{}{}

	// 关键词搜索（不区分大小写）
	if keyword != "" {
		query += ` AND (content LIKE ? COLLATE NOCASE)`
		args = append(args, "%"+keyword+"%")
	}

	// 类型过滤（支持中文）
	if filterType != "" && filterType != "所有类型" {
		var dbContentType string
		switch filterType {
		case "文本":
			dbContentType = "Text"
		case "图片":
			dbContentType = "Image"
		case "文件":
			dbContentType = "File"
		case "URL":
			dbContentType = "URL"
		case "颜色":
			dbContentType = "Color"
		default:
			dbContentType = filterType // 如果是英文直接使用
		}
		query += ` AND content_type = ?`
		args = append(args, dbContentType)
	}

	query += ` ORDER BY timestamp DESC LIMIT ?`
	args = append(args, limit)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("搜索剪贴板项目失败: %v", err)
	}
	defer rows.Close()

	var items []ClipboardItem
	for rows.Next() {
		var item ClipboardItem
		err := rows.Scan(
			&item.ID,
			&item.Content,
			&item.ContentType,
			&item.ContentHash,
			&item.ImageData,
			&item.FilePaths,
			&item.FileInfo,
			&item.Timestamp,
			&item.Source,
			&item.CharCount,
			&item.WordCount,
		)
		if err != nil {
			log.Printf("扫描行失败: %v", err)
			continue
		}
		items = append(items, item)
	}

	return items, nil
}

// GetStatistics 获取统计信息
func GetStatistics() (map[string]interface{}, error) {
	if DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	stats := make(map[string]interface{})

	// 总数量
	var totalCount int
	err := DB.QueryRow("SELECT COUNT(*) FROM clipboard_items").Scan(&totalCount)
	if err != nil {
		return nil, err
	}
	stats["total_count"] = totalCount

	// 今天的数量
	today := time.Now().Format("2006-01-02")
	var todayCount int
	err = DB.QueryRow("SELECT COUNT(*) FROM clipboard_items WHERE DATE(timestamp) = ?", today).Scan(&todayCount)
	if err != nil {
		return nil, err
	}
	stats["today_count"] = todayCount

	// 按类型统计
	rows, err := DB.Query("SELECT content_type, COUNT(*) FROM clipboard_items GROUP BY content_type")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	typeStats := make(map[string]int)
	for rows.Next() {
		var contentType string
		var count int
		if err := rows.Scan(&contentType, &count); err == nil {
			typeStats[contentType] = count
		}
	}
	stats["type_stats"] = typeStats

	return stats, nil
}

// initDefaultSettings 初始化默认设置
func initDefaultSettings() error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 检查设置是否已存在
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM app_settings WHERE key = 'app_settings'").Scan(&count)
	if err != nil {
		return fmt.Errorf("检查设置失败: %v", err)
	}

	// 如果设置已存在，不覆盖
	if count > 0 {
		log.Println("设置已存在，跳过初始化")
		return nil
	}

	// 创建默认设置 JSON（密码默认为空，表示不需要密码，快捷键默认为 Command+Option++c）
	defaultSettings := `{"autoClean":true,"retentionDays":30,"pageSize":100,"password":"","hotkey":"Command+Option+c"}`

	insertSQL := `
	INSERT INTO app_settings (key, value, updated_at) 
	VALUES ('app_settings', ?, datetime('now'))
	`

	_, err = DB.Exec(insertSQL, defaultSettings)
	if err != nil {
		return fmt.Errorf("初始化默认设置失败: %v", err)
	}

	log.Printf("✅ 已初始化默认设置: %s", defaultSettings)
	return nil
}

// initDefaultTextRecord 初始化默认文本记录
func initDefaultTextRecord() error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 创建默认文本记录
	defaultText := "剪存：Command+Option+c 唤之即来"
	timestamp := time.Now()
	item := ClipboardItem{
		ID:          fmt.Sprintf("%d", timestamp.UnixNano()),
		Content:     defaultText,
		ContentType: "Text",
		Timestamp:   timestamp,
		Source:      "系统初始化",
		CharCount:   len([]rune(defaultText)),
		WordCount:   countWords(defaultText),
	}

	// 计算内容哈希
	item.ContentHash = calculateContentHash(&item)

	// 保存到数据库
	if err := SaveClipboardItem(&item); err != nil {
		return fmt.Errorf("保存默认文本记录失败: %v", err)
	}

	log.Printf("✅ 已初始化默认文本记录: %s", defaultText)
	return nil
}

// SaveSetting 保存单个设置项
func SaveSetting(key string, value string) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	insertSQL := `
	INSERT INTO app_settings (key, value, updated_at) 
	VALUES (?, ?, datetime('now'))
	ON CONFLICT(key) DO UPDATE SET 
		value = excluded.value,
		updated_at = datetime('now')
	`

	_, err := DB.Exec(insertSQL, key, value)
	if err != nil {
		return fmt.Errorf("保存设置失败: %v", err)
	}

	log.Printf("已保存设置: %s", key)
	return nil
}

// GetSetting 获取单个设置项
func GetSetting(key string) (string, error) {
	if DB == nil {
		return "", fmt.Errorf("数据库未初始化")
	}

	query := `SELECT value FROM app_settings WHERE key = ?`
	var value string
	err := DB.QueryRow(query, key).Scan(&value)

	if err == sql.ErrNoRows {
		return "", nil // 设置不存在返回空字符串
	}
	if err != nil {
		return "", fmt.Errorf("获取设置失败: %v", err)
	}

	return value, nil
}

// GetAllSettings 获取所有设置
func GetAllSettings() (map[string]string, error) {
	if DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	query := `SELECT key, value FROM app_settings`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询设置失败: %v", err)
	}
	defer rows.Close()

	settings := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			log.Printf("扫描设置行失败: %v", err)
			continue
		}
		settings[key] = value
	}

	return settings, nil
}

// migrateContentHash 为现有数据添加哈希值
func migrateContentHash() error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 查找所有没有哈希值的记录
	query := `
	SELECT id, content, content_type, image_data, file_paths
	FROM clipboard_items 
	WHERE content_hash IS NULL OR content_hash = ''
	`

	rows, err := DB.Query(query)
	if err != nil {
		return fmt.Errorf("查询需要迁移的记录失败: %v", err)
	}
	defer rows.Close()

	var migratedCount int
	updateSQL := `UPDATE clipboard_items SET content_hash = ? WHERE id = ?`

	for rows.Next() {
		var item ClipboardItem
		err := rows.Scan(
			&item.ID,
			&item.Content,
			&item.ContentType,
			&item.ImageData,
			&item.FilePaths,
		)
		if err != nil {
			log.Printf("扫描迁移记录失败: %v", err)
			continue
		}

		// 计算哈希值
		contentHash := calculateContentHash(&item)
		if contentHash == "" {
			log.Printf("计算哈希失败，跳过记录: %s", item.ID)
			continue
		}

		// 更新数据库
		_, err = DB.Exec(updateSQL, contentHash, item.ID)
		if err != nil {
			log.Printf("更新记录哈希失败: %v (ID: %s)", err, item.ID)
			continue
		}

		migratedCount++
	}

	if migratedCount > 0 {
		log.Printf("✅ 成功为 %d 条现有记录添加了哈希值", migratedCount)
	} else {
		log.Printf("✅ 所有记录都已有哈希值，无需迁移")
	}

	return nil
}

// checkAndAddNewFields 检查并添加新字段（兼容老用户）
func checkAndAddNewFields() error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 检查content_hash字段是否存在
	checkSQL := `SELECT COUNT(*) FROM pragma_table_info('clipboard_items') WHERE name = 'content_hash'`
	var count int
	err := DB.QueryRow(checkSQL).Scan(&count)
	if err != nil {
		return fmt.Errorf("检查content_hash字段是否存在失败: %v", err)
	}

	if count == 0 {
		// 字段不存在，需要添加
		log.Printf("🔧 检测到老版本数据库，正在添加content_hash字段...")

		// 添加content_hash字段
		alterSQL := `ALTER TABLE clipboard_items ADD COLUMN content_hash TEXT`
		_, err := DB.Exec(alterSQL)
		if err != nil {
			return fmt.Errorf("添加content_hash字段失败: %v", err)
		}
		log.Printf("✅ 已添加content_hash字段")

		// 创建索引
		indexSQL := `CREATE INDEX IF NOT EXISTS idx_content_hash ON clipboard_items(content_hash, content_type)`
		_, err = DB.Exec(indexSQL)
		if err != nil {
			return fmt.Errorf("创建content_hash索引失败: %v", err)
		}
		log.Printf("✅ 已创建content_hash索引")

		// 为现有数据添加哈希值（只在字段刚添加时进行）
		if err := migrateContentHash(); err != nil {
			log.Printf("⚠️ 警告: 为现有数据添加哈希值失败: %v", err)
			// 不返回错误，允许应用继续运行
		}
	} else {
		log.Printf("✅ content_hash字段已存在")
		// 字段已存在，但检查是否有未设置哈希值的记录（处理意外情况）
		if err := migrateContentHash(); err != nil {
			log.Printf("⚠️ 警告: 检查并更新哈希值失败: %v", err)
			// 不返回错误，允许应用继续运行
		}
	}

	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
