package oss

import (
	"fmt"
	"io"
	"path"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// Config holds Alibaba Cloud OSS configuration.
type Config struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	Bucket          string `yaml:"bucket"`
	CDNDomain       string `yaml:"cdn_domain"` // e.g. https://img.zioran.com
	UploadPrefix    string `yaml:"upload_prefix"` // e.g. uploads/
}

// Client wraps the Alibaba Cloud OSS bucket operations.
type Client struct {
	cfg    Config
	bucket *oss.Bucket
}

// NewClient creates a new OSS client. Returns nil if config is incomplete.
func NewClient(cfg Config) (*Client, error) {
	if cfg.Endpoint == "" || cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" || cfg.Bucket == "" {
		return nil, fmt.Errorf("oss: incomplete config")
	}
	client, err := oss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
	if err != nil {
		return nil, fmt.Errorf("oss: create client: %w", err)
	}
	bucket, err := client.Bucket(cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("oss: get bucket: %w", err)
	}
	return &Client{cfg: cfg, bucket: bucket}, nil
}

// Upload uploads a file to OSS and returns the CDN URL.
// objectKey is the path within the bucket (e.g. "uploads/1719820000000.jpg").
func (c *Client) Upload(objectKey string, reader io.Reader) (string, error) {
	err := c.bucket.PutObject(objectKey, reader)
	if err != nil {
		return "", fmt.Errorf("oss: put object: %w", err)
	}
	return c.URL(objectKey), nil
}

// URL returns the CDN URL for a given object key.
func (c *Client) URL(objectKey string) string {
	if c.cfg.CDNDomain != "" {
		return c.cfg.CDNDomain + "/" + objectKey
	}
	return fmt.Sprintf("https://%s.%s/%s", c.cfg.Bucket, c.cfg.Endpoint, objectKey)
}

// GenerateObjectKey creates a unique object key with the configured prefix.
func (c *Client) GenerateObjectKey(ext string) string {
	prefix := c.cfg.UploadPrefix
	if prefix == "" {
		prefix = "uploads/"
	}
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	return path.Join(prefix, filename)
}

// Delete removes an object from the bucket.
func (c *Client) Delete(objectKey string) error {
	return c.bucket.DeleteObject(objectKey)
}
