; 用 Inno Setup 编译本脚本可生成 Windows 安装程序（安装到 程序文件 + 开始菜单 + 桌面快捷方式）。
; 前置：先运行 build-win.bat 生成 dist\bmp2tif-win.exe。
; 可选图标：把 assets 下的图标导出为 installer.ico，取消下面 SetupIconFile 注释。

[Setup]
AppName=Bmp2Tif
AppVersion=1.0.0
AppPublisher=Your Name
DefaultDirName={autopf}\Bmp2Tif
DefaultGroupName=Bmp2Tif
OutputBaseFilename=Bmp2Tif-Setup
Compression=lzma2
SolidCompression=yes
SetupIconFile=assets\AppIcon.ico
WizardStyle=modern

[Files]
Source: "dist\bmp2tif-win.exe"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\Bmp2Tif"; Filename: "{app}\bmp2tif-win.exe"
Name: "{commondesktop}\Bmp2Tif"; Filename: "{app}\bmp2tif-win.exe"

[Run]
Filename: "{app}\bmp2tif-win.exe"; Description: "运行 Bmp2Tif"; Flags: nowait postinstall skipifsilent
