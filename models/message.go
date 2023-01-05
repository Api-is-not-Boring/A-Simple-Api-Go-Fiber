package models

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RequestMethod string

type V2Response struct {
	Method string `json:"method"`
	Car    Car    `json:"car"`
}

type v2Error struct {
	Method RequestMethod `json:"method"`
	Error  string        `json:"error"`
}

var V2RecordNotFound = v2Error{Method: "", Error: "404 Record not found"}
var V2BadRequest = v2Error{Method: "", Error: "405 Bad Request"}
var V2InternalError = v2Error{Method: "", Error: "500 Internal Server Error"}

func (m *RequestMethod) SetMethod(c *fiber.Ctx, p ...string) {
	*m = RequestMethod(SetMethod(c, p...))
}

func SetMethod(c *fiber.Ctx, p ...string) string {
	var v2MessagePrefix = "[v2] -> "
	if p == nil {
		p = []string{""}
	}
	switch p[0] {
	case "JSON":
		return v2MessagePrefix + c.Method() + " with " + "JSON"
	case "Query":
		return v2MessagePrefix + c.Method() + " with " + "Query Parameter"
	case "Path":
		return v2MessagePrefix + c.Method() + " with " + "Path Parameter"
	default:
		return v2MessagePrefix + c.Method()
	}
}

type V3Message struct {
	Message string `json:"message"`
}

func (v *V3Message) SetMessage(c *fiber.Ctx, message string) {
	var v3MessagePrefix = "[v3] -> "
	v.Message = v3MessagePrefix + strconv.Itoa(c.Response().StatusCode()) + " " + message
}

type V3Success struct {
	V3Message
	Token string `json:"token"`
}
