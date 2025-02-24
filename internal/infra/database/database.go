package database

import (
	"b8boost/backend/config"
	"b8boost/backend/internal/adapters/repo"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(conf config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Moscow", conf.DatabaseHost, conf.DatabaseUser, conf.DatabasePassword, conf.DatabaseDB, conf.DatabasePort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	repo.Migrate(db)
	return db
}
