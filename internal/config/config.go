package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	HTTPServer  `yaml:"http_server"`
	RaceLimiter `yaml:"race_limit"`
	Cache       `yaml:"cache"`
}

type HTTPServer struct {
	Address      string        `yaml:"address" env-required:"true"`
	Timeout      time.Duration `yaml: "timeout" env-default:"10s"`
	IddleTimeout time.Duration `yaml:"iddle_timeout" env-default: "2m"`
}

type RaceLimiter struct {
	BucketCap    int    `yaml:"bucket_cap" env-default: "50"`
	FillingSpeed int    `yaml:"filling_speed" env-default: "10"`
	ReqWeight    int    `yaml:"req_weight" env-default: "1"`
	Strategy     string `yaml:"strategy" env-default: "decine"`
}

type Cache struct {
	Addr        string        `yaml:"addr" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true"`
	DB          int           `yaml:"db" env-required:"true"`
	MaxRetries  int           `yaml:"max_retries" env-required:"true"`
	DialTimeout time.Duration `yaml:"dial_timeout" env-required:"true"`
	Timeout     time.Duration `yaml:"timeout" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
