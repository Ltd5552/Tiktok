package middleware

import (
	"Tiktok/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// 验证token有效性，无效则终止后续操作，有效则传递id值
func VerfyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query("token")
		id, err := jwt.ParseToken(token)
		if err !=nil {
			ctx.Abort()
		}
		ctx.Set("ID", id)
		ctx.Next()
	}
}