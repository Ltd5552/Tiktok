package service

import "Tiktok/controller/common"

func Info(userID string) (common.User, bool) {
	return common.User{}, true
}

func Register(map[string]interface{}) (int64, bool) {
	return 1, true
}

func Exist(username string) bool {
	return true
}

func Login(map[string]interface{}) (int64, bool) {
	return 1, true
}
