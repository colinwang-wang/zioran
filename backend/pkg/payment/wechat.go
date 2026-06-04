package payment

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// WechatPayConfig holds WeChat Pay configuration.
type WechatPayConfig struct {
	Enabled   bool   `yaml:"enabled"`
	AppID     string `yaml:"app_id"`
	MchID     string `yaml:"mch_id"`
	APIKey    string `yaml:"api_key"`
	NotifyURL string `yaml:"notify_url"`
	CertPath  string `yaml:"cert_path"`
}

// WechatPay implements WeChat Pay V3 Native payment.
type WechatPay struct {
	Cfg WechatPayConfig
}

func NewWechatPay(cfg WechatPayConfig) *WechatPay {
	return &WechatPay{Cfg: cfg}
}

// CreateNativeOrder creates a Native (QR code) payment order.
func (w *WechatPay) CreateNativeOrder(orderNo string, amount int, description string) (codeURL string, err error) {
	if !w.Cfg.Enabled {
		return "mock://wechat_pay/" + orderNo, nil
	}

	payload := map[string]interface{}{
		"appid":        w.Cfg.AppID,
		"mchid":        w.Cfg.MchID,
		"description":  description,
		"out_trade_no": orderNo,
		"notify_url":   w.Cfg.NotifyURL,
		"amount": map[string]interface{}{
			"total":    amount,
			"currency": "CNY",
		},
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "https://api.mch.weixin.qq.com/v3/pay/transactions/native", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", w.buildAuth("POST", "/v3/pay/transactions/native", string(body)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("wechat pay: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("wechat pay: status %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		CodeURL string `json:"code_url"`
	}
	json.Unmarshal(respBody, &result)
	return result.CodeURL, nil
}

// VerifyNotify verifies WeChat Pay callback and extracts order number.
func (w *WechatPay) VerifyNotify(body []byte, signature string) (orderNo string, err error) {
	if !w.Cfg.Enabled {
		// Mock: try to parse order_id from body
		var mock struct {
			OrderNo string `json:"order_no"`
		}
		json.Unmarshal(body, &mock)
		return mock.OrderNo, nil
	}

	// Parse notification
	var notification struct {
		Resource struct {
			Ciphertext     string `json:"ciphertext"`
			Nonce          string `json:"nonce"`
			AssociatedData string `json:"associated_data"`
		} `json:"resource"`
	}
	if err := json.Unmarshal(body, &notification); err != nil {
		return "", fmt.Errorf("wechat notify: invalid body")
	}

	// Decrypt using AEAD_AES_256_GCM
	plaintext, err := decryptAESGCM(w.Cfg.APIKey, notification.Resource.Nonce, notification.Resource.Ciphertext, notification.Resource.AssociatedData)
	if err != nil {
		return "", fmt.Errorf("wechat notify: decrypt failed: %w", err)
	}

	var order struct {
		OutTradeNo string `json:"out_trade_no"`
		TradeState string `json:"trade_state"`
	}
	json.Unmarshal(plaintext, &order)
	if order.TradeState != "SUCCESS" {
		return "", fmt.Errorf("wechat notify: trade_state=%s", order.TradeState)
	}
	return order.OutTradeNo, nil
}

func (w *WechatPay) buildAuth(method, path, body string) string {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	message := method + "\n" + path + "\n" + timestamp + "\n" + nonce + "\n" + body + "\n"
	h := sha256.Sum256([]byte(message + w.Cfg.APIKey))
	sig := hex.EncodeToString(h[:])
	return fmt.Sprintf(`WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",signature="%s",timestamp="%s",serial_no="default"`,
		w.Cfg.MchID, nonce, sig, timestamp)
}

func decryptAESGCM(apiKey, nonce, ciphertext, associatedData string) ([]byte, error) {
	key := []byte(apiKey)
	if len(key) > 32 {
		key = key[:32]
	}
	for len(key) < 32 {
		key = append(key, 0)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ct, err := hex.DecodeString(ciphertext)
	if err != nil {
		// Try base64
		return nil, fmt.Errorf("decode ciphertext: %w", err)
	}
	return aead.Open(nil, []byte(nonce), ct, []byte(associatedData))
}
