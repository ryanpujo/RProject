package infrastructure

import (
	"helper"
	"os"
)

type config struct {
	Port     int
	GrpcPort int
	Env      string
	Api      string
	Host     string
	Dsn      string
}

func (c *config) Setup() {
	c.Port = helper.GetEnvInt("PORT")
	c.Dsn = os.Getenv("DSN")
	c.Env = os.Getenv("ENV")
	c.Api = os.Getenv("API")
	c.Host = os.Getenv("HOST")
}
