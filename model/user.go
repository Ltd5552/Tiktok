package model

import (
	"gorm.io/gorm"
)

type User struct {
	Name           string  `json:"name"`
	Password       string  `json:"password"`
	FavoriteVideos []Video `gorm:"many2many:Favorite" json:"video_list"`
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

func ValidateUser(data map[string]interface{}) (uint, error) {
	user := &User{
		Name:     data["name"].(string),
		Password: data["password"].(string),
	}
	err := DB.Where("name = ? AND password = ?", user.Name, user.Password).Find(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, err
	}
	return user.ID, nil
}

func ExistUser(username string) bool {
	var user User
	err := DB.First(&user, "username = ?", username).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false
	}
	return true
}
