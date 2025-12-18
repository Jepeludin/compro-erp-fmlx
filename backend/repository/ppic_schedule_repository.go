package repository

import (
	"database/sql"
	"fmt"
	"ganttpro-backend/models"
	"time"
)

type PPICScheduleRepository struct {
	db *sql.DB
}

func NewPPICScheduleRepository(db *sql.DB) *PPICScheduleRepository {
	return &PPICScheduleRepository{db: db}
}

// Create creates a new PPIC schedule with machine assignments
func (r *PPICScheduleRepository) Create(req *models.CreatePPICScheduleRequest, createdBy int64, startDate, finishDate time.Time) (*models.PPICSchedule, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert schedule
	query := `
		INSERT INTO ppic_schedules (njo, part_name, priority, priority_alpha, material_status, start_date, finish_date, ppic_notes, created_by, status, progress, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 'pending', 0, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	var schedule models.PPICSchedule
	err = tx.QueryRow(query, req.NJO, req.PartName, req.Priority, req.PriorityAlpha, req.MaterialStatus, startDate, finishDate, req.PPICNotes, createdBy).
		Scan(&schedule.ID, &schedule.CreatedAt, &schedule.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	// Set basic fields
	schedule.NJO = req.NJO
	schedule.PartName = req.PartName
	schedule.Priority = req.Priority
	schedule.PriorityAlpha = req.PriorityAlpha
	schedule.MaterialStatus = req.MaterialStatus
	schedule.StartDate = startDate
	schedule.FinishDate = finishDate
	schedule.PPICNotes = req.PPICNotes
	schedule.CreatedBy = createdBy
	schedule.Status = "pending"
	schedule.Progress = 0

	// Insert machine assignments
	for _, ma := range req.MachineAssignments {
		assignment, err := r.createMachineAssignment(tx, schedule.ID, &ma)
		if err != nil {
			return nil, err
		}
		schedule.MachineAssignments = append(schedule.MachineAssignments, *assignment)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &schedule, nil
}

func (r *PPICScheduleRepository) createMachineAssignment(tx *sql.Tx, scheduleID int64, req *models.CreateMachineAssignmentRequest) (*models.MachineAssignment, error) {
	// Get machine info
	var machineName, machineCode string
	err := tx.QueryRow("SELECT machine_name, machine_code FROM machines WHERE id = $1", req.MachineID).Scan(&machineName, &machineCode)
	if err != nil {
		return nil, fmt.Errorf("machine not found: %w", err)
	}

	query := `
		INSERT INTO machine_assignments (schedule_id, machine_id, sequence, target_hours, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 'pending', NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	var assignment models.MachineAssignment
	err = tx.QueryRow(query, scheduleID, req.MachineID, req.Sequence, req.TargetHours).
		Scan(&assignment.ID, &assignment.CreatedAt, &assignment.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create machine assignment: %w", err)
	}

	assignment.ScheduleID = scheduleID
	assignment.MachineID = req.MachineID
	assignment.MachineName = machineName
	assignment.MachineCode = machineCode
	assignment.Sequence = req.Sequence
	assignment.TargetHours = req.TargetHours
	assignment.Status = "pending"

	return &assignment, nil
}

// GetByID retrieves a schedule by ID with its machine assignments
func (r *PPICScheduleRepository) GetByID(id int64) (*models.PPICSchedule, error) {
	query := `
		SELECT id, njo, part_name, priority, priority_alpha, material_status, status, progress, 
		       start_date, finish_date, ppic_notes, created_by, created_at, updated_at
		FROM ppic_schedules
		WHERE id = $1 AND deleted_at IS NULL
	`

	var schedule models.PPICSchedule
	err := r.db.QueryRow(query, id).Scan(
		&schedule.ID, &schedule.NJO, &schedule.PartName, &schedule.Priority, &schedule.PriorityAlpha,
		&schedule.MaterialStatus, &schedule.Status, &schedule.Progress, &schedule.StartDate,
		&schedule.FinishDate, &schedule.PPICNotes, &schedule.CreatedBy, &schedule.CreatedAt, &schedule.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Get machine assignments
	assignments, err := r.getMachineAssignments(schedule.ID)
	if err != nil {
		return nil, err
	}
	schedule.MachineAssignments = assignments

	return &schedule, nil
}

// GetByNJO retrieves a schedule by NJO
func (r *PPICScheduleRepository) GetByNJO(njo string) (*models.PPICSchedule, error) {
	query := `
		SELECT id FROM ppic_schedules WHERE njo = $1 AND deleted_at IS NULL
	`
	var id int64
	err := r.db.QueryRow(query, njo).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return r.GetByID(id)
}

// GetAll retrieves all schedules
func (r *PPICScheduleRepository) GetAll() ([]models.PPICSchedule, error) {
	query := `
		SELECT id, njo, part_name, priority, priority_alpha, material_status, status, progress, 
		       start_date, finish_date, ppic_notes, created_by, created_at, updated_at
		FROM ppic_schedules
		WHERE deleted_at IS NULL
		ORDER BY 
			CASE priority 
				WHEN 'Top Urgent' THEN 1 
				WHEN 'Urgent' THEN 2 
				WHEN 'Medium' THEN 3 
				WHEN 'Low' THEN 4 
			END,
			start_date ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.PPICSchedule
	for rows.Next() {
		var s models.PPICSchedule
		err := rows.Scan(
			&s.ID, &s.NJO, &s.PartName, &s.Priority, &s.PriorityAlpha,
			&s.MaterialStatus, &s.Status, &s.Progress, &s.StartDate,
			&s.FinishDate, &s.PPICNotes, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		assignments, err := r.getMachineAssignments(s.ID)
		if err != nil {
			return nil, err
		}
		s.MachineAssignments = assignments
		schedules = append(schedules, s)
	}

	return schedules, nil
}

// GetWithFilters retrieves schedules with filters
// GetWithFilters retrieves schedules with filters
func (r *PPICScheduleRepository) GetWithFilters(filter models.GanttFilterRequest) ([]models.PPICSchedule, error) {
	query := `
		SELECT ps.id, ps.njo, ps.part_name, ps.priority, ps.priority_alpha, ps.material_status, 
		       ps.status, ps.progress, ps.start_date, ps.finish_date, ps.ppic_notes, ps.created_by, 
		       ps.created_at, ps.updated_at
		FROM ppic_schedules ps
		WHERE ps.deleted_at IS NULL
	`

	var args []interface{}
	argNum := 1

	if filter.StartDate != "" {
		query += fmt.Sprintf(" AND ps.start_date >= $%d", argNum)
		args = append(args, filter.StartDate)
		argNum++
	}
	if filter.EndDate != "" {
		query += fmt.Sprintf(" AND ps.finish_date <= $%d", argNum)
		args = append(args, filter.EndDate)
		argNum++
	}
	if filter.Priority != "" {
		query += fmt.Sprintf(" AND ps.priority = $%d", argNum)
		args = append(args, filter.Priority)
		argNum++
	}
	if filter.Status != "" {
		query += fmt.Sprintf(" AND ps.status = $%d", argNum)
		args = append(args, filter.Status)
		argNum++
	}
	if filter.MachineID > 0 {
		query += fmt.Sprintf(" AND ps.id IN (SELECT schedule_id FROM machine_assignments WHERE machine_id = $%d)", argNum)
		args = append(args, filter.MachineID)
		argNum++
	}

	query += ` ORDER BY 
		CASE ps.priority 
			WHEN 'Top Urgent' THEN 1 
			WHEN 'Urgent' THEN 2 
			WHEN 'Medium' THEN 3 
			WHEN 'Low' THEN 4 
		END,
		ps.start_date ASC`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.PPICSchedule
	for rows.Next() {
		var s models.PPICSchedule
		err := rows.Scan(
			&s.ID, &s.NJO, &s.PartName, &s.Priority, &s.PriorityAlpha,
			&s.MaterialStatus, &s.Status, &s.Progress, &s.StartDate,
			&s.FinishDate, &s.PPICNotes, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		assignments, err := r.getMachineAssignments(s.ID)
		if err != nil {
			return nil, err
		}
		s.MachineAssignments = assignments
		schedules = append(schedules, s)
	}

	return schedules, nil
}

// Update updates a schedule
func (r *PPICScheduleRepository) Update(id int64, req *models.UpdatePPICScheduleRequest, startDate, finishDate *time.Time) (*models.PPICSchedule, error) {
	query := "UPDATE ppic_schedules SET updated_at = NOW()"
	var args []interface{}
	argNum := 1

	if req.PartName != "" {
		query += fmt.Sprintf(", part_name = $%d", argNum)
		args = append(args, req.PartName)
		argNum++
	}
	if req.Priority != "" {
		query += fmt.Sprintf(", priority = $%d", argNum)
		args = append(args, req.Priority)
		argNum++
	}
	if req.PriorityAlpha != "" {
		query += fmt.Sprintf(", priority_alpha = $%d", argNum)
		args = append(args, req.PriorityAlpha)
		argNum++
	}
	if req.MaterialStatus != "" {
		query += fmt.Sprintf(", material_status = $%d", argNum)
		args = append(args, req.MaterialStatus)
		argNum++
	}
	if req.Status != "" {
		query += fmt.Sprintf(", status = $%d", argNum)
		args = append(args, req.Status)
		argNum++
	}
	if req.Progress != nil {
		query += fmt.Sprintf(", progress = $%d", argNum)
		args = append(args, *req.Progress)
		argNum++
	}
	if startDate != nil {
		query += fmt.Sprintf(", start_date = $%d", argNum)
		args = append(args, *startDate)
		argNum++
	}
	if finishDate != nil {
		query += fmt.Sprintf(", finish_date = $%d", argNum)
		args = append(args, *finishDate)
		argNum++
	}
	if req.PPICNotes != "" {
		query += fmt.Sprintf(", ppic_notes = $%d", argNum)
		args = append(args, req.PPICNotes)
		argNum++
	}

	query += fmt.Sprintf(" WHERE id = $%d AND deleted_at IS NULL", argNum)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	// Update machine assignments if provided
	for _, ma := range req.MachineAssignments {
		if ma.ID > 0 {
			err = r.updateMachineAssignment(&ma)
			if err != nil {
				return nil, err
			}
		}
	}

	return r.GetByID(id)
}

func (r *PPICScheduleRepository) updateMachineAssignment(req *models.UpdateMachineAssignmentRequest) error {
	query := "UPDATE machine_assignments SET updated_at = NOW()"
	var args []interface{}
	argNum := 1

	if req.Sequence > 0 {
		query += fmt.Sprintf(", sequence = $%d", argNum)
		args = append(args, req.Sequence)
		argNum++
	}
	if req.TargetHours > 0 {
		query += fmt.Sprintf(", target_hours = $%d", argNum)
		args = append(args, req.TargetHours)
		argNum++
	}
	if req.Status != "" {
		query += fmt.Sprintf(", status = $%d", argNum)
		args = append(args, req.Status)
		argNum++
	}
	if req.ActualStart != nil {
		query += fmt.Sprintf(", actual_start = $%d", argNum)
		args = append(args, req.ActualStart)
		argNum++
	}
	if req.ActualEnd != nil {
		query += fmt.Sprintf(", actual_end = $%d", argNum)
		args = append(args, req.ActualEnd)
		argNum++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argNum)
	args = append(args, req.ID)

	_, err := r.db.Exec(query, args...)
	return err
}

// Delete soft deletes a schedule
func (r *PPICScheduleRepository) Delete(id int64) error {
	_, err := r.db.Exec("UPDATE ppic_schedules SET deleted_at = NOW() WHERE id = $1", id)
	return err
}

// GetAllMachines returns all machines
func (r *PPICScheduleRepository) GetAllMachines() ([]models.Machine, error) {
	query := `SELECT id, machine_code, machine_name, machine_type, location, status, created_at, updated_at 
	          FROM machines WHERE deleted_at IS NULL ORDER BY machine_name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var machines []models.Machine
	for rows.Next() {
		var m models.Machine
		err := rows.Scan(&m.ID, &m.MachineCode, &m.MachineName, &m.MachineType, &m.Location, &m.Status, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			return nil, err
		}
		machines = append(machines, m)
	}
	return machines, nil
}

// GetSummary returns summary statistics
func (r *PPICScheduleRepository) GetSummary() (*models.GanttSummary, error) {
	summary := &models.GanttSummary{}

	// Total and status counts
	err := r.db.QueryRow(`
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'completed') as completed,
			COUNT(*) FILTER (WHERE status = 'in_progress') as in_progress,
			COUNT(*) FILTER (WHERE status = 'pending') as pending
		FROM ppic_schedules WHERE deleted_at IS NULL
	`).Scan(&summary.TotalTasks, &summary.CompletedTasks, &summary.InProgressTasks, &summary.PendingTasks)
	if err != nil {
		return nil, err
	}

	// Priority counts
	err = r.db.QueryRow(`
		SELECT 
			COUNT(*) FILTER (WHERE priority = 'Top Urgent') as top_urgent,
			COUNT(*) FILTER (WHERE priority = 'Urgent') as urgent,
			COUNT(*) FILTER (WHERE priority = 'Medium') as medium,
			COUNT(*) FILTER (WHERE priority = 'Low') as low
		FROM ppic_schedules WHERE deleted_at IS NULL
	`).Scan(&summary.TopUrgentCount, &summary.UrgentCount, &summary.MediumCount, &summary.LowCount)
	if err != nil {
		return nil, err
	}

	// Material status counts
	err = r.db.QueryRow(`
		SELECT 
			COUNT(*) FILTER (WHERE material_status = 'Ready') as ready,
			COUNT(*) FILTER (WHERE material_status != 'Ready') as not_ready
		FROM ppic_schedules WHERE deleted_at IS NULL
	`).Scan(&summary.MaterialReady, &summary.MaterialNotReady)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

// GetSchedulesByMachine gets schedules for a specific machine
func (r *PPICScheduleRepository) GetSchedulesByMachine(machineID int64) ([]models.PPICSchedule, error) {
	query := `
		SELECT DISTINCT ps.id, ps.njo, ps.part_name, ps.priority, ps.priority_alpha, ps.material_status, 
		       ps.status, ps.progress, ps.start_date, ps.finish_date, ps.ppic_notes, ps.created_by, 
		       ps.created_at, ps.updated_at
		FROM ppic_schedules ps
		JOIN machine_assignments ma ON ps.id = ma.schedule_id
		WHERE ma.machine_id = $1 AND ps.deleted_at IS NULL
		ORDER BY ps.start_date
	`

	rows, err := r.db.Query(query, machineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.PPICSchedule
	for rows.Next() {
		var s models.PPICSchedule
		err := rows.Scan(
			&s.ID, &s.NJO, &s.PartName, &s.Priority, &s.PriorityAlpha,
			&s.MaterialStatus, &s.Status, &s.Progress, &s.StartDate,
			&s.FinishDate, &s.PPICNotes, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		assignments, err := r.getMachineAssignments(s.ID)
		if err != nil {
			return nil, err
		}
		s.MachineAssignments = assignments
		schedules = append(schedules, s)
	}

	return schedules, nil
}

// AddMachineAssignment adds a machine assignment to a schedule
func (r *PPICScheduleRepository) AddMachineAssignment(req *models.CreateMachineAssignmentRequest, scheduleID int64) (*models.MachineAssignment, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	assignment, err := r.createMachineAssignment(tx, scheduleID, req)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return assignment, nil
}

// DeleteMachineAssignment removes a machine assignment
func (r *PPICScheduleRepository) DeleteMachineAssignment(assignmentID int64) error {
	_, err := r.db.Exec("DELETE FROM machine_assignments WHERE id = $1", assignmentID)
	return err
}

// Helper function to get machine assignments for a schedule
func (r *PPICScheduleRepository) getMachineAssignments(scheduleID int64) ([]models.MachineAssignment, error) {
	query := `
		SELECT ma.id, ma.schedule_id, ma.machine_id, m.machine_name, m.machine_code,
		       ma.sequence, ma.target_hours, ma.scheduled_start, ma.scheduled_end,
		       ma.actual_start, ma.actual_end, ma.status, ma.created_at, ma.updated_at
		FROM machine_assignments ma
		JOIN machines m ON ma.machine_id = m.id
		WHERE ma.schedule_id = $1
		ORDER BY ma.sequence
	`

	rows, err := r.db.Query(query, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []models.MachineAssignment
	for rows.Next() {
		var a models.MachineAssignment
		err := rows.Scan(
			&a.ID, &a.ScheduleID, &a.MachineID, &a.MachineName, &a.MachineCode,
			&a.Sequence, &a.TargetHours, &a.ScheduledStart, &a.ScheduledEnd,
			&a.ActualStart, &a.ActualEnd, &a.Status, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}

	return assignments, nil
}
