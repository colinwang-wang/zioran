package model

import "time"

type Course struct {
	ID             int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Title          string     `json:"title" gorm:"size:500;not null"`
	Subtitle       string     `json:"subtitle" gorm:"size:500"`
	Slug           string     `json:"slug" gorm:"size:200;not null;uniqueIndex"`
	CategoryID     int        `json:"category_id" gorm:"not null"`
	QualityLabel   string     `json:"quality_label" gorm:"size:50"`
	CoverImage     string     `json:"cover_image" gorm:"size:500"`
	Content        string     `json:"content" gorm:"type:text"`
	DetailImages   string     `json:"detail_images" gorm:"type:text"`
	DetailTitle    string     `json:"detail_title" gorm:"size:500"`
	DetailSubtitle string     `json:"detail_subtitle" gorm:"size:500"`
	Price          int        `json:"price" gorm:"not null;default:0"`
	VipPrice       int        `json:"vip_price" gorm:"not null;default:0"`
	Status         string     `json:"status" gorm:"size:20;not null;default:draft"`
	ViewCount      int        `json:"view_count" gorm:"not null;default:0"`
	LikeCount      int        `json:"like_count" gorm:"not null;default:0"`
	DownloadCount  int        `json:"download_count" gorm:"not null;default:0"`
	CommentCount   int        `json:"comment_count" gorm:"not null;default:0"`
	AuthorID       *int64     `json:"author_id"`
	PublishedAt    *time.Time `json:"published_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	Category  *Category        `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Tags      []Tag            `json:"tags,omitempty" gorm:"many2many:course_tags;"`
	Resources []CourseResource `json:"resources,omitempty" gorm:"foreignKey:CourseID"`
}

func (Course) TableName() string { return "courses" }

type Category struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:50;not null"`
	Slug        string    `json:"slug" gorm:"size:50;not null;uniqueIndex"`
	Description string    `json:"description,omitempty" gorm:"type:text"`
	ParentID    *int      `json:"parent_id"`
	SortOrder   int       `json:"sort_order" gorm:"not null;default:0"`
	CourseCount int       `json:"course_count" gorm:"not null;default:0"`
	IsActive    bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedAt   time.Time `json:"created_at"`

	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

func (Category) TableName() string { return "categories" }

type Tag struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:50;not null"`
	Slug        string    `json:"slug" gorm:"size:50;not null;uniqueIndex"`
	CourseCount int       `json:"course_count" gorm:"not null;default:0"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Tag) TableName() string { return "tags" }

type CourseResource struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	CourseID  int64     `json:"course_id" gorm:"not null"`
	Name      string    `json:"name" gorm:"size:200"`
	URL       string    `json:"url" gorm:"size:1000;not null"`
	Password  string    `json:"password" gorm:"size:100"`
	SortOrder int       `json:"sort_order" gorm:"not null;default:0"`
	CreatedAt time.Time `json:"created_at"`
}

func (CourseResource) TableName() string { return "course_resources" }

type UserFavorite struct {
	UserID    int64     `json:"user_id" gorm:"primaryKey"`
	CourseID  int64     `json:"course_id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserFavorite) TableName() string { return "user_favorites" }
