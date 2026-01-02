package services

import (
	"errors"
	"fmt"
	"ganttpro-backend/models"
	"ganttpro-backend/repository"
)

type PPICLinkService struct {
	linkRepo     *repository.PPICLinkRepository
	scheduleRepo *repository.PPICScheduleRepository
}

func NewPPICLinkService(linkRepo *repository.PPICLinkRepository, scheduleRepo *repository.PPICScheduleRepository) *PPICLinkService {
	return &PPICLinkService{
		linkRepo:     linkRepo,
		scheduleRepo: scheduleRepo,
	}
}

// CreateLink creates a new PPIC link
func (s *PPICLinkService) CreateLink(req *models.CreatePPICLinkRequest) (*models.PPICLink, error) {
	// Validate that source and target are different
	if req.SourceScheduleID == req.TargetScheduleID {
		return nil, errors.New("source and target schedules must be different")
	}

	// Default link type to '0' (finish-to-start) if not specified
	if req.LinkType == "" {
		req.LinkType = "0"
	}

	// Fetch source and target schedules
	sourceSchedule, err := s.scheduleRepo.GetByID(req.SourceScheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch source schedule: %w", err)
	}
	if sourceSchedule == nil {
		return nil, errors.New("source schedule not found")
	}

	targetSchedule, err := s.scheduleRepo.GetByID(req.TargetScheduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch target schedule: %w", err)
	}
	if targetSchedule == nil {
		return nil, errors.New("target schedule not found")
	}

	// Validate that both schedules have the same machine
	if err := s.validateSameMachine(sourceSchedule, targetSchedule); err != nil {
		return nil, err
	}

	// Auto-reschedule target task if there's a date conflict (only for finish-to-start links)
	if req.LinkType == "0" {
		if err := s.autoRescheduleIfNeeded(sourceSchedule, targetSchedule); err != nil {
			return nil, fmt.Errorf("failed to auto-reschedule: %w", err)
		}
	}

	return s.linkRepo.Create(req)
}

// validateSameMachine checks if both schedules have at least one common machine
func (s *PPICLinkService) validateSameMachine(source, target *models.PPICSchedule) error {
	// Get machine IDs from source schedule
	sourceMachines := make(map[int64]bool)
	for _, ma := range source.MachineAssignments {
		sourceMachines[ma.MachineID] = true
	}

	// Check if target has any common machine
	hasCommonMachine := false
	for _, ma := range target.MachineAssignments {
		if sourceMachines[ma.MachineID] {
			hasCommonMachine = true
			break
		}
	}

	// If no machine assignments, both must have no machines
	if len(source.MachineAssignments) == 0 && len(target.MachineAssignments) == 0 {
		return errors.New("cannot link tasks without machine assignments")
	}

	if !hasCommonMachine {
		return errors.New("tasks can only be linked if they have the same machine")
	}

	return nil
}

// autoRescheduleIfNeeded reschedules the target task if it starts before source finishes
func (s *PPICLinkService) autoRescheduleIfNeeded(source, target *models.PPICSchedule) error {
	// For finish-to-start: target should start after source finishes
	// Check if target starts before or on the same day as source finishes
	if target.StartDate.Before(source.FinishDate) || target.StartDate.Equal(source.FinishDate) {
		// Calculate the duration of the target task
		duration := target.FinishDate.Sub(target.StartDate)

		// New start date is the day after source finishes
		newStartDate := source.FinishDate.AddDate(0, 0, 1)
		newFinishDate := newStartDate.Add(duration)

		// Update the target schedule
		updateReq := &models.UpdatePPICScheduleRequest{
			StartDate:  newStartDate.Format("2006-01-02"),
			FinishDate: newFinishDate.Format("2006-01-02"),
		}

		_, err := s.scheduleRepo.Update(target.ID, updateReq, &newStartDate, &newFinishDate)
		if err != nil {
			return fmt.Errorf("failed to update target schedule dates: %w", err)
		}
	}

	return nil
}

// GetAllLinks returns all PPIC links
func (s *PPICLinkService) GetAllLinks() ([]models.PPICLink, error) {
	return s.linkRepo.GetAll()
}

// DeleteLink deletes a PPIC link by ID
func (s *PPICLinkService) DeleteLink(id int64) error {
	return s.linkRepo.Delete(id)
}
