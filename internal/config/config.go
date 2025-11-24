package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database string           `yaml:"database" env-required:"true"`
	Server   HTTPServerConfig `yaml:"server" env-required:"false"`
}

type HTTPServerConfig struct {
	ReadTimeout time.Duration `yaml:"read_timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	Host        string        `yaml:"host" env-default:"127.0.0.1"`
	Port        int           `yaml:"port" env-default:"8080"`
}

func (s *HTTPServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func MustLoad(path string) *Config {
	if path == "" {
		panic("configuration file path is not set")
	}

	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}
