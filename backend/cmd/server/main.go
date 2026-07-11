package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/zioran/backend/internal/api"
	"github.com/zioran/backend/internal/repository"
	"github.com/zioran/backend/internal/service"
	"github.com/zioran/backend/pkg/config"
	"github.com/zioran/backend/pkg/email"
	"github.com/zioran/backend/pkg/oauth"
	ossClient "github.com/zioran/backend/pkg/oss"
	"github.com/zioran/backend/pkg/payment"
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

	// External services
	emailSender := email.NewSender(cfg.Email)
	wechatPay := payment.NewWechatPay(cfg.Payment.Wechat)
	alipayClient := payment.NewAlipayClient(cfg.Payment.Alipay)
	wechatOAuth := oauth.NewWechatOAuth(cfg.OAuth.Wechat)

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
	authSvc.SetEmailSender(emailSender)
	courseSvc := service.NewCourseService(courseRepo, catRepo, tagRepo, favRepo)
	paySvc := service.NewPaymentService(payRepo, courseRepo, userRepo, wechatPay, alipayClient)
	commSvc := service.NewCommunityService(commRepo)
	ticketSvc := service.NewTicketService(ticketRepo, userRepo)

	// Handlers
	authHandler := api.NewAuthHandler(authSvc)
	courseHandler := api.NewCourseHandler(courseSvc)
	adminHandler := api.NewAdminHandler(courseSvc)
	payHandler := api.NewPaymentHandler(paySvc)
	commHandler := api.NewCommunityHandler(commSvc)
	adminPayHandler := api.NewAdminPaymentHandler(paySvc, commSvc)
	uploadDir, err := resolveUploadDir(cfg.Server.UploadDir)
	if err != nil {
		log.Fatalf("resolve upload dir: %v", err)
	}
	log.Printf("Upload dir: %s", uploadDir)

	// OSS client (nil if not configured, falls back to local storage)
	var oss *ossClient.Client
	if cfg.OSS.Endpoint != "" {
		oss, err = ossClient.NewClient(cfg.OSS)
		if err != nil {
			log.Printf("WARNING: OSS init failed: %v (falling back to local storage)", err)
		} else {
			log.Printf("OSS enabled: bucket=%s, cdn=%s", cfg.OSS.Bucket, cfg.OSS.CDNDomain)
		}
	}

	// Inject OSS into services that need file cleanup
	if oss != nil {
		courseSvc.SetOSS(oss)
	}

	uploadHandler := api.NewUploadHandler(uploadDir, oss)
	ticketHandler := api.NewTicketHandler(ticketSvc, authSvc, paySvc, wechatOAuth, cfg.JWT.Secret, cfg.JWT.Expire, uploadDir, oss)

	r := api.SetupRouter(authHandler, courseHandler, adminHandler, payHandler, commHandler, adminPayHandler, uploadHandler, ticketHandler, cfg.JWT.Secret, db)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func resolveUploadDir(uploadDir string) (string, error) {
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	if filepath.IsAbs(uploadDir) {
		return filepath.Clean(uploadDir), nil
	}
	abs, err := filepath.Abs(uploadDir)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(filepath.Dir(abs)); err != nil {
		return "", err
	}
	return abs, nil
}
