package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zioran/backend/pkg/errcode"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "ok", Data: data})
}

func Error(c *gin.Context, err *errcode.Error) {
	c.JSON(http.StatusOK, Response{Code: err.Code, Message: err.Message})
}

func ErrorMsg(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{Code: code, Message: msg})
}
