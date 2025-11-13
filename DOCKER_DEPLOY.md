# 🐳 Docker部署指南

完整的Docker容器化部署文档，适用于生产环境和开发环境。

## ✅ **重要：Web配置可以持久化保存！**

通过Docker Volume挂载，**Web页面修改的配置可以持久化保存**到宿主机，容器重启后配置不会丢失！

### 配置持久化原理

```yaml
# docker-compose.yml 中的关键配置
volumes:
  # 配置文件（可读写，支持Web保存）⚠️ 不能加 :ro
  - ./config_stock.json:/app/config_stock.json

  # 日志目录（持久化分析日志）
  - ./stock_analysis_logs:/app/stock_analysis_logs

  # Web前端文件
  - ./web:/app/web
```

### 工作流程

1. **容器启动** → 读取宿主机的 `config_stock.json`
2. **Web修改配置** → 保存到容器内的 `/app/config_stock.json`
3. **Volume同步** → Docker自动同步到宿主机的 `./config_stock.json`
4. **配置备份** → 备份文件也保存在宿主机（`config_stock.json.backup.YYYYMMDDHHMMSS`）
5. **容器重启** → 配置依然存在（从宿主机读取）

**使用步骤**：
1. 访问 `http://your-server-ip:9090` 打开Web配置页面
2. 修改AI模型、股票列表等配置
3. 点击"💾 保存配置"
4. 执行 `docker-compose restart stock-analyzer` 重启容器
5. 配置生效！✅

---

## 📋 目录

- [系统要求](#系统要求)
- [快速开始](#快速开始)
- [详细配置](#详细配置)
- [部署步骤](#部署步骤)
- [服务管理](#服务管理)
- [故障排查](#故障排查)
- [性能优化](#性能优化)
- [安全建议](#安全建议)

---

## 🔧 系统要求

### 硬件要求

| 组件 | 最低配置 | 推荐配置 |
|-----|---------|---------|
| CPU | 1核 | 2核+ |
| 内存 | 512MB | 1GB+ |
| 磁盘 | 2GB | 5GB+ |
| 网络 | 1Mbps | 10Mbps+ |

### 软件要求

- **Docker**: 20.10+ 或更高版本
- **Docker Compose**: 2.0+ 或更高版本（或内置的 `docker compose`）
- **操作系统**: 
  - Linux (推荐 Ubuntu 20.04+, CentOS 7+)
  - Windows 10/11 with Docker Desktop
  - macOS with Docker Desktop

---

## 🚀 快速开始

### 1. 安装Docker

#### Linux (Ubuntu/Debian)

```bash
# 安装Docker
curl -fsSL https://get.docker.com | bash

# 启动Docker服务
sudo systemctl start docker
sudo systemctl enable docker

# 添加当前用户到docker组（可选）
sudo usermod -aG docker $USER
newgrp docker

# 安装Docker Compose（如果需要）
sudo apt-get update
sudo apt-get install docker-compose-plugin
```

#### Windows/macOS

下载并安装 [Docker Desktop](https://www.docker.com/products/docker-desktop/)

### 2. 克隆或下载项目

```bash
git clone <your-repo-url>
cd nofx-stock
```

### 3. 配置系统

```bash
# 复制配置示例
cp config_stock.json.example config_stock.json

# 编辑配置文件
nano config_stock.json  # Linux/macOS
notepad config_stock.json  # Windows
```

**必须配置的字段**：

```json
{
  "tdx_api_url": "http://your-tdx-api:5000",
  "ai_config": {
    "deepseek": {
      "api_key": "your-deepseek-api-key",
      "enabled": true
    }
  },
  "stocks": [
    {
      "code": "600519",
      "name": "贵州茅台",
      "enabled": true
    }
  ]
}
```

### 4. 启动服务

#### Linux/macOS

```bash
# 添加执行权限
chmod +x docker-start.sh

# 启动服务
./docker-start.sh start
```

#### Windows

```cmd
docker-start.bat start
```

或直接双击 `docker-start.bat`

### 5. 访问系统

- **Web界面**: http://localhost
- **API接口**: http://localhost:8080/api/stocks

---

## ⚙️ 详细配置

### Docker Compose 配置

`docker-compose.yml` 文件结构：

```yaml
services:
  stock-analyzer:    # 后端服务
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./config_stock.json:/app/config_stock.json:ro
      - ./logs:/app/logs
    
  stock-web:         # 前端服务
    image: nginx:1.25-alpine
    ports:
      - "80:80"
    depends_on:
      - stock-analyzer
```

### 环境变量配置

创建 `.env` 文件（可选）：

```bash
cp .env.example .env
nano .env
```

支持的环境变量：

```bash
# 时区
TZ=Asia/Shanghai

# 日志级别
LOG_LEVEL=info

# 端口映射
API_PORT=8080
WEB_PORT=80

# 网络配置
SUBNET=172.20.0.0/16
```

### 卷（Volumes）配置

| 卷路径 | 宿主机路径 | 说明 | 是否必需 |
|-------|-----------|------|---------|
| `/app/config_stock.json` | `./config_stock.json` | 配置文件 | ✅ 必需 |
| `/app/logs` | `./logs` | 日志目录 | ⭐ 推荐 |
| `/etc/localtime` | `/etc/localtime` | 时区同步 | 可选 |

### 端口映射

| 容器端口 | 宿主机端口 | 服务 | 说明 |
|---------|-----------|-----|------|
| 8080 | 8080 | API服务 | 后端API接口 |
| 80 | 80 | Web服务 | 前端界面 |

如需修改宿主机端口，编辑 `docker-compose.yml`：

```yaml
ports:
  - "8888:8080"  # 将API映射到8888端口
  - "8000:80"    # 将Web映射到8000端口
```

---

## 📦 部署步骤

### 开发环境部署

```bash
# 1. 构建镜像
docker-compose build

# 2. 启动服务（前台运行，方便调试）
docker-compose up

# 3. 查看日志
docker-compose logs -f
```

### 生产环境部署

```bash
# 1. 准备配置文件
cp config_stock.json.example config_stock.json
vim config_stock.json

# 2. 创建必要目录
mkdir -p logs

# 3. 构建并启动（后台运行）
docker-compose up -d --build

# 4. 检查服务状态
docker-compose ps

# 5. 查看日志
docker-compose logs -f --tail=100
```

### 使用启动脚本（推荐）

#### Linux/macOS

```bash
# 启动
./docker-start.sh start

# 查看状态
./docker-start.sh status

# 查看日志
./docker-start.sh logs

# 停止
./docker-start.sh stop

# 重启
./docker-start.sh restart
```

#### Windows

```cmd
REM 启动
docker-start.bat start

REM 查看状态
docker-start.bat status

REM 查看日志
docker-start.bat logs

REM 停止
docker-start.bat stop
```

---

## 🎛️ 服务管理

### 启动服务

```bash
# 方式1：使用脚本
./docker-start.sh start                    # Linux/macOS
docker-start.bat start                     # Windows

# 方式2：使用docker-compose
docker-compose up -d

# 方式3：指定配置文件
docker-compose -f docker-compose.yml up -d
```

### 停止服务

```bash
# 方式1：使用脚本
./docker-start.sh stop

# 方式2：使用docker-compose
docker-compose down

# 方式3：停止但不删除容器
docker-compose stop
```

### 重启服务

```bash
# 完全重启
./docker-start.sh restart

# 仅重启特定服务
docker-compose restart stock-analyzer
docker-compose restart stock-web
```

### 查看日志

```bash
# 实时日志（所有服务）
docker-compose logs -f

# 实时日志（特定服务）
docker-compose logs -f stock-analyzer

# 最近100行日志
docker-compose logs --tail=100

# 查看文件日志
tail -f logs/stock_analyzer.log
```

### 查看状态

```bash
# 使用脚本
./docker-start.sh status

# 使用docker-compose
docker-compose ps

# 详细状态
docker ps --filter "name=stock-"
```

### 进入容器

```bash
# 进入后端容器
docker exec -it stock-analyzer sh

# 进入Web容器
docker exec -it stock-web sh

# 使用脚本
./docker-start.sh shell
```

### 更新服务

```bash
# 使用脚本
./docker-start.sh update

# 手动更新
git pull                           # 更新代码
docker-compose down                # 停止服务
docker-compose build --no-cache    # 重新构建
docker-compose up -d               # 启动服务
```

### 清理数据

```bash
# 完全清理（危险操作！）
./docker-start.sh clean

# 清理未使用的镜像
docker image prune -a

# 清理未使用的卷
docker volume prune

# 清理所有未使用资源
docker system prune -a --volumes
```

---

## 🔍 故障排查

### 常见问题

#### 1. 端口已被占用

**错误信息**：
```
Error starting userland proxy: listen tcp 0.0.0.0:8080: bind: address already in use
```

**解决方法**：

```bash
# Linux/macOS - 查找占用端口的进程
sudo lsof -i :8080
sudo kill -9 <PID>

# Windows - 查找占用端口的进程
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# 或修改docker-compose.yml中的端口映射
ports:
  - "8888:8080"  # 改用8888端口
```

#### 2. 配置文件未找到

**错误信息**：
```
Error: Config file not found: config_stock.json
```

**解决方法**：

```bash
# 检查配置文件
ls -l config_stock.json

# 如果不存在，复制示例
cp config_stock.json.example config_stock.json

# 编辑配置
vim config_stock.json

# 重启服务
docker-compose restart
```

#### 3. 内存不足

**错误信息**：
```
OOMKilled
```

**解决方法**：

在 `docker-compose.yml` 中增加内存限制：

```yaml
services:
  stock-analyzer:
    mem_limit: 1g
    mem_reservation: 512m
```

或释放系统内存：

```bash
# Linux
sudo sh -c 'echo 3 > /proc/sys/vm/drop_caches'

# 清理Docker缓存
docker system prune -a
```

#### 4. 网络连接问题

**错误信息**：
```
dial tcp: lookup stock-analyzer on 127.0.0.11:53: no such host
```

**解决方法**：

```bash
# 重建网络
docker-compose down
docker network prune
docker-compose up -d

# 或指定DNS
# 在docker-compose.yml中添加：
services:
  stock-analyzer:
    dns:
      - 8.8.8.8
      - 114.114.114.114
```

#### 5. 构建失败

**错误信息**：
```
failed to solve: process "/bin/sh -c go mod download" did not complete successfully
```

**解决方法**：

```bash
# 清理构建缓存
docker builder prune -a

# 使用国内镜像
# 已在Dockerfile中配置了国内源

# 重新构建
docker-compose build --no-cache
```

### 健康检查

```bash
# 检查容器健康状态
docker ps --format "table {{.Names}}\t{{.Status}}"

# 检查API健康
curl http://localhost:8080/api/stocks

# 检查Web健康
curl http://localhost/health

# 查看容器资源使用
docker stats stock-analyzer stock-web
```

### 日志分析

```bash
# 查看错误日志
docker-compose logs stock-analyzer | grep -i error

# 查看最近的警告
docker-compose logs stock-analyzer | grep -i warn | tail -20

# 导出日志
docker-compose logs > docker-logs.txt

# 查看容器内的日志文件
docker exec stock-analyzer cat /app/logs/stock_analyzer.log
```

---

## ⚡ 性能优化

### 1. 镜像优化

- ✅ 使用多阶段构建减少镜像体积
- ✅ 使用Alpine Linux作为基础镜像
- ✅ 删除构建缓存和临时文件
- ✅ 只复制必要的文件

当前镜像大小：

```bash
docker images | grep stock
# stock-analyzer  ~20-30MB
```

### 2. 资源限制

在 `docker-compose.yml` 中配置：

```yaml
services:
  stock-analyzer:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

### 3. 网络优化

```yaml
networks:
  stock-network:
    driver: bridge
    driver_opts:
      com.docker.network.driver.mtu: 1500
```

### 4. 日志管理

```yaml
logging:
  driver: "json-file"
  options:
    max-size: "10m"
    max-file: "3"
```

### 5. 缓存优化

```bash
# 构建时使用缓存
docker-compose build

# 不使用缓存（完全重建）
docker-compose build --no-cache

# 使用BuildKit加速
DOCKER_BUILDKIT=1 docker-compose build
```

---

## 🔒 安全建议

### 1. 使用非root用户

Dockerfile 中已配置：

```dockerfile
RUN addgroup -g 1000 stockapp && \
    adduser -D -u 1000 -G stockapp stockapp
USER stockapp
```

### 2. 只读挂载敏感文件

```yaml
volumes:
  - ./config_stock.json:/app/config_stock.json:ro  # 只读
```

### 3. 网络隔离

```yaml
networks:
  stock-network:
    internal: true  # 内部网络，不暴露到外部
```

### 4. 限制容器能力

```yaml
services:
  stock-analyzer:
    cap_drop:
      - ALL
    cap_add:
      - NET_BIND_SERVICE
```

### 5. 使用secrets管理敏感信息

```yaml
secrets:
  api_key:
    file: ./secrets/api_key.txt

services:
  stock-analyzer:
    secrets:
      - api_key
```

### 6. 定期更新

```bash
# 更新基础镜像
docker pull golang:1.24-alpine
docker pull alpine:3.19
docker pull nginx:1.25-alpine

# 重新构建
docker-compose build --no-cache
```

### 7. 扫描漏洞

```bash
# 使用Trivy扫描镜像
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
  aquasec/trivy image stock-analyzer:latest
```

---

## 📚 附加资源

### Docker命令速查

| 命令 | 说明 |
|-----|------|
| `docker-compose up -d` | 后台启动 |
| `docker-compose down` | 停止并删除 |
| `docker-compose ps` | 查看状态 |
| `docker-compose logs -f` | 实时日志 |
| `docker-compose restart` | 重启 |
| `docker-compose exec <service> sh` | 进入容器 |
| `docker-compose build --no-cache` | 重新构建 |
| `docker system prune -a` | 清理系统 |

### 目录结构

```
nofx-stock/
├── Dockerfile              # 构建镜像
├── docker-compose.yml      # 服务编排
├── .dockerignore          # 忽略文件
├── nginx.conf             # Nginx配置
├── docker-start.sh        # Linux启动脚本
├── docker-start.bat       # Windows启动脚本
├── .env.example           # 环境变量示例
├── config_stock.json      # 配置文件
└── logs/                  # 日志目录
```

### 相关文档

- [Docker官方文档](https://docs.docker.com/)
- [Docker Compose文档](https://docs.docker.com/compose/)
- [Docker最佳实践](https://docs.docker.com/develop/dev-best-practices/)
- [项目主文档](README_STOCK.md)
- [快速开始指南](使用说明.md)

---

## ❓ 获取帮助

如果遇到问题：

1. **查看日志**：`./docker-start.sh logs` 或 `docker-compose logs -f`
2. **检查状态**：`./docker-start.sh status` 或 `docker-compose ps`
3. **查看本文档**的故障排查章节
4. **提交Issue**到项目仓库

---

**🎉 部署完成后，访问 http://localhost 开始使用股票分析系统！**

