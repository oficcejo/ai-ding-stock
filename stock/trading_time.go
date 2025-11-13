package stock

import (
	"time"
)

// TradingTimeConfig 交易时间配置
type TradingTimeConfig struct {
	EnableTradingTimeCheck bool     `json:"enable_trading_time_check"` // 是否启用交易时间检查
	TradingHours           []string `json:"trading_hours"`             // 交易时段（如：["09:30-11:30", "13:00-15:00"]）
	Timezone               string   `json:"timezone"`                  // 时区（如：Asia/Shanghai）
}

// DefaultTradingTimeConfig 默认交易时间配置（A股）
func DefaultTradingTimeConfig() TradingTimeConfig {
	return TradingTimeConfig{
		EnableTradingTimeCheck: true,
		TradingHours: []string{
			"09:30-11:30", // 上午盘
			"13:00-15:00", // 下午盘
		},
		Timezone: "Asia/Shanghai",
	}
}

// TradingTimeChecker 交易时间检查器
type TradingTimeChecker struct {
	Config   TradingTimeConfig
	Location *time.Location
}

// NewTradingTimeChecker 创建交易时间检查器
func NewTradingTimeChecker(config TradingTimeConfig) (*TradingTimeChecker, error) {
	loc, err := time.LoadLocation(config.Timezone)
	if err != nil {
		// 如果加载时区失败，使用本地时区
		loc = time.Local
	}

	return &TradingTimeChecker{
		Config:   config,
		Location: loc,
	}, nil
}

// IsTradingDay 判断是否是交易日
func (tc *TradingTimeChecker) IsTradingDay(t time.Time) bool {
	// 转换到配置的时区
	t = t.In(tc.Location)

	// 周六日不交易
	weekday := t.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}

	// TODO: 可以添加节假日判断
	// 这里可以集成节假日API或使用本地节假日数据
	if tc.isHoliday(t) {
		return false
	}

	return true
}

// IsTradingTime 判断是否在交易时段内
func (tc *TradingTimeChecker) IsTradingTime(t time.Time) bool {
	// 如果未启用交易时间检查，总是返回true
	if !tc.Config.EnableTradingTimeCheck {
		return true
	}

	// 首先判断是否是交易日
	if !tc.IsTradingDay(t) {
		return false
	}

	// 转换到配置的时区
	t = t.In(tc.Location)

	// 获取当前时间（只看时分）
	currentTime := t.Format("15:04")

	// 检查是否在任一交易时段内
	for _, period := range tc.Config.TradingHours {
		if tc.isInTimePeriod(currentTime, period) {
			return true
		}
	}

	return false
}

// isInTimePeriod 判断时间是否在指定时段内
func (tc *TradingTimeChecker) isInTimePeriod(current string, period string) bool {
	// 解析时间段，格式：HH:MM-HH:MM
	if len(period) < 11 {
		return false
	}

	// 简单的字符串比较（因为格式固定为HH:MM）
	start := period[:5]
	end := period[6:]

	return current >= start && current <= end
}

// isHoliday 判断是否是节假日
func (tc *TradingTimeChecker) isHoliday(t time.Time) bool {
	// 这里可以实现节假日判断逻辑
	// 可以集成节假日API或使用本地数据

	// 简单示例：2025年节假日（需要更新）
	holidays := map[string]bool{
		"2025-01-01": true, // 元旦
		"2025-01-28": true, // 春节
		"2025-01-29": true, // 春节
		"2025-01-30": true, // 春节
		"2025-01-31": true, // 春节
		"2025-02-01": true, // 春节
		"2025-02-02": true, // 春节
		"2025-02-03": true, // 春节
		"2025-04-04": true, // 清明节
		"2025-04-05": true, // 清明节
		"2025-04-06": true, // 清明节
		"2025-05-01": true, // 劳动节
		"2025-05-02": true, // 劳动节
		"2025-05-03": true, // 劳动节
		"2025-06-10": true, // 端午节
		"2025-10-01": true, // 国庆节
		"2025-10-02": true, // 国庆节
		"2025-10-03": true, // 国庆节
		"2025-10-04": true, // 国庆节
		"2025-10-05": true, // 国庆节
		"2025-10-06": true, // 国庆节
		"2025-10-07": true, // 国庆节
	}

	dateStr := t.Format("2006-01-02")
	return holidays[dateStr]
}

// GetNextTradingTime 获取下一个交易时间
func (tc *TradingTimeChecker) GetNextTradingTime(t time.Time) time.Time {
	t = t.In(tc.Location)

	// 如果未启用交易时间检查，返回当前时间
	if !tc.Config.EnableTradingTimeCheck {
		return t
	}

	// 如果当前就在交易时段，返回当前时间
	if tc.IsTradingTime(t) {
		return t
	}

	// 尝试找到今天的下一个交易时段
	today := t.Format("2006-01-02")
	currentTime := t.Format("15:04")

	if tc.IsTradingDay(t) {
		for _, period := range tc.Config.TradingHours {
			start := period[:5]
			if currentTime < start {
				// 找到今天的下一个交易时段
				nextTime, _ := time.ParseInLocation("2006-01-02 15:04", today+" "+start, tc.Location)
				return nextTime
			}
		}
	}

	// 如果今天没有下一个交易时段，找下一个交易日的第一个时段
	nextDay := t.AddDate(0, 0, 1)
	for {
		if tc.IsTradingDay(nextDay) {
			// 返回下一个交易日的第一个交易时段开始时间
			if len(tc.Config.TradingHours) > 0 {
				dateStr := nextDay.Format("2006-01-02")
				start := tc.Config.TradingHours[0][:5]
				nextTime, _ := time.ParseInLocation("2006-01-02 15:04", dateStr+" "+start, tc.Location)
				return nextTime
			}
		}
		nextDay = nextDay.AddDate(0, 0, 1)

		// 防止无限循环，最多查找30天
		if nextDay.Sub(t) > 30*24*time.Hour {
			return t.Add(24 * time.Hour)
		}
	}
}

// GetTradingTimeStatus 获取交易时间状态信息
func (tc *TradingTimeChecker) GetTradingTimeStatus(t time.Time) map[string]interface{} {
	t = t.In(tc.Location)

	status := map[string]interface{}{
		"current_time":    t.Format("2006-01-02 15:04:05"),
		"is_trading_day":  tc.IsTradingDay(t),
		"is_trading_time": tc.IsTradingTime(t),
		"weekday":         t.Weekday().String(),
		"timezone":        tc.Config.Timezone,
		"check_enabled":   tc.Config.EnableTradingTimeCheck,
	}

	if !tc.IsTradingTime(t) {
		nextTime := tc.GetNextTradingTime(t)
		status["next_trading_time"] = nextTime.Format("2006-01-02 15:04:05")
		duration := nextTime.Sub(t)
		status["wait_duration"] = duration.String()
	}

	return status
}
