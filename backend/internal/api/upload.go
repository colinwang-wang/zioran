package api

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	ossClient "github.com/zioran/backend/pkg/oss"
	"github.com/zioran/backend/pkg/errcode"
	"github.com/zioran/backend/pkg/response"
)

type UploadHandler struct {
	uploadDir string
	oss       *ossClient.Client
}

func NewUploadHandler(uploadDir string, ossClient *ossClient.Client) *UploadHandler {
	if ossClient == nil {
		ensureUploadDir(uploadDir)
	}
	return &UploadHandler{uploadDir: uploadDir, oss: ossClient}
}

func (h *UploadHandler) ImageUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	url, saveErr := h.saveImage(c, file)
	if saveErr != nil {
		if e, ok := saveErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, gin.H{"url": url})
}

const maxImageUploadSize = 5 << 20

var allowedImageExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
}

func (h *UploadHandler) saveImage(c *gin.Context, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExts[ext] {
		return "", errcode.New(40001, "仅支持 jpg、png、webp、gif 图片")
	}
	if file.Size > maxImageUploadSize {
		return "", errcode.New(40001, "图片大小不能超过 5MB")
	}

	// 如果配置了 OSS，上传到 OSS
	if h.oss != nil {
		src, err := file.Open()
		if err != nil {
			return "", err
		}
		defer src.Close()

		objectKey := h.oss.GenerateObjectKey(ext)
		url, err := h.oss.Upload(objectKey, src)
		if err != nil {
			return "", err
		}
		return url, nil
	}

	// 否则保存到本地（兼容旧逻辑）
	return saveUploadedImageLocal(c, h.uploadDir, file)
}

func saveUploadedImageLocal(c *gin.Context, uploadDir string, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if err := ensureUploadDir(uploadDir); err != nil {
		return "", err
	}
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}
	if err := os.Chmod(dst, 0644); err != nil {
		return "", err
	}
	return "/uploads/" + filename, nil
}

func ensureUploadDir(uploadDir string) error {
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return err
	}
	return os.Chmod(uploadDir, 0755)
}
