package stock

import (
	"fmt"
	"log"
	"math"
	"nofx/mcp"
	"nofx/notifier"
	"strings"
	"time"
)

// StockAnalyzer 股票分析器
type StockAnalyzer struct {
	TDXClient          *TDXClient
	MCPClient          *mcp.Client
	Notifier           notifier.Notifier
	AnalysisConfig     *AnalysisConfig
	TradingTimeChecker *TradingTimeChecker
}

// AnalysisConfig 分析配置
type AnalysisConfig struct {
	StockCode          string        // 股票代码
	StockName          string        // 股票名称
	ScanInterval       time.Duration // 扫描间隔
	EnableNotification bool          // 是否启用通知
	MinConfidence      int           // 最小信心度阈值（低于此值不发送通知）
}

// NewStockAnalyzer 创建股票分析器
func NewStockAnalyzer(tdxClient *TDXClient, mcpClient *mcp.Client, notif notifier.Notifier, config *AnalysisConfig, tradingTimeChecker *TradingTimeChecker) *StockAnalyzer {
	return &StockAnalyzer{
		TDXClient:          tdxClient,
		MCPClient:          mcpClient,
		Notifier:           notif,
		AnalysisConfig:     config,
		TradingTimeChecker: tradingTimeChecker,
	}
}

// AnalysisResult 分析结果
type AnalysisResult struct {
	StockCode     string                 `json:"stock_code"`
	StockName     string                 `json:"stock_name"`
	CurrentPrice  float64                `json:"current_price"`
	Signal        string                 `json:"signal"` // BUY/SELL/HOLD
	Confidence    int                    `json:"confidence"`
	Reasoning     string                 `json:"reasoning"`
	TargetPrice   float64                `json:"target_price,omitempty"`
	StopLoss      float64                `json:"stop_loss,omitempty"`
	RiskReward    string                 `json:"risk_reward,omitempty"`
	TechnicalData map[string]interface{} `json:"technical_data"`
	Timestamp     time.Time              `json:"timestamp"`
}

// Analyze 执行单次分析
func (a *StockAnalyzer) Analyze() (*AnalysisResult, error) {
	// 0. 检查是否在交易时间内
	if a.TradingTimeChecker != nil && !a.TradingTimeChecker.IsTradingTime(time.Now()) {
		status := a.TradingTimeChecker.GetTradingTimeStatus(time.Now())
		log.Printf("⏸️  非交易时段，跳过分析 | 下次交易时间: %v", status["next_trading_time"])
		return nil, fmt.Errorf("非交易时段")
	}

	log.Printf("📊 开始分析股票 %s(%s)...", a.AnalysisConfig.StockName, a.AnalysisConfig.StockCode)

	// 1. 获取实时行情
	quote, err := a.TDXClient.GetQuote(a.AnalysisConfig.StockCode)
	if err != nil {
		return nil, fmt.Errorf("获取行情失败: %w", err)
	}

	// 2. 获取日K线数据（最近60天）
	dayKline, err := a.TDXClient.GetKline(a.AnalysisConfig.StockCode, "day", 60)
	if err != nil {
		return nil, fmt.Errorf("获取日K线失败: %w", err)
	}

	// 3. 获取30分钟K线数据（最近100条）
	min30Kline, err := a.TDXClient.GetKline(a.AnalysisConfig.StockCode, "minute30", 100)
	if err != nil {
		return nil, fmt.Errorf("获取30分钟K线失败: %w", err)
	}

	// 4. 获取今日分时数据
	minuteData, err := a.TDXClient.GetMinute(a.AnalysisConfig.StockCode, "")
	if err != nil {
		log.Printf("⚠️  获取分时数据失败（可能非交易时间）: %v", err)
		minuteData = nil // 非交易时间可能获取不到，设为nil
	}

	// 5. 计算技术指标
	technicalData := a.calculateTechnicalIndicators(quote, dayKline, min30Kline)

	// 6. 构建AI分析提示词
	prompt := a.buildAnalysisPrompt(quote, dayKline, min30Kline, minuteData, technicalData)

	// 7. 调用AI进行分析
	log.Printf("🤖 调用AI进行深度分析...")
	systemPrompt := "你是一位专业的A股分析师，精通技术分析和市场研判。"
	aiResponse, err := a.MCPClient.CallWithMessages(systemPrompt, prompt)
	if err != nil {
		return nil, fmt.Errorf("AI分析失败: %w", err)
	}

	// 8. 解析AI响应
	result, err := a.parseAIResponse(aiResponse, quote, technicalData)
	if err != nil {
		return nil, fmt.Errorf("解析AI响应失败: %w", err)
	}

	// 9. 发送通知（如果启用且信心度达到阈值）
	if a.AnalysisConfig.EnableNotification &&
		result.Confidence >= a.AnalysisConfig.MinConfidence &&
		(result.Signal == "BUY" || result.Signal == "SELL") {
		a.sendNotification(result)
	}

	return result, nil
}

// calculateTechnicalIndicators 计算技术指标
func (a *StockAnalyzer) calculateTechnicalIndicators(quote *QuoteData, dayKline *KlineData, min30Kline *KlineData) map[string]interface{} {
	data := make(map[string]interface{})

	// 当前价格信息
	currentPrice := PriceToYuan(quote.K.Close)
	data["current_price"] = currentPrice
	data["open_price"] = PriceToYuan(quote.K.Open)
	data["high_price"] = PriceToYuan(quote.K.High)
	data["low_price"] = PriceToYuan(quote.K.Low)
	data["prev_close"] = PriceToYuan(quote.K.Last)

	// 涨跌幅
	if quote.K.Last > 0 {
		changePercent := (float64(quote.K.Close-quote.K.Last) / float64(quote.K.Last)) * 100
		data["change_percent"] = fmt.Sprintf("%.2f%%", changePercent)
	}

	// 成交量和成交额
	data["volume"] = VolumeToShares(quote.TotalHand)
	data["amount"] = AmountToYuan(quote.Amount)

	// 内外盘比
	if quote.InsideDish+quote.OuterDisc > 0 {
		outerRatio := float64(quote.OuterDisc) / float64(quote.InsideDish+quote.OuterDisc) * 100
		data["outer_ratio"] = fmt.Sprintf("%.1f%%", outerRatio)
	}

	// 买卖盘力度
	if len(quote.BuyLevel) > 0 && len(quote.SellLevel) > 0 {
		buyPower := 0
		sellPower := 0
		for _, level := range quote.BuyLevel {
			buyPower += level.Number
		}
		for _, level := range quote.SellLevel {
			sellPower += level.Number
		}
		data["buy_sell_ratio"] = fmt.Sprintf("%.2f", float64(buyPower)/float64(sellPower))
	}

	// 日K线指标（简化版MA和趋势）
	// 注意：K线数据List按时间升序排列，List[0]是最旧的，List[len-1]是最新的
	// 因此计算MA时需要从末尾开始取数据
	if len(dayKline.List) >= 5 {
		listLen := len(dayKline.List)

		// 计算5日均价（使用最近5天）
		sum5 := 0
		for i := listLen - 5; i < listLen; i++ {
			sum5 += dayKline.List[i].Close
		}
		ma5 := PriceToYuan(sum5 / 5)
		data["ma5"] = ma5

		// 计算10日均价
		if len(dayKline.List) >= 10 {
			sum10 := 0
			for i := listLen - 10; i < listLen; i++ {
				sum10 += dayKline.List[i].Close
			}
			ma10 := PriceToYuan(sum10 / 10)
			data["ma10"] = ma10
		}

		// 计算20日均价
		if len(dayKline.List) >= 20 {
			sum20 := 0
			for i := listLen - 20; i < listLen; i++ {
				sum20 += dayKline.List[i].Close
			}
			ma20 := PriceToYuan(sum20 / 20)
			data["ma20"] = ma20
		}

		// 计算60日均价（季线）
		if len(dayKline.List) >= 60 {
			sum60 := 0
			for i := listLen - 60; i < listLen; i++ {
				sum60 += dayKline.List[i].Close
			}
			ma60 := PriceToYuan(sum60 / 60)
			data["ma60"] = ma60
		}
	}

	// 计算简化RSI（相对强弱指标）
	if len(dayKline.List) >= 14 {
		rsi14 := a.calculateRSI(dayKline.List, 14)
		data["rsi14"] = fmt.Sprintf("%.2f", rsi14)
	}

	// 计算近期波动率
	if len(dayKline.List) >= 20 {
		volatility := a.calculateVolatility(dayKline.List, 20)
		data["volatility_20d"] = fmt.Sprintf("%.2f%%", volatility*100)
	}

	return data
}

// calculateRSI 计算RSI指标（简化版）
func (a *StockAnalyzer) calculateRSI(klines []KlineItem, period int) float64 {
	if len(klines) < period+1 {
		return 50.0 // 数据不足返回中性值
	}

	gains := 0.0
	losses := 0.0

	// K线数据按时间升序排列，从末尾往前计算最近period天的RSI
	listLen := len(klines)
	for i := listLen - period; i < listLen; i++ {
		// 当前K线的收盘价与前一根K线的收盘价比较
		if i > 0 {
			change := float64(klines[i].Close - klines[i-1].Close)
			if change > 0 {
				gains += change
			} else {
				losses += -change
			}
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	if avgLoss == 0 {
		return 100.0
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	return rsi
}

// calculateVolatility 计算波动率（标准差）
func (a *StockAnalyzer) calculateVolatility(klines []KlineItem, period int) float64 {
	if len(klines) < period+1 {
		return 0
	}

	// K线数据按时间升序排列，计算最近period天的波动率
	listLen := len(klines)
	returns := make([]float64, period)

	// 计算收益率
	for i := 0; i < period; i++ {
		idx := listLen - period + i
		prevIdx := idx - 1
		if prevIdx >= 0 && klines[prevIdx].Close != 0 {
			returns[i] = float64(klines[idx].Close-klines[prevIdx].Close) / float64(klines[prevIdx].Close)
		} else {
			returns[i] = 0
		}
	}

	// 计算均值
	mean := 0.0
	for _, r := range returns {
		mean += r
	}
	mean /= float64(period)

	// 计算标准差
	variance := 0.0
	for _, r := range returns {
		variance += math.Pow(r-mean, 2)
	}
	variance /= float64(period)

	return math.Sqrt(variance)
}

// buildAnalysisPrompt 构建AI分析提示词
func (a *StockAnalyzer) buildAnalysisPrompt(quote *QuoteData, dayKline *KlineData, min30Kline *KlineData, minuteData *MinuteData, technical map[string]interface{}) string {
	prompt := fmt.Sprintf(`# 股票深度分析任务

你是一位专业的A股分析师，请对以下股票进行深度技术分析，并给出明确的操作建议。

## 基本信息
- **股票代码**: %s
- **股票名称**: %s
- **分析时间**: %s

## 实时行情数据
- **当前价格**: %.2f元
- **今日开盘**: %.2f元
- **最高价**: %.2f元
- **最低价**: %.2f元
- **昨收价**: %.2f元
- **涨跌幅**: %s
- **成交量**: %d股
- **成交额**: %.2f万元
- **外盘占比**: %s（外盘越高说明买盘越强）
- **买卖盘比**: %s（>1说明买盘强于卖盘）

## 五档盘口
**买盘**:
`,
		a.AnalysisConfig.StockCode,
		a.AnalysisConfig.StockName,
		time.Now().Format("2006-01-02 15:04:05"),
		technical["current_price"].(float64),
		technical["open_price"].(float64),
		technical["high_price"].(float64),
		technical["low_price"].(float64),
		technical["prev_close"].(float64),
		technical["change_percent"].(string),
		technical["volume"].(int64),
		AmountToYuan(quote.Amount)/10000,
		technical["outer_ratio"].(string),
		technical["buy_sell_ratio"].(string),
	)

	// 添加买五档
	for i, level := range quote.BuyLevel {
		prompt += fmt.Sprintf("- 买%d: %.2f元 x %d股\n", i+1, PriceToYuan(level.Price), level.Number)
	}

	prompt += "\n**卖盘**:\n"
	// 添加卖五档
	for i, level := range quote.SellLevel {
		prompt += fmt.Sprintf("- 卖%d: %.2f元 x %d股\n", i+1, PriceToYuan(level.Price), level.Number)
	}

	// 添加技术指标
	prompt += fmt.Sprintf(`
## 技术指标
- **MA5**: %.2f元
- **MA10**: %.2f元
- **MA20**: %.2f元
- **MA60**: %.2f元（季线）
- **RSI(14)**: %s
- **近20日波动率**: %s

`,
		technical["ma5"].(float64),
		technical["ma10"].(float64),
		technical["ma20"].(float64),
		technical["ma60"].(float64),
		technical["rsi14"].(string),
		technical["volatility_20d"].(string),
	)

	// 添加K线概况
	prompt += fmt.Sprintf(`## K线数据概况
- **日K线**: 最近%d个交易日数据
- **30分钟K线**: 最近%d条数据
`,
		len(dayKline.List),
		len(min30Kline.List),
	)

	// 添加近期价格趋势（从最近5天开始，从新到旧显示）
	if len(dayKline.List) >= 5 {
		prompt += "\n**近5日收盘价趋势**:\n"
		listLen := len(dayKline.List)
		// 从最新的一天开始倒序显示
		for i := listLen - 1; i >= listLen-5 && i >= 0; i-- {
			kline := dayKline.List[i]
			prompt += fmt.Sprintf("- %s: %.2f元 (成交量: %d手)\n",
				kline.Time.Format("01-02"),
				PriceToYuan(kline.Close),
				kline.Volume)
		}
	}

	// 分析要求
	prompt += `
## 分析要求

请基于以上数据进行**全面的技术分析**，并给出明确的操作建议。分析时请考虑：

1. **趋势分析**: 当前价格与均线的关系，是否处于上升/下降/盘整趋势
2. **量价关系**: 成交量的变化是否支持价格走势
3. **盘口分析**: 买卖盘力量对比，大单情况
4. **技术指标**: RSI是否超买超卖，均线排列情况
5. **风险评估**: 当前位置的风险收益比

## 输出格式

请严格按照以下JSON格式输出（只输出JSON，不要其他文字）:

` + "```json" + `
{
  "signal": "BUY 或 SELL 或 HOLD",
  "confidence": 0-100的整数（信心度，越高越确定）,
  "reasoning": "详细的分析理由，包含关键技术指标和逻辑",
  "target_price": 目标价格（元，数字），如果是SELL或HOLD可以为0,
  "stop_loss": 止损价格（元，数字），如果是HOLD可以为0,
  "risk_reward": "风险回报比，例如 1:2 或 1:3"
}
` + "```" + `

**注意事项**:
- signal只能是 "BUY"、"SELL" 或 "HOLD" 三个值之一
- confidence是0-100的整数，代表你的信心程度
- reasoning要详细说明你的分析逻辑和关键依据
- 如果是BUY信号，必须给出target_price和stop_loss
- 如果是SELL信号，应该给出止损建议
- 如果是HOLD，说明原因（如趋势不明、等待突破等）
`

	return prompt
}

// parseAIResponse 解析AI响应
func (a *StockAnalyzer) parseAIResponse(aiResponse string, quote *QuoteData, technical map[string]interface{}) (*AnalysisResult, error) {
	// 1. 解析AI响应中的JSON决策
	aiDecision, err := ParseAIResponse(aiResponse)
	if err != nil {
		// 如果解析失败，记录完整响应并返回默认HOLD信号
		log.Printf("⚠️  AI响应解析失败: %v", err)
		log.Printf("AI原始响应:\n%s", aiResponse)

		return &AnalysisResult{
			StockCode:     a.AnalysisConfig.StockCode,
			StockName:     a.AnalysisConfig.StockName,
			CurrentPrice:  technical["current_price"].(float64),
			Signal:        "HOLD",
			Confidence:    30,
			Reasoning:     fmt.Sprintf("AI响应解析失败，建议观望。原始响应: %s", aiResponse),
			TechnicalData: technical,
			Timestamp:     time.Now(),
		}, nil
	}

	// 2. 验证决策合理性
	currentPrice := technical["current_price"].(float64)
	warnings := ValidateDecision(aiDecision, currentPrice)
	if len(warnings) > 0 {
		log.Printf("⚠️  决策验证警告:")
		for _, warning := range warnings {
			log.Printf("   - %s", warning)
		}
		// 将警告添加到reasoning中
		aiDecision.Reasoning += "\n\n【系统提示】\n" + strings.Join(warnings, "\n")
	}

	// 3. 转换为分析结果
	result := ConvertToAnalysisResult(
		aiDecision,
		a.AnalysisConfig.StockCode,
		a.AnalysisConfig.StockName,
		currentPrice,
		technical,
	)

	// 4. 记录决策日志
	log.Printf("✓ AI决策: %s | 信号: %s | 信心度: %d%%",
		a.AnalysisConfig.StockName,
		result.Signal,
		result.Confidence)

	if result.Signal == "BUY" {
		log.Printf("  目标价: %.2f | 止损价: %.2f | 风险回报比: %s",
			result.TargetPrice, result.StopLoss, result.RiskReward)
	}

	return result, nil
}

// sendNotification 发送通知
func (a *StockAnalyzer) sendNotification(result *AnalysisResult) {
	if a.Notifier == nil {
		return
	}

	signal := &notifier.TradingSignal{
		StockCode:     result.StockCode,
		StockName:     result.StockName,
		Signal:        result.Signal,
		Price:         result.CurrentPrice,
		Confidence:    result.Confidence,
		Reasoning:     result.Reasoning,
		TargetPrice:   result.TargetPrice,
		StopLoss:      result.StopLoss,
		RiskReward:    result.RiskReward,
		Timestamp:     result.Timestamp,
		TechnicalData: result.TechnicalData,
	}

	if err := a.Notifier.SendSignal(signal); err != nil {
		log.Printf("❌ 发送通知失败: %v", err)
	} else {
		log.Printf("✅ 已发送%s信号通知", result.Signal)
	}
}

// StartMonitoring 启动持续监控
func (a *StockAnalyzer) StartMonitoring(stopChan <-chan struct{}) {
	ticker := time.NewTicker(a.AnalysisConfig.ScanInterval)
	defer ticker.Stop()

	log.Printf("🚀 开始监控股票 %s(%s)，扫描间隔: %v",
		a.AnalysisConfig.StockName,
		a.AnalysisConfig.StockCode,
		a.AnalysisConfig.ScanInterval)

	// 立即执行一次分析
	if _, err := a.Analyze(); err != nil {
		log.Printf("❌ 分析失败: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if _, err := a.Analyze(); err != nil {
				log.Printf("❌ 分析失败: %v", err)
			}
		case <-stopChan:
			log.Printf("⏹️  停止监控股票 %s", a.AnalysisConfig.StockCode)
			return
		}
	}
}
