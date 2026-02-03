package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	Environment     string `env:"ENVIRONMENT" envDefault:"production"`
	HostName        string `env:"HOST_NAME" envDefault:"localhost"`
	HostPort        string `env:"HOST_PORT" envDefault:"8080"`
	DbUser			string `env:"DB_PASSWORD" envDefault:"/"`
	DbPassword 		string `env:"DB_PASSWORD" envDefault:"/"`
	DbHost			string `env:"DB_HOST" envDefault:"/"`
	DbPort			string `env:"DB_PORT" envDefault:"/"`
	DbName			string `env:"DB_NAME" envDefault:"/"`
	DbURL			string `env:"DB_URL" envDefault:"/"`
}

func LoadEnv() (*Env, error) {
	_ = godotenv.Load()

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "production"
	}

	return &Env{
		Environment:     environment,
		HostName:        os.Getenv("HOST_NAME"),
		HostPort:        os.Getenv("HOST_PORT"),
		DbUser: 		 os.Getenv("DB_USER"),
		DbPassword: 	 os.Getenv("DB_PASSWORD"),
		DbHost:			 os.Getenv("DB_HOST"),
		DbPort:			 os.Getenv("DB_PORT"),
		DbName:			 os.Getenv("DB_NAME"),
		DbURL:			 os.Getenv("DB_URL"),
	}, nil
}