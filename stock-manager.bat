@echo off
REM ================================
REM 股票监测快速管理工具
REM ================================

setlocal enabledelayedexpansion

set "CONFIG_FILE=config_stock.json"
set "INFO=[94m[INFO][0m"
set "SUCCESS=[92m[SUCCESS][0m"
set "WARN=[93m[WARN][0m"

echo ================================
echo   股票监测快速管理工具
echo ================================
echo.

REM 检查配置文件
if not exist "%CONFIG_FILE%" (
    echo %WARN% 配置文件不存在: %CONFIG_FILE%
    pause
    exit /b 1
)

:menu
echo 请选择操作：
echo.
echo   1. 查看当前监测的股票
echo   2. 编辑配置文件
echo   3. 快速重启容器（应用配置）
echo   4. 查看容器日志
echo   5. 退出
echo.
set /p choice="请输入选项 (1-5): "

if "%choice%"=="1" goto :show_stocks
if "%choice%"=="2" goto :edit_config
if "%choice%"=="3" goto :restart_container
if "%choice%"=="4" goto :show_logs
if "%choice%"=="5" goto :exit
goto :menu

REM ======= 查看当前监测的股票 =======
:show_stocks
echo.
echo %INFO% 当前监测的股票：
echo.
type %CONFIG_FILE% | findstr /C:"\"code\"" /C:"\"name\"" /C:"\"enabled\""
echo.
pause
goto :menu

REM ======= 编辑配置 =======
:edit_config
echo.
echo %INFO% 打开配置文件编辑器...
echo.
echo 修改提示：
echo   - 添加股票：复制一个股票配置块，修改code和name
echo   - 暂停监测：将 "enabled": true 改为 false
echo   - 删除股票：删除整个股票配置块（注意JSON格式）
echo   - 修改间隔：调整 scan_interval_minutes 值
echo.
echo 编辑完成后保存并关闭，然后选择"快速重启容器"应用配置
echo.
pause
notepad %CONFIG_FILE%
goto :menu

REM ======= 重启容器 =======
:restart_container
echo.
echo %INFO% 正在快速重启容器（无需重新构建）...
docker-compose restart stock-analyzer

if errorlevel 1 (
    echo %WARN% 重启失败，尝试完整重启...
    docker-compose down
    timeout /t 2 /nobreak >nul
    docker-compose up -d
)

echo.
echo %SUCCESS% 容器已重启！新配置已生效
echo.
echo 提示：可以通过"查看容器日志"确认新股票已加载
echo.
pause
goto :menu

REM ======= 查看日志 =======
:show_logs
echo.
echo %INFO% 显示最近日志（按 Ctrl+C 停止）...
echo.
timeout /t 2 /nobreak >nul
docker-compose logs --tail=50 stock-analyzer
echo.
pause
goto :menu

REM ======= 退出 =======
:exit
echo.
echo 👋 谢谢使用！
timeout /t 1 /nobreak >nul
exit /b 0

