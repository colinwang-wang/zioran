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
		&model.Ticket{}, &model.TicketReply{}, &model.TicketAttachment{}, &model.Setting{},
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

func Test_UserDownloads_ReturnsOrderMeta(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	cat := &model.Category{Name: "编程", Slug: "code", IsActive: true}
	db.Create(cat)
	course := &model.Course{Title: "Go 实战", Slug: "go", CategoryID: cat.ID, Status: "published", Price: 20}
	db.Create(course)
	targetID := int(course.ID)
	order := &model.Order{OrderNo: "ORD_DOWNLOAD_1", UserID: 1, Type: "course", TargetID: &targetID, TargetName: course.Title, Amount: 20, Status: "paid"}
	db.Create(order)
	db.Create(&model.Purchase{UserID: 1, CourseID: course.ID, OrderID: &order.ID})
	db.Create(&model.UserDownload{UserID: 1, CourseID: course.ID})

	result := authedGet(ts.URL+"/api/v1/user/downloads", token)
	assert.Equal(t, 0, result.Code)
	var page struct {
		Items []model.DownloadResponse `json:"items"`
	}
	json.Unmarshal(result.Data, &page)
	if assert.Len(t, page.Items, 1) {
		assert.Equal(t, "ORD_DOWNLOAD_1", page.Items[0].OrderNo)
		assert.Equal(t, 20, page.Items[0].Amount)
	}
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

func Test_AdminUserDetail_ReturnsVipAndStats(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	cat := &model.Category{Name: "编程", Slug: "code", IsActive: true}
	db.Create(cat)
	course := &model.Course{Title: "Go 实战", Slug: "go", CategoryID: cat.ID, Status: "published"}
	db.Create(course)
	db.Create(&model.CoinAccount{UserID: 1, Balance: 88})
	db.Create(&model.VipPackage{ID: 1, Name: "终身VIP", Price: 99, IsActive: true})
	db.Create(&model.UserVip{UserID: 1, PackageID: 1, StartedAt: time.Now(), IsActive: true})
	targetID := int(course.ID)
	order := &model.Order{OrderNo: "ORD_ADMIN_USER_1", UserID: 1, Type: "course", TargetID: &targetID, TargetName: course.Title, Amount: 20, Status: "paid"}
	db.Create(order)
	db.Create(&model.Purchase{UserID: 1, CourseID: course.ID, OrderID: &order.ID})
	db.Create(&model.UserFavorite{UserID: 1, CourseID: course.ID})

	result := authedGet(ts.URL+"/api/v1/admin/users/1", token)
	assert.Equal(t, 0, result.Code)
	var user model.AdminUserResponse
	json.Unmarshal(result.Data, &user)
	assert.True(t, user.IsVip)
	assert.Equal(t, 88, user.Balance)
	assert.Equal(t, int64(1), user.PurchasedCount)
	assert.Equal(t, int64(1), user.FavoriteCount)
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

// === Bug1: 退款逻辑修复 (TDD-RED) ===

func Test_AdminRefund_Coin订单_余额不足时拒绝全额退款(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	// 充值10元（mock自动完成，余额为10）
	authedPost(ts.URL+"/api/v1/coins/recharge", token, map[string]interface{}{
		"amount": 10, "pay_method": "wechat",
	})

	// 消费4金币：创建课程并购买
	cat := model.Category{Name: "test", Slug: "test", IsActive: true}
	db.Create(&cat)
	course := model.Course{Title: "课程A", Slug: "course-a", CategoryID: cat.ID, Price: 4, Status: "published",
		Resources: []model.CourseResource{{Name: "R1", URL: "http://example.com"}}}
	db.Create(&course)
	authedPost(ts.URL+"/api/v1/orders", token, map[string]interface{}{
		"type": "course", "target_id": course.ID,
	})

	// 余额应该是 10 - 4 = 6
	balResult := authedGet(ts.URL+"/api/v1/coins/balance", token)
	var bal model.CoinBalanceResponse
	json.Unmarshal(balResult.Data, &bal)
	assert.Equal(t, 6, bal.Balance)

	// 查找充值订单（第一个订单）
	var rechargeOrder model.Order
	db.Where("type = ? AND status = ?", "coin", "paid").First(&rechargeOrder)

	// 退款充值订单 — 应该只能退6（当前余额），而非10（充值金额）
	// 或者应该拒绝退款因为金额已部分消费
	refundResult := authedPost(fmt.Sprintf("%s/api/v1/admin/orders/%d/refund", ts.URL, rechargeOrder.ID), token, nil)

	// 退款后余额应为0（退了6金币），而非变成负数
	balResult2 := authedGet(ts.URL+"/api/v1/coins/balance", token)
	var bal2 model.CoinBalanceResponse
	json.Unmarshal(balResult2.Data, &bal2)
	assert.Equal(t, 0, bal2.Balance, "退款后余额应为0")
	assert.Equal(t, 0, refundResult.Code, "退款操作应成功")

	// 订单状态应变为 refunded
	var updatedOrder model.Order
	db.First(&updatedOrder, rechargeOrder.ID)
	assert.Equal(t, "refunded", updatedOrder.Status)
}

func Test_AdminRefund_Coin订单_按实际到账金币扣减(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	db.Create(&model.Setting{Key: "coinRechargeRatio", Value: "10"})

	authedPost(ts.URL+"/api/v1/coins/recharge", token, map[string]interface{}{
		"amount": 10, "pay_method": "wechat",
	})

	var rechargeOrder model.Order
	db.Where("type = ? AND status = ?", "coin", "paid").First(&rechargeOrder)

	refundResult := authedPost(fmt.Sprintf("%s/api/v1/admin/orders/%d/refund", ts.URL, rechargeOrder.ID), token, nil)
	assert.Equal(t, 0, refundResult.Code)

	balResult := authedGet(ts.URL+"/api/v1/coins/balance", token)
	var bal model.CoinBalanceResponse
	json.Unmarshal(balResult.Data, &bal)
	assert.Equal(t, 0, bal.Balance, "充值10元按10倍到账100金币，退款应扣减100金币")
}

func Test_AdminRefund_Course订单_应退还金币到用户账户(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	// 充值20
	authedPost(ts.URL+"/api/v1/coins/recharge", token, map[string]interface{}{
		"amount": 20, "pay_method": "wechat",
	})

	// 购买课程花费10
	cat := model.Category{Name: "test", Slug: "test-cat", IsActive: true}
	db.Create(&cat)
	course := model.Course{Title: "课程B", Slug: "course-b", CategoryID: cat.ID, Price: 10, Status: "published",
		Resources: []model.CourseResource{{Name: "R1", URL: "http://example.com"}}}
	db.Create(&course)
	authedPost(ts.URL+"/api/v1/orders", token, map[string]interface{}{
		"type": "course", "target_id": course.ID,
	})

	// 余额 = 20 - 10 = 10
	balResult := authedGet(ts.URL+"/api/v1/coins/balance", token)
	var bal model.CoinBalanceResponse
	json.Unmarshal(balResult.Data, &bal)
	assert.Equal(t, 10, bal.Balance)

	// 找到课程购买订单
	var courseOrder model.Order
	db.Where("type = ? AND status = ?", "course", "paid").First(&courseOrder)

	// 退款课程订单 — 应退还10金币
	refundResult := authedPost(fmt.Sprintf("%s/api/v1/admin/orders/%d/refund", ts.URL, courseOrder.ID), token, nil)
	assert.Equal(t, 0, refundResult.Code, "课程订单退款应成功")

	// 退款后余额 = 10 + 10 = 20
	balResult2 := authedGet(ts.URL+"/api/v1/coins/balance", token)
	var bal2 model.CoinBalanceResponse
	json.Unmarshal(balResult2.Data, &bal2)
	assert.Equal(t, 20, bal2.Balance, "课程退款应退还金币")
}

func Test_AdminRefund_Course订单_重复退款不重复退金币(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	authedPost(ts.URL+"/api/v1/coins/recharge", token, map[string]interface{}{
		"amount": 20, "pay_method": "wechat",
	})

	cat := model.Category{Name: "test", Slug: "repeat-refund", IsActive: true}
	db.Create(&cat)
	course := model.Course{Title: "重复退款课程", Slug: "repeat-refund-course", CategoryID: cat.ID, Price: 10, Status: "published"}
	db.Create(&course)
	authedPost(ts.URL+"/api/v1/orders", token, map[string]interface{}{
		"type": "course", "target_id": course.ID,
	})

	var courseOrder model.Order
	db.Where("type = ? AND status = ?", "course", "paid").First(&courseOrder)

	firstRefund := authedPost(fmt.Sprintf("%s/api/v1/admin/orders/%d/refund", ts.URL, courseOrder.ID), token, nil)
	secondRefund := authedPost(fmt.Sprintf("%s/api/v1/admin/orders/%d/refund", ts.URL, courseOrder.ID), token, nil)
	assert.Equal(t, 0, firstRefund.Code)
	assert.NotEqual(t, 0, secondRefund.Code)

	balResult := authedGet(ts.URL+"/api/v1/coins/balance", token)
	var bal model.CoinBalanceResponse
	json.Unmarshal(balResult.Data, &bal)
	assert.Equal(t, 20, bal.Balance, "重复退款后余额不应超过原始充值余额")
}

func Test_AdminRefund_VIP订单_应退还金币(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	db.Create(&model.VipPackage{Name: "终身VIP", Price: 99, IsActive: true})

	// 充值100
	authedPost(ts.URL+"/api/v1/coins/recharge", token, map[string]interface{}{
		"amount": 100, "pay_method": "wechat",
	})

	// 购买VIP花费99
	authedPost(ts.URL+"/api/v1/vip/purchase", token, map[string]interface{}{"package_id": 1})

	// 余额 = 100 - 99 = 1
	balResult := authedGet(ts.URL+"/api/v1/coins/balance", token)
	var bal model.CoinBalanceResponse
	json.Unmarshal(balResult.Data, &bal)
	assert.Equal(t, 1, bal.Balance)

	// 找到VIP购买订单
	var vipOrder model.Order
	db.Where("type = ? AND status = ?", "vip", "paid").First(&vipOrder)

	// 退款VIP订单 — 应退还99金币
	refundResult := authedPost(fmt.Sprintf("%s/api/v1/admin/orders/%d/refund", ts.URL, vipOrder.ID), token, nil)
	assert.Equal(t, 0, refundResult.Code, "VIP订单退款应成功")

	// 退款后余额 = 1 + 99 = 100
	balResult2 := authedGet(ts.URL+"/api/v1/coins/balance", token)
	var bal2 model.CoinBalanceResponse
	json.Unmarshal(balResult2.Data, &bal2)
	assert.Equal(t, 100, bal2.Balance, "VIP退款应退还金币")
}

// === Bug2: 用户管理VIP筛选 (TDD-RED) ===

func Test_AdminUserList_VIP筛选_只返回VIP用户(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	// 创建普通用户和VIP用户
	normalUser := &model.User{Username: "normaluser", Phone: "13900000002", PasswordHash: "$2a$10$dummy", Role: "user", Status: "active"}
	db.Create(normalUser)
	vipUser := &model.User{Username: "vipuser", Phone: "13900000003", PasswordHash: "$2a$10$dummy", Role: "user", Status: "active"}
	db.Create(vipUser)

	// 给vipUser设置VIP
	db.Create(&model.UserVip{UserID: vipUser.ID, PackageID: 1, StartedAt: time.Now(), IsActive: true})

	// 用vipFilter=vip筛选
	result := authedGet(ts.URL+"/api/v1/admin/users?vipFilter=vip", token)
	assert.Equal(t, 0, result.Code)

	var page struct {
		Items []model.AdminUserResponse `json:"items"`
		Total int64                     `json:"total"`
	}
	json.Unmarshal(result.Data, &page)

	// 应该只包含VIP用户（vipUser），不包含normalUser和默认testuser
	for _, u := range page.Items {
		assert.True(t, u.IsVip, "vipFilter=vip 应只返回VIP用户，但返回了: "+u.Username)
	}
	assert.True(t, page.Total >= 1, "至少应有1个VIP用户")
}

func Test_AdminUserList_Normal筛选_只返回非VIP用户(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	normalUser := &model.User{Username: "normaluser", Phone: "13900000002", PasswordHash: "$2a$10$dummy", Role: "user", Status: "active"}
	db.Create(normalUser)
	vipUser := &model.User{Username: "vipuser", Phone: "13900000003", PasswordHash: "$2a$10$dummy", Role: "user", Status: "active"}
	db.Create(vipUser)
	db.Create(&model.UserVip{UserID: vipUser.ID, PackageID: 1, StartedAt: time.Now(), IsActive: true})

	// 用vipFilter=normal筛选
	result := authedGet(ts.URL+"/api/v1/admin/users?vipFilter=normal", token)
	assert.Equal(t, 0, result.Code)

	var page struct {
		Items []model.AdminUserResponse `json:"items"`
		Total int64                     `json:"total"`
	}
	json.Unmarshal(result.Data, &page)

	// 不应包含VIP用户
	for _, u := range page.Items {
		assert.False(t, u.IsVip, "vipFilter=normal 应只返回非VIP用户，但返回了VIP: "+u.Username)
	}
}

// === Bug3: 订单管理搜索功能 (TDD-RED) ===

func Test_AdminOrders_Keyword搜索_按订单号(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	user := &model.User{Username: "buyer", Phone: "13900000010", PasswordHash: "$2a$10$dummy", Role: "user", Status: "active"}
	db.Create(user)
	now := time.Now()
	db.Create(&model.Order{OrderNo: "ORD_SEARCH_001", UserID: user.ID, Type: "coin", TargetName: "充值100金币", Amount: 100, Status: "paid", CreatedAt: now})
	db.Create(&model.Order{OrderNo: "ORD_SEARCH_002", UserID: user.ID, Type: "course", TargetName: "Go实战课程", Amount: 20, Status: "paid", CreatedAt: now})
	db.Create(&model.Order{OrderNo: "ORD_OTHER_003", UserID: user.ID, Type: "vip", TargetName: "终身VIP", Amount: 99, Status: "paid", CreatedAt: now})

	// 按订单号搜索
	result := authedGet(ts.URL+"/api/v1/admin/orders?keyword=SEARCH_001", token)
	assert.Equal(t, 0, result.Code)
	var page struct {
		Items []model.OrderResponse `json:"items"`
		Total int64                 `json:"total"`
	}
	json.Unmarshal(result.Data, &page)
	assert.Equal(t, int64(1), page.Total, "按订单号搜索应只匹配1条")
	if assert.Len(t, page.Items, 1) {
		assert.Equal(t, "ORD_SEARCH_001", page.Items[0].OrderNo)
	}
}

func Test_AdminOrders_Keyword搜索_按商品名称(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	user := &model.User{Username: "buyer2", Phone: "13900000011", PasswordHash: "$2a$10$dummy", Role: "user", Status: "active"}
	db.Create(user)
	now := time.Now()
	db.Create(&model.Order{OrderNo: "ORD_A1", UserID: user.ID, Type: "course", TargetName: "Go实战课程", Amount: 20, Status: "paid", CreatedAt: now})
	db.Create(&model.Order{OrderNo: "ORD_A2", UserID: user.ID, Type: "coin", TargetName: "充值50金币", Amount: 50, Status: "paid", CreatedAt: now})

	// 按商品名称搜索
	result := authedGet(ts.URL+"/api/v1/admin/orders?keyword=Go实战", token)
	assert.Equal(t, 0, result.Code)
	var page struct {
		Items []model.OrderResponse `json:"items"`
		Total int64                 `json:"total"`
	}
	json.Unmarshal(result.Data, &page)
	assert.Equal(t, int64(1), page.Total, "按商品名搜索应只匹配1条")
	if assert.Len(t, page.Items, 1) {
		assert.Equal(t, "Go实战课程", page.Items[0].TargetName)
	}
}

func Test_AdminOrders_Keyword搜索_按用户名(t *testing.T) {
	db, ts, token := setupPhase34Router(t)
	defer ts.Close()

	user1 := &model.User{Username: "zhangsan", Phone: "13900000020", PasswordHash: "$2a$10$dummy", Role: "user", Status: "active"}
	db.Create(user1)
	user2 := &model.User{Username: "lisi", Phone: "13900000021", PasswordHash: "$2a$10$dummy", Role: "user", Status: "active"}
	db.Create(user2)
	now := time.Now()
	db.Create(&model.Order{OrderNo: "ORD_U1", UserID: user1.ID, Type: "coin", TargetName: "充值", Amount: 10, Status: "paid", CreatedAt: now})
	db.Create(&model.Order{OrderNo: "ORD_U2", UserID: user2.ID, Type: "coin", TargetName: "充值", Amount: 20, Status: "paid", CreatedAt: now})

	// 按用户名搜索
	result := authedGet(ts.URL+"/api/v1/admin/orders?keyword=zhangsan", token)
	assert.Equal(t, 0, result.Code)
	var page struct {
		Items []struct {
			OrderNo  string `json:"order_no"`
			UserName string `json:"user_name"`
		} `json:"items"`
		Total int64 `json:"total"`
	}
	json.Unmarshal(result.Data, &page)
	assert.Equal(t, int64(1), page.Total, "按用户名搜索应只匹配该用户的订单")
	if assert.Len(t, page.Items, 1) {
		assert.Equal(t, "zhangsan", page.Items[0].UserName)
	}
}
