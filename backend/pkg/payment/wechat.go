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
	"math/big"
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
	PublicKeyPath    string `yaml:"public_key_path"`
	PublicKeyID      string `yaml:"public_key_id"`
}

// WechatPay implements WeChat Pay V3 Native payment.
type WechatPay struct {
	Cfg WechatPayConfig
}

type WechatNotifyHeaders struct {
	Signature string
	Timestamp string
	Nonce     string
	Serial    string
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

	req, err := http.NewRequest("POST", "https://api.mch.weixin.qq.com/v3/pay/transactions/native", strings.NewReader(string(body)))
	if err != nil {
		return "", fmt.Errorf("wechat pay: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
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
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("wechat pay: invalid response")
	}
	if result.CodeURL == "" {
		return "", fmt.Errorf("wechat pay: empty code_url")
	}
	return result.CodeURL, nil
}

// VerifyNotify verifies WeChat Pay callback and extracts order number.
func (w *WechatPay) VerifyNotify(body []byte, headers WechatNotifyHeaders) (orderNo string, err error) {
	if !w.Cfg.Enabled {
		// Mock: try to parse order_id from body
		var mock struct {
			OrderNo string `json:"order_no"`
		}
		json.Unmarshal(body, &mock)
		return mock.OrderNo, nil
	}

	if err := w.verifyNotifySignature(body, headers); err != nil {
		return "", err
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
	if err := json.Unmarshal(plaintext, &order); err != nil {
		return "", fmt.Errorf("wechat notify: invalid resource")
	}
	if order.TradeState != "SUCCESS" {
		return "", fmt.Errorf("wechat notify: trade_state=%s", order.TradeState)
	}
	if order.OutTradeNo == "" {
		return "", fmt.Errorf("wechat notify: empty out_trade_no")
	}
	return order.OutTradeNo, nil
}

func (w *WechatPay) verifyNotifySignature(body []byte, headers WechatNotifyHeaders) error {
	if headers.Signature == "" || headers.Timestamp == "" || headers.Nonce == "" {
		return fmt.Errorf("wechat notify: missing signature headers")
	}
	publicKey, err := w.loadNotifyPublicKey(headers.Serial)
	if err != nil {
		return err
	}
	message := headers.Timestamp + "\n" + headers.Nonce + "\n" + string(body) + "\n"
	sigBytes, err := base64.StdEncoding.DecodeString(headers.Signature)
	if err != nil {
		return fmt.Errorf("wechat notify: invalid signature")
	}
	h := sha256.Sum256([]byte(message))
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, h[:], sigBytes); err != nil {
		return fmt.Errorf("wechat notify: verify signature: %w", err)
	}
	return nil
}

func (w *WechatPay) loadNotifyPublicKey(headerSerial string) (*rsa.PublicKey, error) {
	if w.Cfg.PublicKeyPath != "" {
		if w.Cfg.PublicKeyID != "" && headerSerial != "" && w.Cfg.PublicKeyID != headerSerial {
			return nil, fmt.Errorf("wechat notify: public key id mismatch")
		}
		return loadWechatPublicKey(w.Cfg.PublicKeyPath)
	}
	cert, err := loadWechatPlatformCert(w.Cfg.CertPath)
	if err != nil {
		return nil, err
	}
	if headerSerial != "" && !serialEqual(cert.SerialNumber, headerSerial) {
		return nil, fmt.Errorf("wechat notify: serial mismatch")
	}
	pub, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("wechat notify: platform cert is not RSA")
	}
	return pub, nil
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

func loadWechatPlatformCert(path string) (*x509.Certificate, error) {
	if path == "" {
		return nil, fmt.Errorf("wechat notify: cert_path is required")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("wechat notify: read platform cert: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("wechat notify: invalid platform cert pem")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("wechat notify: parse platform cert: %w", err)
	}
	return cert, nil
}

func loadWechatPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("wechat notify: read public key: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("wechat notify: invalid public key pem")
	}
	parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("wechat notify: parse public key: %w", err)
	}
	publicKey, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("wechat notify: public key is not RSA")
	}
	return publicKey, nil
}

func serialEqual(certSerial *big.Int, headerSerial string) bool {
	certValue := strings.TrimLeft(strings.ToUpper(certSerial.Text(16)), "0")
	headerValue := strings.TrimLeft(strings.ToUpper(headerSerial), "0")
	return certValue == headerValue
}

func decryptAESGCM(apiKey, nonce, ciphertext, associatedData string) ([]byte, error) {
	key := []byte(apiKey)
	if len(key) != 32 {
		return nil, fmt.Errorf("api_v3_key must be 32 bytes")
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
