package main

import (
	"fmt"
	"log"

	"github.com/zioran/backend/internal/api"
	"github.com/zioran/backend/internal/repository"
	"github.com/zioran/backend/internal/service"
	"github.com/zioran/backend/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}

	// Repositories
	userRepo := repository.NewUserRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	catRepo := repository.NewCategoryRepository(db)
	tagRepo := repository.NewTagRepository(db)
	favRepo := repository.NewFavoriteRepository(db)
	payRepo := repository.NewPaymentRepository(db)
	commRepo := repository.NewCommunityRepository(db)
	ticketRepo := repository.NewTicketRepository(db)

	// Services
	authSvc := service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.Expire)
	courseSvc := service.NewCourseService(courseRepo, catRepo, tagRepo, favRepo)
	paySvc := service.NewPaymentService(payRepo, courseRepo, userRepo)
	commSvc := service.NewCommunityService(commRepo)
	ticketSvc := service.NewTicketService(ticketRepo, userRepo)

	// Handlers
	authHandler := api.NewAuthHandler(authSvc)
	courseHandler := api.NewCourseHandler(courseSvc)
	adminHandler := api.NewAdminHandler(courseSvc)
	payHandler := api.NewPaymentHandler(paySvc)
	commHandler := api.NewCommunityHandler(commSvc)
	adminPayHandler := api.NewAdminPaymentHandler(paySvc, commSvc)
	uploadHandler := api.NewUploadHandler("./uploads")
	ticketHandler := api.NewTicketHandler(ticketSvc, authSvc, cfg.JWT.Secret, cfg.JWT.Expire, "./uploads")

	r := api.SetupRouter(authHandler, courseHandler, adminHandler, payHandler, commHandler, adminPayHandler, uploadHandler, ticketHandler, cfg.JWT.Secret)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}
