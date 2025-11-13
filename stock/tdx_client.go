package stock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TDXClient TDX股票数据API客户端
type TDXClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewTDXClient 创建新的TDX客户端
func NewTDXClient(baseURL string) *TDXClient {
	return &TDXClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// APIResponse 统一API响应格式
type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// QuoteData 五档行情数据
type QuoteData struct {
	Exchange   int     `json:"Exchange"`
	Code       string  `json:"Code"`
	Active1    int     `json:"Active1"`
	K          KData   `json:"K"`
	ServerTime string  `json:"ServerTime"`
	TotalHand  int64   `json:"TotalHand"`  // 总手数
	Intuition  int     `json:"Intuition"`  // 现量
	Amount     float64 `json:"Amount"`     // 成交额（厘，可能是浮点数）
	InsideDish int64   `json:"InsideDish"` // 内盘
	OuterDisc  int64   `json:"OuterDisc"`  // 外盘
	BuyLevel   []Level `json:"BuyLevel"`   // 买五档
	SellLevel  []Level `json:"SellLevel"`  // 卖五档
	Rate       float64 `json:"Rate"`
	Active2    int     `json:"Active2"`
}

// KData K线基础数据
type KData struct {
	Last  int `json:"Last"`  // 昨收价（厘）
	Open  int `json:"Open"`  // 开盘价（厘）
	High  int `json:"High"`  // 最高价（厘）
	Low   int `json:"Low"`   // 最低价（厘）
	Close int `json:"Close"` // 收盘价/最新价（厘）
}

// Level 盘口档位
type Level struct {
	Buy    bool `json:"Buy"`
	Price  int  `json:"Price"`  // 价格（厘）
	Number int  `json:"Number"` // 挂单量（股）
}

// KlineData K线数据
type KlineData struct {
	Count int         `json:"Count"`
	List  []KlineItem `json:"List"`
}

// KlineItem K线单条数据
type KlineItem struct {
	Last      int       `json:"Last"`   // 昨收价（厘）
	Open      int       `json:"Open"`   // 开盘价（厘）
	High      int       `json:"High"`   // 最高价（厘）
	Low       int       `json:"Low"`    // 最低价（厘）
	Close     int       `json:"Close"`  // 收盘价（厘）
	Volume    int64     `json:"Volume"` // 成交量（手）
	Amount    float64   `json:"Amount"` // 成交额（厘，可能是浮点数）
	Time      time.Time `json:"Time"`
	UpCount   int       `json:"UpCount"`   // 上涨数
	DownCount int       `json:"DownCount"` // 下跌数
}

// MinuteData 分时数据
type MinuteData struct {
	Count int          `json:"Count"`
	List  []MinuteItem `json:"List"`
}

// MinuteItem 分时单条数据
type MinuteItem struct {
	Time   string `json:"Time"`   // 时间（HH:MM）
	Price  int    `json:"Price"`  // 价格（厘）
	Number int    `json:"Number"` // 成交量（手）
}

// SearchResult 搜索结果
type SearchResult struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// GetQuote 获取五档行情
func (c *TDXClient) GetQuote(code string) (*QuoteData, error) {
	url := fmt.Sprintf("%s/api/quote?code=%s", c.BaseURL, code)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if apiResp.Code != 0 {
		return nil, fmt.Errorf("API错误: %s", apiResp.Message)
	}

	var quotes []QuoteData
	if err := json.Unmarshal(apiResp.Data, &quotes); err != nil {
		return nil, fmt.Errorf("解析行情数据失败: %w", err)
	}

	if len(quotes) == 0 {
		return nil, fmt.Errorf("未获取到行情数据")
	}

	return &quotes[0], nil
}

// GetKline 获取K线数据
// adjust参数: 0=不复权(默认), 1=前复权, 2=后复权
// 为了与实时行情价格一致，默认使用不复权数据(adjust=0)
func (c *TDXClient) GetKline(code string, klineType string, limit int) (*KlineData, error) {
	url := fmt.Sprintf("%s/api/kline?code=%s&type=%s&adjust=0", c.BaseURL, code, klineType)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if apiResp.Code != 0 {
		return nil, fmt.Errorf("API错误: %s", apiResp.Message)
	}

	var klineData KlineData
	if err := json.Unmarshal(apiResp.Data, &klineData); err != nil {
		return nil, fmt.Errorf("解析K线数据失败: %w", err)
	}

	// 限制返回数量（取最近的limit条，而不是最旧的limit条）
	if limit > 0 && len(klineData.List) > limit {
		klineData.List = klineData.List[len(klineData.List)-limit:]
		klineData.Count = limit
	}

	return &klineData, nil
}

// GetMinute 获取分时数据
func (c *TDXClient) GetMinute(code string, date string) (*MinuteData, error) {
	urlStr := fmt.Sprintf("%s/api/minute?code=%s", c.BaseURL, code)
	if date != "" {
		urlStr += "&date=" + date
	}

	resp, err := c.HTTPClient.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if apiResp.Code != 0 {
		return nil, fmt.Errorf("API错误: %s", apiResp.Message)
	}

	var minuteData MinuteData
	if err := json.Unmarshal(apiResp.Data, &minuteData); err != nil {
		return nil, fmt.Errorf("解析分时数据失败: %w", err)
	}

	return &minuteData, nil
}

// SearchStock 搜索股票
func (c *TDXClient) SearchStock(keyword string) ([]SearchResult, error) {
	urlStr := fmt.Sprintf("%s/api/search?keyword=%s", c.BaseURL, url.QueryEscape(keyword))
	resp, err := c.HTTPClient.Get(urlStr)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if apiResp.Code != 0 {
		return nil, fmt.Errorf("API错误: %s", apiResp.Message)
	}

	var results []SearchResult
	if err := json.Unmarshal(apiResp.Data, &results); err != nil {
		return nil, fmt.Errorf("解析搜索结果失败: %w", err)
	}

	return results, nil
}

// BatchGetQuote 批量获取行情
func (c *TDXClient) BatchGetQuote(codes []string) ([]QuoteData, error) {
	// 使用逗号分隔的方式批量获取
	codeStr := strings.Join(codes, ",")
	url := fmt.Sprintf("%s/api/quote?code=%s", c.BaseURL, codeStr)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if apiResp.Code != 0 {
		return nil, fmt.Errorf("API错误: %s", apiResp.Message)
	}

	var quotes []QuoteData
	if err := json.Unmarshal(apiResp.Data, &quotes); err != nil {
		return nil, fmt.Errorf("解析行情数据失败: %w", err)
	}

	return quotes, nil
}

// PriceToYuan 将厘转换为元
func PriceToYuan(li int) float64 {
	return float64(li) / 1000.0
}

// VolumeToShares 将手转换为股
func VolumeToShares(hands int64) int64 {
	return hands * 100
}

// AmountToYuan 将成交额（厘）转换为元
func AmountToYuan(amount float64) float64 {
	return amount / 1000.0
}
