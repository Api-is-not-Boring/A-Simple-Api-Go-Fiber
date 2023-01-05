package handler

import (
	"errors"
	"log"
	"strconv"

	m "A-Simple-Api-Go-Fiber/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var validate = validator.New()

func DbInit() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&m.Car{})
	if err != nil {
		panic("failed to migrate database")
	}

	err = db.AutoMigrate(&m.User{})

	// Create Data
	db.CreateInBatches(&m.InitCars, 8)
	password, _ := bcrypt.GenerateFromPassword([]byte("password"), 12)
	db.Create(&m.User{Username: "admin", Password: string(password)})
	return db, err
}

func DbReset(db *gorm.DB) {
	err := db.Migrator().DropTable(&m.Car{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&m.Car{})
	if err != nil {
		return
	}
	db.CreateInBatches(&m.InitCars, 8)
}

func GetCars(c *fiber.Ctx, db *gorm.DB) interface{} {
	if id, err := strconv.Atoi(c.Query("id")); id != 0 && err == nil {
		var inputCar m.Car
		db.First(&inputCar, id)
		return m.V2Response{
			Method: m.SetMethod(c, "Query"),
			Car:    inputCar,
		}
	} else if id, err := c.ParamsInt("id"); id != 0 && err == nil {
		var inputCar m.Car
		db.First(&inputCar, id)
		return m.V2Response{
			Method: m.SetMethod(c, "Path"),
			Car:    inputCar,
		}
	} else {
		var Cars []m.Car
		db.Find(&Cars)
		return m.AllCars{Total: len(Cars), Cars: Cars}
	}
}

func CreateCar(c *fiber.Ctx, db *gorm.DB) interface{} {
	var inputCar m.Car
	var records int64
	if err := c.BodyParser(&inputCar); err != nil {
		c.Status(fiber.StatusBadRequest)
		m.V2BadRequest.Method.SetMethod(c)
		return m.V2BadRequest
	}
	if err := validate.Struct(inputCar); err != nil {
		c.Status(fiber.StatusBadRequest)
		m.V2InternalError.Method.SetMethod(c)
		return m.V2InternalError
	}
	db.Model(&[]m.Car{}).Count(&records)
	if records >= 20 {
		log.Println("[v2] -> Maximum number of records reached. Resetting database.")
		DbReset(db)
	}
	db.Create(&inputCar)
	c.Status(fiber.StatusCreated)
	return m.V2Response{
		Method: m.SetMethod(c, "JSON"),
		Car:    inputCar,
	}
}

func UpdateCar(c *fiber.Ctx, db *gorm.DB) interface{} {
	var inputCar m.Car
	if err := c.BodyParser(&inputCar); err != nil {
		c.Status(fiber.StatusBadRequest)
		m.V2BadRequest.Method.SetMethod(c)
		return m.V2BadRequest
	}
	if err := validate.Struct(inputCar); err != nil {
		c.Status(fiber.StatusBadRequest)
		m.V2InternalError.Method.SetMethod(c)
		return m.V2InternalError
	}
	if id, err := strconv.Atoi(c.Query("id")); id != 0 && err == nil {
		if err := db.Take(&inputCar, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			m.V2RecordNotFound.Method.SetMethod(c, "Query")
			c.Status(fiber.StatusNotFound)
			return m.V2RecordNotFound
		}
		db.Model(&inputCar).Clauses(clause.Returning{}).Where("id = ?", id).Updates(inputCar)
		return &m.V2Response{
			Method: m.SetMethod(c, "Query"),
			Car:    inputCar,
		}
	} else if id, err := c.ParamsInt("id"); id != 0 && err == nil {
		if err := db.Take(&inputCar, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			m.V2RecordNotFound.Method.SetMethod(c, "Path")
			c.Status(fiber.StatusNotFound)
			return m.V2RecordNotFound
		}
		db.Model(&inputCar).Clauses(clause.Returning{}).Where("id = ?", id).Updates(inputCar)
		return &m.V2Response{
			Method: m.SetMethod(c, "Path"),
			Car:    inputCar,
		}
	} else {
		db.Model(&inputCar).Updates(inputCar)
		return &m.V2Response{
			Method: m.SetMethod(c, "JSON"),
			Car:    inputCar,
		}
	}
}

func DeleteCar(c *fiber.Ctx, db *gorm.DB) interface{} {
	var inputCar m.Car
	if id, err := strconv.Atoi(c.Query("id")); id != 0 && err == nil {
		if err := db.Take(&inputCar, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			m.V2RecordNotFound.Method.SetMethod(c, "Query")
			c.Status(fiber.StatusNotFound)
			return m.V2RecordNotFound
		}
		db.Delete(inputCar)
		return &m.V2Response{
			Method: m.SetMethod(c, "Query"),
			Car:    inputCar,
		}
	} else if id, err := c.ParamsInt("id"); id != 0 && err == nil {
		if err := db.Take(&inputCar, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			m.V2RecordNotFound.Method.SetMethod(c, "Path")
			c.Status(fiber.StatusNotFound)
			return m.V2RecordNotFound
		}
		db.Delete(inputCar)
		return &m.V2Response{
			Method: m.SetMethod(c, "Path"),
			Car:    inputCar,
		}
	} else {
		c.Status(fiber.StatusBadRequest)
		m.V2BadRequest.Method.SetMethod(c)
		return m.V2BadRequest
	}
}
