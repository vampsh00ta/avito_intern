package config

import (
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	HTTPServer `yaml:"http_server"`
	DBConfig   `yaml:"db"`
	TTL        `yaml:"ttl"`
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
type TTL struct {
	TimeUpdate time.Duration `yaml:"timeupdate" env-default:"5s"`
}

func Load() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	currPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	filePath := currPath + os.Getenv("path") + "/" + os.Getenv("env") + ".yml"
	file, err := os.Open(filePath)
	defer file.Close()

	d := yaml.NewDecoder(file)
	var cfg *Config

	if err := d.Decode(&cfg); err != nil {
		return nil, err
	}
	return cfg, nil

}
