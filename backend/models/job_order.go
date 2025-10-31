package models

import "time"

// JobOrder represents a production job order
type JobOrder struct {
	ID           int64      `json:"id"`
	MachineID    int64      `json:"machine_id"`
	MachineName  string     `json:"machine_name,omitempty"` // For JOIN queries
	NJO          string     `json:"njo"`
	Project      string     `json:"project"`
	Item         string     `json:"item"`
	Note         string     `json:"note"`
	Deadline     string     `json:"deadline"`
	OperatorID   *int64     `json:"operator_id,omitempty"`
	OperatorName string     `json:"operator_name,omitempty"` // For JOIN queries
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`

	// For detailed view with stages
	ProcessStages []ProcessStage `json:"process_stages,omitempty"`
}

// ProcessStage represents a stage in the production process
type ProcessStage struct {
	ID              int64      `json:"id"`
	JobOrderID      int64      `json:"job_order_id"`
	StageName       string     `json:"stage_name"` // 'setting', 'proses', 'cmm', 'kalibrasi'
	StartTime       *time.Time `json:"start_time,omitempty"`
	FinishTime      *time.Time `json:"finish_time,omitempty"`
	DurationMinutes *float64   `json:"duration_minutes,omitempty"` // Auto-calculated
	OperatorID      *int64     `json:"operator_id,omitempty"`
	OperatorName    string     `json:"operator_name,omitempty"` // For JOIN queries
	Notes           string     `json:"notes"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// CreateJobOrderRequest for creating new job order
type CreateJobOrderRequest struct {
	MachineID  int64  `json:"machine_id" binding:"required"`
	NJO        string `json:"njo" binding:"required"`
	Project    string `json:"project"`
	Item       string `json:"item"`
	Note       string `json:"note"`
	Deadline   string `json:"deadline"`
	OperatorID *int64 `json:"operator_id"`
}

// UpdateJobOrderRequest for updating job order
type UpdateJobOrderRequest struct {
	Project    string `json:"project"`
	Item       string `json:"item"`
	Note       string `json:"note"`
	Deadline   string `json:"deadline"`
	OperatorID *int64 `json:"operator_id"`
	Status     string `json:"status"`
}

// UpdateProcessStageRequest for updating process stage
type UpdateProcessStageRequest struct {
	StartTime  *time.Time `json:"start_time"`
	FinishTime *time.Time `json:"finish_time"`
	OperatorID *int64     `json:"operator_id"`
	Notes      string     `json:"notes"`
}

// JobOrderWithStages for detailed view
type JobOrderWithStages struct {
	JobOrder
	Stages []ProcessStage `json:"stages"`
}
