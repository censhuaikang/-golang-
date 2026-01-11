package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DBConn string
}

func LoadConfig() *Config {
	err := godotenv.Load() // Load .env file for environment variables
	if err != nil {
		panic("Error loading .env file")
	}
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD") //密码
	host := os.Getenv("DB_HOST")         //数据库地址，可以是Ip或者域名
	port := os.Getenv("DB_PORT")         //数据库端口
	Dbname := os.Getenv("DB_NAME")       //数据库名
	timeout := os.Getenv("DB_TIMEOUT")   //连接超时，30s
	dbConn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	return &Config{
		DBConn: dbConn,
	}
}
