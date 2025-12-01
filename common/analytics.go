package common

import (
	"fmt"
	"log"
	"os"
	gRuntime "runtime"
	"time"

	"github.com/posthog/posthog-go"
)

const (
	postHogAPIKey   = "phc_Bz1UlC85EL3cC4yVW20QpGxsybg3bqzqNPCpvwn9An0"
	postHogEndpoint = "https://us.i.posthog.com"
)

var (
	userID        string
	posthogClient posthog.Client
)

// InitAnalytics 初始化统计模块
func InitAnalytics() error {
	// 初始化 PostHog 客户端
	client, err := posthog.NewWithConfig(postHogAPIKey, posthog.Config{
		Endpoint: postHogEndpoint,
	})
	if err != nil {
		log.Printf("初始化 PostHog 客户端失败: %v", err)
		return err
	}

	posthogClient = client

	// 获取或生成用户唯一ID
	userID, err = GetOrCreateUserID()
	if err != nil {
		log.Printf("获取用户ID失败: %v", err)
		return err
	}

	// 发送应用启动事件
	go TrackEvent("app_started", nil)
	// 发送初始活跃事件
	go TrackEvent("app_active", nil)

	go activeTickerWorker()

	return nil
}

func activeTickerWorker() {
	ticker := time.NewTicker(6 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		if posthogClient == nil {
			return
		}
		TrackEvent("app_active", nil)
	}
}

// GetOrCreateUserID 获取或创建用户唯一ID
func GetOrCreateUserID() (string, error) {
	// 尝试从数据库获取用户ID
	userID, err := GetSetting("analytics_user_id")
	if err == nil && userID != "" {
		return userID, nil
	}

	// 如果不存在，生成新的UUID
	newUserID := generateUUID()
	if err := SaveSetting("analytics_user_id", newUserID); err != nil {
		return "", fmt.Errorf("保存用户ID失败: %v", err)
	}

	return newUserID, nil
}

// generateUUID 生成简单的UUID（用于统计）
func generateUUID() string {
	// 使用时间戳、系统信息和随机数生成简单的唯一ID
	timestamp := time.Now().UnixNano()
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}
	return fmt.Sprintf("%x-%s-%s-%x", timestamp, hostname, gRuntime.GOOS, time.Now().Unix())
}

// TrackEvent 追踪事件
func TrackEvent(eventName string, properties map[string]interface{}) {
	if posthogClient == nil {
		return
	}
	// 异步发送，不阻塞主流程
	go func() {
		// 添加默认属性
		if properties == nil {
			properties = make(map[string]interface{})
		}
		properties["$lib"] = "clip-save"
		properties["$lib_version"] = "2.0.7"
		properties["$os"] = gRuntime.GOOS

		// 使用 PostHog SDK 发送事件
		err := posthogClient.Enqueue(posthog.Capture{
			DistinctId: userID,
			Event:      eventName,
			Properties: properties,
		})

		if err != nil {
			log.Printf("发送统计事件失败: %v", err)
		}
	}()
}
