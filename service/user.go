package user

import (
	"Tiktok/controller/handlers"
	"Tiktok/model"
	"fmt"
)

// Exist 判断用户名是否已经存在
func Exist(username string) bool {
	return model.ExistUser(username)
}

// Register 注册
func Register(account map[string]interface{}) (uint, error) {
	userID, err := model.CreateUser(account)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// Login 登录
func Login(map[string]interface{}) (uint, bool) {
	return 1, true
}

// Info 通过userID获取用户信息
func Info(userID string) (handlers.User, bool) {
	var user handlers.User
	userData, err := model.ReadUser(userID)
	user.Id = userData.ID
	user.Name = userData.Name
	fmt.Println(userData)
	if err != nil {
		return handlers.User{}, false
	}
	return user, true
}
