package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Car struct {
	ID    int    `gorm:"primaryKey" gorm:"autoIncrement" json:"id"`
	Name  string `gorm:"not null" json:"name"`
	Price int    `gorm:"not null" json:"price"`
}

type AllCars struct {
	Total int   `json:"total"`
	Cars  []Car `json:"cars"`
}

var cars = []Car{
	{Name: "Audi", Price: 52642},
	{Name: "Mercedes", Price: 57127},
	{Name: "Skoda", Price: 9000},
	{Name: "Volvo", Price: 29000},
	{Name: "Bentley", Price: 350000},
	{Name: "Citroen", Price: 21000},
	{Name: "Hummer", Price: 41400},
	{Name: "Volkswagen", Price: 21600},
}

func DbInit() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&Car{})
	if err != nil {
		panic("failed to migrate database")
	}

	// Create Data
	db.CreateInBatches(&cars, 8)

	return db, err
}

func DbReset(db *gorm.DB) {
	err := db.Migrator().DropTable(&Car{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&Car{})
	if err != nil {
		return
	}
	db.CreateInBatches(&cars, 8)
}
