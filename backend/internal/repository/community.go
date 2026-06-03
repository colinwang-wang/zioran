package repository

import (
	"context"

	"github.com/zioran/backend/internal/model"
	"gorm.io/gorm"
)

type CommunityRepository struct {
	db *gorm.DB
}

func NewCommunityRepository(db *gorm.DB) *CommunityRepository {
	return &CommunityRepository{db: db}
}

// Guestbook

func (r *CommunityRepository) GuestbookList(ctx context.Context, page, pageSize int) ([]model.Guestbook, int64, error) {
	var items []model.Guestbook
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Guestbook{}).Where("status = ?", "visible")
	query.Count(&total)
	err := query.Preload("User").Order("is_pinned DESC, created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *CommunityRepository) GuestbookCreate(ctx context.Context, g *model.Guestbook) error {
	return r.db.WithContext(ctx).Create(g).Error
}

func (r *CommunityRepository) GuestbookGetByID(ctx context.Context, id int64) (*model.Guestbook, error) {
	var g model.Guestbook
	err := r.db.WithContext(ctx).First(&g, id).Error
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *CommunityRepository) GuestbookDelete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Guestbook{}, id).Error
}

func (r *CommunityRepository) GuestbookLike(ctx context.Context, userID, guestbookID int64) (bool, error) {
	var count int64
	r.db.WithContext(ctx).Model(&model.GuestbookLike{}).
		Where("user_id = ? AND guestbook_id = ?", userID, guestbookID).Count(&count)

	if count > 0 {
		r.db.WithContext(ctx).Where("user_id = ? AND guestbook_id = ?", userID, guestbookID).Delete(&model.GuestbookLike{})
		r.db.WithContext(ctx).Model(&model.Guestbook{}).Where("id = ?", guestbookID).
			Update("like_count", gorm.Expr("like_count - 1"))
		return false, nil
	}
	r.db.WithContext(ctx).Create(&model.GuestbookLike{UserID: userID, GuestbookID: guestbookID})
	r.db.WithContext(ctx).Model(&model.Guestbook{}).Where("id = ?", guestbookID).
		Update("like_count", gorm.Expr("like_count + 1"))
	return true, nil
}

func (r *CommunityRepository) IsGuestbookLiked(ctx context.Context, userID, guestbookID int64) bool {
	var count int64
	r.db.WithContext(ctx).Model(&model.GuestbookLike{}).
		Where("user_id = ? AND guestbook_id = ?", userID, guestbookID).Count(&count)
	return count > 0
}

// Admin guestbook
func (r *CommunityRepository) AdminGuestbookList(ctx context.Context, page, pageSize int) ([]model.Guestbook, int64, error) {
	var items []model.Guestbook
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Guestbook{})
	query.Count(&total)
	err := query.Preload("User").Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *CommunityRepository) GuestbookUpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&model.Guestbook{}).Where("id = ?", id).Update("status", status).Error
}

func (r *CommunityRepository) GuestbookPin(ctx context.Context, id int64, pinned bool) error {
	return r.db.WithContext(ctx).Model(&model.Guestbook{}).Where("id = ?", id).Update("is_pinned", pinned).Error
}

// Comments

func (r *CommunityRepository) CommentList(ctx context.Context, targetType string, targetID int64, page, pageSize int) ([]model.Comment, int64, error) {
	var items []model.Comment
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("target_type = ? AND target_id = ? AND status = ? AND parent_id IS NULL", targetType, targetID, "visible")
	query.Count(&total)
	err := query.Preload("User").Preload("Children", "status = ?", "visible").Preload("Children.User").
		Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *CommunityRepository) CommentCreate(ctx context.Context, c *model.Comment) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *CommunityRepository) CommentGetByID(ctx context.Context, id int64) (*model.Comment, error) {
	var c model.Comment
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CommunityRepository) CommentDelete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Comment{}, id).Error
}

// Admin comments
func (r *CommunityRepository) AdminCommentList(ctx context.Context, page, pageSize int) ([]model.Comment, int64, error) {
	var items []model.Comment
	var total int64
	query := r.db.WithContext(ctx).Model(&model.Comment{})
	query.Count(&total)
	err := query.Preload("User").Order("created_at DESC").
		Offset((page - 1) * pageSize).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *CommunityRepository) CommentUpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&model.Comment{}).Where("id = ?", id).Update("status", status).Error
}

// NavItems

func (r *CommunityRepository) NavItemList(ctx context.Context) ([]model.NavItem, error) {
	var items []model.NavItem
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("sort_order ASC").Find(&items).Error
	return items, err
}

func (r *CommunityRepository) AdminNavItemList(ctx context.Context) ([]model.NavItem, error) {
	var items []model.NavItem
	err := r.db.WithContext(ctx).Order("sort_order ASC").Find(&items).Error
	return items, err
}

func (r *CommunityRepository) NavItemCreate(ctx context.Context, item *model.NavItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *CommunityRepository) NavItemUpdate(ctx context.Context, item *model.NavItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *CommunityRepository) NavItemGetByID(ctx context.Context, id int) (*model.NavItem, error) {
	var item model.NavItem
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CommunityRepository) NavItemDelete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.NavItem{}, id).Error
}

// Banners

func (r *CommunityRepository) BannerList(ctx context.Context) ([]model.Banner, error) {
	var items []model.Banner
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("sort_order ASC").Find(&items).Error
	return items, err
}

func (r *CommunityRepository) AdminBannerList(ctx context.Context) ([]model.Banner, error) {
	var items []model.Banner
	err := r.db.WithContext(ctx).Order("sort_order ASC").Find(&items).Error
	return items, err
}

func (r *CommunityRepository) BannerCreate(ctx context.Context, item *model.Banner) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *CommunityRepository) BannerUpdate(ctx context.Context, item *model.Banner) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *CommunityRepository) BannerGetByID(ctx context.Context, id int) (*model.Banner, error) {
	var item model.Banner
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *CommunityRepository) BannerDelete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.Banner{}, id).Error
}
