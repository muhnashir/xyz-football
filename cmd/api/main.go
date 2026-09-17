package main

import (
	"log"

	"github.com/xyz-corp/xyz-football-api/config"
	"github.com/xyz-corp/xyz-football-api/docs"
	"github.com/xyz-corp/xyz-football-api/internal/handler"
	"github.com/xyz-corp/xyz-football-api/internal/repository"
	"github.com/xyz-corp/xyz-football-api/internal/router"
	"github.com/xyz-corp/xyz-football-api/internal/service"
	"github.com/xyz-corp/xyz-football-api/pkg/jwt"
	"github.com/xyz-corp/xyz-football-api/pkg/uploader"
	"go.uber.org/zap"
)

// @title       XYZ Football Team Management API
// @version     1.0
// @description Backend API untuk manajemen tim, pemain, jadwal, dan hasil pertandingan sepak bola amatir.
// @host        localhost:8080
// @BasePath    /api/v1
// @securityDefinitions.apikey BearerAuth
// @in          header
// @name        Authorization
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Host swagger dinamis mengikuti APP_PORT agar "Try it out"
	// tetap benar saat port diubah (mis. 8080 vs 8081).
	docs.SwaggerInfo.Host = "localhost:" + cfg.AppPort

	var logger *zap.Logger
	if cfg.AppEnv == "development" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync()

	db, err := config.NewDatabase(cfg)
	if err != nil {
		logger.Fatal("failed to connect database", zap.Error(err))
	}

	jwtManager := jwt.NewManager(cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	up := uploader.New(cfg.UploadDir, cfg.BaseURL)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	playerRepo := repository.NewPlayerRepository(db)
	matchRepo := repository.NewMatchRepository(db)
	matchEventRepo := repository.NewMatchEventRepository(db)
	reportRepo := repository.NewReportRepository(db)

	// Services
	authSvc := service.NewAuthService(userRepo, jwtManager)
	teamSvc := service.NewTeamService(teamRepo, up)
	playerSvc := service.NewPlayerService(playerRepo, teamRepo)
	matchSvc := service.NewMatchService(matchRepo, teamRepo)
	matchEventSvc := service.NewMatchEventService(matchRepo, playerRepo, matchEventRepo)
	reportSvc := service.NewReportService(matchRepo, reportRepo, cfg)

	// Handlers
	handlers := router.Handlers{
		Auth:   handler.NewAuthHandler(authSvc),
		Team:   handler.NewTeamHandler(teamSvc, playerSvc),
		Player: handler.NewPlayerHandler(playerSvc),
		Match:  handler.NewMatchHandler(matchSvc, matchEventSvc),
		Report: handler.NewReportHandler(reportSvc),
	}

	engine := router.New(cfg, jwtManager, logger, handlers)

	logger.Info("server starting", zap.String("port", cfg.AppPort))
	if err := engine.Run(":" + cfg.AppPort); err != nil {
		logger.Fatal("server failed to start", zap.Error(err))
	}
}
