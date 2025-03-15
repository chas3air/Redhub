package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string        `yaml:"env" env-default:"local"`
	AccessTokenTTL   time.Duration `yaml:"access_token_ttl" env-default:"15m"`
	RefreshTokenTTL  time.Duration `yaml:"refresh_token_ttl" env-default:"5h"`
	UsersStorageHost string        `yaml:"usersStorageHost" env-default:"usersManageService"`
	UsersStoragePort int           `yaml:"usersStoragePort" env-default:"50051"`
	Grpc             GrpcConfig    `yaml:"grpc"`
}

type GrpcConfig struct {
	Port    int           `yaml:"port" env-default:"50051"`
	Timeout time.Duration `yaml:"timeout" env-default:"10h"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// Проверка существования файла
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	// --config=./config/local.yaml
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
