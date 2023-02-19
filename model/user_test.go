package model

import (
	"Tiktok/config"
	"Tiktok/pkg/hash"
	"fmt"
	"strconv"
	"testing"
)

func TestCreatAndReadUser(t *testing.T) {
	//测试的时候要把配置文件放在当前目录下，测完记得把多的文件（配置文件、log文件）删掉
	config.InitViper()
	InitDB()
	// 创建account
	account := map[string]interface{}{
		"name":     "test",
		"password": hash.Md5WithSalt("password", config.AuthSetting.Md5Salt),
	}
	userID, err := CreateUser(account)
	if err == nil {
		fmt.Println(userID)
	}
	user, err := ReadUser(strconv.Itoa(int(userID)))
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(user)
}
