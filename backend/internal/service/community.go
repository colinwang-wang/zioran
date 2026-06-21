package service

import (
	"context"

	"github.com/zioran/backend/internal/model"
	"github.com/zioran/backend/internal/repository"
	"github.com/zioran/backend/pkg/errcode"
)

type CommunityService struct {
	repo *repository.CommunityRepository
}

func NewCommunityService(repo *repository.CommunityRepository) *CommunityService {
	return &CommunityService{repo: repo}
}

// Guestbook

func (s *CommunityService) GuestbookList(ctx context.Context, userID int64, page, pageSize int) (*model.PaginatedList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	items, total, err := s.repo.GuestbookList(ctx, page, pageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	results := make([]model.GuestbookResponse, len(items))
	for i, g := range items {
		results[i] = model.GuestbookResponse{
			ID: g.ID, UserID: g.UserID, Content: g.Content,
			LikeCount: g.LikeCount, IsPinned: g.IsPinned, CreatedAt: g.CreatedAt,
		}
		if g.User != nil {
			results[i].Username = g.User.Username
			results[i].Avatar = g.User.AvatarURL
		}
		if userID > 0 {
			results[i].IsLiked = s.repo.IsGuestbookLiked(ctx, userID, g.ID)
		}
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: results, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}, nil
}

func (s *CommunityService) GuestbookCreate(ctx context.Context, userID int64, req *model.GuestbookCreateRequest) (*model.Guestbook, error) {
	g := &model.Guestbook{UserID: userID, Content: req.Content, Status: "visible"}
	if err := s.repo.GuestbookCreate(ctx, g); err != nil {
		return nil, errcode.ErrInternal
	}
	return g, nil
}

func (s *CommunityService) GuestbookLike(ctx context.Context, userID, id int64) (bool, error) {
	_, err := s.repo.GuestbookGetByID(ctx, id)
	if err != nil {
		return false, errcode.ErrNotFound
	}
	return s.repo.GuestbookLike(ctx, userID, id)
}

func (s *CommunityService) GuestbookDelete(ctx context.Context, userID, id int64, isAdmin bool) error {
	g, err := s.repo.GuestbookGetByID(ctx, id)
	if err != nil {
		return errcode.ErrNotFound
	}
	if !isAdmin && g.UserID != userID {
		return errcode.ErrForbidden
	}
	return s.repo.GuestbookDelete(ctx, id)
}

// Comments

func (s *CommunityService) CommentList(ctx context.Context, targetType string, targetID int64, page, pageSize int) (*model.PaginatedList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	items, total, err := s.repo.CommentList(ctx, targetType, targetID, page, pageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	results := make([]model.CommentResponse, len(items))
	for i, c := range items {
		results[i] = toCommentResponse(&c)
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: results, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}, nil
}

func (s *CommunityService) CommentCreate(ctx context.Context, userID int64, req *model.CommentCreateRequest) (*model.Comment, error) {
	c := &model.Comment{
		UserID:     userID,
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		ParentID:   req.ParentID,
		Content:    req.Content,
		Status:     "visible",
	}
	if err := s.repo.CommentCreate(ctx, c); err != nil {
		return nil, errcode.ErrInternal
	}
	return c, nil
}

func (s *CommunityService) CommentDelete(ctx context.Context, userID, id int64, isAdmin bool) error {
	c, err := s.repo.CommentGetByID(ctx, id)
	if err != nil {
		return errcode.ErrNotFound
	}
	if !isAdmin && c.UserID != userID {
		return errcode.ErrForbidden
	}
	return s.repo.CommentDelete(ctx, id)
}

func (s *CommunityService) UserCommentList(ctx context.Context, userID int64, page, pageSize int) (*model.PaginatedList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	items, total, err := s.repo.UserCommentList(ctx, userID, page, pageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	results := make([]model.CommentResponse, len(items))
	for i, c := range items {
		results[i] = toCommentResponse(&c)
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: results, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}, nil
}

// Home config

func (s *CommunityService) NavItems(ctx context.Context) ([]model.NavItem, error) {
	return s.repo.NavItemList(ctx)
}

func (s *CommunityService) Banners(ctx context.Context) ([]model.Banner, error) {
	return s.repo.BannerList(ctx)
}

// Admin guestbook

func (s *CommunityService) AdminGuestbookList(ctx context.Context, page, pageSize int) (*model.PaginatedList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	items, total, err := s.repo.AdminGuestbookList(ctx, page, pageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: items, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}, nil
}

func (s *CommunityService) AdminGuestbookUpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := s.repo.GuestbookGetByID(ctx, id)
	if err != nil {
		return errcode.ErrNotFound
	}
	return s.repo.GuestbookUpdateStatus(ctx, id, status)
}

func (s *CommunityService) AdminGuestbookPin(ctx context.Context, id int64, pinned bool) error {
	_, err := s.repo.GuestbookGetByID(ctx, id)
	if err != nil {
		return errcode.ErrNotFound
	}
	return s.repo.GuestbookPin(ctx, id, pinned)
}

// Admin comments

func (s *CommunityService) AdminCommentList(ctx context.Context, page, pageSize int) (*model.PaginatedList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	items, total, err := s.repo.AdminCommentList(ctx, page, pageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: items, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}, nil
}

func (s *CommunityService) AdminCommentUpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := s.repo.CommentGetByID(ctx, id)
	if err != nil {
		return errcode.ErrNotFound
	}
	return s.repo.CommentUpdateStatus(ctx, id, status)
}

// Admin nav items

func (s *CommunityService) AdminNavItemList(ctx context.Context) ([]model.NavItem, error) {
	return s.repo.AdminNavItemList(ctx)
}

func (s *CommunityService) AdminNavItemCreate(ctx context.Context, req *model.NavItemRequest) (*model.NavItem, error) {
	item := &model.NavItem{Title: req.Title, Icon: req.Icon, URL: req.URL, SortOrder: req.SortOrder, IsActive: true}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}
	if err := s.repo.NavItemCreate(ctx, item); err != nil {
		return nil, errcode.ErrInternal
	}
	return item, nil
}

func (s *CommunityService) AdminNavItemUpdate(ctx context.Context, id int, req *model.NavItemRequest) (*model.NavItem, error) {
	item, err := s.repo.NavItemGetByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	item.Title = req.Title
	item.Icon = req.Icon
	item.URL = req.URL
	item.SortOrder = req.SortOrder
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}
	if err := s.repo.NavItemUpdate(ctx, item); err != nil {
		return nil, errcode.ErrInternal
	}
	return item, nil
}

func (s *CommunityService) AdminNavItemDelete(ctx context.Context, id int) error {
	_, err := s.repo.NavItemGetByID(ctx, id)
	if err != nil {
		return errcode.ErrNotFound
	}
	return s.repo.NavItemDelete(ctx, id)
}

// Admin banners

func (s *CommunityService) AdminBannerList(ctx context.Context) ([]model.Banner, error) {
	return s.repo.AdminBannerList(ctx)
}

func (s *CommunityService) AdminBannerCreate(ctx context.Context, req *model.BannerRequest) (*model.Banner, error) {
	item := &model.Banner{Title: req.Title, ImageURL: req.ImageURL, LinkURL: req.LinkURL, SortOrder: req.SortOrder, IsActive: true}
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}
	if err := s.repo.BannerCreate(ctx, item); err != nil {
		return nil, errcode.ErrInternal
	}
	return item, nil
}

func (s *CommunityService) AdminBannerUpdate(ctx context.Context, id int, req *model.BannerRequest) (*model.Banner, error) {
	item, err := s.repo.BannerGetByID(ctx, id)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	item.Title = req.Title
	item.ImageURL = req.ImageURL
	item.LinkURL = req.LinkURL
	item.SortOrder = req.SortOrder
	if req.IsActive != nil {
		item.IsActive = *req.IsActive
	}
	if err := s.repo.BannerUpdate(ctx, item); err != nil {
		return nil, errcode.ErrInternal
	}
	return item, nil
}

func (s *CommunityService) AdminBannerDelete(ctx context.Context, id int) error {
	_, err := s.repo.BannerGetByID(ctx, id)
	if err != nil {
		return errcode.ErrNotFound
	}
	return s.repo.BannerDelete(ctx, id)
}

func toCommentResponse(c *model.Comment) model.CommentResponse {
	resp := model.CommentResponse{
		ID: c.ID, UserID: c.UserID, TargetType: c.TargetType, TargetID: c.TargetID,
		Content: c.Content, ParentID: c.ParentID, Status: c.Status, CreatedAt: c.CreatedAt,
	}
	if c.User != nil {
		resp.Username = c.User.Username
		resp.Avatar = c.User.AvatarURL
	}
	if len(c.Children) > 0 {
		resp.Children = make([]model.CommentResponse, len(c.Children))
		for i, child := range c.Children {
			resp.Children[i] = toCommentResponse(&child)
		}
	}
	return resp
}
