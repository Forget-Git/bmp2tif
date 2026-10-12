package converter

import (
	"image"
	"image/color"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

// TestConvertFile 验证：用 x/image/bmp 生成测试图 -> ConvertFile 转 TIF(Deflate)
// -> 再用 x/image/tiff 解码，确认输出的 TIF 可被正确读取。
func TestConvertFile(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.bmp")

	f, err := os.Create(src)
	if err != nil {
		t.Fatal(err)
	}
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255})
	if err := bmp.Encode(f, img); err != nil {
		t.Fatal(err)
	}
	f.Close()

	dst := filepath.Join(dir, "a.tif")
	if err := ConvertFile(src, dst); err != nil {
		t.Fatalf("ConvertFile: %v", err)
	}

	df, err := os.Open(dst)
	if err != nil {
		t.Fatal(err)
	}
	defer df.Close()

	dec, err := tiff.Decode(df)
	if err != nil {
		t.Fatalf("tiff decode: %v", err)
	}
	if dec.Bounds().Dx() != 8 || dec.Bounds().Dy() != 8 {
		t.Fatalf("unexpected bounds: %v", dec.Bounds())
	}
	if c := dec.At(0, 0); !colorsEqual(c, color.RGBA{255, 0, 0, 255}) {
		t.Fatalf("pixel mismatch: %v", c)
	}
}

func colorsEqual(a, b color.Color) bool {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return ar == br && ag == bg && ab == bb && aa == ba
}
