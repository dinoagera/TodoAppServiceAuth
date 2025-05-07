package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	StoragePath string        `env:"STORAGE_PATH" env-required:"true"`
	TokenTTL    time.Duration `env:"TOKEN_TTL"`
	GRPC        GRPC
	JWTSecret   string `env:"JWT_SECRET" env-required:"true"`
}
type GRPC struct {
	Port    string        `env:"GRPC_PORT_API_AUTH" env-default:"8082"`
	Timeout time.Duration `env:"GRPC_TIME_OUT_API_AUTH"`
}

var cfg *Config

func InitConfig(log *slog.Logger) *Config {
	cfgPath := ".env"
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		log.Debug("config path is empty.", "Path:", cfgPath, "Error:", err.Error())
		log.Info("config path is empty.")
		os.Exit(1)
	}
	var localCfg Config
	err := cleanenv.ReadConfig(cfgPath, &localCfg)
	if err != nil {
		log.Debug("reading config to failed.", "Error:", err.Error())
		log.Info("reading config to failed.")
		os.Exit(1)
	}
	cfg = &localCfg
	return cfg
}
func GetConfig() *Config {
	if cfg == nil {
		panic("config is not init")
	}
	return cfg
}
