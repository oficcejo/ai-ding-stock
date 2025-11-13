# 📈 AI股票分析系统

> 基于DeepSeek/Qwen大模型的智能股票分析系统，实时监控、AI分析、自动通知

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://www.docker.com/)

---

## 🌟 项目特点

- 🤖 **AI驱动分析** - 使用DeepSeek/Qwen大模型进行深度技术分析
- 📊 **实时监控** - 自动监控多只股票，定时分析
- 🎯 **智能信号** - 提供BUY/SELL/HOLD明确信号和目标价
- 📱 **即时通知** - 支持钉钉、飞书Webhook推送
- 🌐 **Web界面** - 实时查看分析结果和历史记录
- 🔌 **RESTful API** - 完整的API接口，易于集成
- 🐳 **容器化部署** - Docker一键部署，开箱即用
- 📈 **技术指标** - 支持MA、RSI、波动率等多种技术指标

---

## 📋 目录

- [功能概述](#功能概述)
- [系统架构](#系统架构)
- [技术栈](#技术栈)
- [快速开始](#快速开始)
  - [Docker部署（推荐）](#docker部署推荐)
  - [本地运行](#本地运行)
- [配置说明](#配置说明)
- [使用指南](#使用指南)
- [API文档](#api文档)
- [通知配置](#通知配置)
- [常见问题](#常见问题)
- [更新日志](#更新日志)

---

## 🎯 功能概述

### 核心功能

#### 1. 实时股票监控
- ✅ 支持多只股票同时监控
- ✅ 可配置扫描间隔（1-60分钟）
- ✅ 自动获取实时行情和K线数据
- ✅ 24/7不间断运行

#### 2. AI深度分析
- ✅ 基于大语言模型的技术分析
- ✅ 综合考虑趋势、量价、盘口等多维度
- ✅ 给出BUY/SELL/HOLD明确信号
- ✅ 提供信心度评分（0-100）
- ✅ 给出目标价和止损价建议

#### 3. 技术指标计算
- ✅ **均线系统**: MA5、MA10、MA20、MA60
- ✅ **相对强弱**: RSI(14)
- ✅ **波动率**: 20日标准差
- ✅ **量价分析**: 成交量、成交额、内外盘比
- ✅ **盘口分析**: 买卖五档、委比

#### 4. 智能通知
- ✅ 钉钉机器人推送
- ✅ 飞书机器人推送
- ✅ 可配置信心度阈值
- ✅ 支持自定义关键词
- ✅ 富文本/Markdown格式

#### 5. Web监控界面
- ✅ 实时显示分析结果
- ✅ 股票列表和状态
- ✅ AI分析详情
- ✅ 历史信号记录
- ✅ 响应式设计

#### 6. RESTful API
- ✅ 获取所有股票状态
- ✅ 查询单个股票分析
- ✅ 健康检查端点
- ✅ 标准JSON响应

---

## 🏗️ 系统架构

```
┌─────────────────────────────────────────────────────────┐
│                    AI股票分析系统                        │
└─────────────────────────────────────────────────────────┘

┌─────────────┐      ┌──────────────┐      ┌─────────────┐
│  TDX API    │─────▶│  数据采集层   │◀────▶│   数据源    │
│  (行情接口)  │      │  (实时行情)   │      │  (K线数据)  │
└─────────────┘      └──────────────┘      └─────────────┘
                            │
                            ▼
                     ┌──────────────┐
                     │  技术指标层   │
                     │ (MA/RSI/波动) │
                     └──────────────┘
                            │
                            ▼
┌─────────────┐      ┌──────────────┐      ┌─────────────┐
│  DeepSeek   │◀────▶│   AI分析层   │─────▶│   决策引擎  │
│   / Qwen    │      │  (智能分析)   │      │ (BUY/SELL)  │
└─────────────┘      └──────────────┘      └─────────────┘
                            │
              ┌─────────────┼─────────────┐
              ▼             ▼             ▼
       ┌───────────┐ ┌───────────┐ ┌───────────┐
       │ 钉钉通知  │ │ 飞书通知  │ │  Web界面  │
       └───────────┘ └───────────┘ └───────────┘
```

### 数据流程

```
1. 定时任务触发
   └─▶ 获取实时行情（五档、价格、成交量）
       └─▶ 获取K线数据（日K、30分钟K）
           └─▶ 计算技术指标（MA、RSI、波动率）
               └─▶ 构建AI分析提示词
                   └─▶ 调用AI模型分析
                       └─▶ 解析AI决策
                           └─▶ 验证决策合理性
                               └─▶ 存储分析结果
                                   └─▶ 判断是否通知
                                       └─▶ 发送推送通知
```

---

## 💻 技术栈

### 后端
- **语言**: Go 1.24+
- **框架**: Gin (Web框架)
- **AI**: DeepSeek API / Qwen API
- **并发**: Goroutine + Channel
- **配置**: JSON

### 前端
- **技术**: HTML5 + CSS3 + Vanilla JavaScript
- **设计**: 响应式布局
- **风格**: 现代化渐变UI

### 数据源
- **行情API**: 通达信TDX API
- **数据格式**: JSON
- **更新频率**: 实时

### 部署
- **容器**: Docker + Docker Compose
- **Web服务器**: Nginx
- **网络模式**: Host Network
- **日志**: 文件 + 控制台

---

## 🚀 快速开始

### 前置要求

- Docker 20.10+ 和 Docker Compose 2.0+（Docker部署）
- 或 Go 1.24+（本地运行）
- TDX股票数据API服务
- DeepSeek 或 Qwen API密钥
- docker本地部署tdxapi，项目地址https://github.com/oficcejo/tdx-api
---

## 🐳 Docker部署（推荐）

### 1️⃣ 克隆项目

```bash
git clone https://github.com/oficcejo/ai-ding-stock.git
cd ai-ding-stock
```

### 2️⃣ 配置系统

```bash
# 复制配置示例
cp config_stock.json.example config_stock.json

# 编辑配置文件
nano config_stock.json
```

**最小配置**：

```json
{
  "tdx_api_url": "http://your-tdx-api:8181",  # TDX实时行情 API项目地址https://github.com/oficcejo/tdx-api
  "ai_config": {
    "provider": "deepseek",
    "deepseek_key": "sk-your-deepseek-api-key",
    "qwen_key": "",
    "custom_api_url": "",
    "custom_api_key": "",
    "custom_model_name": ""
  },
  "stocks": [
    {
      "code": "600519",
      "name": "贵州茅台",
      "enabled": true,
      "scan_interval_minutes": 5,
      "min_confidence": 70
    }
  ],
  "notification": {
    "enabled": true,
    "dingtalk": {
      "enabled": true,
      "webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=YOUR_TOKEN",
      "secret": "关键词（如：股票通知）"
    },
    "feishu": {
      "enabled": false,
      "webhook_url": "",
      "secret": ""
    }
  },
  "api_server_port": 9090,
  "log_dir": "logs"
}
```

### 3️⃣ 启动服务

#### Linux/macOS

```bash
# 添加执行权限
chmod +x docker-start.sh

# 启动服务
bash docker-start.sh start

# 或使用docker-compose
docker-compose up -d --build
```

#### Windows

```cmd
REM 直接双击运行
docker-start.bat

REM 或命令行
docker-start.bat start
```

### 4️⃣ 访问系统

- **Web界面**: http://localhost:8090
- **API接口**: http://localhost:9090/api/stocks

### 5️⃣ 管理服务

```bash
# 查看状态
bash docker-start.sh status

# 查看日志
bash docker-start.sh logs

# 重启服务
bash docker-start.sh restart

# 停止服务
bash docker-start.sh stop
```

---

## 💻 本地运行

### 1️⃣ 安装依赖

```bash
# 安装Go 1.24+
# https://golang.org/dl/

# 验证安装
go version
```

### 2️⃣ 下载项目

```bash
git clone https://github.com/your-repo/nofx-stock.git
cd nofx-stock
```

### 3️⃣ 安装依赖

```bash
# 使用国内代理（可选）
go env -w GOPROXY=https://goproxy.cn,direct

# 下载依赖
go mod download
```

### 4️⃣ 配置系统

```bash
# 复制配置
cp config_stock.json.example config_stock.json

# 编辑配置
vim config_stock.json
```

### 5️⃣ 编译运行

```bash
# 编译
go build -o stock_analyzer main_stock.go

# 运行
./stock_analyzer

# Windows
stock_analyzer.exe
```

### 6️⃣ 访问系统

- **Web界面**: http://localhost:8090/stock_dashboard.html
- **API接口**: http://localhost:9090/api/stocks

---

## ⚙️ 配置说明

### 配置文件结构

```json
{
  "tdx_api_url": "TDX API地址",
  "ai_config": {
    "provider": "AI提供商: deepseek/qwen/custom",
    "deepseek_key": "DeepSeek API密钥",
    "qwen_key": "Qwen API密钥",
    "custom_api_url": "自定义API地址（可选）",
    "custom_api_key": "自定义API密钥（可选）",
    "custom_model_name": "自定义模型名称（可选）"
  },
  "stocks": [
    {
      "code": "股票代码",
      "name": "股票名称",
      "enabled": true,
      "scan_interval_minutes": 5,
      "min_confidence": 70
    }
  ],
  "notification": {
    "enabled": true,
    "dingtalk": {
      "enabled": true,
      "webhook_url": "钉钉机器人Webhook",
      "secret": "安全关键词"
    },
    "feishu": {
      "enabled": false,
      "webhook_url": "飞书机器人Webhook",
      "secret": "签名密钥"
    }
  },
  "api_server_port": 9090,
  "log_dir": "logs"
}
```

### 配置项说明

#### TDX API配置
- `tdx_api_url`: TDX股票数据API的基础URL

#### AI配置
- `provider`: AI提供商，支持：
  - `deepseek`: 使用DeepSeek API
  - `qwen`: 使用通义千问API
  - `custom`: 使用自定义API
- `deepseek_key`: DeepSeek API密钥
- `qwen_key`: 通义千问API密钥

#### 股票配置
- `code`: 股票代码（如：600519）
- `name`: 股票名称（如：贵州茅台）
- `enabled`: 是否启用监控
- `scan_interval_minutes`: 扫描间隔（分钟），建议5-60
- `min_confidence`: 最小信心度阈值（0-100），低于此值不发送通知

#### 通知配置
- `enabled`: 是否启用通知
- `dingtalk.webhook_url`: 钉钉机器人Webhook地址
- `dingtalk.secret`: 钉钉机器人关键词（用于安全验证）
- `feishu.webhook_url`: 飞书机器人Webhook地址

#### 服务配置
- `api_server_port`: API服务器端口（默认9090）
- `log_dir`: 日志目录

---

## 📖 使用指南

### 添加监控股票

编辑 `config_stock.json`：

```json
{
  "stocks": [
    {
      "code": "600519",
      "name": "贵州茅台",
      "enabled": true,
      "scan_interval_minutes": 5,
      "min_confidence": 70
    },
    {
      "code": "000001",
      "name": "平安银行",
      "enabled": true,
      "scan_interval_minutes": 10,
      "min_confidence": 75
    }
  ]
}
```

重启服务：

```bash
# Docker
docker-compose restart stock-analyzer

# 本地
# Ctrl+C 停止，然后重新运行
./stock_analyzer
```

### 调整分析频率

修改 `scan_interval_minutes`：

```json
{
  "code": "600519",
  "scan_interval_minutes": 15  // 改为15分钟
}
```

### 调整通知阈值

修改 `min_confidence`：

```json
{
  "code": "600519",
  "min_confidence": 80  // 提高到80%才通知
}
```

### 临时禁用股票

设置 `enabled` 为 `false`：

```json
{
  "code": "600519",
  "enabled": false  // 暂时不监控
}
```

---

## 🔌 API文档

### 基础信息

- **Base URL**: `http://localhost:9090`
- **Content-Type**: `application/json`
- **响应格式**: JSON

### API端点

#### 1. 获取所有股票状态

```http
GET /api/stocks
```

**响应示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 2,
    "stocks": [
      {
        "code": "600519",
        "name": "贵州茅台",
        "enabled": true
      },
      {
        "code": "000001",
        "name": "平安银行",
        "enabled": true
      }
    ]
  }
}
```

#### 2. 获取单个股票分析

```http
GET /api/stock/:code
```

**参数**：
- `code`: 股票代码

**响应示例**：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "stock_code": "600519",
    "stock_name": "贵州茅台",
    "current_price": 1234.56,
    "signal": "BUY",
    "confidence": 85,
    "reasoning": "技术面分析：当前价格位于MA5之上...",
    "target_price": 1300.00,
    "stop_loss": 1200.00,
    "risk_reward": "1:3",
    "timestamp": "2025-11-04T15:30:00Z",
    "technical_data": {
      "ma5": 1230.00,
      "ma10": 1220.00,
      "rsi14": "58.50",
      "volatility_20d": "2.50%"
    }
  }
}
```

#### 3. 健康检查

```http
GET /health
```

**响应示例**：

```json
{
  "status": "ok",
  "timestamp": "2025-11-04T15:30:00Z"
}
```

---

## 📱 通知配置

### 钉钉机器人

#### 1. 创建机器人

1. 打开钉钉群
2. 群设置 → 智能群助手 → 添加机器人
3. 选择"自定义"机器人
4. 设置机器人名称和头像
5. **安全设置**：选择"自定义关键词"，填写关键词（如：`股票通知`）
6. 复制Webhook地址

#### 2. 配置系统

```json
{
  "notification": {
    "dingtalk": {
      "enabled": true,
      "webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=YOUR_TOKEN",
      "secret": "股票通知"  // 与机器人安全设置的关键词一致
    }
  }
}
```

#### 3. 通知效果

```markdown
# ⚠️ SELL信号 - 贵州茅台(600519)

> **【股票通知】AI股票分析系统**

---

**当前价格**: 1234.56元
**信心度**: 85%
**目标价格**: 1200.00元
**止损价格**: 1250.00元
**风险回报比**: 1:2

---

**分析原因**:
技术面分析显示...

---

**时间**: 2025-11-04 15:30:00
```

### 飞书机器人

#### 1. 创建机器人

1. 打开飞书群
2. 群设置 → 群机器人 → 添加机器人
3. 选择"自定义机器人"
4. 设置机器人名称
5. 复制Webhook URL

#### 2. 配置系统

```json
{
  "notification": {
    "feishu": {
      "enabled": true,
      "webhook_url": "https://open.feishu.cn/open-apis/bot/v2/hook/YOUR_TOKEN",
      "secret": ""
    }
  }
}
```

---

## 🎨 Web界面

### 功能特性

- 📊 **实时数据**：显示所有监控股票的最新分析
- 🎯 **信号展示**：清晰的BUY/SELL/HOLD标识
- 💯 **信心度**：可视化的信心度百分比
- 📈 **技术指标**：MA、RSI、波动率等
- 📝 **分析详情**：AI的完整推理过程
- ⏰ **时间戳**：最后更新时间
- 🔄 **自动刷新**：定期自动更新数据

### 访问地址

- **Docker部署**: http://服务器IP:8090
- **本地运行**: http://localhost:8090/stock_dashboard.html

---

## ❓ 常见问题

### 1. Docker构建失败

**问题**: `go: module xxx requires go >= 1.24.0`

**解决**:
```bash
# 确认Docker镜像使用正确的Go版本
grep "golang:" Dockerfile
# 应该显示: FROM golang:1.24-alpine
```

### 2. TDX API连接超时

**问题**: `获取行情失败: 超时`

**解决**:
```bash
# 测试TDX API是否可访问
curl http://your-tdx-api:8181

# 如果在Docker中，确保使用host网络模式
# docker-compose.yml中应该有:
# network_mode: host
```

### 3. AI API调用失败

**问题**: `AI分析失败: API key invalid`

**解决**:
- 检查API密钥是否正确
- 确认AI服务是否正常
- 检查网络连接

### 4. 钉钉通知失败

**问题**: `关键词不匹配`

**解决**:
- 确保配置文件中的`secret`字段与钉钉机器人的关键词一致
- 检查Webhook URL是否正确

### 5. Web界面无法访问

**问题**: `403 Forbidden`

**解决**:
```bash
# 修改web目录权限
sudo chmod -R 755 web/

# 重启服务
docker-compose restart stock-web
```

### 6. 端口冲突

**问题**: `address already in use`

**解决**:
```bash
# 查找占用端口的进程
sudo lsof -i :8090
sudo lsof -i :9090

# 修改docker-compose.yml中的端口映射
# 或停止占用端口的程序
```

---

## 🔧 开发指南

### 项目结构

```
nofx-stock/
├── main_stock.go           # 主程序入口
├── config_stock.json       # 配置文件
├── go.mod                  # Go模块定义
├── go.sum                  # 依赖校验
├── Dockerfile              # Docker镜像
├── docker-compose.yml      # Docker编排
├── nginx.conf              # Nginx配置
├── docker-start.sh         # 启动脚本(Linux)
├── docker-start.bat        # 启动脚本(Windows)
│
├── api/                    # API服务
│   └── stock_server.go     # API路由和处理
│
├── config/                 # 配置管理
│   └── stock_config.go     # 配置加载
│
├── stock/                  # 股票分析核心
│   ├── tdx_client.go       # TDX API客户端
│   ├── analyzer.go         # 分析引擎
│   └── ai_parser.go        # AI响应解析
│
├── notifier/               # 通知系统
│   └── webhook.go          # Webhook通知
│
├── mcp/                    # AI通信
│   └── client.go           # AI API客户端
│
├── web/                    # Web前端
│   └── stock_dashboard.html
│
└── logs/                   # 日志目录
    └── *.log
```

### 添加新功能

1. **添加新的技术指标**

编辑 `stock/analyzer.go`：

```go
// 在getTechnicalData方法中添加
if len(dayKline.List) >= 50 {
    // 计算MACD或其他指标
    macd := calculateMACD(dayKline.List)
    data["macd"] = macd
}
```

2. **添加新的通知渠道**

在 `notifier/webhook.go` 中实现新的Notifier接口：

```go
type WeChatNotifier struct {
    WebhookURL string
}

func (w *WeChatNotifier) SendSignal(signal *TradingSignal) error {
    // 实现企业微信通知
}
```

3. **自定义AI提示词**

编辑 `stock/analyzer.go` 的 `buildPrompt` 方法。

---

## 🔄 更新日志

### v2.0.0 (2025-11-04)

#### 🎉 重大更新
- ✨ 完全重构为股票分析系统
- 🚀 移除所有加密货币相关代码
- 🐳 优化Docker部署方案

#### ✨ 新功能
- ➕ 集成TDX股票数据API
- ➕ 支持DeepSeek/Qwen AI分析
- ➕ 钉钉/飞书Webhook通知
- ➕ Web监控界面
- ➕ RESTful API接口
- ➕ 多股票并发监控
- ➕ 技术指标计算（MA/RSI/波动率）

#### 🔧 优化
- ⚡ 使用国内镜像源加速构建
- ⚡ Host网络模式提升性能
- ⚡ 精简Docker镜像（20-30MB）
- ⚡ 优化AI提示词
- ⚡ 改进错误处理

#### 🐛 修复
- 🔨 修复Amount字段类型问题（int64→float64）
- 🔨 修复Docker网络配置
- 🔨 修复Web文件权限问题
- 🔨 修复钉钉关键词验证

#### 📝 文档
- 📚 全新README
- 📚 Docker部署指南
- 📚 快速开始指南
- 📚 API文档
- 📚 常见问题解答

---

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

---

## 🤝 贡献

欢迎贡献代码！请遵循以下步骤：

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

---

## 💬 支持

- **问题反馈**: [GitHub Issues](https://github.com/your-repo/nofx-stock/issues)
- **功能建议**: [GitHub Discussions](https://github.com/your-repo/nofx-stock/discussions)

---

## ⚠️ 免责声明

本系统提供的分析结果仅供参考，不构成投资建议。

- ❌ 不保证分析准确性
- ❌ 不承担投资损失责任
- ❌ AI分析存在局限性
- ✅ 请独立思考，谨慎决策
- ✅ 投资有风险，入市需谨慎

---

## 🌟 Star History

如果这个项目对你有帮助，请给个 ⭐️ Star！

---

**Made with ❤️ by NOFX Team**
