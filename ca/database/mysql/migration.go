package mysql

import (
	"fmt"

	"github.com/cloudslit/cloudslit/ca/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	lo := logger.Named("migration")
	sql, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get DB instance: %v", err)
	}
	driver, err := mysql.WithInstance(sql, &mysql.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/mysql/migrations/",
		"mysql", driver)
	if err != nil {
		return fmt.Errorf("migrate instance error: %v", err)
	}
	if err = m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			lo.Info("no changes.")
			return nil
		}
		return fmt.Errorf("MySQL migration exception: %v", err)
	}
	lo.Info("Migrations success.")
	return nil
}
