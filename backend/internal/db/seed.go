package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func SeedData(ctx context.Context, pool *pgxpool.Pool) error {
	// 1. Check if users already exist
	var count int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		log.Println("Database already seeded, skipping.")
		return nil
	}

	log.Println("Seeding database...")

	// 2. Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("SecurePassword123!"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 3. Create departments
	var engID, hrID, opsID int64
	err = pool.QueryRow(ctx, "INSERT INTO departments (name) VALUES ('Engineering') RETURNING id").Scan(&engID)
	if err != nil {
		return err
	}
	err = pool.QueryRow(ctx, "INSERT INTO departments (name) VALUES ('HR') RETURNING id").Scan(&hrID)
	if err != nil {
		return err
	}
	err = pool.QueryRow(ctx, "INSERT INTO departments (name) VALUES ('Operations') RETURNING id").Scan(&opsID)
	if err != nil {
		return err
	}

	// 4. Create users
	var adminID, managerID, employeeID, headID int64
	err = pool.QueryRow(ctx, `
		INSERT INTO users (name, email, password, role, gender, department_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`, "Admin User", "admin@company.com", string(passwordHash), "ADMIN", "Male", engID, "ACTIVE").Scan(&adminID)
	if err != nil {
		return err
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO users (name, email, password, role, gender, department_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`, "Asset Manager", "manager@company.com", string(passwordHash), "ASSET_MANAGER", "Female", opsID, "ACTIVE").Scan(&managerID)
	if err != nil {
		return err
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO users (name, email, password, role, gender, department_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`, "Employee User", "employee@company.com", string(passwordHash), "EMPLOYEE", "Male", engID, "ACTIVE").Scan(&employeeID)
	if err != nil {
		return err
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO users (name, email, password, role, gender, department_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`, "Department Head", "head@company.com", string(passwordHash), "DEPARTMENT_HEAD", "Female", engID, "ACTIVE").Scan(&headID)
	if err != nil {
		return err
	}

	// 5. Update departments head_id
	_, err = pool.Exec(ctx, "UPDATE departments SET head_id = $1 WHERE id = $2", headID, engID)
	if err != nil {
		return err
	}
	_, err = pool.Exec(ctx, "UPDATE departments SET head_id = $1 WHERE id = $2", managerID, opsID)
	if err != nil {
		return err
	}

	// 6. Create categories
	var elecID, furnID, vehID int64
	err = pool.QueryRow(ctx, "INSERT INTO asset_categories (name) VALUES ('Electronics') RETURNING id").Scan(&elecID)
	if err != nil {
		return err
	}
	err = pool.QueryRow(ctx, "INSERT INTO asset_categories (name) VALUES ('Furniture') RETURNING id").Scan(&furnID)
	if err != nil {
		return err
	}
	err = pool.QueryRow(ctx, "INSERT INTO asset_categories (name) VALUES ('Vehicles') RETURNING id").Scan(&vehID)
	if err != nil {
		return err
	}

	// 7. Create assets
	var mProID, roomID, vanID, monitorID, deskID int64
	err = pool.QueryRow(ctx, `
		INSERT INTO assets (tag, name, category_id, status, location, expected_location, condition, is_sharable, is_bookable, acquisition_date, acquisition_cost)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`, "AF-0001", "MacBook Pro 16", elecID, "AVAILABLE", "HQ - Floor 2", "HQ - Floor 2", "New", true, false, time.Now().AddDate(0, -6, 0), 2400.00).Scan(&mProID)
	if err != nil {
		return err
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO assets (tag, name, category_id, status, location, expected_location, condition, is_sharable, is_bookable, acquisition_date, acquisition_cost)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`, "AF-0002", "Conference Room B2", elecID, "AVAILABLE", "HQ - Floor 1", "HQ - Floor 1", "Excellent", true, true, time.Now().AddDate(-1, 0, 0), 0.00).Scan(&roomID)
	if err != nil {
		return err
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO assets (tag, name, category_id, status, location, expected_location, condition, is_sharable, is_bookable, acquisition_date, acquisition_cost)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`, "AF-0003", "Company Van", vehID, "AVAILABLE", "HQ Garage", "HQ Garage", "Good", true, true, time.Now().AddDate(-2, 0, 0), 35000.00).Scan(&vanID)
	if err != nil {
		return err
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO assets (tag, name, category_id, status, location, expected_location, condition, is_sharable, is_bookable, acquisition_date, acquisition_cost)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`, "AF-0004", "Dell Monitor 27", elecID, "ALLOCATED", "HQ - Floor 2", "HQ - Floor 2", "Good", false, false, time.Now().AddDate(0, -3, 0), 400.00).Scan(&monitorID)
	if err != nil {
		return err
	}

	err = pool.QueryRow(ctx, `
		INSERT INTO assets (tag, name, category_id, status, location, expected_location, condition, is_sharable, is_bookable, acquisition_date, acquisition_cost)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`, "AF-0005", "Office Desk", furnID, "AVAILABLE", "HQ - Floor 2", "HQ - Floor 2", "New", false, false, time.Now().AddDate(0, -9, 0), 300.00).Scan(&deskID)
	if err != nil {
		return err
	}

	// Generate 100 extra assets for testing
	var newCatID int64
	err = pool.QueryRow(ctx, "INSERT INTO asset_categories (name) VALUES ('Test Category') RETURNING id").Scan(&newCatID)
	if err != nil {
		return err
	}
	
	for i := 1; i <= 100; i++ {
		tag := fmt.Sprintf("AF-TEST-%03d", i)
		name := fmt.Sprintf("Test Asset %d", i)
		condition := "New"
		if i%2 == 0 {
			condition = "Good"
		} else if i%3 == 0 {
			condition = "Fair"
		}

		status := "AVAILABLE"
		var isBookable bool
		if i%4 == 0 {
			status = "ALLOCATED"
			isBookable = false
		} else if i%5 == 0 {
			isBookable = true
		}

		_, err = pool.Exec(ctx, `
			INSERT INTO assets (tag, name, category_id, status, location, expected_location, condition, is_sharable, is_bookable, acquisition_date, acquisition_cost)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`, 
			tag, name, newCatID, status, "Warehouse", "Warehouse", condition, true, isBookable, time.Now().AddDate(0, -1, -i), float64(100 + i*10))
		if err != nil {
			return err
		}
	}

	// 8. Create allocations
	_, err = pool.Exec(ctx, `
		INSERT INTO allocations (asset_id, from_user_id, to_user_id, reason, expected_return_date)
		VALUES ($1, $2, $3, $4, $5)`, monitorID, managerID, employeeID, "Initial IT Allocation", time.Now().AddDate(0, 6, 0))
	if err != nil {
		return err
	}

	// 9. Create bookings
	var bookingID int64
	err = pool.QueryRow(ctx, `
		INSERT INTO resource_bookings (asset_id, booked_by, booking_status)
		VALUES ($1, $2, $3)
		RETURNING id`, roomID, employeeID, "UPCOMING").Scan(&bookingID)
	if err != nil {
		return err
	}

	// Create booking slots for room
	startTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 14, 0, 0, 0, time.Local)
	endTime := startTime.Add(1 * time.Hour)
	_, err = pool.Exec(ctx, `
		INSERT INTO booking_slots (booking_id, start_time, end_time)
		VALUES ($1, $2, $3)`, bookingID, startTime, endTime)
	if err != nil {
		return err
	}

	// 10. Create maintenance request
	_, err = pool.Exec(ctx, `
		INSERT INTO maintenance (asset_id, raised_by, priority, description, status)
		VALUES ($1, $2, $3, $4, $5)`, monitorID, employeeID, "MEDIUM", "Screen flickering occasionally", "PENDING")
	if err != nil {
		return err
	}

	// 11. Create activity logs
	_, err = pool.Exec(ctx, `
		INSERT INTO activity_logs (user_id, action, entity_type, entity_id)
		VALUES ($1, $2, $3, $4)`, adminID, "Created asset MacBook Pro 16", "ASSET", mProID)
	if err != nil {
		return err
	}
	_, err = pool.Exec(ctx, `
		INSERT INTO activity_logs (user_id, action, entity_type, entity_id)
		VALUES ($1, $2, $3, $4)`, managerID, "Allocated Dell Monitor 27 to Employee User", "ALLOCATION", monitorID)
	if err != nil {
		return err
	}

	// 12. Create open audit cycle
	var auditID int64
	err = pool.QueryRow(ctx, `
		INSERT INTO audits (auditors, status, scope, scope_department_id, scope_location, from_date, to_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`, []int64{adminID}, "OPEN", "Engineering Department Assets", engID, "HQ - Floor 2", time.Now().AddDate(0, 0, -5), time.Now().AddDate(0, 0, 5)).Scan(&auditID)
	if err != nil {
		return err
	}

	log.Println("Database successfully seeded.")
	return nil
}
