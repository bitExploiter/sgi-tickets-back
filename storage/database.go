package storage

import (
	"fmt"
	"log"
	"os"

	"sgi-tickets-back/migrations"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnection() {
	var DSN string = fmt.Sprint(
		"host=", os.Getenv("DB_HOST"),
		" port=", os.Getenv("DB_PORT"),
		" user=", os.Getenv("DB_USER"),
		" password=", os.Getenv("DB_PASSWORD"),
		" dbname=", os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Base de datos conectada")

	if err := migrations.RunMigrations(DB); err != nil {
		log.Fatal(err)
	}
}
