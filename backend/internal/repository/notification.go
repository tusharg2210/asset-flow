package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"asset-flow/internal/model"
)

var ErrNotificationNotFound = errors.New("notification not found")

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{db: db}
}


func (r *NotificationRepository) Create(ctx context.Context, tx pgx.Tx, n *model.Notification) error {
	query := `
		INSERT INTO notifications (user_id, type, message, entity_type, entity_id, is_read)
		VALUES ($1, $2, $3, $4, $5, false)
		RETURNING id, created_at`

	if tx != nil {
		return tx.QueryRow(ctx, query, n.UserID, n.Type, n.Message, n.EntityType, n.EntityID).
			Scan(&n.ID, &n.CreatedAt)
	}
	return r.db.QueryRow(ctx, query, n.UserID, n.Type, n.Message, n.EntityType, n.EntityID).
		Scan(&n.ID, &n.CreatedAt)
}

func (r *NotificationRepository) ListByUser(ctx context.Context, userID int64, unreadOnly bool, page, limit int) ([]model.Notification, int, error) {
	where := "user_id = $1"
	args := []any{userID}
	if unreadOnly {
		where += " AND is_read = false"
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM notifications WHERE "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, type, message, entity_type, entity_id, is_read, created_at
		FROM notifications
		WHERE `+where+`
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var n model.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Message, &n.EntityType, &n.EntityID, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, 0, err
		}
		notifications = append(notifications, n)
	}
	return notifications, total, rows.Err()
}

func (r *NotificationRepository) MarkRead(ctx context.Context, id, userID int64) error {
	cmd, err := r.db.Exec(ctx, `
		UPDATE notifications SET is_read = true
		WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return ErrNotificationNotFound
	}
	return nil
}