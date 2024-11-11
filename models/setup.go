package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Konfigurasi koneksi MySQL
	dsn := "root:gibran123@tcp(localhost:3306)/belajar-go"

	// Membuka koneksi ke MySQL
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	database.AutoMigrate(&Product{}, &User{}, &Transaction{}, &Wallet{})
	DB = database

	SeedAdminUser()
}
