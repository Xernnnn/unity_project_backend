package main

// @title Learning Session Backend API
// @version 1.0
// @description Dokumentasi API Backend
// @host localhost:8080
// @BasePath /

import (
	"fmt"
	"os"

	"github.com/RaFYWStud/LearningSessionBackend/config"
	dbConfig "github.com/RaFYWStud/LearningSessionBackend/config/database"
	"github.com/RaFYWStud/LearningSessionBackend/config/pkg/token"
	"github.com/RaFYWStud/LearningSessionBackend/config/server"
	dbMigration "github.com/RaFYWStud/LearningSessionBackend/database"
)

func main() {
	config.Load()
	token.Load()

	// Handle CLI commands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			runMigrations()
			return
		case "reset":
			runReset()
			return
		case "seed":
			runSeedOnly()
			return
		default:
			fmt.Println("Unknown command. Use: migrate | reset | seed")
			return
		}
	}

	// ðŸ”¥ AUTO-MIGRATE saat server start
	db, _, err := dbConfig.ConnectDB()
	if err != nil {
		panic(fmt.Errorf("failed to connect database: %w", err))
	}

	// Check if migration needed
	if !db.Migrator().HasTable(&dbMigration.User{}) {
		fmt.Println("ðŸ”„ First run detected, running auto-migration...")
		if err := dbMigration.RunMigration(db); err != nil {
			panic(fmt.Errorf("auto-migration failed: %w", err))
		}
	} else {
		fmt.Println("âœ… Database already migrated")
	}

	server.Run()
}

func runMigrations() {
	db, _, err := dbConfig.ConnectDB()
	if err != nil {
		panic(err)
	}

	if err := dbMigration.RunMigration(db); err != nil {
		panic(err)
	}
}

func runReset() {
	db, _, err := dbConfig.ConnectDB()
	if err != nil {
		panic(err)
	}

	fmt.Println("ðŸ—‘ï¸  Dropping all tables...")
	err = db.Migrator().DropTable(
		&dbMigration.User{},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("ðŸ”„ Recreating tables with AutoMigrate...")
	if err := dbMigration.RunMigration(db); err != nil {
		panic(err)
	}

	fmt.Println("âœ… Database reset completed")
}

func runSeedOnly() {
	db, _, err := dbConfig.ConnectDB()
	if err != nil {
		panic(err)
	}

	fmt.Println("ðŸŒ± Running seed only...")
	if err := dbMigration.Seed(db); err != nil {
		panic(err)
	}
	fmt.Println("âœ… Seeding completed")
}
