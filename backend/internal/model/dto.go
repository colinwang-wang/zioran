package model

import "time"

// Auth DTOs

type CaptchaResponse struct {
	CaptchaKey   string `json:"captcha_key"`
	CaptchaImage string `json:"captcha_image"`
}

type SendEmailRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Captcha    string `json:"captcha" binding:"required"`
	CaptchaKey string `json:"captcha_key" binding:"required"`
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"omitempty,max=50"`
	Email     string `json:"email" binding:"required,email"`
	EmailCode string `json:"email_code" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	Captcha    string `json:"captcha" binding:"required"`
	CaptchaKey string `json:"captcha_key" binding:"required"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	IsVip    bool   `json:"is_vip"`
}

type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email    string `json:"email" binding:"omitempty,email"`
}

// Admin login DTOs

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdminLoginResponse struct {
	Token string        `json:"token"`
	Admin AdminUserInfo `json:"admin"`
}

type AdminUserInfo struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
