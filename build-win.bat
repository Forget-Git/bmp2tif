@echo off
REM 在 Windows 上构建 bmp2tif 可执行文件（含图标 / 版本信息 / 高 DPI 清单）。
REM 前置：
REM   1) 安装 Go (>=1.24)
REM   2) 安装 C 编译器：MSYS2 + MinGW-w64 或 TDM-GCC（Fyne 依赖 CGO）
REM   3) 联网（首次会通过 go run 拉取 go-winres 资源工具）
REM 步骤：go-winres 生成 resource.syso -> go build 自动嵌入图标

cd /d %~dp0

REM 由源图 assets/AppIcon.png 刷新 Windows 图标 winres/AppIcon.ico（需先准备该 PNG）
go run ./tools/makico
if not %errorlevel%==0 (
    echo [提示] makico 生成图标失败，将使用已有的 winres/AppIcon.ico
)

REM 用 go-winres 生成资源（图标/清单/版本），写入 resource.syso
go run github.com/tc-hib/go-winres@latest make
if not %errorlevel%==0 (
    echo [错误] go-winres 生成资源失败，请确认已联网且 Go 版本 >=1.24
    pause
    exit /b 1
)

go build -o dist\bmp2tif-win.exe .
if %errorlevel%==0 (
    echo 已生成 dist\bmp2tif-win.exe（含应用图标）
) else (
    echo 构建失败：请确认已安装 Go 与 C 编译器(MinGW-w64 / TDM-GCC)
)
pause
