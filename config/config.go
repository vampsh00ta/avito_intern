package config

import (
	"time"
)

type Config struct {
	HTTPServer `yaml:"http_server"`
	DBConfig   `yaml:"db"`
}
type DBConfig struct {
	Username string `yaml:"username" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	Name     string `yaml:"name" env-default:"postgres"`
}
type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func Load() *Config {
	var cfg Config

	cfg = Config{
		HTTPServer{
			Address:     "localhost:8000",
			Timeout:     time.Second * 4,
			IdleTimeout: time.Second * 60,
		},
		DBConfig{
			Username: "avito",
			Password: "avito",
			Host:     "localhost",
			Port:     "5432",
			Name:     "avito",
		},
	}
	return &cfg
}
