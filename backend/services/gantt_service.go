package services

import (
	"errors"
	"fmt"
	"ganttpro-backend/models"
	"ganttpro-backend/repository"
	"time"
)

type GanttService struct {
	ppicRepo     *repository.PPICScheduleRepository
	ppicLinkRepo *repository.PPICLinkRepository
}

func NewGanttService(ppicRepo *repository.PPICScheduleRepository, ppicLinkRepo *repository.PPICLinkRepository) *GanttService {
	return &GanttService{
		ppicRepo:     ppicRepo,
		ppicLinkRepo: ppicLinkRepo,
	}
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

	// Validate machine assignments count (0-5)
	if len(req.MachineAssignments) > 5 {
		return nil, errors.New("cannot have more than 5 machine assignments")
	}

	// Validate unique sequences if there are machine assignments
	if len(req.MachineAssignments) > 0 {
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

	// Validate that the new dates don't conflict with predecessor tasks (if this task is a target of any links)
	if startDate != nil {
		if err := s.validateNoPredecessorConflict(id, *startDate); err != nil {
			return nil, err
		}
	}

	// Update the schedule
	updated, err := s.ppicRepo.Update(id, req, startDate, finishDate)
	if err != nil {
		return nil, err
	}

	// If dates were updated, cascade the changes to dependent tasks
	if startDate != nil || finishDate != nil {
		// Get the updated schedule to use for cascading
		scheduleAfterUpdate, err := s.ppicRepo.GetByID(id)
		if err != nil {
			return updated, nil // Return the update result even if cascade fails
		}

		// Cascade reschedule to all dependent tasks
		if err := s.cascadeReschedule(scheduleAfterUpdate); err != nil {
			// Log error but don't fail the update
			fmt.Printf("Warning: Failed to cascade reschedule: %v\n", err)
		}
	}

	return updated, nil
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

	// Get all PPIC links
	ppicLinks, err := s.ppicLinkRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Build response
	response := &models.GanttChartResponse{
		Machines: machines,
		Links:    s.convertToGanttLinks(ppicLinks),
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

// validateNoPredecessorConflict checks if the new start date conflicts with any predecessor tasks
func (s *GanttService) validateNoPredecessorConflict(scheduleID int64, newStartDate time.Time) error {
	// Get all links where this schedule is the target (predecessor links)
	predecessorLinks, err := s.ppicLinkRepo.GetByTargetScheduleID(scheduleID)
	if err != nil {
		return fmt.Errorf("failed to get predecessor links: %w", err)
	}

	// If no predecessors, no conflict possible
	if len(predecessorLinks) == 0 {
		return nil
	}

	// Check each predecessor link
	for _, link := range predecessorLinks {
		// Only validate finish-to-start links
		if link.LinkType != "0" {
			continue
		}

		// Get the source (predecessor) schedule
		sourceSchedule, err := s.ppicRepo.GetByID(link.SourceScheduleID)
		if err != nil {
			return fmt.Errorf("failed to get predecessor schedule %d: %w", link.SourceScheduleID, err)
		}
		if sourceSchedule == nil {
			continue
		}

		// For finish-to-start: target must start after source finishes
		// newStartDate must be > sourceSchedule.FinishDate
		if newStartDate.Before(sourceSchedule.FinishDate) || newStartDate.Equal(sourceSchedule.FinishDate) {
			return fmt.Errorf("task tidak bisa dimajuin ke tanggal %s karena terhubung dengan task '%s' yang selesai di tanggal %s. Task ini harus mulai minimal tanggal %s",
				newStartDate.Format("2006-01-02"),
				sourceSchedule.PartName,
				sourceSchedule.FinishDate.Format("2006-01-02"),
				sourceSchedule.FinishDate.AddDate(0, 0, 1).Format("2006-01-02"))
		}
	}

	return nil
}

// cascadeReschedule recursively reschedules all dependent tasks when a source task is updated
func (s *GanttService) cascadeReschedule(sourceSchedule *models.PPICSchedule) error {
	// Get all links where this schedule is the source (tasks that depend on this one)
	dependentLinks, err := s.ppicLinkRepo.GetBySourceScheduleID(sourceSchedule.ID)
	if err != nil {
		return fmt.Errorf("failed to get dependent links: %w", err)
	}

	// If no dependent tasks, nothing to do
	if len(dependentLinks) == 0 {
		return nil
	}

	// For each dependent task, reschedule it
	for _, link := range dependentLinks {
		// Only handle finish-to-start links for now
		if link.LinkType != "0" {
			continue
		}

		// Get the target schedule
		targetSchedule, err := s.ppicRepo.GetByID(link.TargetScheduleID)
		if err != nil {
			return fmt.Errorf("failed to get target schedule %d: %w", link.TargetScheduleID, err)
		}
		if targetSchedule == nil {
			continue
		}

		// Calculate new dates for the target task
		// For finish-to-start: target should start after source finishes
		newStartDate := sourceSchedule.FinishDate.AddDate(0, 0, 1) // Day after source finishes

		// Calculate task duration to maintain it
		duration := targetSchedule.FinishDate.Sub(targetSchedule.StartDate)
		newFinishDate := newStartDate.Add(duration)

		// Update the target schedule
		updateReq := &models.UpdatePPICScheduleRequest{
			StartDate:  newStartDate.Format("2006-01-02"),
			FinishDate: newFinishDate.Format("2006-01-02"),
		}

		_, err = s.ppicRepo.Update(targetSchedule.ID, updateReq, &newStartDate, &newFinishDate)
		if err != nil {
			return fmt.Errorf("failed to update dependent schedule %d: %w", targetSchedule.ID, err)
		}

		// Get the updated target schedule
		updatedTarget, err := s.ppicRepo.GetByID(targetSchedule.ID)
		if err != nil {
			return fmt.Errorf("failed to get updated target schedule %d: %w", targetSchedule.ID, err)
		}

		// Recursively cascade to tasks that depend on this target
		if err := s.cascadeReschedule(updatedTarget); err != nil {
			return fmt.Errorf("failed to cascade reschedule for schedule %d: %w", updatedTarget.ID, err)
		}
	}

	return nil
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

func (s *GanttService) convertToGanttLinks(ppicLinks []models.PPICLink) []models.GanttLink {
	var ganttLinks []models.GanttLink

	for _, link := range ppicLinks {
		ganttLink := models.GanttLink{
			ID:     link.ID,
			Source: fmt.Sprintf("task-%d", link.SourceScheduleID),
			Target: fmt.Sprintf("task-%d", link.TargetScheduleID),
			Type:   link.LinkType,
		}
		ganttLinks = append(ganttLinks, ganttLink)
	}

	return ganttLinks
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
