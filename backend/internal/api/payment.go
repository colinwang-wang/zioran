package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zioran/backend/internal/model"
	"github.com/zioran/backend/internal/service"
	"github.com/zioran/backend/pkg/errcode"
	"github.com/zioran/backend/pkg/response"
)

type PaymentHandler struct {
	paySvc *service.PaymentService
}

func NewPaymentHandler(paySvc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paySvc: paySvc}
}

// Coins

func (h *PaymentHandler) CoinBalance(c *gin.Context) {
	userID := c.GetInt64("user_id")
	result, err := h.paySvc.GetBalance(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *PaymentHandler) CoinTransactions(c *gin.Context) {
	userID := c.GetInt64("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	result, err := h.paySvc.GetTransactions(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *PaymentHandler) RechargeConfig(c *gin.Context) {
	result, err := h.paySvc.RechargeConfig(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *PaymentHandler) Recharge(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.RechargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.paySvc.Recharge(c.Request.Context(), userID, &req)
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

// VIP

func (h *PaymentHandler) VipPackages(c *gin.Context) {
	result, err := h.paySvc.VipPackages(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *PaymentHandler) VipStatus(c *gin.Context) {
	userID := c.GetInt64("user_id")
	result, err := h.paySvc.VipStatus(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *PaymentHandler) VipPurchase(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.VipPurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.paySvc.PurchaseVip(c.Request.Context(), userID, &req)
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

// Orders

func (h *PaymentHandler) CreateOrder(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	// Route to specific purchase flow based on type
	switch req.Type {
	case "course":
		result, err := h.paySvc.PurchaseCourse(c.Request.Context(), userID, int64(req.TargetID))
		if err != nil {
			if e, ok := err.(*errcode.Error); ok {
				response.Error(c, e)
				return
			}
			response.Error(c, errcode.ErrInternal)
			return
		}
		response.Success(c, result)
	default:
		response.Error(c, errcode.ErrParam)
	}
}

func (h *PaymentHandler) GetOrder(c *gin.Context) {
	userID := c.GetInt64("user_id")
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, svcErr := h.paySvc.GetOrder(c.Request.Context(), userID, orderID)
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

func (h *PaymentHandler) Download(c *gin.Context) {
	userID := c.GetInt64("user_id")
	courseID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, svcErr := h.paySvc.Download(c.Request.Context(), userID, courseID)
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

// User center

func (h *PaymentHandler) UserOrders(c *gin.Context) {
	userID := c.GetInt64("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	result, err := h.paySvc.UserOrders(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *PaymentHandler) UserDownloads(c *gin.Context) {
	userID := c.GetInt64("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	result, err := h.paySvc.UserDownloads(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *PaymentHandler) UserFavorites(c *gin.Context) {
	userID := c.GetInt64("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	result, err := h.paySvc.UserFavorites(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *PaymentHandler) AddFavorite(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.FavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if err := h.paySvc.AddFavorite(c.Request.Context(), userID, req.CourseID); err != nil {
		if e, ok := err.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

func (h *PaymentHandler) RemoveFavorite(c *gin.Context) {
	userID := c.GetInt64("user_id")
	courseID, err := strconv.ParseInt(c.Param("courseId"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.paySvc.RemoveFavorite(c.Request.Context(), userID, courseID); svcErr != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

func (h *PaymentHandler) ChangePassword(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if err := h.paySvc.ChangePassword(c.Request.Context(), userID, &req); err != nil {
		if e, ok := err.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}
