package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zioran/backend/internal/middleware"
	"gorm.io/gorm"
)

func SetupRouter(
	authHandler *AuthHandler,
	courseHandler *CourseHandler,
	adminHandler *AdminHandler,
	payHandler *PaymentHandler,
	commHandler *CommunityHandler,
	adminPayHandler *AdminPaymentHandler,
	uploadHandler *UploadHandler,
	ticketHandler *TicketHandler,
	jwtSecret string,
	db ...*gorm.DB,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())
	r.Static("/uploads", uploadHandler.uploadDir)

	v1 := r.Group("/api/v1")
	{
		// Auth
		auth := v1.Group("/auth")
		{
			auth.POST("/captcha", authHandler.Captcha)
			auth.POST("/email/send", authHandler.SendEmail)
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/forgot-password", ticketHandler.ForgotPassword)
			// OAuth
			auth.GET("/oauth/wechat", ticketHandler.OAuthWechat)
			auth.GET("/oauth/wechat/callback", ticketHandler.OAuthWechatCallback)
			auth.POST("/oauth/wechat/callback", ticketHandler.OAuthWechatCallback)
		}

		// Public course routes
		v1.GET("/courses/latest", courseHandler.Latest)
		v1.GET("/courses", courseHandler.List)
		v1.GET("/courses/:slug", courseHandler.Detail)
		v1.GET("/categories", courseHandler.Categories)
		v1.GET("/tags", courseHandler.Tags)
		v1.GET("/search", courseHandler.Search)

		// Public home config
		v1.GET("/home/nav-items", commHandler.NavItems)
		v1.GET("/home/banners", commHandler.Banners)
		v1.GET("/home/config", ticketHandler.HomeConfig)

		// Public VIP packages
		v1.GET("/vip/packages", payHandler.VipPackages)

		// Public recharge config
		v1.GET("/coins/recharge-config", payHandler.RechargeConfig)

		// Public guestbook (read)
		v1.GET("/guestbook", commHandler.GuestbookList)

		// Public comments (read)
		v1.GET("/comments", commHandler.CommentList)

		// Payment notify (MOCK: 待接入真实服务)
		v1.POST("/pay/notify/wechat", ticketHandler.WechatNotify)
		v1.POST("/pay/notify/alipay", ticketHandler.AlipayNotify)

		// Authed routes
		authed := v1.Group("")
		authed.Use(middleware.JWTAuth(jwtSecret))
		{
			// Token refresh
			authed.POST("/auth/refresh", ticketHandler.RefreshToken)

			// Upload
			authed.POST("/upload/image", uploadHandler.ImageUpload)
			authed.POST("/upload/images", ticketHandler.BatchImageUpload)

			// User profile
			authed.GET("/user/profile", authHandler.Profile)
			authed.PUT("/user/profile", authHandler.UpdateProfile)
			authed.PUT("/user/password", payHandler.ChangePassword)
			authed.GET("/user/orders", payHandler.UserOrders)
			authed.GET("/user/orders/:id", payHandler.GetOrder)
			authed.GET("/user/downloads", payHandler.UserDownloads)
			authed.GET("/user/favorites", payHandler.UserFavorites)
			authed.POST("/user/favorites", payHandler.AddFavorite)
			authed.DELETE("/user/favorites/:courseId", payHandler.RemoveFavorite)
			authed.GET("/user/comments", commHandler.UserCommentList)

			// Courses (authed)
			authed.POST("/courses/:id/like", courseHandler.Like)
			authed.POST("/courses/:id/download", payHandler.Download)

			// Coins
			authed.GET("/coins/balance", payHandler.CoinBalance)
			authed.GET("/coins/transactions", payHandler.CoinTransactions)
			authed.POST("/coins/recharge", payHandler.Recharge)

			// VIP
			authed.GET("/vip/status", payHandler.VipStatus)
			authed.POST("/vip/purchase", payHandler.VipPurchase)

			// Orders
			authed.POST("/orders", payHandler.CreateOrder)
			authed.GET("/orders/:id", payHandler.GetOrder)
			authed.POST("/orders/:id/cancel", ticketHandler.CancelOrder)

			// Guestbook (write)
			authed.POST("/guestbook", commHandler.GuestbookCreate)
			authed.POST("/guestbook/:id/like", commHandler.GuestbookLike)
			authed.DELETE("/guestbook/:id", commHandler.GuestbookDelete)

			// Comments (write)
			authed.POST("/comments", commHandler.CommentCreate)
			authed.DELETE("/comments/:id", commHandler.CommentDelete)

			// Tickets (user)
			authed.GET("/tickets", ticketHandler.TicketList)
			authed.POST("/tickets", ticketHandler.TicketCreate)
			authed.GET("/tickets/:id", ticketHandler.TicketDetail)
			authed.POST("/tickets/:id/reply", ticketHandler.TicketReply)
		}

		// Admin login (no JWT required)
		v1.POST("/admin/login", authHandler.AdminLogin)

		// Admin routes
		admin := v1.Group("/admin")
		admin.Use(middleware.JWTAuth(jwtSecret))
		if len(db) > 0 && db[0] != nil {
			admin.Use(middleware.AdminRequired(db[0]))
		}
		{
			// === 所有管理角色可访问（仪表盘） ===
			admin.GET("/dashboard/stats", adminPayHandler.DashboardStats)
			admin.GET("/dashboard/charts", adminPayHandler.DashboardCharts)

			// === 客服+运营+管理员+超管 可访问 ===
			// Tickets
			admin.GET("/tickets", ticketHandler.AdminTicketList)
			admin.GET("/tickets/:id", ticketHandler.AdminTicketDetail)
			admin.PUT("/tickets/:id/status", ticketHandler.AdminTicketUpdateStatus)
			admin.POST("/tickets/:id/reply", ticketHandler.AdminTicketReply)
			// Guestbook
			admin.GET("/guestbook", adminPayHandler.GuestbookList)
			admin.PUT("/guestbook/:id/status", adminPayHandler.GuestbookUpdateStatus)
			admin.PUT("/guestbook/:id/pin", adminPayHandler.GuestbookPin)
			admin.DELETE("/guestbook/:id", adminPayHandler.GuestbookDelete)
			// Comments
			admin.GET("/comments", adminPayHandler.CommentList)
			admin.PUT("/comments/:id/status", adminPayHandler.CommentUpdateStatus)
			admin.DELETE("/comments/:id", adminPayHandler.CommentDelete)
			admin.POST("/comments/:id/reply", ticketHandler.AdminCommentReply)

			// === 运营+管理员+超管 可访问 ===
			ops := admin.Group("")
			ops.Use(middleware.PermissionRequired("courses", "orders", "home_config"))
			{
				// Course management
				ops.GET("/courses", adminHandler.CourseList)
				ops.POST("/courses", adminHandler.CourseCreate)
				ops.PUT("/courses/:id", adminHandler.CourseUpdate)
				ops.DELETE("/courses/:id", adminHandler.CourseDelete)
				ops.PUT("/courses/:id/status", adminHandler.CourseUpdateStatus)
				ops.POST("/courses/batch", adminHandler.CourseBatch)
				// Category management
				ops.GET("/categories", adminHandler.CategoryList)
				ops.POST("/categories", adminHandler.CategoryCreate)
				ops.PUT("/categories/:id", adminHandler.CategoryUpdate)
				ops.DELETE("/categories/:id", adminHandler.CategoryDelete)
				ops.PUT("/categories/:id/status", adminHandler.CategoryUpdateStatus)
				// Tag management
				ops.GET("/tags", adminHandler.TagList)
				ops.POST("/tags", adminHandler.TagCreate)
				ops.PUT("/tags/:id", adminHandler.TagUpdate)
				ops.DELETE("/tags/:id", adminHandler.TagDelete)
				// Order management
				ops.GET("/orders", adminPayHandler.OrderList)
				ops.GET("/orders/:id", adminPayHandler.OrderDetail)
				ops.POST("/orders/:id/refund", adminPayHandler.OrderRefund)
				// VIP package management
				ops.GET("/vip/packages", adminPayHandler.VipPackageList)
				ops.PUT("/vip/packages/:id", adminPayHandler.VipPackageUpdate)
				// Nav items management
				ops.GET("/nav-items", adminPayHandler.NavItemList)
				ops.POST("/nav-items", adminPayHandler.NavItemCreate)
				ops.PUT("/nav-items/:id", adminPayHandler.NavItemUpdate)
				ops.DELETE("/nav-items/:id", adminPayHandler.NavItemDelete)
				// Banner management
				ops.GET("/banners", adminPayHandler.BannerList)
				ops.POST("/banners", adminPayHandler.BannerCreate)
				ops.PUT("/banners/:id", adminPayHandler.BannerUpdate)
				ops.DELETE("/banners/:id", adminPayHandler.BannerDelete)
				// Data
				ops.GET("/finance/summary", ticketHandler.FinanceSummary)
				ops.GET("/finance/withdrawals", ticketHandler.FinanceWithdrawals)
				ops.GET("/logs/operations", ticketHandler.OperationLogs)
				ops.GET("/logs/payments", ticketHandler.PaymentLogs)
			}

			// === 管理员+超管 可访问 ===
			mgr := admin.Group("")
			mgr.Use(middleware.PermissionRequired("users"))
			{
				// User management
				mgr.GET("/users", adminPayHandler.UserList)
				mgr.GET("/users/:id", adminPayHandler.UserDetail)
				mgr.PUT("/users/:id/status", adminPayHandler.UserUpdateStatus)
				mgr.POST("/users/:id/recharge", adminPayHandler.UserRecharge)
			}

			// === 仅超管可访问 ===
			superOnly := admin.Group("")
			superOnly.Use(middleware.PermissionRequired("settings", "admins"))
			{
				// Settings
				superOnly.GET("/settings", ticketHandler.GetSettings)
				superOnly.PUT("/settings", ticketHandler.UpdateSettings)
				// Admin account management
				superOnly.GET("/admins", ticketHandler.AdminList)
				superOnly.POST("/admins", ticketHandler.AdminCreate)
				superOnly.PUT("/admins/:id", ticketHandler.AdminUpdate)
				superOnly.DELETE("/admins/:id", ticketHandler.AdminDelete)
				// Permission management
				superOnly.GET("/permissions/all", ticketHandler.GetAllPermissions)
				superOnly.GET("/permissions/:role", ticketHandler.GetRolePermissions)
				superOnly.PUT("/permissions/:role", ticketHandler.UpdateRolePermissions)
			}
		}
	}

	return r
}

func corsMiddleware() gin.HandlerFunc {
	allowedOrigins := map[string]bool{
		"https://zioran.com":       true,
		"https://www.zioran.com":   true,
		"https://admin.zioran.com": true,
		"http://localhost:3000":    true,
		"http://127.0.0.1:3000":    true,
		"http://localhost:5173":    true,
		"http://127.0.0.1:5173":    true,
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
