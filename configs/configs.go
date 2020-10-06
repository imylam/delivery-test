package configs

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	KeyAppPort         string = "APP_PORT"
	KeyMysqlHost       string = "MYSQL_HOST"
	KeyMysqlDbBame     string = "MYSQL_DBNAME"
	KeyMysqlUser       string = "MYSQL_USER"
	KeyMysqlPw         string = "MYSQL_PASSWORD"
	KeyGoogleMapAPIKey string = "GOOGLE_MAP_API_KEY"
)

var configs map[string]string

// Init loading the configuration
func Init() error {
	err := loadConfigs()

	return err
}

// Get get value from configs
func Get(key string) string {
	return os.Getenv(key)
}

func loadConfigs() error {
	err := godotenv.Load()

	return err
}
