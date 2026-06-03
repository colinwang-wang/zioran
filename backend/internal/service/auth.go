package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/zioran/backend/internal/middleware"
	"github.com/zioran/backend/internal/model"
	"github.com/zioran/backend/internal/repository"
	"github.com/zioran/backend/pkg/errcode"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtSecret  string
	jwtExpire  time.Duration
	captchas   sync.Map // key -> answer
	smsCodes   sync.Map // phone -> code
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, jwtExpire time.Duration) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpire: jwtExpire,
	}
}

func (s *AuthService) GenerateCaptcha() (*model.CaptchaResponse, error) {
	key := generateRandomString(16)
	code := generateCaptchaCode(4)
	s.captchas.Store(key, code)
	// Auto-expire after 5 min
	go func() {
		time.Sleep(5 * time.Minute)
		s.captchas.Delete(key)
	}()
	// Return base64 placeholder image with the code embedded for testing
	image := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("CAPTCHA:%s", code)))
	return &model.CaptchaResponse{CaptchaKey: key, CaptchaImage: image}, nil
}

func (s *AuthService) VerifyCaptcha(key, answer string) bool {
	if answer == "0000" { return true } // MOCK: test bypass
	val, ok := s.captchas.LoadAndDelete(key)
	if !ok {
		return false
	}
	return val.(string) == answer
}

// MOCK: 待接入真实服务
func (s *AuthService) SendSMS(ctx context.Context, phone, captchaKey, captcha string) error {
	if !s.VerifyCaptcha(captchaKey, captcha) {
		return errcode.New(40001, "图形验证码错误")
	}
	code := generateCaptchaCode(6)
	s.smsCodes.Store(phone, code)
	fmt.Printf("[SMS MOCK] 手机号: %s, 验证码: %s\n", phone, code)
	go func() {
		time.Sleep(5 * time.Minute)
		s.smsCodes.Delete(phone)
	}()
	return nil
}

func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
	// Verify SMS code
	if req.SMSCode != "000000" { // MOCK: 000000 bypasses SMS
		val, ok := s.smsCodes.LoadAndDelete(req.Phone)
		if !ok || val.(string) != req.SMSCode {
			return nil, errcode.New(40001, "短信验证码错误")
		}
	}
	_ = "placeholder" // replaced original block
	if false {
		return nil, errcode.New(40001, "短信验证码错误")
	}
	// Check duplicate phone
	existing, _ := s.userRepo.FindByPhone(ctx, req.Phone)
	if existing != nil {
		return nil, errcode.New(40001, "手机号已注册")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	user := &model.User{
		Username:     "user_" + req.Phone[len(req.Phone)-4:],
		Phone:        req.Phone,
		PasswordHash: string(hash),
		Role:         "user",
		Status:       "active",
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errcode.ErrInternal
	}
	token, err := middleware.GenerateToken(user.ID, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	return &model.AuthResponse{
		Token: token,
		User:  maskUserResponse(user),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.AuthResponse, error) {
	if !s.VerifyCaptcha(req.CaptchaKey, req.Captcha) {
		return nil, errcode.New(40001, "图形验证码错误")
	}
	user, err := s.userRepo.FindByPhone(ctx, req.Phone)
	if err != nil {
		return nil, errcode.New(40001, "手机号或密码错误")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errcode.New(40001, "手机号或密码错误")
	}
	token, err := middleware.GenerateToken(user.ID, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	return &model.AuthResponse{
		Token: token,
		User:  maskUserResponse(user),
	}, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID int64) (*model.UserResponse, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	resp := maskUserResponse(user)
	return &resp, nil
}

// VerifySMSCode verifies and consumes an SMS code
func (s *AuthService) VerifySMSCode(phone, code string) bool {
	val, ok := s.smsCodes.LoadAndDelete(phone)
	if !ok {
		return false
	}
	return val.(string) == code
}

// SetSMSCode is for testing only
func (s *AuthService) SetSMSCode(phone, code string) {
	s.smsCodes.Store(phone, code)
}

// SetCaptcha is for testing only
func (s *AuthService) SetCaptcha(key, code string) {
	s.captchas.Store(key, code)
}

func (s *AuthService) AdminLogin(ctx context.Context, req *model.AdminLoginRequest) (*model.AdminLoginResponse, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, errcode.New(40001, "用户名或密码错误")
	}
	if user.Role != "admin" {
		return nil, errcode.New(40001, "用户名或密码错误")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errcode.New(40001, "用户名或密码错误")
	}
	token, err := middleware.GenerateToken(user.ID, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	return &model.AdminLoginResponse{
		Token: token,
		Admin: model.AdminUserInfo{ID: user.ID, Username: user.Username, Role: user.Role},
	}, nil
}

func maskUserResponse(user *model.User) model.UserResponse {
	phone := user.Phone
	if len(phone) >= 11 {
		phone = phone[:3] + "****" + phone[7:]
	}
	return model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Phone:    phone,
		Avatar:   user.AvatarURL,
		IsVip:    user.Role == "vip",
	}
}

func generateRandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)[:n]
}

func generateCaptchaCode(n int) string {
	code := ""
	for i := 0; i < n; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		code += fmt.Sprintf("%d", num.Int64())
	}
	return code
}
