package main

import (
	_ "embed"
	"os"
	"path/filepath"
)

// embedFont 内嵌 assets/font.ttf（当前为 Noto Sans SC，OFL 授权），
// 供 Mac / Windows 通用渲染，无需任何网络拉取或系统字体依赖。
// 替换字体时只需覆盖 assets/font.ttf 并重新构建，无需改动代码。
//
//go:embed assets/font.ttf
var embedFont []byte

// setupEmbeddedFont 将内嵌字体写入临时文件，并通过 FYNE_FONT 环境变量
// 注入给 Fyne（必须在 app.New() 之前调用）。
func setupEmbeddedFont() {
	if len(embedFont) == 0 {
		return
	}
	dir, err := os.MkdirTemp("", "bmp2tif-font")
	if err != nil {
		return
	}
	p := filepath.Join(dir, "font.ttf")
	if err := os.WriteFile(p, embedFont, 0o644); err != nil {
		return
	}
	os.Setenv("FYNE_FONT", p)
}
