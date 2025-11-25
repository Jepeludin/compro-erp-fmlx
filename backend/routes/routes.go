package routes

import (
	"net/http"

	"ganttpro-backend/config"
	"ganttpro-backend/handlers"
	"ganttpro-backend/middleware"
	"ganttpro-backend/repository"
	"ganttpro-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configures all application routes
func SetupRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Apply CORS middleware
	router.Use(middleware.CORS(cfg))

	// Get underlying SQL database for repositories
	sqlDB, err := db.DB()
	if err != nil {
		panic("Failed to get SQL database: " + err.Error())
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	machineRepo := repository.NewMachineRepository(sqlDB)
	jobOrderRepo := repository.NewJobOrderRepository(sqlDB)
	tokenBlacklistRepo := repository.NewTokenBlacklistRepository(db)
	// Initialize services
	authService := services.NewAuthService(userRepo, tokenBlacklistRepo, cfg)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(userRepo)
	machineHandler := handlers.NewMachineHandler(machineRepo)
	jobOrderHandler := handlers.NewJobOrderHandler(jobOrderRepo)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "GanttPro API is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/logout", middleware.AuthMiddleware(authService), authHandler.Logout)
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// User profile
			protected.GET("/auth/profile", authHandler.GetProfile)

			// Machine routes (all authenticated users can view)
			machines := protected.Group("/machines")
			{
				machines.GET("", machineHandler.GetAllMachines)
				machines.GET("/:id", machineHandler.GetMachine)
			}

			// Job Order routes
			jobOrders := protected.Group("/job-orders")
			{
				jobOrders.GET("", jobOrderHandler.GetAllJobOrders)
				jobOrders.GET("/machine/:machine_id", jobOrderHandler.GetJobOrdersByMachine)
				jobOrders.GET("/:id", jobOrderHandler.GetJobOrder)
				jobOrders.POST("", jobOrderHandler.CreateJobOrder)
				jobOrders.PUT("/:id", jobOrderHandler.UpdateJobOrder)
				jobOrders.DELETE("/:id", jobOrderHandler.DeleteJobOrder)
			}

			// Process Stage routes
			protected.PUT("/process-stages/:id", jobOrderHandler.UpdateProcessStage)
		}

		// Admin only routes
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(authService))
		admin.Use(middleware.RequireAdmin())
		{
			// User management
			admin.GET("/users", adminHandler.GetAllUsers)
			admin.PUT("/users/:id", adminHandler.UpdateUser)
			admin.DELETE("/users/:id", adminHandler.DeleteUser)

			// Machine management (admin only can create/update/delete)
			admin.POST("/machines", machineHandler.CreateMachine)
			admin.PUT("/machines/:id", machineHandler.UpdateMachine)
			admin.DELETE("/machines/:id", machineHandler.DeleteMachine)
		}
	}
}
