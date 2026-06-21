package email

import (
	"crypto/tls"
	"fmt"
	"io"
	"mime"
	"net/smtp"
	"strings"
)

type Sender interface {
	Send(to, code string) error
}

type EmailConfig struct {
	Provider string     `yaml:"provider"` // mock | smtp
	SMTP     SMTPConfig `yaml:"smtp"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
	Subject  string `yaml:"subject"`
}

func NewSender(cfg EmailConfig) Sender {
	switch cfg.Provider {
	case "smtp":
		return &SMTPSender{cfg: cfg.SMTP}
	default:
		return &MockSender{}
	}
}

type MockSender struct{}

func (s *MockSender) Send(to, code string) error {
	fmt.Printf("[EMAIL MOCK] 邮箱: %s, 验证码: %s\n", to, code)
	return nil
}

type SMTPSender struct {
	cfg SMTPConfig
}

func (s *SMTPSender) Send(to, code string) error {
	if s.cfg.Host == "" || s.cfg.Port == 0 || s.cfg.From == "" {
		return fmt.Errorf("smtp email: incomplete config")
	}
	subject := s.cfg.Subject
	if subject == "" {
		subject = "知猿验证码"
	}
	body := fmt.Sprintf("您的验证码是：%s，5分钟内有效。如非本人操作，请忽略本邮件。", code)
	message := strings.Join([]string{
		"From: " + s.cfg.From,
		"To: " + to,
		"Subject: " + mime.QEncoding.Encode("UTF-8", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	var auth smtp.Auth
	if s.cfg.Username != "" || s.cfg.Password != "" {
		auth = smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	}
	if s.cfg.Port == 465 {
		return sendMailTLS(addr, s.cfg.Host, auth, s.cfg.From, []string{to}, []byte(message))
	}
	return smtp.SendMail(addr, auth, s.cfg.From, []string{to}, []byte(message))
}

func sendMailTLS(addr, host string, auth smtp.Auth, from string, to []string, msg []byte) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12})
	if err != nil {
		return fmt.Errorf("smtp tls: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Close()

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("smtp from: %w", err)
	}
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("smtp rcpt: %w", err)
		}
	}
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := writer.Write(msg); err != nil {
		writer.Close()
		return fmt.Errorf("smtp write: %w", err)
	}
	if err := writer.Close(); err != nil && err != io.EOF {
		return fmt.Errorf("smtp close: %w", err)
	}
	return client.Quit()
}
