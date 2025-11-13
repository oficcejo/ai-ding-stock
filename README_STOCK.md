# 📈 AI股票实时分析系统

> 基于原NOFX加密货币交易系统改造，专注于A股市场的实时分析与信号通知

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 🎯 项目简介

这是一个**AI驱动的股票实时分析系统**，将原NOFX系统的AI分析能力适配到A股市场，提供：

- ✅ **实时股票分析** - 基于TDX行情数据进行深度技术分析
- ✅ **AI决策引擎** - 使用DeepSeek/Qwen等大模型进行智能分析
- ✅ **买卖信号推送** - 钉钉/飞书webhook即时通知
- ✅ **技术指标计算** - MA、RSI、波动率等多维度分析
- ✅ **Web监控面板** - 实时查看分析结果和信号

---

## 🚀 快速开始

### 前置条件

1. **Go 1.23+** 环境
2. **TDX股票数据API服务** (参考 `API_接口文档.md`)
3. **AI API密钥** (DeepSeek/Qwen/自定义)
4. **通知webhook** (钉钉/飞书，可选)

### 安装步骤

#### 1. 克隆项目

```bash
cd F:\webtest\nofx-stock
```

#### 2. 安装依赖

```bash
go mod download
```

#### 3. 配置系统

复制配置模板并编辑：

```bash
copy config_stock.json.example config_stock.json
notepad config_stock.json
```

**配置说明**:

```json
{
  "tdx_api_url": "http://localhost:8181",  // TDX API地址
  "ai_config": {
    "provider": "deepseek",                 // AI提供商: deepseek/qwen/custom
    "deepseek_key": "sk-your-api-key"       // AI API密钥
  },
  "stocks": [
    {
      "code": "000001",                     // 股票代码
      "name": "平安银行",                    // 股票名称
      "enabled": true,                       // 是否启用
      "scan_interval_minutes": 5,           // 扫描间隔（分钟）
      "min_confidence": 70                   // 最小信心度阈值
    }
  ],
  "notification": {
    "enabled": true,                        // 是否启用通知
    "dingtalk": {
      "enabled": true,
      "webhook_url": "https://oapi.dingtalk.com/robot/send?access_token=YOUR_TOKEN"
    },
    "feishu": {
      "enabled": false,
      "webhook_url": "https://open.feishu.cn/open-apis/bot/v2/hook/YOUR_TOKEN"
    }
  },
  "api_server_port": 9090,                  // API服务端口
  "log_dir": "stock_analysis_logs"          // 日志目录
}
```

#### 4. 启动系统

```bash
# 编译程序
go build -o stock_analyzer main_stock.go

# 运行
.\stock_analyzer.exe
```

或直接运行：

```bash
go run main_stock.go
```

#### 5. 访问监控面板

打开浏览器访问:

```
http://localhost:9090
```

或打开Web面板：

```
web\stock_dashboard.html
```

---

## 📊 系统架构

```
stock-analyzer/
├── main_stock.go              # 股票分析主程序
├── config_stock.json          # 股票分析配置
│
├── stock/                     # 股票分析模块
│   ├── tdx_client.go          # TDX API客户端
│   └── analyzer.go            # 股票分析器核心
│
├── notifier/                  # 通知模块
│   └── webhook.go             # 钉钉/飞书通知器
│
├── config/                    # 配置模块
│   └── stock_config.go        # 配置加载器
│
├── api/                       # API服务
│   └── stock_server.go        # HTTP API服务器
│
├── mcp/                       # AI通信层（复用）
│   └── client.go              # AI API客户端
│
└── web/                       # Web界面
    └── stock_dashboard.html   # 监控面板
```

---

## 🧠 AI分析流程

```
每个扫描周期（默认5分钟）:

1. 📊 获取实时数据
   ├── 五档行情（买卖盘口）
   ├── 日K线数据（最近60天）
   ├── 30分钟K线（最近100条）
   └── 今日分时数据

2. 🔢 计算技术指标
   ├── 均线: MA5, MA10, MA20, MA60
   ├── RSI(14): 相对强弱指标
   ├── 波动率: 20日标准差
   ├── 涨跌幅: 当日涨跌
   └── 买卖盘力量对比

3. 🤖 AI深度分析
   ├── 趋势分析（上升/下降/盘整）
   ├── 量价关系分析
   ├── 盘口分析（大单流向）
   ├── 技术指标综合研判
   └── 风险收益评估

4. 📍 输出交易信号
   ├── 信号类型: BUY/SELL/HOLD
   ├── 信心度: 0-100%
   ├── 目标价格（BUY时）
   ├── 止损价格（BUY时）
   ├── 风险回报比
   └── 详细分析理由

5. 🔔 推送通知（可选）
   └── 信心度≥阈值时发送到钉钉/飞书
```

---

## 📱 通知配置

### 钉钉机器人

1. 进入钉钉群 → 群设置 → 智能群助手 → 添加机器人
2. 选择"自定义"机器人
3. 设置名称（如"AI股票分析"）
4. 复制Webhook地址
5. 填入 `config_stock.json` 的 `notification.dingtalk.webhook_url`

**通知示例**:

```markdown
🚀 BUY信号 - 平安银行(000001)

---
当前价格: 12.50元
信心度: 85%
目标价格: 13.80元
止损价格: 11.90元
风险回报比: 1:2.3

---
分析原因:
技术面：突破MA20均线，MACD金叉，RSI(14)=58处于健康区间...
```

### 飞书机器人

1. 进入飞书群 → 群设置 → 群机器人 → 添加机器人
2. 选择"自定义机器人"
3. 复制Webhook地址
4. 填入 `config_stock.json` 的 `notification.feishu.webhook_url`

---

## 🔧 API接口

### 获取所有监控股票

```
GET http://localhost:9090/api/stocks
```

### 获取最新分析结果

```
GET http://localhost:9090/api/stock/{code}/latest
```

示例: `http://localhost:9090/api/stock/000001/latest`

### 获取历史分析记录

```
GET http://localhost:9090/api/stock/{code}/history
```

### 手动触发分析

```
POST http://localhost:9090/api/stock/{code}/analyze
```

### 系统统计

```
GET http://localhost:9090/api/statistics
```

---

## 📝 配置说明

### AI配置

#### 配置文件位置

所有AI配置都在 `config_stock.json` 文件的 `ai_config` 部分：

```json
{
  "ai_config": {
    "provider": "deepseek",           // AI提供商：deepseek/qwen/custom
    "deepseek_key": "sk-xxx",         // DeepSeek API密钥
    "qwen_key": "sk-xxx",             // Qwen API密钥
    "custom_api_url": "",             // 自定义API地址
    "custom_api_key": "",             // 自定义API密钥
    "custom_model_name": ""           // 自定义模型名称
  }
}
```

#### 支持的AI提供商

系统支持三种AI提供商，通过修改 `provider` 字段切换：

##### 1️⃣ DeepSeek（默认推荐）

```json
{
  "provider": "deepseek",
  "deepseek_key": "sk-your-deepseek-api-key"
}
```

- **获取密钥**: [https://platform.deepseek.com](https://platform.deepseek.com)
- **优势**: 性价比高，中文理解能力强
- **适用场景**: A股分析、中文技术指标解读

##### 2️⃣ Qwen（通义千问）

```json
{
  "provider": "qwen",
  "qwen_key": "sk-your-qwen-api-key"
}
```

- **获取密钥**: [https://dashscope.aliyun.com](https://dashscope.aliyun.com)
- **优势**: 阿里云生态，稳定性好
- **适用场景**: 需要阿里云集成的场景

##### 3️⃣ 自定义API（Custom）

支持任何兼容OpenAI格式的API接口：

```json
{
  "provider": "custom",
  "custom_api_url": "https://api.openai.com/v1",
  "custom_api_key": "sk-your-api-key",
  "custom_model_name": "gpt-4o"
}
```

**常见自定义配置示例**:

**OpenAI GPT-4**:
```json
{
  "provider": "custom",
  "custom_api_url": "https://api.openai.com/v1",
  "custom_api_key": "sk-proj-xxxxxxxxxxxxx",
  "custom_model_name": "gpt-4o"
}
```

**Claude (通过兼容接口)**:
```json
{
  "provider": "custom",
  "custom_api_url": "https://api.anthropic.com/v1",
  "custom_api_key": "sk-ant-xxxxxxxxxxxxx",
  "custom_model_name": "claude-3-5-sonnet-20241022"
}
```

**本地模型（Ollama）**:
```json
{
  "provider": "custom",
  "custom_api_url": "http://localhost:11434/v1",
  "custom_api_key": "ollama",
  "custom_model_name": "llama3"
}
```

**Azure OpenAI**:
```json
{
  "provider": "custom",
  "custom_api_url": "https://your-resource.openai.azure.com/openai/deployments/your-deployment",
  "custom_api_key": "your-azure-api-key",
  "custom_model_name": "gpt-4"
}
```

**国内中转API**:
```json
{
  "provider": "custom",
  "custom_api_url": "https://api.your-proxy.com/v1",
  "custom_api_key": "sk-xxxxxxxxxxxxx",
  "custom_model_name": "gpt-4o"
}
```

#### 配置字段说明

| 字段 | 说明 | 示例 |
|-----|------|------|
| `provider` | AI提供商 | `deepseek`, `qwen`, `custom` |
| `deepseek_key` | DeepSeek API密钥 | `sk-xxx` |
| `qwen_key` | Qwen API密钥 | `sk-xxx` |
| `custom_api_url` | 自定义API基础地址（不含 `/chat/completions`） | `https://api.openai.com/v1` |
| `custom_api_key` | 自定义API密钥 | `sk-proj-xxx` |
| `custom_model_name` | 自定义模型名称 | `gpt-4o`, `claude-3-5-sonnet` |

#### 切换AI提供商步骤

1. 打开 `config_stock.json` 文件
2. 修改 `ai_config.provider` 字段为目标提供商
3. 填写对应的API密钥和配置
4. 保存文件
5. 重启程序（`stock_analyzer.exe`）

启动时会显示当前使用的AI提供商：
```
✓ AI客户端已初始化 (DEEPSEEK)
✓ AI客户端已初始化 (QWEN)
✓ AI客户端已初始化 (CUSTOM)
```

#### 注意事项

- ⚠️ 自定义API必须兼容OpenAI的 `/v1/chat/completions` 接口格式
- ⚠️ `custom_api_url` 应该是基础URL，程序会自动拼接 `/chat/completions`
- ⚠️ 确保 `custom_model_name` 是API支持的有效模型名
- ⚠️ 修改配置后必须重启程序才能生效
- ⚠️ 确保服务器能访问到API地址（检查防火墙和网络）

### 股票配置

| 字段 | 说明 | 默认值 |
|-----|------|--------|
| `code` | 股票代码 | 必填 |
| `name` | 股票名称 | 必填 |
| `enabled` | 是否启用 | `true` |
| `scan_interval_minutes` | 扫描间隔 | `5`分钟 |
| `min_confidence` | 最小信心阈值 | `70`% |

### 通知配置

| 字段 | 说明 | 默认值 |
|-----|------|--------|
| `enabled` | 是否启用通知 | `true` |
| `dingtalk.enabled` | 钉钉通知开关 | `true` |
| `dingtalk.webhook_url` | 钉钉Webhook | 必填 |
| `feishu.enabled` | 飞书通知开关 | `false` |
| `feishu.webhook_url` | 飞书Webhook | 可选 |

---

## ⚠️ 风险提示

1. **本系统仅供学习研究使用**，AI分析结果不构成投资建议
2. **股票投资有风险**，请根据自身风险承受能力谨慎决策
3. **建议纸面模拟测试**后再考虑实盘应用
4. **技术指标有滞后性**，市场变化可能超出AI预期
5. **系统依赖数据源稳定性**，数据异常可能影响分析质量

---

## 🔍 常见问题

### 1. TDX API连接失败

**问题**: `获取行情失败: 请求失败`

**解决**:
- 检查TDX API服务是否启动
- 确认 `tdx_api_url` 配置正确
- 测试API连接: `curl http://localhost:8181/api/quote?code=000001`

### 2. AI分析超时

**问题**: `AI分析失败: timeout`

**解决**:
- 检查网络连接
- 确认AI API密钥有效
- 增加AI客户端超时时间（修改 `mcp/client.go`）

### 3. 通知发送失败

**问题**: `发送通知失败: xxx`

**解决**:
- 检查webhook地址是否正确
- 测试webhook: `curl -X POST <webhook_url> -H "Content-Type: application/json" -d '{"msg_type":"text","content":{"text":"test"}}'`
- 查看钉钉/飞书机器人是否被禁用

### 4. 非交易时间分析

**问题**: 非交易时间无法获取分时数据

**解决**: 系统会自动跳过分时数据获取失败，仅使用K线数据分析

---

## 📚 技术指标说明

### MA均线
- **MA5**: 5日均线（短期趋势）
- **MA10**: 10日均线
- **MA20**: 20日均线（月线）
- **MA60**: 60日均线（季线）

### RSI相对强弱指标
- **< 30**: 超卖区（可能反弹）
- **30-70**: 正常区间
- **> 70**: 超买区（可能回调）

### 波动率
- 反映股价波动剧烈程度
- 高波动率 = 高风险高收益
- 低波动率 = 稳定但收益有限

---

## 🛠️ 开发计划

- [ ] 增加更多技术指标（MACD、KDJ、布林带）
- [ ] 支持自定义分析策略
- [ ] 历史回测功能
- [ ] 多股票组合分析
- [ ] 微信、邮件通知
- [ ] 实时推送（WebSocket）
- [ ] 移动端适配

---

## 📄 许可证

MIT License

---

## 🤝 贡献

欢迎提交Issue和Pull Request！

---

## 📞 联系方式

- 项目原作者: [@nofx_ai](https://x.com/nofx_ai)
- 股票版改造: 基于NOFX开源项目

---

**Happy Trading! 📈**

**⚠️ 投资有风险，入市需谨慎！本系统仅供学习研究，不构成投资建议！**

