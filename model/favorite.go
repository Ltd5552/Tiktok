package model

import (
	"gorm.io/gorm"
)

func Like(userID int, videoID int) error {
	favorite := Favorite{
		UserId:  userID,
		VideoId: videoID,
	}
	if err := DB.Create(&favorite).Error; err != nil {
		return err
	}
	return nil
}

func Dislike(userID int, videoID int) error {
	favorite := Favorite{
		UserId:  userID,
		VideoId: videoID,
	}
	if err := DB.Where("user_id = ?", userID).Where("video_id = ?", videoID).Delete(favorite).Error; err != nil {
		return err
	}
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
