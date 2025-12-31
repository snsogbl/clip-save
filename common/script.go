package common

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// UserScript 用户自定义脚本
type UserScript struct {
	ID          string
	Name        string
	Enabled     bool
	Trigger     string   // "before_save", "after_save", "on_copy", "manual"
	ContentType []string // 触发的内容类型（空数组表示所有类型）
	Keywords    []string // 关键词过滤（空数组表示不过滤）
	Script      string   // JavaScript 脚本代码
	Description string   // 脚本描述
	SortOrder   int      // 排序顺序
	PluginID    string   // 在线插件的 ID（如果是从在线插件安装的）
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ScriptEventCallback 用于发送脚本执行事件的回调函数类型
type ScriptEventCallback func(eventName string, data interface{})

// 全局脚本事件回调函数（由 app.go 设置）
var globalScriptEventCallback ScriptEventCallback

// SetScriptEventCallback 设置全局脚本事件回调函数
func SetScriptEventCallback(callback ScriptEventCallback) {
	globalScriptEventCallback = callback
}

// checkAndAddScriptTable 检查并添加脚本表
func checkAndAddScriptTable() error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 检查表是否存在
	checkSQL := `SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='user_scripts'`
	var count int
	err := DB.QueryRow(checkSQL).Scan(&count)
	if err != nil {
		return fmt.Errorf("检查脚本表失败: %v", err)
	}

	if count == 0 {
		log.Printf("🔧 正在创建 user_scripts 表...")
		createTableSQL := `
		CREATE TABLE IF NOT EXISTS user_scripts (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			enabled INTEGER DEFAULT 1,
			trigger TEXT NOT NULL,
			content_types TEXT, -- JSON 数组
			keywords TEXT, -- JSON 数组
			script TEXT NOT NULL,
			description TEXT,
			sort_order INTEGER DEFAULT 0,
			plugin_id TEXT, -- 在线插件的 ID
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_script_trigger ON user_scripts(trigger);
		CREATE INDEX IF NOT EXISTS idx_script_enabled ON user_scripts(enabled);
		CREATE INDEX IF NOT EXISTS idx_script_sort_order ON user_scripts(sort_order);
		CREATE INDEX IF NOT EXISTS idx_script_plugin_id ON user_scripts(plugin_id);
		`

		_, err := DB.Exec(createTableSQL)
		if err != nil {
			return fmt.Errorf("创建脚本表失败: %v", err)
		}
		log.Printf("✅ 已创建 user_scripts 表")
	}

	// 检查并添加 plugin_id 字段（兼容老用户）
	checkPluginIDSQL := `SELECT COUNT(*) FROM pragma_table_info('user_scripts') WHERE name = 'plugin_id'`
	var pluginIDCount int
	err = DB.QueryRow(checkPluginIDSQL).Scan(&pluginIDCount)
	if err != nil {
		log.Printf("⚠️ 检查 plugin_id 字段失败: %v", err)
	} else if pluginIDCount == 0 {
		log.Printf("🔧 正在添加 plugin_id 字段...")
		_, err = DB.Exec("ALTER TABLE user_scripts ADD COLUMN plugin_id TEXT")
		if err != nil {
			log.Printf("⚠️ 添加 plugin_id 字段失败: %v", err)
		} else {
			log.Printf("✅ 已添加 plugin_id 字段")
			// 添加索引
			_, err = DB.Exec("CREATE INDEX IF NOT EXISTS idx_script_plugin_id ON user_scripts(plugin_id)")
			if err != nil {
				log.Printf("⚠️ 创建 plugin_id 索引失败: %v", err)
			}
		}
	}

	return nil
}

// GetAllUserScripts 获取所有用户脚本
func GetAllUserScripts() ([]UserScript, error) {
	if DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	query := `SELECT id, name, enabled, trigger, content_types, keywords, 
	                 script, description, sort_order, COALESCE(plugin_id, '') as plugin_id, created_at, updated_at
	          FROM user_scripts
	          ORDER BY sort_order DESC, created_at DESC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询脚本失败: %v", err)
	}
	defer rows.Close()

	var scripts []UserScript
	for rows.Next() {
		var script UserScript
		var contentTypesJSON, keywordsJSON string

		err := rows.Scan(
			&script.ID, &script.Name, &script.Enabled, &script.Trigger,
			&contentTypesJSON, &keywordsJSON, &script.Script,
			&script.Description, &script.SortOrder, &script.PluginID, &script.CreatedAt, &script.UpdatedAt,
		)
		if err != nil {
			log.Printf("扫描脚本行失败: %v", err)
			continue
		}

		// 解析 JSON 数组
		if contentTypesJSON != "" {
			json.Unmarshal([]byte(contentTypesJSON), &script.ContentType)
		}
		if keywordsJSON != "" {
			json.Unmarshal([]byte(keywordsJSON), &script.Keywords)
		}

		scripts = append(scripts, script)
	}

	return scripts, nil
}

// GetEnabledUserScripts 获取启用的脚本
func GetEnabledUserScripts(trigger string) ([]UserScript, error) {
	if DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	query := `SELECT id, name, enabled, trigger, content_types, keywords, 
	                 script, description, sort_order, COALESCE(plugin_id, '') as plugin_id, created_at, updated_at
	          FROM user_scripts
	          WHERE enabled = 1 AND trigger = ?
	          ORDER BY sort_order DESC, created_at DESC`

	rows, err := DB.Query(query, trigger)
	if err != nil {
		return nil, fmt.Errorf("查询脚本失败: %v", err)
	}
	defer rows.Close()

	var scripts []UserScript
	for rows.Next() {
		var script UserScript
		var contentTypesJSON, keywordsJSON string

		err := rows.Scan(
			&script.ID, &script.Name, &script.Enabled, &script.Trigger,
			&contentTypesJSON, &keywordsJSON, &script.Script,
			&script.Description, &script.SortOrder, &script.PluginID, &script.CreatedAt, &script.UpdatedAt,
		)
		if err != nil {
			continue
		}

		if contentTypesJSON != "" {
			json.Unmarshal([]byte(contentTypesJSON), &script.ContentType)
		}
		if keywordsJSON != "" {
			json.Unmarshal([]byte(keywordsJSON), &script.Keywords)
		}

		scripts = append(scripts, script)
	}

	return scripts, nil
}

// SaveUserScript 保存用户脚本
func SaveUserScript(script *UserScript) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	// 如果没有 ID，生成一个（新脚本）
	isNewScript := script.ID == ""
	if isNewScript {
		script.ID = fmt.Sprintf("%d", time.Now().UnixNano())
		// 新脚本放在最前面：查询当前最大的 sort_order，然后设置为最大值 + 1
		if script.SortOrder == 0 {
			var maxSortOrder sql.NullInt64
			err := DB.QueryRow("SELECT MAX(sort_order) FROM user_scripts").Scan(&maxSortOrder)
			if err != nil {
				// 如果查询失败或没有记录，设置为 1
				script.SortOrder = 1
			} else if maxSortOrder.Valid {
				// 设置为最大值 + 1，确保新脚本在最前面（排序值大的在前）
				script.SortOrder = int(maxSortOrder.Int64) + 1
			} else {
				// 如果没有现有脚本，设置为 1
				script.SortOrder = 1
			}
		}
	}

	contentTypesJSON, _ := json.Marshal(script.ContentType)
	keywordsJSON, _ := json.Marshal(script.Keywords)

	enabled := 0
	if script.Enabled {
		enabled = 1
	}

	insertSQL := `
	INSERT INTO user_scripts 
	(id, name, enabled, trigger, content_types, keywords, script, description, sort_order, plugin_id, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))
	ON CONFLICT(id) DO UPDATE SET
		name = excluded.name,
		enabled = excluded.enabled,
		trigger = excluded.trigger,
		content_types = excluded.content_types,
		keywords = excluded.keywords,
		script = excluded.script,
		description = excluded.description,
		sort_order = excluded.sort_order,
		plugin_id = excluded.plugin_id,
		updated_at = datetime('now')
	`

	_, err := DB.Exec(insertSQL,
		script.ID, script.Name, enabled, script.Trigger,
		string(contentTypesJSON), string(keywordsJSON),
		script.Script, script.Description, script.SortOrder, script.PluginID,
	)

	if err != nil {
		return fmt.Errorf("保存脚本失败: %v", err)
	}

	log.Printf("✅ 已保存脚本: %s", script.Name)
	return nil
}

// DeleteUserScript 删除用户脚本
func DeleteUserScript(id string) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	deleteSQL := `DELETE FROM user_scripts WHERE id = ?`
	result, err := DB.Exec(deleteSQL, id)
	if err != nil {
		return fmt.Errorf("删除脚本失败: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("未找到要删除的脚本")
	}

	log.Printf("✅ 已删除脚本: %s", id)
	return nil
}

// GetUserScriptByID 根据 ID 获取脚本
func GetUserScriptByID(id string) (*UserScript, error) {
	if DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	query := `SELECT id, name, enabled, trigger, content_types, keywords, 
	                 script, description, sort_order, COALESCE(plugin_id, '') as plugin_id, created_at, updated_at
	          FROM user_scripts WHERE id = ?`

	var script UserScript
	var contentTypesJSON, keywordsJSON string

	err := DB.QueryRow(query, id).Scan(
		&script.ID, &script.Name, &script.Enabled, &script.Trigger,
		&contentTypesJSON, &keywordsJSON, &script.Script,
		&script.Description, &script.SortOrder, &script.PluginID, &script.CreatedAt, &script.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("未找到脚本")
	}
	if err != nil {
		return nil, fmt.Errorf("查询脚本失败: %v", err)
	}

	if contentTypesJSON != "" {
		json.Unmarshal([]byte(contentTypesJSON), &script.ContentType)
	}
	if keywordsJSON != "" {
		json.Unmarshal([]byte(keywordsJSON), &script.Keywords)
	}

	return &script, nil
}

// GetUserScriptsByIDs 根据 ID 列表批量获取脚本
func GetUserScriptsByIDs(ids []string) ([]UserScript, error) {
	if DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}

	if len(ids) == 0 {
		return []UserScript{}, nil
	}

	// 构建 IN 查询的占位符
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT id, name, enabled, trigger, content_types, keywords, 
	                            script, description, sort_order, COALESCE(plugin_id, '') as plugin_id, created_at, updated_at
	                     FROM user_scripts WHERE id IN (%s)`,
		strings.Join(placeholders, ","))

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("批量查询脚本失败: %v", err)
	}
	defer rows.Close()

	var scripts []UserScript
	for rows.Next() {
		var script UserScript
		var contentTypesJSON, keywordsJSON string

		err := rows.Scan(
			&script.ID, &script.Name, &script.Enabled, &script.Trigger,
			&contentTypesJSON, &keywordsJSON, &script.Script,
			&script.Description, &script.SortOrder, &script.PluginID, &script.CreatedAt, &script.UpdatedAt,
		)
		if err != nil {
			log.Printf("扫描脚本行失败: %v", err)
			continue
		}

		if contentTypesJSON != "" {
			json.Unmarshal([]byte(contentTypesJSON), &script.ContentType)
		}
		if keywordsJSON != "" {
			json.Unmarshal([]byte(keywordsJSON), &script.Keywords)
		}

		scripts = append(scripts, script)
	}

	return scripts, nil
}

// UpdateUserScriptOrder 更新单个脚本顺序
func UpdateUserScriptOrder(scriptID string, sortOrder int) error {
	if DB == nil {
		return fmt.Errorf("数据库未初始化")
	}

	updateSQL := `UPDATE user_scripts SET sort_order = ?, updated_at = datetime('now') WHERE id = ?`
	_, err := DB.Exec(updateSQL, sortOrder, scriptID)
	if err != nil {
		return fmt.Errorf("更新脚本 %s 顺序失败: %v", scriptID, err)
	}

	return nil
}
