package routes

import (
	"ganttpro-backend/handlers"
	"ganttpro-backend/middleware"
	"ganttpro-backend/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	machineHandler *handlers.MachineHandler,
	jobOrderHandler *handlers.JobOrderHandler,
	adminHandler *handlers.AdminHandler,
	opPlanHandler *handlers.OperationPlanHandler,
	gcodeHandler *handlers.GCodeHandler,
	authService *services.AuthService,
) {
	api := router.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/logout", middleware.AuthMiddleware(authService), authHandler.Logout)
		auth.GET("/profile", middleware.AuthMiddleware(authService), authHandler.GetProfile)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(authService))
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
			opPlans.POST("", opPlanHandler.CreateOperationPlan)                  // PEM creates plan
			opPlans.GET("", opPlanHandler.GetAllOperationPlans)                  // View all plans
			opPlans.GET("/:id", opPlanHandler.GetOperationPlan)                  // View specific plan
			opPlans.POST("/:id/submit", opPlanHandler.SubmitForApproval)         // Submit for approval
			opPlans.POST("/:id/approve", opPlanHandler.ApproveOperationPlan)     // Approve plan
			opPlans.GET("/pending-approvals", opPlanHandler.GetPendingApprovals) // Get pending approvals
			opPlans.POST("/:id/start", opPlanHandler.StartExecution)             // Start execution
			opPlans.POST("/:id/finish", opPlanHandler.FinishExecution)           // Finish execution
			opPlans.DELETE("/:id", opPlanHandler.DeleteOperationPlan)            // Delete plan (draft only)
		}

		// G-Code routes
		gcodes := protected.Group("/g-codes")
		{
			gcodes.POST("/upload", gcodeHandler.UploadGCode)               // Upload file
			gcodes.GET("/plan/:plan_id", gcodeHandler.GetGCodeFilesByPlan) // Get files by plan
			gcodes.GET("/:id/download", gcodeHandler.DownloadGCode)        // Download file
			gcodes.DELETE("/:id", gcodeHandler.DeleteGCode)                // Delete file
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
}
