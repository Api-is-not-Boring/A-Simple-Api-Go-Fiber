package models

type Car struct {
	ID    int    `gorm:"primaryKey,autoIncrement" json:"id"`
	Name  string `gorm:"not null" json:"name" validate:"required"`
	Price int    `gorm:"not null" json:"price" validate:"required,number"`
}

type AllCars struct {
	Total int   `json:"total"`
	Cars  []Car `json:"cars"`
}

var InitCars = []Car{
	{Name: "Audi", Price: 52642},
	{Name: "Mercedes", Price: 57127},
	{Name: "Skoda", Price: 9000},
	{Name: "Volvo", Price: 29000},
	{Name: "Bentley", Price: 350000},
	{Name: "Citroen", Price: 21000},
	{Name: "Hummer", Price: 41400},
	{Name: "Volkswagen", Price: 21600},
}
