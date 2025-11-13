#!/bin/bash

# NOFX Stock Analyzer - Docker 快速部署脚本

set -e

echo "🐳 NOFX Stock Analyzer - Docker 部署脚本"
echo "=========================================="
echo ""

# 检查Docker是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ 错误：未检测到 Docker，请先安装 Docker"
    exit 1
fi

# 检查Docker Compose是否安装
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "❌ 错误：未检测到 Docker Compose，请先安装 Docker Compose"
    exit 1
fi

# 检查配置文件
if [ ! -f "config_stock.json" ]; then
    echo "⚠️  配置文件不存在，创建默认配置..."
    cat > config_stock.json <<EOF
{
  "tdx_api_url": "http://192.168.1.222:8181",
  "ai_config": {
    "provider": "deepseek",
    "deepseek_key": "",
    "qwen_key": "",
    "custom_api_url": "",
    "custom_api_key": "",
    "custom_model_name": ""
  },
  "stocks": [
    {
      "code": "000001",
      "name": "平安银行",
      "enabled": true,
      "scan_interval_minutes": 5,
      "min_confidence": 70
    }
  ],
  "notification": {
    "enabled": false,
    "dingtalk": {
      "enabled": false,
      "webhook_url": "",
      "secret": ""
    },
    "feishu": {
      "enabled": false,
      "webhook_url": "",
      "secret": ""
    }
  },
  "trading_time": {
    "enable_check": true,
    "trading_hours": ["09:30-11:30", "13:00-15:00"],
    "timezone": "Asia/Shanghai"
  },
  "api_server_port": 9090,
  "log_dir": "stock_analysis_logs"
}
EOF
    echo "✅ 已创建默认配置文件: config_stock.json"
    echo "📝 请编辑配置文件或稍后通过 Web 界面修改"
    echo ""
fi

# 创建日志目录
if [ ! -d "stock_analysis_logs" ]; then
    mkdir -p stock_analysis_logs
    echo "✅ 已创建日志目录: stock_analysis_logs"
fi

# 停止旧容器
echo ""
echo "🛑 停止旧容器..."
docker-compose down 2>/dev/null || docker compose down 2>/dev/null || true

# 构建镜像
echo ""
echo "🔨 构建 Docker 镜像..."
if command -v docker-compose &> /dev/null; then
    docker-compose build
else
    docker compose build
fi

# 启动容器
echo ""
echo "🚀 启动容器..."
if command -v docker-compose &> /dev/null; then
    docker-compose up -d
else
    docker compose up -d
fi

# 等待服务启动
echo ""
echo "⏳ 等待服务启动..."
sleep 5

# 检查容器状态
echo ""
echo "📊 容器状态："
if command -v docker-compose &> /dev/null; then
    docker-compose ps
else
    docker compose ps
fi

# 显示日志
echo ""
echo "📝 最近日志："
if command -v docker-compose &> /dev/null; then
    docker-compose logs --tail=20 stock-analyzer
else
    docker compose logs --tail=20 stock-analyzer
fi

# 完成提示
echo ""
echo "=========================================="
echo "✅ 部署完成！"
echo ""
echo "🌐 Web 配置页面: http://localhost:9090"
echo "📡 API 接口: http://localhost:9090/api/stocks"
echo ""
echo "📋 常用命令："
echo "  查看日志: docker-compose logs -f stock-analyzer"
echo "  重启服务: docker-compose restart stock-analyzer"
echo "  停止服务: docker-compose down"
echo "  进入容器: docker-compose exec stock-analyzer sh"
echo ""
echo "📝 配置文件位置: ./config_stock.json"
echo "📂 日志目录: ./stock_analysis_logs"
echo ""
echo "⚠️  修改配置后需要重启容器: docker-compose restart stock-analyzer"
echo "=========================================="

