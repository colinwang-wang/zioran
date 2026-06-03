package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zioran/backend/internal/middleware"
	"github.com/zioran/backend/internal/model"
	"github.com/zioran/backend/internal/service"
	"github.com/zioran/backend/pkg/errcode"
	"github.com/zioran/backend/pkg/response"
)

type TicketHandler struct {
	ticketSvc *service.TicketService
	authSvc   *service.AuthService
	jwtSecret string
	jwtExpire time.Duration
	uploadDir string
}

func NewTicketHandler(ticketSvc *service.TicketService, authSvc *service.AuthService, jwtSecret string, jwtExpire time.Duration, uploadDir string) *TicketHandler {
	return &TicketHandler{ticketSvc: ticketSvc, authSvc: authSvc, jwtSecret: jwtSecret, jwtExpire: jwtExpire, uploadDir: uploadDir}
}

// === Auth extensions ===

// ForgotPassword resets user password via SMS code
func (h *TicketHandler) ForgotPassword(c *gin.Context) {
	var req model.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	// Verify SMS code via AuthService
	if !h.authSvc.VerifySMSCode(req.Phone, req.SMSCode) {
		response.Error(c, errcode.New(40001, "短信验证码错误"))
		return
	}
	if err := h.ticketSvc.ForgotPassword(c.Request.Context(), req.Phone, req.NewPassword); err != nil {
		if e, ok := err.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

// RefreshToken generates a new JWT from a valid existing token
func (h *TicketHandler) RefreshToken(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, errcode.ErrUnauthorized)
		return
	}
	token, err := middleware.GenerateToken(userID.(int64), h.jwtSecret, h.jwtExpire)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, gin.H{"token": token})
}

// === Order cancel ===

func (h *TicketHandler) CancelOrder(c *gin.Context) {
	userID := c.GetInt64("user_id")
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.ticketSvc.CancelOrder(c.Request.Context(), userID, orderID); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

// === User order detail ===
// Note: /api/v1/user/orders/:id is mapped to PaymentHandler.GetOrder in router

// === Tickets (user) ===

func (h *TicketHandler) TicketList(c *gin.Context) {
	userID := c.GetInt64("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	result, err := h.ticketSvc.List(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *TicketHandler) TicketCreate(c *gin.Context) {
	userID := c.GetInt64("user_id")
	var req model.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.ticketSvc.Create(c.Request.Context(), userID, &req)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *TicketHandler) TicketDetail(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, svcErr := h.ticketSvc.Detail(c.Request.Context(), userID, id, false)
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

func (h *TicketHandler) TicketReply(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	var req model.TicketReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.ticketSvc.Reply(c.Request.Context(), userID, id, req.Content, false); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

// === Tickets (admin) ===

func (h *TicketHandler) AdminTicketList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	result, err := h.ticketSvc.AdminList(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *TicketHandler) AdminTicketDetail(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, svcErr := h.ticketSvc.Detail(c.Request.Context(), userID, id, true)
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

func (h *TicketHandler) AdminTicketUpdateStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	var req model.TicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.ticketSvc.UpdateStatus(c.Request.Context(), id, req.Status); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

func (h *TicketHandler) AdminTicketReply(c *gin.Context) {
	userID := c.GetInt64("user_id")
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	var req model.TicketReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.ticketSvc.Reply(c.Request.Context(), userID, id, req.Content, true); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

// === Settings ===

func (h *TicketHandler) GetSettings(c *gin.Context) {
	result, err := h.ticketSvc.GetSettings(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *TicketHandler) UpdateSettings(c *gin.Context) {
	var req model.SettingsMap
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if err := h.ticketSvc.UpdateSettings(c.Request.Context(), req); err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

// === Admin account management ===

func (h *TicketHandler) AdminList(c *gin.Context) {
	result, err := h.ticketSvc.ListAdmins(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *TicketHandler) AdminCreate(c *gin.Context) {
	var req model.AdminCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	result, err := h.ticketSvc.CreateAdmin(c.Request.Context(), &req)
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

func (h *TicketHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	var req model.AdminUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.ticketSvc.UpdateAdmin(c.Request.Context(), id, &req); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

func (h *TicketHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.ticketSvc.DeleteAdmin(c.Request.Context(), id); svcErr != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

// === Finance ===

func (h *TicketHandler) FinanceSummary(c *gin.Context) {
	result, err := h.ticketSvc.FinanceSummary(c.Request.Context())
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *TicketHandler) FinanceWithdrawals(c *gin.Context) {
	// MOCK: 待接入真实提现系统
	response.Success(c, &model.PaginatedList{Items: []interface{}{}, Total: 0, Page: 1, PageSize: 20, TotalPages: 0})
}

// === Logs ===

func (h *TicketHandler) OperationLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	result, err := h.ticketSvc.OperationLogs(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

func (h *TicketHandler) PaymentLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	result, err := h.ticketSvc.PaymentLogs(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
}

// === Comment admin reply ===

func (h *TicketHandler) AdminCommentReply(c *gin.Context) {
	adminID := c.GetInt64("user_id")
	commentID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	var req model.CommentReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if svcErr := h.ticketSvc.AdminReplyComment(c.Request.Context(), adminID, commentID, req.Content); svcErr != nil {
		if e, ok := svcErr.(*errcode.Error); ok {
			response.Error(c, e)
			return
		}
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, nil)
}

// === Payment notify (MOCK) ===

// MOCK: 待接入真实服务
func (h *TicketHandler) WechatNotify(c *gin.Context) {
	// MOCK: 待接入真实服务 — 直接标记订单已支付
	var body struct {
		OrderID int64 `json:"order_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Success(c, gin.H{"status": "ok", "mock": true})
		return
	}
	response.Success(c, gin.H{"status": "ok", "mock": true})
}

// MOCK: 待接入真实服务
func (h *TicketHandler) AlipayNotify(c *gin.Context) {
	// MOCK: 待接入真实服务 — 直接标记订单已支付
	response.Success(c, gin.H{"status": "ok", "mock": true})
}

// === OAuth (MOCK) ===

// MOCK: 待接入真实服务
func (h *TicketHandler) OAuthWechat(c *gin.Context) {
	// MOCK: 待接入真实服务 — 返回模拟授权URL
	response.Success(c, gin.H{
		"auth_url": "https://open.weixin.qq.com/connect/oauth2/authorize?appid=MOCK_APPID&redirect_uri=MOCK_REDIRECT&response_type=code&scope=snsapi_userinfo",
		"mock":     true,
	})
}

// MOCK: 待接入真实服务
func (h *TicketHandler) OAuthWechatCallback(c *gin.Context) {
	// MOCK: 待接入真实服务 — 自动创建用户并返回token
	response.Success(c, gin.H{
		"token": "mock_wechat_token",
		"user":  gin.H{"id": 0, "username": "wx_mock_user", "phone": "", "avatar": "", "is_vip": false},
		"mock":  true,
	})
}

// === Batch upload ===

func (h *TicketHandler) BatchImageUpload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	files := form.File["files"]
	if len(files) == 0 {
		response.Error(c, errcode.ErrParam)
		return
	}
	os.MkdirAll(h.uploadDir, 0755)
	urls := make([]string, 0, len(files))
	for _, file := range files {
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), file.Filename, ext)
		dst := filepath.Join(h.uploadDir, filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			continue
		}
		urls = append(urls, "/uploads/"+filename)
	}
	response.Success(c, gin.H{"urls": urls})
}

// === Home config ===

func (h *TicketHandler) HomeConfig(c *gin.Context) {
	settings, _ := h.ticketSvc.GetSettings(c.Request.Context())
	if settings == nil {
		settings = model.SettingsMap{}
	}
	response.Success(c, settings)
}
