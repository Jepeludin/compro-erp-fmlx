package models

import "time"

// Priority constants
const (
	PriorityLow       = "Low"
	PriorityMedium    = "Medium"
	PriorityUrgent    = "Urgent"
	PriorityTopUrgent = "Top Urgent"
)

// Material status constants
const (
	MaterialReady    = "Ready"
	MaterialPending  = "Pending"
	MaterialOrdered  = "Ordered"
	MaterialNotReady = "Not Ready"
)

// Schedule status constants
const (
	ScheduleStatusPending    = "pending"
	ScheduleStatusInProgress = "in_progress"
	ScheduleStatusCompleted  = "completed"
	ScheduleStatusOnHold     = "on_hold"
)

// PPICSchedule represents a PPIC schedule entry for Gantt chart
type PPICSchedule struct {
	ID                 int64               `json:"id"`
	NJO                string              `json:"njo"`
	PartName           string              `json:"part_name"`
	Priority           string              `json:"priority"`
	PriorityAlpha      string              `json:"priority_alpha"`
	MaterialStatus     string              `json:"material_status"`
	Status             string              `json:"status"`
	Progress           int                 `json:"progress"`
	StartDate          time.Time           `json:"start_date"`
	FinishDate         time.Time           `json:"finish_date"`
	PPICNotes          string              `json:"ppic_notes"`
	CreatedBy          int64               `json:"created_by"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
	MachineAssignments []MachineAssignment `json:"machine_assignments"`
}

// MachineAssignment represents a machine assigned to a schedule
type MachineAssignment struct {
	ID             int64      `json:"id"`
	ScheduleID     int64      `json:"schedule_id"`
	MachineID      int64      `json:"machine_id"`
	MachineName    string     `json:"machine_name"`
	MachineCode    string     `json:"machine_code"`
	Sequence       int        `json:"sequence"`
	TargetHours    float64    `json:"target_hours"`
	ScheduledStart *time.Time `json:"scheduled_start"`
	ScheduledEnd   *time.Time `json:"scheduled_end"`
	ActualStart    *time.Time `json:"actual_start"`
	ActualEnd      *time.Time `json:"actual_end"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// Request DTOs

type CreatePPICScheduleRequest struct {
	NJO                string                           `json:"njo" binding:"required"`
	PartName           string                           `json:"part_name" binding:"required"`
	Priority           string                           `json:"priority" binding:"required"`
	PriorityAlpha      string                           `json:"priority_alpha"`
	MaterialStatus     string                           `json:"material_status" binding:"required"`
	StartDate          string                           `json:"start_date" binding:"required"`
	FinishDate         string                           `json:"finish_date" binding:"required"`
	PPICNotes          string                           `json:"ppic_notes"`
	MachineAssignments []CreateMachineAssignmentRequest `json:"machine_assignments" binding:"max=5"`
}

type CreateMachineAssignmentRequest struct {
	MachineID      int64   `json:"machine_id" binding:"required"`
	Sequence       int     `json:"sequence" binding:"required,min=1,max=5"`
	TargetHours    float64 `json:"target_hours"`
	ScheduledStart string  `json:"scheduled_start"`
	ScheduledEnd   string  `json:"scheduled_end"`
}

type UpdatePPICScheduleRequest struct {
	PartName           string                           `json:"part_name"`
	Priority           string                           `json:"priority"`
	PriorityAlpha      string                           `json:"priority_alpha"`
	MaterialStatus     string                           `json:"material_status"`
	Status             string                           `json:"status"`
	Progress           *int                             `json:"progress"`
	StartDate          string                           `json:"start_date"`
	FinishDate         string                           `json:"finish_date"`
	PPICNotes          string                           `json:"ppic_notes"`
	MachineAssignments []UpdateMachineAssignmentRequest `json:"machine_assignments"`
}

type UpdateMachineAssignmentRequest struct {
	ID          int64      `json:"id"`
	MachineID   int64      `json:"machine_id"`
	Sequence    int        `json:"sequence"`
	TargetHours float64    `json:"target_hours"`
	Status      string     `json:"status"`
	ActualStart *time.Time `json:"actual_start"`
	ActualEnd   *time.Time `json:"actual_end"`
}

// Gantt Chart DTOs

type GanttFilterRequest struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Priority  string `form:"priority"`
	Status    string `form:"status"`
	MachineID int64  `form:"machine_id"`
	GroupBy   string `form:"group_by"` // "priority", "machine", or empty for all
}

type GanttChartResponse struct {
	Sections []GanttSection      `json:"sections"`
	Machines []Machine           `json:"machines"`
	Links    []GanttLink         `json:"links"`
	Summary  GanttSummary        `json:"summary"`
	Filters  GanttFiltersApplied `json:"filters_applied"`
}

type GanttSection struct {
	SectionID   string      `json:"section_id"`
	SectionName string      `json:"section_name"`
	Tasks       []GanttTask `json:"tasks"`
}

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
	Color          string             `json:"color"`
	Machines       []GanttMachineInfo `json:"machines"`
}

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

type GanttSummary struct {
	TotalTasks       int `json:"total_tasks"`
	CompletedTasks   int `json:"completed_tasks"`
	InProgressTasks  int `json:"in_progress_tasks"`
	PendingTasks     int `json:"pending_tasks"`
	TopUrgentCount   int `json:"top_urgent_count"`
	UrgentCount      int `json:"urgent_count"`
	MediumCount      int `json:"medium_count"`
	LowCount         int `json:"low_count"`
	MaterialReady    int `json:"material_ready"`
	MaterialNotReady int `json:"material_not_ready"`
}

type GanttLink struct {
	ID     int64  `json:"id"`
	Source string `json:"source"` // Format: "task-{schedule_id}"
	Target string `json:"target"` // Format: "task-{schedule_id}"
	Type   string `json:"type"`   // Link type (0, 1, 2, 3)
}

type GanttFiltersApplied struct {
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Priority  string     `json:"priority"`
	Status    string     `json:"status"`
	MachineID *int64     `json:"machine_id"`
}

// Validation functions

func ValidatePriority(priority string) bool {
	validPriorities := []string{PriorityLow, PriorityMedium, PriorityUrgent, PriorityTopUrgent}
	for _, p := range validPriorities {
		if p == priority {
			return true
		}
	}
	return false
}

func ValidateMaterialStatus(status string) bool {
	validStatuses := []string{MaterialReady, MaterialPending, MaterialOrdered, MaterialNotReady}
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}

func GetPriorityColor(priority string) string {
	colors := map[string]string{
		PriorityTopUrgent: "#dc3545", // Red
		PriorityUrgent:    "#fd7e14", // Orange
		PriorityMedium:    "#ffc107", // Yellow
		PriorityLow:       "#28a745", // Green
	}
	if color, ok := colors[priority]; ok {
		return color
	}
	return "#6c757d" // Gray default
}

// PPICLink represents a dependency link between schedules
type PPICLink struct {
	ID               int64     `json:"id" gorm:"primaryKey"`
	SourceScheduleID int64     `json:"source_schedule_id" gorm:"not null"`
	TargetScheduleID int64     `json:"target_schedule_id" gorm:"not null"`
	LinkType         string    `json:"link_type" gorm:"size:20;default:'0'"` // 0=finish-to-start, 1=start-to-start, 2=finish-to-finish, 3=start-to-finish
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreatePPICLinkRequest struct {
	SourceScheduleID int64  `json:"source_schedule_id" binding:"required"`
	TargetScheduleID int64  `json:"target_schedule_id" binding:"required"`
	LinkType         string `json:"link_type"`
}
