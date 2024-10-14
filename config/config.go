package config

import (
	"os"
	"strconv"

	"github.com/amedoeyes/hadath/internal/utility"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost string
	ServerPort int

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

var config *Config

func Load(files ...string) error {
	root, err := utility.FindProjectRoot()
	if err != nil {
		return err
	}

	var filenames []string
	for _, file := range files {
		filenames = append(filenames, root+"/"+file)
	}

	err = godotenv.Load(filenames...)
	if err != nil {
		return err
	}

	config = &Config{
		ServerHost: getEnv("SERVER_HOST", "localhost"),
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnvAsInt("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", ""),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", ""),
	}

	return nil
}

func Get() *Config {
	return config
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
