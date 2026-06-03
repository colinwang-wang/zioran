package model

import "time"

// Course list/detail DTOs

type CourseListRequest struct {
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"pageSize,default=16"`
	CategoryID int    `form:"categoryId"`
	TagID      int    `form:"tagId"`
	Sort       string `form:"sort,default=latest"`
	Keyword    string `form:"keyword"`
}

type CourseListItem struct {
	ID           int64           `json:"id"`
	Title        string          `json:"title"`
	Subtitle     string          `json:"subtitle"`
	Slug         string          `json:"slug"`
	Cover        string          `json:"cover"`
	Category     *CategoryBrief  `json:"category"`
	Price        int             `json:"price"`
	VipPrice     int             `json:"vip_price"`
	RelativeTime string          `json:"relative_time"`
	PublishedAt  *time.Time      `json:"published_at"`
}

type CategoryBrief struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type PaginatedList struct {
	Items      interface{} `json:"items"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalPages int         `json:"totalPages"`
}

type CourseDetailResponse struct {
	ID             int64            `json:"id"`
	Title          string           `json:"title"`
	Subtitle       string           `json:"subtitle"`
	Slug           string           `json:"slug"`
	Cover          string           `json:"cover"`
	Content        string           `json:"content"`
	DetailTitle    string           `json:"detail_title"`
	DetailSubtitle string           `json:"detail_subtitle"`
	Price          int              `json:"price"`
	VipPrice       int              `json:"vip_price"`
	QualityLabel   string           `json:"quality_label"`
	Category       *CategoryBrief   `json:"category"`
	Tags           []TagBrief       `json:"tags"`
	LikeCount      int              `json:"like_count"`
	CommentCount   int              `json:"comment_count"`
	PublishedAt    *time.Time       `json:"published_at"`
	PrevCourse     *CourseBrief     `json:"prev_course"`
	NextCourse     *CourseBrief     `json:"next_course"`
	RelatedCourses []CourseListItem `json:"related_courses"`
	UserAccess     *UserAccess      `json:"user_access"`
	Resources      []ResourceItem   `json:"resources"`
}

type TagBrief struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type CourseBrief struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

type UserAccess struct {
	CanDownload  bool `json:"can_download"`
	HasPurchased bool `json:"has_purchased"`
	IsVip        bool `json:"is_vip"`
	IsFavorited  bool `json:"is_favorited"`
}

type ResourceItem struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	URL      string `json:"url"`
	Password string `json:"password"`
}

type LikeResponse struct {
	Liked     bool `json:"liked"`
	LikeCount int  `json:"like_count"`
}

// Search
type SearchRequest struct {
	Q          string `form:"q" binding:"required"`
	CategoryID int    `form:"categoryId"`
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"pageSize,default=16"`
}

// Admin course DTOs
type AdminCourseRequest struct {
	Title          string       `json:"title" binding:"required"`
	Subtitle       string       `json:"subtitle"`
	Slug           string       `json:"slug" binding:"required"`
	CategoryID     int          `json:"category_id" binding:"required"`
	QualityLabel   string       `json:"quality_label"`
	CoverImage     string       `json:"cover_image"`
	Content        string       `json:"content"`
	DetailTitle    string       `json:"detail_title"`
	DetailSubtitle string       `json:"detail_subtitle"`
	Price          int          `json:"price"`
	VipPrice       int          `json:"vip_price"`
	TagIDs         []int        `json:"tag_ids"`
	Resources      []ResourceIn `json:"resources"`
}

type ResourceIn struct {
	Name      string `json:"name"`
	URL       string `json:"url" binding:"required"`
	Password  string `json:"password"`
	SortOrder int    `json:"sort_order"`
}

type AdminCourseStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// Admin category DTOs
type AdminCategoryRequest struct {
	Name      string `json:"name" binding:"required"`
	Slug      string `json:"slug" binding:"required"`
	ParentID  *int   `json:"parent_id"`
	SortOrder int    `json:"sort_order"`
	IsActive  *bool  `json:"is_active"`
}

// Admin tag DTOs
type AdminTagRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}
