package models

import "time"

// Priority levels for PPIC scheduling
const (
	PriorityLow       = "Low"
	PriorityMedium    = "Medium"
	PriorityUrgent    = "Urgent"
	PriorityTopUrgent = "Top Urgent"
)

// Material status constants
const (
	MaterialStatusReady    = "Ready"
	MaterialStatusPending  = "Pending"
	MaterialStatusOrdered  = "Ordered"
	MaterialStatusNotReady = "Not Ready"
)

// PPICSchedule represents a PPIC scheduling entry for Gantt chart
type PPICSchedule struct {
	ID             int64      `json:"id"`
	NJO            string     `json:"njo"`             // Order Number
	PartName       string     `json:"part_name"`       // Part name
	StartDate      time.Time  `json:"start_date"`      // Scheduled start
	FinishDate     time.Time  `json:"finish_date"`     // Scheduled finish
	Priority       string     `json:"priority"`        // Low, Medium, Urgent, Top Urgent
	PriorityAlpha  string     `json:"priority_alpha"`  // A, B, C, D, etc.
	MaterialStatus string     `json:"material_status"` // Ready, Pending, etc.
	PPICNotes      string     `json:"ppic_notes"`      // Notes from PPIC
	Status         string     `json:"status"`          // pending, in_progress, completed
	Progress       int        `json:"progress"`        // 0-100 percentage
	CreatedBy      int64      `json:"created_by"`      // User who created
	CreatedByName  string     `json:"created_by_name,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`

	// Machine assignments (loaded separately)
	MachineAssignments []MachineAssignment `json:"machine_assignments,omitempty"`
}

// MachineAssignment represents a machine assignment for a PPIC schedule
type MachineAssignment struct {
	ID             int64      `json:"id"`
	PPICScheduleID int64      `json:"ppic_schedule_id"`
	MachineID      int64      `json:"machine_id"`
	MachineName    string     `json:"machine_name,omitempty"`
	MachineCode    string     `json:"machine_code,omitempty"`
	TargetHours    float64    `json:"target_hours"`    // Duration in hours
	ScheduledStart *time.Time `json:"scheduled_start"` // Machine-specific start
	ScheduledEnd   *time.Time `json:"scheduled_end"`   // Machine-specific end
	ActualStart    *time.Time `json:"actual_start"`
	ActualEnd      *time.Time `json:"actual_end"`
	Status         string     `json:"status"`   // pending, in_progress, completed
	Sequence       int        `json:"sequence"` // Order of machine in process (1-5)
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// CreatePPICScheduleRequest for creating new PPIC schedule
type CreatePPICScheduleRequest struct {
	NJO                string                           `json:"njo" binding:"required"`
	PartName           string                           `json:"part_name" binding:"required"`
	StartDate          string                           `json:"start_date" binding:"required"` // Format: 2006-01-02
	FinishDate         string                           `json:"finish_date" binding:"required"`
	Priority           string                           `json:"priority" binding:"required"`
	PriorityAlpha      string                           `json:"priority_alpha"`
	MaterialStatus     string                           `json:"material_status" binding:"required"`
	PPICNotes          string                           `json:"ppic_notes"`
	MachineAssignments []CreateMachineAssignmentRequest `json:"machine_assignments" binding:"required,min=1,max=5"`
}

// CreateMachineAssignmentRequest for machine assignment in schedule
type CreateMachineAssignmentRequest struct {
	MachineID   int64   `json:"machine_id" binding:"required"`
	TargetHours float64 `json:"target_hours" binding:"required,gt=0"`
	Sequence    int     `json:"sequence" binding:"required,min=1,max=5"`
}

// UpdatePPICScheduleRequest for updating PPIC schedule
type UpdatePPICScheduleRequest struct {
	PartName           string                           `json:"part_name"`
	StartDate          string                           `json:"start_date"`
	FinishDate         string                           `json:"finish_date"`
	Priority           string                           `json:"priority"`
	PriorityAlpha      string                           `json:"priority_alpha"`
	MaterialStatus     string                           `json:"material_status"`
	PPICNotes          string                           `json:"ppic_notes"`
	Status             string                           `json:"status"`
	Progress           *int                             `json:"progress"`
	MachineAssignments []UpdateMachineAssignmentRequest `json:"machine_assignments"`
}

// UpdateMachineAssignmentRequest for updating machine assignment
type UpdateMachineAssignmentRequest struct {
	ID             int64      `json:"id"`
	MachineID      int64      `json:"machine_id"`
	TargetHours    float64    `json:"target_hours"`
	ScheduledStart *time.Time `json:"scheduled_start"`
	ScheduledEnd   *time.Time `json:"scheduled_end"`
	ActualStart    *time.Time `json:"actual_start"`
	ActualEnd      *time.Time `json:"actual_end"`
	Status         string     `json:"status"`
	Sequence       int        `json:"sequence"`
}

// GanttChartResponse represents the response for Gantt chart display
type GanttChartResponse struct {
	Sections []GanttSection      `json:"sections"`
	Machines []MachineInfo       `json:"machines"`
	Summary  GanttSummary        `json:"summary"`
	Filters  GanttFiltersApplied `json:"filters_applied"`
}

// GanttSection represents a section/group in Gantt chart
type GanttSection struct {
	SectionID   string      `json:"section_id"`
	SectionName string      `json:"section_name"`
	Tasks       []GanttTask `json:"tasks"`
}

// GanttTask represents a single task bar in Gantt chart
type GanttTask struct {
	TaskID         string             `json:"task_id"`
	TaskName       string             `json:"task_name"`
	NJO            string             `json:"njo"`
	PartName       string             `json:"part_name"`
	Start          time.Time          `json:"start"`
	End            time.Time          `json:"end"`
	Priority       string             `json:"priority"`
	PriorityAlpha  string             `json:"priority_alpha"`
	MaterialStatus string             `json:"material_status"`
	Status         string             `json:"status"`
	Progress       int                `json:"progress"`
	PPICNotes      string             `json:"ppic_notes"`
	Machines       []GanttMachineInfo `json:"machines"`
	Color          string             `json:"color"` // Color based on priority
}

// GanttMachineInfo represents machine info in Gantt task
type GanttMachineInfo struct {
	MachineID      int64      `json:"machine_id"`
	MachineName    string     `json:"machine_name"`
	MachineCode    string     `json:"machine_code"`
	DurationHours  float64    `json:"duration_hours"`
	ScheduledStart *time.Time `json:"scheduled_start"`
	ScheduledEnd   *time.Time `json:"scheduled_end"`
	Status         string     `json:"status"`
	Sequence       int        `json:"sequence"`
}

// MachineInfo for listing available machines
type MachineInfo struct {
	ID          int64  `json:"id"`
	MachineCode string `json:"machine_code"`
	MachineName string `json:"machine_name"`
	Status      string `json:"status"`
}

// GanttSummary for summary statistics
type GanttSummary struct {
	TotalJobs     int `json:"total_jobs"`
	Pending       int `json:"pending"`
	InProgress    int `json:"in_progress"`
	Completed     int `json:"completed"`
	TopUrgent     int `json:"top_urgent"`
	Urgent        int `json:"urgent"`
	Medium        int `json:"medium"`
	Low           int `json:"low"`
	MaterialReady int `json:"material_ready"`
}

// GanttFiltersApplied shows which filters were applied
type GanttFiltersApplied struct {
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	MachineID *int64     `json:"machine_id,omitempty"`
	Priority  string     `json:"priority,omitempty"`
	Status    string     `json:"status,omitempty"`
}

// GanttFilterRequest for filtering Gantt chart data
type GanttFilterRequest struct {
	StartDate string `form:"start_date"` // Format: 2006-01-02
	EndDate   string `form:"end_date"`
	MachineID int64  `form:"machine_id"`
	Priority  string `form:"priority"`
	Status    string `form:"status"`
	GroupBy   string `form:"group_by"` // project, machine, priority
}

// GetPriorityColor returns color based on priority
func GetPriorityColor(priority string) string {
	switch priority {
	case PriorityTopUrgent:
		return "#dc2626" // Red
	case PriorityUrgent:
		return "#ea580c" // Orange
	case PriorityMedium:
		return "#ca8a04" // Yellow
	case PriorityLow:
		return "#16a34a" // Green
	default:
		return "#6b7280" // Gray
	}
}

// ValidatePriority checks if priority is valid
func ValidatePriority(priority string) bool {
	return priority == PriorityLow ||
		priority == PriorityMedium ||
		priority == PriorityUrgent ||
		priority == PriorityTopUrgent
}

// ValidateMaterialStatus checks if material status is valid
func ValidateMaterialStatus(status string) bool {
	return status == MaterialStatusReady ||
		status == MaterialStatusPending ||
		status == MaterialStatusOrdered ||
		status == MaterialStatusNotReady
}


