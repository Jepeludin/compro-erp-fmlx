package models

import "time"

// Machine represents a production machine
type Machine struct {
	ID          int64      `json:"id"`
	MachineCode string     `json:"machine_code"`
	MachineName string     `json:"machine_name"`
	MachineType string     `json:"machine_type"`
	Location    string     `json:"location"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

// CreateMachineRequest for creating new machine
type CreateMachineRequest struct {
	MachineCode string `json:"machine_code" binding:"required"`
	MachineName string `json:"machine_name" binding:"required"`
	MachineType string `json:"machine_type"`
	Location    string `json:"location"`
	Status      string `json:"status"`
}

// UpdateMachineRequest for updating machine
type UpdateMachineRequest struct {
	MachineName string `json:"machine_name"`
	MachineType string `json:"machine_type"`
	Location    string `json:"location"`
	Status      string `json:"status"`
}
