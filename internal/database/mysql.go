package database

import (
	"context"
	"fmt"
	"time"

	"bookshare/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMySQL(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnLifetime)

	return db, nil
}

func NewTestMySQL() (*gorm.DB, error) {
	dsn := "root:root@tcp(localhost:3306)/bookshare_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

type DB struct {
	*gorm.DB
}

func (db *DB) Transaction(fn func(tx *gorm.DB) error) error {
	return db.DB.Transaction(fn)
}

func (db *DB) WithContext(timeout time.Duration) *gorm.DB {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return db.DB.WithContext(ctx)
}
