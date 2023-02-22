package model

import (
	"gorm.io/gorm"
)

func Like(userID int, videoID int) error {
	favorite := Favorite{
		UserId:  userID,
		VideoId: videoID,
	}
	// 开始事务
	tx := DB.Begin()

	if err := tx.Create(&favorite).Error; err != nil {
		tx.Rollback()
		return err
	}
	var video Video
	if err := tx.Where("id = ?", videoID).Find(&video).Error; err != nil {
		return err
	}
	// UPDATE "favorite_count" SET "favorite_count" = favorite_count + 1 WHERE "id" = videoID;
	if err := tx.Model(&video).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}
	//提交事务
	tx.Commit()
	return nil
}

func Dislike(userID int, videoID int) error {
	favorite := Favorite{
		UserId:  userID,
		VideoId: videoID,
	}
	// 开始事务
	tx := DB.Begin()
	//删除操作
	if err := tx.Where("user_id = ?", userID).Where("video_id = ?", videoID).Delete(favorite).Error; err != nil {
		tx.Rollback()
		return err
	}

	var video Video
	if err := tx.Where("id = ?", videoID).Find(&video).Error; err != nil {
		return err
	}
	// UPDATE "favorite_count" SET "favorite_count" = favorite_count - 1 WHERE "id" = videoID;
	if err := tx.Model(&video).UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}
	//提交事务
	tx.Commit()
	return nil
}

func GetVideoByLike(userID int) ([]*Video, error) {
	var videos []*Video
	//select *
	//form video
	//where video.id = (
	//select videoID
	//form favorites
	//where userID = userID)
	err := DB.Where("id = (?)", DB.Where("user_id = ?", userID).Table("favorites").Select("video_id")).Find(&videos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return videos, nil
}
