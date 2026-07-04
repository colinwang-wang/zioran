package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zioran/backend/internal/api"
	"github.com/zioran/backend/internal/model"
	"github.com/zioran/backend/internal/repository"
	"github.com/zioran/backend/internal/service"
	"github.com/zioran/backend/pkg/oauth"
	"github.com/zioran/backend/pkg/payment"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const testJWTSecret = "test-secret"

func setupTestRouter(t *testing.T) (*service.AuthService, *httptest.Server) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&model.User{}, &model.Course{}, &model.Category{}, &model.Tag{}, &model.CourseResource{}, &model.UserFavorite{},
		&model.CoinAccount{}, &model.CoinTransaction{}, &model.VipPackage{}, &model.UserVip{}, &model.Order{}, &model.Purchase{},
		&model.Guestbook{}, &model.GuestbookLike{}, &model.Comment{}, &model.NavItem{}, &model.Banner{}, &model.UserDownload{},
		&model.Ticket{}, &model.TicketReply{}, &model.TicketAttachment{}, &model.Setting{}, &model.OperationLog{}, &model.PaymentLog{}, &model.WithdrawalRequest{}))

	userRepo := repository.NewUserRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	catRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)
	favRepo := repository.NewFavoriteRepository(db)
	payRepo := repository.NewPaymentRepository(db)
	commRepo := repository.NewCommunityRepository(db)
	ticketRepo := repository.NewTicketRepository(db)

	authSvc := service.NewAuthService(userRepo, testJWTSecret, 72*time.Hour)
	courseSvc := service.NewCourseService(courseRepo, catRepo, tagRepo, favRepo)
	paySvc := service.NewPaymentService(payRepo, courseRepo, userRepo, payment.NewWechatPay(payment.WechatPayConfig{}), payment.NewAlipayClient(payment.AlipayConfig{}))
	commSvc := service.NewCommunityService(commRepo)
	ticketSvc := service.NewTicketService(ticketRepo, userRepo)

	authHandler := api.NewAuthHandler(authSvc)
	courseHandler := api.NewCourseHandler(courseSvc)
	adminHandler := api.NewAdminHandler(courseSvc)
	payHandler := api.NewPaymentHandler(paySvc)
	commHandler := api.NewCommunityHandler(commSvc)
	adminPayHandler := api.NewAdminPaymentHandler(paySvc, commSvc)
	ticketHandler := api.NewTicketHandler(ticketSvc, authSvc, paySvc, oauth.NewWechatOAuth(oauth.WechatOAuthConfig{}), testJWTSecret, 72*time.Hour, t.TempDir(), nil)

	r := api.SetupRouter(authHandler, courseHandler, adminHandler, payHandler, commHandler, adminPayHandler, api.NewUploadHandler(t.TempDir(), nil), ticketHandler, testJWTSecret)
	ts := httptest.NewServer(r)
	return authSvc, ts
}

type apiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func postJSON(url string, body interface{}) (*http.Response, *apiResponse) {
	b, _ := json.Marshal(body)
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(b))
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	return resp, &result
}

func getJSON(url string, token string) (*http.Response, *apiResponse) {
	req, _ := http.NewRequest("GET", url, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	return resp, &result
}

// === RED-GREEN Tests ===

func Test_Captcha_返回验证码(t *testing.T) {
	_, ts := setupTestRouter(t)
	defer ts.Close()

	_, result := postJSON(ts.URL+"/api/v1/auth/captcha", nil)
	assert.Equal(t, 0, result.Code)

	var data model.CaptchaResponse
	json.Unmarshal(result.Data, &data)
	assert.NotEmpty(t, data.CaptchaKey)
	assert.NotEmpty(t, data.CaptchaImage)
}

func Test_SendEmail_正确验证码发送成功(t *testing.T) {
	authSvc, ts := setupTestRouter(t)
	defer ts.Close()

	authSvc.SetCaptcha("email-key", "1234")
	_, result := postJSON(ts.URL+"/api/v1/auth/email/send", map[string]string{
		"email":       "user@example.com",
		"captcha":     "1234",
		"captcha_key": "email-key",
	})
	assert.Equal(t, 0, result.Code)
}

func Test_Register_缺少必填字段返回40001(t *testing.T) {
	_, ts := setupTestRouter(t)
	defer ts.Close()

	_, result := postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"email": "user@example.com",
	})
	assert.Equal(t, 40001, result.Code)
}

func Test_Register_邮箱验证码错误返回40001(t *testing.T) {
	_, ts := setupTestRouter(t)
	defer ts.Close()

	_, result := postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"email":      "user@example.com",
		"email_code": "wrong",
		"password":   "123456",
	})
	assert.Equal(t, 40001, result.Code)
}

func Test_Register_成功注册返回token(t *testing.T) {
	authSvc, ts := setupTestRouter(t)
	defer ts.Close()

	authSvc.SetEmailCode("user@example.com", "123456")
	_, result := postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"username":   "paint_user",
		"email":      "user@example.com",
		"email_code": "123456",
		"password":   "password123",
	})
	assert.Equal(t, 0, result.Code)

	var data model.AuthResponse
	json.Unmarshal(result.Data, &data)
	assert.NotEmpty(t, data.Token)
	assert.Equal(t, "user@example.com", data.User.Email)
	assert.Equal(t, "paint_user", data.User.Username)
}

func Test_Register_邮箱重复注册返回40001(t *testing.T) {
	authSvc, ts := setupTestRouter(t)
	defer ts.Close()

	// First registration
	authSvc.SetEmailCode("user@example.com", "111111")
	postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"email":      "user@example.com",
		"email_code": "111111",
		"password":   "password123",
	})

	// Duplicate registration
	authSvc.SetEmailCode("user@example.com", "222222")
	_, result := postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"email":      "user@example.com",
		"email_code": "222222",
		"password":   "password123",
	})
	assert.Equal(t, 40001, result.Code)
	assert.Contains(t, result.Message, "已注册")
}

func Test_Register_用户名重复返回40001(t *testing.T) {
	authSvc, ts := setupTestRouter(t)
	defer ts.Close()

	authSvc.SetEmailCode("user1@example.com", "111111")
	postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"username":   "same_name",
		"email":      "user1@example.com",
		"email_code": "111111",
		"password":   "password123",
	})

	authSvc.SetEmailCode("user2@example.com", "222222")
	_, result := postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"username":   "same_name",
		"email":      "user2@example.com",
		"email_code": "222222",
		"password":   "password123",
	})
	assert.Equal(t, 40001, result.Code)
	assert.Contains(t, result.Message, "用户名已存在")
}

func Test_Login_密码错误返回40001(t *testing.T) {
	authSvc, ts := setupTestRouter(t)
	defer ts.Close()

	// Register first
	authSvc.SetEmailCode("login@example.com", "111111")
	postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"email":      "login@example.com",
		"email_code": "111111",
		"password":   "correctpass",
	})

	// Login with wrong password
	authSvc.SetCaptcha("cap-key", "abcd")
	_, result := postJSON(ts.URL+"/api/v1/auth/login", map[string]string{
		"email":       "login@example.com",
		"password":    "wrongpass",
		"captcha":     "abcd",
		"captcha_key": "cap-key",
	})
	assert.Equal(t, 40001, result.Code)
}

func Test_Login_成功登录返回token(t *testing.T) {
	authSvc, ts := setupTestRouter(t)
	defer ts.Close()

	authSvc.SetEmailCode("login@example.com", "111111")
	postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"email":      "login@example.com",
		"email_code": "111111",
		"password":   "mypassword",
	})

	authSvc.SetCaptcha("cap-key", "9999")
	_, result := postJSON(ts.URL+"/api/v1/auth/login", map[string]string{
		"email":       "login@example.com",
		"password":    "mypassword",
		"captcha":     "9999",
		"captcha_key": "cap-key",
	})
	assert.Equal(t, 0, result.Code)

	var data model.AuthResponse
	json.Unmarshal(result.Data, &data)
	assert.NotEmpty(t, data.Token)
}

func Test_Profile_未携带Token返回40101(t *testing.T) {
	_, ts := setupTestRouter(t)
	defer ts.Close()

	_, result := getJSON(ts.URL+"/api/v1/user/profile", "")
	assert.Equal(t, 40101, result.Code)
}

func Test_Profile_携带Token返回用户信息(t *testing.T) {
	authSvc, ts := setupTestRouter(t)
	defer ts.Close()

	authSvc.SetEmailCode("profile@example.com", "111111")
	_, regResult := postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"email":      "profile@example.com",
		"email_code": "111111",
		"password":   "mypassword",
	})

	var authData model.AuthResponse
	json.Unmarshal(regResult.Data, &authData)

	_, result := getJSON(ts.URL+"/api/v1/user/profile", authData.Token)
	assert.Equal(t, 0, result.Code)

	var user model.UserResponse
	json.Unmarshal(result.Data, &user)
	assert.Equal(t, "profile@example.com", user.Email)
}

func Test_UpdateProfile_更新邮箱(t *testing.T) {
	authSvc, ts := setupTestRouter(t)
	defer ts.Close()

	authSvc.SetEmailCode("old@example.com", "111111")
	_, regResult := postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"email":      "old@example.com",
		"email_code": "111111",
		"password":   "mypassword",
	})

	var authData model.AuthResponse
	json.Unmarshal(regResult.Data, &authData)

	b, _ := json.Marshal(map[string]string{
		"username": "newname",
		"email":    "new@example.com",
	})
	req, _ := http.NewRequest("PUT", ts.URL+"/api/v1/user/profile", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authData.Token)
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	assert.Equal(t, 0, result.Code)

	var user model.UserResponse
	json.Unmarshal(result.Data, &user)
	assert.Equal(t, "newname", user.Username)
	assert.Equal(t, "new@example.com", user.Email)
}

func Test_全流程_注册登录获取Profile(t *testing.T) {
	authSvc, ts := setupTestRouter(t)
	defer ts.Close()

	// 1. Register
	authSvc.SetEmailCode("flow@example.com", "666666")
	_, regResult := postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"email":      "flow@example.com",
		"email_code": "666666",
		"password":   "testpass123",
	})
	assert.Equal(t, 0, regResult.Code)

	// 2. Login
	authSvc.SetCaptcha("flow-cap", "1111")
	_, loginResult := postJSON(ts.URL+"/api/v1/auth/login", map[string]string{
		"email":       "flow@example.com",
		"password":    "testpass123",
		"captcha":     "1111",
		"captcha_key": "flow-cap",
	})
	assert.Equal(t, 0, loginResult.Code)

	var authData model.AuthResponse
	json.Unmarshal(loginResult.Data, &authData)
	assert.NotEmpty(t, authData.Token)

	// 3. Get Profile
	_, profileResult := getJSON(ts.URL+"/api/v1/user/profile", authData.Token)
	assert.Equal(t, 0, profileResult.Code)

	var user model.UserResponse
	json.Unmarshal(profileResult.Data, &user)
	assert.Equal(t, "flow@example.com", user.Email)
	assert.NotEmpty(t, user.Username)
}
