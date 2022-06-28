package configs

import (
	"os"
)

const (
	KeyAppEnv          string = "APP_ENV"
	KeyAppPort         string = "APP_PORT"
	KeyMysqlDbBame     string = "MYSQL_DBNAME"
	KeyMysqlHost       string = "MYSQL_HOST"
	KeyMysqlUser       string = "MYSQL_USER"
	KeyMysqlPw         string = "MYSQL_PASSWORD"
	KeyGoogleMapAPIKey string = "GOOGLE_MAP_API_KEY"
)

// Get get value from configs
func Get(key string) string {
	return os.Getenv(key)
}
