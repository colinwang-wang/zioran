package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zioran/backend/internal/middleware"
)

func SetupRouter(
	authHandler *AuthHandler,
	courseHandler *CourseHandler,
	adminHandler *AdminHandler,
	payHandler *PaymentHandler,
	commHandler *CommunityHandler,
	adminPayHandler *AdminPaymentHandler,
	jwtSecret string,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		// Auth
		auth := v1.Group("/auth")
		{
			auth.POST("/captcha", authHandler.Captcha)
			auth.POST("/sms/send", authHandler.SendSMS)
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
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

		// Public VIP packages
		v1.GET("/vip/packages", payHandler.VipPackages)

		// Public guestbook (read)
		v1.GET("/guestbook", commHandler.GuestbookList)

		// Public comments (read)
		v1.GET("/comments", commHandler.CommentList)

		// Authed routes
		authed := v1.Group("")
		authed.Use(middleware.JWTAuth(jwtSecret))
		{
			// User profile
			authed.GET("/user/profile", authHandler.Profile)
			authed.PUT("/user/password", payHandler.ChangePassword)
			authed.GET("/user/orders", payHandler.UserOrders)
			authed.GET("/user/downloads", payHandler.UserDownloads)
			authed.GET("/user/favorites", payHandler.UserFavorites)
			authed.POST("/user/favorites", payHandler.AddFavorite)
			authed.DELETE("/user/favorites/:courseId", payHandler.RemoveFavorite)

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

			// Guestbook (write)
			authed.POST("/guestbook", commHandler.GuestbookCreate)
			authed.POST("/guestbook/:id/like", commHandler.GuestbookLike)
			authed.DELETE("/guestbook/:id", commHandler.GuestbookDelete)

			// Comments (write)
			authed.POST("/comments", commHandler.CommentCreate)
			authed.DELETE("/comments/:id", commHandler.CommentDelete)
		}

		// Admin routes
		admin := v1.Group("/admin")
		admin.Use(middleware.JWTAuth(jwtSecret))
		{
			// Course management
			admin.GET("/courses", adminHandler.CourseList)
			admin.POST("/courses", adminHandler.CourseCreate)
			admin.PUT("/courses/:id", adminHandler.CourseUpdate)
			admin.DELETE("/courses/:id", adminHandler.CourseDelete)
			admin.PUT("/courses/:id/status", adminHandler.CourseUpdateStatus)

			// Category management
			admin.GET("/categories", adminHandler.CategoryList)
			admin.POST("/categories", adminHandler.CategoryCreate)
			admin.PUT("/categories/:id", adminHandler.CategoryUpdate)
			admin.DELETE("/categories/:id", adminHandler.CategoryDelete)

			// Tag management
			admin.GET("/tags", adminHandler.TagList)
			admin.POST("/tags", adminHandler.TagCreate)
			admin.PUT("/tags/:id", adminHandler.TagUpdate)
			admin.DELETE("/tags/:id", adminHandler.TagDelete)

			// Order management
			admin.GET("/orders", adminPayHandler.OrderList)
			admin.POST("/orders/:id/refund", adminPayHandler.OrderRefund)

			// User management
			admin.GET("/users", adminPayHandler.UserList)
			admin.PUT("/users/:id/status", adminPayHandler.UserUpdateStatus)
			admin.POST("/users/:id/recharge", adminPayHandler.UserRecharge)

			// Guestbook management
			admin.GET("/guestbook", adminPayHandler.GuestbookList)
			admin.PUT("/guestbook/:id/status", adminPayHandler.GuestbookUpdateStatus)
			admin.PUT("/guestbook/:id/pin", adminPayHandler.GuestbookPin)
			admin.DELETE("/guestbook/:id", adminPayHandler.GuestbookDelete)

			// Comment management
			admin.GET("/comments", adminPayHandler.CommentList)
			admin.PUT("/comments/:id/status", adminPayHandler.CommentUpdateStatus)
			admin.DELETE("/comments/:id", adminPayHandler.CommentDelete)

			// Nav items management
			admin.GET("/nav-items", adminPayHandler.NavItemList)
			admin.POST("/nav-items", adminPayHandler.NavItemCreate)
			admin.PUT("/nav-items/:id", adminPayHandler.NavItemUpdate)
			admin.DELETE("/nav-items/:id", adminPayHandler.NavItemDelete)

			// Banner management
			admin.GET("/banners", adminPayHandler.BannerList)
			admin.POST("/banners", adminPayHandler.BannerCreate)
			admin.PUT("/banners/:id", adminPayHandler.BannerUpdate)
			admin.DELETE("/banners/:id", adminPayHandler.BannerDelete)

			// Dashboard
			admin.GET("/dashboard/stats", adminPayHandler.DashboardStats)
		}
	}

	return r
}
