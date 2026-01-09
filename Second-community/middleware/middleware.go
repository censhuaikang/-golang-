package middleware

import (
	"community/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录或授权格式错误"})
			c.Abort()
			return
		}

		tokenString := authHeader[7:]
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌无效或已过期"})
			c.Abort()
			return
		}

		// 将用户ID存入上下文供后续使用
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
