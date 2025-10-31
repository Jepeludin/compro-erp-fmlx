package handlers

import (
	"ganttpro-backend/models"
	"ganttpro-backend/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MachineHandler struct {
	repo *repository.MachineRepository
}

func NewMachineHandler(repo *repository.MachineRepository) *MachineHandler {
	return &MachineHandler{repo: repo}
}

// GetAllMachines godoc
// @Summary Get all machines
// @Description Get all active machines
// @Tags machines
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/machines [get]
func (h *MachineHandler) GetAllMachines(c *gin.Context) {
	machines, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch machines"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"machines": machines,
		"count":    len(machines),
	})
}

// GetMachine godoc
// @Summary Get machine by ID
// @Description Get a specific machine by ID
// @Tags machines
// @Produce json
// @Param id path int true "Machine ID"
// @Success 200 {object} models.Machine
// @Router /api/v1/machines/{id} [get]
func (h *MachineHandler) GetMachine(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine ID"})
		return
	}

	machine, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch machine"})
		return
	}

	if machine == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Machine not found"})
		return
	}

	c.JSON(http.StatusOK, machine)
}

// CreateMachine godoc
// @Summary Create new machine
// @Description Create a new machine
// @Tags machines
// @Accept json
// @Produce json
// @Param machine body models.CreateMachineRequest true "Machine data"
// @Success 201 {object} models.Machine
// @Router /api/v1/machines [post]
func (h *MachineHandler) CreateMachine(c *gin.Context) {
	var req models.CreateMachineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if machine code already exists
	existing, err := h.repo.GetByCode(req.MachineCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check machine code"})
		return
	}
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Machine code already exists"})
		return
	}

	machine, err := h.repo.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create machine"})
		return
	}

	c.JSON(http.StatusCreated, machine)
}

// UpdateMachine godoc
// @Summary Update machine
// @Description Update an existing machine
// @Tags machines
// @Accept json
// @Produce json
// @Param id path int true "Machine ID"
// @Param machine body models.UpdateMachineRequest true "Machine data"
// @Success 200 {object} models.Machine
// @Router /api/v1/machines/{id} [put]
func (h *MachineHandler) UpdateMachine(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine ID"})
		return
	}

	var req models.UpdateMachineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	machine, err := h.repo.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update machine"})
		return
	}

	if machine == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Machine not found"})
		return
	}

	c.JSON(http.StatusOK, machine)
}

// DeleteMachine godoc
// @Summary Delete machine
// @Description Soft delete a machine
// @Tags machines
// @Produce json
// @Param id path int true "Machine ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/machines/{id} [delete]
func (h *MachineHandler) DeleteMachine(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid machine ID"})
		return
	}

	err = h.repo.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete machine"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Machine deleted successfully"})
}
