#!/usr/bin/env bash
# 构建 macOS 版本并打包为可安装到"应用程序"的 .app。
# 前置：Go >= 1.24；Xcode 命令行工具（xcode-select --install）。
set -e

mkdir -p dist

echo "==> 编译 macOS (Apple Silicon)"
GOOS=darwin GOARCH=arm64 go build -o dist/bmp2tif-mac .

echo "==> 打包为 Bmp2Tif.app"
bash package-mac.sh

echo "完成：dist/Bmp2Tif.app（可拖入 /Applications 安装）"
