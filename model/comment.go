package model

import (
	"gorm.io/gorm"
)

type Comment struct {
	Text        string `gorm:"text"`
	CommenterId uint   `gorm:"commenter_id"`
	VideoId     uint   `gorm:"video_id"`
	gorm.Model
}

func CreateComment(userId int, videoId int, text string) (Comment, error) {
	comment := &Comment{
		Text:        text,
		CommenterId: uint(userId),
		VideoId:     uint(videoId),
	}

	if err := DB.Create(&comment).Error; err != nil {
		return Comment{}, err
	}
	return *comment, nil
}

func DeleteComment(commentId uint) (Comment, error) {
	var comment Comment
	err := DB.Where("ID = ?", commentId).Find(&comment).Error
	if err != nil {
		return Comment{}, err
	}
	err = DB.Delete(&Comment{}, commentId).Error
	if err != nil {
		return Comment{}, err
	}
	return comment, nil
}

func GetComment(videoId string) ([]Comment, error) {
	var commentList []Comment
	if err := DB.Order("created_at desc").Where("video_id = ?", videoId).Find(&commentList).Error; err != nil {
		return nil, err
	}
	return commentList, nil
}
