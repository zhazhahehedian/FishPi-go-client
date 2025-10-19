package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 配置
type Config struct {
	BaseURL   string `json:"base_url"`
	UserAgent string `json:"user_agent"`
	APIKey    string `json:"api_key"`
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		BaseURL:   "https://fishpi.cn",
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
	}
}

// GetConfigPath 获取配置文件路径(创建保存密码的文件信息)
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".fishpi")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &config, nil
}

// SaveConfig 保存配置
func SaveConfig(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("保存配置文件失败: %w", err)
	}

	return nil
}

// SaveAPIKey 保存API Key到配置
func SaveAPIKey(apiKey string) error {
	config, err := LoadConfig()
	if err != nil {
		config = DefaultConfig()
	}

	config.APIKey = apiKey
	return SaveConfig(config)
}

// GetAPIKey 从配置获取API Key
func GetAPIKey() (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", err
	}

	return config.APIKey, nil
}

