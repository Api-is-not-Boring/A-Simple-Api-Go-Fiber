package handler

import (
	"errors"
	"time"

	m "A-Simple-Api-Go-Fiber/models"
	"github.com/gofiber/fiber/v2"
	j "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jaevor/go-nanoid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var v3Secret = generateSecret()

var V3Config = j.Config{
	SigningMethod: "HS512",
	SigningKey:    []byte(v3Secret),
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		ctx.Status(fiber.StatusUnauthorized)
		m := new(m.V3Message)
		m.SetMessage(ctx, "Invalid or expired JWT")
		return ctx.JSON(m)
	},
}

var V3Middleware = j.New(V3Config)

func generateSecret() string {
	secret, err := nanoid.Standard(32)
	if err != nil {
		panic(err)
	}
	return secret()
}

func generateToken() (string, error) {
	// Create the Claims
	claims := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Second * 60).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	t, err := token.SignedString([]byte(v3Secret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func Announce(c *fiber.Ctx) interface{} {
	message := new(m.V3Message)
	message.SetMessage(c, "Login with Post Request")
	return message
}

func Login(c *fiber.Ctx, db *gorm.DB) interface{} {
	loginUser := new(m.User)
	if err := db.Take(&loginUser, "Username = ?", c.FormValue("username")).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		message := new(m.V3Message)
		c.Status(fiber.StatusUnauthorized)
		message.SetMessage(c, "Unauthorized")
		return message
	}
	if err := bcrypt.CompareHashAndPassword([]byte(loginUser.Password), []byte(c.FormValue("password"))); err != nil {
		message := new(m.V3Message)
		c.Status(fiber.StatusUnauthorized)
		message.SetMessage(c, "Unauthorized")
		return message
	}
	token, _ := generateToken()
	message := m.V3Success{Token: token}
	message.V3Message.SetMessage(c, "Login Successful !!!")
	return message
}

func Check(c *fiber.Ctx) interface{} {
	message := new(m.V3Message)
	message.SetMessage(c, "JWT Token validation successful!")
	return message
}

func Secure() map[string]interface{} {
	return fiber.Map{"s3cr5t": v3Secret}
}
