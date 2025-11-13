# 📡 TDX股票数据API接口文档

## 🌐 基础信息

**Base URL**: `http://your-server:8181`  
**Content-Type**: `application/json; charset=utf-8`  
**编码**: UTF-8

---

## 📋 响应格式

所有接口统一返回格式：

```json
{
  "code": 0,           // 0=成功, -1=失败
  "message": "success", // 提示信息
  "data": {}           // 数据内容
}
```

---

## 📊 API接口列表

### 1. 获取五档行情

**接口**: `GET /api/quote`

**描述**: 获取股票实时五档买卖盘口数据

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| code | string | 是 | 股票代码（如：000001）支持多个，逗号分隔 |

**请求示例**:
```
GET /api/quote?code=000001
GET /api/quote?code=000001,600519
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "Exchange": 0,
      "Code": "000001",
      "Active1": 2843,
      "K": {
        "Last": 12250,    // 昨收价（厘）
        "Open": 12300,    // 开盘价（厘）
        "High": 12600,    // 最高价（厘）
        "Low": 12280,     // 最低价（厘）
        "Close": 12500    // 收盘价/最新价（厘）
      },
      "ServerTime": "1730617200",
      "TotalHand": 1235000,    // 总手数
      "Intuition": 100,        // 现量
      "Amount": 156000000,     // 成交额
      "InsideDish": 520000,    // 内盘
      "OuterDisc": 715000,     // 外盘
      "BuyLevel": [            // 买五档
        {
          "Buy": true,
          "Price": 12500,      // 买一价（厘）
          "Number": 35000      // 挂单量（股）
        },
        // ... 买二到买五
      ],
      "SellLevel": [           // 卖五档
        {
          "Buy": false,
          "Price": 12510,      // 卖一价（厘）
          "Number": 30000      // 挂单量（股）
        },
        // ... 卖二到卖五
      ],
      "Rate": 0.0,
      "Active2": 2843
    }
  ]
}
```

**数据说明**:
- 价格单位：厘（1元 = 1000厘）
- 成交量单位：手（1手 = 100股）
- 挂单量单位：股

---

### 2. 获取K线数据

**接口**: `GET /api/kline`

**描述**: 获取股票K线数据（OHLC + 成交量成交额）

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| code | string | 是 | 股票代码（如：000001） |
| type | string | 否 | K线类型，默认day |

**K线类型(type)**:
- `minute1` - 1分钟K线（最多24000条）
- `minute5` - 5分钟K线
- `minute15` - 15分钟K线
- `minute30` - 30分钟K线
- `hour` - 60分钟/小时K线
- `day` - 日K线（默认）
- `week` - 周K线
- `month` - 月K线

**请求示例**:
```
GET /api/kline?code=000001&type=day
GET /api/kline?code=600519&type=minute30
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "Count": 100,
    "List": [
      {
        "Last": 12250,      // 昨收价（厘）
        "Open": 12300,      // 开盘价（厘）
        "High": 12600,      // 最高价（厘）
        "Low": 12280,       // 最低价（厘）
        "Close": 12500,     // 收盘价（厘）
        "Volume": 1235000,  // 成交量（手）
        "Amount": 156000000,// 成交额（厘）
        "Time": "2024-11-03T00:00:00Z",
        "UpCount": 0,       // 上涨数（指数有效）
        "DownCount": 0      // 下跌数（指数有效）
      }
      // ... 更多K线数据
    ]
  }
}
```

**数据说明**:
- 数据按时间倒序排列（最新的在前）
- 价格单位：厘
- 成交量单位：手
- 成交额单位：厘

---

### 3. 获取分时数据

**接口**: `GET /api/minute`

**描述**: 获取股票分时走势数据

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| code | string | 是 | 股票代码（如：000001） |
| date | string | 否 | 日期（YYYYMMDD格式），默认当天 |

**请求示例**:
```
GET /api/minute?code=000001
GET /api/minute?code=000001&date=20241103
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "Count": 240,
    "List": [
      {
        "Time": "09:31",
        "Price": 12300,    // 价格（厘）
        "Number": 1500     // 成交量（手）
      },
      {
        "Time": "09:32",
        "Price": 12310,
        "Number": 1200
      }
      // ... 240个数据点（9:30-11:30, 13:00-15:00）
    ]
  }
}
```

**数据说明**:
- 交易时段：9:30-11:30（120分钟）, 13:00-15:00（120分钟）
- 共240个数据点
- 价格单位：厘

---

### 4. 获取分时成交

**接口**: `GET /api/trade`

**描述**: 获取股票逐笔成交明细

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| code | string | 是 | 股票代码（如：000001） |
| date | string | 否 | 日期（YYYYMMDD格式），默认当天 |

**请求示例**:
```
GET /api/trade?code=000001
GET /api/trade?code=000001&date=20241103
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "Count": 1800,
    "List": [
      {
        "Time": "2024-11-03T14:59:58Z",
        "Price": 12500,    // 成交价（厘）
        "Volume": 100,     // 成交量（手）
        "Status": 0,       // 0=买入, 1=卖出, 2=中性
        "Number": 5        // 成交单数
      },
      {
        "Time": "2024-11-03T14:59:55Z",
        "Price": 12490,
        "Volume": 50,
        "Status": 1,
        "Number": 3
      }
      // ... 更多成交记录
    ]
  }
}
```

**数据说明**:
- Status: 0=主动买入(红色), 1=主动卖出(绿色), 2=中性
- 当日最多返回1800条
- 历史日期最多返回2000条
- 价格单位：厘
- 成交量单位：手

---

### 5. 搜索股票代码

**接口**: `GET /api/search`

**描述**: 根据关键词搜索股票代码和名称

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| keyword | string | 是 | 搜索关键词（代码或名称） |

**请求示例**:
```
GET /api/search?keyword=平安
GET /api/search?keyword=000001
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "code": "000001",
      "name": "平安银行"
    },
    {
      "code": "601318",
      "name": "中国平安"
    }
    // ... 最多50条结果
  ]
}
```

**数据说明**:
- 支持代码和名称模糊搜索
- 最多返回50条结果
- 仅返回A股（过滤指数等）

---

### 6. 获取股票综合信息

**接口**: `GET /api/stock-info`

**描述**: 一次性获取股票的多种数据（五档行情+日K线+分时）

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| code | string | 是 | 股票代码（如：000001） |

**请求示例**:
```
GET /api/stock-info?code=000001
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "quote": {
      // 五档行情数据（同/api/quote）
    },
    "kline_day": {
      // 最近30天日K线（同/api/kline?type=day）
    },
    "minute": {
      // 今日分时数据（同/api/minute）
    }
  }
}
```

**数据说明**:
- 整合了三个接口的数据
- 适合快速获取股票概览
- 减少API调用次数

---

## 🔧 扩展接口（高级功能）

### 7. 获取股票列表

**接口**: `GET /api/codes`

**描述**: 获取指定交易所的所有股票代码列表

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| exchange | string | 否 | 交易所代码，默认all |

**交易所代码**:
- `sh` - 上海证券交易所
- `sz` - 深圳证券交易所
- `bj` - 北京证券交易所
- `all` - 全部（默认）

**请求示例**:
```
GET /api/codes
GET /api/codes?exchange=sh
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "total": 5234,
    "exchanges": {
      "sh": 2156,
      "sz": 2845,
      "bj": 233
    },
    "codes": [
      {
        "code": "000001",
        "name": "平安银行",
        "exchange": "sz"
      }
      // ... 更多股票
    ]
  }
}
```

---

### 8. 批量获取行情

**接口**: `POST /api/batch-quote`

**描述**: 批量获取多只股票的实时行情

**请求参数** (JSON Body):
```json
{
  "codes": ["000001", "600519", "601318"]
}
```

**请求示例**:
```bash
curl -X POST http://localhost:8080/api/batch-quote \
  -H "Content-Type: application/json" \
  -d '{"codes":["000001","600519","601318"]}'
```

**响应示例**:
```json
{
  "code": 0,
  "message": "success",
  "data": [
    // 数组，每个元素同/api/quote的单个股票数据
  ]
}
```

---

### 9. 获取历史K线

**接口**: `GET /api/kline-history`

**描述**: 获取指定时间范围的K线数据

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| code | string | 是 | 股票代码 |
| type | string | 是 | K线类型 |
| start_date | string | 否 | 开始日期（YYYYMMDD） |
| end_date | string | 否 | 结束日期（YYYYMMDD） |
| limit | int | 否 | 返回条数，默认100，最大800 |

**请求示例**:
```
GET /api/kline-history?code=000001&type=day&limit=30
GET /api/kline-history?code=000001&type=day&start_date=20241001&end_date=20241101
```

---

### 10. 获取指数数据

**接口**: `GET /api/index`

**描述**: 获取指数K线数据（如上证指数、深证成指）

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|-----|------|------|------|
| code | string | 是 | 指数代码（如：sh000001） |
| type | string | 否 | K线类型，默认day |

**常用指数代码**:
- `sh000001` - 上证指数
- `sz399001` - 深证成指
- `sz399006` - 创业板指
- `sh000300` - 沪深300

**请求示例**:
```
GET /api/index?code=sh000001&type=day
```

---

## 💡 使用示例

### Python示例

```python
import requests

BASE_URL = "http://your-server:8080"

# 1. 获取五档行情
def get_quote(code):
    url = f"{BASE_URL}/api/quote?code={code}"
    response = requests.get(url)
    data = response.json()
    if data['code'] == 0:
        return data['data']
    return None

# 2. 获取日K线
def get_kline(code, type='day'):
    url = f"{BASE_URL}/api/kline?code={code}&type={type}"
    response = requests.get(url)
    data = response.json()
    if data['code'] == 0:
        return data['data']['List']
    return None

# 3. 搜索股票
def search_stock(keyword):
    url = f"{BASE_URL}/api/search?keyword={keyword}"
    response = requests.get(url)
    data = response.json()
    if data['code'] == 0:
        return data['data']
    return None

# 使用示例
if __name__ == "__main__":
    # 搜索股票
    stocks = search_stock("平安")
    print(f"搜索结果: {stocks}")
    
    # 获取行情
    quote = get_quote("000001")
    print(f"最新价: {quote[0]['K']['Close'] / 1000}元")
    
    # 获取K线
    klines = get_kline("000001", "day")
    print(f"获取到{len(klines)}条K线数据")
```

### JavaScript示例

```javascript
const BASE_URL = 'http://your-server:8080';

// 1. 获取五档行情
async function getQuote(code) {
    const response = await fetch(`${BASE_URL}/api/quote?code=${code}`);
    const data = await response.json();
    if (data.code === 0) {
        return data.data;
    }
    return null;
}

// 2. 获取K线
async function getKline(code, type = 'day') {
    const response = await fetch(`${BASE_URL}/api/kline?code=${code}&type=${type}`);
    const data = await response.json();
    if (data.code === 0) {
        return data.data.List;
    }
    return null;
}

// 3. 批量获取行情
async function batchGetQuote(codes) {
    const response = await fetch(`${BASE_URL}/api/batch-quote`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ codes })
    });
    const data = await response.json();
    return data.data;
}

// 使用示例
(async () => {
    // 获取行情
    const quote = await getQuote('000001');
    console.log('最新价:', quote[0].K.Close / 1000);
    
    // 获取K线
    const klines = await getKline('000001', 'day');
    console.log('K线数据量:', klines.length);
    
    // 批量获取
    const quotes = await batchGetQuote(['000001', '600519', '601318']);
    console.log('批量行情:', quotes.length);
})();
```

### cURL示例

```bash
# 1. 获取五档行情
curl "http://localhost:8080/api/quote?code=000001"

# 2. 获取日K线
curl "http://localhost:8080/api/kline?code=000001&type=day"

# 3. 获取分时数据
curl "http://localhost:8080/api/minute?code=000001"

# 4. 搜索股票
curl "http://localhost:8080/api/search?keyword=平安"

# 5. 批量获取行情
curl -X POST http://localhost:8080/api/batch-quote \
  -H "Content-Type: application/json" \
  -d '{"codes":["000001","600519"]}'
```

---

## 🔒 错误码说明

| code | message | 说明 |
|------|---------|------|
| 0 | success | 请求成功 |
| -1 | 股票代码不能为空 | 缺少必填参数code |
| -1 | 获取行情失败: xxx | 数据获取失败，xxx为具体错误 |
| -1 | 获取K线失败: xxx | K线数据获取失败 |
| -1 | 未找到相关股票 | 搜索无结果 |
| -1 | 搜索关键词不能为空 | 缺少keyword参数 |

---

## 📊 数据单位换算

### 价格单位
- **返回值**：厘（1元 = 1000厘）
- **换算公式**：元 = 厘 / 1000
- **示例**：12500厘 = 12.50元

### 成交量单位
- **返回值**：手（1手 = 100股）
- **换算公式**：股 = 手 × 100
- **示例**：1235手 = 123500股

### 成交额单位
- **返回值**：厘
- **换算公式**：元 = 厘 / 1000
- **示例**：156000000厘 = 156000元 = 15.6万元

---

## 🚀 性能建议

1. **批量请求**：使用批量接口代替多次单个请求
2. **缓存**：对不常变化的数据（如股票列表）做本地缓存
3. **限流**：避免频繁请求，建议间隔>=3秒
4. **压缩**：使用gzip压缩减少传输量

---

## 📝 更新日志

### v1.0.0 (2024-11-03)
- ✅ 实现基础6个API接口
- ✅ 统一响应格式
- ✅ 完整文档和示例

### v1.1.0 (计划中)
- 🔄 批量查询接口
- 🔄 历史K线范围查询
- 🔄 指数数据接口
- 🔄 WebSocket实时推送

---

## 📞 技术支持

- 文档地址：本文件
- API测试：使用Postman或cURL
- 问题反馈：GitHub Issues

---

**Happy Coding!** 🎉

