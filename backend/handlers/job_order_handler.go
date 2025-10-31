package handlers

import (
	"ganttpro-backend/models"
	"ganttpro-backend/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type JobOrderHandler struct {
	repo *repository.JobOrderRepository
}

func NewJobOrderHandler(repo *repository.JobOrderRepository) *JobOrderHandler {
	return &JobOrderHandler{repo: repo}
}

// GetAllJobOrders godoc
// @Summary Get all job orders
// @Description Get all job orders with machine and operator info
// @Tags job_orders
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/job-orders [get]
func (h *JobOrderHandler) GetAllJobOrders(c *gin.Context) {
	jobs, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"job_orders": jobs,
		"count":      len(jobs),
	})
}

// GetJobOrder godoc
// @Summary Get job order by ID
// @Description Get a specific job order with process stages
// @Tags job_orders
// @Produce json
// @Param id path int true "Job Order ID"
// @Success 200 {object} models.JobOrder
// @Router /api/v1/job-orders/{id} [get]
func (h *JobOrderHandler) GetJobOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job order ID"})
		return
	}

	job, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job order"})
		return
	}

	if job == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job order not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

// GetJobOrdersByMachine godoc
// @Summary Get job orders by machine ID
// @Description Get all job orders for a specific machine
// @Tags job_orders
// @Produce json
// @Param machine_id path int true "Machine ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/machines/{machine_id}/job-orders [get]
func (h *JobOrderHandler) GetJobOrdersByMachine(c *gin.Context) {
	machineID, err := strconv.ParseInt(c.Param("machine_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine ID"})
		return
	}

	jobs, err := h.repo.GetByMachineID(machineID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"job_orders": jobs,
		"count":      len(jobs),
	})
}

// CreateJobOrder godoc
// @Summary Create new job order
// @Description Create a new job order with default process stages
// @Tags job_orders
// @Accept json
// @Produce json
// @Param job_order body models.CreateJobOrderRequest true "Job Order data"
// @Success 201 {object} models.JobOrder
// @Router /api/v1/job-orders [post]
func (h *JobOrderHandler) CreateJobOrder(c *gin.Context) {
	var req models.CreateJobOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := h.repo.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job order"})
		return
	}

	c.JSON(http.StatusCreated, job)
}

// UpdateJobOrder godoc
// @Summary Update job order
// @Description Update an existing job order
// @Tags job_orders
// @Accept json
// @Produce json
// @Param id path int true "Job Order ID"
// @Param job_order body models.UpdateJobOrderRequest true "Job Order data"
// @Success 200 {object} models.JobOrder
// @Router /api/v1/job-orders/{id} [put]
func (h *JobOrderHandler) UpdateJobOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job order ID"})
		return
	}

	var req models.UpdateJobOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := h.repo.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update job order"})
		return
	}

	if job == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job order not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

// DeleteJobOrder godoc
// @Summary Delete job order
// @Description Soft delete a job order
// @Tags job_orders
// @Produce json
// @Param id path int true "Job Order ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/job-orders/{id} [delete]
func (h *JobOrderHandler) DeleteJobOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job order ID"})
		return
	}

	err = h.repo.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete job order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Job order deleted successfully"})
}

// UpdateProcessStage godoc
// @Summary Update process stage
// @Description Update a process stage (start_time, finish_time, operator, notes)
// @Tags process_stages
// @Accept json
// @Produce json
// @Param id path int true "Process Stage ID"
// @Param stage body models.UpdateProcessStageRequest true "Process Stage data"
// @Success 200 {object} models.ProcessStage
// @Router /api/v1/process-stages/{id} [put]
func (h *JobOrderHandler) UpdateProcessStage(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid process stage ID"})
		return
	}

	var req models.UpdateProcessStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stage, err := h.repo.UpdateProcessStage(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update process stage"})
		return
	}

	if stage == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Process stage not found"})
		return
	}

	c.JSON(http.StatusOK, stage)
}
