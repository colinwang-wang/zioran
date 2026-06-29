package service

import (
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/zioran/backend/internal/middleware"
	"github.com/zioran/backend/internal/model"
	"github.com/zioran/backend/internal/repository"
	"github.com/zioran/backend/pkg/email"
	"github.com/zioran/backend/pkg/errcode"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    *repository.UserRepository
	jwtSecret   string
	jwtExpire   time.Duration
	emailSender email.Sender
	captchas    sync.Map // key -> answer
	emailCodes  sync.Map // email -> code
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, jwtExpire time.Duration) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		jwtExpire:   jwtExpire,
		emailSender: &email.MockSender{},
	}
}

func (s *AuthService) SetEmailSender(sender email.Sender) {
	if sender != nil {
		s.emailSender = sender
	}
}

func (s *AuthService) GenerateCaptcha() (*model.CaptchaResponse, error) {
	key := generateRandomString(16)
	code := generateCaptchaCode(4)
	s.captchas.Store(key, code)
	go func() {
		time.Sleep(5 * time.Minute)
		s.captchas.Delete(key)
	}()
	svg := generateCaptchaSVG(code)
	image := "data:image/svg+xml;base64," + base64.StdEncoding.EncodeToString([]byte(svg))
	return &model.CaptchaResponse{CaptchaKey: key, CaptchaImage: image}, nil
}

func generateCaptchaSVG(code string) string {
	colors := []string{"#e60023", "#333", "#0066cc", "#009933"}
	svg := `<svg xmlns="http://www.w3.org/2000/svg" width="120" height="40" viewBox="0 0 120 40"><rect width="120" height="40" fill="#f5f5f5" rx="4"/>`
	for i, c := range code {
		x := 15 + i*25
		svg += fmt.Sprintf(`<text x="%d" y="28" font-size="22" font-family="Arial" fill="%s" transform="rotate(%d %d 20)">%c</text>`, x, colors[i%4], i*7-10, x, c)
	}
	svg += `<line x1="10" y1="15" x2="110" y2="25" stroke="#ddd"/><line x1="20" y1="32" x2="100" y2="8" stroke="#eee"/></svg>`
	return svg
}

func (s *AuthService) VerifyCaptcha(key, answer string) bool {
	val, ok := s.captchas.LoadAndDelete(key)
	if !ok {
		return false
	}
	return val.(string) == answer
}

func (s *AuthService) SendEmail(ctx context.Context, emailAddress, captchaKey, captcha string) error {
	if !s.VerifyCaptcha(captchaKey, captcha) {
		return errcode.New(40001, "图形验证码错误")
	}
	emailAddress = normalizeEmail(emailAddress)
	code := generateCaptchaCode(6)
	s.emailCodes.Store(emailAddress, code)
	go func() {
		time.Sleep(5 * time.Minute)
		s.emailCodes.Delete(emailAddress)
	}()
	return s.emailSender.Send(emailAddress, code)
}

func (s *AuthService) Register(ctx context.Context, req *model.RegisterRequest) (*model.AuthResponse, error) {
	req.Email = normalizeEmail(req.Email)
	val, ok := s.emailCodes.LoadAndDelete(req.Email)
	if !ok || val.(string) != req.EmailCode {
		return nil, errcode.New(40001, "邮箱验证码错误")
	}
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errcode.New(40001, "邮箱已注册")
	}
	username := strings.TrimSpace(req.Username)
	if username == "" {
		username = usernameFromEmail(req.Email)
	} else {
		existing, _ := s.userRepo.FindByUsername(ctx, username)
		if existing != nil {
			return nil, errcode.New(40001, "用户名已存在")
		}
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	user := &model.User{
		Username:     username,
		Phone:        phonePlaceholderFromEmail(req.Email),
		Email:        req.Email,
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
	user, err := s.userRepo.FindByEmail(ctx, normalizeEmail(req.Email))
	if err != nil {
		return nil, errcode.New(40001, "邮箱或密码错误")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errcode.New(40001, "邮箱或密码错误")
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

func (s *AuthService) UpdateProfile(ctx context.Context, userID int64, req *model.UpdateProfileRequest) (*model.UserResponse, error) {
	updates := map[string]interface{}{}
	username := strings.TrimSpace(req.Username)
	if username != "" {
		updates["username"] = username
	}
	emailAddress := normalizeEmail(req.Email)
	if emailAddress != "" {
		existing, _ := s.userRepo.FindByEmail(ctx, emailAddress)
		if existing != nil && existing.ID != userID {
			return nil, errcode.New(40001, "邮箱已绑定")
		}
		updates["email"] = emailAddress
	}
	if len(updates) > 0 {
		if err := s.userRepo.UpdateProfile(ctx, userID, updates); err != nil {
			return nil, errcode.ErrInternal
		}
	}
	return s.GetProfile(ctx, userID)
}

func (s *AuthService) VerifyEmailCode(emailAddress, code string) bool {
	val, ok := s.emailCodes.LoadAndDelete(normalizeEmail(emailAddress))
	if !ok {
		return false
	}
	return val.(string) == code
}

func (s *AuthService) SetEmailCode(emailAddress, code string) {
	s.emailCodes.Store(normalizeEmail(emailAddress), code)
}

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
	return model.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.AvatarURL,
		IsVip:    user.Role == "vip",
	}
}

func normalizeEmail(emailAddress string) string {
	return strings.ToLower(strings.TrimSpace(emailAddress))
}

func usernameFromEmail(emailAddress string) string {
	local := strings.SplitN(emailAddress, "@", 2)[0]
	var b strings.Builder
	for _, r := range local {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	name := strings.Trim(b.String(), "_")
	if name == "" {
		name = "user"
	}
	if len(name) > 24 {
		name = name[:24]
	}
	return name + "_" + emailHash(emailAddress)[:8]
}

func phonePlaceholderFromEmail(emailAddress string) string {
	return "email_" + emailHash(emailAddress)[:14]
}

func emailHash(emailAddress string) string {
	sum := sha1.Sum([]byte(emailAddress))
	return hex.EncodeToString(sum[:])
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
