package config

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var conf Config

type Config struct {
	DefaultUsername string
	DefaultPassword string
	Port            string
	JwtKey          string
	DSN             string
	initialized     bool
}

func GetConfig() (config *Config) {

	config = &conf
	if config.initialized {
		return
	}
	conf = Config{}
	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatal("$PORT should be set")
	}
	conf.Port = port

	username := os.Getenv("USER")
	if len(username) == 0 {
		log.Fatal("$USER should be set")
	}
	conf.DefaultUsername = username

	password := os.Getenv("PASSWORD")
	if len(password) == 0 {
		log.Fatal("$PASSWORD should be set")
	}
	conf.DefaultPassword = password

	key := os.Getenv("JWT_KEY")
	if len(key) == 0 {
		log.Fatal("$JWT_KEY should be set")
	}
	conf.JwtKey = key

	dbUser := os.Getenv("MYSQL_USER")
	if len(dbUser) == 0 {
		log.Fatal("$MYSQL_USER should be set")
	}

	dbPassword := os.Getenv("MYSQL_PASSWORD")
	if len(dbPassword) == 0 {
		log.Fatal("$MYSQL_PASSWORD should be set")
	}

	dbName := os.Getenv("MYSQL_DATABASE")
	if len(dbName) == 0 {
		log.Fatal("$MYSQL_DATABASE should be set")
	}

	dbHost := os.Getenv("MYSQL_HOST")
	if len(dbHost) == 0 {
		log.Fatal("$MYSQL_HOST should be set")
	}

	dbPort := os.Getenv("MYSQL_PORT")
	if len(dbPort) == 0 {
		log.Fatal("$MYSQL_PORT should be set")
	}

	dbConf := mysql.NewConfig()
	dbConf.User = dbUser
	dbConf.Passwd = dbPassword
	dbConf.DBName = dbName
	dbConf.Net = "tcp"
	dbConf.Addr = fmt.Sprintf("%s:%s", dbHost, dbPort)
	dbConf.ParseTime = true
	dbConf.Params = map[string]string{"charset": "utf8", "time_zone":"'UTC'"}
	fmt.Println(dbConf.FormatDSN())
	conf.DSN = dbConf.FormatDSN()
	return
}
