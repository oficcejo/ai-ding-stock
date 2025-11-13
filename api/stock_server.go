package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// StockAPIServer 股票分析API服务器
type StockAPIServer struct {
	router  *gin.Engine
	manager AnalyzerManagerInterface
	port    int
}

// AnalyzerManagerInterface 分析器管理器接口
type AnalyzerManagerInterface interface {
	GetAnalyzer(code string) interface{}
	GetAllAnalyzers() map[string]interface{}
}

// NewStockAPIServer 创建股票API服务器
func NewStockAPIServer(manager AnalyzerManagerInterface, port int) *StockAPIServer {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// 配置CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	server := &StockAPIServer{
		router:  router,
		manager: manager,
		port:    port,
	}

	server.setupRoutes()
	return server
}

// setupRoutes 设置路由
func (s *StockAPIServer) setupRoutes() {
	// 健康检查
	s.router.GET("/health", s.handleHealth)

	// 静态文件服务
	s.router.Static("/static", "./web/static")
	s.router.StaticFile("/", "./web/config.html")
	s.router.StaticFile("/config", "./web/config.html")

	// API路由组
	api := s.router.Group("/api")
	{
		// 配置管理接口
		api.GET("/config", s.handleGetConfig)
		api.POST("/config", s.handleSaveConfig)

		// 获取所有监控股票列表
		api.GET("/stocks", s.handleGetStocks)

		// 获取单个股票的最新分析结果
		api.GET("/stock/:code/latest", s.handleGetLatestAnalysis)

		// 获取单个股票的历史分析记录
		api.GET("/stock/:code/history", s.handleGetAnalysisHistory)

		// 手动触发分析
		api.POST("/stock/:code/analyze", s.handleTriggerAnalysis)

		// 获取系统统计信息
		api.GET("/statistics", s.handleGetStatistics)
	}
}

// handleHealth 健康检查
func (s *StockAPIServer) handleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format("2006-01-02 15:04:05"),
	})
}

// handleGetStocks 获取所有监控股票
func (s *StockAPIServer) handleGetStocks(c *gin.Context) {
	analyzers := s.manager.GetAllAnalyzers()

	stocks := []gin.H{}
	for code := range analyzers {
		// TODO: 获取每个分析器的配置信息
		stocks = append(stocks, gin.H{
			"code":    code,
			"name":    "", // 需要从analyzer获取
			"enabled": true,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total":  len(stocks),
			"stocks": stocks,
		},
	})
}

// handleGetLatestAnalysis 获取最新分析结果
func (s *StockAPIServer) handleGetLatestAnalysis(c *gin.Context) {
	code := c.Param("code")

	analyzer := s.manager.GetAnalyzer(code)
	if analyzer == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "未找到该股票的分析器",
		})
		return
	}

	// TODO: 从analyzer获取最新分析结果
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"stock_code": code,
			"timestamp":  time.Now(),
			// 更多数据...
		},
	})
}

// handleGetAnalysisHistory 获取历史分析记录
func (s *StockAPIServer) handleGetAnalysisHistory(c *gin.Context) {
	code := c.Param("code")
	limit := 20 // 默认返回最近20条

	analyzer := s.manager.GetAnalyzer(code)
	if analyzer == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "未找到该股票的分析器",
		})
		return
	}

	// TODO: 从日志文件读取历史记录
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"stock_code": code,
			"count":      0,
			"limit":      limit,
			"records":    []gin.H{},
		},
	})
}

// handleTriggerAnalysis 手动触发分析
func (s *StockAPIServer) handleTriggerAnalysis(c *gin.Context) {
	code := c.Param("code")

	analyzer := s.manager.GetAnalyzer(code)
	if analyzer == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": "未找到该股票的分析器",
		})
		return
	}

	// TODO: 触发立即分析
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "分析任务已提交",
		"data": gin.H{
			"stock_code": code,
			"triggered":  true,
		},
	})
}

// handleGetStatistics 获取系统统计
func (s *StockAPIServer) handleGetStatistics(c *gin.Context) {
	analyzers := s.manager.GetAllAnalyzers()

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total_stocks":   len(analyzers),
			"system_uptime":  "", // TODO: 计算运行时间
			"total_analysis": 0,  // TODO: 统计总分析次数
		},
	})
}

// handleGetConfig 获取配置
func (s *StockAPIServer) handleGetConfig(c *gin.Context) {
	// 读取配置文件
	configFile := "config_stock.json"
	data, err := os.ReadFile(configFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": fmt.Sprintf("读取配置文件失败: %v", err),
		})
		return
	}

	// 解析为JSON对象
	var config map[string]interface{}
	if err := json.Unmarshal(data, &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": fmt.Sprintf("解析配置文件失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    config,
	})
}

// handleSaveConfig 保存配置
func (s *StockAPIServer) handleSaveConfig(c *gin.Context) {
	var config map[string]interface{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    -1,
			"message": fmt.Sprintf("请求数据格式错误: %v", err),
		})
		return
	}

	// 转换为格式化的JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": fmt.Sprintf("序列化配置失败: %v", err),
		})
		return
	}

	// 备份原配置文件
	configFile := "config_stock.json"
	backupFile := fmt.Sprintf("config_stock.json.backup.%s", time.Now().Format("20060102150405"))
	if err := os.Rename(configFile, backupFile); err != nil {
		log.Printf("⚠️  备份配置文件失败: %v", err)
	} else {
		log.Printf("✓ 配置文件已备份: %s", backupFile)
	}

	// 写入新配置
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    -1,
			"message": fmt.Sprintf("保存配置文件失败: %v", err),
		})
		return
	}

	log.Printf("✓ 配置文件已更新: %s", configFile)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "配置保存成功，请重启程序使配置生效",
		"data": gin.H{
			"backup_file": backupFile,
		},
	})
}

// Start 启动服务器
func (s *StockAPIServer) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("🚀 股票分析API服务器启动在端口 %d", s.port)
	return s.router.Run(addr)
}
