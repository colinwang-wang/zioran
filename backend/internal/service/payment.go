package service

import (
	"context"
	"fmt"
	"time"

	"github.com/zioran/backend/internal/model"
	"github.com/zioran/backend/internal/repository"
	"github.com/zioran/backend/pkg/errcode"
	"golang.org/x/crypto/bcrypt"
)

type PaymentService struct {
	payRepo    *repository.PaymentRepository
	courseRepo *repository.CourseRepository
	userRepo   *repository.UserRepository
}

func NewPaymentService(payRepo *repository.PaymentRepository, courseRepo *repository.CourseRepository, userRepo *repository.UserRepository) *PaymentService {
	return &PaymentService{payRepo: payRepo, courseRepo: courseRepo, userRepo: userRepo}
}

// Coins

func (s *PaymentService) GetBalance(ctx context.Context, userID int64) (*model.CoinBalanceResponse, error) {
	acc, err := s.payRepo.GetOrCreateAccount(ctx, userID)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	return &model.CoinBalanceResponse{Balance: acc.Balance, TotalEarned: acc.TotalEarned, TotalSpent: acc.TotalSpent}, nil
}

func (s *PaymentService) GetTransactions(ctx context.Context, userID int64, page, pageSize int) (*model.PaginatedList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	txs, total, err := s.payRepo.CoinTransactions(ctx, userID, page, pageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: txs, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}, nil
}

func (s *PaymentService) Recharge(ctx context.Context, userID int64, req *model.RechargeRequest) (*model.RechargeResponse, error) {
	// Create coin order
	order := &model.Order{
		OrderNo:    generateOrderNo(),
		UserID:     userID,
		Type:       "coin",
		TargetName: fmt.Sprintf("充值%d金币", req.Amount),
		Amount:     req.Amount,
		PayMethod:  req.PayMethod,
		Status:     "pending",
	}
	if err := s.payRepo.CreateOrder(ctx, order); err != nil {
		return nil, errcode.ErrInternal
	}

	// MOCK: simulate payment callback (real env should wait for payment gateway callback)
	s.payRepo.UpdateOrderStatus(ctx, order.ID, "paid")
	s.payRepo.Recharge(ctx, userID, req.Amount, order.ID)

	return &model.RechargeResponse{
		OrderID: order.ID,
		OrderNo: order.OrderNo,
		PayURL:  "mock://pay",
	}, nil
}

// VIP

func (s *PaymentService) VipPackages(ctx context.Context) ([]model.VipPackage, error) {
	return s.payRepo.VipPackages(ctx)
}

func (s *PaymentService) VipStatus(ctx context.Context, userID int64) (*model.VipStatusResponse, error) {
	vip, err := s.payRepo.GetVipStatus(ctx, userID)
	if err != nil {
		return &model.VipStatusResponse{IsVip: false}, nil
	}
	pkg, _ := s.payRepo.GetVipPackage(ctx, vip.PackageID)
	name := ""
	if pkg != nil {
		name = pkg.Name
	}
	return &model.VipStatusResponse{IsVip: true, ExpiresAt: vip.ExpiresAt, Package: name}, nil
}

func (s *PaymentService) PurchaseVip(ctx context.Context, userID int64, req *model.VipPurchaseRequest) (*model.OrderResponse, error) {
	// Check if already VIP
	if s.payRepo.IsVip(ctx, userID) {
		return nil, errcode.New(40001, "已是VIP会员")
	}

	pkg, err := s.payRepo.GetVipPackage(ctx, req.PackageID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}

	// Create order
	order := &model.Order{
		OrderNo:    generateOrderNo(),
		UserID:     userID,
		Type:       "vip",
		TargetID:   &pkg.ID,
		TargetName: pkg.Name,
		Amount:     pkg.Price,
		PayMethod:  "coin",
		Status:     "pending",
	}
	if err := s.payRepo.CreateOrder(ctx, order); err != nil {
		return nil, errcode.ErrInternal
	}

	// Deduct coins
	if err := s.payRepo.DeductCoins(ctx, userID, pkg.Price, "vip", "购买"+pkg.Name, &order.ID); err != nil {
		return nil, errcode.New(40001, "金币余额不足")
	}

	// Activate VIP
	var expiresAt *time.Time
	if pkg.DurationDays != nil {
		t := time.Now().Add(time.Duration(*pkg.DurationDays) * 24 * time.Hour)
		expiresAt = &t
	}
	vip := &model.UserVip{
		UserID:    userID,
		PackageID: pkg.ID,
		StartedAt: time.Now(),
		ExpiresAt: expiresAt,
		IsActive:  true,
	}
	if err := s.payRepo.CreateVip(ctx, vip); err != nil {
		return nil, errcode.ErrInternal
	}

	// Update order status
	s.payRepo.UpdateOrderStatus(ctx, order.ID, "paid")
	order.Status = "paid"
	now := time.Now()
	order.PaidAt = &now

	return &model.OrderResponse{
		ID: order.ID, OrderNo: order.OrderNo, Type: order.Type,
		TargetName: order.TargetName, Amount: order.Amount,
		Status: order.Status, CreatedAt: order.CreatedAt, PaidAt: order.PaidAt,
	}, nil
}

// Orders

func (s *PaymentService) GetOrder(ctx context.Context, userID, orderID int64) (*model.OrderResponse, error) {
	order, err := s.payRepo.GetOrder(ctx, orderID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	if order.UserID != userID {
		return nil, errcode.ErrForbidden
	}
	return &model.OrderResponse{
		ID: order.ID, OrderNo: order.OrderNo, Type: order.Type,
		TargetName: order.TargetName, Amount: order.Amount,
		Status: order.Status, CreatedAt: order.CreatedAt, PaidAt: order.PaidAt,
	}, nil
}

func (s *PaymentService) UserOrders(ctx context.Context, userID int64, page, pageSize int) (*model.PaginatedList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	orders, total, err := s.payRepo.UserOrders(ctx, userID, page, pageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	items := make([]model.OrderResponse, len(orders))
	for i, o := range orders {
		items[i] = model.OrderResponse{
			ID: o.ID, OrderNo: o.OrderNo, Type: o.Type,
			TargetName: o.TargetName, Amount: o.Amount,
			Status: o.Status, CreatedAt: o.CreatedAt, PaidAt: o.PaidAt,
		}
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: items, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}, nil
}

// Purchase Course

func (s *PaymentService) PurchaseCourse(ctx context.Context, userID int64, courseID int64) (*model.OrderResponse, error) {
	// Check already purchased
	if s.payRepo.HasPurchased(ctx, userID, courseID) {
		return nil, errcode.New(40001, "已购买该课程")
	}

	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}

	// VIP can download for free
	if s.payRepo.IsVip(ctx, userID) {
		// Create free order + purchase record
		order := &model.Order{
			OrderNo: generateOrderNo(), UserID: userID, Type: "course",
			TargetID: intPtr(int(courseID)), TargetName: course.Title,
			Amount: 0, PayMethod: "vip", Status: "paid",
		}
		now := time.Now()
		order.PaidAt = &now
		s.payRepo.CreateOrder(ctx, order)
		s.payRepo.CreatePurchase(ctx, &model.Purchase{UserID: userID, CourseID: courseID, OrderID: &order.ID})
		return &model.OrderResponse{
			ID: order.ID, OrderNo: order.OrderNo, Type: order.Type,
			TargetName: order.TargetName, Amount: 0, Status: "paid",
			CreatedAt: order.CreatedAt, PaidAt: order.PaidAt,
		}, nil
	}

	// Non-VIP: charge coins
	price := course.Price
	order := &model.Order{
		OrderNo: generateOrderNo(), UserID: userID, Type: "course",
		TargetID: intPtr(int(courseID)), TargetName: course.Title,
		Amount: price, PayMethod: "coin", Status: "pending",
	}
	if err := s.payRepo.CreateOrder(ctx, order); err != nil {
		return nil, errcode.ErrInternal
	}

	if err := s.payRepo.DeductCoins(ctx, userID, price, "purchase", "购买课程: "+course.Title, &order.ID); err != nil {
		return nil, errcode.New(40001, "金币余额不足")
	}

	s.payRepo.CreatePurchase(ctx, &model.Purchase{UserID: userID, CourseID: courseID, OrderID: &order.ID})
	s.payRepo.UpdateOrderStatus(ctx, order.ID, "paid")
	order.Status = "paid"
	now := time.Now()
	order.PaidAt = &now

	return &model.OrderResponse{
		ID: order.ID, OrderNo: order.OrderNo, Type: order.Type,
		TargetName: order.TargetName, Amount: order.Amount,
		Status: "paid", CreatedAt: order.CreatedAt, PaidAt: order.PaidAt,
	}, nil
}

// Download

func (s *PaymentService) Download(ctx context.Context, userID, courseID int64) (*model.CourseDownloadResponse, error) {
	// Check access: purchased or VIP
	hasPurchased := s.payRepo.HasPurchased(ctx, userID, courseID)
	isVip := s.payRepo.IsVip(ctx, userID)

	if !hasPurchased && !isVip {
		return nil, errcode.New(40301, "未购买该课程")
	}

	course, err := s.courseRepo.FindByID(ctx, courseID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}

	// Record download
	s.payRepo.RecordDownload(ctx, userID, courseID)

	resources := make([]model.ResourceItem, len(course.Resources))
	for i, r := range course.Resources {
		resources[i] = model.ResourceItem{ID: r.ID, Name: r.Name, URL: r.URL, Password: r.Password}
	}
	return &model.CourseDownloadResponse{Resources: resources}, nil
}

// User center

func (s *PaymentService) UserDownloads(ctx context.Context, userID int64, page, pageSize int) (*model.PaginatedList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	downloads, total, err := s.payRepo.UserDownloads(ctx, userID, page, pageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	items := make([]model.DownloadResponse, len(downloads))
	for i, d := range downloads {
		items[i] = model.DownloadResponse{ID: d.ID, CourseID: d.CourseID, CreatedAt: d.CreatedAt}
		if d.Course != nil {
			items[i].Title = d.Course.Title
			items[i].Cover = d.Course.CoverImage
		}
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: items, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}, nil
}

func (s *PaymentService) UserFavorites(ctx context.Context, userID int64, page, pageSize int) (*model.PaginatedList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	favs, total, err := s.payRepo.UserFavorites(ctx, userID, page, pageSize)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	courseIDs := make([]int64, len(favs))
	for i, f := range favs {
		courseIDs[i] = f.CourseID
	}
	courses, _ := s.payRepo.GetCoursesByIDs(ctx, courseIDs)
	courseMap := make(map[int64]*model.Course)
	for i := range courses {
		courseMap[courses[i].ID] = &courses[i]
	}
	items := make([]model.FavoriteResponse, len(favs))
	for i, f := range favs {
		items[i] = model.FavoriteResponse{CourseID: f.CourseID, CreatedAt: f.CreatedAt}
		if c, ok := courseMap[f.CourseID]; ok {
			items[i].Title = c.Title
			items[i].Cover = c.CoverImage
			items[i].Slug = c.Slug
		}
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: items, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}, nil
}

func (s *PaymentService) ChangePassword(ctx context.Context, userID int64, req *model.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errcode.ErrNotFound
	}
	if err := comparePassword(user.PasswordHash, req.OldPassword); err != nil {
		return errcode.New(40001, "原密码错误")
	}
	hash, err := hashPassword(req.NewPassword)
	if err != nil {
		return errcode.ErrInternal
	}
	user.PasswordHash = hash
	return s.userRepo.UpdatePassword(ctx, userID, hash)
}

// Admin

func (s *PaymentService) AdminOrders(ctx context.Context, page, pageSize int, status string) (*model.PaginatedList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	orders, total, err := s.payRepo.AdminOrders(ctx, page, pageSize, status)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	items := make([]model.OrderResponse, len(orders))
	for i, o := range orders {
		items[i] = model.OrderResponse{
			ID: o.ID, OrderNo: o.OrderNo, Type: o.Type,
			TargetName: o.TargetName, Amount: o.Amount,
			Status: o.Status, CreatedAt: o.CreatedAt, PaidAt: o.PaidAt,
		}
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: items, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}, nil
}

func (s *PaymentService) AdminRefund(ctx context.Context, orderID int64) error {
	order, err := s.payRepo.GetOrder(ctx, orderID)
	if err != nil {
		return errcode.ErrNotFound
	}
	if order.Status != "paid" {
		return errcode.New(40001, "订单状态不允许退款")
	}
	// Refund coins
	s.payRepo.Recharge(ctx, order.UserID, order.Amount, order.ID)
	return s.payRepo.UpdateOrderStatus(ctx, orderID, "refunded")
}

func (s *PaymentService) AdminUsers(ctx context.Context, page, pageSize int, keyword string) (*model.PaginatedList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	users, total, err := s.payRepo.AdminUsers(ctx, page, pageSize, keyword)
	if err != nil {
		return nil, errcode.ErrInternal
	}
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return &model.PaginatedList{Items: users, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}, nil
}

func (s *PaymentService) AdminUpdateUserStatus(ctx context.Context, userID int64, status string) error {
	return s.payRepo.UpdateUserStatus(ctx, userID, status)
}

func (s *PaymentService) AdminRecharge(ctx context.Context, userID int64, req *model.AdminRechargeRequest) error {
	order := &model.Order{
		OrderNo: generateOrderNo(), UserID: userID, Type: "coin",
		TargetName: "管理员充值", Amount: req.Amount, PayMethod: "admin", Status: "paid",
	}
	now := time.Now()
	order.PaidAt = &now
	if err := s.payRepo.CreateOrder(ctx, order); err != nil {
		return errcode.ErrInternal
	}
	return s.payRepo.Recharge(ctx, userID, req.Amount, order.ID)
}

func (s *PaymentService) DashboardStats(ctx context.Context) (*model.DashboardStats, error) {
	return s.payRepo.DashboardStats(ctx)
}

func (s *PaymentService) AdminGetUser(ctx context.Context, userID int64) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errcode.ErrNotFound
	}
	return user, nil
}

func (s *PaymentService) DashboardCharts(ctx context.Context, period string) *model.DashboardChartsResponse {
	// MOCK: return sample chart data
	days := 7
	if period == "month" {
		days = 30
	}
	users := make([]model.ChartPoint, days)
	orders := make([]model.ChartPoint, days)
	for i := 0; i < days; i++ {
		date := time.Now().AddDate(0, 0, -days+1+i).Format("2006-01-02")
		users[i] = model.ChartPoint{Date: date, Value: int64(10 + i*2)}
		orders[i] = model.ChartPoint{Date: date, Value: int64(5 + i)}
	}
	return &model.DashboardChartsResponse{Users: users, Orders: orders}
}

// Favorites

func (s *PaymentService) AddFavorite(ctx context.Context, userID, courseID int64) error {
	return s.payRepo.AddFavorite(ctx, userID, courseID)
}

func (s *PaymentService) RemoveFavorite(ctx context.Context, userID, courseID int64) error {
	return s.payRepo.RemoveFavorite(ctx, userID, courseID)
}

// helpers

func generateOrderNo() string {
	return fmt.Sprintf("ORD%d", time.Now().UnixNano())
}

func intPtr(v int) *int { return &v }

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func comparePassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
