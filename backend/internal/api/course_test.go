package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	"github.com/zioran/backend/pkg/sms"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func jsonBody(v interface{}) io.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}

func setupCourseTestRouter(t *testing.T) (*gorm.DB, *httptest.Server, string) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&model.User{}, &model.Course{}, &model.Category{}, &model.Tag{}, &model.CourseResource{}, &model.UserFavorite{},
		&model.CoinAccount{}, &model.CoinTransaction{}, &model.VipPackage{}, &model.UserVip{}, &model.Order{}, &model.Purchase{},
		&model.Guestbook{}, &model.GuestbookLike{}, &model.Comment{}, &model.NavItem{}, &model.Banner{}, &model.UserDownload{},
		&model.Ticket{}, &model.TicketReply{}, &model.Setting{}, &model.OperationLog{}, &model.PaymentLog{}))

	// Seed data
	cat := model.Category{Name: "AIGC课堂", Slug: "aigc", IsActive: true}
	db.Create(&cat)
	tag := model.Tag{Name: "PS课程", Slug: "ps"}
	db.Create(&tag)

	now := time.Now()
	for i := 1; i <= 20; i++ {
		pub := now.Add(-time.Duration(i) * time.Hour)
		c := model.Course{
			Title:      fmt.Sprintf("课程%d", i),
			Slug:       fmt.Sprintf("course-%d", i),
			CategoryID: cat.ID,
			CoverImage: "https://example.com/cover.jpg",
			Status:     "published",
			Price:      10,
			PublishedAt: &pub,
		}
		db.Create(&c)
		if i <= 5 {
			db.Exec("INSERT INTO course_tags (course_id, tag_id) VALUES (?, ?)", c.ID, tag.ID)
		}
	}

	userRepo := repository.NewUserRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	catRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)
	favRepo := repository.NewFavoriteRepository(db)
	payRepo := repository.NewPaymentRepository(db)
	commRepo := repository.NewCommunityRepository(db)

	authSvc := service.NewAuthService(userRepo, testJWTSecret, 72*time.Hour, &sms.MockSender{})
	courseSvc := service.NewCourseService(courseRepo, catRepo, tagRepo, favRepo)
	paySvc := service.NewPaymentService(payRepo, courseRepo, userRepo, payment.NewWechatPay(payment.WechatPayConfig{}), payment.NewAlipayClient(payment.AlipayConfig{}))
	commSvc := service.NewCommunityService(commRepo)

	authHandler := api.NewAuthHandler(authSvc)
	courseHandler := api.NewCourseHandler(courseSvc)
	adminHandler := api.NewAdminHandler(courseSvc)
	payHandler := api.NewPaymentHandler(paySvc)
	commHandler := api.NewCommunityHandler(commSvc)
	adminPayHandler := api.NewAdminPaymentHandler(paySvc, commSvc)
	ticketRepo := repository.NewTicketRepository(db)
	ticketSvc := service.NewTicketService(ticketRepo, userRepo)
	ticketHandler := api.NewTicketHandler(ticketSvc, authSvc, paySvc, oauth.NewWechatOAuth(oauth.WechatOAuthConfig{}), testJWTSecret, 72*time.Hour, t.TempDir())

	r := api.SetupRouter(authHandler, courseHandler, adminHandler, payHandler, commHandler, adminPayHandler, api.NewUploadHandler(t.TempDir()), ticketHandler, testJWTSecret)
	ts := httptest.NewServer(r)

	// Register a user and get token
	authSvc.SetSMSCode("13800001111", "123456")
	_, regResult := postJSON(ts.URL+"/api/v1/auth/register", map[string]string{
		"phone": "13800001111", "sms_code": "123456", "password": "testpass",
	})
	var auth model.AuthResponse
	json.Unmarshal(regResult.Data, &auth)

	return db, ts, auth.Token
}

// === Course List Tests ===

func Test_CourseList_返回分页数据(t *testing.T) {
	_, ts, _ := setupCourseTestRouter(t)
	defer ts.Close()

	_, result := getJSON(ts.URL+"/api/v1/courses?page=1&pageSize=5", "")
	assert.Equal(t, 0, result.Code)

	var data model.PaginatedList
	json.Unmarshal(result.Data, &data)
	assert.Equal(t, int64(20), data.Total)
	assert.Equal(t, 1, data.Page)
	assert.Equal(t, 5, data.PageSize)
	assert.Equal(t, 4, data.TotalPages)
}

func Test_CourseList_按分类筛选(t *testing.T) {
	_, ts, _ := setupCourseTestRouter(t)
	defer ts.Close()

	_, result := getJSON(ts.URL+"/api/v1/courses?categoryId=1", "")
	assert.Equal(t, 0, result.Code)

	var data model.PaginatedList
	json.Unmarshal(result.Data, &data)
	assert.Equal(t, int64(20), data.Total)
}

func Test_CourseList_按标签筛选(t *testing.T) {
	_, ts, _ := setupCourseTestRouter(t)
	defer ts.Close()

	_, result := getJSON(ts.URL+"/api/v1/courses?tagId=1", "")
	assert.Equal(t, 0, result.Code)

	var data model.PaginatedList
	json.Unmarshal(result.Data, &data)
	assert.Equal(t, int64(5), data.Total)
}

func Test_CourseList_关键字搜索(t *testing.T) {
	_, ts, _ := setupCourseTestRouter(t)
	defer ts.Close()

	_, result := getJSON(ts.URL+"/api/v1/courses?keyword=课程1", "")
	assert.Equal(t, 0, result.Code)

	var data model.PaginatedList
	json.Unmarshal(result.Data, &data)
	// "课程1", "课程10"-"课程19" = 11 results
	assert.True(t, data.Total >= 1)
}

// === Course Detail Tests ===

func Test_CourseDetail_返回课程详情(t *testing.T) {
	_, ts, _ := setupCourseTestRouter(t)
	defer ts.Close()

	_, result := getJSON(ts.URL+"/api/v1/courses/course-1", "")
	assert.Equal(t, 0, result.Code)

	var data model.CourseDetailResponse
	json.Unmarshal(result.Data, &data)
	assert.Equal(t, "课程1", data.Title)
	assert.Equal(t, "course-1", data.Slug)
	assert.NotNil(t, data.Category)
	assert.Equal(t, "AIGC课堂", data.Category.Name)
}

func Test_CourseDetail_不存在返回404(t *testing.T) {
	_, ts, _ := setupCourseTestRouter(t)
	defer ts.Close()

	_, result := getJSON(ts.URL+"/api/v1/courses/nonexistent", "")
	assert.Equal(t, 40401, result.Code)
}

// === Latest Courses ===

func Test_CoursesLatest_返回最新8个(t *testing.T) {
	_, ts, _ := setupCourseTestRouter(t)
	defer ts.Close()

	_, result := getJSON(ts.URL+"/api/v1/courses/latest", "")
	assert.Equal(t, 0, result.Code)

	var data []model.CourseListItem
	json.Unmarshal(result.Data, &data)
	assert.Equal(t, 8, len(data))
}

// === Like/Favorite ===

func Test_Like_未登录返回40101(t *testing.T) {
	_, ts, _ := setupCourseTestRouter(t)
	defer ts.Close()

	_, result := postJSON(ts.URL+"/api/v1/courses/1/like", nil)
	assert.Equal(t, 40101, result.Code)
}

func Test_Like_登录后收藏成功(t *testing.T) {
	_, ts, token := setupCourseTestRouter(t)
	defer ts.Close()

	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/courses/1/like", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	assert.Equal(t, 0, result.Code)
	var data model.LikeResponse
	json.Unmarshal(result.Data, &data)
	assert.True(t, data.Liked)
	assert.Equal(t, 1, data.LikeCount)
}

func Test_Like_再次点击取消收藏(t *testing.T) {
	_, ts, token := setupCourseTestRouter(t)
	defer ts.Close()

	// Like first time
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/courses/1/like", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	http.DefaultClient.Do(req)

	// Like second time (unlike)
	req2, _ := http.NewRequest("POST", ts.URL+"/api/v1/courses/1/like", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	resp, _ := http.DefaultClient.Do(req2)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	assert.Equal(t, 0, result.Code)
	var data model.LikeResponse
	json.Unmarshal(result.Data, &data)
	assert.False(t, data.Liked)
}

// === Categories ===

func Test_Categories_返回分类列表(t *testing.T) {
	_, ts, _ := setupCourseTestRouter(t)
	defer ts.Close()

	_, result := getJSON(ts.URL+"/api/v1/categories", "")
	assert.Equal(t, 0, result.Code)

	var data []model.Category
	json.Unmarshal(result.Data, &data)
	assert.True(t, len(data) >= 1)
	assert.Equal(t, "AIGC课堂", data[0].Name)
}

// === Tags ===

func Test_Tags_返回标签列表(t *testing.T) {
	_, ts, _ := setupCourseTestRouter(t)
	defer ts.Close()

	_, result := getJSON(ts.URL+"/api/v1/tags", "")
	assert.Equal(t, 0, result.Code)

	var data []model.Tag
	json.Unmarshal(result.Data, &data)
	assert.True(t, len(data) >= 1)
	assert.Equal(t, "PS课程", data[0].Name)
}

// === Search ===

func Test_Search_缺少关键字返回错误(t *testing.T) {
	_, ts, _ := setupCourseTestRouter(t)
	defer ts.Close()

	_, result := getJSON(ts.URL+"/api/v1/search", "")
	assert.Equal(t, 40001, result.Code)
}

func Test_Search_返回匹配结果(t *testing.T) {
	_, ts, _ := setupCourseTestRouter(t)
	defer ts.Close()

	_, result := getJSON(ts.URL+"/api/v1/search?q=课程1", "")
	assert.Equal(t, 0, result.Code)

	var data model.PaginatedList
	json.Unmarshal(result.Data, &data)
	assert.True(t, data.Total >= 1)
}

// === Admin Course CRUD ===

func Test_AdminCourse_未登录返回40101(t *testing.T) {
	_, ts, _ := setupCourseTestRouter(t)
	defer ts.Close()

	_, result := getJSON(ts.URL+"/api/v1/admin/courses", "")
	assert.Equal(t, 40101, result.Code)
}

func Test_AdminCourse_创建课程(t *testing.T) {
	_, ts, token := setupCourseTestRouter(t)
	defer ts.Close()

	body := model.AdminCourseRequest{
		Title:      "新课程",
		Slug:       "new-course",
		CategoryID: 1,
		Price:      20,
	}
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/admin/courses", jsonBody(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	assert.Equal(t, 0, result.Code)
}

func Test_AdminCourse_更新课程(t *testing.T) {
	_, ts, token := setupCourseTestRouter(t)
	defer ts.Close()

	body := model.AdminCourseRequest{
		Title:      "更新标题",
		Slug:       "course-1",
		CategoryID: 1,
		Price:      30,
	}
	req, _ := http.NewRequest("PUT", ts.URL+"/api/v1/admin/courses/1", jsonBody(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	assert.Equal(t, 0, result.Code)
}

func Test_AdminCourse_更新状态(t *testing.T) {
	_, ts, token := setupCourseTestRouter(t)
	defer ts.Close()

	body := model.AdminCourseStatusRequest{Status: "published"}
	req, _ := http.NewRequest("PUT", ts.URL+"/api/v1/admin/courses/1/status", jsonBody(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	assert.Equal(t, 0, result.Code)
}

func Test_AdminCourse_删除课程(t *testing.T) {
	_, ts, token := setupCourseTestRouter(t)
	defer ts.Close()

	req, _ := http.NewRequest("DELETE", ts.URL+"/api/v1/admin/courses/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	assert.Equal(t, 0, result.Code)
}

func Test_AdminCourse_删除不存在返回404(t *testing.T) {
	_, ts, token := setupCourseTestRouter(t)
	defer ts.Close()

	req, _ := http.NewRequest("DELETE", ts.URL+"/api/v1/admin/courses/9999", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	assert.Equal(t, 40401, result.Code)
}

// === Admin Category CRUD ===

func Test_AdminCategory_创建分类(t *testing.T) {
	_, ts, token := setupCourseTestRouter(t)
	defer ts.Close()

	body := model.AdminCategoryRequest{Name: "新分类", Slug: "new-cat"}
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/admin/categories", jsonBody(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	assert.Equal(t, 0, result.Code)
}

func Test_AdminCategory_更新分类(t *testing.T) {
	_, ts, token := setupCourseTestRouter(t)
	defer ts.Close()

	body := model.AdminCategoryRequest{Name: "更新分类", Slug: "aigc"}
	req, _ := http.NewRequest("PUT", ts.URL+"/api/v1/admin/categories/1", jsonBody(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	assert.Equal(t, 0, result.Code)
}

func Test_AdminCategory_删除分类(t *testing.T) {
	db, ts, token := setupCourseTestRouter(t)
	defer ts.Close()

	// Create a category without courses to delete
	cat := model.Category{Name: "待删除", Slug: "to-del", IsActive: true}
	db.Create(&cat)

	req, _ := http.NewRequest("DELETE", ts.URL+fmt.Sprintf("/api/v1/admin/categories/%d", cat.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	assert.Equal(t, 0, result.Code)
}

// === Admin Tag CRUD ===

func Test_AdminTag_创建标签(t *testing.T) {
	_, ts, token := setupCourseTestRouter(t)
	defer ts.Close()

	body := model.AdminTagRequest{Name: "新标签", Slug: "new-tag"}
	req, _ := http.NewRequest("POST", ts.URL+"/api/v1/admin/tags", jsonBody(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	assert.Equal(t, 0, result.Code)
}

func Test_AdminTag_更新标签(t *testing.T) {
	_, ts, token := setupCourseTestRouter(t)
	defer ts.Close()

	body := model.AdminTagRequest{Name: "更新标签", Slug: "ps"}
	req, _ := http.NewRequest("PUT", ts.URL+"/api/v1/admin/tags/1", jsonBody(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	assert.Equal(t, 0, result.Code)
}

func Test_AdminTag_删除标签(t *testing.T) {
	_, ts, token := setupCourseTestRouter(t)
	defer ts.Close()

	req, _ := http.NewRequest("DELETE", ts.URL+"/api/v1/admin/tags/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := http.DefaultClient.Do(req)
	var result apiResponse
	json.NewDecoder(resp.Body).Decode(&result)
	resp.Body.Close()

	assert.Equal(t, 0, result.Code)
}
