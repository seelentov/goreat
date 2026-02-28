package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSQLiteInMemoryDB() (*gorm.DB, error) {
	db, err := newDB(sqlite.Open("file::memory:?cache=shared"))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewSQLiteFileDB(filepath string) (*gorm.DB, error) {
	db, err := newDB(sqlite.Open(filepath))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func newDB(dialector gorm.Dialector) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	err = db.AutoMigrate(dbModels...)
	if err != nil {
		return nil, err
	}

	return db, nil
}
