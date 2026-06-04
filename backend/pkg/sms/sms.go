package sms

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Sender sends SMS verification codes.
type Sender interface {
	Send(phone, code string) error
}

// SMSConfig holds provider selection and credentials.
type SMSConfig struct {
	Provider string       `yaml:"provider"` // mock | aliyun | tencent
	Aliyun   AliyunConfig `yaml:"aliyun"`
	Tencent  TencentConfig `yaml:"tencent"`
}

type AliyunConfig struct {
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	SignName         string `yaml:"sign_name"`
	TemplateCode    string `yaml:"template_code"`
}

type TencentConfig struct {
	SecretId   string `yaml:"secret_id"`
	SecretKey  string `yaml:"secret_key"`
	AppId      string `yaml:"app_id"`
	SignName   string `yaml:"sign_name"`
	TemplateId string `yaml:"template_id"`
}

// NewSender creates a Sender based on config.
func NewSender(cfg SMSConfig) Sender {
	switch cfg.Provider {
	case "aliyun":
		return &AliyunSender{cfg: cfg.Aliyun}
	case "tencent":
		return &TencentSender{cfg: cfg.Tencent}
	default:
		return &MockSender{}
	}
}

// MockSender prints SMS to console.
type MockSender struct{}

func (s *MockSender) Send(phone, code string) error {
	fmt.Printf("[SMS MOCK] 手机号: %s, 验证码: %s\n", phone, code)
	return nil
}

// AliyunSender sends SMS via Alibaba Cloud.
type AliyunSender struct {
	cfg AliyunConfig
}

func (s *AliyunSender) Send(phone, code string) error {
	params := map[string]string{
		"AccessKeyId":      s.cfg.AccessKeyID,
		"Action":           "SendSms",
		"Format":           "JSON",
		"PhoneNumbers":     phone,
		"SignName":         s.cfg.SignName,
		"SignatureMethod":  "HMAC-SHA1",
		"SignatureNonce":   strconv.FormatInt(time.Now().UnixNano(), 10),
		"SignatureVersion": "1.0",
		"TemplateCode":     s.cfg.TemplateCode,
		"TemplateParam":    fmt.Sprintf(`{"code":"%s"}`, code),
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"Version":          "2017-05-25",
	}
	// Build signature
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var query strings.Builder
	for i, k := range keys {
		if i > 0 {
			query.WriteByte('&')
		}
		query.WriteString(url.QueryEscape(k) + "=" + url.QueryEscape(params[k]))
	}
	stringToSign := "GET&" + url.QueryEscape("/") + "&" + url.QueryEscape(query.String())
	mac := hmac.New(sha256.New, []byte(s.cfg.AccessKeySecret+"&"))
	mac.Write([]byte(stringToSign))
	sig := url.QueryEscape(hex.EncodeToString(mac.Sum(nil)))

	reqURL := "https://dysmsapi.aliyuncs.com/?" + query.String() + "&Signature=" + sig
	resp, err := http.Get(reqURL)
	if err != nil {
		return fmt.Errorf("aliyun sms: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Code    string `json:"Code"`
		Message string `json:"Message"`
	}
	json.Unmarshal(body, &result)
	if result.Code != "OK" {
		return fmt.Errorf("aliyun sms: %s - %s", result.Code, result.Message)
	}
	return nil
}

// TencentSender sends SMS via Tencent Cloud.
type TencentSender struct {
	cfg TencentConfig
}

func (s *TencentSender) Send(phone, code string) error {
	host := "sms.tencentcloudapi.com"
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	date := time.Now().UTC().Format("2006-01-02")

	payload := fmt.Sprintf(`{"PhoneNumberSet":["+86%s"],"SmsSdkAppId":"%s","SignName":"%s","TemplateId":"%s","TemplateParamSet":["%s"]}`,
		phone, s.cfg.AppId, s.cfg.SignName, s.cfg.TemplateId, code)

	// TC3-HMAC-SHA256 signing
	hashedPayload := sha256Hex([]byte(payload))
	canonicalRequest := "POST\n/\n\ncontent-type:application/json\nhost:" + host + "\n\ncontent-type;host\n" + hashedPayload
	credentialScope := date + "/sms/tc3_request"
	stringToSign := "TC3-HMAC-SHA256\n" + timestamp + "\n" + credentialScope + "\n" + sha256Hex([]byte(canonicalRequest))

	secretDate := hmacSHA256([]byte("TC3"+s.cfg.SecretKey), []byte(date))
	secretService := hmacSHA256(secretDate, []byte("sms"))
	secretSigning := hmacSHA256(secretService, []byte("tc3_request"))
	signature := hex.EncodeToString(hmacSHA256(secretSigning, []byte(stringToSign)))

	auth := fmt.Sprintf("TC3-HMAC-SHA256 Credential=%s/%s, SignedHeaders=content-type;host, Signature=%s",
		s.cfg.SecretId, credentialScope, signature)

	req, _ := http.NewRequest("POST", "https://"+host, strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", host)
	req.Header.Set("X-TC-Action", "SendSms")
	req.Header.Set("X-TC-Version", "2021-01-11")
	req.Header.Set("X-TC-Timestamp", timestamp)
	req.Header.Set("Authorization", auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("tencent sms: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Response struct {
			Error *struct {
				Code    string `json:"Code"`
				Message string `json:"Message"`
			} `json:"Error"`
		} `json:"Response"`
	}
	json.Unmarshal(body, &result)
	if result.Response.Error != nil {
		return fmt.Errorf("tencent sms: %s - %s", result.Response.Error.Code, result.Response.Error.Message)
	}
	return nil
}

func sha256Hex(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func hmacSHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}
