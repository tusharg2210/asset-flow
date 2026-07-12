package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"asset-flow/internal/model"
)

var (
	ErrDepartmentExists   = errors.New("department name already exists")
	ErrDepartmentNotFound = errors.New("department not found")
)

type DepartmentRepository struct {
	db *pgxpool.Pool
}

func NewDepartmentRepository(db *pgxpool.Pool) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) Create(ctx context.Context, d *model.Department) error {
	query := `
		INSERT INTO departments (name, parent_department_id, head_id, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		d.Name, d.ParentDepartmentID, d.HeadID, d.Status,
	).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDepartmentExists
		}
		return err
	}

	return nil
}

func (r *DepartmentRepository) GetByID(ctx context.Context, id int64) (*model.Department, error) {
	query := `
		SELECT id, name, parent_department_id, head_id, status, created_at, updated_at, deleted_at
		FROM departments
		WHERE id = $1 AND deleted_at IS NULL`

	return r.scanOne(r.db.QueryRow(ctx, query, id))
}


func (r *DepartmentRepository) List(ctx context.Context, status string) ([]model.Department, error) {
	query := `
		SELECT id, name, parent_department_id, head_id, status, created_at, updated_at, deleted_at
		FROM departments
		WHERE deleted_at IS NULL`
	args := []any{}

	if status != "" {
		args = append(args, status)
		query += ` AND status = $1`
	}
	query += ` ORDER BY name ASC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []model.Department
	for rows.Next() {
		var d model.Department
		if err := rows.Scan(
			&d.ID, &d.Name, &d.ParentDepartmentID, &d.HeadID, &d.Status,
			&d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
		); err != nil {
			return nil, err
		}
		departments = append(departments, d)
	}
	return departments, rows.Err()
}


func (r *DepartmentRepository) Update(ctx context.Context, id int64, name *string, headID *int64, parentID *int64, status *string) (*model.Department, error) {
	query := `
		UPDATE departments
		SET name = COALESCE($2, name),
		    head_id = COALESCE($3, head_id),
		    parent_department_id = COALESCE($4, parent_department_id),
		    status = COALESCE($5, status),
		    updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, name, parent_department_id, head_id, status, created_at, updated_at, deleted_at`

	return r.scanOne(r.db.QueryRow(ctx, query, id, name, headID, parentID, status))
}

func (r *DepartmentRepository) scanOne(row pgx.Row) (*model.Department, error) {
	var d model.Department
	err := row.Scan(
		&d.ID, &d.Name, &d.ParentDepartmentID, &d.HeadID, &d.Status,
		&d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDepartmentNotFound
		}
		return nil, err
	}
	return &d, nil
}