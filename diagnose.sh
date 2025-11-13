#!/bin/bash

# ================================
# 股票分析系统诊断脚本
# ================================

echo "========================================"
echo "  股票分析系统诊断工具"
echo "========================================"
echo ""

# 检查权限
if [ "$EUID" -ne 0 ]; then 
    echo "⚠️  请使用sudo运行此脚本"
    echo "   sudo bash diagnose.sh"
    exit 1
fi

echo "📋 诊断信息收集中..."
echo ""

# 1. 检查容器状态
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "1️⃣  容器运行状态"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker compose ps
echo ""

# 2. 查看最近的MA值
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "2️⃣  日志中的MA值"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker compose logs --tail=200 stock-analyzer | grep -E "MA5|MA10|MA20|当前价格" | tail -10
echo ""

# 3. 测试TDX API
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "3️⃣  TDX API返回的最近5天数据"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 测试adjust=0参数
echo "📊 测试 adjust=0 (不复权):"
curl -s "http://192.168.1.222:8181/api/kline?code=300438&type=day&adjust=0" 2>/dev/null | \
    jq -r '.data.List[-5:] | .[] | "  \(.Time[:10]) 收盘: \(.Close/1000 | tonumber | . * 100 | round / 100)元"' 2>/dev/null || \
    echo "  ⚠️  jq未安装，无法解析JSON"

echo ""

# 测试默认参数
echo "📊 测试默认参数（无adjust）:"
curl -s "http://192.168.1.222:8181/api/kline?code=300438&type=day" 2>/dev/null | \
    jq -r '.data.List[-5:] | .[] | "  \(.Time[:10]) 收盘: \(.Close/1000 | tonumber | . * 100 | round / 100)元"' 2>/dev/null || \
    echo "  ⚠️  jq未安装，无法解析JSON"

echo ""

# 4. 检查镜像构建时间
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "4️⃣  Docker镜像信息"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker images | grep -E "REPOSITORY|nofx-stock-stock-analyzer"
echo ""

# 5. 查看最近的完整日志
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "5️⃣  最近的分析日志（最后30行）"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
docker compose logs --tail=30 stock-analyzer
echo ""

# 6. 手动计算MA5
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "6️⃣  手动计算MA5验证"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 获取最近5天收盘价
echo "📊 获取最近5天收盘价..."
PRICES=$(curl -s "http://192.168.1.222:8181/api/kline?code=300438&type=day&adjust=0" 2>/dev/null | \
    jq -r '.data.List[-5:] | .[].Close' 2>/dev/null)

if [ ! -z "$PRICES" ]; then
    echo "  最近5天收盘价（厘）："
    echo "$PRICES" | while read price; do
        yuan=$(echo "scale=2; $price / 1000" | bc)
        echo "    $price 厘 = $yuan 元"
    done
    
    # 计算平均值
    SUM=$(echo "$PRICES" | awk '{sum+=$1} END {print sum}')
    AVG=$(echo "scale=2; $SUM / 5 / 1000" | bc)
    echo ""
    echo "  📈 手动计算MA5: $AVG 元"
    echo ""
    echo "  ✅ 如果程序显示MA5接近 $AVG 元，说明正常"
    echo "  ❌ 如果程序显示MA5是 11-17 元，说明还有问题"
else
    echo "  ⚠️  无法获取数据"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ 诊断完成"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📝 分析结果："
echo "  1. 检查 '2️⃣ 日志中的MA值' 是否正常"
echo "  2. 对比 '3️⃣ TDX API数据' 和 '6️⃣ 手动计算' 的结果"
echo "  3. 如果API数据正常但程序MA异常，需要进一步排查代码"
echo ""
echo "🔧 如果问题依然存在，请将上述输出发送给开发人员"
echo ""

