package oss

import (
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// Config holds Alibaba Cloud OSS configuration.
type Config struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	Bucket          string `yaml:"bucket"`
	CDNDomain       string `yaml:"cdn_domain"`    // e.g. https://img.zioran.com
	UploadPrefix    string `yaml:"upload_prefix"` // e.g. uploads/
}

// Client wraps the Alibaba Cloud OSS bucket operations.
type Client struct {
	cfg    Config
	bucket *oss.Bucket
}

// NewClient creates a new OSS client. Returns error if config is incomplete.
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
func (c *Client) Upload(objectKey string, reader io.Reader) (string, error) {
	err := c.bucket.PutObject(objectKey, reader, oss.ObjectACL(oss.ACLPublicRead))
	if err != nil {
		return "", fmt.Errorf("oss: put object: %w", err)
	}
	return c.URL(objectKey), nil
}

// URL returns the public CDN URL for a given object key.
func (c *Client) URL(objectKey string) string {
	if c.cfg.CDNDomain != "" {
		return strings.TrimRight(c.cfg.CDNDomain, "/") + "/" + objectKey
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

// DeleteMultiple removes multiple objects from the bucket.
func (c *Client) DeleteMultiple(objectKeys []string) error {
	if len(objectKeys) == 0 {
		return nil
	}
	_, err := c.bucket.DeleteObjects(objectKeys)
	return err
}

// --- 图片处理 (阿里云 OSS 图片处理服务) ---

// ThumbnailURL returns a URL with OSS image resize (fit mode, max width/height).
// Example output: https://img.zioran.com/uploads/xxx.jpg?x-oss-process=image/resize,m_lfit,w_300,h_300
func (c *Client) ThumbnailURL(objectKey string, width, height int) string {
	process := fmt.Sprintf("image/resize,m_lfit,w_%d,h_%d", width, height)
	return c.processURL(objectKey, process)
}

// CropURL returns a URL with OSS image center-crop to exact dimensions.
// Example output: https://img.zioran.com/uploads/xxx.jpg?x-oss-process=image/resize,m_fill,w_300,h_200
func (c *Client) CropURL(objectKey string, width, height int) string {
	process := fmt.Sprintf("image/resize,m_fill,w_%d,h_%d", width, height)
	return c.processURL(objectKey, process)
}

// WebpURL returns a URL that converts image to webp format (with optional resize).
// If width/height are 0, only format conversion is applied.
func (c *Client) WebpURL(objectKey string, width, height int) string {
	var process string
	if width > 0 && height > 0 {
		process = fmt.Sprintf("image/resize,m_lfit,w_%d,h_%d/format,webp", width, height)
	} else {
		process = "image/format,webp"
	}
	return c.processURL(objectKey, process)
}

// processURL appends x-oss-process query parameter to the object URL.
func (c *Client) processURL(objectKey, process string) string {
	base := c.URL(objectKey)
	sep := "?"
	if strings.Contains(base, "?") {
		sep = "&"
	}
	return base + sep + "x-oss-process=" + process
}

// --- 签名 URL (私有资源临时访问) ---

// SignURL generates a pre-signed URL for private object download.
// The URL is valid for the specified duration.
func (c *Client) SignURL(objectKey string, expires time.Duration) (string, error) {
	expireSeconds := int64(expires.Seconds())
	if expireSeconds < 1 {
		expireSeconds = 3600 // default 1 hour
	}
	signedURL, err := c.bucket.SignURL(objectKey, oss.HTTPGet, expireSeconds)
	if err != nil {
		return "", fmt.Errorf("oss: sign url: %w", err)
	}
	return signedURL, nil
}

// --- 工具方法 ---

// ObjectKeyFromURL extracts the object key from a full OSS/CDN URL.
// Returns empty string if the URL doesn't match the configured domain.
func (c *Client) ObjectKeyFromURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}

	// Try CDN domain first
	if c.cfg.CDNDomain != "" {
		cdnBase := strings.TrimRight(c.cfg.CDNDomain, "/") + "/"
		if strings.HasPrefix(rawURL, cdnBase) {
			key := strings.TrimPrefix(rawURL, cdnBase)
			// Remove query parameters
			if idx := strings.Index(key, "?"); idx != -1 {
				key = key[:idx]
			}
			return key
		}
	}

	// Try OSS direct URL
	ossBase := fmt.Sprintf("https://%s.%s/", c.cfg.Bucket, c.cfg.Endpoint)
	if strings.HasPrefix(rawURL, ossBase) {
		key := strings.TrimPrefix(rawURL, ossBase)
		if idx := strings.Index(key, "?"); idx != -1 {
			key = key[:idx]
		}
		return key
	}

	// Try parsing as URL and extract path
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	// Check if host matches
	ossHost := c.cfg.Bucket + "." + c.cfg.Endpoint
	if parsed.Host == ossHost || (c.cfg.CDNDomain != "" && strings.Contains(c.cfg.CDNDomain, parsed.Host)) {
		return strings.TrimPrefix(parsed.Path, "/")
	}

	return ""
}

// IsOSSURL checks whether the given URL belongs to this OSS bucket (either CDN or direct).
func (c *Client) IsOSSURL(rawURL string) bool {
	return c.ObjectKeyFromURL(rawURL) != ""
}
