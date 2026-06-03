package model

import "time"

type User struct {
	ID           int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	Username     string     `json:"username" gorm:"size:50;not null;uniqueIndex"`
	Phone        string     `json:"phone" gorm:"size:20;not null;uniqueIndex"`
	Email        string     `json:"email" gorm:"size:100"`
	PasswordHash string     `json:"-" gorm:"size:255;not null"`
	AvatarURL    string     `json:"avatar_url" gorm:"size:500"`
	Bio          string     `json:"bio" gorm:"type:text"`
	Role         string     `json:"role" gorm:"size:20;not null;default:user"`
	Status       string     `json:"status" gorm:"size:20;not null;default:active"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (User) TableName() string { return "users" }
