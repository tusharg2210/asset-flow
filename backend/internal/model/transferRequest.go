package model

import "time"

type TransferRequest struct {
	ID          int64      `json:"id" db:"id"`
	AssetID     int64      `json:"asset_id" db:"asset_id"`
	FromUserID  int64      `json:"from_user_id" db:"from_user_id"`
	ToUserID    int64      `json:"to_user_id" db:"to_user_id"`
	RequestedBy int64      `json:"requested_by" db:"requested_by"`
	ApprovedBy  *int64     `json:"approved_by,omitempty" db:"approved_by"`
	Status      string     `json:"status" db:"status"`
	Reason      string     `json:"reason" db:"reason"`
	RequestedAt time.Time  `json:"requested_at" db:"requested_at"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty" db:"approved_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}


const (
	TransferRequested = "REQUESTED"
	TransferApproved  = "APPROVED"
	TransferRejected  = "REJECTED"
)