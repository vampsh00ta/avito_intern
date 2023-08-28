package config

import (
	"time"
)

type Config struct {
	HTTPServer `yaml:"http_server"`
	DBConfig   `yaml:"db"`
	TTL        `yaml:"ttl"`
	Redis      `yaml:"redis"`
}
type DBConfig struct {
	Username string `yaml:"username" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	Name     string `yaml:"name" env-default:"postgres"`
}
type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
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
	//err := godotenv.Load(".env")
	//if err != nil {
	//	return nil, err
	//}
	//dir, err := os.Getwd()
	//if err != nil {
	//	return nil, err
	//}
	//
	//configFile := fmt.Sprintf(dir+"/config/%s.yaml", os.Getenv("env"))
	//f, err := os.Open(configFile)
	//if err != nil {
	//	return nil, err
	//}
	//defer f.Close()
	//var cfg Config
	//decoder := yaml.NewDecoder(f)
	//err = decoder.Decode(&cfg)
	//if err != nil {
	//	return nil, err
	//}
	var cfg *Config
	cfg = &Config{
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
		TTL{5},
		Redis{DB: 0, Password: "", Addr: "localhost:6379"},
	}
	return cfg, nil

}
