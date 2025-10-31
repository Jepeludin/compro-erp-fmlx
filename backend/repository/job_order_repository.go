package repository

import (
	"database/sql"
	"ganttpro-backend/models"
	"time"
)

type JobOrderRepository struct {
	db *sql.DB
}

func NewJobOrderRepository(db *sql.DB) *JobOrderRepository {
	return &JobOrderRepository{db: db}
}

// GetAll retrieves all job orders with machine and operator info
func (r *JobOrderRepository) GetAll() ([]models.JobOrder, error) {
	query := `
		SELECT 
			jo.id, jo.machine_id, m.machine_name, jo.njo, jo.project, jo.item, 
			jo.note, jo.deadline, jo.operator_id, u.username, jo.status, 
			jo.created_at, jo.completed_at, jo.updated_at
		FROM job_orders jo
		LEFT JOIN machines m ON m.id = jo.machine_id
		LEFT JOIN users u ON u.id = jo.operator_id
		WHERE jo.deleted_at IS NULL
		ORDER BY jo.created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.JobOrder
	for rows.Next() {
		var j models.JobOrder
		err := rows.Scan(
			&j.ID,
			&j.MachineID,
			&j.MachineName,
			&j.NJO,
			&j.Project,
			&j.Item,
			&j.Note,
			&j.Deadline,
			&j.OperatorID,
			&j.OperatorName,
			&j.Status,
			&j.CreatedAt,
			&j.CompletedAt,
			&j.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}

	return jobs, nil
}

// GetByID retrieves a job order by ID with stages
func (r *JobOrderRepository) GetByID(id int64) (*models.JobOrder, error) {
	query := `
		SELECT 
			jo.id, jo.machine_id, m.machine_name, jo.njo, jo.project, jo.item, 
			jo.note, jo.deadline, jo.operator_id, u.username, jo.status, 
			jo.created_at, jo.completed_at, jo.updated_at
		FROM job_orders jo
		LEFT JOIN machines m ON m.id = jo.machine_id
		LEFT JOIN users u ON u.id = jo.operator_id
		WHERE jo.id = $1 AND jo.deleted_at IS NULL
	`

	var j models.JobOrder
	err := r.db.QueryRow(query, id).Scan(
		&j.ID,
		&j.MachineID,
		&j.MachineName,
		&j.NJO,
		&j.Project,
		&j.Item,
		&j.Note,
		&j.Deadline,
		&j.OperatorID,
		&j.OperatorName,
		&j.Status,
		&j.CreatedAt,
		&j.CompletedAt,
		&j.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Get process stages
	stages, err := r.GetProcessStages(id)
	if err != nil {
		return nil, err
	}
	j.ProcessStages = stages

	return &j, nil
}

// GetByMachineID retrieves all job orders for a specific machine
func (r *JobOrderRepository) GetByMachineID(machineID int64) ([]models.JobOrder, error) {
	query := `
		SELECT 
			jo.id, jo.machine_id, m.machine_name, jo.njo, jo.project, jo.item, 
			jo.note, jo.deadline, jo.operator_id, u.username, jo.status, 
			jo.created_at, jo.completed_at, jo.updated_at
		FROM job_orders jo
		LEFT JOIN machines m ON m.id = jo.machine_id
		LEFT JOIN users u ON u.id = jo.operator_id
		WHERE jo.machine_id = $1 AND jo.deleted_at IS NULL
		ORDER BY jo.created_at DESC
	`

	rows, err := r.db.Query(query, machineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []models.JobOrder
	for rows.Next() {
		var j models.JobOrder
		err := rows.Scan(
			&j.ID,
			&j.MachineID,
			&j.MachineName,
			&j.NJO,
			&j.Project,
			&j.Item,
			&j.Note,
			&j.Deadline,
			&j.OperatorID,
			&j.OperatorName,
			&j.Status,
			&j.CreatedAt,
			&j.CompletedAt,
			&j.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}

	return jobs, nil
}

// Create creates a new job order with default process stages
func (r *JobOrderRepository) Create(req *models.CreateJobOrderRequest) (*models.JobOrder, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Insert job order
	query := `
		INSERT INTO job_orders (machine_id, njo, project, item, note, deadline, operator_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'pending')
		RETURNING id, machine_id, njo, project, item, note, deadline, operator_id, status, created_at, updated_at
	`

	var j models.JobOrder
	err = tx.QueryRow(
		query,
		req.MachineID,
		req.NJO,
		req.Project,
		req.Item,
		req.Note,
		req.Deadline,
		req.OperatorID,
	).Scan(
		&j.ID,
		&j.MachineID,
		&j.NJO,
		&j.Project,
		&j.Item,
		&j.Note,
		&j.Deadline,
		&j.OperatorID,
		&j.Status,
		&j.CreatedAt,
		&j.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Create default process stages (setting, proses, cmm, kalibrasi)
	stageQuery := `
		INSERT INTO process_stages (job_order_id, stage_name)
		VALUES ($1, $2)
	`

	stages := []string{"setting", "proses", "cmm", "kalibrasi"}
	for _, stageName := range stages {
		_, err = tx.Exec(stageQuery, j.ID, stageName)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &j, nil
}

// Update updates a job order
func (r *JobOrderRepository) Update(id int64, req *models.UpdateJobOrderRequest) (*models.JobOrder, error) {
	query := `
		UPDATE job_orders
		SET project = $1, item = $2, note = $3, deadline = $4, operator_id = $5, status = $6, updated_at = $7
		WHERE id = $8 AND deleted_at IS NULL
		RETURNING id, machine_id, njo, project, item, note, deadline, operator_id, status, created_at, updated_at
	`

	var j models.JobOrder
	err := r.db.QueryRow(
		query,
		req.Project,
		req.Item,
		req.Note,
		req.Deadline,
		req.OperatorID,
		req.Status,
		time.Now(),
		id,
	).Scan(
		&j.ID,
		&j.MachineID,
		&j.NJO,
		&j.Project,
		&j.Item,
		&j.Note,
		&j.Deadline,
		&j.OperatorID,
		&j.Status,
		&j.CreatedAt,
		&j.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &j, nil
}

// Delete soft deletes a job order
func (r *JobOrderRepository) Delete(id int64) error {
	query := `
		UPDATE job_orders
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

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

// GetProcessStages retrieves all process stages for a job order
func (r *JobOrderRepository) GetProcessStages(jobOrderID int64) ([]models.ProcessStage, error) {
	query := `
		SELECT 
			ps.id, ps.job_order_id, ps.stage_name, ps.start_time, ps.finish_time, 
			ps.duration_minutes, ps.operator_id, u.username, ps.notes, ps.created_at, ps.updated_at
		FROM process_stages ps
		LEFT JOIN users u ON u.id = ps.operator_id
		WHERE ps.job_order_id = $1
		ORDER BY 
			CASE ps.stage_name
				WHEN 'setting' THEN 1
				WHEN 'proses' THEN 2
				WHEN 'cmm' THEN 3
				WHEN 'kalibrasi' THEN 4
				ELSE 5
			END
	`

	rows, err := r.db.Query(query, jobOrderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stages []models.ProcessStage
	for rows.Next() {
		var s models.ProcessStage
		err := rows.Scan(
			&s.ID,
			&s.JobOrderID,
			&s.StageName,
			&s.StartTime,
			&s.FinishTime,
			&s.DurationMinutes,
			&s.OperatorID,
			&s.OperatorName,
			&s.Notes,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		stages = append(stages, s)
	}

	return stages, nil
}

// UpdateProcessStage updates a process stage
func (r *JobOrderRepository) UpdateProcessStage(stageID int64, req *models.UpdateProcessStageRequest) (*models.ProcessStage, error) {
	query := `
		UPDATE process_stages
		SET start_time = $1, finish_time = $2, operator_id = $3, notes = $4, updated_at = $5
		WHERE id = $6
		RETURNING id, job_order_id, stage_name, start_time, finish_time, duration_minutes, operator_id, notes, created_at, updated_at
	`

	var s models.ProcessStage
	err := r.db.QueryRow(
		query,
		req.StartTime,
		req.FinishTime,
		req.OperatorID,
		req.Notes,
		time.Now(),
		stageID,
	).Scan(
		&s.ID,
		&s.JobOrderID,
		&s.StageName,
		&s.StartTime,
		&s.FinishTime,
		&s.DurationMinutes,
		&s.OperatorID,
		&s.Notes,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &s, nil
}
