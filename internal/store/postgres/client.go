package postgres

import (
	"fmt"
	"log/slog"

	gormpg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresStore(url string) (*gorm.DB, error) {
	db, err := gorm.Open(gormpg.Open(url))

	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database is not reachable: %w", err)
	}
	slog.Info("database connection established")
	return db, nil
}
