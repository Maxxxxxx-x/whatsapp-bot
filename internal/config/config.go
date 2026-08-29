package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv         string
	LogLevel       string
	CommandPrefix  string
	WhatsappDBPath string
	AppDBPath      string
	LogDir         string
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

func Load() *Config {
	godotenv.Load()
	return &Config{
		AppEnv:         getEnv("APP_ENV", "dev"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		CommandPrefix:  getEnv("COMMAND_PREFIX", "!"),
		WhatsappDBPath: getEnv("WHATAPP_DB_PATH", "data/whatsmeow.db"),
		AppDBPath:      getEnv("APP_DB_PATH", "data/app.db"),
		LogDir:         getEnv("LOG_DIR", "logs"),
	}
}

func (c *Config) LogFilePath() string {
	return filepath.Join(c.LogDir, "whatsapp-bot.log")
}
