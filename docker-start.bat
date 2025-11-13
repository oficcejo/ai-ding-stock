@echo off
REM ================================
REM 股票分析系统 Docker 启动脚本 (Windows)
REM ================================

setlocal enabledelayedexpansion

REM 设置颜色代码（需要Windows 10+）
set "INFO=[94m[INFO][0m"
set "SUCCESS=[92m[SUCCESS][0m"
set "WARN=[93m[WARN][0m"
set "ERROR=[91m[ERROR][0m"

REM 检查Docker是否安装
:check_docker
docker --version >nul 2>&1
if errorlevel 1 (
    echo %ERROR% Docker 未安装，请先安装 Docker Desktop
    pause
    exit /b 1
)

docker-compose --version >nul 2>&1
if errorlevel 1 (
    docker compose version >nul 2>&1
    if errorlevel 1 (
        echo %ERROR% Docker Compose 未安装
        pause
        exit /b 1
    )
    set "COMPOSE_CMD=docker compose"
) else (
    set "COMPOSE_CMD=docker-compose"
)

REM 检查参数
if "%1"=="" goto :start
if "%1"=="start" goto :start
if "%1"=="stop" goto :stop
if "%1"=="restart" goto :restart
if "%1"=="logs" goto :logs
if "%1"=="status" goto :status
if "%1"=="clean" goto :clean
if "%1"=="update" goto :update
if "%1"=="shell" goto :shell
if "%1"=="help" goto :usage
goto :usage

REM ======= 启动服务 =======
:start
echo %INFO% 启动股票分析系统...

REM 检查配置文件
if not exist "config_stock.json" (
    echo %WARN% 配置文件 config_stock.json 不存在
    if exist "config_stock.json.example" (
        echo %INFO% 复制示例配置文件...
        copy config_stock.json.example config_stock.json >nul
        echo %WARN% 请编辑 config_stock.json 填写您的配置
        echo.
        echo 至少需要配置：
        echo   - TDX API URL
        echo   - AI API配置（DeepSeek或Qwen）
        echo   - 要监控的股票代码
        echo.
        pause
        notepad config_stock.json
    ) else (
        echo %ERROR% 找不到配置示例文件
        pause
        exit /b 1
    )
)

REM 创建日志目录
if not exist "logs" mkdir logs

echo %INFO% 构建并启动容器...
%COMPOSE_CMD% up -d --build

if errorlevel 1 (
    echo %ERROR% 启动失败
    pause
    exit /b 1
)

echo %SUCCESS% 服务启动成功！
echo.
echo 🎉 股票分析系统已启动
echo.
echo 📊 访问地址：
echo   - Web界面: http://localhost
echo   - API接口: http://localhost:8080/api/stocks
echo.
echo 📝 查看日志：
echo   - 实时日志: %~nx0 logs
echo   - 文件日志: .\logs\
echo.
pause
exit /b 0

REM ======= 停止服务 =======
:stop
echo %INFO% 停止股票分析系统...
%COMPOSE_CMD% down
echo %SUCCESS% 服务已停止
pause
exit /b 0

REM ======= 重启服务 =======
:restart
echo %INFO% 重启股票分析系统...
call :stop
timeout /t 2 /nobreak >nul
call :start
exit /b 0

REM ======= 查看日志 =======
:logs
echo %INFO% 查看实时日志（按 Ctrl+C 退出）...
%COMPOSE_CMD% logs -f --tail=100
exit /b 0

REM ======= 查看状态 =======
:status
echo %INFO% 服务状态：
%COMPOSE_CMD% ps
echo.
echo %INFO% 容器健康状态：
docker ps --filter "name=stock-" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
pause
exit /b 0

REM ======= 清理 =======
:clean
echo %WARN% 这将删除所有容器、镜像和卷数据！
set /p confirm="确定要继续吗？(yes/no): "
if /i not "%confirm%"=="yes" (
    echo %INFO% 已取消
    pause
    exit /b 0
)

echo %INFO% 清理中...
%COMPOSE_CMD% down -v --rmi all
echo %SUCCESS% 清理完成
pause
exit /b 0

REM ======= 更新 =======
:update
echo %INFO% 更新股票分析系统...
%COMPOSE_CMD% pull
%COMPOSE_CMD% up -d --build
echo %SUCCESS% 更新完成
pause
exit /b 0

REM ======= 进入容器 =======
:shell
echo %INFO% 进入后端容器...
docker exec -it stock-analyzer sh
exit /b 0

REM ======= 使用说明 =======
:usage
echo ================================
echo 股票分析系统 Docker 管理脚本
echo ================================
echo.
echo 用法: %~nx0 [命令]
echo.
echo 命令：
echo   start   - 启动服务（默认）
echo   stop    - 停止服务
echo   restart - 重启服务
echo   logs    - 查看实时日志
echo   status  - 查看服务状态
echo   clean   - 清理所有数据（危险操作）
echo   update  - 更新并重启服务
echo   shell   - 进入后端容器
echo   help    - 显示帮助信息
echo.
echo 示例：
echo   %~nx0          启动服务
echo   %~nx0 start    启动服务
echo   %~nx0 logs     查看日志
echo   %~nx0 status   查看状态
echo.
pause
exit /b 0

