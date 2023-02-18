package handlers

import (
	"Tiktok/config"
	"Tiktok/controller"
	"Tiktok/pkg/hash"
	"Tiktok/pkg/log"
	"Tiktok/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

type UserLoginResponse struct {
	controller.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	controller.Response
	controller.User `json:"user"`
}

func UserRegister(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 合法性校验
	if !ValidateAccount(username, password) {
		log.Infos(c, "Account validate failed")
		// 这部分参考示例实现
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: controller.Response{
				StatusCode: 1,
				StatusMsg:  "Account validate failed"}})
		return
	}

	// 唯一性校验
	if exist := user.Exist(username); exist {
		log.Infos(c, "User already exist")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: controller.Response{
				StatusCode: 1,
				StatusMsg:  "User already exist"}})
		return
	}

	// 创建account
	account := map[string]interface{}{
		"name":     username,
		"password": hash.Md5WithSalt(password, config.AuthSetting.Md5Salt),
	}

	// 待添加token
	if userID, ok := user.Register(account); ok {
		log.Infos(c, "User register success")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: controller.Response{StatusCode: 0, StatusMsg: "User already exist"},
			UserId:   userID,
			Token:    "token"})
	} else {
		log.Infos(c, "User register err")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: controller.Response{
				StatusCode: 1,
				StatusMsg:  "User register err"}})
		return
	}
}

func UserLogin(c *gin.Context) {

	username := c.Query("username")
	password := c.Query("password")

	// 合法性校验
	if !ValidateAccount(username, password) {
		log.Infos(c, "Account validate failed")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: controller.Response{StatusCode: 1, StatusMsg: "Account validate failed"}})
		return
	}

	account := map[string]interface{}{
		"name":     username,
		"password": hash.Md5WithSalt(password, config.AuthSetting.Md5Salt),
	}

	if id, exist := user.Login(account); exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: controller.Response{StatusCode: 0},
			UserId:   id,
			Token:    "token",
		})
	} else {
		log.Infos(c, "User doesn't exist")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: controller.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func GetUserInfo(c *gin.Context) {

	//token := c.Query("token")

	userID := c.Query("user_id")

	if User, exist := user.Info(userID); exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: controller.Response{StatusCode: 0},
			User:     User,
		})
	} else {
		log.Infos(c, "User doesn't exist")
		c.JSON(http.StatusOK, UserResponse{
			Response: controller.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func ValidateAccount(username string, password string) bool {
	NameRegExp := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{3,32}$`) // 字母开头，允许3-32字节，允许字母数字下划线
	PasswordRegExp := regexp.MustCompile(`^[a-zA-Z]\w{6,32}$`)       // 字母开头，长度在6~32字节，只能包含字母、数字和下划线
	return NameRegExp.MatchString(username) && PasswordRegExp.MatchString(password)
}
