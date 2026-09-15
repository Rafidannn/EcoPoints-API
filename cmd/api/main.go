package main

import (
	"fmt"
	"log"
	"net/http"

	"ecopoints-go-api/internal/config"
	"ecopoints-go-api/internal/dto"
	"ecopoints-go-api/internal/handler"
	"ecopoints-go-api/internal/middleware"
	"ecopoints-go-api/internal/repository"
	"ecopoints-go-api/internal/service"

	docs "ecopoints-go-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title EcoPoints Golang REST API
// @version 1.0
// @description REST API for EcoPoints application built with Go Gin and connected to existing Laravel MySQL database.
// @termsOfService http://swagger.io/terms/

// @contact.name EcoPoints Dev Team
// @contact.email support@ecopoints.test

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer " followed by your JWT token. Example: "Bearer eyJhbGciOi..."

func main() {
	// Dynamically use current host and port
	docs.SwaggerInfo.Host = ""

	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Set Gin mode
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 3. Connect to existing Laravel database
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 4. Initialize Repositories
	userRepo := repository.NewUserRepository(db)
	wasteTypeRepo := repository.NewWasteTypeRepository(db)
	dropPointRepo := repository.NewDropPointRepository(db)
	rewardRepo := repository.NewRewardRepository(db)
	leaderboardRepo := repository.NewLeaderboardRepository(db)
	wasteDepositRepo := repository.NewWasteDepositRepository(db)

	// 5. Initialize Services
	jwtService := service.NewJWTService(cfg.JWTSecret, cfg.JWTExpirationHours)
	authService := service.NewAuthService(userRepo, jwtService)
	wasteTypeService := service.NewWasteTypeService(wasteTypeRepo)
	dropPointService := service.NewDropPointService(dropPointRepo)
	rewardService := service.NewRewardService(rewardRepo)
	leaderboardService := service.NewLeaderboardService(leaderboardRepo)
	wasteDepositService := service.NewWasteDepositService(wasteDepositRepo, wasteTypeRepo)

	// 6. Initialize Handlers
	authHandler := handler.NewAuthHandler(authService)
	wasteTypeHandler := handler.NewWasteTypeHandler(wasteTypeService)
	dropPointHandler := handler.NewDropPointHandler(dropPointService)
	rewardHandler := handler.NewRewardHandler(rewardService)
	leaderboardHandler := handler.NewLeaderboardHandler(leaderboardService)
	wasteDepositHandler := handler.NewWasteDepositHandler(wasteDepositService)
	adminUserHandler := handler.NewAdminUserHandler(userRepo)

	// 7. Initialize Gin Router
	router := gin.Default()
	_ = router.SetTrustedProxies(nil)

	// CORS Middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Root route: auto-redirect to Swagger documentation
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, dto.SuccessResponse("EcoPoints Go API is healthy and connected to database", gin.H{
			"app_name": cfg.AppName,
			"database": cfg.DBDatabase,
		}))
	})

	// Swagger documentation route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 Routes
	v1 := router.Group("/api/v1")
	{
		// 1. Auth Routes
		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.GET("/me", middleware.AuthMiddleware(jwtService), authHandler.Me)
			authGroup.PUT("/change-password", middleware.AuthMiddleware(jwtService), authHandler.ChangePassword)
		}

		// 2. Waste Types Routes (CRUD)
		wasteTypesGroup := v1.Group("/waste-types")
		{
			wasteTypesGroup.GET("", wasteTypeHandler.GetAll)
			wasteTypesGroup.GET("/:id", wasteTypeHandler.GetByID)

			// Admin/Petugas only write routes
			adminWaste := wasteTypesGroup.Group("")
			adminWaste.Use(middleware.AuthMiddleware(jwtService), middleware.RoleMiddleware("admin", "petugas"))
			{
				adminWaste.POST("", wasteTypeHandler.Create)
				adminWaste.PUT("/:id", wasteTypeHandler.Update)
				adminWaste.DELETE("/:id", wasteTypeHandler.Delete)
			}
		}

		// 3. Drop Points Routes (CRUD)
		dropPointsGroup := v1.Group("/drop-points")
		{
			dropPointsGroup.GET("", dropPointHandler.GetAll)
			dropPointsGroup.GET("/:id", dropPointHandler.GetByID)

			// Admin/Petugas only write routes
			adminDropPoints := dropPointsGroup.Group("")
			adminDropPoints.Use(middleware.AuthMiddleware(jwtService), middleware.RoleMiddleware("admin", "petugas"))
			{
				adminDropPoints.POST("", dropPointHandler.Create)
				adminDropPoints.PUT("/:id", dropPointHandler.Update)
				adminDropPoints.DELETE("/:id", dropPointHandler.Delete)
			}
		}

		// 4. Rewards Routes (CRUD)
		rewardsGroup := v1.Group("/rewards")
		{
			rewardsGroup.GET("", rewardHandler.GetAll)
			rewardsGroup.GET("/:id", rewardHandler.GetByID)

			// Authenticated user routes (redeem & my-redemptions)
			authRewards := rewardsGroup.Group("")
			authRewards.Use(middleware.AuthMiddleware(jwtService))
			{
				authRewards.POST("/:id/redeem", rewardHandler.Redeem)
				authRewards.GET("/my-redemptions", rewardHandler.GetMyRedemptions)

				// Admin/Petugas redemption management routes
				staffRedemptions := authRewards.Group("")
				staffRedemptions.Use(middleware.RoleMiddleware("admin", "petugas"))
				{
					staffRedemptions.GET("/redemptions", rewardHandler.GetAllRedemptions)
					staffRedemptions.PUT("/redemptions/:id/verify", rewardHandler.CompleteRedemption)
					staffRedemptions.PUT("/redemptions/:id/reject", rewardHandler.RejectRedemption)
				}
			}

			// Admin only write routes
			adminRewards := rewardsGroup.Group("")
			adminRewards.Use(middleware.AuthMiddleware(jwtService), middleware.RoleMiddleware("admin"))
			{
				adminRewards.POST("", rewardHandler.Create)
				adminRewards.PUT("/:id", rewardHandler.Update)
				adminRewards.DELETE("/:id", rewardHandler.Delete)
			}
		}

		adminUsers := v1.Group("/admin/users")
		adminUsers.Use(middleware.AuthMiddleware(jwtService), middleware.RoleMiddleware("admin"))
		{
			adminUsers.GET("", adminUserHandler.GetAll)
			adminUsers.POST("", adminUserHandler.Create)
			adminUsers.PUT("/:id", adminUserHandler.Update)
			adminUsers.DELETE("/:id", adminUserHandler.Delete)
		}

		// 5. Leaderboard Route (public)
		v1.GET("/leaderboard", leaderboardHandler.GetLeaderboard)

		// 6. Waste Deposits Routes
		wasteDepositsGroup := v1.Group("/waste-deposits")
		wasteDepositsGroup.Use(middleware.AuthMiddleware(jwtService))
		{
			wasteDepositsGroup.POST("", wasteDepositHandler.Create)
			wasteDepositsGroup.GET("", wasteDepositHandler.GetAll)
			wasteDepositsGroup.GET("/:id", wasteDepositHandler.GetByID)
			wasteDepositsGroup.PUT("/:id/verify", middleware.RoleMiddleware("admin", "petugas"), wasteDepositHandler.Verify)
			wasteDepositsGroup.PUT("/:id/reject", middleware.RoleMiddleware("admin", "petugas"), wasteDepositHandler.Reject)
			wasteDepositsGroup.PUT("/:id/cancel", wasteDepositHandler.Cancel)
		}
	}

	// 8. Start server
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server starting on port %s (Swagger docs: http://localhost:%s/swagger/index.html)", cfg.AppPort, cfg.AppPort)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
