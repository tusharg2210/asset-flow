package repository

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"asset-flow/internal/model"
)

type ActivityLogRepository struct {
	db *pgxpool.Pool
}

func NewActivityLogRepository(db *pgxpool.Pool) *ActivityLogRepository {
	return &ActivityLogRepository{db: db}
}

func (r *ActivityLogRepository) Create(ctx context.Context, tx pgx.Tx, l *model.ActivityLog) error {
	query := `
		INSERT INTO activity_logs (user_id, action, entity_type, entity_id, metadata)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	if tx != nil {
		return tx.QueryRow(ctx, query, l.UserID, l.Action, l.EntityType, l.EntityID, l.Metadata).
			Scan(&l.ID, &l.CreatedAt)
	}
	return r.db.QueryRow(ctx, query, l.UserID, l.Action, l.EntityType, l.EntityID, l.Metadata).
		Scan(&l.ID, &l.CreatedAt)
}

type ActivityLogFilter struct {
	UserID     *int64
	EntityType string
	Page       int
	Limit      int
}

func (r *ActivityLogRepository) List(ctx context.Context, f ActivityLogFilter) ([]model.ActivityLog, int, error) {
	where := "1=1"
	var args []any

	if f.UserID != nil {
		args = append(args, *f.UserID)
		where += " AND user_id = $" + itoa(len(args))
	}
	if f.EntityType != "" {
		args = append(args, f.EntityType)
		where += " AND entity_type = $" + itoa(len(args))
	}

	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 {
		f.Limit = 50
	}
	offset := (f.Page - 1) * f.Limit

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM activity_logs WHERE "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, f.Limit, offset)
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, action, entity_type, entity_id, metadata, created_at
		FROM activity_logs
		WHERE `+where+`
		ORDER BY created_at DESC
		LIMIT $`+itoa(len(args)-1)+` OFFSET $`+itoa(len(args)), args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []model.ActivityLog
	for rows.Next() {
		var l model.ActivityLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.Action, &l.EntityType, &l.EntityID, &l.Metadata, &l.CreatedAt); err != nil {
			return nil, 0, err
		}
		logs = append(logs, l)
	}
	return logs, total, rows.Err()
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
