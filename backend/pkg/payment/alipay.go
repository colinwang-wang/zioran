package payment

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

// AlipayConfig holds Alipay configuration.
type AlipayConfig struct {
	Enabled          bool   `yaml:"enabled"`
	MockAutoComplete bool   `yaml:"mock_auto_complete"`
	AppID            string `yaml:"app_id"`
	PrivateKey       string `yaml:"private_key"`
	AlipayPublicKey  string `yaml:"alipay_public_key"`
	NotifyURL        string `yaml:"notify_url"`
	ReturnURL        string `yaml:"return_url"`
}

// AlipayClient implements Alipay PC page payment.
type AlipayClient struct {
	Cfg AlipayConfig
}

func NewAlipayClient(cfg AlipayConfig) *AlipayClient {
	return &AlipayClient{Cfg: cfg}
}

// CreatePagePay creates a PC page payment URL.
func (a *AlipayClient) CreatePagePay(orderNo string, amount int, subject string) (payURL string, err error) {
	if !a.Cfg.Enabled {
		return "", fmt.Errorf("alipay: disabled")
	}

	bizContent, err := json.Marshal(map[string]string{
		"out_trade_no": orderNo,
		"total_amount": fmt.Sprintf("%.2f", float64(amount)/100),
		"subject":      subject,
		"product_code": "FAST_INSTANT_TRADE_PAY",
	})
	if err != nil {
		return "", fmt.Errorf("alipay biz_content: %w", err)
	}

	params := map[string]string{
		"app_id":      a.Cfg.AppID,
		"method":      "alipay.trade.page.pay",
		"charset":     "utf-8",
		"sign_type":   "RSA2",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"notify_url":  a.Cfg.NotifyURL,
		"biz_content": string(bizContent),
	}
	if a.Cfg.ReturnURL != "" {
		params["return_url"] = a.Cfg.ReturnURL
	}

	sign, err := a.signParams(params)
	if err != nil {
		return "", fmt.Errorf("alipay sign: %w", err)
	}
	params["sign"] = sign

	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	return "https://openapi.alipay.com/gateway.do?" + query.Encode(), nil
}

// VerifyNotify verifies Alipay async notification and returns order number.
func (a *AlipayClient) VerifyNotify(params map[string]string) (orderNo string, err error) {
	if !a.Cfg.Enabled {
		return params["out_trade_no"], nil
	}

	sign := params["sign"]
	delete(params, "sign")
	delete(params, "sign_type")

	// Build sign string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k + "=" + params[k])
	}

	// Verify RSA2 signature
	if err := a.verifySign(buf.String(), sign); err != nil {
		return "", fmt.Errorf("alipay verify: %w", err)
	}

	if params["trade_status"] != "TRADE_SUCCESS" && params["trade_status"] != "TRADE_FINISHED" {
		return "", fmt.Errorf("alipay: trade_status=%s", params["trade_status"])
	}
	return params["out_trade_no"], nil
}

func (a *AlipayClient) signParams(params map[string]string) (string, error) {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buf strings.Builder
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k + "=" + params[k])
	}

	key, err := parseAlipayPrivateKey(a.Cfg.PrivateKey)
	if err != nil {
		return "", err
	}

	h := sha256.Sum256([]byte(buf.String()))
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, h[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func (a *AlipayClient) verifySign(content, sign string) error {
	block, _ := pem.Decode([]byte(formatPublicKey(a.Cfg.AlipayPublicKey)))
	if block == nil {
		return fmt.Errorf("invalid public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	sigBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	h := sha256.Sum256([]byte(content))
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, h[:], sigBytes)
}

func parseAlipayPrivateKey(key string) (*rsa.PrivateKey, error) {
	var candidates []string
	if strings.Contains(key, "BEGIN") {
		candidates = append(candidates, key)
	} else {
		candidates = append(candidates,
			"-----BEGIN PRIVATE KEY-----\n"+key+"\n-----END PRIVATE KEY-----",
			"-----BEGIN RSA PRIVATE KEY-----\n"+key+"\n-----END RSA PRIVATE KEY-----",
		)
	}
	var lastErr error
	for _, candidate := range candidates {
		block, _ := pem.Decode([]byte(candidate))
		if block == nil {
			lastErr = fmt.Errorf("invalid private key")
			continue
		}
		parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err == nil {
			privateKey, ok := parsed.(*rsa.PrivateKey)
			if !ok {
				return nil, fmt.Errorf("private key is not RSA")
			}
			return privateKey, nil
		}
		lastErr = err
		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err == nil {
			return privateKey, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("parse private key: %w", lastErr)
}

func formatPublicKey(key string) string {
	if strings.Contains(key, "BEGIN") {
		return key
	}
	return "-----BEGIN PUBLIC KEY-----\n" + key + "\n-----END PUBLIC KEY-----"
}
