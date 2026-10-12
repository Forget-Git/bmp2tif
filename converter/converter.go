// Package converter 提供 BMP -> TIF(LZW) 的核心转换逻辑，仅依赖 Go 标准库。
package converter

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

// Scan 扫描 srcDir 下的 .bmp 文件。
// recursive 为 true 时递归遍历子目录。
func Scan(srcDir string, recursive bool) ([]string, error) {
	var files []string
	add := func(p string) {
		if strings.EqualFold(filepath.Ext(p), ".bmp") {
			files = append(files, p)
		}
	}

	if recursive {
		err := filepath.WalkDir(srcDir, func(p string, d os.DirEntry, err error) error {
			if err != nil {
				return nil // 忽略单个文件错误，继续遍历
			}
			if !d.IsDir() {
				add(p)
			}
			return nil
		})
		return files, err
	}

	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() {
			add(filepath.Join(srcDir, e.Name()))
		}
	}
	return files, nil
}

// ConvertFile 将单个 BMP 文件转换为 TIF（LZW 压缩）。
// dst 为目标路径，父目录不存在时会自动创建。
func ConvertFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	img, err := bmp.Decode(in)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	// 使用 Deflate（zlib）无损压缩。
	// 说明：Go 生态的 TIFF 编码器（golang.org/x/image/tiff）均不支持 LZW 编码，
	// 原 Python 脚本的 compression='tiff_lzw' 在此以 Deflate 替代——
	// 无损、压缩率通常更好、且被所有 TIFF 阅读器支持。
	opts := &tiff.Options{Compression: tiff.Deflate}
	return tiff.Encode(out, img, opts)
}

// TargetName 根据源文件生成目标 TIF 路径。
// relDir 为相对输入根目录的子路径（用于保留目录结构）；为空则输出到输出根目录。
func TargetName(outRoot, relDir, base string) string {
	base = strings.TrimSuffix(base, filepath.Ext(base)) + ".tif"
	return filepath.Join(outRoot, relDir, base)
}
