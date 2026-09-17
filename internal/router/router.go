package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/xyz-corp/xyz-football-api/config"
	"github.com/xyz-corp/xyz-football-api/internal/handler"
	"github.com/xyz-corp/xyz-football-api/internal/middleware"
	"github.com/xyz-corp/xyz-football-api/pkg/jwt"
	"go.uber.org/zap"
)

type Handlers struct {
	Auth   *handler.AuthHandler
	Team   *handler.TeamHandler
	Player *handler.PlayerHandler
	Match  *handler.MatchHandler
	Report *handler.ReportHandler
}

func New(cfg *config.Config, jwtManager *jwt.Manager, logger *zap.Logger, h Handlers) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.CORS())
	r.Use(middleware.SecurityHeaders())

	r.MaxMultipartMemory = 8 << 20 // 8 MB

	r.Static("/uploads", cfg.UploadDir)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "OK"})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")

	authGroup := api.Group("/auth")
	authGroup.Use(middleware.RateLimiter(cfg.RateLimit))
	{
		authGroup.POST("/register", h.Auth.Register)
		authGroup.POST("/login", h.Auth.Login)
		authGroup.POST("/refresh", h.Auth.Refresh)
	}
	api.GET("/auth/me", middleware.JWTAuth(jwtManager), h.Auth.Me)

	protected := api.Group("")
	protected.Use(middleware.JWTAuth(jwtManager))
	{
		teams := protected.Group("/teams")
		teams.POST("", h.Team.Create)
		teams.GET("", h.Team.List)
		teams.GET("/:uuid", h.Team.Detail)
		teams.PUT("/:uuid", h.Team.Update)
		teams.DELETE("/:uuid", h.Team.Delete)
		teams.POST("/:uuid/logo", h.Team.UploadLogo)
		teams.GET("/:uuid/players", h.Team.Players)

		players := protected.Group("/players")
		players.POST("", h.Player.Create)
		players.GET("", h.Player.List)
		players.GET("/:uuid", h.Player.Detail)
		players.PUT("/:uuid", h.Player.Update)
		players.DELETE("/:uuid", h.Player.Delete)

		matches := protected.Group("/matches")
		matches.POST("", h.Match.Create)
		matches.GET("", h.Match.List)
		matches.GET("/:uuid", h.Match.Detail)
		matches.PUT("/:uuid", h.Match.Update)
		matches.PATCH("/:uuid/status", h.Match.UpdateStatus)
		matches.DELETE("/:uuid", h.Match.Delete)
		matches.POST("/:uuid/result", h.Match.SubmitResult)
		matches.PUT("/:uuid/result", h.Match.CorrectResult)
		matches.GET("/:uuid/events", h.Match.ListEvents)
		matches.POST("/:uuid/events", h.Match.AddEvent)
		matches.DELETE("/:uuid/events/:event_uuid", h.Match.DeleteEvent)

		reports := protected.Group("/reports")
		reports.GET("/matches", h.Report.List)
		reports.GET("/matches/:uuid", h.Report.Detail)
	}

	return r
}
