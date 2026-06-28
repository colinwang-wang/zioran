package api

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zioran/backend/internal/middleware"
	"github.com/zioran/backend/internal/model"
	"github.com/zioran/backend/internal/service"
	"github.com/zioran/backend/pkg/errcode"
	"github.com/zioran/backend/pkg/oauth"
	"github.com/zioran/backend/pkg/payment"
	"github.com/zioran/backend/pkg/response"
)

type TicketHandler struct {
	ticketSvc   *service.TicketService
	authSvc     *service.AuthService
	paySvc      *service.PaymentService
	wechatOAuth *oauth.WechatOAuth
	jwtSecret   string
	jwtExpire   time.Duration
	uploadDir   string
}

func NewTicketHandler(ticketSvc *service.TicketService, authSvc *service.AuthService, paySvc *service.PaymentService, wechatOAuth *oauth.WechatOAuth, jwtSecret string, jwtExpire time.Duration, uploadDir string) *TicketHandler {
	return &TicketHandler{ticketSvc: ticketSvc, authSvc: authSvc, paySvc: paySvc, wechatOAuth: wechatOAuth, jwtSecret: jwtSecret, jwtExpire: jwtExpire, uploadDir: uploadDir}
}

// === Auth extensions ===

// ForgotPassword resets user password via email code.
func (h *TicketHandler) ForgotPassword(c *gin.Context) {
	var req model.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.ErrParam)
		return
	}
	if req.Email == "" || req.EmailCode == "" {
		response.Error(c, errcode.ErrParam)
		return
	}
	if !h.authSvc.VerifyEmailCode(req.Email, req.EmailCode) {
		response.Error(c, errcode.New(40001, "邮箱验证码错误"))
		return
	}
	err := h.ticketSvc.ForgotPasswordByEmail(c.Request.Context(), req.Email, req.NewPassword)
	if err != nil {
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
	status := c.Query("status")
	result, err := h.ticketSvc.AdminList(c.Request.Context(), page, pageSize, status)
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
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	status := c.Query("status")
	result, err := h.ticketSvc.FinanceWithdrawals(c.Request.Context(), page, pageSize, status)
	if err != nil {
		response.Error(c, errcode.ErrInternal)
		return
	}
	response.Success(c, result)
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

// === Payment notify ===

func (h *TicketHandler) WechatNotify(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(200, gin.H{"code": "FAIL", "message": "read body failed"})
		return
	}
	headers := payment.WechatNotifyHeaders{
		Signature: c.GetHeader("Wechatpay-Signature"),
		Timestamp: c.GetHeader("Wechatpay-Timestamp"),
		Nonce:     c.GetHeader("Wechatpay-Nonce"),
		Serial:    c.GetHeader("Wechatpay-Serial"),
	}
	if svcErr := h.paySvc.WechatNotifyCallback(c.Request.Context(), body, headers); svcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "FAIL", "message": svcErr.Error()})
		return
	}
	c.JSON(200, gin.H{"code": "SUCCESS", "message": "ok"})
}

func (h *TicketHandler) AlipayNotify(c *gin.Context) {
	c.Request.ParseForm()
	params := make(map[string]string)
	for k, v := range c.Request.PostForm {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	if svcErr := h.paySvc.AlipayNotifyCallback(c.Request.Context(), params); svcErr != nil {
		c.String(http.StatusInternalServerError, "fail")
		return
	}
	c.String(200, "success")
}

// === OAuth ===

func (h *TicketHandler) OAuthWechat(c *gin.Context) {
	state := c.DefaultQuery("state", "login")
	authURL := h.wechatOAuth.GetAuthURL(state)
	response.Success(c, gin.H{"auth_url": authURL})
}

func (h *TicketHandler) OAuthWechatCallback(c *gin.Context) {
	code := c.Query("code")
	if c.Request.Method == http.MethodGet {
		if code == "" {
			response.Error(c, errcode.ErrParam)
			return
		}
		if callbackURL := h.wechatOAuth.GetFrontendCallbackURL(code, c.Query("state")); callbackURL != "" {
			c.Redirect(http.StatusFound, callbackURL)
			return
		}
	} else {
		var req struct {
			Code string `json:"code" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, errcode.ErrParam)
			return
		}
		code = req.Code
	}
	result, err := h.completeWechatLogin(c, code)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, result)
}

func (h *TicketHandler) completeWechatLogin(c *gin.Context, code string) (gin.H, *errcode.Error) {
	wxUser, err := h.wechatOAuth.GetUserInfo(code)
	if err != nil {
		return nil, errcode.New(40001, "微信授权失败")
	}
	// Find or create user by wechat openid
	user, svcErr := h.ticketSvc.FindOrCreateByWechat(c.Request.Context(), wxUser.OpenID, wxUser.Nickname, wxUser.Avatar)
	if svcErr != nil {
		return nil, errcode.ErrInternal
	}
	token, tokenErr := middleware.GenerateToken(user.ID, h.jwtSecret, h.jwtExpire)
	if tokenErr != nil {
		return nil, errcode.ErrInternal
	}
	return gin.H{
		"token": token,
		"user":  gin.H{"id": user.ID, "username": user.Username, "email": user.Email, "avatar": user.AvatarURL, "is_vip": user.Role == "vip"},
	}, nil
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
	urls := make([]string, 0, len(files))
	for _, file := range files {
		url, saveErr := saveUploadedImage(c, h.uploadDir, file)
		if saveErr != nil {
			if e, ok := saveErr.(*errcode.Error); ok {
				response.Error(c, e)
				return
			}
			response.Error(c, errcode.ErrInternal)
			return
		}
		urls = append(urls, url)
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
