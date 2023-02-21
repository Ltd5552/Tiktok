package model

import (
	"gorm.io/gorm"
)

type User struct {
	Name           string  `json:"name"`
	Password       string  `json:"password"`
	FavoriteVideos []Video `gorm:"many2many:Favorites" json:"video_list"`
	gorm.Model
}

func CreateUser(data map[string]interface{}) (uint, error) {
	user := &User{
		Name:     data["name"].(string),
		Password: data["password"].(string),
	}
	if err := DB.Create(&user).Error; err != nil {
		return 0, err
	}
	return user.ID, nil
}

func ReadUser(id string) (User, error) {
	var user User
	err := DB.First(&user, id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return user, err
	}
	return user, nil
}

func ValidateUser(data map[string]interface{}) uint {
	user := &User{
		Name:     data["name"].(string),
		Password: data["password"].(string),
	}
	DB.Where("name = ?", user.Name).Where("password = ?", user.Password).Find(&user)
	return user.ID
}

func ExistUser(username string) bool {
	var user User
	err := DB.First(&user, "name = ?", username).Error
	//若没有这条记录就返回false
	if err == gorm.ErrRecordNotFound {
		return false
	}
	return true
}
