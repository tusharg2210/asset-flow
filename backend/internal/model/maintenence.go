package model

import "time"

type Maintenance struct {
	ID                   int64      `json:"id" db:"id"`
	AssetID              int64      `json:"asset_id" db:"asset_id"`
	RaisedBy             *int64     `json:"raised_by,omitempty" db:"raised_by"`
	AssignedTechnicianID *int64     `json:"assigned_technician_id,omitempty" db:"assigned_technician_id"`
	Priority             string     `json:"priority" db:"priority"`
	Description          string     `json:"description" db:"description"`
	Images               []string   `json:"images" db:"images"`
	Status               string     `json:"status" db:"status"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Priorities, mirroring the maintenance_priority enum.
const (
	PriorityLow      = "LOW"
	PriorityMedium   = "MEDIUM"
	PriorityHigh     = "HIGH"
	PriorityCritical = "CRITICAL"
)

// Maintenance workflow states, mirroring the maintenance_status enum.
const (
	MaintenancePending             = "PENDING"
	MaintenanceApproved            = "APPROVED"
	MaintenanceRejected            = "REJECTED"
	MaintenanceTechnicianAssigned  = "TECHNICIAN_ASSIGNED"
	MaintenanceInProgress          = "IN_PROGRESS"
	MaintenanceResolved            = "RESOLVED"
)

type MaintenanceWorkflow struct {
	ID            int64      `json:"id" db:"id"`
	MaintenanceID int64      `json:"maintenance_id" db:"maintenance_id"`
	Status        string     `json:"status" db:"status"`
	Description   string     `json:"description" db:"description"`
	UpdatedBy     *int64     `json:"updated_by,omitempty" db:"updated_by"`
	WorkflowDate  time.Time  `json:"workflow_date" db:"workflow_date"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}