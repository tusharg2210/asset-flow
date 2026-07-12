package model

import "time"

type Asset struct {
	ID                int64      `json:"id" db:"id"`
	Tag               string     `json:"tag" db:"tag"`
	Name              string     `json:"name" db:"name"`
	CategoryID        *int64     `json:"category_id,omitempty" db:"category_id"`
	SerialNumber      string     `json:"serial_number" db:"serial_number"`
	QRCode            string     `json:"qr_code" db:"qr_code"`
	Status            string     `json:"status" db:"status"`
	Location          string     `json:"location" db:"location"`
	ExpectedLocation  string     `json:"expected_location" db:"expected_location"`
	Condition         string     `json:"condition" db:"condition"`
	PhotosDocs        []string   `json:"photos_docs" db:"photos_docs"`
	CustomFieldValues []byte     `json:"custom_field_values" db:"custom_field_values"`
	IsSharable        bool       `json:"is_sharable" db:"is_sharable"`
	IsBookable        bool       `json:"is_bookable" db:"is_bookable"`
	AcquisitionDate   *time.Time `json:"acquisition_date,omitempty" db:"acquisition_date"`
	AcquisitionCost   float64    `json:"acquisition_cost" db:"acquisition_cost"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

const (
	AssetAvailable        = "AVAILABLE"
	AssetAllocated        = "ALLOCATED"
	AssetReserved         = "RESERVED"
	AssetUnderMaintenance = "UNDER_MAINTENANCE"
	AssetLost             = "LOST"
	AssetRetired          = "RETIRED"
	AssetDisposed         = "DISPOSED"
)

type AssetCategory struct {
	ID                 int64      `json:"id" db:"id"`
	Name               string     `json:"name" db:"name"`
	CustomFieldsSchema []byte     `json:"custom_fields_schema" db:"custom_fields_schema"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}