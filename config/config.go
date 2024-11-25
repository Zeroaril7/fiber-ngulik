package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type envConfig struct {
	AppName            string
	AppPort            string
	PostgreSQLHost     string
	PostgreSQLPort     string
	PostgreSQLUsername string
	PostgreSQLPassword string
	PostgreSQLDBName   string
	BasicAuthUsername  string
	BasicAuthPassword  string
	PrivateKey         string
	PublicKey          string
}

var envCfg envConfig

func init() {
	LoadConfig()
}

func LoadConfig() {
	err := godotenv.Load()

	if err != nil {
		println(err.Error())
	}

	envCfg = envConfig{
		AppName:            os.Getenv("APP_NAME"),
		AppPort:            os.Getenv("APP_PORT"),
		PostgreSQLHost:     os.Getenv("POSTGRESQL_HOST"),
		PostgreSQLPort:     os.Getenv("POSTGRESQL_PORT"),
		PostgreSQLUsername: os.Getenv("POSTGRESQL_USERNAME"),
		PostgreSQLPassword: os.Getenv("POSTGRESQL_PASSWORD"),
		PostgreSQLDBName:   os.Getenv("POSTGRESQL_DB_NAME"),
		BasicAuthUsername:  os.Getenv("BASIC_AUTH_USERNAME"),
		BasicAuthPassword:  os.Getenv("BASIC_AUTH_PASSWORD"),
		PrivateKey:         os.Getenv("PRIVATE_KEY"),
		PublicKey:          os.Getenv("PUBLIC_KEY"),
	}
}

func (e envConfig) PostgreDSN() (string, string) {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			envCfg.PostgreSQLHost,
			envCfg.PostgreSQLUsername,
			envCfg.PostgreSQLPassword,
			envCfg.PostgreSQLDBName,
			envCfg.PostgreSQLPort),
		envCfg.PostgreSQLDBName
}

func Config() *envConfig {
	return &envCfg
}
