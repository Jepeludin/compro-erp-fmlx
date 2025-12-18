package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	ppicScheduleRepo := repository.NewPPICScheduleRepository(sqlDB)
	tokenBlacklistRepo := repository.NewTokenBlacklistRepository(db)
	opPlanRepo := repository.NewOperationPlanRepository(db)
	gcodeRepo := repository.NewGCodeFileRepository(db)

	uploadPath := "./uploads/gcodes"

	// Initialize services
	authService := services.NewAuthService(userRepo, tokenBlacklistRepo, cfg)
	emailService := services.NewEmailService(cfg)
	opPlanService := services.NewOperationPlanService(opPlanRepo, gcodeRepo, jobOrderRepo, userRepo, emailService)
	gcodeService := services.NewGCodeService(gcodeRepo, opPlanRepo, uploadPath)
	ganttService := services.NewGanttService(ppicScheduleRepo)

	// Initialize and start cleanup service (cleans expired tokens every hour)
	cleanupService := services.NewCleanupService(tokenBlacklistRepo, services.DefaultCleanupConfig())
	cleanupService.Start()
	log.Println("Cleanup service: STARTED")

	// Log email service status
	if emailService.IsConfigured() {
		log.Println("Email service: CONFIGURED")
	} else {
		log.Println("Email service: NOT CONFIGURED (set SMTP_USERNAME, SMTP_PASSWORD, SMTP_FROM_ADDR)")
	}

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	machineHandler := handlers.NewMachineHandler(machineRepo)
	jobOrderHandler := handlers.NewJobOrderHandler(jobOrderRepo)
	adminHandler := handlers.NewAdminHandler(userRepo)
	opPlanHandler := handlers.NewOperationPlanHandler(opPlanService)
	gcodeHandler := handlers.NewGCodeHandler(gcodeService)
	ganttHandler := handlers.NewGanttHandler(ganttService)
	emailHandler := handlers.NewEmailHandler(emailService, opPlanRepo, userRepo)

	// Setup Gin router
	router := gin.Default()
	router.Use(middleware.CORS(cfg))

	// Setup routes and get rate limiters for graceful shutdown
	rateLimiters := routes.SetupRoutes(
		router,
		authHandler,
		machineHandler,
		jobOrderHandler,
		adminHandler,
		opPlanHandler,
		gcodeHandler,
		ganttHandler,
		emailHandler,
		authService,
	)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Stop cleanup service
	if err := cleanupService.Stop(ctx); err != nil {
		log.Printf("Cleanup service shutdown error: %v", err)
	}
	log.Println("Cleanup service stopped")

	// Stop rate limiters
	rateLimiters.Stop()
	log.Println("Rate limiters stopped")

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)

	}

	// Close database connection
	if err := sqlDB.Close(); err != nil {
		log.Printf("Database close error: %v", err)
	}

	log.Println("Server exited gracefully")
}
