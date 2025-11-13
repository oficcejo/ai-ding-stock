package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// StockConfig 股票分析系统配置
type StockConfig struct {
	TDXAPIUrl     string             `json:"tdx_api_url"`
	AIConfig      AIConfig           `json:"ai_config"`
	Stocks        []StockItem        `json:"stocks"`
	Notification  NotificationConfig `json:"notification"`
	TradingTime   TradingTimeConfig  `json:"trading_time"`
	APIServerPort int                `json:"api_server_port"`
	LogDir        string             `json:"log_dir"`
}

// TradingTimeConfig 交易时间配置
type TradingTimeConfig struct {
	EnableCheck  bool     `json:"enable_check"`  // 是否启用交易时间检查
	TradingHours []string `json:"trading_hours"` // 交易时段（如：["09:30-11:30", "13:00-15:00"]）
	Timezone     string   `json:"timezone"`      // 时区（如：Asia/Shanghai）
}

// AIConfig AI配置
type AIConfig struct {
	Provider        string `json:"provider"` // "deepseek", "qwen", "custom"
	DeepSeekKey     string `json:"deepseek_key"`
	QwenKey         string `json:"qwen_key"`
	CustomAPIURL    string `json:"custom_api_url"`
	CustomAPIKey    string `json:"custom_api_key"`
	CustomModelName string `json:"custom_model_name"`
}

// StockItem 股票配置项
type StockItem struct {
	Code                string `json:"code"`
	Name                string `json:"name"`
	Enabled             bool   `json:"enabled"`
	ScanIntervalMinutes int    `json:"scan_interval_minutes"`
	MinConfidence       int    `json:"min_confidence"` // 最小信心度阈值
}

// NotificationConfig 通知配置
type NotificationConfig struct {
	Enabled  bool           `json:"enabled"`
	DingTalk DingTalkConfig `json:"dingtalk"`
	Feishu   FeishuConfig   `json:"feishu"`
}

// DingTalkConfig 钉钉配置
type DingTalkConfig struct {
	Enabled    bool   `json:"enabled"`
	WebhookURL string `json:"webhook_url"`
	Secret     string `json:"secret"`
}

// FeishuConfig 飞书配置
type FeishuConfig struct {
	Enabled    bool   `json:"enabled"`
	WebhookURL string `json:"webhook_url"`
	Secret     string `json:"secret"`
}

// LoadStockConfig 加载股票分析配置
func LoadStockConfig(filename string) (*StockConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config StockConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return &config, nil
}

// Validate 验证配置
func (c *StockConfig) Validate() error {
	// 验证TDX API URL
	if c.TDXAPIUrl == "" {
		return fmt.Errorf("tdx_api_url不能为空")
	}

	// 验证AI配置
	if c.AIConfig.Provider == "" {
		return fmt.Errorf("ai_config.provider不能为空")
	}
	if c.AIConfig.Provider != "deepseek" && c.AIConfig.Provider != "qwen" && c.AIConfig.Provider != "custom" {
		return fmt.Errorf("ai_config.provider必须是 'deepseek', 'qwen' 或 'custom'")
	}

	// 验证对应的API密钥
	if c.AIConfig.Provider == "deepseek" && c.AIConfig.DeepSeekKey == "" {
		return fmt.Errorf("使用DeepSeek时必须配置deepseek_key")
	}
	if c.AIConfig.Provider == "qwen" && c.AIConfig.QwenKey == "" {
		return fmt.Errorf("使用Qwen时必须配置qwen_key")
	}
	if c.AIConfig.Provider == "custom" {
		if c.AIConfig.CustomAPIURL == "" || c.AIConfig.CustomAPIKey == "" || c.AIConfig.CustomModelName == "" {
			return fmt.Errorf("使用自定义API时必须配置custom_api_url, custom_api_key和custom_model_name")
		}
	}

	// 验证股票列表
	if len(c.Stocks) == 0 {
		return fmt.Errorf("至少需要配置一只股票")
	}

	stockCodes := make(map[string]bool)
	enabledCount := 0
	for i, stock := range c.Stocks {
		if stock.Code == "" {
			return fmt.Errorf("stocks[%d]: code不能为空", i)
		}
		if stock.Name == "" {
			return fmt.Errorf("stocks[%d]: name不能为空", i)
		}
		if stockCodes[stock.Code] {
			return fmt.Errorf("stocks[%d]: 股票代码 '%s' 重复", i, stock.Code)
		}
		stockCodes[stock.Code] = true

		if stock.Enabled {
			enabledCount++
		}

		// 设置默认值
		if stock.ScanIntervalMinutes <= 0 {
			c.Stocks[i].ScanIntervalMinutes = 5 // 默认5分钟
		}
		if stock.MinConfidence <= 0 {
			c.Stocks[i].MinConfidence = 70 // 默认70%信心度
		}
	}

	if enabledCount == 0 {
		return fmt.Errorf("至少需要启用一只股票")
	}

	// 设置默认API端口
	if c.APIServerPort <= 0 {
		c.APIServerPort = 9090
	}

	// 设置默认日志目录
	if c.LogDir == "" {
		c.LogDir = "stock_analysis_logs"
	}

	// 设置默认交易时间配置
	if c.TradingTime.Timezone == "" {
		c.TradingTime.Timezone = "Asia/Shanghai"
	}
	if len(c.TradingTime.TradingHours) == 0 {
		c.TradingTime.TradingHours = []string{"09:30-11:30", "13:00-15:00"} // A股默认交易时段
	}

	// 验证通知配置
	if c.Notification.Enabled {
		if !c.Notification.DingTalk.Enabled && !c.Notification.Feishu.Enabled {
			return fmt.Errorf("启用通知时至少需要配置一个通知渠道（钉钉或飞书）")
		}
		if c.Notification.DingTalk.Enabled && c.Notification.DingTalk.WebhookURL == "" {
			return fmt.Errorf("启用钉钉通知时必须配置webhook_url")
		}
		if c.Notification.Feishu.Enabled && c.Notification.Feishu.WebhookURL == "" {
			return fmt.Errorf("启用飞书通知时必须配置webhook_url")
		}
	}

	return nil
}

// GetScanInterval 获取扫描间隔
func (s *StockItem) GetScanInterval() time.Duration {
	return time.Duration(s.ScanIntervalMinutes) * time.Minute
}
