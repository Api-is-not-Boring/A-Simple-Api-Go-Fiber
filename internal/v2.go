package internal

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"strconv"
)

var validate = validator.New()

type car struct {
	ID    int    `gorm:"primaryKey,autoIncrement" json:"id"`
	Name  string `gorm:"not null" json:"name" validate:"required"`
	Price int    `gorm:"not null" json:"price" validate:"required,number"`
}

type allCars struct {
	Total int   `json:"total"`
	Cars  []car `json:"cars"`
}

type user struct {
	ID       int    `gorm:"primaryKey,autoIncrement"`
	Username string `gorm:"not null" form:"username"`
	Password string `gorm:"not null" form:"password"`
}

type requestMethod string

type v2Response struct {
	Method requestMethod `json:"method"`
	Car    car           `json:"car"`
}

type v2Error struct {
	Method requestMethod `json:"method"`
	Error  string        `json:"error"`
}

var v2MessagePrefix = "[v2] -> "

func setMethod(c *fiber.Ctx, p ...string) requestMethod {
	if p == nil {
		p = []string{""}
	}
	switch p[0] {
	case "JSON":
		return requestMethod(v2MessagePrefix + c.Method() + " with " + "JSON")
	case "Query":
		return requestMethod(v2MessagePrefix + c.Method() + " with " + "Query Parameter")
	case "Path":
		return requestMethod(v2MessagePrefix + c.Method() + " with " + "Path Parameter")
	default:
		return requestMethod(v2MessagePrefix + c.Method())
	}
}

func (m *requestMethod) setMethod(c *fiber.Ctx, p ...string) {
	*m = setMethod(c, p...)
}

var v2RecordNotFound = v2Error{Method: "", Error: "404 Record not found"}
var v2BadRequest = v2Error{Method: "", Error: "405 Bad Request"}
var v2InternalError = v2Error{Method: "", Error: "500 Internal Server Error"}

var initCars = []car{
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
	err = db.AutoMigrate(&car{})
	if err != nil {
		panic("failed to migrate database")
	}

	err = db.AutoMigrate(&user{})

	// Create Data
	db.CreateInBatches(&initCars, 8)
	password, _ := bcrypt.GenerateFromPassword([]byte("password"), 12)
	db.Create(&user{Username: "admin", Password: string(password)})
	return db, err
}

func DbReset(db *gorm.DB) {
	err := db.Migrator().DropTable(&car{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&car{})
	if err != nil {
		return
	}
	db.CreateInBatches(&initCars, 8)
}

func GetCars(c *fiber.Ctx, db *gorm.DB) interface{} {
	if id, err := strconv.Atoi(c.Query("id")); id != 0 && err == nil {
		var inputCar car
		db.First(&inputCar, id)
		return v2Response{
			Method: setMethod(c, "Query"),
			Car:    inputCar,
		}
	} else if id, err := c.ParamsInt("id"); id != 0 && err == nil {
		var inputCar car
		db.First(&inputCar, id)
		return v2Response{
			Method: setMethod(c, "Path"),
			Car:    inputCar,
		}
	} else {
		var cars []car
		db.Find(&cars)
		return allCars{Total: len(cars), Cars: cars}
	}
}

func CreateCar(c *fiber.Ctx, db *gorm.DB) interface{} {
	var inputCar car
	var records int64
	if err := c.BodyParser(&inputCar); err != nil {
		c.Status(fiber.StatusBadRequest)
		v2BadRequest.Method.setMethod(c)
		return v2BadRequest
	}
	if err := validate.Struct(inputCar); err != nil {
		c.Status(fiber.StatusBadRequest)
		v2InternalError.Method.setMethod(c)
		return v2InternalError
	}
	db.Model(&[]car{}).Count(&records)
	if records >= 20 {
		log.Println("[v2] -> Maximum number of records reached. Resetting database.")
		DbReset(db)
	}
	db.Create(&inputCar)
	c.Status(fiber.StatusCreated)
	return v2Response{
		Method: setMethod(c, "JSON"),
		Car:    inputCar,
	}
}

func UpdateCar(c *fiber.Ctx, db *gorm.DB) interface{} {
	var inputCar car
	if err := c.BodyParser(&inputCar); err != nil {
		c.Status(fiber.StatusBadRequest)
		v2BadRequest.Method.setMethod(c)
		return v2BadRequest
	}
	if err := validate.Struct(inputCar); err != nil {
		c.Status(fiber.StatusBadRequest)
		v2InternalError.Method.setMethod(c)
		return v2InternalError
	}
	if id, err := strconv.Atoi(c.Query("id")); id != 0 && err == nil {
		if err := db.Take(&inputCar, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			v2RecordNotFound.Method.setMethod(c, "Query")
			c.Status(fiber.StatusNotFound)
			return v2RecordNotFound
		}
		db.Model(&inputCar).Clauses(clause.Returning{}).Where("id = ?", id).Updates(inputCar)
		return &v2Response{
			Method: setMethod(c, "Query"),
			Car:    inputCar,
		}
	} else if id, err := c.ParamsInt("id"); id != 0 && err == nil {
		if err := db.Take(&inputCar, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			v2RecordNotFound.Method.setMethod(c, "Path")
			c.Status(fiber.StatusNotFound)
			return v2RecordNotFound
		}
		db.Model(&inputCar).Clauses(clause.Returning{}).Where("id = ?", id).Updates(inputCar)
		return &v2Response{
			Method: setMethod(c, "Path"),
			Car:    inputCar,
		}
	} else {
		db.Model(&inputCar).Updates(inputCar)
		return &v2Response{
			Method: setMethod(c, "JSON"),
			Car:    inputCar,
		}
	}
}

func DeleteCar(c *fiber.Ctx, db *gorm.DB) interface{} {
	var inputCar car
	if id, err := strconv.Atoi(c.Query("id")); id != 0 && err == nil {
		if err := db.Take(&inputCar, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			v2RecordNotFound.Method.setMethod(c, "Query")
			c.Status(fiber.StatusNotFound)
			return v2RecordNotFound
		}
		db.Delete(inputCar)
		return &v2Response{
			Method: setMethod(c, "Query"),
			Car:    inputCar,
		}
	} else if id, err := c.ParamsInt("id"); id != 0 && err == nil {
		if err := db.Take(&inputCar, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			v2RecordNotFound.Method.setMethod(c, "Path")
			c.Status(fiber.StatusNotFound)
			return v2RecordNotFound
		}
		db.Delete(inputCar)
		return &v2Response{
			Method: setMethod(c, "Path"),
			Car:    inputCar,
		}
	} else {
		c.Status(fiber.StatusBadRequest)
		v2BadRequest.Method.setMethod(c)
		return v2BadRequest
	}
}
