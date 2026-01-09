//package main
//
//import (
//	routes "community/router"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"gorm.io/gorm"
//	"net/http"
//	"time"
//)
//
//var users = make(map[string]string) // 存储用户名和密码的映射关系
//
//func string1(c *gin.Context) {
//	c.String(200, "hello")
//}
//
//func main() {
//	http.HandleFunc("/register", registerHandler)
//	http.HandleFunc("/login", loginHandler)
//	http.ListenAndServe(":8080", nil)
//}
//
//func registerHandler(w http.ResponseWriter, r *http.Request) {
//	if r.Method == "POST" {
//		username := r.FormValue("username")
//		password := r.FormValue("password")
//
//		if _, ok := users[username]; ok {
//			// 用户名已存在
//			http.Error(w, "Username already exists", 200)
//			return
//		}
//		users[username] = password
//		fmt.Fprintf(w, "Registration successful!")
//	}
//}
//
//func loginHandler(w http.ResponseWriter, r *http.Request) {
//	if r.Method == "POST" {
//		username := r.FormValue("username")
//		password := r.FormValue("password")
//
//		storedPassword, ok := users[username]
//		if !ok || storedPassword != password {
//			// 用户名不存在或密码不匹配
//			http.Error(w, "Invalid username or password", 200)
//			return
//		}
//
//		// 设置cookie记录登录状态
//		cookie := &http.Cookie{
//			Name:    "session",
//			Value:   username,
//			Path:    "/",
//			Expires: time.Now().Add(24 * time.Hour),
//		}
//		http.SetCookie(w, cookie)
//		fmt.Fprintf(w, "Login successful!")
//	}
//}

// package main
//
// import (
//
//	"community/config"
//	"community/routes"
//	"github.com/gin-gonic/gin"
//	"gorm.io/gorm"
//
// )
//
//	func main() {
//		config := configs.LoadConfig()
//		db, err := gorm.Open(gorm.Dialectic(), config.DBConn)
//		if err != nil {
//			panic("Failed to connect database")
//		}
//
//		r := gin.Default()
//		routes.SetupUserRoutes(r)
//		// 设置其他路由...
//
//		err1 := r.Run(":8080")
//		if err1 != nil {
//			return
//		}
//
// }
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

	config := configs.LoadConfig()
	db, err := gorm.Open(mysql.Open(config.DBConn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database")
	}

	// 迁移数据库模式
	err1 := db.AutoMigrate(&models.User{}, &models.Question{}, &models.Answer{}, &models.Follow{})
	if err1 != nil {
		return
	}

	// 设置路由
	r := gin.Default()
	routes.SetupUserRoutes(r)

	// 运行服务器
	err2 := r.Run(":8080")
	if err2 != nil {
		return
	}
}
