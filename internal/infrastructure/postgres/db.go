package postgres

import (
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Open creates a GORM connection to PostgreSQL.
func Open(databaseURL string) (*gorm.DB, error) {
	return gorm.Open(gormpostgres.Open(databaseURL), &gorm.Config{})
}
