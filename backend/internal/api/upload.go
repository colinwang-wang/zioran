package api

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zioran/backend/pkg/errcode"
	"github.com/zioran/backend/pkg/response"
)

type UploadHandler struct {
	uploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
	ensureUploadDir(uploadDir)
	return &UploadHandler{uploadDir: uploadDir}
}

func (h *UploadHandler) ImageUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	url, saveErr := saveUploadedImage(c, h.uploadDir, file)
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

func saveUploadedImage(c *gin.Context, uploadDir string, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExts[ext] {
		return "", errcode.New(40001, "仅支持 jpg、png、webp、gif 图片")
	}
	if file.Size > maxImageUploadSize {
		return "", errcode.New(40001, "图片大小不能超过 5MB")
	}
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
