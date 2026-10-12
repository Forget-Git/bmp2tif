# Bmp2Tif

将 BMP 图片批量转换为 TIF 的桌面客户端。用 Go + [Fyne](https://fyne.io/) 编写，
**纯单文件分发、无需任何环境配置**，提供 macOS 与 Windows 版本，支持系统原生的文件夹选择框。

> 原 Python 版（`bmp2tif.py`）依赖 Pillow（底层 C 库），安装麻烦、跨平台易出错；
> 本 Go 版本改用纯标准库 + `golang.org/x/image` 完成转换，编译出一个二进制即可双击运行。

## 功能

- 选择输入文件夹，批量将其中所有 `.bmp` 转为 `.tif`
- 输出文件夹可自定义（留空则输出到输入文件夹）
- 递归处理子目录（可选）
- 保留子目录结构输出（可选）
- 重名文件冲突策略：覆盖 / 跳过（可选）
- 实时转换日志 + 进度条
- 记住上次选择的路径（写入 `~/.config/bmp2tif/config.json`）
- 内嵌中文字体（Noto Sans SC），界面中文显示正常，不依赖系统字体
- 文件夹选择调用**系统原生对话框**（macOS `NSOpenPanel` / Windows `IFileDialog`）

## 技术说明

| 项 | 方案 |
|----|------|
| 语言 / GUI | Go 1.24 + Fyne v2 |
| 图像转换 | `golang.org/x/image/bmp` 解码 + `golang.org/x/image/tiff` 编码 |
| 压缩方式 | **Deflate（zlib，无损）** |
| 文件夹选择 | `github.com/sqweek/dialog`（原生系统对话框） |
| 中文字体 | 内嵌 `assets/font.ttf`（Noto Sans SC，OFL 授权） |
| 图标 | `assets/AppIcon.png`（源图，带透明通道）→ 自动生成 Mac `AppIcon.icns` / Win `AppIcon.ico` |

> **关于压缩方式**：原 Python 版使用 `compression='tiff_lzw'`。
> Go 生态的 TIFF 编码器（`golang.org/x/image/tiff`）只支持 None / Deflate，**不支持 LZW 编码**，
> 因此改用 Deflate 替代——无损、压缩率通常更优、且所有 TIFF 阅读器（Photoshop / GIMP /
> ImageMagick / Pillow 等）均支持，对二次修图无任何影响。

## 目录结构

```
bmp2tif/
├── go.mod
├── main.go                 # 入口 + GUI 布局
├── embed.go                # 内嵌字体注入（FYNE_FONT）
├── converter/
│   ├── converter.go        # 核心转换逻辑（纯标准库，可单测）
│   └── converter_test.go   # 单元测试
├── config/
│   └── config.go           # 输入输出路径持久化
├── assets/
│   ├── font.ttf            # 内嵌中文字体（Noto Sans SC）
│   └── AppIcon.png         # 应用图标源图（1024×1024，带透明通道；改图只改它）
├── winres/                 # Windows 资源（图标 / 版本 / 清单）
│   ├── winres.json
│   ├── app.manifest
│   └── AppIcon.ico
├── build-mac.sh            # 编译 + 打包为 .app
├── package-mac.sh          # 仅打包 .app（被 build-mac.sh 调用）
├── build-win.bat           # Windows 编译（生成带图标的 .exe）
├── installer.iss           # Inno Setup 安装脚本
└── dist/                   # 构建产物（bmp2tif-mac / Bmp2Tif.app / bmp2tif-win.exe）
```

## macOS 构建与安装

前置：Go ≥ 1.24，且已安装 Xcode 命令行工具（`xcode-select --install`）。

```bash
# 一条命令：编译 + 打包为 Bmp2Tif.app
bash build-mac.sh

# 安装到应用程序
cp -r dist/Bmp2Tif.app /Applications/
```

之后在 Launchpad / 应用程序里即可启动。脚本已对 `.app` 做 ad-hoc 签名并移除隔离属性，本机可直接打开。

> 分发给他人时，需使用 Apple 开发者证书签名 + 公证（`xcrun notarytool`），否则他人首次打开会被 Gatekeeper 拦截。

## Windows 构建与安装

前置（在 Windows 机器上）：

1. 安装 [Go ≥ 1.24](https://go.dev/dl/)
2. 安装 C 编译器（Fyne 依赖 CGO，**二选一**）：
   - MSYS2 + MinGW-w64：`pacman -S mingw-w64-x86_64-gcc`，并把其 `bin` 加入 PATH
   - 或 [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)
3. 将整个项目目录拷贝到 Windows（或 git 拉取）

```bat
:: 编译（首次会联网拉取 go-winres 把图标嵌入 .exe）
build-win.bat
```

生成 `dist\bmp2tif-win.exe`（自包含，字体已内嵌，无需系统字体）。

**做成安装包**（可选）：用 [Inno Setup](https://jrsoftware.org/isdl.php) 打开 `installer.iss` 编译，
得到 `Bmp2Tif-Setup.exe`（安装到 程序文件 + 开始菜单 + 桌面快捷方式）。

> 本项目在 macOS 上无法通过 MinGW 交叉编译 Windows 版（沙箱无交叉编译器且 Fyne 依赖 CGO），
> 因此在 Windows 机器上构建。如需在本机交叉编译，需先安装 MinGW-w64 交叉工具链。

## 替换字体 / 图标（可选）

- **换字体**：把任意含中文的 `.ttf` 重命名为 `font.ttf` 覆盖 `assets/font.ttf`，重新构建即可（OFL 授权的 Noto Sans SC 已内置）。
- **换图标**：只改一张源图即可，两边图标会自动重新生成（务必使用**带透明通道**的 PNG，否则边角会变白）：
  - 用你的图标覆盖 `assets/AppIcon.png`（建议 1024×1024）
  - Mac：重跑 `package-mac.sh`（会基于 PNG 重新生成 `assets/AppIcon.icns`）
  - Windows：重跑 `build-win.bat`（`tools/makico` 会基于 PNG 重新生成 `winres/AppIcon.ico`）

## 常见问题

**Q：运行时中文显示成问号？**
A：内嵌字体未生效。确认 `assets/font.ttf` 存在且为含中文的 TrueType 字体，重新构建。

**Q：选择文件夹弹不出系统对话框？**
A：Mac 上首次可能需在「系统设置 → 隐私与安全性」允许；Windows 上 `sqweek/dialog` 走原生 COM，无需额外设置。

**Q：Windows 构建报 `gcc: command not found`？**
A：未安装 C 编译器，参见上文 Windows 前置条件。
