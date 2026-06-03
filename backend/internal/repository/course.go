package repository

import (
	"context"

	"github.com/zioran/backend/internal/model"
	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) List(ctx context.Context, req *model.CourseListRequest) ([]model.Course, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Course{}).Where("status = ?", "published")

	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}
	if req.TagID > 0 {
		query = query.Where("id IN (SELECT course_id FROM course_tags WHERE tag_id = ?)", req.TagID)
	}
	if req.Keyword != "" {
		query = query.Where("title LIKE ? OR subtitle LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	var total int64
	query.Count(&total)

	switch req.Sort {
	case "popular":
		query = query.Order("view_count DESC")
	case "downloads":
		query = query.Order("download_count DESC")
	default:
		query = query.Order("published_at DESC")
	}

	var courses []model.Course
	offset := (req.Page - 1) * req.PageSize
	err := query.Preload("Category").Offset(offset).Limit(req.PageSize).Find(&courses).Error
	return courses, total, err
}

func (r *CourseRepository) FindBySlug(ctx context.Context, slug string) (*model.Course, error) {
	var course model.Course
	err := r.db.WithContext(ctx).Preload("Category").Preload("Tags").Preload("Resources", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort_order ASC")
	}).Where("slug = ?", slug).First(&course).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *CourseRepository) FindByID(ctx context.Context, id int64) (*model.Course, error) {
	var course model.Course
	err := r.db.WithContext(ctx).Preload("Category").Preload("Tags").Preload("Resources").First(&course, id).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *CourseRepository) Latest(ctx context.Context, limit int) ([]model.Course, error) {
	var courses []model.Course
	err := r.db.WithContext(ctx).Preload("Category").
		Where("status = ?", "published").
		Order("published_at DESC").Limit(limit).Find(&courses).Error
	return courses, err
}

func (r *CourseRepository) GetPrevNext(ctx context.Context, courseID int64, publishedAt interface{}) (*model.Course, *model.Course, error) {
	var prev, next model.Course
	var prevPtr, nextPtr *model.Course

	err := r.db.WithContext(ctx).
		Where("status = ? AND published_at < ? AND id != ?", "published", publishedAt, courseID).
		Order("published_at DESC").First(&prev).Error
	if err == nil {
		prevPtr = &prev
	}

	err = r.db.WithContext(ctx).
		Where("status = ? AND published_at > ? AND id != ?", "published", publishedAt, courseID).
		Order("published_at ASC").First(&next).Error
	if err == nil {
		nextPtr = &next
	}

	return prevPtr, nextPtr, nil
}

func (r *CourseRepository) Related(ctx context.Context, categoryID int, excludeID int64, limit int) ([]model.Course, error) {
	var courses []model.Course
	err := r.db.WithContext(ctx).Preload("Category").
		Where("status = ? AND category_id = ? AND id != ?", "published", categoryID, excludeID).
		Order("published_at DESC").Limit(limit).Find(&courses).Error
	return courses, err
}

func (r *CourseRepository) Search(ctx context.Context, req *model.SearchRequest) ([]model.Course, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Course{}).Where("status = ?", "published")
	query = query.Where("title LIKE ? OR subtitle LIKE ?", "%"+req.Q+"%", "%"+req.Q+"%")
	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}

	var total int64
	query.Count(&total)

	var courses []model.Course
	offset := (req.Page - 1) * req.PageSize
	err := query.Preload("Category").Order("published_at DESC").Offset(offset).Limit(req.PageSize).Find(&courses).Error
	return courses, total, err
}

func (r *CourseRepository) Create(ctx context.Context, course *model.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

func (r *CourseRepository) Update(ctx context.Context, course *model.Course) error {
	return r.db.WithContext(ctx).Save(course).Error
}

func (r *CourseRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&model.Course{}, id).Error
}

func (r *CourseRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&model.Course{}).Where("id = ?", id).Update("status", status).Error
}

func (r *CourseRepository) AdminList(ctx context.Context, page, pageSize int, categoryID int, status, keyword string) ([]model.Course, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Course{})
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var courses []model.Course
	offset := (page - 1) * pageSize
	err := query.Preload("Category").Preload("Tags").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&courses).Error
	return courses, total, err
}

func (r *CourseRepository) ReplaceTags(ctx context.Context, courseID int64, tagIDs []int) error {
	// Remove existing
	r.db.WithContext(ctx).Exec("DELETE FROM course_tags WHERE course_id = ?", courseID)
	// Insert new
	for _, tagID := range tagIDs {
		r.db.WithContext(ctx).Exec("INSERT INTO course_tags (course_id, tag_id) VALUES (?, ?)", courseID, tagID)
	}
	return nil
}

func (r *CourseRepository) ReplaceResources(ctx context.Context, courseID int64, resources []model.CourseResource) error {
	r.db.WithContext(ctx).Where("course_id = ?", courseID).Delete(&model.CourseResource{})
	for i := range resources {
		resources[i].CourseID = courseID
		r.db.WithContext(ctx).Create(&resources[i])
	}
	return nil
}

func (r *CourseRepository) IncrementLikeCount(ctx context.Context, id int64, delta int) error {
	return r.db.WithContext(ctx).Model(&model.Course{}).Where("id = ?", id).
		Update("like_count", gorm.Expr("like_count + ?", delta)).Error
}
