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

	tx := DB.Begin()
	if err := tx.Create(&comment).Error; err != nil {
		tx.Rollback()
		return Comment{}, err
	}
	var video Video
	if err := tx.Where("id = ?", videoId).Find(&video).Error; err != nil {
		return Comment{}, err
	}
	// UPDATE "comment_count" SET "comment_count" = comment_count + 1 WHERE "id" = videoID;
	if err := tx.Model(&video).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		return Comment{}, err
	}
	//提交事务
	tx.Commit()
	return *comment, nil
}

func DeleteComment(commentId uint) (Comment, error) {
	var comment Comment
	tx := DB.Begin()
	err := tx.Where("id = ?", commentId).Find(&comment).Error
	if err != nil {
		tx.Rollback()
		return Comment{}, err
	}
	var video Video
	video.ID = comment.VideoId
	err = DB.Delete(&Comment{}, commentId).Error
	if err != nil {
		tx.Rollback()
		return Comment{}, err
	}

	// UPDATE "comment_count" SET "comment_count" = comment_count - 1 WHERE "id" = videoID;
	if err := tx.Model(&video).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
		tx.Rollback()
		return Comment{}, err
	}
	//提交事务
	tx.Commit()

	return comment, nil
}

func GetComment(videoId string) ([]Comment, error) {
	var commentList []Comment
	if err := DB.Order("created_at desc").Where("video_id = ?", videoId).Find(&commentList).Error; err != nil {
		return nil, err
	}
	return commentList, nil
}
