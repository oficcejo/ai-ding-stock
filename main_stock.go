package main

import (
	"fmt"
	"log"
	"nofx/api"
	"nofx/config"
	"nofx/mcp"
	"nofx/notifier"
	"nofx/stock"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║    📈 AI股票分析系统 - 实时分析与信号通知               ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// 加载配置文件
	configFile := "config_stock.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	log.Printf("📋 加载配置文件: %s", configFile)
	cfg, err := config.LoadStockConfig(configFile)
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}

	log.Printf("✓ 配置加载成功")
	fmt.Println()

	// 创建TDX客户端
	tdxClient := stock.NewTDXClient(cfg.TDXAPIUrl)
	log.Printf("✓ TDX API客户端已初始化: %s", cfg.TDXAPIUrl)

	// 创建AI客户端
	mcpClient, err := createMCPClient(&cfg.AIConfig)
	if err != nil {
		log.Fatalf("❌ 创建AI客户端失败: %v", err)
	}
	log.Printf("✓ AI客户端已初始化 (%s)", strings.ToUpper(cfg.AIConfig.Provider))

	// 创建通知器
	var notif notifier.Notifier
	if cfg.Notification.Enabled {
		notif = createNotifier(&cfg.Notification)
		log.Printf("✓ 通知系统已初始化")
	} else {
		log.Printf("⏭️  通知系统未启用")
	}

	// 创建交易时间检查器
	tradingTimeConfig := stock.TradingTimeConfig{
		EnableTradingTimeCheck: cfg.TradingTime.EnableCheck,
		TradingHours:           cfg.TradingTime.TradingHours,
		Timezone:               cfg.TradingTime.Timezone,
	}
	tradingTimeChecker, err := stock.NewTradingTimeChecker(tradingTimeConfig)
	if err != nil {
		log.Printf("⚠️  创建交易时间检查器失败: %v, 将禁用交易时间检查", err)
		tradingTimeChecker = nil
	} else if cfg.TradingTime.EnableCheck {
		log.Printf("✓ 交易时间检查已启用")
		log.Printf("  交易时段: %v", cfg.TradingTime.TradingHours)
		status := tradingTimeChecker.GetTradingTimeStatus(time.Now())
		log.Printf("  当前状态: 交易日=%v, 交易时段=%v",
			status["is_trading_day"], status["is_trading_time"])
	} else {
		log.Printf("⏭️  交易时间检查未启用（将持续分析）")
	}

	// 创建日志目录
	if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
		log.Printf("⚠️  创建日志目录失败: %v", err)
	}

	fmt.Println()
	fmt.Println("📊 监控股票列表:")
	enabledStocks := []config.StockItem{}
	for _, stockItem := range cfg.Stocks {
		if stockItem.Enabled {
			enabledStocks = append(enabledStocks, stockItem)
			fmt.Printf("  • %s(%s) - 扫描间隔: %d分钟, 信心阈值: %d%%\n",
				stockItem.Name, stockItem.Code, stockItem.ScanIntervalMinutes, stockItem.MinConfidence)
		}
	}

	fmt.Println()
	fmt.Println("🤖 AI分析模式:")
	fmt.Println("  • AI将基于实时行情、K线、技术指标进行全面分析")
	fmt.Println("  • 提供BUY/SELL/HOLD明确信号")
	fmt.Println("  • 给出目标价位和止损建议")
	fmt.Println("  • 信心度≥阈值时发送通知")
	fmt.Println()
	fmt.Println("⚠️  风险提示: AI分析仅供参考，投资有风险，决策需谨慎！")
	fmt.Println()
	fmt.Println("按 Ctrl+C 停止运行")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	// 创建分析器管理器
	analyzerManager := &AnalyzerManager{
		analyzers: make(map[string]*stock.StockAnalyzer),
		stopChans: make(map[string]chan struct{}),
	}

	// 为每只启用的股票创建分析器
	for _, stockItem := range enabledStocks {
		analysisConfig := &stock.AnalysisConfig{
			StockCode:          stockItem.Code,
			StockName:          stockItem.Name,
			ScanInterval:       stockItem.GetScanInterval(),
			EnableNotification: cfg.Notification.Enabled,
			MinConfidence:      stockItem.MinConfidence,
		}

		analyzer := stock.NewStockAnalyzer(tdxClient, mcpClient, notif, analysisConfig, tradingTimeChecker)
		analyzerManager.AddAnalyzer(stockItem.Code, analyzer)
	}

	// 创建并启动API服务器
	apiServer := api.NewStockAPIServer(analyzerManager, cfg.APIServerPort)
	go func() {
		if err := apiServer.Start(); err != nil {
			log.Printf("❌ API服务器错误: %v", err)
		}
	}()
	log.Printf("✓ API服务器已启动: http://localhost:%d", cfg.APIServerPort)
	fmt.Println()

	// 设置优雅退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 启动所有分析器
	analyzerManager.StartAll()

	// 等待退出信号
	<-sigChan
	fmt.Println()
	fmt.Println()
	log.Println("📛 收到退出信号，正在停止所有分析器...")
	analyzerManager.StopAll()

	fmt.Println()
	fmt.Println("👋 感谢使用AI股票分析系统！")
}

// createMCPClient 创建MCP客户端
func createMCPClient(aiConfig *config.AIConfig) (*mcp.Client, error) {
	client := mcp.New()

	switch aiConfig.Provider {
	case "deepseek":
		client.SetDeepSeekAPIKey(aiConfig.DeepSeekKey)
	case "qwen":
		client.SetQwenAPIKey(aiConfig.QwenKey, "")
	case "custom":
		client.SetCustomAPI(aiConfig.CustomAPIURL, aiConfig.CustomAPIKey, aiConfig.CustomModelName)
	default:
		return nil, fmt.Errorf("不支持的AI提供商: %s", aiConfig.Provider)
	}

	return client, nil
}

// createNotifier 创建通知器
func createNotifier(notifConfig *config.NotificationConfig) notifier.Notifier {
	var notifiers []notifier.Notifier

	if notifConfig.DingTalk.Enabled {
		ding := notifier.NewDingTalkNotifier(
			notifConfig.DingTalk.WebhookURL,
			notifConfig.DingTalk.Secret,
		)
		notifiers = append(notifiers, ding)
		log.Printf("  ✓ 钉钉通知已启用")
	}

	if notifConfig.Feishu.Enabled {
		feishu := notifier.NewFeishuNotifier(
			notifConfig.Feishu.WebhookURL,
			notifConfig.Feishu.Secret,
		)
		notifiers = append(notifiers, feishu)
		log.Printf("  ✓ 飞书通知已启用")
	}

	if len(notifiers) == 0 {
		return nil
	}

	if len(notifiers) == 1 {
		return notifiers[0]
	}

	return notifier.NewMultiNotifier(notifiers...)
}

// AnalyzerManager 分析器管理器
type AnalyzerManager struct {
	analyzers map[string]*stock.StockAnalyzer
	stopChans map[string]chan struct{}
	mutex     sync.RWMutex
}

// AddAnalyzer 添加分析器
func (m *AnalyzerManager) AddAnalyzer(code string, analyzer *stock.StockAnalyzer) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.analyzers[code] = analyzer
	m.stopChans[code] = make(chan struct{})
}

// GetAnalyzer 获取分析器
func (m *AnalyzerManager) GetAnalyzer(code string) interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.analyzers[code]
}

// StartAll 启动所有分析器
func (m *AnalyzerManager) StartAll() {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for code, analyzer := range m.analyzers {
		stopChan := m.stopChans[code]
		go analyzer.StartMonitoring(stopChan)
	}
}

// StopAll 停止所有分析器
func (m *AnalyzerManager) StopAll() {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, stopChan := range m.stopChans {
		close(stopChan)
	}
}

// GetAllAnalyzers 获取所有分析器
func (m *AnalyzerManager) GetAllAnalyzers() map[string]interface{} {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[string]interface{})
	for code, analyzer := range m.analyzers {
		result[code] = analyzer
	}
	return result
}
