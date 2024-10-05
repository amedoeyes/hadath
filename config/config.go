package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	serverHost string
	serverPort int

	dbHost     string
	dbPort     int
	dbUser     string
	dbPassword string
	dbName     string
}

func (c *Config) ServerHost() string {
	return c.serverHost
}

func (c *Config) ServerPort() int {
	return c.serverPort
}

func (c *Config) DBHost() string {
	return c.dbHost
}

func (c *Config) DBPort() int {
	return c.dbPort
}

func (c *Config) DBUser() string {
	return c.dbUser
}

func (c *Config) DBPassword() string {
	return c.dbPassword
}

func (c *Config) DBName() string {
	return c.dbName
}

var config *Config

func Load() {
	godotenv.Load()
	config = &Config{
		serverHost: getEnv("SERVER_HOST", "localhost"),
		serverPort: getEnvAsInt("SERVER_PORT", 8080),
		dbHost:     getEnv("DB_HOST", "localhost"),
		dbPort:     getEnvAsInt("DB_PORT", 5432),
		dbUser:     getEnv("DB_USER", ""),
		dbPassword: getEnv("DB_PASSWORD", ""),
		dbName:     getEnv("DB_NAME", ""),
	}
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
