package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServ    `yaml:"http_server"`
	Postgresql  `yaml:"postgresql"`
}

type HTTPServ struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type Postgresql struct {
	DB_Address string `yaml:"db_address" env-default:"localhost:5432"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	DbName     string `yaml:"db_name"`
}

func MustLoad() *Config {
	configFilePath := `C:\Users\artmp\go-projects\url-performer\config\local.yaml`

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		log.Fatalln("config file does not exist ", configFilePath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configFilePath, &config); err != nil {
		log.Fatalln("cannot read config", err)
	}

	return &config
}
