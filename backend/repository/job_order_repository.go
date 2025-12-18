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

// nullableJobOrder holds nullable fields for safe scanning
type nullableJobOrder struct {
	ID           int64
	MachineID    int64
	MachineName  sql.NullString
	NJO          string
	Project      string
	Item         string
	Note         sql.NullString
	Deadline     sql.NullString
	OperatorID   sql.NullInt64
	OperatorName sql.NullString
	Status       string
	CreatedAt    time.Time
	CompletedAt  sql.NullTime
	UpdatedAt    time.Time
}

// toJobOrder converts nullable fields to the model
func (n *nullableJobOrder) toJobOrder() models.JobOrder {
	j := models.JobOrder{
		ID:        n.ID,
		MachineID: n.MachineID,
		NJO:       n.NJO,
		Project:   n.Project,
		Item:      n.Item,
		Status:    n.Status,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}

	if n.MachineName.Valid {
		j.MachineName = n.MachineName.String
	}
	if n.Note.Valid {
		j.Note = n.Note.String
	}
	if n.Deadline.Valid {
		j.Deadline = n.Deadline.String
	}
	if n.OperatorID.Valid {
		id := n.OperatorID.Int64
		j.OperatorID = &id
	}
	if n.OperatorName.Valid {
		j.OperatorName = n.OperatorName.String
	}
	if n.CompletedAt.Valid {
		j.CompletedAt = &n.CompletedAt.Time
	}

	return j
}

// scanJobOrder scans a row into a nullableJobOrder and converts it
func scanJobOrder(scanner interface{ Scan(...any) error }) (models.JobOrder, error) {
	var n nullableJobOrder
	err := scanner.Scan(
		&n.ID,
		&n.MachineID,
		&n.MachineName,
		&n.NJO,
		&n.Project,
		&n.Item,
		&n.Note,
		&n.Deadline,
		&n.OperatorID,
		&n.OperatorName,
		&n.Status,
		&n.CreatedAt,
		&n.CompletedAt,
		&n.UpdatedAt,
	)
	if err != nil {
		return models.JobOrder{}, err
	}
	return n.toJobOrder(), nil
}

// GetAll retrieves all job orders with machine and operator info
func (r *JobOrderRepository) GetAll() ([]models.JobOrder, error) {
	query := `
		SELECT 
			jo.id, jo.machine_id, COALESCE(m.machine_name, ''), jo.njo, jo.project, jo.item, 
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
		j, err := scanJobOrder(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

// GetByID retrieves a job order by ID with stages
func (r *JobOrderRepository) GetByID(id int64) (*models.JobOrder, error) {
	query := `
		SELECT 
			jo.id, jo.machine_id, COALESCE(m.machine_name, ''), jo.njo, jo.project, jo.item, 
			jo.note, jo.deadline, jo.operator_id, u.username, jo.status, 
			jo.created_at, jo.completed_at, jo.updated_at
		FROM job_orders jo
		LEFT JOIN machines m ON m.id = jo.machine_id
		LEFT JOIN users u ON u.id = jo.operator_id
		WHERE jo.id = $1 AND jo.deleted_at IS NULL
	`

	j, err := scanJobOrder(r.db.QueryRow(query, id))
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
			jo.id, jo.machine_id, COALESCE(m.machine_name, ''), jo.njo, jo.project, jo.item, 
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
		j, err := scanJobOrder(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}

	if err = rows.Err(); err != nil {
		return nil, err
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

	// Use nullable types for scanning RETURNING clause
	var (
		id         int64
		machineID  int64
		njo        string
		project    string
		item       string
		note       sql.NullString
		deadline   sql.NullString
		operatorID sql.NullInt64
		status     string
		createdAt  time.Time
		updatedAt  time.Time
	)

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
		&id,
		&machineID,
		&njo,
		&project,
		&item,
		&note,
		&deadline,
		&operatorID,
		&status,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Build JobOrder from scanned values
	j := models.JobOrder{
		ID:        id,
		MachineID: machineID,
		NJO:       njo,
		Project:   project,
		Item:      item,
		Status:    status,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	if note.Valid {
		j.Note = note.String
	}
	if deadline.Valid {
		j.Deadline = deadline.String
	}
	if operatorID.Valid {
		opID := operatorID.Int64
		j.OperatorID = &opID
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

	// Use nullable types for scanning RETURNING clause
	var (
		jobID      int64
		machineID  int64
		njo        string
		project    string
		item       string
		note       sql.NullString
		deadline   sql.NullString
		operatorID sql.NullInt64
		status     string
		createdAt  time.Time
		updatedAt  time.Time
	)

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
		&jobID,
		&machineID,
		&njo,
		&project,
		&item,
		&note,
		&deadline,
		&operatorID,
		&status,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Build JobOrder from scanned values
	j := models.JobOrder{
		ID:        jobID,
		MachineID: machineID,
		NJO:       njo,
		Project:   project,
		Item:      item,
		Status:    status,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	if note.Valid {
		j.Note = note.String
	}
	if deadline.Valid {
		j.Deadline = deadline.String
	}
	if operatorID.Valid {
		opID := operatorID.Int64
		j.OperatorID = &opID
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

// nullableProcessStage holds nullable fields for safe scanning
type nullableProcessStage struct {
	ID              int64
	JobOrderID      int64
	StageName       string
	StartTime       sql.NullTime
	FinishTime      sql.NullTime
	DurationMinutes sql.NullFloat64
	OperatorID      sql.NullInt64
	OperatorName    sql.NullString
	Notes           sql.NullString
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// toProcessStage converts nullable fields to the model
func (n *nullableProcessStage) toProcessStage() models.ProcessStage {
	s := models.ProcessStage{
		ID:         n.ID,
		JobOrderID: n.JobOrderID,
		StageName:  n.StageName,
		CreatedAt:  n.CreatedAt,
		UpdatedAt:  n.UpdatedAt,
	}

	if n.StartTime.Valid {
		s.StartTime = &n.StartTime.Time
	}
	if n.FinishTime.Valid {
		s.FinishTime = &n.FinishTime.Time
	}
	if n.DurationMinutes.Valid {
		s.DurationMinutes = &n.DurationMinutes.Float64
	}
	if n.OperatorID.Valid {
		id := n.OperatorID.Int64
		s.OperatorID = &id
	}
	if n.OperatorName.Valid {
		s.OperatorName = n.OperatorName.String
	}
	if n.Notes.Valid {
		s.Notes = n.Notes.String
	}

	return s
}

// scanProcessStage scans a row into a nullableProcessStage and converts it
func scanProcessStage(scanner interface{ Scan(...any) error }) (models.ProcessStage, error) {
	var n nullableProcessStage
	err := scanner.Scan(
		&n.ID,
		&n.JobOrderID,
		&n.StageName,
		&n.StartTime,
		&n.FinishTime,
		&n.DurationMinutes,
		&n.OperatorID,
		&n.OperatorName,
		&n.Notes,
		&n.CreatedAt,
		&n.UpdatedAt,
	)
	if err != nil {
		return models.ProcessStage{}, err
	}
	return n.toProcessStage(), nil
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
		s, err := scanProcessStage(rows)
		if err != nil {
			return nil, err
		}
		stages = append(stages, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
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

	// Use nullable scanner for RETURNING clause
	var n nullableProcessStage
	err := r.db.QueryRow(
		query,
		req.StartTime,
		req.FinishTime,
		req.OperatorID,
		req.Notes,
		time.Now(),
		stageID,
	).Scan(
		&n.ID,
		&n.JobOrderID,
		&n.StageName,
		&n.StartTime,
		&n.FinishTime,
		&n.DurationMinutes,
		&n.OperatorID,
		&n.Notes,
		&n.CreatedAt,
		&n.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	s := n.toProcessStage()
	return &s, nil
}
