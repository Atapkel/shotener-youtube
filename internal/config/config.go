package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string     `yaml:"env" env-required:"true"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	HttpServer  HttpServer `yaml:"http_server"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"10s"`
}

func MustLoad() *Config {
	var configPath string

	flag.StringVar(&configPath, "path", "./config/local.yaml",
		"enter path of yaml that stores config")
	flag.Parse()

	if configPath == "" {
		log.Fatal("Config path is empty")
	}
	if _, err := os.Stat(configPath); os.IsExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config %s", err)
	}
	return &cfg
}
