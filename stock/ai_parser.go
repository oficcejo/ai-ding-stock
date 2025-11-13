package stock

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// AIDecisionResponse AI决策响应结构
type AIDecisionResponse struct {
	Signal      string  `json:"signal"`       // BUY/SELL/HOLD
	Confidence  int     `json:"confidence"`   // 0-100
	Reasoning   string  `json:"reasoning"`    // 分析理由
	TargetPrice float64 `json:"target_price"` // 目标价格
	StopLoss    float64 `json:"stop_loss"`    // 止损价格
	RiskReward  string  `json:"risk_reward"`  // 风险回报比
}

// ParseAIResponse 解析AI响应，提取JSON决策
func ParseAIResponse(response string) (*AIDecisionResponse, error) {
	// 尝试多种方式提取JSON

	// 方式1: 查找```json ... ```代码块
	jsonPattern := regexp.MustCompile("(?s)```json\\s*({.*?})\\s*```")
	matches := jsonPattern.FindStringSubmatch(response)

	var jsonStr string
	if len(matches) >= 2 {
		jsonStr = matches[1]
	} else {
		// 方式2: 查找{ ... }对象
		objectPattern := regexp.MustCompile("(?s){[^{}]*\"signal\"[^{}]*}")
		matches = objectPattern.FindStringSubmatch(response)
		if len(matches) >= 1 {
			jsonStr = matches[0]
		} else {
			// 方式3: 尝试直接解析整个响应
			jsonStr = response
		}
	}

	// 清理JSON字符串
	jsonStr = strings.TrimSpace(jsonStr)

	// 解析JSON
	var decision AIDecisionResponse
	if err := json.Unmarshal([]byte(jsonStr), &decision); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w\n原始响应:\n%s", err, response)
	}

	// 验证必填字段
	if decision.Signal == "" {
		return nil, fmt.Errorf("AI响应缺少signal字段")
	}

	// 规范化signal值
	decision.Signal = strings.ToUpper(strings.TrimSpace(decision.Signal))
	if decision.Signal != "BUY" && decision.Signal != "SELL" && decision.Signal != "HOLD" {
		return nil, fmt.Errorf("无效的signal值: %s (必须是BUY/SELL/HOLD)", decision.Signal)
	}

	// 验证信心度范围
	if decision.Confidence < 0 || decision.Confidence > 100 {
		// 尝试修正
		if decision.Confidence < 0 {
			decision.Confidence = 0
		} else if decision.Confidence > 100 {
			decision.Confidence = 100
		}
	}

	// 验证BUY信号必须有目标价和止损
	if decision.Signal == "BUY" {
		if decision.TargetPrice == 0 {
			return nil, fmt.Errorf("BUY信号必须设置target_price")
		}
		if decision.StopLoss == 0 {
			return nil, fmt.Errorf("BUY信号必须设置stop_loss")
		}
	}

	return &decision, nil
}

// ConvertToAnalysisResult 将AI决策转换为分析结果
func ConvertToAnalysisResult(aiDecision *AIDecisionResponse, stockCode, stockName string, currentPrice float64, technical map[string]interface{}) *AnalysisResult {
	return &AnalysisResult{
		StockCode:     stockCode,
		StockName:     stockName,
		CurrentPrice:  currentPrice,
		Signal:        aiDecision.Signal,
		Confidence:    aiDecision.Confidence,
		Reasoning:     aiDecision.Reasoning,
		TargetPrice:   aiDecision.TargetPrice,
		StopLoss:      aiDecision.StopLoss,
		RiskReward:    aiDecision.RiskReward,
		TechnicalData: technical,
		Timestamp:     time.Now(),
	}
}

// ValidateDecision 验证决策的合理性
func ValidateDecision(decision *AIDecisionResponse, currentPrice float64) []string {
	var warnings []string

	// 检查BUY信号
	if decision.Signal == "BUY" {
		// 目标价应该高于当前价
		if decision.TargetPrice <= currentPrice {
			warnings = append(warnings, fmt.Sprintf("目标价%.2f应该高于当前价%.2f", decision.TargetPrice, currentPrice))
		}

		// 止损价应该低于当前价
		if decision.StopLoss >= currentPrice {
			warnings = append(warnings, fmt.Sprintf("止损价%.2f应该低于当前价%.2f", decision.StopLoss, currentPrice))
		}

		// 检查风险回报比是否合理（至少1:1.5）
		if decision.TargetPrice > currentPrice && decision.StopLoss < currentPrice {
			reward := decision.TargetPrice - currentPrice
			risk := currentPrice - decision.StopLoss
			ratio := reward / risk
			if ratio < 1.5 {
				warnings = append(warnings, fmt.Sprintf("风险回报比%.2f偏低，建议至少1:1.5", ratio))
			}
		}
	}

	// 检查SELL信号
	if decision.Signal == "SELL" {
		// SELL通常意味着要减仓或止损
		if decision.StopLoss > 0 && decision.StopLoss <= currentPrice {
			warnings = append(warnings, fmt.Sprintf("SELL信号的止损价%.2f应该高于当前价%.2f", decision.StopLoss, currentPrice))
		}
	}

	// 检查信心度
	if decision.Confidence < 50 && (decision.Signal == "BUY" || decision.Signal == "SELL") {
		warnings = append(warnings, fmt.Sprintf("信心度%d%%偏低，建议谨慎操作", decision.Confidence))
	}

	return warnings
}
