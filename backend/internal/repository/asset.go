package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"asset-flow/internal/model"
)

var (
	ErrAssetTagExists = errors.New("asset tag already exists")
	ErrAssetNotFound  = errors.New("asset not found")
)

type AssetRepository struct {
	db *pgxpool.Pool
}

func NewAssetRepository(db *pgxpool.Pool) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) Create(ctx context.Context, a *model.Asset) error {
	query := `
		INSERT INTO assets (
			tag, name, category_id, serial_number, qr_code, status, location,
			expected_location, condition, photos_docs, custom_field_values,
			is_sharable, is_bookable, acquisition_date, acquisition_cost
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		a.Tag, a.Name, a.CategoryID, a.SerialNumber, a.QRCode, a.Status, a.Location,
		a.ExpectedLocation, a.Condition, a.PhotosDocs, a.CustomFieldValues,
		a.IsSharable, a.IsBookable, a.AcquisitionDate, a.AcquisitionCost,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrAssetTagExists
		}
		return err
	}

	return nil
}

func (r *AssetRepository) GetByID(ctx context.Context, id int64) (*model.Asset, error) {
	query := `
		SELECT id, tag, name, category_id, serial_number, qr_code, status, location,
		       expected_location, condition, photos_docs, custom_field_values,
		       is_sharable, is_bookable, acquisition_date, acquisition_cost,
		       created_at, updated_at, deleted_at
		FROM assets
		WHERE id = $1 AND deleted_at IS NULL`

	return r.scanOne(r.db.QueryRow(ctx, query, id))
}

func (r *AssetRepository) GetByIDForUpdate(ctx context.Context, tx pgx.Tx, id int64) (*model.Asset, error) {
	query := `
		SELECT id, tag, name, category_id, serial_number, qr_code, status, location,
		       expected_location, condition, photos_docs, custom_field_values,
		       is_sharable, is_bookable, acquisition_date, acquisition_cost,
		       created_at, updated_at, deleted_at
		FROM assets
		WHERE id = $1 AND deleted_at IS NULL
		FOR UPDATE`

	return r.scanOne(tx.QueryRow(ctx, query, id))
}

type AssetFilter struct {
	Tag          string
	SerialNumber string
	QRCode       string
	CategoryID   *int64
	Status       string
	DepartmentID *int64 
	Location     string
	Page         int
	Limit        int
}

func (r *AssetRepository) List(ctx context.Context, f AssetFilter) ([]model.Asset, int, error) {
	var conditions []string
	var args []any
	conditions = append(conditions, "a.deleted_at IS NULL")

	addCond := func(cond string, val any) {
		args = append(args, val)
		conditions = append(conditions, fmt.Sprintf(cond, len(args)))
	}

	if f.Tag != "" {
		addCond("a.tag ILIKE $%d", "%"+f.Tag+"%")
	}
	if f.SerialNumber != "" {
		addCond("a.serial_number = $%d", f.SerialNumber)
	}
	if f.QRCode != "" {
		addCond("a.qr_code = $%d", f.QRCode)
	}
	if f.CategoryID != nil {
		addCond("a.category_id = $%d", *f.CategoryID)
	}
	if f.Status != "" {
		addCond("a.status = $%d", f.Status)
	}
	if f.Location != "" {
		addCond("a.location ILIKE $%d", "%"+f.Location+"%")
	}

	where := strings.Join(conditions, " AND ")

	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 20
	}
	offset := (f.Page - 1) * f.Limit

	countQuery := "SELECT COUNT(*) FROM assets a WHERE " + where
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, f.Limit, offset)
	dataQuery := fmt.Sprintf(`
		SELECT a.id, a.tag, a.name, a.category_id, a.serial_number, a.qr_code, a.status,
		       a.location, a.expected_location, a.condition, a.photos_docs,
		       a.custom_field_values, a.is_sharable, a.is_bookable, a.acquisition_date,
		       a.acquisition_cost, a.created_at, a.updated_at, a.deleted_at
		FROM assets a
		WHERE %s
		ORDER BY a.id DESC
		LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args))

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var assets []model.Asset
	for rows.Next() {
		var a model.Asset
		if err := rows.Scan(
			&a.ID, &a.Tag, &a.Name, &a.CategoryID, &a.SerialNumber, &a.QRCode, &a.Status,
			&a.Location, &a.ExpectedLocation, &a.Condition, &a.PhotosDocs,
			&a.CustomFieldValues, &a.IsSharable, &a.IsBookable, &a.AcquisitionDate,
			&a.AcquisitionCost, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
		); err != nil {
			return nil, 0, err
		}
		assets = append(assets, a)
	}
	return assets, total, rows.Err()
}

func (r *AssetRepository) UpdateStatus(ctx context.Context, tx pgx.Tx, id int64, status string) error {
	query := `UPDATE assets SET status = $2, updated_at = now() WHERE id = $1 AND deleted_at IS NULL`

	var cmd pgconn.CommandTag
	var err error
	if tx != nil {
		cmd, err = tx.Exec(ctx, query, id, status)
	} else {
		cmd, err = r.db.Exec(ctx, query, id, status)
	}
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrAssetNotFound
	}
	return nil
}

func (r *AssetRepository) Update(ctx context.Context, id int64, a *model.Asset) (*model.Asset, error) {
	query := `
		UPDATE assets
		SET name = $2, category_id = $3, location = $4, condition = $5,
		    photos_docs = $6, is_sharable = $7, is_bookable = $8, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, tag, name, category_id, serial_number, qr_code, status, location,
		          expected_location, condition, photos_docs, custom_field_values,
		          is_sharable, is_bookable, acquisition_date, acquisition_cost,
		          created_at, updated_at, deleted_at`

	return r.scanOne(r.db.QueryRow(ctx, query,
		id, a.Name, a.CategoryID, a.Location, a.Condition,
		a.PhotosDocs, a.IsSharable, a.IsBookable,
	))
}

func (r *AssetRepository) scanOne(row pgx.Row) (*model.Asset, error) {
	var a model.Asset
	err := row.Scan(
		&a.ID, &a.Tag, &a.Name, &a.CategoryID, &a.SerialNumber, &a.QRCode, &a.Status,
		&a.Location, &a.ExpectedLocation, &a.Condition, &a.PhotosDocs,
		&a.CustomFieldValues, &a.IsSharable, &a.IsBookable, &a.AcquisitionDate,
		&a.AcquisitionCost, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAssetNotFound
		}
		return nil, err
	}
	return &a, nil
}