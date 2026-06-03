package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zioran/backend/internal/model"
	"github.com/zioran/backend/internal/service"
	"github.com/zioran/backend/pkg/errcode"
	"github.com/zioran/backend/pkg/response"
)

type AdminPaymentHandler struct {
	paySvc  *service.PaymentService
	commSvc *service.CommunityService
}

func NewAdminPaymentHandler(paySvc *service.PaymentService, commSvc *service.CommunityService) *AdminPaymentHandler {
	return &AdminPaymentHandler{paySvc: paySvc, commSvc: commSvc}
}

// Orders

func (h *AdminPaymentHandler) OrderList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	status := c.Query("status")
	result, err := h.paySvc.AdminOrders(c.Request.Context(), page, pageSize, status)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AdminPaymentHandler) OrderRefund(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.paySvc.AdminRefund(c.Request.Context(), id); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

// Users

func (h *AdminPaymentHandler) UserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	keyword := c.Query("keyword")
	result, err := h.paySvc.AdminUsers(c.Request.Context(), page, pageSize, keyword)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AdminPaymentHandler) UserUpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	var req model.AdminUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.paySvc.AdminUpdateUserStatus(c.Request.Context(), id, req.Status); svcErr != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

func (h *AdminPaymentHandler) UserRecharge(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	var req model.AdminRechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.paySvc.AdminRecharge(c.Request.Context(), id, &req); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

// Dashboard

func (h *AdminPaymentHandler) DashboardStats(c *gin.Context) {
	result, err := h.paySvc.DashboardStats(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AdminPaymentHandler) DashboardCharts(c *gin.Context) {
	period := c.DefaultQuery("period", "week")
	result := h.paySvc.DashboardCharts(c.Request.Context(), period)
	response.Success(c, result)
}

func (h *AdminPaymentHandler) UserDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	user, svcErr := h.paySvc.AdminGetUser(c.Request.Context(), id)
	if svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, user)
}

// Guestbook admin

func (h *AdminPaymentHandler) GuestbookList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	result, err := h.commSvc.AdminGuestbookList(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AdminPaymentHandler) GuestbookUpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	var req model.AdminGuestbookStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.commSvc.AdminGuestbookUpdateStatus(c.Request.Context(), id, req.Status); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

func (h *AdminPaymentHandler) GuestbookPin(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	// Toggle pin
	if svcErr := h.commSvc.AdminGuestbookPin(c.Request.Context(), id, true); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

func (h *AdminPaymentHandler) GuestbookDelete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.commSvc.GuestbookDelete(c.Request.Context(), 0, id, true); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

// Comments admin

func (h *AdminPaymentHandler) CommentList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	result, err := h.commSvc.AdminCommentList(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AdminPaymentHandler) CommentUpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	var req model.AdminCommentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.commSvc.AdminCommentUpdateStatus(c.Request.Context(), id, req.Status); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

func (h *AdminPaymentHandler) CommentDelete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.commSvc.CommentDelete(c.Request.Context(), 0, id, true); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

// Nav Items admin

func (h *AdminPaymentHandler) NavItemList(c *gin.Context) {
	result, err := h.commSvc.AdminNavItemList(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AdminPaymentHandler) NavItemCreate(c *gin.Context) {
	var req model.NavItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.commSvc.AdminNavItemCreate(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AdminPaymentHandler) NavItemUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	var req model.NavItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, svcErr := h.commSvc.AdminNavItemUpdate(c.Request.Context(), id, &req)
	if svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AdminPaymentHandler) NavItemDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.commSvc.AdminNavItemDelete(c.Request.Context(), id); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

// Banners admin

func (h *AdminPaymentHandler) BannerList(c *gin.Context) {
	result, err := h.commSvc.AdminBannerList(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AdminPaymentHandler) BannerCreate(c *gin.Context) {
	var req model.BannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.commSvc.AdminBannerCreate(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AdminPaymentHandler) BannerUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	var req model.BannerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, svcErr := h.commSvc.AdminBannerUpdate(c.Request.Context(), id, &req)
	if svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *AdminPaymentHandler) BannerDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.commSvc.AdminBannerDelete(c.Request.Context(), id); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}
