package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zioran/backend/internal/model"
	"github.com/zioran/backend/internal/service"
	"github.com/zioran/backend/pkg/errcode"
	"github.com/zioran/backend/pkg/response"
)

type CommunityHandler struct {
	commSvc *service.CommunityService
}

func NewCommunityHandler(commSvc *service.CommunityService) *CommunityHandler {
	return &CommunityHandler{commSvc: commSvc}
}

// Guestbook

func (h *CommunityHandler) GuestbookList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	var userID int64
	if id, exists := c.Get("user_id"); exists {
		userID = id.(int64)
	}
	result, err := h.commSvc.GuestbookList(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *CommunityHandler) GuestbookCreate(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.GuestbookCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.commSvc.GuestbookCreate(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *CommunityHandler) GuestbookLike(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	liked, svcErr := h.commSvc.GuestbookLike(c.Request.Context(), userID, id)
	if svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, gin.H{"liked": liked})
}

func (h *CommunityHandler) GuestbookDelete(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.commSvc.GuestbookDelete(c.Request.Context(), userID, id, false); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

// Comments

func (h *CommunityHandler) CommentList(c *gin.Context) {
	targetType := c.Query("targetType")
	targetID, _ := strconv.ParseInt(c.Query("targetId"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if targetType == "" || targetID == 0 {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.commSvc.CommentList(c.Request.Context(), targetType, targetID, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *CommunityHandler) CommentCreate(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.commSvc.CommentCreate(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *CommunityHandler) CommentDelete(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.commSvc.CommentDelete(c.Request.Context(), userID, id, false); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

// Home config (public)

func (h *CommunityHandler) NavItems(c *gin.Context) {
	result, err := h.commSvc.NavItems(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *CommunityHandler) Banners(c *gin.Context) {
	result, err := h.commSvc.Banners(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}
