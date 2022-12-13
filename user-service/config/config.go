package config

import (
	"helper"
	"os"
)

type Config struct {
	Port int
	Env  string
	Api  string
	Host string
	Dsn  string
}

func (c *Config) Setup() {
	c.Port = helper.GetEnvInt("PORT")
	c.Dsn = os.Getenv("DSN")
	c.Env = os.Getenv("ENV")
	c.Api = os.Getenv("API")
	c.Host = os.Getenv("HOST")
}
