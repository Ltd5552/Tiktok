package handlers

import (
	"Tiktok/config"
	"Tiktok/model"
	"Tiktok/pkg/hash"
	"Tiktok/pkg/jwt"
	"Tiktok/pkg/log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserLoginResponse struct {
	Response
	UserId uint   `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User `json:"user"`
}

func UserRegister(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 合法性校验
	if !ValidateAccount(username, password) {
		log.Infos(c, "Account validate failed")
		// 这部分参考示例实现
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "Account validate failed"}})
		return
	}

	// 唯一性校验
	if exist := model.ExistUser(username); exist {
		log.Infos(c, "User already exist")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User already exist"}})
		return
	}

	// 创建account
	account := map[string]interface{}{
		"name":     username,
		"password": hash.Md5WithSalt(password, config.AuthSetting.Md5Salt),
	}

	if userID, err := model.CreateUser(account); err == nil {
		token, err := jwt.CreateToken(strconv.Itoa(int(userID)), username)
		if err != nil {
			log.Error("Create token error", zap.Error(err))
		}
		log.Infos(c, "User register success")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0, StatusMsg: "User register success"},
			UserId:   userID,
			Token:    token})
	} else {
		log.Infos(c, "User register err", zap.Error(err))
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
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
			Response: Response{StatusCode: 1, StatusMsg: "Account validate failed"}})
		return
	}

	account := map[string]interface{}{
		"name":     username,
		"password": hash.Md5WithSalt(password, config.AuthSetting.Md5Salt),
	}

	//存在性校验
	if exist := model.ExistUser(username); !exist {
		log.Infos(c, "User doesn't exist")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist"}})
		return
	}

	id := model.ValidateUser(account)
	if id == 0 {
		log.Infos(c, "username or password err")
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "username or password err"},
		})
		return
	}
	if token, err := jwt.CreateToken(strconv.Itoa(int(id)), username); err == nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   id,
			Token:    token})
	}
}

func GetUserInfo(c *gin.Context) {

	userID := c.Query("user_id")
	token := c.Query("token")
	_, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "please login first"}})
		return
	}
	if ModelUser, err := model.ReadUser(userID); err == nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User: User{
				Name: ModelUser.Name,
				Id:   ModelUser.ID},
		})
	} else {
		log.Infos(c, "User doesn't exist")
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func ValidateAccount(username string, password string) bool {
	NameRegExp := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{2,31}$`) // 字母开头，允许3-32字节，允许字母数字下划线
	PasswordRegExp := regexp.MustCompile(`^[a-zA-Z]\w{5,31}$`)       // 字母开头，长度在6~32字节，只能包含字母、数字和下划线
	return NameRegExp.MatchString(username) && PasswordRegExp.MatchString(password)
}
