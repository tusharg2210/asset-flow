package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"asset-flow/internal/config"
	"asset-flow/internal/db"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Pool.Close()

	ctx := context.Background()
	log.Println("Connected to database, starting to seed 100 extra assets...")

	var newCatID int64
	err = database.Pool.QueryRow(ctx, "INSERT INTO asset_categories (name) VALUES ('Bulk Test Category') RETURNING id").Scan(&newCatID)
	if err != nil {
		log.Fatalf("failed to create category: %v", err)
	}

	log.Printf("Created new category with ID: %d", newCatID)

	for i := 1; i <= 100; i++ {
		tag := fmt.Sprintf("AF-BULK-%03d", i)
		name := fmt.Sprintf("Bulk Test Asset %d", i)
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

		_, err = database.Pool.Exec(ctx, `
			INSERT INTO assets (tag, name, category_id, status, location, expected_location, condition, is_sharable, is_bookable, acquisition_date, acquisition_cost)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			tag, name, newCatID, status, "HQ Warehouse", "HQ Warehouse", condition, true, isBookable, time.Now().AddDate(0, -1, -i), float64(100+i*10))
		if err != nil {
			log.Fatalf("failed to insert asset %d: %v", i, err)
		}
	}

	log.Println("Successfully seeded 100 assets into Supabase!")
}
