package models

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

type Pong struct {
	Agent   string `json:"agent"`
	Date    string `json:"date"`
	Message string `json:"message"`
}

func GetPong(c *fiber.Ctx) Pong {
	req := c.GetReqHeaders()
	return Pong{
		Agent:   req[fiber.HeaderUserAgent],
		Date:    time.Now().Format(http.TimeFormat),
		Message: "pong",
	}
}

type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Url         string `json:"url"`
	GitHash     string `json:"git hash"`
	Version     string `json:"version"`
}

type Info struct {
	Project Project `json:"project"`
}

func getGitHash() string {
	var gitHash string
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				gitHash = setting.Value[:7]
			}
		}
	}
	return gitHash
}

func GetInfo() Info {
	return Info{
		Project: Project{
			Name:        os.Getenv("NAME"),
			Description: os.Getenv("DESCRIPTION"),
			Language:    os.Getenv("LANGUAGE"),
			Url:         os.Getenv("URL"),
			GitHash:     getGitHash(),
			Version:     os.Getenv("VERSION"),
		},
	}
}

type Client struct {
	Connections []struct {
		Id       string `json:"id"`
		Protocol string `json:"protocol"`
		Type     string `json:"type"`
		Local    string `json:"local"`
		Remote   string `json:"remote"`
	} `json:"connections"`
}
