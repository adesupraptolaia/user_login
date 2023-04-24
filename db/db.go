package db

import (
	"fmt"
	"time"

	"github.com/adesupraptolaia/user_login/config"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDatabase(migrationPath string) (*gorm.DB, error) {
	cfg := config.Config

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %s", err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error getting database connection: %s", err.Error())
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Perform any database migrations or setup here
	err = goose.Up(sqlDB, migrationPath, goose.WithNoVersioning())
	if err != nil {
		return nil, fmt.Errorf("error when perform database migrations, err: %s", err.Error())
	}

	return db, nil
}
