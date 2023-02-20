package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	Text         string `gorm:"text"`
	Commenter_id uint   `gorm:"commenter_id"`
	Video_id     uint   `gorm:"video_id"`
	gorm.Model
}

func CreateComment(user_id uint, video_id uint, text string) (Comment, error){
	comment := &Comment{
		Text: text,
		Commenter_id: user_id,
		Video_id: video_id,
	}

	if err := DB.Create(&comment).Error; err !=nil {
		return Comment{}, err
	}
	return *comment, nil
}

func DeleteComment(comment_id uint) (Comment, error){
	var comment Comment
	err := DB.Where("ID = ?", comment_id).Find(&comment).Error
	if err !=nil {
		return Comment{}, err
	}
	err = DB.Delete(&Comment{}, comment_id).Error
	if err !=nil {
		return Comment{}, err
	}
	return comment, nil
}

func GetComment(video_id string) ([]Comment, error){
	var comment_list []Comment
	if err := DB.Order("create_at desc").Where("video_id = ?", video_id).Find(&comment_list).Error; err != nil {
		return nil, err
	}
	return comment_list, nil
}