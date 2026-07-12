package model

import "time"

type ActivityLog struct {
	ID         int64     `json:"id" db:"id"`
	UserID     *int64    `json:"user_id,omitempty" db:"user_id"`
	Action     string    `json:"action" db:"action"`
	EntityType string    `json:"entity_type" db:"entity_type"`
	EntityID   *int64    `json:"entity_id,omitempty" db:"entity_id"`
	Metadata   []byte    `json:"metadata" db:"metadata"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}