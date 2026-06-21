package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zioran/backend/internal/api"
	"github.com/zioran/backend/internal/middleware"
	"github.com/zioran/backend/internal/model"
	"github.com/zioran/backend/internal/repository"
	"github.com/zioran/backend/internal/service"
	"github.com/zioran/backend/pkg/oauth"
	"github.com/zioran/backend/pkg/payment"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTicketTestRouter(t *testing.T) (*gorm.DB, *httptest.Server, *service.AuthService, string) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(
		&model.User{}, &model.Course{}, &model.Category{}, &model.Tag{},
		&model.CourseResource{}, &model.UserFavorite{},
		&model.CoinAccount{}, &model.CoinTransaction{}, &model.VipPackage{},
		&model.UserVip{}, &model.Order{}, &model.Purchase{},
		&model.Guestbook{}, &model.GuestbookLike{}, &model.Comment{},
		&model.NavItem{}, &model.Banner{}, &model.UserDownload{},
		&model.Ticket{}, &model.TicketReply{}, &model.Setting{},
		&model.OperationLog{}, &model.PaymentLog{}, &model.WithdrawalRequest{},
	))

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
	ticketHandler := api.NewTicketHandler(ticketSvc, authSvc, paySvc, oauth.NewWechatOAuth(oauth.WechatOAuthConfig{}), testJWTSecret, 72*time.Hour, t.TempDir())

	r := api.SetupRouter(authHandler, courseHandler, adminHandler, payHandler, commHandler, adminPayHandler, api.NewUploadHandler(t.TempDir()), ticketHandler, testJWTSecret)
	ts := httptest.NewServer(r)

	authSvc.SetCaptcha("cap1", "1234")
	return db, ts, authSvc, ts.URL
}

func createTestUser(t *testing.T, db *gorm.DB) (int64, string) {
	user := &model.User{
		Username: "testuser", Phone: "13800000001", Email: "test@example.com",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // "password"
		Role:         "user", Status: "active",
	}
	db.Create(user)
	token, _ := middleware.GenerateToken(user.ID, testJWTSecret, 72*time.Hour)
	return user.ID, token
}

func createTestAdmin(t *testing.T, db *gorm.DB) (int64, string) {
	user := &model.User{
		Username: "admin", Phone: "13800000002",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
		Role:         "admin", Status: "active",
	}
	db.Create(user)
	token, _ := middleware.GenerateToken(user.ID, testJWTSecret, 72*time.Hour)
	return user.ID, token
}

func doPost(url, token string, body interface{}) *apiResponse {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	return &result
}

func doGet(url, token string) *apiResponse {
	req, _ := http.NewRequest("GET", url, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	return &result
}

func doPut(url, token string, body interface{}) *apiResponse {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	return &result
}

func doDelete(url, token string) *apiResponse {
	req, _ := http.NewRequest("DELETE", url, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	return &result
}

// === Test ForgotPassword ===

func TestForgotPassword_ByEmail(t *testing.T) {
	db, ts, authSvc, baseURL := setupTicketTestRouter(t)
	defer ts.Close()
	createTestUser(t, db)

	authSvc.SetEmailCode("test@example.com", "654321")

	resp := doPost(baseURL+"/api/v1/auth/forgot-password", "", map[string]string{
		"email": " Test@Example.com ", "email_code": "654321", "new_password": "newpass123",
	})
	assert.Equal(t, 0, resp.Code)
}

func TestForgotPassword_WrongCode(t *testing.T) {
	db, ts, _, baseURL := setupTicketTestRouter(t)
	defer ts.Close()
	createTestUser(t, db)

	resp := doPost(baseURL+"/api/v1/auth/forgot-password", "", map[string]string{
		"email": "test@example.com", "email_code": "000000", "new_password": "newpass123",
	})
	assert.NotEqual(t, 0, resp.Code)
}

// === Test RefreshToken ===

func TestRefreshToken(t *testing.T) {
	db, ts, _, baseURL := setupTicketTestRouter(t)
	defer ts.Close()
	_, token := createTestUser(t, db)

	resp := doPost(baseURL+"/api/v1/auth/refresh", token, nil)
	assert.Equal(t, 0, resp.Code)
	var data map[string]string
	json.Unmarshal(resp.Data, &data)
	assert.NotEmpty(t, data["token"])
}

// === Test Order Cancel ===

func TestOrderCancel(t *testing.T) {
	db, ts, _, baseURL := setupTicketTestRouter(t)
	defer ts.Close()
	userID, token := createTestUser(t, db)

	// Create a pending order
	order := &model.Order{
		OrderNo: "ORD_TEST_001", UserID: userID, Type: "course",
		TargetName: "test", Amount: 100, Status: "pending",
	}
	db.Create(order)

	resp := doPost(fmt.Sprintf("%s/api/v1/orders/%d/cancel", baseURL, order.ID), token, nil)
	assert.Equal(t, 0, resp.Code)

	// Verify order is cancelled
	var updated model.Order
	db.First(&updated, order.ID)
	assert.Equal(t, "cancelled", updated.Status)
}

func TestOrderCancel_NotPending(t *testing.T) {
	db, ts, _, baseURL := setupTicketTestRouter(t)
	defer ts.Close()
	userID, token := createTestUser(t, db)

	order := &model.Order{
		OrderNo: "ORD_TEST_002", UserID: userID, Type: "course",
		TargetName: "test", Amount: 100, Status: "paid",
	}
	db.Create(order)

	resp := doPost(fmt.Sprintf("%s/api/v1/orders/%d/cancel", baseURL, order.ID), token, nil)
	assert.NotEqual(t, 0, resp.Code)
}

// === Test Tickets (user) ===

func TestTicketCRUD(t *testing.T) {
	db, ts, _, baseURL := setupTicketTestRouter(t)
	defer ts.Close()
	_, token := createTestUser(t, db)

	// Create
	resp := doPost(baseURL+"/api/v1/tickets", token, map[string]string{
		"title": "test ticket", "content": "help me",
	})
	assert.Equal(t, 0, resp.Code)

	// List
	resp = doGet(baseURL+"/api/v1/tickets", token)
	assert.Equal(t, 0, resp.Code)

	// Detail
	resp = doGet(baseURL+"/api/v1/tickets/1", token)
	assert.Equal(t, 0, resp.Code)

	// Reply
	resp = doPost(baseURL+"/api/v1/tickets/1/reply", token, map[string]string{
		"content": "more info here",
	})
	assert.Equal(t, 0, resp.Code)
}

// === Test Tickets (admin) ===

func TestAdminTicketOperations(t *testing.T) {
	db, ts, _, baseURL := setupTicketTestRouter(t)
	defer ts.Close()
	_, userToken := createTestUser(t, db)
	_, adminToken := createTestAdmin(t, db)

	// User creates ticket
	doPost(baseURL+"/api/v1/tickets", userToken, map[string]string{
		"title": "need help", "content": "my order failed",
	})

	// Admin lists all tickets
	resp := doGet(baseURL+"/api/v1/admin/tickets", adminToken)
	assert.Equal(t, 0, resp.Code)

	// Admin views detail
	resp = doGet(baseURL+"/api/v1/admin/tickets/1", adminToken)
	assert.Equal(t, 0, resp.Code)

	// Admin updates status
	resp = doPut(baseURL+"/api/v1/admin/tickets/1/status", adminToken, map[string]string{
		"status": "processing",
	})
	assert.Equal(t, 0, resp.Code)

	// Admin replies
	resp = doPost(baseURL+"/api/v1/admin/tickets/1/reply", adminToken, map[string]string{
		"content": "we are looking into this",
	})
	assert.Equal(t, 0, resp.Code)
}

func TestAdminTicketList_StatusFilter(t *testing.T) {
	db, ts, _, baseURL := setupTicketTestRouter(t)
	defer ts.Close()
	userID, _ := createTestUser(t, db)
	_, adminToken := createTestAdmin(t, db)

	db.Create(&model.Ticket{UserID: userID, Title: "open ticket", Content: "a", Status: "open"})
	db.Create(&model.Ticket{UserID: userID, Title: "closed ticket", Content: "b", Status: "closed"})

	resp := doGet(baseURL+"/api/v1/admin/tickets?status=pending", adminToken)
	assert.Equal(t, 0, resp.Code)
	var page struct {
		Items []model.TicketResponse `json:"items"`
		Total int64                  `json:"total"`
	}
	json.Unmarshal(resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)
	if assert.Len(t, page.Items, 1) {
		assert.Equal(t, "open", page.Items[0].Status)
	}
}

// === Test Settings ===

func TestSettings(t *testing.T) {
	_, ts, _, baseURL := setupTicketTestRouter(t)
	defer ts.Close()
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	_ = db

	// Use the running ts
	db2, ts2, _, baseURL2 := setupTicketTestRouter(t)
	defer ts2.Close()
	_, adminToken := createTestAdmin(t, db2)

	// Update settings
	resp := doPut(baseURL2+"/api/v1/admin/settings", adminToken, map[string]string{
		"site_name": "知猿", "site_description": "知猿学堂",
	})
	assert.Equal(t, 0, resp.Code)

	// Get settings
	resp = doGet(baseURL2+"/api/v1/admin/settings", adminToken)
	assert.Equal(t, 0, resp.Code)

	_ = baseURL // suppress unused
}

// === Test Admin Account Management ===

func TestAdminAccountManagement(t *testing.T) {
	db, ts, _, baseURL := setupTicketTestRouter(t)
	defer ts.Close()
	_, adminToken := createTestAdmin(t, db)

	// Create admin
	resp := doPost(baseURL+"/api/v1/admin/admins", adminToken, map[string]string{
		"username": "newadmin", "password": "pass123456", "role": "admin",
	})
	assert.Equal(t, 0, resp.Code)

	// List admins
	resp = doGet(baseURL+"/api/v1/admin/admins", adminToken)
	assert.Equal(t, 0, resp.Code)

	// Update admin
	resp = doPut(baseURL+"/api/v1/admin/admins/2", adminToken, map[string]string{
		"role": "admin",
	})
	assert.Equal(t, 0, resp.Code)

	// Delete admin
	resp = doDelete(baseURL+"/api/v1/admin/admins/2", adminToken)
	assert.Equal(t, 0, resp.Code)
}

// === Test Finance ===

func TestFinanceEndpoints(t *testing.T) {
	db, ts, _, baseURL := setupTicketTestRouter(t)
	defer ts.Close()
	userID, _ := createTestUser(t, db)
	_, adminToken := createTestAdmin(t, db)
	db.Create(&model.WithdrawalRequest{
		UserID: userID, Amount: 100, AccountName: "tester", AccountNo: "NO1",
		BankName: "bank", Status: "pending",
	})
	db.Create(&model.WithdrawalRequest{
		UserID: userID, Amount: 200, AccountName: "tester", AccountNo: "NO2",
		BankName: "bank", Status: "approved",
	})

	resp := doGet(baseURL+"/api/v1/admin/finance/summary", adminToken)
	assert.Equal(t, 0, resp.Code)

	resp = doGet(baseURL+"/api/v1/admin/finance/withdrawals?status=pending", adminToken)
	assert.Equal(t, 0, resp.Code)
	var page struct {
		Items []model.FinanceWithdrawalResponse `json:"items"`
		Total int64                             `json:"total"`
	}
	json.Unmarshal(resp.Data, &page)
	assert.Equal(t, int64(1), page.Total)
	if assert.Len(t, page.Items, 1) {
		assert.Equal(t, 100, page.Items[0].Amount)
		assert.Equal(t, "testuser", page.Items[0].Username)
	}
}

// === Test Logs ===

func TestLogEndpoints(t *testing.T) {
	db, ts, _, baseURL := setupTicketTestRouter(t)
	defer ts.Close()
	_, adminToken := createTestAdmin(t, db)

	resp := doGet(baseURL+"/api/v1/admin/logs/operations", adminToken)
	assert.Equal(t, 0, resp.Code)

	resp = doGet(baseURL+"/api/v1/admin/logs/payments", adminToken)
	assert.Equal(t, 0, resp.Code)
}

// === Test Comment Admin Reply ===

func TestAdminCommentReply(t *testing.T) {
	db, ts, _, baseURL := setupTicketTestRouter(t)
	defer ts.Close()
	userID, _ := createTestUser(t, db)
	_, adminToken := createTestAdmin(t, db)

	// Create a comment first
	comment := &model.Comment{
		UserID: userID, TargetType: "course", TargetID: 1,
		Content: "nice course!", Status: "visible",
	}
	db.Create(comment)

	resp := doPost(fmt.Sprintf("%s/api/v1/admin/comments/%d/reply", baseURL, comment.ID), adminToken, map[string]string{
		"content": "thanks for your feedback!",
	})
	assert.Equal(t, 0, resp.Code)
}

// === Test Payment Notify (MOCK) ===

func TestPaymentNotifyMock(t *testing.T) {
	_, ts, _, baseURL := setupTicketTestRouter(t)
	defer ts.Close()

	// Wechat notify
	resp := doPost(baseURL+"/api/v1/pay/notify/wechat", "", map[string]int64{"order_id": 1})
	assert.Equal(t, 0, resp.Code)

	// Alipay notify
	resp = doPost(baseURL+"/api/v1/pay/notify/alipay", "", nil)
	assert.Equal(t, 0, resp.Code)
}

// === Test OAuth Mock ===

func TestOAuthMock(t *testing.T) {
	_, ts, _, baseURL := setupTicketTestRouter(t)
	defer ts.Close()

	// Get wechat auth URL
	resp := doGet(baseURL+"/api/v1/auth/oauth/wechat", "")
	assert.Equal(t, 0, resp.Code)
	var authPayload map[string]string
	assert.NoError(t, json.Unmarshal(resp.Data, &authPayload))
	authURL, err := url.Parse(authPayload["auth_url"])
	assert.NoError(t, err)
	assert.Equal(t, "/connect/qrconnect", authURL.Path)
	assert.Equal(t, "snsapi_login", authURL.Query().Get("scope"))

	// Wechat callback
	resp = doPost(baseURL+"/api/v1/auth/oauth/wechat/callback", "", map[string]string{"code": "test"})
	assert.Equal(t, 0, resp.Code)
	var loginPayload map[string]interface{}
	assert.NoError(t, json.Unmarshal(resp.Data, &loginPayload))
	assert.NotEmpty(t, loginPayload["token"])

	resp = doGet(baseURL+"/api/v1/auth/oauth/wechat/callback?code=test&state=login", "")
	assert.Equal(t, 0, resp.Code)
}

// === Test Home Config ===

func TestHomeConfig(t *testing.T) {
	_, ts, _, baseURL := setupTicketTestRouter(t)
	defer ts.Close()

	resp := doGet(baseURL+"/api/v1/home/config", "")
	assert.Equal(t, 0, resp.Code)
}

// === Test User Order Detail ===

func TestUserOrderDetail(t *testing.T) {
	db, ts, _, baseURL := setupTicketTestRouter(t)
	defer ts.Close()
	userID, token := createTestUser(t, db)

	order := &model.Order{
		OrderNo: "ORD_TEST_003", UserID: userID, Type: "course",
		TargetName: "test course", Amount: 50, Status: "paid",
	}
	db.Create(order)

	resp := doGet(fmt.Sprintf("%s/api/v1/user/orders/%d", baseURL, order.ID), token)
	assert.Equal(t, 0, resp.Code)
}
