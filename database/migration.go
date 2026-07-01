package database

import (
	"fmt"

	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) error {
	fmt.Println("Running migrations...")

	if err := db.AutoMigrate(
		&User{},
		&LearningNote{},
		&Todo{},
	); err != nil {
		return fmt.Errorf("gagal migrasi: %w", err)
	}

	fmt.Println("Migrations completed")

	fmt.Println("Seeding database...")
	if err := Seed(db); err != nil {
		return fmt.Errorf("gagal seeding: %w", err)
	}

	return nil
}
