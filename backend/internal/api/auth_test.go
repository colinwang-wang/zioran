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
	"github.com/zioran/backend/pkg/response"
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
		&model.Ticket{}, &model.TicketReply{}, &model.Setting{}, &model.OperationLog{}, &model.PaymentLog{}))

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
	paySvc := service.NewPaymentService(payRepo, courseRepo, userRepo)
	commSvc := service.NewCommunityService(commRepo)
	ticketSvc := service.NewTicketService(ticketRepo, userRepo)

	authHandler := api.NewAuthHandler(authSvc)
	courseHandler := api.NewCourseHandler(courseSvc)
	adminHandler := api.NewAdminHandler(courseSvc)
	payHandler := api.NewPaymentHandler(paySvc)
	commHandler := api.NewCommunityHandler(commSvc)
	adminPayHandler := api.NewAdminPaymentHandler(paySvc, commSvc)
	ticketHandler := api.NewTicketHandler(ticketSvc, authSvc, testJWTSecret, 72*time.Hour, t.TempDir())

	r := api.SetupRouter(authHandler, courseHandler, adminHandler, payHandler, commHandler, adminPayHandler, api.NewUploadHandler(t.TempDir()), ticketHandler, testJWTSecret)
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

func Test_SendSMS_验证码错误返回40001(t *testing.T) {
	_, ts := setupTestRouter(t)
	defer ts.Close()

	_, result := postJSON(ts.URL+"/api/v1/auth/sms/send", map[string]string{
		"phone":       "13800138000",
		"captcha":     "wrong",
		"captcha_key": "nonexistent",
	})
	assert.Equal(t, 40001, result.Code)
}

func Test_SendSMS_正确验证码发送成功(t *testing.T) {
	authSvc, ts := setupTestRouter(t)
	defer ts.Close()

	authSvc.SetCaptcha("test-key", "1234")
	_, result := postJSON(ts.URL+"/api/v1/auth/sms/send", map[string]string{
		"phone":       "13800138000",
		"captcha":     "1234",
		"captcha_key": "test-key",
	})
	assert.Equal(t, 0, result.Code)
}

func Test_Register_缺少必填字段返回40001(t *testing.T) {
	_, ts := setupTestRouter(t)
	defer ts.Close()

	_, result := postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"phone": "13800138000",
	})
	assert.Equal(t, 40001, result.Code)
}

func Test_Register_短信验证码错误返回40001(t *testing.T) {
	_, ts := setupTestRouter(t)
	defer ts.Close()

	_, result := postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"phone":    "13800138000",
		"sms_code": "wrong",
		"password": "123456",
	})
	assert.Equal(t, 40001, result.Code)
}

func Test_Register_成功注册返回token(t *testing.T) {
	authSvc, ts := setupTestRouter(t)
	defer ts.Close()

	authSvc.SetSMSCode("13800138000", "123456")
	_, result := postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"phone":    "13800138000",
		"sms_code": "123456",
		"password": "password123",
	})
	assert.Equal(t, 0, result.Code)

	var data model.AuthResponse
	json.Unmarshal(result.Data, &data)
	assert.NotEmpty(t, data.Token)
	assert.Equal(t, "138****8000", data.User.Phone)
}

func Test_Register_手机号重复注册返回40001(t *testing.T) {
	authSvc, ts := setupTestRouter(t)
	defer ts.Close()

	// First registration
	authSvc.SetSMSCode("13800138000", "111111")
	postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"phone":    "13800138000",
		"sms_code": "111111",
		"password": "password123",
	})

	// Duplicate registration
	authSvc.SetSMSCode("13800138000", "222222")
	_, result := postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"phone":    "13800138000",
		"sms_code": "222222",
		"password": "password123",
	})
	assert.Equal(t, 40001, result.Code)
	assert.Contains(t, result.Message, "已注册")
}

func Test_Login_密码错误返回40001(t *testing.T) {
	authSvc, ts := setupTestRouter(t)
	defer ts.Close()

	// Register first
	authSvc.SetSMSCode("13800138000", "111111")
	postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"phone":    "13800138000",
		"sms_code": "111111",
		"password": "correctpass",
	})

	// Login with wrong password
	authSvc.SetCaptcha("cap-key", "abcd")
	_, result := postJSON(ts.URL+"/api/v1/auth/login", map[string]string{
		"phone":       "13800138000",
		"password":    "wrongpass",
		"captcha":     "abcd",
		"captcha_key": "cap-key",
	})
	assert.Equal(t, 40001, result.Code)
}

func Test_Login_成功登录返回token(t *testing.T) {
	authSvc, ts := setupTestRouter(t)
	defer ts.Close()

	authSvc.SetSMSCode("13800138000", "111111")
	postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"phone":    "13800138000",
		"sms_code": "111111",
		"password": "mypassword",
	})

	authSvc.SetCaptcha("cap-key", "9999")
	_, result := postJSON(ts.URL+"/api/v1/auth/login", map[string]string{
		"phone":       "13800138000",
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

	authSvc.SetSMSCode("13800138000", "111111")
	_, regResult := postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"phone":    "13800138000",
		"sms_code": "111111",
		"password": "mypassword",
	})

	var authData model.AuthResponse
	json.Unmarshal(regResult.Data, &authData)

	_, result := getJSON(ts.URL+"/api/v1/user/profile", authData.Token)
	assert.Equal(t, 0, result.Code)

	var profileData response.Response
	json.Unmarshal(result.Data, &profileData)

	var user model.UserResponse
	json.Unmarshal(result.Data, &user)
	assert.Equal(t, "138****8000", user.Phone)
}

func Test_全流程_注册登录获取Profile(t *testing.T) {
	authSvc, ts := setupTestRouter(t)
	defer ts.Close()

	// 1. Register
	authSvc.SetSMSCode("13900139000", "666666")
	_, regResult := postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"phone":    "13900139000",
		"sms_code": "666666",
		"password": "testpass123",
	})
	assert.Equal(t, 0, regResult.Code)

	// 2. Login
	authSvc.SetCaptcha("flow-cap", "1111")
	_, loginResult := postJSON(ts.URL+"/api/v1/auth/login", map[string]string{
		"phone":       "13900139000",
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
	assert.Equal(t, "139****9000", user.Phone)
	assert.Equal(t, "user_9000", user.Username)
}
