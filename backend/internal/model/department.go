package model

import "time"

type Department struct {
	ID                 int64      `json:"id" db:"id"`
	Name               string     `json:"name" db:"name"`
	ParentDepartmentID *int64     `json:"parent_department_id,omitempty" db:"parent_department_id"`
	HeadID             *int64     `json:"head_id,omitempty" db:"head_id"`
	Status             string     `json:"status" db:"status"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Department reuses model.StatusActive / model.StatusInactive (defined in
// user.go) since department_status has the same two values as user_status.