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
	ErrCategoryExists   = errors.New("category name already exists")
	ErrCategoryNotFound = errors.New("asset category not found")
)

type AssetCategoryRepository struct {
	db *pgxpool.Pool
}

func NewAssetCategoryRepository(db *pgxpool.Pool) *AssetCategoryRepository {
	return &AssetCategoryRepository{db: db}
}

func (r *AssetCategoryRepository) Create(ctx context.Context, c *model.AssetCategory) error {
	query := `
		INSERT INTO asset_categories (name, custom_fields_schema)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query, c.Name, c.CustomFieldsSchema).
		Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrCategoryExists
		}
		return err
	}

	return nil
}

func (r *AssetCategoryRepository) GetByID(ctx context.Context, id int64) (*model.AssetCategory, error) {
	query := `
		SELECT id, name, custom_fields_schema, created_at, updated_at, deleted_at
		FROM asset_categories
		WHERE id = $1 AND deleted_at IS NULL`

	return r.scanOne(r.db.QueryRow(ctx, query, id))
}

func (r *AssetCategoryRepository) List(ctx context.Context) ([]model.AssetCategory, error) {
	query := `
		SELECT id, name, custom_fields_schema, created_at, updated_at, deleted_at
		FROM asset_categories
		WHERE deleted_at IS NULL
		ORDER BY name ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.AssetCategory
	for rows.Next() {
		var c model.AssetCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.CustomFieldsSchema, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *AssetCategoryRepository) Update(ctx context.Context, id int64, name *string, schema []byte) (*model.AssetCategory, error) {
	query := `
		UPDATE asset_categories
		SET name = COALESCE($2, name),
		    custom_fields_schema = COALESCE($3, custom_fields_schema),
		    updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, name, custom_fields_schema, created_at, updated_at, deleted_at`

	return r.scanOne(r.db.QueryRow(ctx, query, id, name, schema))
}

func (r *AssetCategoryRepository) scanOne(row pgx.Row) (*model.AssetCategory, error) {
	var c model.AssetCategory
	err := row.Scan(&c.ID, &c.Name, &c.CustomFieldsSchema, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	return &c, nil
}