package model

import "time"

type Audit struct {
	ID                int64      `json:"id" db:"id"`
	Auditors          []int64    `json:"auditors" db:"auditors"`
	Status            string     `json:"status" db:"status"`
	Scope             string     `json:"scope" db:"scope"`
	ScopeDepartmentID *int64     `json:"scope_department_id,omitempty" db:"scope_department_id"`
	ScopeLocation     string     `json:"scope_location" db:"scope_location"`
	FromDate          time.Time  `json:"from_date" db:"from_date"`
	ToDate            time.Time  `json:"to_date" db:"to_date"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Audit cycle statuses, mirroring the audit_cycle_status enum.
const (
	AuditDraft  = "DRAFT"
	AuditOpen   = "OPEN"
	AuditClosed = "CLOSED"
)

type AuditReport struct {
	ID                 int64      `json:"id" db:"id"`
	AuditID            int64      `json:"audit_id" db:"audit_id"`
	AssetID            int64      `json:"asset_id" db:"asset_id"`
	VerificationStatus string     `json:"verification_status" db:"verification_status"`
	Remarks            string     `json:"remarks" db:"remarks"`
	VerifiedBy         *int64     `json:"verified_by,omitempty" db:"verified_by"`
	VerifiedAt         *time.Time `json:"verified_at,omitempty" db:"verified_at"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Verification statuses, mirroring the audit_verification_status enum.
const (
	VerificationPending  = "PENDING"
	VerificationVerified = "VERIFIED"
	VerificationMissing  = "MISSING"
	VerificationDamaged  = "DAMAGED"
)