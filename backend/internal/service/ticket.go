package service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"strings"

	"github.com/zioran/backend/internal/model"
	"github.com/zioran/backend/internal/repository"
	"github.com/zioran/backend/pkg/errcode"
	"golang.org/x/crypto/bcrypt"
)

type TicketService struct {
	repo     *repository.TicketRepository
	userRepo *repository.UserRepository
}

func NewTicketService(repo *repository.TicketRepository, userRepo *repository.UserRepository) *TicketService {
	return &TicketService{repo: repo, userRepo: userRepo}
}

// User ticket operations

func (s *TicketService) Create(ctx context.Context, userID int64, req *model.CreateTicketRequest) (*model.TicketResponse, error) {
	ticket := &model.Ticket{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
		Status:  "open",
	}
	if err := s.repo.Create(ctx, ticket); err != nil {
		return nil, errcode.ErrInternal
	}
	attachments := make([]model.TicketAttachment, 0, len(req.Attachments))
	for _, url := range req.Attachments {
		url = strings.TrimSpace(url)
		if url == "" {
			continue
		}
		attachments = append(attachments, model.TicketAttachment{TicketID: ticket.ID, URL: url})
	}
	if err := s.repo.CreateAttachments(ctx, attachments); err != nil {
		return nil, errcode.ErrInternal
	}
	return &model.TicketResponse{
		ID: ticket.ID, UserID: userID, Title: ticket.Title,
		Content: ticket.Content, Status: ticket.Status, CreatedAt: ticket.CreatedAt,
		Attachments: req.Attachments,
	}, nil
}

func (s *TicketService) List(ctx context.Context, userID int64, page, pageSize int) (*model.PaginatedList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	items, total, err := s.repo.ListByUser(ctx, userID, page, pageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	results := make([]model.TicketResponse, len(items))
	for i, t := range items {
		results[i] = model.TicketResponse{
			ID: t.ID, UserID: t.UserID, Title: t.Title,
			Content: t.Content, Status: t.Status, CreatedAt: t.CreatedAt,
		}
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: results, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}, nil
}

func (s *TicketService) Detail(ctx context.Context, userID int64, ticketID int64, isAdmin bool) (*model.TicketResponse, error) {
	ticket, err := s.repo.FindByID(ctx, ticketID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	if !isAdmin && ticket.UserID != userID {
		return nil, errcode.ErrForbidden
	}
	resp := &model.TicketResponse{
		ID: ticket.ID, UserID: ticket.UserID, Title: ticket.Title,
		Content: ticket.Content, Status: ticket.Status, CreatedAt: ticket.CreatedAt,
	}
	if ticket.User != nil {
		resp.Username = ticket.User.Username
	}
	resp.Attachments = make([]string, 0, len(ticket.Attachments))
	for _, attachment := range ticket.Attachments {
		resp.Attachments = append(resp.Attachments, attachment.URL)
	}
	resp.Replies = make([]model.TicketReplyResponse, len(ticket.Replies))
	for i, r := range ticket.Replies {
		resp.Replies[i] = model.TicketReplyResponse{
			ID: r.ID, UserID: r.UserID, Content: r.Content,
			IsAdmin: r.IsAdmin, CreatedAt: r.CreatedAt,
		}
		if r.User != nil {
			resp.Replies[i].Username = r.User.Username
		}
	}
	return resp, nil
}

func (s *TicketService) Reply(ctx context.Context, userID int64, ticketID int64, content string, isAdmin bool) error {
	ticket, err := s.repo.FindByID(ctx, ticketID)
	if err != nil {
		return errcode.ErrNotFound
	}
	if !isAdmin && ticket.UserID != userID {
		return errcode.ErrForbidden
	}
	reply := &model.TicketReply{
		TicketID: ticketID, UserID: userID, Content: content, IsAdmin: isAdmin,
	}
	if err := s.repo.CreateReply(ctx, reply); err != nil {
		return errcode.ErrInternal
	}
	if isAdmin {
		s.repo.UpdateStatus(ctx, ticketID, "replied")
	}
	return nil
}

// Admin ticket operations

func (s *TicketService) AdminList(ctx context.Context, page, pageSize int, status string) (*model.PaginatedList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	items, total, err := s.repo.ListAll(ctx, page, pageSize, status)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	results := make([]model.TicketResponse, len(items))
	for i, t := range items {
		results[i] = model.TicketResponse{
			ID: t.ID, UserID: t.UserID, Title: t.Title,
			Content: t.Content, Status: t.Status, CreatedAt: t.CreatedAt,
		}
		if t.User != nil {
			results[i].Username = t.User.Username
		}
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: results, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}, nil
}

func (s *TicketService) UpdateStatus(ctx context.Context, ticketID int64, status string) error {
	valid := map[string]bool{"processing": true, "replied": true, "closed": true}
	if !valid[status] {
		return errcode.ErrParam
	}
	return s.repo.UpdateStatus(ctx, ticketID, status)
}

// Settings

func (s *TicketService) GetSettings(ctx context.Context) (model.SettingsMap, error) {
	settings, err := s.repo.GetAllSettings(ctx)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	result := make(model.SettingsMap)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result, nil
}

func (s *TicketService) UpdateSettings(ctx context.Context, settings model.SettingsMap) error {
	return s.repo.UpsertSettings(ctx, settings)
}

// Admin accounts

func (s *TicketService) ListAdmins(ctx context.Context) ([]model.AdminUserInfo, error) {
	users, err := s.repo.ListAdmins(ctx)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	result := make([]model.AdminUserInfo, len(users))
	for i, u := range users {
		result[i] = model.AdminUserInfo{ID: u.ID, Username: u.Username, Role: u.Role}
	}
	return result, nil
}

func (s *TicketService) CreateAdmin(ctx context.Context, req *model.AdminCreateRequest) (*model.AdminUserInfo, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	user := &model.User{
		Username:     req.Username,
		Phone:        "admin_" + req.Username,
		PasswordHash: string(hash),
		Role:         req.Role,
		Status:       "active",
	}
	if err := s.repo.CreateAdmin(ctx, user); err != nil {
		return nil, errcode.ErrInternal
	}
	return &model.AdminUserInfo{ID: user.ID, Username: user.Username, Role: user.Role}, nil
}

func (s *TicketService) UpdateAdmin(ctx context.Context, id int64, req *model.AdminUpdateRequest) error {
	updates := map[string]interface{}{}
	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Role != "" {
		updates["role"] = req.Role
	}
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return errcode.ErrInternal
		}
		updates["password_hash"] = string(hash)
	}
	if len(updates) == 0 {
		return nil
	}
	return s.repo.UpdateAdmin(ctx, id, updates)
}

func (s *TicketService) DeleteAdmin(ctx context.Context, id int64) error {
	return s.repo.DeleteAdmin(ctx, id)
}

// Finance

func (s *TicketService) FinanceSummary(ctx context.Context) (*model.FinanceSummary, error) {
	return s.repo.FinanceSummary(ctx)
}

func (s *TicketService) FinanceWithdrawals(ctx context.Context, page, pageSize int, status string) (*model.PaginatedList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	items, total, err := s.repo.FinanceWithdrawals(ctx, page, pageSize, status)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	results := make([]model.FinanceWithdrawalResponse, len(items))
	for i, item := range items {
		results[i] = model.FinanceWithdrawalResponse{
			ID: item.ID, UserID: item.UserID, Amount: item.Amount,
			AccountName: item.AccountName, AccountNo: item.AccountNo, BankName: item.BankName,
			Status: item.Status, Remark: item.Remark, ProcessedAt: item.ProcessedAt, CreatedAt: item.CreatedAt,
		}
		if item.User != nil {
			results[i].Username = item.User.Username
		}
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: results, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}, nil
}

// Logs

func (s *TicketService) OperationLogs(ctx context.Context, page, pageSize int) (*model.PaginatedList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	items, total, err := s.repo.OperationLogs(ctx, page, pageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: items, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}, nil
}

func (s *TicketService) PaymentLogs(ctx context.Context, page, pageSize int) (*model.PaginatedList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	items, total, err := s.repo.PaymentLogs(ctx, page, pageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: items, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}, nil
}

// Comment admin reply

func (s *TicketService) AdminReplyComment(ctx context.Context, adminID int64, commentID int64, content string) error {
	parent, err := s.repo.FindComment(ctx, commentID)
	if err != nil {
		return errcode.ErrNotFound
	}
	reply := &model.Comment{
		UserID:     adminID,
		TargetType: parent.TargetType,
		TargetID:   parent.TargetID,
		ParentID:   &commentID,
		Content:    content,
		Status:     "visible",
	}
	return s.repo.CreateComment(ctx, reply)
}

// Order cancel

func (s *TicketService) CancelOrder(ctx context.Context, userID, orderID int64) error {
	order, err := s.repo.GetOrder(ctx, orderID)
	if err != nil {
		return errcode.ErrNotFound
	}
	if order.UserID != userID {
		return errcode.ErrForbidden
	}
	if order.Status != "pending" {
		return errcode.New(40001, "仅待支付订单可取消")
	}
	return s.repo.UpdateOrderStatus(ctx, orderID, "cancelled")
}

// Forgot password

func (s *TicketService) ForgotPasswordByEmail(ctx context.Context, email, newPassword string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return errcode.ErrNotFound
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errcode.ErrInternal
	}
	return s.userRepo.UpdatePassword(ctx, user.ID, string(hash))
}

// FindOrCreateByWechat finds a user by wechat openid or creates one.
func (s *TicketService) FindOrCreateByWechat(ctx context.Context, openID, nickname, avatar string) (*model.User, error) {
	user, err := s.userRepo.FindByWechatOpenID(ctx, openID)
	if err == nil {
		return user, nil
	}
	// Create new user
	openIDHash := shortSHA1(openID, 17)
	user = &model.User{
		Username:     buildWechatUsername(nickname, openIDHash[:8]),
		Phone:        "wx_" + openIDHash,
		PasswordHash: "-",
		AvatarURL:    avatar,
		Role:         "user",
		Status:       "active",
		WechatOpenID: openID,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func buildWechatUsername(nickname, suffix string) string {
	if nickname == "" {
		return "wx_" + suffix
	}
	const maxUsernameRunes = 50
	const suffixRunes = 9
	nameRunes := []rune(nickname)
	limit := maxUsernameRunes - suffixRunes
	if len(nameRunes) > limit {
		nameRunes = nameRunes[:limit]
	}
	return string(nameRunes) + "_" + suffix
}

func shortSHA1(value string, size int) string {
	sum := sha1.Sum([]byte(value))
	encoded := hex.EncodeToString(sum[:])
	if size > len(encoded) {
		size = len(encoded)
	}
	return encoded[:size]
}

// Refresh token - validated in handler using middleware.GenerateToken directly
