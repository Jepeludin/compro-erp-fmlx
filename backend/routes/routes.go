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
	googleSheetsHandler *handlers.GoogleSheetsHandler,
	pemPlanHandler *handlers.PEMOperationPlanHandler,
	toolpatherFileHandler *handlers.ToolpatherFileHandler,
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

		// Google Sheets routes
		googleSheets := protected.Group("/google-sheets")
		{
			googleSheets.GET("/part-name/:orderNumber", googleSheetsHandler.GetPartNameByOrderNumber) // Get part name by order number
			googleSheets.GET("/all-data", googleSheetsHandler.GetAllSheetData)                        // Get all Google Sheets data
		}

		// Users routes (for approver selection - accessible by all authenticated users)
		protected.GET("/users", adminHandler.GetUsers) // Get active users for approver selection

		// PEM Operation Plans routes
		pemPlans := protected.Group("/pem-operation-plans")
		{
			pemPlans.GET("", pemPlanHandler.GetAllPEMPlans)                           // Get all PEM plans
			pemPlans.GET("/:id", pemPlanHandler.GetPEMPlan)                           // Get single PEM plan
			pemPlans.POST("", pemPlanHandler.CreatePEMPlan)                           // Create PEM plan
			pemPlans.PUT("/:id", pemPlanHandler.UpdatePEMPlan)                        // Update PEM plan
			pemPlans.DELETE("/:id", pemPlanHandler.DeletePEMPlan)                     // Delete PEM plan

			// Steps management
			pemPlans.POST("/:id/steps", pemPlanHandler.AddPlanStep)                  // Add step to plan
			pemPlans.PUT("/steps/:step_id", pemPlanHandler.UpdatePlanStep)           // Update step
			pemPlans.DELETE("/steps/:step_id", pemPlanHandler.DeletePlanStep)        // Delete step

			// Image upload
			pemPlans.POST("/steps/:step_id/image", pemPlanHandler.UploadStepImage)   // Upload step image
			pemPlans.DELETE("/steps/:step_id/image", pemPlanHandler.DeleteStepImage) // Delete step image

			// Approval workflow
			pemPlans.POST("/:id/assign-approvers", pemPlanHandler.AssignApprovers)   // Assign all 5 approvers
			pemPlans.POST("/:id/submit", pemPlanHandler.SubmitPlanForApproval)       // Submit for approval
			pemPlans.POST("/:id/approve", pemPlanHandler.ApprovePlan)                // Approve plan (requires ?role= param)
			pemPlans.POST("/:id/reject", pemPlanHandler.RejectPlan)                  // Reject plan (requires ?role= param)

			// Filter by PPIC schedule
			pemPlans.GET("/ppic-schedule/:schedule_id", pemPlanHandler.GetPlansByPPICSchedule) // Get plans by PPIC schedule
			pemPlans.GET("/pending-approvals", pemPlanHandler.GetPendingApprovals)   // Get pending approvals for current user
		}

		// Toolpather File Upload routes
		toolpatherFiles := protected.Group("/toolpather-files")
		{
			toolpatherFiles.POST("/upload", toolpatherFileHandler.UploadFiles)                       // Upload multiple .txt files
			toolpatherFiles.GET("", toolpatherFileHandler.GetAllFiles)                               // Get all files with filters
			toolpatherFiles.GET("/my-files", toolpatherFileHandler.GetMyFiles)                       // Get current user's files
			toolpatherFiles.GET("/:id", toolpatherFileHandler.GetFileByID)                           // Get single file
			toolpatherFiles.GET("/order/:orderNumber", toolpatherFileHandler.GetFilesByOrderNumber)  // Get files by order number
			toolpatherFiles.GET("/:id/download", toolpatherFileHandler.DownloadFile)                 // Download file
			toolpatherFiles.DELETE("/:id", toolpatherFileHandler.DeleteFile)                         // Delete file
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
