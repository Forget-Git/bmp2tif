// Package config 负责读写本地配置文件，记住用户上次选择的输入/输出路径。
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config 保存用户的界面偏好。
type Config struct {
	InputDir   string `json:"input_dir"`
	OutputDir  string `json:"output_dir"`
	Recursive  bool   `json:"recursive"`
	KeepStruct bool   `json:"keep_struct"` // 输出时保留子目录结构
}

// Path 返回配置文件完整路径：~/.config/bmp2tif/config.json
func Path() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "bmp2tif", "config.json"), nil
}

// Load 读取配置；文件不存在或解析失败时返回零值配置（不报错）。
func Load() *Config {
	cfg := &Config{}
	p, err := Path()
	if err != nil {
		return cfg
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(data, cfg)
	return cfg
}

// Save 将配置写回磁盘。
func (c *Config) Save() error {
	p, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o644)
}
