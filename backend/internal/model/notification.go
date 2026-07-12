package model

import "time"

type Notification struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	Type       string    `json:"type" db:"type"`
	Message    string    `json:"message" db:"message"`
	EntityType string    `json:"entity_type" db:"entity_type"`
	EntityID   *int64    `json:"entity_id,omitempty" db:"entity_id"`
	IsRead     bool      `json:"is_read" db:"is_read"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

const (
	NotifAssetAssigned         = "ASSET_ASSIGNED"
	NotifMaintenanceApproved   = "MAINTENANCE_APPROVED"
	NotifMaintenanceRejected   = "MAINTENANCE_REJECTED"
	NotifBookingConfirmed      = "BOOKING_CONFIRMED"
	NotifBookingCancelled      = "BOOKING_CANCELLED"
	NotifBookingReminder       = "BOOKING_REMINDER"
	NotifTransferApproved      = "TRANSFER_APPROVED"
	NotifOverdueReturn         = "OVERDUE_RETURN"
	NotifAuditDiscrepancy      = "AUDIT_DISCREPANCY"
)