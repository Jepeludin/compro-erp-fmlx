package routes

import (
	"ganttpro-backend/handlers"
	"ganttpro-backend/middleware"
	"ganttpro-backend/services"

	"github.com/gin-gonic/gin"
)

// RateLimiters holds rate limiter instances for graceful shutdown
type RateLimiters struct {
	Auth *middleware.RateLimiter
	API  *middleware.RateLimiter
}

// Stop gracefully stops all rate limiters
func (r *RateLimiters) Stop() {
	if r.Auth != nil {
		r.Auth.Stop()
	}
	if r.API != nil {
		r.API.Stop()
	}
}

func SetupRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	machineHandler *handlers.MachineHandler,
	jobOrderHandler *handlers.JobOrderHandler,
	adminHandler *handlers.AdminHandler,
	opPlanHandler *handlers.OperationPlanHandler,
	gcodeHandler *handlers.GCodeHandler,
	ganttHandler *handlers.GanttHandler,
	ppicLinkHandler *handlers.PPICLinkHandler,
	emailHandler *handlers.EmailHandler,
	authService *services.AuthService,
) *RateLimiters {
	// Initialize rate limiters
	authRateLimiter := middleware.DefaultAuthRateLimiter() // 5 requests per minute for auth
	apiRateLimiter := middleware.DefaultAPIRateLimiter()   // 100 requests per minute for API

	api := router.Group("/api/v1")

	// Public routes with strict rate limiting (prevent brute force)
	auth := api.Group("/auth")
	{
		// Login and Register have strict rate limiting (5 attempts per minute)
		auth.POST("/login", authRateLimiter.RateLimit(), authHandler.Login)
		auth.POST("/register", authRateLimiter.RateLimit(), authHandler.Register)
		// Logout and Profile are protected by auth and have normal rate limiting
		auth.POST("/logout", middleware.AuthMiddleware(authService), apiRateLimiter.RateLimit(), authHandler.Logout)
		auth.GET("/profile", middleware.AuthMiddleware(authService), apiRateLimiter.RateLimit(), authHandler.GetProfile)
	}

	// Protected routes with authentication and rate limiting
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(authService))
	protected.Use(apiRateLimiter.RateLimit())
	{
		// Machine routes
		machines := protected.Group("/machines")
		{
			machines.GET("", machineHandler.GetAllMachines)
			machines.GET("/:id", machineHandler.GetMachine)
		}

		// Job Order routes
		jobOrders := protected.Group("/job-orders")
		{
			jobOrders.GET("", jobOrderHandler.GetAllJobOrders)
			jobOrders.GET("/:id", jobOrderHandler.GetJobOrder)
			jobOrders.GET("/machine/:machine_id", jobOrderHandler.GetJobOrdersByMachine)
			jobOrders.POST("", jobOrderHandler.CreateJobOrder)
			jobOrders.PUT("/:id", jobOrderHandler.UpdateJobOrder)
			jobOrders.DELETE("/:id", jobOrderHandler.DeleteJobOrder)
		}

		// Process Stage routes
		processStages := protected.Group("/process-stages")
		{
			processStages.PUT("/:id", jobOrderHandler.UpdateProcessStage)
		}

		// Operation Plan routes
		opPlans := protected.Group("/operation-plans")
		{
			opPlans.POST("", opPlanHandler.CreateOperationPlan)                   // PEM creates plan
			opPlans.GET("", opPlanHandler.GetAllOperationPlans)                   // View all plans
			opPlans.GET("/:id", opPlanHandler.GetOperationPlan)                   // View specific plan
			opPlans.POST("/:id/submit", opPlanHandler.SubmitForApproval)          // Submit for approval
			opPlans.POST("/:id/approve", opPlanHandler.ApproveOperationPlan)      // Approve plan
			opPlans.GET("/pending-approvals", opPlanHandler.GetPendingApprovals)  // Get pending approvals
			opPlans.POST("/:id/start", opPlanHandler.StartExecution)              // Start execution
			opPlans.POST("/:id/finish", opPlanHandler.FinishExecution)            // Finish execution
			opPlans.DELETE("/:id", opPlanHandler.DeleteOperationPlan)             // Delete plan (draft only)
			opPlans.POST("/:id/send-reminder", emailHandler.SendApprovalReminder) // Send reminder emails (PEM only)
		}

		// Email routes
		email := protected.Group("/email")
		{
			email.GET("/status", emailHandler.CheckEmailConfig) // Check email configuration
		}

		// G-Code routes
		gcodes := protected.Group("/g-codes")
		{
			gcodes.POST("/upload", gcodeHandler.UploadGCode)               // Upload file
			gcodes.GET("/plan/:plan_id", gcodeHandler.GetGCodeFilesByPlan) // Get files by plan
			gcodes.GET("/:id/download", gcodeHandler.DownloadGCode)        // Download file
			gcodes.DELETE("/:id", gcodeHandler.DeleteGCode)                // Delete file
		}

		// Gantt Chart routes
		gantt := protected.Group("/gantt-chart")
		{
			gantt.GET("", ganttHandler.GetGanttChart) // Get Gantt chart data with filters
		}

		// PPIC Schedule routes (for Gantt chart data management)
		ppic := protected.Group("/ppic-schedules")
		{
			ppic.GET("", ganttHandler.GetAllPPICSchedules)                                              // Get all schedules
			ppic.GET("/:id", ganttHandler.GetPPICSchedule)                                              // Get single schedule
			ppic.POST("", ganttHandler.CreatePPICSchedule)                                              // Create schedule
			ppic.PUT("/:id", ganttHandler.UpdatePPICSchedule)                                           // Update schedule
			ppic.DELETE("/:id", ganttHandler.DeletePPICSchedule)                                        // Delete schedule
			ppic.GET("/machine/:machine_id", ganttHandler.GetSchedulesByMachine)                        // Get by machine
			ppic.POST("/:id/machines", ganttHandler.AddMachineAssignment)                               // Add machine
			ppic.DELETE("/:id/machines/:assignment_id", ganttHandler.RemoveMachineAssignment)           // Remove machine
			ppic.PUT("/:id/machines/:assignment_id/status", ganttHandler.UpdateMachineAssignmentStatus) // Update status
		}

		// PPIC Links routes (for Gantt chart dependencies/arrows)
		ppicLinks := protected.Group("/ppic-links")
		{
			ppicLinks.GET("", ppicLinkHandler.GetAllPPICLinks)      // Get all links
			ppicLinks.POST("", ppicLinkHandler.CreatePPICLink)      // Create link
			ppicLinks.DELETE("/:id", ppicLinkHandler.DeletePPICLink) // Delete link
		}

		// Admin routes
		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRole("Admin"))
		{
			// User management
			admin.GET("/users", adminHandler.GetAllUsers)
			admin.PUT("/users/:id", adminHandler.UpdateUser)
			admin.DELETE("/users/:id", adminHandler.DeleteUser)

			// Machine management
			admin.POST("/machines", machineHandler.CreateMachine)
			admin.PUT("/machines/:id", machineHandler.UpdateMachine)
			admin.DELETE("/machines/:id", machineHandler.DeleteMachine)
		}
	}

	// Return rate limiters for graceful shutdown
	return &RateLimiters{
		Auth: authRateLimiter,
		API:  apiRateLimiter,
	}
}
