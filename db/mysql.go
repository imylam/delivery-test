package db

import (
	"fmt"
	"log"

	"github.com/imylam/delivery-test/configs"

	"github.com/jmoiron/sqlx"

	// For mysql connection
	_ "github.com/go-sql-driver/mysql"
)

var mysqlConn *sqlx.DB

// InitDBConn initialize database connection
func InitDBConn() {
	if mysqlConn == nil {
		mysqlConn = connectMysql()
	}

	err := mysqlConn.Ping()
	if err != nil {
		log.Fatal("Error on connection to database:", err)
	}

	log.Print("Success in connecting to database")
}

// GetDBConnection get database connection object
func GetDBConnection() *sqlx.DB {
	return mysqlConn
}

// connectMysql connects to mysql/mariadb database
func connectMysql() *sqlx.DB {
	dbHost, dbName, dbPassword, dbUser := getDbConfigs()
	//connectionStr := fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True", dbUser, dbPassword, dbName)
	connectionStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True", dbUser, dbPassword, dbHost, dbName)

	mysqlCon, err := sqlx.Open("mysql", connectionStr)
	if err != nil {
		log.Fatal(err)
	}

	return mysqlCon
}

func getDbConfigs() (dbHost, dbName, dbPassword, dbUser string) {
	dbHost = configs.Get(configs.KeyMysqlHost)
	dbName = configs.Get(configs.KeyMysqlDbBame)
	dbPassword = configs.Get(configs.KeyMysqlPw)
	dbUser = configs.Get(configs.KeyMysqlUser)

	return
}
