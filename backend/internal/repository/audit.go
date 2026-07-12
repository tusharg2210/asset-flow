package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"asset-flow/internal/model"
)

var (
	ErrAuditNotFound   = errors.New("audit not found")
	ErrAuditNotOpen    = errors.New("audit is not open")
)

type AuditRepository struct {
	db *pgxpool.Pool
}

func NewAuditRepository(db *pgxpool.Pool) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(ctx context.Context, a *model.Audit) error {
	query := `
		INSERT INTO audits (auditors, status, scope, scope_department_id, scope_location, from_date, to_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		a.Auditors, model.AuditOpen, a.Scope, a.ScopeDepartmentID, a.ScopeLocation, a.FromDate, a.ToDate,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

func (r *AuditRepository) GetByID(ctx context.Context, id int64) (*model.Audit, error) {
	query := `
		SELECT id, auditors, status, scope, scope_department_id, scope_location,
		       from_date, to_date, created_at, updated_at, deleted_at
		FROM audits
		WHERE id = $1 AND deleted_at IS NULL`

	return r.scanOne(r.db.QueryRow(ctx, query, id))
}

func (r *AuditRepository) GetByIDForUpdate(ctx context.Context, tx pgx.Tx, id int64) (*model.Audit, error) {
	query := `
		SELECT id, auditors, status, scope, scope_department_id, scope_location,
		       from_date, to_date, created_at, updated_at, deleted_at
		FROM audits
		WHERE id = $1 AND deleted_at IS NULL
		FOR UPDATE`

	return r.scanOne(tx.QueryRow(ctx, query, id))
}

func (r *AuditRepository) ScopedAssetIDs(ctx context.Context, audit *model.Audit) ([]model.Asset, error) {
	query := `
		SELECT DISTINCT a.id, a.tag, a.name, a.category_id, a.serial_number, a.qr_code,
		       a.status, a.location, a.expected_location, a.condition, a.photos_docs,
		       a.custom_field_values, a.is_sharable, a.is_bookable, a.acquisition_date,
		       a.acquisition_cost, a.created_at, a.updated_at, a.deleted_at
		FROM assets a
		LEFT JOIN allocations al ON al.asset_id = a.id AND al.actual_return_date IS NULL AND al.deleted_at IS NULL
		LEFT JOIN users u ON u.id = al.to_user_id
		WHERE a.deleted_at IS NULL
		  AND ($1::bigint IS NULL OR u.department_id = $1)
		  AND ($2 = '' OR a.location = $2)
		ORDER BY a.id`

	rows, err := r.db.Query(ctx, query, audit.ScopeDepartmentID, audit.ScopeLocation)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []model.Asset
	for rows.Next() {
		var a model.Asset
		if err := rows.Scan(
			&a.ID, &a.Tag, &a.Name, &a.CategoryID, &a.SerialNumber, &a.QRCode, &a.Status,
			&a.Location, &a.ExpectedLocation, &a.Condition, &a.PhotosDocs,
			&a.CustomFieldValues, &a.IsSharable, &a.IsBookable, &a.AcquisitionDate,
			&a.AcquisitionCost, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
		); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, rows.Err()
}


func (r *AuditRepository) CreateReport(ctx context.Context, rep *model.AuditReport) error {
	query := `
		INSERT INTO audit_reports (audit_id, asset_id, verification_status, remarks, verified_by, verified_at)
		VALUES ($1, $2, $3, $4, $5, now())
		ON CONFLICT (audit_id, asset_id) WHERE deleted_at IS NULL
		DO UPDATE SET verification_status = EXCLUDED.verification_status,
		              remarks = EXCLUDED.remarks,
		              verified_by = EXCLUDED.verified_by,
		              verified_at = now(),
		              updated_at = now()
		RETURNING id, verified_at, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		rep.AuditID, rep.AssetID, rep.VerificationStatus, rep.Remarks, rep.VerifiedBy,
	).Scan(&rep.ID, &rep.VerifiedAt, &rep.CreatedAt, &rep.UpdatedAt)
}


func (r *AuditRepository) ListReports(ctx context.Context, auditID int64) ([]model.AuditReport, error) {
	query := `
		SELECT id, audit_id, asset_id, verification_status, remarks, verified_by,
		       verified_at, created_at, updated_at, deleted_at
		FROM audit_reports
		WHERE audit_id = $1 AND deleted_at IS NULL
		ORDER BY verified_at DESC NULLS LAST`

	rows, err := r.db.Query(ctx, query, auditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []model.AuditReport
	for rows.Next() {
		var rep model.AuditReport
		if err := rows.Scan(
			&rep.ID, &rep.AuditID, &rep.AssetID, &rep.VerificationStatus, &rep.Remarks,
			&rep.VerifiedBy, &rep.VerifiedAt, &rep.CreatedAt, &rep.UpdatedAt, &rep.DeletedAt,
		); err != nil {
			return nil, err
		}
		reports = append(reports, rep)
	}
	return reports, rows.Err()
}


func (r *AuditRepository) ListFlaggedAssetIDs(ctx context.Context, tx pgx.Tx, auditID int64, status string) ([]int64, error) {
	rows, err := tx.Query(ctx, `
		SELECT asset_id FROM audit_reports
		WHERE audit_id = $1 AND verification_status = $2 AND deleted_at IS NULL`,
		auditID, status,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}


func (r *AuditRepository) Close(ctx context.Context, tx pgx.Tx, id int64) (*model.Audit, error) {
	query := `
		UPDATE audits
		SET status = $2, updated_at = now()
		WHERE id = $1 AND status = $3 AND deleted_at IS NULL
		RETURNING id, auditors, status, scope, scope_department_id, scope_location,
		          from_date, to_date, created_at, updated_at, deleted_at`

	a, err := r.scanOne(tx.QueryRow(ctx, query, id, model.AuditClosed, model.AuditOpen))
	if err != nil {
		if errors.Is(err, ErrAuditNotFound) {
			return nil, ErrAuditNotOpen
		}
		return nil, err
	}
	return a, nil
}

func (r *AuditRepository) scanOne(row pgx.Row) (*model.Audit, error) {
	var a model.Audit
	err := row.Scan(
		&a.ID, &a.Auditors, &a.Status, &a.Scope, &a.ScopeDepartmentID, &a.ScopeLocation,
		&a.FromDate, &a.ToDate, &a.CreatedAt, &a.UpdatedAt, &a.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAuditNotFound
		}
		return nil, err
	}
	return &a, nil
}