package repository

import (
	"context"
	"time"

	"github.com/zioran/backend/internal/model"
	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) IsFavorited(ctx context.Context, userID, courseID int64) bool {
	var count int64
	r.db.WithContext(ctx).Model(&model.UserFavorite{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).Count(&count)
	return count > 0
}

func (r *FavoriteRepository) Toggle(ctx context.Context, userID, courseID int64) (bool, error) {
	if r.IsFavorited(ctx, userID, courseID) {
		err := r.db.WithContext(ctx).
			Where("user_id = ? AND course_id = ?", userID, courseID).
			Delete(&model.UserFavorite{}).Error
		return false, err
	}
	err := r.db.WithContext(ctx).Create(&model.UserFavorite{
		UserID: userID, CourseID: courseID,
	}).Error
	return true, err
}

func (r *FavoriteRepository) HasPurchased(ctx context.Context, userID, courseID int64) bool {
	var count int64
	r.db.WithContext(ctx).Model(&model.Purchase{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).Count(&count)
	return count > 0
}

func (r *FavoriteRepository) IsVip(ctx context.Context, userID int64) bool {
	var count int64
	r.db.WithContext(ctx).Model(&model.UserVip{}).
		Where("user_id = ? AND is_active = ? AND (expires_at IS NULL OR expires_at > ?)", userID, true, time.Now()).
		Count(&count)
	return count > 0
}
