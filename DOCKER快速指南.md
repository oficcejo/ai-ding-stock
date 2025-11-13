# 🐳 Docker快速部署指南

5分钟快速部署股票分析系统！

---

## 📋 前置条件

确保已安装Docker和Docker Compose：

```bash
# 检查Docker版本
docker --version

# 检查Docker Compose版本
docker-compose --version
# 或
docker compose version
```

---

## 🚀 三步部署

### 第1步：配置系统

```bash
# 复制配置示例
cp config_stock.json.example config_stock.json

# 编辑配置（填写API密钥和股票代码）
vim config_stock.json  # Linux/macOS
notepad config_stock.json  # Windows
```

**最小配置示例**：

```json
{
  "tdx_api_url": "http://your-tdx-api:5000",
  "ai_config": {
    "deepseek": {
      "api_key": "sk-xxxxxxxxxxxxxxxx",
      "enabled": true
    }
  },
  "stocks": [
    {
      "code": "600519",
      "name": "贵州茅台",
      "enabled": true,
      "scan_interval_minutes": 5
    }
  ]
}
```

### 第2步：启动服务

#### 🐧 Linux/macOS

```bash
# 添加执行权限
chmod +x docker-start.sh

# 启动
./docker-start.sh start
```

#### 🪟 Windows

直接双击运行：
```
docker-start.bat
```

或命令行：
```cmd
docker-start.bat start
```

### 第3步：访问系统

服务启动后，在浏览器中打开：

- **📊 Web界面**: http://localhost
- **🔌 API接口**: http://localhost:8080/api/stocks

---

## 🎛️ 常用命令

### 查看服务状态

```bash
# Linux/macOS
./docker-start.sh status

# Windows
docker-start.bat status
```

### 查看实时日志

```bash
# Linux/macOS
./docker-start.sh logs

# Windows
docker-start.bat logs
```

### 停止服务

```bash
# Linux/macOS
./docker-start.sh stop

# Windows
docker-start.bat stop
```

### 重启服务

```bash
# Linux/macOS
./docker-start.sh restart

# Windows
docker-start.bat restart
```

---

## 🔧 手动操作（不使用脚本）

如果不想使用启动脚本，可以直接使用docker-compose：

```bash
# 启动（后台运行）
docker-compose up -d

# 查看状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 重启服务
docker-compose restart
```

---

## 📁 文件结构

```
nofx-stock/
├── Dockerfile              # Docker镜像定义
├── docker-compose.yml      # 服务编排配置
├── nginx.conf             # Web服务器配置
├── docker-start.sh        # Linux启动脚本
├── docker-start.bat       # Windows启动脚本
├── config_stock.json      # 配置文件（需自行创建）
└── logs/                  # 日志目录（自动创建）
```

---

## ❓ 常见问题

### 1. 端口被占用

**问题**：启动时报错 `address already in use`

**解决**：

- 方法1：修改 `docker-compose.yml` 中的端口映射
  ```yaml
  ports:
    - "8888:8080"  # 改为8888端口
  ```

- 方法2：停止占用端口的程序
  ```bash
  # Linux/macOS
  sudo lsof -i :8080
  
  # Windows
  netstat -ano | findstr :8080
  ```

### 2. 配置文件错误

**问题**：启动后服务不工作

**解决**：

1. 检查 `config_stock.json` 格式是否正确
2. 确保填写了必需的API密钥
3. 查看日志排查问题：`./docker-start.sh logs`

### 3. 无法访问Web界面

**问题**：浏览器打不开 http://localhost

**解决**：

1. 检查服务是否正常运行：`docker-compose ps`
2. 检查防火墙是否拦截了80端口
3. 尝试访问 http://127.0.0.1 或 http://本机IP

### 4. 容器无法启动

**问题**：Docker容器启动失败

**解决**：

```bash
# 查看详细日志
docker-compose logs stock-analyzer

# 重新构建镜像
docker-compose build --no-cache

# 清理并重启
docker-compose down -v
docker-compose up -d --build
```

---

## 🔄 更新系统

当代码更新后，重新部署：

```bash
# 1. 拉取最新代码
git pull

# 2. 使用脚本更新
./docker-start.sh update

# 或手动更新
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

---

## 🧹 清理数据

如需完全清理（注意：会删除所有容器和镜像）：

```bash
# 使用脚本
./docker-start.sh clean

# 或手动清理
docker-compose down -v --rmi all
docker system prune -a
```

---

## 📚 更多帮助

- **详细文档**: [DOCKER_DEPLOY.md](DOCKER_DEPLOY.md)
- **使用指南**: [README_STOCK.md](README_STOCK.md)
- **快速开始**: [使用说明.md](使用说明.md)

---

## 🎉 就是这么简单！

三步完成部署：
1. ✅ 配置 `config_stock.json`
2. ✅ 运行 `./docker-start.sh start`
3. ✅ 访问 http://localhost

**Happy Trading! 🚀**

