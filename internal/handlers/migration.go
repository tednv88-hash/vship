package handlers

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"vship/internal/database"
)

func RunMigrations() error {
	// Find all migration SQL files
	files, err := filepath.Glob("./migrations/*.sql")
	if err != nil {
		log.Printf("Error finding migration files: %v", err)
		return err
	}

	if len(files) == 0 {
		log.Println("No migration files found")
		return nil
	}

	// Sort to ensure execution order
	sort.Strings(files)

	sqlDB, err := database.DB.DB()
	if err != nil {
		return err
	}

	for _, f := range files {
		sqlBytes, err := os.ReadFile(f)
		if err != nil {
			log.Printf("Error reading migration file %s: %v", f, err)
			return err
		}

		_, err = sqlDB.Exec(string(sqlBytes))
		if err != nil {
			log.Printf("Migration error in %s: %v", f, err)
			return err
		}

		log.Printf("Migration applied: %s", filepath.Base(f))
	}

	log.Println("All migrations completed successfully")
	return nil
}
