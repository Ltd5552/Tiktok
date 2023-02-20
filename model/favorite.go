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

func GetList(userID int) ([]*Video, error) {
	var videos []*Video
	//select *
	//form video
	//where video.id = (
	//select videoID
	//form favorite
	//where userID = userID)
	subQuery := DB.Select("video_id").Where("user_id = ?", userID).Table("favorite")
	err := DB.Where("video_id = ?").Where(subQuery).Find(&videos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return videos, nil
}
