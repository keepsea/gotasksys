// internal/config/config.go

package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 结构体用于映射 config.yaml 的内容
type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
	JWT struct {
		Secret          string `yaml:"secret"`
		ExpirationHours int    `yaml:"expiration_hours"`
	} `yaml:"jwt"`
}

// LoadConfig 从 config.yaml 文件加载配置
func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
