package repository

import (
	"database/sql"
	"ganttpro-backend/models"
	"time"
)

type MachineRepository struct {
	db *sql.DB
}

func NewMachineRepository(db *sql.DB) *MachineRepository {
	return &MachineRepository{db: db}
}

// GetAll retrieves all active machines (without pagination, for backward compatibility)
func (r *MachineRepository) GetAll() ([]models.Machine, error) {
	query := `
		SELECT id, machine_code, machine_name, machine_type, location, status, created_at, updated_at
		FROM machines
		WHERE deleted_at IS NULL
		ORDER BY machine_name ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var machines []models.Machine
	for rows.Next() {
		var m models.Machine
		err := rows.Scan(
			&m.ID,
			&m.MachineCode,
			&m.MachineName,
			&m.MachineType,
			&m.Location,
			&m.Status,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		machines = append(machines, m)
	}

	return machines, nil
}

// GetAllPaginated retrieves machines with pagination
func (r *MachineRepository) GetAllPaginated(limit, offset int, sortField, sortOrder string) ([]models.Machine, int64, error) {
	// Whitelist allowed sort fields to prevent SQL injection
	allowedSortFields := map[string]bool{
		"id": true, "machine_code": true, "machine_name": true,
		"machine_type": true, "location": true, "status": true,
		"created_at": true, "updated_at": true,
	}

	if !allowedSortFields[sortField] {
		sortField = "machine_name"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM machines WHERE deleted_at IS NULL`
	err := r.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data
	query := `
		SELECT id, machine_code, machine_name, machine_type, location, status, created_at, updated_at
		FROM machines
		WHERE deleted_at IS NULL
		ORDER BY ` + sortField + ` ` + sortOrder + `
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var machines []models.Machine
	for rows.Next() {
		var m models.Machine
		err := rows.Scan(
			&m.ID,
			&m.MachineCode,
			&m.MachineName,
			&m.MachineType,
			&m.Location,
			&m.Status,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		machines = append(machines, m)
	}

	return machines, total, nil
}

// GetByID retrieves a machine by ID
func (r *MachineRepository) GetByID(id int64) (*models.Machine, error) {
	query := `
		SELECT id, machine_code, machine_name, machine_type, location, status, created_at, updated_at
		FROM machines
		WHERE id = $1 AND deleted_at IS NULL
	`

	var m models.Machine
	err := r.db.QueryRow(query, id).Scan(
		&m.ID,
		&m.MachineCode,
		&m.MachineName,
		&m.MachineType,
		&m.Location,
		&m.Status,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// Create creates a new machine
func (r *MachineRepository) Create(req *models.CreateMachineRequest) (*models.Machine, error) {
	query := `
		INSERT INTO machines (machine_code, machine_name, machine_type, location, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, machine_code, machine_name, machine_type, location, status, created_at, updated_at
	`

	status := req.Status
	if status == "" {
		status = "active"
	}

	var m models.Machine
	err := r.db.QueryRow(
		query,
		req.MachineCode,
		req.MachineName,
		req.MachineType,
		req.Location,
		status,
	).Scan(
		&m.ID,
		&m.MachineCode,
		&m.MachineName,
		&m.MachineType,
		&m.Location,
		&m.Status,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

// Update updates a machine
func (r *MachineRepository) Update(id int64, req *models.UpdateMachineRequest) (*models.Machine, error) {
	query := `
		UPDATE machines
		SET machine_name = $1, machine_type = $2, location = $3, status = $4, updated_at = $5
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING id, machine_code, machine_name, machine_type, location, status, created_at, updated_at
	`

	var m models.Machine
	err := r.db.QueryRow(
		query,
		req.MachineName,
		req.MachineType,
		req.Location,
		req.Status,
		time.Now(),
		id,
	).Scan(
		&m.ID,
		&m.MachineCode,
		&m.MachineName,
		&m.MachineType,
		&m.Location,
		&m.Status,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// Delete soft deletes a machine
func (r *MachineRepository) Delete(id int64) error {
	query := `
		UPDATE machines
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

// GetByCode retrieves a machine by machine code
func (r *MachineRepository) GetByCode(code string) (*models.Machine, error) {
	query := `
		SELECT id, machine_code, machine_name, machine_type, location, status, created_at, updated_at
		FROM machines
		WHERE machine_code = $1 AND deleted_at IS NULL
	`

	var m models.Machine
	err := r.db.QueryRow(query, code).Scan(
		&m.ID,
		&m.MachineCode,
		&m.MachineName,
		&m.MachineType,
		&m.Location,
		&m.Status,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &m, nil
}
