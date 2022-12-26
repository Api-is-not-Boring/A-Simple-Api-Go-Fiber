package internal

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	j "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jaevor/go-nanoid"
	"gorm.io/gorm"
	"strconv"
	"time"
)

var v3Secret = generateSecret()

var V3Config = j.Config{
	SigningMethod: "HS512",
	SigningKey:    []byte(v3Secret),
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		ctx.Status(fiber.StatusUnauthorized)
		m := new(v3Message)
		m.setMessage(ctx, "Invalid or expired JWT")
		return ctx.JSON(m)
	},
}

var V3Middleware = j.New(V3Config)

var v3MessagePrefix = "[v3] -> "

type v3Message struct {
	Message string `json:"message"`
}

func (v *v3Message) setMessage(c *fiber.Ctx, message string) {
	v.Message = v3MessagePrefix + strconv.Itoa(c.Response().StatusCode()) + " " + message
}

type v3Success struct {
	v3Message
	Token string `json:"token"`
}

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
	m := new(v3Message)
	m.setMessage(c, "Login with Post Request")
	return m
}

func Login(c *fiber.Ctx, db *gorm.DB) interface{} {
	loginUser := new(user)
	if err := c.BodyParser(loginUser); err != nil {
		return err
	}
	if err := db.Take(&user{}, "Username = ? AND Password = ?", &loginUser.Username, &loginUser.Password).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		m := new(v3Message)
		c.Status(fiber.StatusUnauthorized)
		m.setMessage(c, "Unauthorized")
		return m
	}
	token, _ := generateToken()
	message := v3Success{Token: token}
	message.v3Message.setMessage(c, "Login Successful !!!")
	return message
}

func Check(c *fiber.Ctx) interface{} {
	m := new(v3Message)
	m.setMessage(c, "JWT Token validation successful!")
	return m
}

func Secure() map[string]interface{} {
	return fiber.Map{"s3cr5t": v3Secret}
}
