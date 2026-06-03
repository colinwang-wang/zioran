package api

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zioran/backend/pkg/errcode"
	"github.com/zioran/backend/pkg/response"
)

type UploadHandler struct {
	uploadDir string
}

func NewUploadHandler(uploadDir string) *UploadHandler {
	os.MkdirAll(uploadDir, 0755)
	return &UploadHandler{uploadDir: uploadDir}
}

func (h *UploadHandler) ImageUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(h.uploadDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	url := "/uploads/" + filename
	response.Success(c, gin.H{"url": url})
}
