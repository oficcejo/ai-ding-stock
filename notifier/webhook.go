package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Notifier 通知器接口
type Notifier interface {
	SendSignal(signal *TradingSignal) error
	SendMessage(message string) error
}

// TradingSignal 交易信号
type TradingSignal struct {
	StockCode     string                 `json:"stock_code"`               // 股票代码
	StockName     string                 `json:"stock_name"`               // 股票名称
	Signal        string                 `json:"signal"`                   // 信号类型: BUY/SELL/HOLD
	Price         float64                `json:"price"`                    // 当前价格
	Confidence    int                    `json:"confidence"`               // 信心度 (0-100)
	Reasoning     string                 `json:"reasoning"`                // 推理原因
	TargetPrice   float64                `json:"target_price"`             // 目标价格
	StopLoss      float64                `json:"stop_loss"`                // 止损价格
	RiskReward    string                 `json:"risk_reward"`              // 风险回报比
	Timestamp     time.Time              `json:"timestamp"`                // 时间戳
	TechnicalData map[string]interface{} `json:"technical_data,omitempty"` // 技术指标数据
}

// DingTalkNotifier 钉钉通知器
type DingTalkNotifier struct {
	WebhookURL string
	Secret     string // 加签密钥（可选）
}

// NewDingTalkNotifier 创建钉钉通知器
func NewDingTalkNotifier(webhookURL string, secret string) *DingTalkNotifier {
	return &DingTalkNotifier{
		WebhookURL: webhookURL,
		Secret:     secret,
	}
}

// SendSignal 发送交易信号到钉钉
func (d *DingTalkNotifier) SendSignal(signal *TradingSignal) error {
	// 构建Markdown格式的消息
	markdown := d.formatSignalMarkdown(signal)

	// 钉钉消息格式
	message := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": fmt.Sprintf("【%s】%s %s", signal.Signal, signal.StockName, signal.StockCode),
			"text":  markdown,
		},
		"at": map[string]interface{}{
			"isAtAll": false,
		},
	}

	return d.sendRequest(message)
}

// SendMessage 发送普通消息到钉钉
func (d *DingTalkNotifier) SendMessage(message string) error {
	msg := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": message,
		},
	}
	return d.sendRequest(msg)
}

// formatSignalMarkdown 格式化信号为Markdown
func (d *DingTalkNotifier) formatSignalMarkdown(signal *TradingSignal) string {
	var emoji string
	switch signal.Signal {
	case "BUY":
		emoji = "🚀"
	case "SELL":
		emoji = "⚠️"
	case "HOLD":
		emoji = "⏸️"
	default:
		emoji = "📊"
	}

	// 添加关键词以通过钉钉安全验证
	markdown := fmt.Sprintf("# %s %s信号 - %s(%s)\n\n", emoji, signal.Signal, signal.StockName, signal.StockCode)
	markdown += fmt.Sprintf("> **【%s】AI股票分析系统**\n\n", d.Secret)
	markdown += fmt.Sprintf("---\n\n")
	markdown += fmt.Sprintf("**当前价格**: %.2f元\n\n", signal.Price)
	markdown += fmt.Sprintf("**信心度**: %d%%\n\n", signal.Confidence)

	if signal.TargetPrice > 0 {
		markdown += fmt.Sprintf("**目标价格**: %.2f元\n\n", signal.TargetPrice)
	}
	if signal.StopLoss > 0 {
		markdown += fmt.Sprintf("**止损价格**: %.2f元\n\n", signal.StopLoss)
	}
	if signal.RiskReward != "" {
		markdown += fmt.Sprintf("**风险回报比**: %s\n\n", signal.RiskReward)
	}

	markdown += fmt.Sprintf("---\n\n")
	markdown += fmt.Sprintf("**分析原因**:\n\n%s\n\n", signal.Reasoning)
	markdown += fmt.Sprintf("---\n\n")
	markdown += fmt.Sprintf("**时间**: %s\n\n", signal.Timestamp.Format("2006-01-02 15:04:05"))

	return markdown
}

// sendRequest 发送HTTP请求到钉钉
func (d *DingTalkNotifier) sendRequest(message map[string]interface{}) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %w", err)
	}

	// TODO: 如果有Secret，需要进行加签处理
	// 钉钉加签文档: https://open.dingtalk.com/document/robots/custom-robot-access

	resp, err := http.Post(d.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if errcode, ok := result["errcode"].(float64); ok && errcode != 0 {
		return fmt.Errorf("钉钉API错误: %v", result["errmsg"])
	}

	return nil
}

// FeishuNotifier 飞书通知器
type FeishuNotifier struct {
	WebhookURL string
	Secret     string // 签名密钥（可选）
}

// NewFeishuNotifier 创建飞书通知器
func NewFeishuNotifier(webhookURL string, secret string) *FeishuNotifier {
	return &FeishuNotifier{
		WebhookURL: webhookURL,
		Secret:     secret,
	}
}

// SendSignal 发送交易信号到飞书
func (f *FeishuNotifier) SendSignal(signal *TradingSignal) error {
	// 构建富文本消息
	content := f.formatSignalRichText(signal)

	// 飞书消息格式
	message := map[string]interface{}{
		"msg_type": "interactive",
		"card":     content,
	}

	return f.sendRequest(message)
}

// SendMessage 发送普通消息到飞书
func (f *FeishuNotifier) SendMessage(message string) error {
	msg := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": message,
		},
	}
	return f.sendRequest(msg)
}

// formatSignalRichText 格式化信号为飞书卡片
func (f *FeishuNotifier) formatSignalRichText(signal *TradingSignal) map[string]interface{} {
	var emoji string
	var color string
	switch signal.Signal {
	case "BUY":
		emoji = "🚀"
		color = "red"
	case "SELL":
		emoji = "⚠️"
		color = "green"
	case "HOLD":
		emoji = "⏸️"
		color = "yellow"
	default:
		emoji = "📊"
		color = "grey"
	}

	// 飞书卡片消息
	card := map[string]interface{}{
		"config": map[string]bool{
			"wide_screen_mode": true,
		},
		"header": map[string]interface{}{
			"title": map[string]interface{}{
				"tag":     "plain_text",
				"content": fmt.Sprintf("%s %s信号 - %s(%s)", emoji, signal.Signal, signal.StockName, signal.StockCode),
			},
			"template": color,
		},
		"elements": []map[string]interface{}{
			{
				"tag": "div",
				"fields": []map[string]interface{}{
					{
						"is_short": true,
						"text": map[string]string{
							"tag":     "lark_md",
							"content": fmt.Sprintf("**当前价格**\n%.2f元", signal.Price),
						},
					},
					{
						"is_short": true,
						"text": map[string]string{
							"tag":     "lark_md",
							"content": fmt.Sprintf("**信心度**\n%d%%", signal.Confidence),
						},
					},
				},
			},
		},
	}

	// 添加目标价格和止损
	if signal.TargetPrice > 0 || signal.StopLoss > 0 {
		fields := []map[string]interface{}{}
		if signal.TargetPrice > 0 {
			fields = append(fields, map[string]interface{}{
				"is_short": true,
				"text": map[string]string{
					"tag":     "lark_md",
					"content": fmt.Sprintf("**目标价格**\n%.2f元", signal.TargetPrice),
				},
			})
		}
		if signal.StopLoss > 0 {
			fields = append(fields, map[string]interface{}{
				"is_short": true,
				"text": map[string]string{
					"tag":     "lark_md",
					"content": fmt.Sprintf("**止损价格**\n%.2f元", signal.StopLoss),
				},
			})
		}
		card["elements"] = append(card["elements"].([]map[string]interface{}), map[string]interface{}{
			"tag":    "div",
			"fields": fields,
		})
	}

	// 添加分割线
	card["elements"] = append(card["elements"].([]map[string]interface{}), map[string]interface{}{
		"tag": "hr",
	})

	// 添加分析原因
	card["elements"] = append(card["elements"].([]map[string]interface{}), map[string]interface{}{
		"tag": "div",
		"text": map[string]string{
			"tag":     "lark_md",
			"content": fmt.Sprintf("**分析原因**\n%s", signal.Reasoning),
		},
	})

	// 添加时间戳
	card["elements"] = append(card["elements"].([]map[string]interface{}), map[string]interface{}{
		"tag": "note",
		"elements": []map[string]string{
			{
				"tag":     "plain_text",
				"content": signal.Timestamp.Format("2006-01-02 15:04:05"),
			},
		},
	})

	return card
}

// sendRequest 发送HTTP请求到飞书
func (f *FeishuNotifier) sendRequest(message map[string]interface{}) error {
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %w", err)
	}

	// TODO: 如果有Secret，需要进行签名处理
	// 飞书签名文档: https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN

	resp, err := http.Post(f.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if code, ok := result["code"].(float64); ok && code != 0 {
		return fmt.Errorf("飞书API错误: %v", result["msg"])
	}

	return nil
}

// MultiNotifier 多通知器（同时发送到多个平台）
type MultiNotifier struct {
	Notifiers []Notifier
}

// NewMultiNotifier 创建多通知器
func NewMultiNotifier(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{
		Notifiers: notifiers,
	}
}

// SendSignal 发送信号到所有通知器
func (m *MultiNotifier) SendSignal(signal *TradingSignal) error {
	var errors []error
	for _, notifier := range m.Notifiers {
		if err := notifier.SendSignal(signal); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("部分通知器发送失败: %v", errors)
	}
	return nil
}

// SendMessage 发送消息到所有通知器
func (m *MultiNotifier) SendMessage(message string) error {
	var errors []error
	for _, notifier := range m.Notifiers {
		if err := notifier.SendMessage(message); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("部分通知器发送失败: %v", errors)
	}
	return nil
}
