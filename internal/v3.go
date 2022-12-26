package internal

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

var v3MessagePrefix = "[v3] -> "

type V3Message struct {
	Message string `json:"message"`
}

func (v *V3Message) setMessage(c *fiber.Ctx, message string) {
	v.Message = v3MessagePrefix + strconv.Itoa(c.Response().StatusCode()) + " " + message
}

func Announce(c *fiber.Ctx) V3Message {
	message := V3Message{}
	message.setMessage(c, "Login with Post Request")
	return message
}
