package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bmp2tif/config"
	"bmp2tif/converter"

	"github.com/sqweek/dialog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// 启用内嵌的中文字体（assets/font.ttf，Noto Sans SC），修复中文显示为问号。
	// 必须在 app.New() 之前调用。
	setupEmbeddedFont()

	myApp := app.New()
	w := myApp.NewWindow("BMP → TIF 转换工具")

	cfg := config.Load()

	inputEntry := widget.NewEntry()
	inputEntry.SetText(cfg.InputDir)
	outputEntry := widget.NewEntry()
	outputEntry.SetText(cfg.OutputDir)

	recursive := cfg.Recursive
	keepStruct := cfg.KeepStruct
	recursiveCheck := widget.NewCheck("递归处理子目录", func(v bool) { recursive = v })
	recursiveCheck.SetChecked(recursive)
	keepStructCheck := widget.NewCheck("保留子目录结构", func(v bool) { keepStruct = v })
	keepStructCheck.SetChecked(keepStruct)

	conflict := widget.NewSelect([]string{"覆盖", "跳过"}, func(string) {})
	conflict.SetSelected("覆盖")

	logEntry := widget.NewMultiLineEntry() // 不禁用，保持黑色文字清晰可读（禁用会变灰色）
	logEntry.SetMinRowsVisible(10)

	progress := widget.NewProgressBar()

	var running bool
	startBtn := widget.NewButtonWithIcon("开始转换", theme.ConfirmIcon(), nil)

	// 调用系统原生的文件夹选择框（macOS 走 NSOpenPanel，Windows 走 IFileDialog）。
	chooseDir := func(target *widget.Entry) {
		dir, err := dialog.Directory().Title("选择文件夹").Browse()
		if err != nil {
			return // 用户取消
		}
		target.SetText(dir)
	}

	inputBtn := widget.NewButtonWithIcon("选择", theme.FolderOpenIcon(), func() { chooseDir(inputEntry) })
	outputBtn := widget.NewButtonWithIcon("选择", theme.FolderOpenIcon(), func() { chooseDir(outputEntry) })

	startBtn.OnTapped = func() {
		if running {
			return
		}
		inDir := strings.TrimSpace(inputEntry.Text)
		if inDir == "" {
			dialog.Message("请先选择输入文件夹").Title("提示").Info()
			return
		}
		outDir := strings.TrimSpace(outputEntry.Text)
		if outDir == "" {
			outDir = inDir // 默认输出到输入文件夹（与原脚本行为一致）
		}

		// 记住配置
		cfg.InputDir = inDir
		cfg.OutputDir = outDir
		cfg.Recursive = recursive
		cfg.KeepStruct = keepStruct
		_ = cfg.Save()

		running = true
		startBtn.Disable()
		logEntry.SetText("")
		progress.SetValue(0)

		go runConversion(inDir, outDir, recursive, keepStruct, conflict.Selected, progress, logEntry, &running, startBtn)
	}

	inputRow := container.NewBorder(nil, nil, nil, inputBtn, inputEntry)
	outputRow := container.NewBorder(nil, nil, nil, outputBtn, outputEntry)

	form := container.NewVBox(
		widget.NewLabelWithStyle("输入文件夹", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		inputRow,
		widget.NewLabel("输出文件夹（留空则输出到输入文件夹）"),
		outputRow,
		container.NewHBox(recursiveCheck, keepStructCheck, layout.NewSpacer(), widget.NewLabel("冲突："), conflict),
	)

	content := container.NewBorder(
		form,
		container.NewVBox(startBtn, progress),
		nil, nil,
		container.NewVBox(
			widget.NewLabelWithStyle("转换日志", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			logEntry,
		),
	)
	w.SetContent(content)
	w.Resize(fyne.NewSize(720, 640))
	w.ShowAndRun()
}

func runConversion(inDir, outDir string, recursive, keepStruct bool, conflict string,
	progress *widget.ProgressBar, logEntry *widget.Entry,
	running *bool, startBtn *widget.Button) {

	// Fyne v2.4 中，widget 的 setter 会内部触发 Refresh 并把重绘
	// 排入主事件循环，因此可直接在 goroutine 中调用以更新进度与日志。
	logf := func(format string, args ...interface{}) {
		logEntry.Append(fmt.Sprintf(format, args...) + "\n")
	}

	files, err := converter.Scan(inDir, recursive)
	if err != nil {
		dialog.Message("%s", err.Error()).Title("错误").Error()
		startBtn.Enable()
		*running = false
		return
	}
	if len(files) == 0 {
		logf("未找到任何 .bmp 文件")
		startBtn.Enable()
		*running = false
		return
	}

	logf("开始转换，共发现 %d 个 BMP 文件", len(files))
	var ok, fail int
	total := float64(len(files))

	for i, f := range files {
		rel := ""
		if keepStruct {
			rel, _ = filepath.Rel(inDir, filepath.Dir(f))
		}
		dst := converter.TargetName(outDir, rel, filepath.Base(f))

		if _, e := os.Stat(dst); e == nil && conflict == "跳过" {
			logf("跳过（已存在）: %s", f)
		} else {
			if e := converter.ConvertFile(f, dst); e != nil {
				fail++
				logf("失败: %s -> %v", f, e)
			} else {
				ok++
				logf("成功: %s", f)
			}
		}
		progress.SetValue(float64(i+1) / total)
	}

	logf("\n转换完成：成功 %d，失败 %d，共 %d", ok, fail, len(files))
	progress.SetValue(1)
	startBtn.Enable()
	*running = false
}
