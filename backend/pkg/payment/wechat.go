package payment

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// WechatPayConfig holds WeChat Pay configuration.
type WechatPayConfig struct {
	Enabled          bool   `yaml:"enabled"`
	MockAutoComplete bool   `yaml:"mock_auto_complete"`
	AppID            string `yaml:"app_id"`
	MchID            string `yaml:"mch_id"`
	APIKey           string `yaml:"api_key"`
	NotifyURL        string `yaml:"notify_url"`
	SerialNo         string `yaml:"serial_no"`
	PrivateKeyPath   string `yaml:"private_key_path"`
	CertPath         string `yaml:"cert_path"`
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
		return "", fmt.Errorf("wechat pay: disabled")
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
	auth, err := w.buildAuth("POST", "/v3/pay/transactions/native", string(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", auth)

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

func (w *WechatPay) buildAuth(method, path, body string) (string, error) {
	if w.Cfg.SerialNo == "" {
		return "", fmt.Errorf("wechat pay: serial_no is required")
	}
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := fmt.Sprintf("%d", time.Now().UnixNano())
	message := method + "\n" + path + "\n" + timestamp + "\n" + nonce + "\n" + body + "\n"
	privateKey, err := loadWechatPrivateKey(w.Cfg.PrivateKeyPath)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256([]byte(message))
	sigBytes, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, h[:])
	if err != nil {
		return "", err
	}
	sig := base64.StdEncoding.EncodeToString(sigBytes)
	return fmt.Sprintf(`WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",signature="%s",timestamp="%s",serial_no="%s"`,
		w.Cfg.MchID, nonce, sig, timestamp, w.Cfg.SerialNo), nil
}

func loadWechatPrivateKey(path string) (*rsa.PrivateKey, error) {
	if path == "" {
		return nil, fmt.Errorf("wechat pay: private_key_path is required")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("wechat pay: read private key: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("wechat pay: invalid private key pem")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		if privateKey, ok := key.(*rsa.PrivateKey); ok {
			return privateKey, nil
		}
		return nil, fmt.Errorf("wechat pay: private key is not RSA")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("wechat pay: parse private key: %w", err)
	}
	return privateKey, nil
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
	ct, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("decode ciphertext: %w", err)
	}
	return aead.Open(nil, []byte(nonce), ct, []byte(associatedData))
}
