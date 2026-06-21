package config

import (
	"os"
	"strconv"
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
	applyEnvOverrides(&cfg)
	return &cfg, nil
}

func applyEnvOverrides(cfg *Config) {
	setIntEnv("SERVER_PORT", &cfg.Server.Port)

	setStringEnv("DB_HOST", &cfg.Database.Host)
	setIntEnv("DB_PORT", &cfg.Database.Port)
	setStringEnv("DB_USER", &cfg.Database.User)
	setStringEnv("DB_PASSWORD", &cfg.Database.Password)
	setStringEnv("DB_NAME", &cfg.Database.DBName)

	setStringEnv("REDIS_HOST", &cfg.Redis.Host)
	setIntEnv("REDIS_PORT", &cfg.Redis.Port)
	setStringEnv("JWT_SECRET", &cfg.JWT.Secret)
	setDurationEnv("JWT_EXPIRE", &cfg.JWT.Expire)

	setBoolEnv("PAYMENT_WECHAT_ENABLED", &cfg.Payment.Wechat.Enabled)
	setBoolEnv("PAYMENT_WECHAT_MOCK_AUTO_COMPLETE", &cfg.Payment.Wechat.MockAutoComplete)
	setStringEnv("PAYMENT_WECHAT_APP_ID", &cfg.Payment.Wechat.AppID)
	setStringEnv("PAYMENT_WECHAT_MCH_ID", &cfg.Payment.Wechat.MchID)
	setStringEnv("PAYMENT_WECHAT_API_KEY", &cfg.Payment.Wechat.APIKey)
	setStringEnv("PAYMENT_WECHAT_NOTIFY_URL", &cfg.Payment.Wechat.NotifyURL)
	setStringEnv("PAYMENT_WECHAT_SERIAL_NO", &cfg.Payment.Wechat.SerialNo)
	setStringEnv("PAYMENT_WECHAT_PRIVATE_KEY_PATH", &cfg.Payment.Wechat.PrivateKeyPath)
	setStringEnv("PAYMENT_WECHAT_CERT_PATH", &cfg.Payment.Wechat.CertPath)
	setStringEnv("PAYMENT_WECHAT_PUBLIC_KEY_PATH", &cfg.Payment.Wechat.PublicKeyPath)
	setStringEnv("PAYMENT_WECHAT_PUBLIC_KEY_ID", &cfg.Payment.Wechat.PublicKeyID)

	setBoolEnv("PAYMENT_ALIPAY_ENABLED", &cfg.Payment.Alipay.Enabled)
	setBoolEnv("PAYMENT_ALIPAY_MOCK_AUTO_COMPLETE", &cfg.Payment.Alipay.MockAutoComplete)
	setStringEnv("PAYMENT_ALIPAY_APP_ID", &cfg.Payment.Alipay.AppID)
	setStringEnv("PAYMENT_ALIPAY_PRIVATE_KEY", &cfg.Payment.Alipay.PrivateKey)
	setStringEnv("PAYMENT_ALIPAY_PUBLIC_KEY", &cfg.Payment.Alipay.AlipayPublicKey)
	setStringEnv("PAYMENT_ALIPAY_NOTIFY_URL", &cfg.Payment.Alipay.NotifyURL)
	setStringEnv("PAYMENT_ALIPAY_RETURN_URL", &cfg.Payment.Alipay.ReturnURL)

	setBoolEnv("OAUTH_WECHAT_ENABLED", &cfg.OAuth.Wechat.Enabled)
	setStringEnv("OAUTH_WECHAT_APP_ID", &cfg.OAuth.Wechat.AppID)
	setStringEnv("OAUTH_WECHAT_APP_SECRET", &cfg.OAuth.Wechat.AppSecret)
	setStringEnv("OAUTH_WECHAT_REDIRECT_URI", &cfg.OAuth.Wechat.RedirectURI)
	setStringEnv("OAUTH_WECHAT_FRONTEND_REDIRECT_URI", &cfg.OAuth.Wechat.FrontendRedirectURI)
}

func setStringEnv(key string, target *string) {
	if value := os.Getenv(key); value != "" {
		*target = value
	}
}

func setBoolEnv(key string, target *bool) {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			*target = parsed
		}
	}
}

func setIntEnv(key string, target *int) {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			*target = parsed
		}
	}
}

func setDurationEnv(key string, target *time.Duration) {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			*target = parsed
		}
	}
}
