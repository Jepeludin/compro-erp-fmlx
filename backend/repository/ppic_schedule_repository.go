package repository

import (
	"database/sql"
	"fmt"
	"ganttpro-backend/models"
	"strings"
	"time"
)

type PPICScheduleRepository struct {
	db *sql.DB
}

func NewPPICScheduleRepository(db *sql.DB) *PPICScheduleRepository {
	return &PPICScheduleRepository{db: db}
}

// GetAll retrieves all PPIC schedules with machine assignments
func (r *PPICScheduleRepository) GetAll() ([]models.PPICSchedule, error) {
	query := `
		SELECT 
			ps.id, ps.njo, ps.part_name, ps.start_date, ps.finish_date,
			ps.priority, ps.priority_alpha, ps.material_status, ps.ppic_notes,
			ps.status, ps.progress, ps.created_by, u.username,
			ps.created_at, ps.updated_at
		FROM ppic_schedules ps
		LEFT JOIN users u ON u.id = ps.created_by
		WHERE ps.deleted_at IS NULL
		ORDER BY 
			CASE ps.priority 
				WHEN 'Top Urgent' THEN 1 
				WHEN 'Urgent' THEN 2 
				WHEN 'Medium' THEN 3 
				WHEN 'Low' THEN 4 
			END,
			ps.start_date ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.PPICSchedule
	for rows.Next() {
		var s models.PPICSchedule
		var createdByName sql.NullString

		err := rows.Scan(
			&s.ID, &s.NJO, &s.PartName, &s.StartDate, &s.FinishDate,
			&s.Priority, &s.PriorityAlpha, &s.MaterialStatus, &s.PPICNotes,
			&s.Status, &s.Progress, &s.CreatedBy, &createdByName,
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if createdByName.Valid {
			s.CreatedByName = createdByName.String
		}

		// Get machine assignments for this schedule
		assignments, err := r.GetMachineAssignments(s.ID)
		if err != nil {
			return nil, err
		}
		s.MachineAssignments = assignments

		schedules = append(schedules, s)
	}

	return schedules, nil
}

// GetByID retrieves a single PPIC schedule by ID
func (r *PPICScheduleRepository) GetByID(id int64) (*models.PPICSchedule, error) {
	query := `
		SELECT 
			ps.id, ps.njo, ps.part_name, ps.start_date, ps.finish_date,
			ps.priority, ps.priority_alpha, ps.material_status, ps.ppic_notes,
			ps.status, ps.progress, ps.created_by, u.username,
			ps.created_at, ps.updated_at
		FROM ppic_schedules ps
		LEFT JOIN users u ON u.id = ps.created_by
		WHERE ps.id = $1 AND ps.deleted_at IS NULL
	`

	var s models.PPICSchedule
	var createdByName sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&s.ID, &s.NJO, &s.PartName, &s.StartDate, &s.FinishDate,
		&s.Priority, &s.PriorityAlpha, &s.MaterialStatus, &s.PPICNotes,
		&s.Status, &s.Progress, &s.CreatedBy, &createdByName,
		&s.CreatedAt, &s.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if createdByName.Valid {
		s.CreatedByName = createdByName.String
	}

	// Get machine assignments
	assignments, err := r.GetMachineAssignments(s.ID)
	if err != nil {
		return nil, err
	}
	s.MachineAssignments = assignments

	return &s, nil
}

// GetByNJO retrieves a PPIC schedule by NJO
func (r *PPICScheduleRepository) GetByNJO(njo string) (*models.PPICSchedule, error) {
	query := `
		SELECT 
			ps.id, ps.njo, ps.part_name, ps.start_date, ps.finish_date,
			ps.priority, ps.priority_alpha, ps.material_status, ps.ppic_notes,
			ps.status, ps.progress, ps.created_by, u.username,
			ps.created_at, ps.updated_at
		FROM ppic_schedules ps
		LEFT JOIN users u ON u.id = ps.created_by
		WHERE ps.njo = $1 AND ps.deleted_at IS NULL
	`

	var s models.PPICSchedule
	var createdByName sql.NullString

	err := r.db.QueryRow(query, njo).Scan(
		&s.ID, &s.NJO, &s.PartName, &s.StartDate, &s.FinishDate,
		&s.Priority, &s.PriorityAlpha, &s.MaterialStatus, &s.PPICNotes,
		&s.Status, &s.Progress, &s.CreatedBy, &createdByName,
		&s.CreatedAt, &s.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if createdByName.Valid {
		s.CreatedByName = createdByName.String
	}

	return &s, nil
}

// GetWithFilters retrieves PPIC schedules with filters
func (r *PPICScheduleRepository) GetWithFilters(filter models.GanttFilterRequest) ([]models.PPICSchedule, error) {
	query := `
		SELECT DISTINCT
			ps.id, ps.njo, ps.part_name, ps.start_date, ps.finish_date,
			ps.priority, ps.priority_alpha, ps.material_status, ps.ppic_notes,
			ps.status, ps.progress, ps.created_by, u.username,
			ps.created_at, ps.updated_at
		FROM ppic_schedules ps
		LEFT JOIN users u ON u.id = ps.created_by
		LEFT JOIN machine_assignments ma ON ma.ppic_schedule_id = ps.id
		WHERE ps.deleted_at IS NULL
	`

	var conditions []string
	var args []interface{}
	argNum := 1

	// Date range filter
	if filter.StartDate != "" {
		conditions = append(conditions, fmt.Sprintf("ps.start_date >= $%d", argNum))
		args = append(args, filter.StartDate)
		argNum++
	}
	if filter.EndDate != "" {
		conditions = append(conditions, fmt.Sprintf("ps.finish_date <= $%d", argNum))
		args = append(args, filter.EndDate)
		argNum++
	}

	// Machine filter
	if filter.MachineID > 0 {
		conditions = append(conditions, fmt.Sprintf("ma.machine_id = $%d", argNum))
		args = append(args, filter.MachineID)
		argNum++
	}

	// Priority filter
	if filter.Priority != "" {
		conditions = append(conditions, fmt.Sprintf("ps.priority = $%d", argNum))
		args = append(args, filter.Priority)
		argNum++
	}

	// Status filter
	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("ps.status = $%d", argNum))
		args = append(args, filter.Status)
		argNum++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += `
		ORDER BY 
			CASE ps.priority 
				WHEN 'Top Urgent' THEN 1 
				WHEN 'Urgent' THEN 2 
				WHEN 'Medium' THEN 3 
				WHEN 'Low' THEN 4 
			END,
			ps.start_date ASC
	`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.PPICSchedule
	for rows.Next() {
		var s models.PPICSchedule
		var createdByName sql.NullString

		err := rows.Scan(
			&s.ID, &s.NJO, &s.PartName, &s.StartDate, &s.FinishDate,
			&s.Priority, &s.PriorityAlpha, &s.MaterialStatus, &s.PPICNotes,
			&s.Status, &s.Progress, &s.CreatedBy, &createdByName,
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if createdByName.Valid {
			s.CreatedByName = createdByName.String
		}

		// Get machine assignments
		assignments, err := r.GetMachineAssignments(s.ID)
		if err != nil {
			return nil, err
		}
		s.MachineAssignments = assignments

		schedules = append(schedules, s)
	}

	return schedules, nil
}

// Create creates a new PPIC schedule with machine assignments
func (r *PPICScheduleRepository) Create(req *models.CreatePPICScheduleRequest, createdBy int64, startDate, finishDate time.Time) (*models.PPICSchedule, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert PPIC schedule
	query := `
		INSERT INTO ppic_schedules (
			njo, part_name, start_date, finish_date, priority, priority_alpha,
			material_status, ppic_notes, status, progress, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'pending', 0, $9)
		RETURNING id, njo, part_name, start_date, finish_date, priority, priority_alpha,
			material_status, ppic_notes, status, progress, created_by, created_at, updated_at
	`

	var s models.PPICSchedule
	err = tx.QueryRow(
		query,
		req.NJO, req.PartName, startDate, finishDate, req.Priority,
		req.PriorityAlpha, req.MaterialStatus, req.PPICNotes, createdBy,
	).Scan(
		&s.ID, &s.NJO, &s.PartName, &s.StartDate, &s.FinishDate,
		&s.Priority, &s.PriorityAlpha, &s.MaterialStatus, &s.PPICNotes,
		&s.Status, &s.Progress, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Insert machine assignments
	assignmentQuery := `
		INSERT INTO machine_assignments (
			ppic_schedule_id, machine_id, target_hours, sequence, status
		) VALUES ($1, $2, $3, $4, 'pending')
		RETURNING id, ppic_schedule_id, machine_id, target_hours, sequence, status, created_at, updated_at
	`

	for _, ma := range req.MachineAssignments {
		var assignment models.MachineAssignment
		err = tx.QueryRow(
			assignmentQuery,
			s.ID, ma.MachineID, ma.TargetHours, ma.Sequence,
		).Scan(
			&assignment.ID, &assignment.PPICScheduleID, &assignment.MachineID,
			&assignment.TargetHours, &assignment.Sequence, &assignment.Status,
			&assignment.CreatedAt, &assignment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		s.MachineAssignments = append(s.MachineAssignments, assignment)
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &s, nil
}

// Update updates a PPIC schedule
func (r *PPICScheduleRepository) Update(id int64, req *models.UpdatePPICScheduleRequest, startDate, finishDate *time.Time) (*models.PPICSchedule, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Build dynamic update query
	updates := []string{"updated_at = $1"}
	args := []interface{}{time.Now()}
	argNum := 2

	if req.PartName != "" {
		updates = append(updates, fmt.Sprintf("part_name = $%d", argNum))
		args = append(args, req.PartName)
		argNum++
	}
	if startDate != nil {
		updates = append(updates, fmt.Sprintf("start_date = $%d", argNum))
		args = append(args, *startDate)
		argNum++
	}
	if finishDate != nil {
		updates = append(updates, fmt.Sprintf("finish_date = $%d", argNum))
		args = append(args, *finishDate)
		argNum++
	}
	if req.Priority != "" {
		updates = append(updates, fmt.Sprintf("priority = $%d", argNum))
		args = append(args, req.Priority)
		argNum++
	}
	if req.PriorityAlpha != "" {
		updates = append(updates, fmt.Sprintf("priority_alpha = $%d", argNum))
		args = append(args, req.PriorityAlpha)
		argNum++
	}
	if req.MaterialStatus != "" {
		updates = append(updates, fmt.Sprintf("material_status = $%d", argNum))
		args = append(args, req.MaterialStatus)
		argNum++
	}
	if req.PPICNotes != "" {
		updates = append(updates, fmt.Sprintf("ppic_notes = $%d", argNum))
		args = append(args, req.PPICNotes)
		argNum++
	}
	if req.Status != "" {
		updates = append(updates, fmt.Sprintf("status = $%d", argNum))
		args = append(args, req.Status)
		argNum++
	}
	if req.Progress != nil {
		updates = append(updates, fmt.Sprintf("progress = $%d", argNum))
		args = append(args, *req.Progress)
		argNum++
	}

	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE ppic_schedules SET %s
		WHERE id = $%d AND deleted_at IS NULL
		RETURNING id, njo, part_name, start_date, finish_date, priority, priority_alpha,
			material_status, ppic_notes, status, progress, created_by, created_at, updated_at
	`, strings.Join(updates, ", "), argNum)

	var s models.PPICSchedule
	err = tx.QueryRow(query, args...).Scan(
		&s.ID, &s.NJO, &s.PartName, &s.StartDate, &s.FinishDate,
		&s.Priority, &s.PriorityAlpha, &s.MaterialStatus, &s.PPICNotes,
		&s.Status, &s.Progress, &s.CreatedBy, &s.CreatedAt, &s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Update machine assignments if provided
	if len(req.MachineAssignments) > 0 {
		for _, ma := range req.MachineAssignments {
			if ma.ID > 0 {
				// Update existing assignment
				_, err = tx.Exec(`
					UPDATE machine_assignments
					SET machine_id = $1, target_hours = $2, sequence = $3,
						scheduled_start = $4, scheduled_end = $5,
						actual_start = $6, actual_end = $7, status = $8, updated_at = $9
					WHERE id = $10
				`, ma.MachineID, ma.TargetHours, ma.Sequence,
					ma.ScheduledStart, ma.ScheduledEnd,
					ma.ActualStart, ma.ActualEnd, ma.Status, time.Now(), ma.ID)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	// Reload with machine assignments
	return r.GetByID(id)
}

// Delete soft deletes a PPIC schedule
func (r *PPICScheduleRepository) Delete(id int64) error {
	query := `UPDATE ppic_schedules SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	result, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetMachineAssignments retrieves all machine assignments for a schedule
func (r *PPICScheduleRepository) GetMachineAssignments(ppicScheduleID int64) ([]models.MachineAssignment, error) {
	query := `
		SELECT 
			ma.id, ma.ppic_schedule_id, ma.machine_id, m.machine_name, m.machine_code,
			ma.target_hours, ma.scheduled_start, ma.scheduled_end,
			ma.actual_start, ma.actual_end, ma.status, ma.sequence,
			ma.created_at, ma.updated_at
		FROM machine_assignments ma
		LEFT JOIN machines m ON m.id = ma.machine_id
		WHERE ma.ppic_schedule_id = $1
		ORDER BY ma.sequence ASC
	`

	rows, err := r.db.Query(query, ppicScheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []models.MachineAssignment
	for rows.Next() {
		var a models.MachineAssignment
		var machineName, machineCode sql.NullString

		err := rows.Scan(
			&a.ID, &a.PPICScheduleID, &a.MachineID, &machineName, &machineCode,
			&a.TargetHours, &a.ScheduledStart, &a.ScheduledEnd,
			&a.ActualStart, &a.ActualEnd, &a.Status, &a.Sequence,
			&a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if machineName.Valid {
			a.MachineName = machineName.String
		}
		if machineCode.Valid {
			a.MachineCode = machineCode.String
		}

		assignments = append(assignments, a)
	}

	return assignments, nil
}

// GetSchedulesByMachine retrieves all schedules assigned to a specific machine
func (r *PPICScheduleRepository) GetSchedulesByMachine(machineID int64) ([]models.PPICSchedule, error) {
	query := `
		SELECT DISTINCT
			ps.id, ps.njo, ps.part_name, ps.start_date, ps.finish_date,
			ps.priority, ps.priority_alpha, ps.material_status, ps.ppic_notes,
			ps.status, ps.progress, ps.created_by, u.username,
			ps.created_at, ps.updated_at
		FROM ppic_schedules ps
		LEFT JOIN users u ON u.id = ps.created_by
		INNER JOIN machine_assignments ma ON ma.ppic_schedule_id = ps.id
		WHERE ma.machine_id = $1 AND ps.deleted_at IS NULL
		ORDER BY ps.start_date ASC
	`

	rows, err := r.db.Query(query, machineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.PPICSchedule
	for rows.Next() {
		var s models.PPICSchedule
		var createdByName sql.NullString

		err := rows.Scan(
			&s.ID, &s.NJO, &s.PartName, &s.StartDate, &s.FinishDate,
			&s.Priority, &s.PriorityAlpha, &s.MaterialStatus, &s.PPICNotes,
			&s.Status, &s.Progress, &s.CreatedBy, &createdByName,
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if createdByName.Valid {
			s.CreatedByName = createdByName.String
		}

		// Get machine assignments
		assignments, err := r.GetMachineAssignments(s.ID)
		if err != nil {
			return nil, err
		}
		s.MachineAssignments = assignments

		schedules = append(schedules, s)
	}

	return schedules, nil
}

// GetSummary returns summary statistics for Gantt chart
func (r *PPICScheduleRepository) GetSummary() (*models.GanttSummary, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending,
			COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
			COUNT(CASE WHEN priority = 'Top Urgent' THEN 1 END) as top_urgent,
			COUNT(CASE WHEN priority = 'Urgent' THEN 1 END) as urgent,
			COUNT(CASE WHEN priority = 'Medium' THEN 1 END) as medium,
			COUNT(CASE WHEN priority = 'Low' THEN 1 END) as low,
			COUNT(CASE WHEN material_status = 'Ready' THEN 1 END) as material_ready
		FROM ppic_schedules
		WHERE deleted_at IS NULL
	`

	var summary models.GanttSummary
	err := r.db.QueryRow(query).Scan(
		&summary.TotalJobs,
		&summary.Pending,
		&summary.InProgress,
		&summary.Completed,
		&summary.TopUrgent,
		&summary.Urgent,
		&summary.Medium,
		&summary.Low,
		&summary.MaterialReady,
	)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}

// GetAllMachines retrieves all active machines
func (r *PPICScheduleRepository) GetAllMachines() ([]models.MachineInfo, error) {
	query := `
		SELECT id, machine_code, machine_name, status
		FROM machines
		WHERE deleted_at IS NULL
		ORDER BY machine_name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var machines []models.MachineInfo
	for rows.Next() {
		var m models.MachineInfo
		err := rows.Scan(&m.ID, &m.MachineCode, &m.MachineName, &m.Status)
		if err != nil {
			return nil, err
		}
		machines = append(machines, m)
	}

	return machines, nil
}

// AddMachineAssignment adds a new machine assignment to an existing schedule
func (r *PPICScheduleRepository) AddMachineAssignment(req *models.CreateMachineAssignmentRequest, ppicScheduleID int64) (*models.MachineAssignment, error) {
	query := `
		INSERT INTO machine_assignments (
			ppic_schedule_id, machine_id, target_hours, sequence, status
		) VALUES ($1, $2, $3, $4, 'pending')
		RETURNING id, ppic_schedule_id, machine_id, target_hours, sequence, status, created_at, updated_at
	`

	var a models.MachineAssignment
	err := r.db.QueryRow(
		query,
		ppicScheduleID, req.MachineID, req.TargetHours, req.Sequence,
	).Scan(
		&a.ID, &a.PPICScheduleID, &a.MachineID,
		&a.TargetHours, &a.Sequence, &a.Status,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

// DeleteMachineAssignment removes a machine assignment
func (r *PPICScheduleRepository) DeleteMachineAssignment(assignmentID int64) error {
	query := `DELETE FROM machine_assignments WHERE id = $1`
	result, err := r.db.Exec(query, assignmentID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}


