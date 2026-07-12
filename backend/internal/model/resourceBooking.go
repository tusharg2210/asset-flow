package model

import "time"

type ResourceBooking struct {
	ID            int64      `json:"id" db:"id"`
	AssetID       int64      `json:"asset_id" db:"asset_id"`
	BookedBy      int64      `json:"booked_by" db:"booked_by"`
	BookingStatus string     `json:"booking_status" db:"booking_status"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

const (
	BookingUpcoming  = "UPCOMING"
	BookingOngoing   = "ONGOING"
	BookingCompleted = "COMPLETED"
	BookingCancelled = "CANCELLED"
)

type BookingSlot struct {
	ID        int64      `json:"id" db:"id"`
	BookingID int64      `json:"booking_id" db:"booking_id"`
	StartTime time.Time  `json:"start_time" db:"start_time"`
	EndTime   time.Time  `json:"end_time" db:"end_time"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}