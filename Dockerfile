# ================================
# 股票分析系统 Dockerfile
# ================================

# 构建阶段
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 配置国内镜像源加速
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk add --no-cache git gcc musl-dev

# 配置Go模块代理
ENV GOPROXY=https://goproxy.cn,direct \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 复制go.mod和go.sum
COPY go.mod go.sum ./

# 下载依赖（允许部分失败，因为某些crypto依赖不需要）
RUN go mod download || echo "Some dependencies failed to download, continuing..."

# 复制源代码
COPY . .

# 构建应用（只构建股票分析系统，使用-mod=mod允许修改模块）
RUN go build -mod=mod -ldflags="-s -w" -o stock_analyzer main_stock.go

# 运行阶段
FROM alpine:3.19

# 配置国内镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update && \
    apk add --no-cache ca-certificates tzdata

# 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 创建非root用户
RUN addgroup -g 1000 stockapp && \
    adduser -D -u 1000 -G stockapp stockapp

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=builder /app/stock_analyzer .

# 复制Web前端文件
COPY --from=builder /app/web ./web

# 复制启动脚本
COPY docker-entrypoint.sh /app/docker-entrypoint.sh

# 创建必要的目录和示例配置文件
RUN mkdir -p /app/stock_analysis_logs && \
    echo '{"tdx_api_url":"http://192.168.1.222:8181","ai_config":{"provider":"deepseek","deepseek_key":"","qwen_key":"","custom_api_url":"","custom_api_key":"","custom_model_name":""},"stocks":[],"notification":{"enabled":false,"dingtalk":{"enabled":false,"webhook_url":"","secret":""},"feishu":{"enabled":false,"webhook_url":"","secret":""}},"trading_time":{"enable_check":true,"trading_hours":["09:30-11:30","13:00-15:00"],"timezone":"Asia/Shanghai"},"api_server_port":9090,"log_dir":"stock_analysis_logs"}' > /app/config_stock.json.example && \
    chmod +x /app/docker-entrypoint.sh && \
    chown -R stockapp:stockapp /app

# 切换到非root用户
USER stockapp

# 暴露端口
EXPOSE 9090

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9090/api/stocks || exit 1

# 使用启动脚本
ENTRYPOINT ["/app/docker-entrypoint.sh"]

