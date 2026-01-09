package routes

import (
	"community/controllers"
	"community/middleware"
	"community/repositories"
	"community/services"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupUserRoutes(r *gin.Engine) {
	dsn := "root:censhuaikang@tcp(127.0.0.1:3306)/问答社区?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 创建服务和控制器
	userService := services.NewUserService(repositories.NewUserRepository(db))
	userCtrl := controllers.NewUserController(userService)

	// --- 公开路由 ---
	r.POST("/user/register", userCtrl.Register)
	r.POST("/user/login", userCtrl.Login)

	// --- 受保护路由组 ---
	auth := r.Group("/")
	auth.Use(middleware.JWTMiddleware())
	{
		auth.POST("/question/create", userCtrl.CreateQuestion)
		auth.POST("/answer/create", userCtrl.CreateAnswer)
		auth.POST("/answer/reply", userCtrl.Answer)
		auth.POST("/modify/question", userCtrl.ModifyQuestion)
		auth.POST("/modify/answer", userCtrl.ModifyAnswer)
		auth.POST("/question/delete", userCtrl.DeleteQuestion)
		auth.POST("/answer/delete", userCtrl.DeleteAnswer)
		auth.POST("/follow/:id", userCtrl.Follow)
		auth.POST("/unfollow/:id", userCtrl.Unfollow)
		auth.POST("/user/delete", userCtrl.DeleteUser)
	}
}
