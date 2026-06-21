package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zioran/backend/internal/model"
	"github.com/zioran/backend/internal/service"
	"github.com/zioran/backend/pkg/errcode"
	"github.com/zioran/backend/pkg/response"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Captcha(c *gin.Context) {
	result, err := h.authSvc.GenerateCaptcha()
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AuthHandler) SendEmail(c *gin.Context) {
	var req model.SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if err := h.authSvc.SendEmail(c.Request.Context(), req.Email, req.CaptchaKey, req.Captcha); err != nil {
		if e, ok := err.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.authSvc.Register(c.Request.Context(), &req)
	if err != nil {
		if e, ok := err.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.authSvc.Login(c.Request.Context(), &req)
	if err != nil {
		if e, ok := err.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AuthHandler) Profile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}
	result, err := h.authSvc.GetProfile(c.Request.Context(), userID.(int64))
	if err != nil {
		if e, ok := err.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}
	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.authSvc.UpdateProfile(c.Request.Context(), userID.(int64), &req)
	if err != nil {
		if e, ok := err.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req model.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.authSvc.AdminLogin(c.Request.Context(), &req)
	if err != nil {
		if e, ok := err.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}
