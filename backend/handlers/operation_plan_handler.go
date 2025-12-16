package handlers

import (
    "ganttpro-backend/models"
    "ganttpro-backend/services"
    "net/http"
    "os"
    "strconv"

    "github.com/gin-gonic/gin"
)

type OperationPlanHandler struct {
    service *services.OperationPlanService
}

func NewOperationPlanHandler(service *services.OperationPlanService) *OperationPlanHandler {
    return &OperationPlanHandler{service: service}
}

// CreateOperationPlan creates a new operation plan
// @Summary Create operation plan
// @Description Create a new operation plan for a job order
// @Tags operation-plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.CreateOperationPlanRequest true "Operation plan details"
// @Success 201 {object} models.OperationPlan
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/operation-plans [post]
func (h *OperationPlanHandler) CreateOperationPlan(c *gin.Context) {
    var req services.CreateOperationPlanRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
        return
    }

    // Get current user from context
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
        return
    }
    currentUser := user.(*models.User)

    plan, err := h.service.Create(req, currentUser.ID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Operation plan created successfully",
        "data":    plan,
    })
}

// GetOperationPlan gets an operation plan by ID
// @Summary Get operation plan
// @Description Get operation plan details by ID
// @Tags operation-plans
// @Produce json
// @Security BearerAuth
// @Param id path int true "Operation Plan ID"
// @Success 200 {object} models.OperationPlanResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/operation-plans/{id} [get]
func (h *OperationPlanHandler) GetOperationPlan(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation plan ID"})
        return
    }

    plan, err := h.service.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Operation plan not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": plan})
}

// GetOperationPlanByJobOrder gets an operation plan by job order ID
// @Summary Get operation plan by job order
// @Description Get operation plan for a specific job order
// @Tags operation-plans
// @Produce json
// @Security BearerAuth
// @Param job_order_id path int true "Job Order ID"
// @Success 200 {object} models.OperationPlanResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/operation-plans/job-order/{job_order_id} [get]
func (h *OperationPlanHandler) GetOperationPlanByJobOrder(c *gin.Context) {
    jobOrderID, err := strconv.ParseUint(c.Param("job_order_id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job order ID"})
        return
    }

    plan, err := h.service.GetByJobOrderID(uint(jobOrderID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Operation plan not found for this job order"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": plan})
}

// GetAllOperationPlans gets all operation plans
// @Summary Get all operation plans
// @Description Get all operation plans with optional filters
// @Tags operation-plans
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status (draft, pending_approval, approved)"
// @Param machine_id query int false "Filter by machine ID"
// @Success 200 {array} models.OperationPlanResponse
// @Router /api/v1/operation-plans [get]
func (h *OperationPlanHandler) GetAllOperationPlans(c *gin.Context) {
    status := c.Query("status")
    machineIDStr := c.Query("machine_id")

    var machineID uint
    if machineIDStr != "" {
        id, err := strconv.ParseUint(machineIDStr, 10, 32)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine ID"})
            return
        }
        machineID = uint(id)
    }

    plans, err := h.service.GetAll(status, machineID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch operation plans"})
        return
    }

    // Apply filters if needed
    var filteredPlans []models.OperationPlanResponse
    for _, plan := range plans {
        if status != "" && plan.Status != status {
            continue
        }
        if machineID != 0 && plan.MachineID != machineID {
            continue
        }
        filteredPlans = append(filteredPlans, plan.ToResponse())
    }

    c.JSON(http.StatusOK, gin.H{"data": filteredPlans})
}

// GetMyOperationPlans gets operation plans created by the current user
// @Summary Get my operation plans
// @Description Get operation plans created by the authenticated user
// @Tags operation-plans
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.OperationPlanResponse
// @Router /api/v1/operation-plans/my-plans [get]
func (h *OperationPlanHandler) GetMyOperationPlans(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
        return
    }
    currentUser := user.(*models.User)

    plans, err := h.service.GetByCreator(currentUser.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch operation plans"})
        return
    }

    // Filter by creator
    var myPlans []models.OperationPlanResponse
    for _, plan := range plans {
        if plan.CreatedBy == currentUser.ID {
            myPlans = append(myPlans, plan.ToResponse())
        }
    }

    c.JSON(http.StatusOK, gin.H{"data": myPlans})
}

// UpdateOperationPlan updates an operation plan
// @Summary Update operation plan
// @Description Update an operation plan (only if status is draft)
// @Tags operation-plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Operation Plan ID"
// @Param request body services.CreateOperationPlanRequest true "Update details"
// @Success 200 {object} models.OperationPlan
// @Failure 400 {object} map[string]string
// @Router /api/v1/operation-plans/{id} [put]
func (h *OperationPlanHandler) UpdateOperationPlan(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation plan ID"})
        return
    }

    var req services.CreateOperationPlanRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
        return
    }

    plan, err := h.service.Update(uint(id), services.UpdateOperationPlanRequest{
        PartQuantity: req.PartQuantity,
        Description:  req.Description,
    })
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Operation plan updated successfully",
        "data":    plan,
    })
}

// DeleteOperationPlan deletes an operation plan
// @Summary Delete operation plan
// @Description Delete an operation plan (only if status is draft)
// @Tags operation-plans
// @Produce json
// @Security BearerAuth
// @Param id path int true "Operation Plan ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/operation-plans/{id} [delete]
func (h *OperationPlanHandler) DeleteOperationPlan(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation plan ID"})
        return
    }

    if err := h.service.Delete(uint(id)); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Operation plan deleted successfully",
    })
}

// SubmitForApproval submits an operation plan for approval
// @Summary Submit for approval
// @Description Submit an operation plan for approval
// @Tags operation-plans
// @Produce json
// @Security BearerAuth
// @Param id path int true "Operation Plan ID"
// @Success 200 {object} models.OperationPlanResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/operation-plans/{id}/submit [post]
func (h *OperationPlanHandler) SubmitForApproval(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation plan ID"})
        return
    }

    // Get current user
    plan, err := h.service.SubmitForApproval(uint(id))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Get updated plan
    plan, err = h.service.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Operation plan not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Operation plan submitted for approval",
        "data":    plan,
    })
}

// ApproveOperationPlan approves an operation plan
// @Summary Approve operation plan
// @Description Approve an operation plan (role-based)
// @Tags operation-plans
// @Produce json
// @Security BearerAuth
// @Param id path int true "Operation Plan ID"
// @Success 200 {object} models.OperationPlanResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /api/v1/operation-plans/{id}/approve [post]
func (h *OperationPlanHandler) ApproveOperationPlan(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation plan ID"})
        return
    }

    // Get current user from context
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
        return
    }
    currentUser := user.(*models.User)

    // Check if user's role is an approver role
    validRoles := []string{
        models.RolePEM,
        models.RolePPIC,
        models.RoleQC,
        models.RoleEngineering,
        models.RoleToolpather,
    }

    isApprover := false
    for _, role := range validRoles {
        if role == currentUser.Role {
            isApprover = true
            break
        }
    }

    if !isApprover {
        c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to approve operation plans"})
        return
    }

    plan, err := h.service.Approve(uint(id), currentUser.ID, currentUser.Role)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Get updated plan
    plan, err = h.service.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Operation plan not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Operation plan approved successfully",
        "data":    plan,
    })
}

// GetApprovalStatus gets the approval status of an operation plan
// @Summary Get approval status
// @Description Get the approval status of an operation plan
// @Tags operation-plans
// @Produce json
// @Security BearerAuth
// @Param id path int true "Operation Plan ID"
// @Success 200 {array} models.OperationPlanApproval
// @Router /api/v1/operation-plans/{id}/approvals [get]
func (h *OperationPlanHandler) GetApprovalStatus(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation plan ID"})
        return
    }

    plan, err := h.service.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Operation plan not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data":             plan.Approvals,
    })
}

// GetPendingApprovals gets operation plans pending approval for current user's role
// @Summary Get pending approvals
// @Description Get operation plans pending approval for your role
// @Tags operation-plans
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.OperationPlanResponse
// @Router /api/v1/operation-plans/pending-approvals [get]
func (h *OperationPlanHandler) GetPendingApprovals(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
        return
    }
    currentUser := user.(*models.User)

    plans, err := h.service.GetByCreator(currentUser.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve pending approvals"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": plans})
}

// UploadGCodeFile uploads a G-code file
// @Summary Upload G-code file
// @Description Upload a G-code file (.txt only) to an operation plan
// @Tags operation-plans
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param id path int true "Operation Plan ID"
// @Param file formData file true "G-code file (.txt)"
// @Success 201 {object} models.GCodeFile
// @Failure 400 {object} map[string]string
// @Router /api/v1/operation-plans/{id}/g-codes [post]
func (h *OperationPlanHandler) UploadGCodeFile(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation plan ID"})
        return
    }

    // Get file from form
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
        return
    }

    // Get current user from context
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
        return
    }
    currentUser := user.(*models.User)

    gCodeFile, err := h.service.UploadGCodeFile(uint(id), file, currentUser.ID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "G-code file uploaded successfully",
        "data":    gCodeFile,
    })
}

// GetGCodeFiles gets all G-code files for an operation plan
// @Summary Get G-code files
// @Description Get all G-code files for an operation plan
// @Tags operation-plans
// @Produce json
// @Security BearerAuth
// @Param id path int true "Operation Plan ID"
// @Success 200 {array} models.GCodeFile
// @Router /api/v1/operation-plans/{id}/g-codes [get]
func (h *OperationPlanHandler) GetGCodeFiles(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation plan ID"})
        return
    }

    plan, err := h.service.GetByID(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Operation plan not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": plan.GCodeFiles})
}

// DownloadGCodeFile downloads a G-code file
// @Summary Download G-code file
// @Description Download a G-code file by ID
// @Tags g-codes
// @Produce octet-stream
// @Security BearerAuth
// @Param file_id path int true "G-code File ID"
// @Success 200 {file} file
// @Failure 404 {object} map[string]string
// @Router /api/v1/g-codes/{file_id}/download [get]
func (h *OperationPlanHandler) DownloadGCodeFile(c *gin.Context) {
    fileID, err := strconv.ParseUint(c.Param("file_id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
        return
    }

    file, err := h.service.GetGCodeFile(uint(fileID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
        return
    }

    // Check if file exists on filesystem
    if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
        c.JSON(http.StatusNotFound, gin.H{"error": "File not found on server"})
        return
    }

    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Disposition", "attachment; filename="+file.OriginalName)
    c.Header("Content-Type", "text/plain")
    c.File(file.FilePath)
}

// DeleteGCodeFile deletes a G-code file
// @Summary Delete G-code file
// @Description Delete a G-code file by ID
// @Tags g-codes
// @Produce json
// @Security BearerAuth
// @Param file_id path int true "G-code File ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/g-codes/{file_id} [delete]
func (h *OperationPlanHandler) DeleteGCodeFile(c *gin.Context) {
    fileID, err := strconv.ParseUint(c.Param("file_id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
        return
    }

    // Get current user from context
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
        return
    }
    currentUser := user.(*models.User)

    if err := h.service.DeleteGCodeFile(uint(fileID), currentUser.ID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "G-code file deleted successfully"})
}

// StartExecution starts the execution of an approved operation plan
// @Summary Start execution
// @Description Start the execution of an approved operation plan (Operator only)
// @Tags operation-plans
// @Produce json
// @Security BearerAuth
// @Param id path int true "Operation Plan ID"
// @Success 200 {object} models.OperationPlanResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/operation-plans/{id}/start [post]
func (h *OperationPlanHandler) StartExecution(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation plan ID"})
        return
    }

    plan, err := h.service.StartExecution(uint(id))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Get updated plan with details
    updatedPlan, _ := h.service.GetByID(plan.ID)

    c.JSON(http.StatusOK, gin.H{
        "message": "Execution started successfully",
        "data":    updatedPlan,
    })
}

// FinishExecution finishes the execution of an operation plan
// @Summary Finish execution
// @Description Finish the execution of an operation plan (Operator only)
// @Tags operation-plans
// @Produce json
// @Security BearerAuth
// @Param id path int true "Operation Plan ID"
// @Success 200 {object} models.OperationPlanResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/operation-plans/{id}/finish [post]
func (h *OperationPlanHandler) FinishExecution(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid operation plan ID"})
        return
    }

    plan, err := h.service.FinishExecution(uint(id))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Get updated plan with details
    updatedPlan, _ := h.service.GetByID(plan.ID)

    c.JSON(http.StatusOK, gin.H{
        "message": "Execution finished successfully",
        "data":    updatedPlan,
    })
}