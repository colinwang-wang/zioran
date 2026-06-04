package config

import (
	"os"
	"time"

	"github.com/zioran/backend/pkg/oauth"
	"github.com/zioran/backend/pkg/payment"
	"github.com/zioran/backend/pkg/sms"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
	SMS      sms.SMSConfig  `yaml:"sms"`
	Payment  PaymentConfig  `yaml:"payment"`
	OAuth    OAuthConfig    `yaml:"oauth"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type JWTConfig struct {
	Secret string        `yaml:"secret"`
	Expire time.Duration `yaml:"expire"`
}

type PaymentConfig struct {
	Wechat payment.WechatPayConfig `yaml:"wechat"`
	Alipay payment.AlipayConfig    `yaml:"alipay"`
}

type OAuthConfig struct {
	Wechat oauth.WechatOAuthConfig `yaml:"wechat"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
