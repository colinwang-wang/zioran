package payment

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAlipayCreatePagePayUsesReturnURLAndRawPKCS8Key(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	keyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		t.Fatalf("marshal key: %v", err)
	}

	client := NewAlipayClient(AlipayConfig{
		Enabled:    true,
		AppID:      "app-1",
		PrivateKey: base64.StdEncoding.EncodeToString(keyDER),
		NotifyURL:  "https://api.zioran.com/api/v1/pay/notify/alipay",
		ReturnURL:  "https://www.zioran.com/user/transactions",
	})
	payURL, err := client.CreatePagePay("ORD123", 1234, "测试订单")
	if err != nil {
		t.Fatalf("create page pay: %v", err)
	}

	parsed, err := url.Parse(payURL)
	if err != nil {
		t.Fatalf("parse pay url: %v", err)
	}
	query := parsed.Query()
	if query.Get("return_url") != "https://www.zioran.com/user/transactions" {
		t.Fatalf("unexpected return_url: %s", query.Get("return_url"))
	}
	if query.Get("notify_url") != "https://api.zioran.com/api/v1/pay/notify/alipay" {
		t.Fatalf("unexpected notify_url: %s", query.Get("notify_url"))
	}
	if query.Get("sign") == "" {
		t.Fatal("missing sign")
	}
}

func TestWechatVerifyNotifyChecksSignatureAndDecrypts(t *testing.T) {
	platformKey, certPath := createPlatformCert(t)
	apiKey := "12345678901234567890123456789012"
	resourceBody := []byte(`{"out_trade_no":"ORD123","trade_state":"SUCCESS"}`)
	ciphertext, resourceNonce := encryptWechatResource(t, apiKey, "transaction", resourceBody)
	body, err := json.Marshal(map[string]interface{}{
		"resource": map[string]string{
			"ciphertext":      ciphertext,
			"nonce":           resourceNonce,
			"associated_data": "transaction",
		},
	})
	if err != nil {
		t.Fatalf("marshal notify: %v", err)
	}

	timestamp := "1710000000"
	nonce := "notify-nonce"
	signature := signWechatNotify(t, platformKey, timestamp+"\n"+nonce+"\n"+string(body)+"\n")
	client := NewWechatPay(WechatPayConfig{
		Enabled:  true,
		APIKey:   apiKey,
		CertPath: certPath,
	})

	orderNo, err := client.VerifyNotify(body, WechatNotifyHeaders{
		Signature: signature,
		Timestamp: timestamp,
		Nonce:     nonce,
		Serial:    "1",
	})
	if err != nil {
		t.Fatalf("verify notify: %v", err)
	}
	if orderNo != "ORD123" {
		t.Fatalf("unexpected order no: %s", orderNo)
	}

	_, err = client.VerifyNotify(body, WechatNotifyHeaders{
		Signature: signature[:len(signature)-2] + "xx",
		Timestamp: timestamp,
		Nonce:     nonce,
		Serial:    "1",
	})
	if err == nil {
		t.Fatal("expected invalid signature error")
	}
}

func TestWechatVerifyNotifyWithPublicKey(t *testing.T) {
	platformKey, publicKeyPath := createPlatformPublicKey(t)
	apiKey := "12345678901234567890123456789012"
	resourceBody := []byte(`{"out_trade_no":"ORD456","trade_state":"SUCCESS"}`)
	ciphertext, resourceNonce := encryptWechatResource(t, apiKey, "transaction", resourceBody)
	body, err := json.Marshal(map[string]interface{}{
		"resource": map[string]string{
			"ciphertext":      ciphertext,
			"nonce":           resourceNonce,
			"associated_data": "transaction",
		},
	})
	if err != nil {
		t.Fatalf("marshal notify: %v", err)
	}

	timestamp := "1710000000"
	nonce := "notify-nonce"
	signature := signWechatNotify(t, platformKey, timestamp+"\n"+nonce+"\n"+string(body)+"\n")
	client := NewWechatPay(WechatPayConfig{
		Enabled:       true,
		APIKey:        apiKey,
		PublicKeyPath: publicKeyPath,
		PublicKeyID:   "PUB_KEY_ID_1",
	})

	orderNo, err := client.VerifyNotify(body, WechatNotifyHeaders{
		Signature: signature,
		Timestamp: timestamp,
		Nonce:     nonce,
		Serial:    "PUB_KEY_ID_1",
	})
	if err != nil {
		t.Fatalf("verify notify: %v", err)
	}
	if orderNo != "ORD456" {
		t.Fatalf("unexpected order no: %s", orderNo)
	}

	_, err = client.VerifyNotify(body, WechatNotifyHeaders{
		Signature: signature,
		Timestamp: timestamp,
		Nonce:     nonce,
		Serial:    "PUB_KEY_ID_2",
	})
	if err == nil {
		t.Fatal("expected public key id mismatch")
	}
}

func createPlatformCert(t *testing.T) (*rsa.PrivateKey, string) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	tpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "Wechat Pay Platform"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	certDER, err := x509.CreateCertificate(rand.Reader, tpl, tpl, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("create cert: %v", err)
	}
	certPath := filepath.Join(t.TempDir(), "wechat-platform.pem")
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	if err := os.WriteFile(certPath, certPEM, 0600); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	return privateKey, certPath
}

func createPlatformPublicKey(t *testing.T) (*rsa.PrivateKey, string) {
	t.Helper()
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	publicKeyDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatalf("marshal public key: %v", err)
	}
	publicKeyPath := filepath.Join(t.TempDir(), "wechat-public-key.pem")
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicKeyDER})
	if err := os.WriteFile(publicKeyPath, publicKeyPEM, 0600); err != nil {
		t.Fatalf("write public key: %v", err)
	}
	return privateKey, publicKeyPath
}

func encryptWechatResource(t *testing.T, apiKey, associatedData string, plaintext []byte) (string, string) {
	t.Helper()
	block, err := aes.NewCipher([]byte(apiKey))
	if err != nil {
		t.Fatalf("cipher: %v", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatalf("gcm: %v", err)
	}
	nonce := "123456789012"
	ciphertext := aead.Seal(nil, []byte(nonce), plaintext, []byte(associatedData))
	return base64.StdEncoding.EncodeToString(ciphertext), nonce
}

func signWechatNotify(t *testing.T, privateKey *rsa.PrivateKey, message string) string {
	t.Helper()
	hash := sha256.Sum256([]byte(message))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	return base64.StdEncoding.EncodeToString(signature)
}
