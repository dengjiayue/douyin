package middleware

import (
	"douyin/internal/gateway/db"
	my_jwt "douyin/pkg/jwt"
	"douyin/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 用户鉴权
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Authorization 字段（JWT Token）
		tokenString := c.GetHeader("Authorization")

		// 检查 Token 是否存在
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
			c.Abort()
			return
		}

		// 解析 Token
		claims, err := my_jwt.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
			c.Abort()
			return
		}

		// 将解析后的用户信息存储到上下文中，以便后续处理使用
		c.Set("user_id", claims.UserID)

		c.Next()
	}
}

// 用户鉴权:使用query参数或者form参数中的token
func AuthMiddlewareQueryOrForm(iGetWay db.IGetWay) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 Authorization 字段（JWT Token）
		tokenString := c.Query("token")
		if len(tokenString) == 0 {
			tokenString = c.PostForm("token")
		}
		id := iGetWay.LogInStatus(tokenString)
		// 将解析后的用户信息存储到上下文中，以便后续处理使用
		c.Set("user_id", id)
		logger.Debugf("token->uid: %d\n", id)
		c.Next()
	}
}
