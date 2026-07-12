package model

import "time"

type User struct {
	ID           int64      `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	Email        string     `json:"email" db:"email"`
	Password     string     `json:"-" db:"password"`
	Role         string     `json:"role" db:"role"`
	Gender       string     `json:"gender" db:"gender"`
	DepartmentID *int64     `json:"department_id,omitempty" db:"department_id"`
	Status       string     `json:"status" db:"status"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Role constants mirror the user_role Postgres enum.
const (
	RoleEmployee       = "EMPLOYEE"
	RoleDepartmentHead = "DEPARTMENT_HEAD"
	RoleAssetManager   = "ASSET_MANAGER"
	RoleAdmin          = "ADMIN"
)

// Status constants mirror the user_status Postgres enum.
const (
	StatusActive   = "ACTIVE"
	StatusInactive = "INACTIVE"
)