package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func setupPhase34Router(t *testing.T) (*gorm.DB, *httptest.Server, string) {
	return setupPhase34RouterWithMockPayment(t, true)
}

func setupPhase34RouterWithMockPayment(t *testing.T, mockAutoComplete bool) (*gorm.DB, *httptest.Server, string) {
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
	paySvc := service.NewPaymentService(
		payRepo,
		courseRepo,
		userRepo,
		payment.NewWechatPay(payment.WechatPayConfig{MockAutoComplete: mockAutoComplete}),
		payment.NewAlipayClient(payment.AlipayConfig{MockAutoComplete: mockAutoComplete}),
	)
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

	// Create test user
	user := &model.User{Username: "testuser", Phone: "13800138000", PasswordHash: "$2a$10$dummy", Role: "user", Status: "active"}
	db.Create(user)
	token, _ := middleware.GenerateToken(user.ID, testJWTSecret, 72*time.Hour)

	return db, ts, token
}

func authedPost(url, token string, body interface{}) *apiResponse {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	return &result
}

func authedGet(url, token string) *apiResponse {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	return &result
}

func authedPut(url, token string, body interface{}) *apiResponse {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	return &result
}

func authedDelete(url, token string) *apiResponse {
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()
	return &result
}

// === 金币 ===

func Test_CoinBalance_初始为零(t *testing.T) {
	_, ts, token := setupPhase34Router(t)
	defer ts.Close()

	result := authedGet(ts.URL+"/api/v1/coins/balance", token)
	assert.Equal(t, 0, result.Code)

	var data model.CoinBalanceResponse
	json.Unmarshal(result.Data, &data)
	assert.Equal(t, 0, data.Balance)
}

func Test_Recharge_充值成功(t *testing.T) {
	_, ts, token := setupPhase34Router(t)
	defer ts.Close()

	result := authedPost(ts.URL+"/api/v1/coins/recharge", token, map[string]interface{}{
		"amount": 100, "pay_method": "wechat",
	})
	assert.Equal(t, 0, result.Code)

	// Check balance
	balResult := authedGet(ts.URL+"/api/v1/coins/balance", token)
	var bal model.CoinBalanceResponse
	json.Unmarshal(balResult.Data, &bal)
	assert.Equal(t, 100, bal.Balance)
}

func Test_Recharge_按后台比例到账(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	db.Create(&model.Setting{Key: "coinRechargeRatio", Value: "10"})

	result := authedPost(ts.URL+"/api/v1/coins/recharge", token, map[string]interface{}{
		"amount": 50, "pay_method": "wechat",
	})
	assert.Equal(t, 0, result.Code)

	var recharge model.RechargeResponse
	json.Unmarshal(result.Data, &recharge)
	assert.Equal(t, 50, recharge.Amount)
	assert.Equal(t, 500, recharge.Coins)

	balResult := authedGet(ts.URL+"/api/v1/coins/balance", token)
	var bal model.CoinBalanceResponse
	json.Unmarshal(balResult.Data, &bal)
	assert.Equal(t, 500, bal.Balance)
}

func Test_RechargeConfig_返回后台配置(t *testing.T) {
	db, ts, _ := setupPhase34Router(t)
	defer ts.Close()

	db.Create(&model.Setting{Key: "coinRechargeRatio", Value: "10"})
	db.Create(&model.Setting{Key: "coinRechargeAmounts", Value: "20,50,100"})

	_, result := getJSON(ts.URL+"/api/v1/coins/recharge-config", "")
	assert.Equal(t, 0, result.Code)

	var config model.RechargeConfigResponse
	json.Unmarshal(result.Data, &config)
	assert.Equal(t, 10, config.Ratio)
	assert.Equal(t, []int{20, 50, 100}, config.Amounts)
}

func Test_Recharge_DisabledPaymentDoesNotCreditBalance(t *testing.T) {
	_, ts, token := setupPhase34RouterWithMockPayment(t, false)
	defer ts.Close()

	result := authedPost(ts.URL+"/api/v1/coins/recharge", token, map[string]interface{}{
		"amount": 100, "pay_method": "wechat",
	})
	assert.NotEqual(t, 0, result.Code)

	balResult := authedGet(ts.URL+"/api/v1/coins/balance", token)
	var bal model.CoinBalanceResponse
	json.Unmarshal(balResult.Data, &bal)
	assert.Equal(t, 0, bal.Balance)
}

func Test_CoinTransactions_有记录(t *testing.T) {
	_, ts, token := setupPhase34Router(t)
	defer ts.Close()

	authedPost(ts.URL+"/api/v1/coins/recharge", token, map[string]interface{}{
		"amount": 50, "pay_method": "alipay",
	})

	result := authedGet(ts.URL+"/api/v1/coins/transactions", token)
	assert.Equal(t, 0, result.Code)
}

// === VIP ===

func Test_VipPackages_公开列表(t *testing.T) {
	db, ts, _ := setupPhase34Router(t)
	defer ts.Close()

	db.Create(&model.VipPackage{Name: "终身VIP", Price: 99, OriginalPrice: 699, IsActive: true})
	_, result := getJSON(ts.URL+"/api/v1/vip/packages", "")
	assert.Equal(t, 0, result.Code)
}

func Test_VipStatus_未开通(t *testing.T) {
	_, ts, token := setupPhase34Router(t)
	defer ts.Close()

	result := authedGet(ts.URL+"/api/v1/vip/status", token)
	assert.Equal(t, 0, result.Code)

	var data model.VipStatusResponse
	json.Unmarshal(result.Data, &data)
	assert.False(t, data.IsVip)
}

func Test_VipPurchase_余额不足(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	db.Create(&model.VipPackage{Name: "终身VIP", Price: 99, IsActive: true})
	result := authedPost(ts.URL+"/api/v1/vip/purchase", token, map[string]interface{}{
		"package_id": 1,
	})
	assert.NotEqual(t, 0, result.Code) // 余额不足
}

func Test_VipPurchase_成功(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	db.Create(&model.VipPackage{Name: "终身VIP", Price: 99, IsActive: true})

	// Recharge first
	authedPost(ts.URL+"/api/v1/coins/recharge", token, map[string]interface{}{
		"amount": 100, "pay_method": "wechat",
	})

	// Purchase VIP
	result := authedPost(ts.URL+"/api/v1/vip/purchase", token, map[string]interface{}{
		"package_id": 1,
	})
	assert.Equal(t, 0, result.Code)

	// Verify VIP status
	statusResult := authedGet(ts.URL+"/api/v1/vip/status", token)
	var status model.VipStatusResponse
	json.Unmarshal(statusResult.Data, &status)
	assert.True(t, status.IsVip)

	// Verify balance deducted
	balResult := authedGet(ts.URL+"/api/v1/coins/balance", token)
	var bal model.CoinBalanceResponse
	json.Unmarshal(balResult.Data, &bal)
	assert.Equal(t, 1, bal.Balance) // 100 - 99 = 1
}

// === 课程购买全流程 ===

func Test_PurchaseCourse_全流程(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	// Setup: create category and course
	cat := model.Category{Name: "设计", Slug: "design", IsActive: true}
	db.Create(&cat)
	course := model.Course{
		Title: "测试课程", Slug: "test-course", CategoryID: cat.ID,
		Price: 10, Status: "published",
		Resources: []model.CourseResource{{Name: "资源1", URL: "https://pan.baidu.com/xxx", Password: "1234"}},
	}
	db.Create(&course)

	// 1. Recharge
	authedPost(ts.URL+"/api/v1/coins/recharge", token, map[string]interface{}{
		"amount": 20, "pay_method": "wechat",
	})

	// 2. Purchase course via orders API
	result := authedPost(ts.URL+"/api/v1/orders", token, map[string]interface{}{
		"type": "course", "target_id": course.ID,
	})
	assert.Equal(t, 0, result.Code)

	// 3. Download
	dlResult := authedPost(fmt.Sprintf("%s/api/v1/courses/%d/download", ts.URL, course.ID), token, nil)
	assert.Equal(t, 0, dlResult.Code)

	var dlData model.CourseDownloadResponse
	json.Unmarshal(dlResult.Data, &dlData)
	assert.Len(t, dlData.Resources, 1)
	assert.Equal(t, "1234", dlData.Resources[0].Password)

	// 4. Verify balance
	balResult := authedGet(ts.URL+"/api/v1/coins/balance", token)
	var bal model.CoinBalanceResponse
	json.Unmarshal(balResult.Data, &bal)
	assert.Equal(t, 10, bal.Balance) // 20 - 10 = 10
}

func Test_Download_未购买被拒(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	cat := model.Category{Name: "设计", Slug: "design", IsActive: true}
	db.Create(&cat)
	course := model.Course{Title: "课程", Slug: "c1", CategoryID: cat.ID, Price: 10, Status: "published"}
	db.Create(&course)

	result := authedPost(fmt.Sprintf("%s/api/v1/courses/%d/download", ts.URL, course.ID), token, nil)
	assert.NotEqual(t, 0, result.Code) // 未购买
}

func Test_VIP免费下载(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	cat := model.Category{Name: "设计", Slug: "design", IsActive: true}
	db.Create(&cat)
	course := model.Course{
		Title: "课程", Slug: "c2", CategoryID: cat.ID, Price: 10, Status: "published",
		Resources: []model.CourseResource{{Name: "R1", URL: "http://dl.com/1"}},
	}
	db.Create(&course)

	// Make user VIP directly
	db.Create(&model.UserVip{UserID: 1, PackageID: 1, StartedAt: time.Now(), IsActive: true})

	// Download should work
	result := authedPost(fmt.Sprintf("%s/api/v1/courses/%d/download", ts.URL, course.ID), token, nil)
	assert.Equal(t, 0, result.Code)
}

// === 留言板 ===

func Test_Guestbook_CRUD(t *testing.T) {
	_, ts, token := setupPhase34Router(t)
	defer ts.Close()

	// Create
	result := authedPost(ts.URL+"/api/v1/guestbook", token, map[string]string{"content": "求课程推荐"})
	assert.Equal(t, 0, result.Code)

	// List
	_, listResult := getJSON(ts.URL+"/api/v1/guestbook", "")
	assert.Equal(t, 0, listResult.Code)

	// Like
	likeResult := authedPost(ts.URL+"/api/v1/guestbook/1/like", token, nil)
	assert.Equal(t, 0, likeResult.Code)

	// Delete
	delResult := authedDelete(ts.URL+"/api/v1/guestbook/1", token)
	assert.Equal(t, 0, delResult.Code)
}

// === 评论 ===

func Test_Comment_CRUD(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	cat := model.Category{Name: "设计", Slug: "design", IsActive: true}
	db.Create(&cat)
	course := model.Course{Title: "课程", Slug: "c3", CategoryID: cat.ID, Status: "published"}
	db.Create(&course)

	// Create comment
	result := authedPost(ts.URL+"/api/v1/comments", token, map[string]interface{}{
		"target_type": "course", "target_id": course.ID, "content": "好课！",
	})
	assert.Equal(t, 0, result.Code)

	// List
	_, listResult := getJSON(fmt.Sprintf("%s/api/v1/comments?target_type=course&target_id=%d", ts.URL, course.ID), "")
	assert.Equal(t, 0, listResult.Code)

	// User comment list
	userComments := authedGet(ts.URL+"/api/v1/user/comments", token)
	assert.Equal(t, 0, userComments.Code)
	var page struct {
		Items []model.CommentResponse `json:"items"`
		Total int64                   `json:"total"`
	}
	json.Unmarshal(userComments.Data, &page)
	assert.Equal(t, int64(1), page.Total)
	if assert.Len(t, page.Items, 1) {
		assert.Equal(t, "course", page.Items[0].TargetType)
		assert.Equal(t, course.ID, page.Items[0].TargetID)
	}

	// Delete
	delResult := authedDelete(ts.URL+"/api/v1/comments/1", token)
	assert.Equal(t, 0, delResult.Code)
}

// === 首页配置 ===

func Test_HomeConfig_NavItems(t *testing.T) {
	db, ts, _ := setupPhase34Router(t)
	defer ts.Close()

	db.Create(&model.NavItem{Title: "最新", Icon: "new", URL: "/courses?sort=latest", IsActive: true, SortOrder: 1})
	_, result := getJSON(ts.URL+"/api/v1/home/nav-items", "")
	assert.Equal(t, 0, result.Code)
}

func Test_HomeConfig_Banners(t *testing.T) {
	db, ts, _ := setupPhase34Router(t)
	defer ts.Close()

	db.Create(&model.Banner{Title: "Banner1", ImageURL: "https://img.com/1.jpg", IsActive: true})
	_, result := getJSON(ts.URL+"/api/v1/home/banners", "")
	assert.Equal(t, 0, result.Code)
}

// === 个人中心 ===

func Test_UserOrders(t *testing.T) {
	_, ts, token := setupPhase34Router(t)
	defer ts.Close()

	// Recharge creates an order
	authedPost(ts.URL+"/api/v1/coins/recharge", token, map[string]interface{}{
		"amount": 50, "pay_method": "wechat",
	})

	result := authedGet(ts.URL+"/api/v1/user/orders", token)
	assert.Equal(t, 0, result.Code)
}

func Test_UserDownloads(t *testing.T) {
	_, ts, token := setupPhase34Router(t)
	defer ts.Close()

	result := authedGet(ts.URL+"/api/v1/user/downloads", token)
	assert.Equal(t, 0, result.Code)
}

func Test_ChangePassword_原密码错误(t *testing.T) {
	_, ts, token := setupPhase34Router(t)
	defer ts.Close()

	result := authedPut(ts.URL+"/api/v1/user/password", token, map[string]string{
		"old_password": "wrong", "new_password": "newpass123",
	})
	assert.NotEqual(t, 0, result.Code)
}

// === 后台管理 ===

func Test_AdminDashboard(t *testing.T) {
	_, ts, token := setupPhase34Router(t)
	defer ts.Close()

	result := authedGet(ts.URL+"/api/v1/admin/dashboard/stats", token)
	assert.Equal(t, 0, result.Code)
}

func Test_AdminOrders(t *testing.T) {
	_, ts, token := setupPhase34Router(t)
	defer ts.Close()

	result := authedGet(ts.URL+"/api/v1/admin/orders", token)
	assert.Equal(t, 0, result.Code)
}

func Test_AdminOrders_FilterByTypeStatusAndDate(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	user := &model.User{Username: "buyer", Phone: "13900000001", PasswordHash: "$2a$10$dummy", Role: "user", Status: "active"}
	db.Create(user)
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)
	db.Create(&model.Order{OrderNo: "ORD_FILTER_1", UserID: user.ID, Type: "course", TargetName: "course", Amount: 10, Status: "paid", CreatedAt: today})
	db.Create(&model.Order{OrderNo: "ORD_FILTER_2", UserID: user.ID, Type: "vip", TargetName: "vip", Amount: 20, Status: "paid", CreatedAt: today})
	db.Create(&model.Order{OrderNo: "ORD_FILTER_3", UserID: user.ID, Type: "course", TargetName: "old course", Amount: 30, Status: "pending", CreatedAt: yesterday})

	result := authedGet(ts.URL+"/api/v1/admin/orders?type=course_purchase&status=paid&startDate="+today.Format("2006-01-02")+"&endDate="+today.Format("2006-01-02"), token)
	assert.Equal(t, 0, result.Code)
	var page struct {
		Items []model.OrderResponse `json:"items"`
		Total int64                 `json:"total"`
	}
	json.Unmarshal(result.Data, &page)
	assert.Equal(t, int64(1), page.Total)
	if assert.Len(t, page.Items, 1) {
		assert.Equal(t, "ORD_FILTER_1", page.Items[0].OrderNo)
	}
}

func Test_DashboardCharts_UsesDatabaseCounts(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	result := authedGet(ts.URL+"/api/v1/admin/dashboard/charts?period=week", token)
	assert.Equal(t, 0, result.Code)
	var chart model.DashboardChartsResponse
	json.Unmarshal(result.Data, &chart)
	var orderSum int64
	for _, dataset := range chart.Datasets {
		if dataset.Label == "订单" {
			for _, value := range dataset.Data {
				orderSum += value
			}
		}
	}
	assert.Equal(t, int64(0), orderSum)

	db.Create(&model.Order{OrderNo: "ORD_CHART_1", UserID: 1, Type: "coin", TargetName: "coin", Amount: 10, Status: "paid"})
	result = authedGet(ts.URL+"/api/v1/admin/dashboard/charts?period=week", token)
	assert.Equal(t, 0, result.Code)
	json.Unmarshal(result.Data, &chart)
	orderSum = 0
	for _, dataset := range chart.Datasets {
		if dataset.Label == "订单" {
			for _, value := range dataset.Data {
				orderSum += value
			}
		}
	}
	assert.Equal(t, int64(1), orderSum)
}

func Test_GuestbookPin_UsesRequestedPinnedValue(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	entry := &model.Guestbook{UserID: 1, Content: "hello", Status: "visible", IsPinned: true}
	db.Create(entry)

	result := authedPut(ts.URL+fmt.Sprintf("/api/v1/admin/guestbook/%d/pin", entry.ID), token, map[string]bool{"pinned": false})
	assert.Equal(t, 0, result.Code)

	var updated model.Guestbook
	db.First(&updated, entry.ID)
	assert.False(t, updated.IsPinned)
}

func Test_AdminUsers(t *testing.T) {
	_, ts, token := setupPhase34Router(t)
	defer ts.Close()

	result := authedGet(ts.URL+"/api/v1/admin/users", token)
	assert.Equal(t, 0, result.Code)
}

func Test_AdminNavItems_CRUD(t *testing.T) {
	_, ts, token := setupPhase34Router(t)
	defer ts.Close()

	// Create
	result := authedPost(ts.URL+"/api/v1/admin/nav-items", token, map[string]interface{}{
		"title": "推荐", "url": "/recommend",
	})
	assert.Equal(t, 0, result.Code)

	// List
	listResult := authedGet(ts.URL+"/api/v1/admin/nav-items", token)
	assert.Equal(t, 0, listResult.Code)

	// Update
	upResult := authedPut(ts.URL+"/api/v1/admin/nav-items/1", token, map[string]interface{}{
		"title": "热门推荐", "url": "/hot",
	})
	assert.Equal(t, 0, upResult.Code)

	// Delete
	delResult := authedDelete(ts.URL+"/api/v1/admin/nav-items/1", token)
	assert.Equal(t, 0, delResult.Code)
}

func Test_AdminBanners_CRUD(t *testing.T) {
	_, ts, token := setupPhase34Router(t)
	defer ts.Close()

	result := authedPost(ts.URL+"/api/v1/admin/banners", token, map[string]interface{}{
		"image_url": "https://img.com/banner.jpg",
	})
	assert.Equal(t, 0, result.Code)

	listResult := authedGet(ts.URL+"/api/v1/admin/banners", token)
	assert.Equal(t, 0, listResult.Code)
}

// === 金币不会为负 ===

func Test_DeductCoins_不允许负数(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	cat := model.Category{Name: "设计", Slug: "design", IsActive: true}
	db.Create(&cat)
	course := model.Course{Title: "贵课程", Slug: "expensive", CategoryID: cat.ID, Price: 999, Status: "published"}
	db.Create(&course)

	// No recharge, try to buy expensive course
	result := authedPost(ts.URL+"/api/v1/orders", token, map[string]interface{}{
		"type": "course", "target_id": course.ID,
	})
	assert.NotEqual(t, 0, result.Code) // Should fail: insufficient balance

	// Verify balance is still 0
	balResult := authedGet(ts.URL+"/api/v1/coins/balance", token)
	var bal model.CoinBalanceResponse
	json.Unmarshal(balResult.Data, &bal)
	assert.Equal(t, 0, bal.Balance)
}
