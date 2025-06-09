package db

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "os"
	"roles/models"
)

var DB *gorm.DB

func ConnectDB() error {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    var err error
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return err
    }

    return DB.AutoMigrate(&models.Role{})
}
