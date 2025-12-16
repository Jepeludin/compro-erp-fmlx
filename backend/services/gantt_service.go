package services

import (
	"errors"
	"fmt"
	"ganttpro-backend/models"
	"ganttpro-backend/repository"
	"time"
)

type GanttService struct {
	ppicRepo *repository.PPICScheduleRepository
}

func NewGanttService(ppicRepo *repository.PPICScheduleRepository) *GanttService {
	return &GanttService{ppicRepo: ppicRepo}
}

// CreatePPICSchedule creates a new PPIC schedule entry
func (s *GanttService) CreatePPICSchedule(req *models.CreatePPICScheduleRequest, createdBy int64) (*models.PPICSchedule, error) {
	// Validate priority
	if !models.ValidatePriority(req.Priority) {
		return nil, errors.New("invalid priority value. Must be: Low, Medium, Urgent, or Top Urgent")
	}

	// Validate material status
	if !models.ValidateMaterialStatus(req.MaterialStatus) {
		return nil, errors.New("invalid material status. Must be: Ready, Pending, Ordered, or Not Ready")
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format. Use YYYY-MM-DD: %v", err)
	}

	finishDate, err := time.Parse("2006-01-02", req.FinishDate)
	if err != nil {
		return nil, fmt.Errorf("invalid finish_date format. Use YYYY-MM-DD: %v", err)
	}

	// Validate date range
	if finishDate.Before(startDate) {
		return nil, errors.New("finish_date must be after start_date")
	}

	// Check if NJO already exists
	existing, err := s.ppicRepo.GetByNJO(req.NJO)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("NJO already exists in PPIC schedule")
	}

	// Validate machine assignments count (1-5)
	if len(req.MachineAssignments) < 1 || len(req.MachineAssignments) > 5 {
		return nil, errors.New("must have between 1 and 5 machine assignments")
	}

	// Validate unique sequences
	sequences := make(map[int]bool)
	for _, ma := range req.MachineAssignments {
		if ma.Sequence < 1 || ma.Sequence > 5 {
			return nil, errors.New("machine sequence must be between 1 and 5")
		}
		if sequences[ma.Sequence] {
			return nil, errors.New("duplicate sequence numbers found")
		}
		sequences[ma.Sequence] = true
	}

	return s.ppicRepo.Create(req, createdBy, startDate, finishDate)
}

// UpdatePPICSchedule updates an existing PPIC schedule
func (s *GanttService) UpdatePPICSchedule(id int64, req *models.UpdatePPICScheduleRequest) (*models.PPICSchedule, error) {
	// Check if exists
	existing, err := s.ppicRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("PPIC schedule not found")
	}

	// Validate priority if provided
	if req.Priority != "" && !models.ValidatePriority(req.Priority) {
		return nil, errors.New("invalid priority value. Must be: Low, Medium, Urgent, or Top Urgent")
	}

	// Validate material status if provided
	if req.MaterialStatus != "" && !models.ValidateMaterialStatus(req.MaterialStatus) {
		return nil, errors.New("invalid material status. Must be: Ready, Pending, Ordered, or Not Ready")
	}

	// Parse dates if provided
	var startDate, finishDate *time.Time

	if req.StartDate != "" {
		sd, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, fmt.Errorf("invalid start_date format. Use YYYY-MM-DD: %v", err)
		}
		startDate = &sd
	}

	if req.FinishDate != "" {
		fd, err := time.Parse("2006-01-02", req.FinishDate)
		if err != nil {
			return nil, fmt.Errorf("invalid finish_date format. Use YYYY-MM-DD: %v", err)
		}
		finishDate = &fd
	}

	// Validate date range if both provided
	if startDate != nil && finishDate != nil && finishDate.Before(*startDate) {
		return nil, errors.New("finish_date must be after start_date")
	}

	// Validate progress if provided
	if req.Progress != nil && (*req.Progress < 0 || *req.Progress > 100) {
		return nil, errors.New("progress must be between 0 and 100")
	}

	return s.ppicRepo.Update(id, req, startDate, finishDate)
}

// DeletePPICSchedule deletes a PPIC schedule
func (s *GanttService) DeletePPICSchedule(id int64) error {
	// Check if exists
	existing, err := s.ppicRepo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("PPIC schedule not found")
	}

	return s.ppicRepo.Delete(id)
}

// GetPPICSchedule gets a single PPIC schedule by ID
func (s *GanttService) GetPPICSchedule(id int64) (*models.PPICSchedule, error) {
	schedule, err := s.ppicRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if schedule == nil {
		return nil, errors.New("PPIC schedule not found")
	}
	return schedule, nil
}

// GetAllPPICSchedules gets all PPIC schedules
func (s *GanttService) GetAllPPICSchedules() ([]models.PPICSchedule, error) {
	return s.ppicRepo.GetAll()
}

// GetGanttChartData returns formatted data for Gantt chart display
func (s *GanttService) GetGanttChartData(filter models.GanttFilterRequest) (*models.GanttChartResponse, error) {
	// Get filtered schedules
	schedules, err := s.ppicRepo.GetWithFilters(filter)
	if err != nil {
		return nil, err
	}

	// Get all machines
	machines, err := s.ppicRepo.GetAllMachines()
	if err != nil {
		return nil, err
	}

	// Get summary
	summary, err := s.ppicRepo.GetSummary()
	if err != nil {
		return nil, err
	}

	// Build response
	response := &models.GanttChartResponse{
		Machines: machines,
		Summary:  *summary,
		Filters:  s.buildFiltersApplied(filter),
	}

	// Group tasks based on groupBy parameter
	switch filter.GroupBy {
	case "priority":
		response.Sections = s.groupByPriority(schedules)
	case "machine":
		response.Sections = s.groupByMachine(schedules)
	default:
		// Default: group by all (single section)
		response.Sections = s.groupAll(schedules)
	}

	return response, nil
}

// GetSchedulesByMachine gets all schedules for a specific machine
func (s *GanttService) GetSchedulesByMachine(machineID int64) ([]models.PPICSchedule, error) {
	return s.ppicRepo.GetSchedulesByMachine(machineID)
}

// AddMachineAssignment adds a machine to an existing schedule
func (s *GanttService) AddMachineAssignment(scheduleID int64, req *models.CreateMachineAssignmentRequest) (*models.MachineAssignment, error) {
	// Check if schedule exists
	schedule, err := s.ppicRepo.GetByID(scheduleID)
	if err != nil {
		return nil, err
	}
	if schedule == nil {
		return nil, errors.New("PPIC schedule not found")
	}

	// Check max 5 machines
	if len(schedule.MachineAssignments) >= 5 {
		return nil, errors.New("maximum 5 machines allowed per schedule")
	}

	// Check sequence uniqueness
	for _, ma := range schedule.MachineAssignments {
		if ma.Sequence == req.Sequence {
			return nil, errors.New("sequence number already exists")
		}
	}

	return s.ppicRepo.AddMachineAssignment(req, scheduleID)
}

// RemoveMachineAssignment removes a machine from a schedule
func (s *GanttService) RemoveMachineAssignment(assignmentID int64) error {
	return s.ppicRepo.DeleteMachineAssignment(assignmentID)
}

// Helper functions

func (s *GanttService) buildFiltersApplied(filter models.GanttFilterRequest) models.GanttFiltersApplied {
	applied := models.GanttFiltersApplied{
		Priority: filter.Priority,
		Status:   filter.Status,
	}

	if filter.StartDate != "" {
		t, _ := time.Parse("2006-01-02", filter.StartDate)
		applied.StartDate = &t
	}
	if filter.EndDate != "" {
		t, _ := time.Parse("2006-01-02", filter.EndDate)
		applied.EndDate = &t
	}
	if filter.MachineID > 0 {
		applied.MachineID = &filter.MachineID
	}

	return applied
}

func (s *GanttService) groupAll(schedules []models.PPICSchedule) []models.GanttSection {
	section := models.GanttSection{
		SectionID:   "all",
		SectionName: "All Tasks",
		Tasks:       s.convertToGanttTasks(schedules),
	}
	return []models.GanttSection{section}
}

func (s *GanttService) groupByPriority(schedules []models.PPICSchedule) []models.GanttSection {
	priorityMap := map[string][]models.PPICSchedule{
		models.PriorityTopUrgent: {},
		models.PriorityUrgent:    {},
		models.PriorityMedium:    {},
		models.PriorityLow:       {},
	}

	for _, schedule := range schedules {
		priorityMap[schedule.Priority] = append(priorityMap[schedule.Priority], schedule)
	}

	var sections []models.GanttSection
	priorities := []string{models.PriorityTopUrgent, models.PriorityUrgent, models.PriorityMedium, models.PriorityLow}

	for _, priority := range priorities {
		if len(priorityMap[priority]) > 0 {
			sections = append(sections, models.GanttSection{
				SectionID:   fmt.Sprintf("priority-%s", priority),
				SectionName: priority,
				Tasks:       s.convertToGanttTasks(priorityMap[priority]),
			})
		}
	}

	return sections
}

func (s *GanttService) groupByMachine(schedules []models.PPICSchedule) []models.GanttSection {
	// Get all machines first
	machines, _ := s.ppicRepo.GetAllMachines()
	machineMap := make(map[int64]string)
	for _, m := range machines {
		machineMap[m.ID] = m.MachineName
	}

	// Group schedules by machine
	machineSchedules := make(map[int64][]models.PPICSchedule)

	for _, schedule := range schedules {
		for _, ma := range schedule.MachineAssignments {
			machineSchedules[ma.MachineID] = append(machineSchedules[ma.MachineID], schedule)
		}
	}

	var sections []models.GanttSection
	for machineID, scheds := range machineSchedules {
		machineName := machineMap[machineID]
		if machineName == "" {
			machineName = fmt.Sprintf("Machine %d", machineID)
		}
		sections = append(sections, models.GanttSection{
			SectionID:   fmt.Sprintf("machine-%d", machineID),
			SectionName: machineName,
			Tasks:       s.convertToGanttTasks(scheds),
		})
	}

	return sections
}

func (s *GanttService) convertToGanttTasks(schedules []models.PPICSchedule) []models.GanttTask {
	var tasks []models.GanttTask

	for _, schedule := range schedules {
		task := models.GanttTask{
			TaskID:         fmt.Sprintf("task-%d", schedule.ID),
			TaskName:       schedule.PartName,
			NJO:            schedule.NJO,
			PartName:       schedule.PartName,
			Start:          schedule.StartDate,
			End:            schedule.FinishDate,
			Priority:       schedule.Priority,
			PriorityAlpha:  schedule.PriorityAlpha,
			MaterialStatus: schedule.MaterialStatus,
			Status:         schedule.Status,
			Progress:       schedule.Progress,
			PPICNotes:      schedule.PPICNotes,
			Color:          models.GetPriorityColor(schedule.Priority),
			Machines:       s.convertToGanttMachines(schedule.MachineAssignments),
		}
		tasks = append(tasks, task)
	}

	return tasks
}

func (s *GanttService) convertToGanttMachines(assignments []models.MachineAssignment) []models.GanttMachineInfo {
	var machines []models.GanttMachineInfo

	for _, ma := range assignments {
		machine := models.GanttMachineInfo{
			MachineID:      ma.MachineID,
			MachineName:    ma.MachineName,
			MachineCode:    ma.MachineCode,
			DurationHours:  ma.TargetHours,
			ScheduledStart: ma.ScheduledStart,
			ScheduledEnd:   ma.ScheduledEnd,
			Status:         ma.Status,
			Sequence:       ma.Sequence,
		}
		machines = append(machines, machine)
	}

	return machines
}

// UpdateMachineAssignmentStatus updates the status of a machine assignment
func (s *GanttService) UpdateMachineAssignmentStatus(scheduleID int64, assignmentID int64, status string, actualStart, actualEnd *time.Time) error {
	schedule, err := s.ppicRepo.GetByID(scheduleID)
	if err != nil {
		return err
	}
	if schedule == nil {
		return errors.New("PPIC schedule not found")
	}

	// Find the assignment
	found := false
	for _, ma := range schedule.MachineAssignments {
		if ma.ID == assignmentID {
			found = true
			break
		}
	}
	if !found {
		return errors.New("machine assignment not found")
	}

	// Update via repository
	updateReq := &models.UpdatePPICScheduleRequest{
		MachineAssignments: []models.UpdateMachineAssignmentRequest{
			{
				ID:          assignmentID,
				Status:      status,
				ActualStart: actualStart,
				ActualEnd:   actualEnd,
			},
		},
	}

	_, err = s.ppicRepo.Update(scheduleID, updateReq, nil, nil)
	return err
}
