package database

import (
	"app/models"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "root"
	port     = 5432
	dbname   = "orders_by"
	db       *gorm.DB
	err      error
)

func init() {
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		host, user, password, dbname, port)

	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error saat melakukan koneksi ke database : %v", err.Error())
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)

	fmt.Println("Koneksi ke database")
	db.Debug().AutoMigrate(models.Order{}, models.Item{})
}

func GetConnection() *gorm.DB {
	return db
}

func CloseConnection() {
	sqlDB, _ := db.DB()
	sqlDB.Close()
}
