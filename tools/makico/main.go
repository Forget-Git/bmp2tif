package main

// makico 将一张带透明通道的 PNG（assets/AppIcon.png）转换为 Windows 可用的
// .ico（以 PNG 形式内嵌，保留 alpha）。Mac 侧由 package-mac.sh 用 iconutil 生成
// .icns，两边共用同一张源图，改图只需替换 assets/AppIcon.png。
import (
	"fmt"
	"io"
	"os"
)

func main() {
	src, dst := "assets/AppIcon.png", "winres/AppIcon.ico"
	in, err := os.Open(src)
	if err != nil {
		fmt.Println("[makico] 打开源图失败:", err)
		os.Exit(1)
	}
	defer in.Close()
	if _, err := io.ReadAll(io.TeeReader(in, io.Discard)); err != nil {
		fmt.Println("[makico] 读取源图失败:", err)
		os.Exit(1)
	}
	if _, err := in.Seek(0, io.SeekStart); err != nil {
		fmt.Println("[makico] 回退文件指针失败:", err)
		os.Exit(1)
	}
	data, err := io.ReadAll(in)
	if err != nil {
		fmt.Println("[makico] 读取源图失败:", err)
		os.Exit(1)
	}

	out, err := os.Create(dst)
	if err != nil {
		fmt.Println("[makico] 创建 ico 失败:", err)
		os.Exit(1)
	}
	defer out.Close()

	// ICONDIR
	out.Write([]byte{0, 0, 1, 0, 1, 0}) // reserved, type=icon, count=1
	// ICONDIRENTRY（尺寸 >255 时宽高记为 0，表示 256）
	entry := []byte{
		0, 0, // width, height (0 => 256)
		0,    // color count
		0,    // reserved
		1, 0, // planes
		32, 0, // bit count
		byte(len(data)), byte(len(data) >> 8), byte(len(data) >> 16), byte(len(data) >> 24),
		22, 0, 0, 0, // image offset (6 + 16)
	}
	out.Write(entry)
	out.Write(data)
	fmt.Println("[makico] 已生成", dst, "(", len(data), "bytes )")
}
