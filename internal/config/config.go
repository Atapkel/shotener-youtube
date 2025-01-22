package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string     `yaml:"env" env-required:"true"` //env-default:"local"
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idleTimeout" env-default:"10s"`
}

func MustLoad() *Config {
	var configPath string

	flag.StringVar(&configPath, "conf", "./config/local.yaml", "path for yaml that stores config")
	flag.Parse()
	if configPath == "" {
		log.Fatal("Config path is empty")
	}
	if _, err := os.Stat(configPath); os.IsExist(err) {
		log.Fatalf("config file %s doesn't exist", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config %s", err)
	}
	return &cfg
}
