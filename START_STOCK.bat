@echo off
chcp 65001 >nul
echo ====================================
echo  AI股票实时分析系统 启动脚本
echo ====================================
echo.

REM 检查配置文件
if not exist "config_stock.json" (
    echo [错误] 找不到配置文件 config_stock.json
    echo.
    echo 请先复制 config_stock.json.example 并重命名为 config_stock.json
    echo 然后编辑配置文件填写您的API密钥
    pause
    exit /b 1
)

echo [1/3] 检查Go环境...
go version >nul 2>&1
if errorlevel 1 (
    echo [错误] 未检测到Go环境，请先安装Go 1.21+
    pause
    exit /b 1
)
echo ✓ Go环境正常

echo.
echo [2/3] 编译程序...
go build -o stock_analyzer.exe main_stock.go
if errorlevel 1 (
    echo [错误] 编译失败，请检查代码
    pause
    exit /b 1
)
echo ✓ 编译成功

echo.
echo [3/3] 启动股票分析系统...
echo.
echo ========================================
echo  系统启动中...
echo  Web监控面板: http://localhost:9090
echo  按 Ctrl+C 停止运行
echo ========================================
echo.

stock_analyzer.exe

pause

