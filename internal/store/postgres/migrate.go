package postgres

import (
	"fmt"
	"log/slog"

	"titanbay/internal/domain"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	slog.Warn("Starting database migration...")

	err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		return fmt.Errorf("failed to create uuid-ossp extension: %w", err)
	}

	// NO need for multiple migration files with AutoMigrate
	err = db.AutoMigrate(
		&domain.Fund{},
		&domain.Investor{},
		&domain.Investment{},
	)
	if err != nil {
		return fmt.Errorf("failed to auto-migrate tables: %w", err)
	}

	slog.Info("Migration completed successfully.")
	return nil
}
