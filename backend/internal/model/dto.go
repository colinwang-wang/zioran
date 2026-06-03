package model

// Auth DTOs

type CaptchaResponse struct {
	CaptchaKey   string `json:"captcha_key"`
	CaptchaImage string `json:"captcha_image"`
}

type SendSMSRequest struct {
	Phone      string `json:"phone" binding:"required"`
	Captcha    string `json:"captcha" binding:"required"`
	CaptchaKey string `json:"captcha_key" binding:"required"`
}

type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	SMSCode  string `json:"sms_code" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Phone      string `json:"phone" binding:"required"`
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
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	IsVip    bool   `json:"is_vip"`
}
