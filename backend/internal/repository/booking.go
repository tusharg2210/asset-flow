package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"asset-flow/internal/model"
)

var (
	ErrBookingNotFound     = errors.New("booking not found")
	ErrBookingSlotOverlap  = errors.New("requested time slot overlaps an existing booking")
)

type BookingRepository struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) *BookingRepository {
	return &BookingRepository{db: db}
}


func (r *BookingRepository) HasOverlap(ctx context.Context, tx pgx.Tx, assetID int64, start, end time.Time) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM booking_slots bs
			JOIN resource_bookings rb ON rb.id = bs.booking_id
			WHERE rb.asset_id = $1
			  AND rb.booking_status IN ($2, $3)
			  AND rb.deleted_at IS NULL
			  AND bs.deleted_at IS NULL
			  AND bs.start_time < $5
			  AND bs.end_time > $4
		)`

	var exists bool
	err := tx.QueryRow(ctx, query,
		assetID, model.BookingUpcoming, model.BookingOngoing, start, end,
	).Scan(&exists)
	return exists, err
}


func (r *BookingRepository) Create(ctx context.Context, tx pgx.Tx, b *model.ResourceBooking, slot *model.BookingSlot) error {
	err := tx.QueryRow(ctx, `
		INSERT INTO resource_bookings (asset_id, booked_by, booking_status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`,
		b.AssetID, b.BookedBy, model.BookingUpcoming,
	).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return err
	}
	b.BookingStatus = model.BookingUpcoming

	slot.BookingID = b.ID
	return tx.QueryRow(ctx, `
		INSERT INTO booking_slots (booking_id, start_time, end_time)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`,
		slot.BookingID, slot.StartTime, slot.EndTime,
	).Scan(&slot.ID, &slot.CreatedAt, &slot.UpdatedAt)
}

func (r *BookingRepository) GetByID(ctx context.Context, id int64) (*model.ResourceBooking, error) {
	query := `
		SELECT id, asset_id, booked_by, booking_status, created_at, updated_at, deleted_at
		FROM resource_bookings
		WHERE id = $1 AND deleted_at IS NULL`

	return r.scanOne(r.db.QueryRow(ctx, query, id))
}


func (r *BookingRepository) ListSlotsByAsset(ctx context.Context, assetID int64) ([]model.BookingSlot, error) {
	query := `
		SELECT bs.id, bs.booking_id, bs.start_time, bs.end_time, bs.created_at, bs.updated_at, bs.deleted_at
		FROM booking_slots bs
		JOIN resource_bookings rb ON rb.id = bs.booking_id
		WHERE rb.asset_id = $1
		  AND rb.booking_status IN ($2, $3)
		  AND rb.deleted_at IS NULL
		  AND bs.deleted_at IS NULL
		ORDER BY bs.start_time ASC`

	rows, err := r.db.Query(ctx, query, assetID, model.BookingUpcoming, model.BookingOngoing)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []model.BookingSlot
	for rows.Next() {
		var s model.BookingSlot
		if err := rows.Scan(&s.ID, &s.BookingID, &s.StartTime, &s.EndTime, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt); err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}
	return slots, rows.Err()
}


func (r *BookingRepository) ListUpcomingWithin(ctx context.Context, window time.Duration) ([]model.BookingSlot, error) {
	query := `
		SELECT bs.id, bs.booking_id, bs.start_time, bs.end_time, bs.created_at, bs.updated_at, bs.deleted_at
		FROM booking_slots bs
		JOIN resource_bookings rb ON rb.id = bs.booking_id
		WHERE rb.booking_status = $1
		  AND rb.deleted_at IS NULL
		  AND bs.deleted_at IS NULL
		  AND bs.start_time BETWEEN now() AND now() + $2::interval`

	rows, err := r.db.Query(ctx, query, model.BookingUpcoming, window.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []model.BookingSlot
	for rows.Next() {
		var s model.BookingSlot
		if err := rows.Scan(&s.ID, &s.BookingID, &s.StartTime, &s.EndTime, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt); err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}
	return slots, rows.Err()
}

func (r *BookingRepository) CountActive(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM resource_bookings
		WHERE booking_status IN ($1, $2) AND deleted_at IS NULL`,
		model.BookingUpcoming, model.BookingOngoing,
	).Scan(&count)
	return count, err
}

func (r *BookingRepository) UpdateStatus(ctx context.Context, id int64, status string) (*model.ResourceBooking, error) {
	query := `
		UPDATE resource_bookings
		SET booking_status = $2, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, asset_id, booked_by, booking_status, created_at, updated_at, deleted_at`

	return r.scanOne(r.db.QueryRow(ctx, query, id, status))
}

func (r *BookingRepository) scanOne(row pgx.Row) (*model.ResourceBooking, error) {
	var b model.ResourceBooking
	err := row.Scan(&b.ID, &b.AssetID, &b.BookedBy, &b.BookingStatus, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBookingNotFound
		}
		return nil, err
	}
	return &b, nil
}