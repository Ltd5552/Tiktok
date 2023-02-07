package user

import "Tiktok/controller"

func Info(userID string) (controller.User, bool) {
	return controller.User{}, true
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
