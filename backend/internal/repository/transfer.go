package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"asset-flow/internal/model"
)

var (
	ErrTransferNotFound     = errors.New("transfer request not found")
	ErrTransferNotRequested = errors.New("transfer request is not in REQUESTED state")
)

type TransferRepository struct {
	db *pgxpool.Pool
}

func NewTransferRepository(db *pgxpool.Pool) *TransferRepository {
	return &TransferRepository{db: db}
}

func (r *TransferRepository) Create(ctx context.Context, t *model.TransferRequest) error {
	query := `
		INSERT INTO transfer_requests (asset_id, from_user_id, to_user_id, requested_by, status, reason, requested_at)
		VALUES ($1, $2, $3, $4, $5, $6, now())
		RETURNING id, requested_at, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		t.AssetID, t.FromUserID, t.ToUserID, t.RequestedBy, model.TransferRequested, t.Reason,
	).Scan(&t.ID, &t.RequestedAt, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TransferRepository) GetByID(ctx context.Context, id int64) (*model.TransferRequest, error) {
	query := `
		SELECT id, asset_id, from_user_id, to_user_id, requested_by, approved_by, status,
		       reason, requested_at, approved_at, created_at, updated_at, deleted_at
		FROM transfer_requests
		WHERE id = $1 AND deleted_at IS NULL`

	return r.scanOne(r.db.QueryRow(ctx, query, id))
}


func (r *TransferRepository) GetByIDForUpdate(ctx context.Context, tx pgx.Tx, id int64) (*model.TransferRequest, error) {
	query := `
		SELECT id, asset_id, from_user_id, to_user_id, requested_by, approved_by, status,
		       reason, requested_at, approved_at, created_at, updated_at, deleted_at
		FROM transfer_requests
		WHERE id = $1 AND deleted_at IS NULL
		FOR UPDATE`

	return r.scanOne(tx.QueryRow(ctx, query, id))
}

func (r *TransferRepository) ListByAsset(ctx context.Context, assetID int64) ([]model.TransferRequest, error) {
	query := `
		SELECT id, asset_id, from_user_id, to_user_id, requested_by, approved_by, status,
		       reason, requested_at, approved_at, created_at, updated_at, deleted_at
		FROM transfer_requests
		WHERE asset_id = $1 AND deleted_at IS NULL
		ORDER BY requested_at DESC`

	rows, err := r.db.Query(ctx, query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []model.TransferRequest
	for rows.Next() {
		var t model.TransferRequest
		if err := rows.Scan(
			&t.ID, &t.AssetID, &t.FromUserID, &t.ToUserID, &t.RequestedBy, &t.ApprovedBy,
			&t.Status, &t.Reason, &t.RequestedAt, &t.ApprovedAt,
			&t.CreatedAt, &t.UpdatedAt, &t.DeletedAt,
		); err != nil {
			return nil, err
		}
		transfers = append(transfers, t)
	}
	return transfers, rows.Err()
}

func (r *TransferRepository) CountPending(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM transfer_requests
		WHERE status = $1 AND deleted_at IS NULL`, model.TransferRequested,
	).Scan(&count)
	return count, err
}


func (r *TransferRepository) UpdateStatus(ctx context.Context, tx pgx.Tx, id int64, status string, approvedBy int64) (*model.TransferRequest, error) {
	query := `
		UPDATE transfer_requests
		SET status = $2, approved_by = $3, approved_at = now(), updated_at = now()
		WHERE id = $1 AND status = $4 AND deleted_at IS NULL
		RETURNING id, asset_id, from_user_id, to_user_id, requested_by, approved_by, status,
		          reason, requested_at, approved_at, created_at, updated_at, deleted_at`

	t, err := r.scanOne(tx.QueryRow(ctx, query, id, status, approvedBy, model.TransferRequested))
	if err != nil {
		if errors.Is(err, ErrTransferNotFound) {
			return nil, ErrTransferNotRequested
		}
		return nil, err
	}
	return t, nil
}

func (r *TransferRepository) scanOne(row pgx.Row) (*model.TransferRequest, error) {
	var t model.TransferRequest
	err := row.Scan(
		&t.ID, &t.AssetID, &t.FromUserID, &t.ToUserID, &t.RequestedBy, &t.ApprovedBy,
		&t.Status, &t.Reason, &t.RequestedAt, &t.ApprovedAt,
		&t.CreatedAt, &t.UpdatedAt, &t.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTransferNotFound
		}
		return nil, err
	}
	return &t, nil
}