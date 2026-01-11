package main

import (
	"community/config"
	"community/models"
	"community/routes"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()                                   //加载配置
	db, err := gorm.Open(mysql.Open(cfg.DBConn), &gorm.Config{}) //链接数据库
	if err != nil {
		panic("Failed to connect database: " + err.Error())
	}

	// 迁移数据库模式
	err = db.AutoMigrate(
		&models.User{},
		&models.Question{},
		&models.Answer{},
		&models.Follow{},
		&models.Reply{},
	)
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}

	r := gin.Default()
	// 将 db 传递进去
	routes.SetupUserRoutes(r, db)

	r.Run(":8080")
}
