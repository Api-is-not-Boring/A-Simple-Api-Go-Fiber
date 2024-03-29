package handler

import (
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/shirou/gopsutil/v3/net"
)

type Pong struct {
	Agent   string `json:"agent"`
	Date    string `json:"datetime"`
	Message string `json:"message"`
}

func GetPong(c *fiber.Ctx) Pong {
	req := c.GetReqHeaders()
	return Pong{
		Agent:   req[fiber.HeaderUserAgent][0],
		Date:    time.Now().Format(http.TimeFormat),
		Message: "pong",
	}
}

type project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    string `json:"language"`
	Url         string `json:"url"`
	GitHash     string `json:"git hash"`
	Version     string `json:"version"`
}

type Info struct {
	Project project `json:"project"`
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
		Project: project{
			Name:        os.Getenv("NAME"),
			Description: os.Getenv("DESCRIPTION"),
			Language:    os.Getenv("LANGUAGE"),
			Url:         os.Getenv("URL"),
			GitHash:     getGitHash(),
			Version:     os.Getenv("VERSION"),
		},
	}
}

type connection struct {
	Id       int    `json:"id"`
	Protocol string `json:"protocol"`
	Type     string `json:"type"`
	Local    string `json:"local"`
	Remote   string `json:"remote"`
}

type Client struct {
	Connections []connection `json:"connections"`
}

func GetConnections() Client {
	r := Client{}
	cons, _ := net.ConnectionsPid("tcp", int32(os.Getpid()))
	for _, con := range cons {
		if con.Status == "LISTEN" {
			con.Status = "LISTENING"
		}
		if con.Laddr.IP == "*" {
			con.Laddr.IP = "0.0.0.0"
		}
		if con.Raddr.IP == "" {
			con.Raddr.IP = "0.0.0.0"
		}
		r.Connections = append(r.Connections, connection{
			Id:       int(con.Fd),
			Protocol: "TCP",
			Type:     con.Status,
			Local:    con.Laddr.IP + ":" + strconv.Itoa(int(con.Laddr.Port)),
			Remote:   con.Raddr.IP + ":" + strconv.Itoa(int(con.Raddr.Port)),
		})
	}
	return r
}
