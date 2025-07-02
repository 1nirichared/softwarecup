package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	AI       AIConfig       `mapstructure:"ai"`
	Xunfei   XunfeiConfig   `mapstructure:"xunfei"`
	LocalAI  LocalAIConfig  `mapstructure:"local_ai"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"` // 过期时间(小时)
}

type AIConfig struct {
	// 主要AI提供商
	Provider string `mapstructure:"provider"`

	// DeepSeek配置
	DeepSeekAPIKey  string `mapstructure:"deepseek_api_key"`
	DeepSeekBaseURL string `mapstructure:"deepseek_base_url"`
	DeepSeekModel   string `mapstructure:"deepseek_model"`

	// OpenAI配置（备用）
	OpenAIAPIKey  string `mapstructure:"openai_api_key"`
	OpenAIBaseURL string `mapstructure:"openai_base_url"`
	OpenAIModel   string `mapstructure:"openai_model"`

	// 通用配置
	MaxTokens   int     `mapstructure:"max_tokens"`
	Temperature float64 `mapstructure:"temperature"`
	Timeout     int     `mapstructure:"timeout"`
}

type XunfeiConfig struct {
	AppID     string `mapstructure:"app_id"`
	APISecret string `mapstructure:"api_secret"`
	APIKey    string `mapstructure:"api_key"`
	BaseURL   string `mapstructure:"base_url"`
	Model     string `mapstructure:"model"`
	MaxTokens int    `mapstructure:"max_tokens"`
	Timeout   int    `mapstructure:"timeout"`
}

type LocalAIConfig struct {
	Enabled     bool    `mapstructure:"enabled"`
	Type        string  `mapstructure:"type"`
	BaseURL     string  `mapstructure:"base_url"`
	Model       string  `mapstructure:"model"`
	MaxTokens   int     `mapstructure:"max_tokens"`
	Temperature float64 `mapstructure:"temperature"`
	Timeout     int     `mapstructure:"timeout"`
}

var GlobalConfig Config

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		return err
	}

	return nil
}
