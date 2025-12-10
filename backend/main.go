package main

import (
	"log"

	"ganttpro-backend/config"
	"ganttpro-backend/database"
	"ganttpro-backend/handlers"
	"ganttpro-backend/middleware"
	"ganttpro-backend/repository"
	"ganttpro-backend/routes"
	"ganttpro-backend/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Auto-migrate models
	database.AutoMigrate(db)

	// Get raw SQL DB connection for repositories that use *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get SQL DB connection:", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	machineRepo := repository.NewMachineRepository(sqlDB)
	jobOrderRepo := repository.NewJobOrderRepository(sqlDB)

	tokenBlacklistRepo := repository.NewTokenBlacklistRepository(db)
	opPlanRepo := repository.NewOperationPlanRepository(db)
	gcodeRepo := repository.NewGCodeFileRepository(db)
	uploadPath := "./uploads/gcodes"

	// Initialize services
	authService := services.NewAuthService(userRepo, tokenBlacklistRepo, cfg)
	opPlanService := services.NewOperationPlanService(opPlanRepo, gcodeRepo, jobOrderRepo)
	gcodeService := services.NewGCodeService(gcodeRepo, opPlanRepo, uploadPath)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	machineHandler := handlers.NewMachineHandler(machineRepo)
	jobOrderHandler := handlers.NewJobOrderHandler(jobOrderRepo)
	adminHandler := handlers.NewAdminHandler(userRepo)
	opPlanHandler := handlers.NewOperationPlanHandler(opPlanService)
	gcodeHandler := handlers.NewGCodeHandler(gcodeService)

	// Setup Gin router
	router := gin.Default()
	router.Use(middleware.CORS(cfg))

	// Setup routes
	routes.SetupRoutes(
		router,
		authHandler,
		machineHandler,
		jobOrderHandler,
		adminHandler,
		opPlanHandler,
		gcodeHandler,
		authService,
	)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
