package model

import "time"

type Allocation struct {
	ID                 int64      `json:"id" db:"id"`
	AssetID            int64      `json:"asset_id" db:"asset_id"`
	FromUserID         *int64     `json:"from_user_id,omitempty" db:"from_user_id"`
	ToUserID           int64      `json:"to_user_id" db:"to_user_id"`
	AllottedDate       time.Time  `json:"allotted_date" db:"allotted_date"`
	ExpectedReturnDate *time.Time `json:"expected_return_date,omitempty" db:"expected_return_date"`
	ActualReturnDate   *time.Time `json:"actual_return_date,omitempty" db:"actual_return_date"`
	Reason             string     `json:"reason" db:"reason"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}