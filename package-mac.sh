#!/usr/bin/env bash
# 将 bmp2tif 打包为 macOS .app，可直接拖到 /Applications 安装。
set -e

APP_NAME="Bmp2Tif"
BUILD_DIR="dist"
APP_BUNDLE="$BUILD_DIR/$APP_NAME.app"
EXE="$BUILD_DIR/bmp2tif-mac"

# 1) 构建（arm64；如需 Intel 版改为 GOARCH=amd64）
go build -o "$EXE" .

# 2) 清理旧 bundle 并建目录
rm -rf "$APP_BUNDLE"
mkdir -p "$APP_BUNDLE/Contents/MacOS" "$APP_BUNDLE/Contents/Resources"

# 3) 放入可执行文件（名称须与 Info.plist 的 CFBundleExecutable 一致）
cp "$EXE" "$APP_BUNDLE/Contents/MacOS/$APP_NAME"

# 3.1) 由源图 assets/AppIcon.png 生成 .icns（含 iconset），并刷新 Windows 的 .ico
if [ -f "assets/AppIcon.png" ]; then
  ICONSET="assets/AppIcon.iconset"
  rm -rf "$ICONSET"; mkdir -p "$ICONSET"
  sips -z 16 16    assets/AppIcon.png --out "$ICONSET/icon_16x16.png"       >/dev/null
  sips -z 32 32    assets/AppIcon.png --out "$ICONSET/icon_16x16@2x.png"   >/dev/null
  sips -z 32 32    assets/AppIcon.png --out "$ICONSET/icon_32x32.png"       >/dev/null
  sips -z 64 64    assets/AppIcon.png --out "$ICONSET/icon_32x32@2x.png"   >/dev/null
  sips -z 128 128  assets/AppIcon.png --out "$ICONSET/icon_128x128.png"     >/dev/null
  sips -z 256 256  assets/AppIcon.png --out "$ICONSET/icon_128x128@2x.png" >/dev/null
  sips -z 256 256  assets/AppIcon.png --out "$ICONSET/icon_256x256.png"     >/dev/null
  sips -z 512 512  assets/AppIcon.png --out "$ICONSET/icon_256x256@2x.png" >/dev/null
  sips -z 512 512  assets/AppIcon.png --out "$ICONSET/icon_512x512.png"     >/dev/null
  sips -z 1024 1024 assets/AppIcon.png --out "$ICONSET/icon_512x512@2x.png" >/dev/null
  iconutil --convert icns "$ICONSET" --out assets/AppIcon.icns
  cp "assets/AppIcon.icns" "$APP_BUNDLE/Contents/Resources/AppIcon.icns"
  # 同步刷新 Windows 图标
  go run ./tools/makico
else
  echo "警告：未找到 assets/AppIcon.png，跳过图标生成"
fi

# 4) 写入 Info.plist
cat > "$APP_BUNDLE/Contents/Info.plist" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleName</key>             <string>$APP_NAME</string>
    <key>CFBundleDisplayName</key>      <string>$APP_NAME</string>
    <key>CFBundleExecutable</key>       <string>$APP_NAME</string>
    <key>CFBundleIdentifier</key>       <string>com.example.bmp2tif</string>
    <key>CFBundleVersion</key>          <string>1.0.0</string>
    <key>CFBundleShortVersionString</key><string>1.0.0</string>
    <key>CFBundlePackageType</key>      <string>APPL</string>
    <key>CFBundleInfoDictionaryVersion</key><string>6.0</string>
    <key>LSMinimumSystemVersion</key>   <string>10.13</string>
    <key>NSHighResolutionCapable</key>  <true/>
    <key>NSPrincipalClass</key>         <string>NSApplication</string>
    <key>CFBundleIconFile</key>         <string>AppIcon</string>
</dict>
</plist>
EOF

# 5) 设可执行权限 + 临时(adr-hoc)签名 + 移除隔离属性，避免"无法验证开发者"拦截
chmod +x "$APP_BUNDLE/Contents/MacOS/$APP_NAME"
codesign --force --deep --sign - "$APP_BUNDLE" 2>/dev/null || true
xattr -dr com.apple.quarantine "$APP_BUNDLE" 2>/dev/null || true

echo "已生成：$APP_BUNDLE"
echo "安装方式（拖到应用程序）："
echo "  cp -r \"$APP_BUNDLE\" /Applications/"
