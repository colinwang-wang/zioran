package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAppliesPaymentEnvOverrides(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	data := []byte(`
server:
  port: 8080
database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: root
  dbname: zioran
redis:
  host: 127.0.0.1
  port: 6379
jwt:
  secret: local
  expire: 72h
payment:
  wechat:
    enabled: false
    mock_auto_complete: false
    app_id: ""
    mch_id: ""
    api_key: ""
    notify_url: ""
    serial_no: ""
    private_key_path: ""
    cert_path: ""
  alipay:
    enabled: false
    mock_auto_complete: false
    app_id: ""
    private_key: ""
    alipay_public_key: ""
    notify_url: ""
    return_url: ""
oauth:
  wechat:
    enabled: false
    app_id: ""
    app_secret: ""
    redirect_uri: ""
    frontend_redirect_uri: ""
`)
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Setenv("PAYMENT_WECHAT_ENABLED", "true")
	t.Setenv("PAYMENT_WECHAT_MCH_ID", "mch-1")
	t.Setenv("PAYMENT_WECHAT_CERT_PATH", "/secure/wechat-platform.pem")
	t.Setenv("PAYMENT_ALIPAY_ENABLED", "true")
	t.Setenv("PAYMENT_ALIPAY_RETURN_URL", "https://www.zioran.com/user/transactions")

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if !cfg.Payment.Wechat.Enabled || cfg.Payment.Wechat.MchID != "mch-1" || cfg.Payment.Wechat.CertPath != "/secure/wechat-platform.pem" {
		t.Fatalf("wechat env overrides not applied: %+v", cfg.Payment.Wechat)
	}
	if !cfg.Payment.Alipay.Enabled || cfg.Payment.Alipay.ReturnURL != "https://www.zioran.com/user/transactions" {
		t.Fatalf("alipay env overrides not applied: %+v", cfg.Payment.Alipay)
	}
}
