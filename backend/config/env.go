package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	Environment     string `env:"ENVIRONMENT" envDefault:"production"`
	HostName        string `env:"HOST_NAME" envDefault:"localhost"`
	HostPort        string `env:"HOST_PORT" envDefault:"8080"`
	DbUser			string `env:"DB_USER" envDefault:""`
	DbPassword 		string `env:"DB_PASSWORD" envDefault:""`
	DbHost			string `env:"DB_HOST" envDefault:""`
	DbPort			string `env:"DB_PORT" envDefault:""`
	DbName			string `env:"DB_NAME" envDefault:""`
	DbURL			string `env:"DB_URL" envDefault:""`
	JWTSecret          string `env:"JWT_SECRET"`
	JWTAccessExpiry    int    `env:"JWT_ACCESS_EXPIRY_HOURS" envDefault:"24"`
	JWTRefreshExpiry   int    `env:"JWT_REFRESH_EXPIRY_DAYS" envDefault:"7"`
	AdminUsername      string `env:"ADMIN_USERNAME" envDefault:"admin"`
	AdminEmail         string `env:"ADMIN_EMAIL" envDefault:"admin@localhost"`
	AdminPassword      string `env:"ADMIN_PASSWORD" envDefault:"admin"`
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func LoadEnv() (*Env, error) {
	_ = godotenv.Load()

	JWTAccessExpiry, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRY_HOURS"))
	JWTRefreshExpiry, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRY_DAYS"))

	return &Env{
		Environment:      getEnv("ENVIRONMENT", "production"),
		HostName:         getEnv("HOST_NAME", "localhost"),
		HostPort:         getEnv("HOST_PORT", "8080"),
		DbUser:           getEnv("DB_USER", ""),
		DbPassword:       getEnv("DB_PASSWORD", ""),
		DbHost:           getEnv("DB_HOST", ""),
		DbPort:           getEnv("DB_PORT", ""),
		DbName:           getEnv("DB_NAME", ""),
		DbURL:            getEnv("DB_URL", ""),
		JWTSecret:        os.Getenv("JWT_SECRET"),
		JWTAccessExpiry:  JWTAccessExpiry,
		JWTRefreshExpiry: JWTRefreshExpiry,
		AdminUsername:    getEnv("ADMIN_USERNAME", "admin"),
		AdminEmail:       getEnv("ADMIN_EMAIL", "admin@localhost"),
		AdminPassword:    getEnv("ADMIN_PASSWORD", "admin"),
	}, nil
}