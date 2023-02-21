package jwt

import (
	"Tiktok/config"
	"Tiktok/pkg/log"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	ID   string
	Name string
	jwt.RegisteredClaims
}

// MySecret 设置secret
var MySecret = []byte(config.AuthSetting.JwtSecret)

// CreateToken 创建token
// 传入参数id，name(数据库ID、name)
// 返回参数token和错误信息
func CreateToken(id string, name string) (string, error) {
	claim := Claims{
		ID:   id,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(MySecret)
}

func Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	}
}

// ParseToken 解析token
// 传入参数token
// 返回数据库id和错误信息
func ParseToken(tokenstr string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenstr, &Claims{}, Secret())
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				log.Error("not a true token")
				return 0, errors.New("not a true token")
			} else {
				log.Error("token error")
				return 0, errors.New("token error")
			}
		}
	}
	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
		id, err := strconv.Atoi(claims.ID)
		if err != nil {
			return 0, errors.New("ID is not int")
		}
		return id, nil
	}
	log.Error("token parse error")
	return 0, errors.New("couldn't parse the token")
}

func VerifyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Query("token")
		if token == "" {
			ctx.Set("Login", false)
			ctx.Next()
		}
		id, err := ParseToken(token)
		if err != nil {
			ctx.Abort()
		}
		ctx.Set("ID", id)
		ctx.Set("Login", true)
		ctx.Next()
	}
}
