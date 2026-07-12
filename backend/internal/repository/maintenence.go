package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"asset-flow/internal/model"
)

var ErrMaintenanceNotFound = errors.New("maintenance request not found")

type MaintenanceRepository struct {
	db *pgxpool.Pool
}

func NewMaintenanceRepository(db *pgxpool.Pool) *MaintenanceRepository {
	return &MaintenanceRepository{db: db}
}

func (r *MaintenanceRepository) Create(ctx context.Context, m *model.Maintenance) error {
	query := `
		INSERT INTO maintenance (asset_id, raised_by, priority, description, images, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		m.AssetID, m.RaisedBy, m.Priority, m.Description, m.Images, model.MaintenancePending,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (r *MaintenanceRepository) GetByID(ctx context.Context, id int64) (*model.Maintenance, error) {
	query := `
		SELECT id, asset_id, raised_by, assigned_technician_id, priority, description,
		       images, status, created_at, updated_at, deleted_at
		FROM maintenance
		WHERE id = $1 AND deleted_at IS NULL`

	return r.scanOne(r.db.QueryRow(ctx, query, id))
}

type MaintenanceFilter struct {
	Status  string
	AssetID *int64
	Page    int
	Limit   int
}

func (r *MaintenanceRepository) List(ctx context.Context, f MaintenanceFilter) ([]model.Maintenance, int, error) {
	var conditions []string
	var args []any
	conditions = append(conditions, "deleted_at IS NULL")

	addCond := func(cond string, val any) {
		args = append(args, val)
		conditions = append(conditions, fmt.Sprintf(cond, len(args)))
	}
	if f.Status != "" {
		addCond("status = $%d", f.Status)
	}
	if f.AssetID != nil {
		addCond("asset_id = $%d", *f.AssetID)
	}
	where := strings.Join(conditions, " AND ")

	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 20
	}
	offset := (f.Page - 1) * f.Limit

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM maintenance WHERE "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, f.Limit, offset)
	dataQuery := fmt.Sprintf(`
		SELECT id, asset_id, raised_by, assigned_technician_id, priority, description,
		       images, status, created_at, updated_at, deleted_at
		FROM maintenance
		WHERE %s
		ORDER BY
		  CASE priority WHEN 'CRITICAL' THEN 0 WHEN 'HIGH' THEN 1 WHEN 'MEDIUM' THEN 2 ELSE 3 END,
		  created_at DESC
		LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args))

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var records []model.Maintenance
	for rows.Next() {
		var m model.Maintenance
		if err := rows.Scan(
			&m.ID, &m.AssetID, &m.RaisedBy, &m.AssignedTechnicianID, &m.Priority,
			&m.Description, &m.Images, &m.Status, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
		); err != nil {
			return nil, 0, err
		}
		records = append(records, m)
	}
	return records, total, rows.Err()
}

func (r *MaintenanceRepository) CountDueToday(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM maintenance
		WHERE status NOT IN ($1, $2) AND deleted_at IS NULL
		  AND created_at::date = CURRENT_DATE`,
		model.MaintenanceResolved, model.MaintenanceRejected,
	).Scan(&count)
	return count, err
}


func (r *MaintenanceRepository) UpdateStatus(ctx context.Context, tx pgx.Tx, id int64, status string, technicianID *int64) (*model.Maintenance, error) {
	query := `
		UPDATE maintenance
		SET status = $2, assigned_technician_id = COALESCE($3, assigned_technician_id), updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, asset_id, raised_by, assigned_technician_id, priority, description,
		          images, status, created_at, updated_at, deleted_at`

	return r.scanOne(tx.QueryRow(ctx, query, id, status, technicianID))
}


func (r *MaintenanceRepository) CreateWorkflowStep(ctx context.Context, tx pgx.Tx, w *model.MaintenanceWorkflow) error {
	query := `
		INSERT INTO maintenance_workflows (maintenance_id, status, description, updated_by, workflow_date)
		VALUES ($1, $2, $3, $4, now())
		RETURNING id, workflow_date, created_at, updated_at`

	return tx.QueryRow(ctx, query,
		w.MaintenanceID, w.Status, w.Description, w.UpdatedBy,
	).Scan(&w.ID, &w.WorkflowDate, &w.CreatedAt, &w.UpdatedAt)
}

func (r *MaintenanceRepository) ListWorkflowHistory(ctx context.Context, maintenanceID int64) ([]model.MaintenanceWorkflow, error) {
	query := `
		SELECT id, maintenance_id, status, description, updated_by, workflow_date, created_at, updated_at, deleted_at
		FROM maintenance_workflows
		WHERE maintenance_id = $1 AND deleted_at IS NULL
		ORDER BY workflow_date ASC`

	rows, err := r.db.Query(ctx, query, maintenanceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var steps []model.MaintenanceWorkflow
	for rows.Next() {
		var w model.MaintenanceWorkflow
		if err := rows.Scan(
			&w.ID, &w.MaintenanceID, &w.Status, &w.Description, &w.UpdatedBy,
			&w.WorkflowDate, &w.CreatedAt, &w.UpdatedAt, &w.DeletedAt,
		); err != nil {
			return nil, err
		}
		steps = append(steps, w)
	}
	return steps, rows.Err()
}

func (r *MaintenanceRepository) ListByAsset(ctx context.Context, assetID int64) ([]model.Maintenance, error) {
	query := `
		SELECT id, asset_id, raised_by, assigned_technician_id, priority, description,
		       images, status, created_at, updated_at, deleted_at
		FROM maintenance
		WHERE asset_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.Maintenance
	for rows.Next() {
		var m model.Maintenance
		if err := rows.Scan(
			&m.ID, &m.AssetID, &m.RaisedBy, &m.AssignedTechnicianID, &m.Priority,
			&m.Description, &m.Images, &m.Status, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
		); err != nil {
			return nil, err
		}
		records = append(records, m)
	}
	return records, rows.Err()
}

func (r *MaintenanceRepository) scanOne(row pgx.Row) (*model.Maintenance, error) {
	var m model.Maintenance
	err := row.Scan(
		&m.ID, &m.AssetID, &m.RaisedBy, &m.AssignedTechnicianID, &m.Priority,
		&m.Description, &m.Images, &m.Status, &m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMaintenanceNotFound
		}
		return nil, err
	}
	return &m, nil
}