package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repositories struct {
	User         *UserRepository
	Department   *DepartmentRepository
	Category     *AssetCategoryRepository
	Asset        *AssetRepository
	Allocation   *AllocationRepository
	Transfer     *TransferRepository
	Booking      *BookingRepository
	Maintenance  *MaintenanceRepository
	Audit        *AuditRepository
	Notification *NotificationRepository
	ActivityLog  *ActivityLogRepository
}

func New(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		User:         NewUserRepository(pool),
		Department:   NewDepartmentRepository(pool),
		Category:     NewAssetCategoryRepository(pool),
		Asset:        NewAssetRepository(pool),
		Allocation:   NewAllocationRepository(pool),
		Transfer:     NewTransferRepository(pool),
		Booking:      NewBookingRepository(pool),
		Maintenance:  NewMaintenanceRepository(pool),
		Audit:        NewAuditRepository(pool),
		Notification: NewNotificationRepository(pool),
		ActivityLog:  NewActivityLogRepository(pool),
	}
}