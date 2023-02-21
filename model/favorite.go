package model

import "gorm.io/gorm"

func Like(data map[string]interface{}) error {
	favorite := &Favorite{
		UserId:  data["userID"].(uint),
		VideoId: data["videoID"].(uint),
	}
	if err := DB.Create(&favorite).Error; err != nil {
		return err
	}
	return nil
}

func Dislike(data map[string]interface{}) error {
	favorite := &Favorite{
		UserId:  data["userID"].(uint),
		VideoId: data["videoID"].(uint),
	}
	if err := DB.Delete(&favorite).Error; err != nil {
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
