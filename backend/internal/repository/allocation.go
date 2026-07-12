package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"asset-flow/internal/model"
)

var (
	ErrAllocationNotFound  = errors.New("allocation not found")
	ErrAssetAlreadyHeld    = errors.New("asset is already allocated to another user")
	ErrAllocationReturned  = errors.New("allocation already returned")
)

type AllocationRepository struct {
	db *pgxpool.Pool
}

func NewAllocationRepository(db *pgxpool.Pool) *AllocationRepository {
	return &AllocationRepository{db: db}
}

func (r *AllocationRepository) Create(ctx context.Context, tx pgx.Tx, a *model.Allocation) error {
	query := `
		INSERT INTO allocations (asset_id, from_user_id, to_user_id, allotted_date, expected_return_date, reason)
		VALUES ($1, $2, $3, now(), $4, $5)
		RETURNING id, allotted_date, created_at, updated_at`

	return tx.QueryRow(ctx, query,
		a.AssetID, a.FromUserID, a.ToUserID, a.ExpectedReturnDate, a.Reason,
	).Scan(&a.ID, &a.AllottedDate, &a.CreatedAt, &a.UpdatedAt)
}

func (r *AllocationRepository) GetByID(ctx context.Context, id int64) (*model.Allocation, error) {
	query := `
		SELECT id, asset_id, from_user_id, to_user_id, allotted_date, expected_return_date,
		       actual_return_date, reason, created_at, updated_at, deleted_at
		FROM allocations
		WHERE id = $1 AND deleted_at IS NULL`

	return r.scanOne(r.db.QueryRow(ctx, query, id))
}

func (r *AllocationRepository) GetActiveByAssetID(ctx context.Context, tx pgx.Tx, assetID int64) (*model.Allocation, error) {
	query := `
		SELECT id, asset_id, from_user_id, to_user_id, allotted_date, expected_return_date,
		       actual_return_date, reason, created_at, updated_at, deleted_at
		FROM allocations
		WHERE asset_id = $1 AND actual_return_date IS NULL AND deleted_at IS NULL
		FOR UPDATE`

	var row pgx.Row
	if tx != nil {
		row = tx.QueryRow(ctx, query, assetID)
	} else {
		row = r.db.QueryRow(ctx, query, assetID)
	}
	return r.scanOne(row)
}

func (r *AllocationRepository) ListByAsset(ctx context.Context, assetID int64) ([]model.Allocation, error) {
	query := `
		SELECT id, asset_id, from_user_id, to_user_id, allotted_date, expected_return_date,
		       actual_return_date, reason, created_at, updated_at, deleted_at
		FROM allocations
		WHERE asset_id = $1 AND deleted_at IS NULL
		ORDER BY allotted_date DESC`

	rows, err := r.db.Query(ctx, query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allocations []model.Allocation
	for rows.Next() {
		var a model.Allocation
		if err := rows.Scan(
			&a.ID, &a.AssetID, &a.FromUserID, &a.ToUserID, &a.AllottedDate,
			&a.ExpectedReturnDate, &a.ActualReturnDate, &a.Reason,
			&a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
		); err != nil {
			return nil, err
		}
		allocations = append(allocations, a)
	}
	return allocations, rows.Err()
}

func (r *AllocationRepository) ListOverdue(ctx context.Context) ([]model.Allocation, error) {
	query := `
		SELECT id, asset_id, from_user_id, to_user_id, allotted_date, expected_return_date,
		       actual_return_date, reason, created_at, updated_at, deleted_at
		FROM allocations
		WHERE actual_return_date IS NULL
		  AND expected_return_date IS NOT NULL
		  AND expected_return_date < now()
		  AND deleted_at IS NULL
		ORDER BY expected_return_date ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allocations []model.Allocation
	for rows.Next() {
		var a model.Allocation
		if err := rows.Scan(
			&a.ID, &a.AssetID, &a.FromUserID, &a.ToUserID, &a.AllottedDate,
			&a.ExpectedReturnDate, &a.ActualReturnDate, &a.Reason,
			&a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
		); err != nil {
			return nil, err
		}
		allocations = append(allocations, a)
	}
	return allocations, rows.Err()
}

func (r *AllocationRepository) MarkReturned(ctx context.Context, tx pgx.Tx, id int64, conditionNotes string) (*model.Allocation, error) {
	query := `
		UPDATE allocations
		SET actual_return_date = now(), reason = COALESCE(reason, '') , updated_at = now()
		WHERE id = $1 AND actual_return_date IS NULL AND deleted_at IS NULL
		RETURNING id, asset_id, from_user_id, to_user_id, allotted_date, expected_return_date,
		          actual_return_date, reason, created_at, updated_at, deleted_at`

	_ = conditionNotes

	a, err := r.scanOne(tx.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, ErrAllocationNotFound) {
			return nil, ErrAllocationReturned
		}
		return nil, err
	}
	return a, nil
}

func (r *AllocationRepository) scanOne(row pgx.Row) (*model.Allocation, error) {
	var a model.Allocation
	err := row.Scan(
		&a.ID, &a.AssetID, &a.FromUserID, &a.ToUserID, &a.AllottedDate,
		&a.ExpectedReturnDate, &a.ActualReturnDate, &a.Reason,
		&a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAllocationNotFound
		}
		return nil, err
	}
	return &a, nil
}