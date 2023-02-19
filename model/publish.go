package model

import "gorm.io/gorm"

func CreatVideo(data map[string]interface{}) error {
	video := &Video{
		AuthorId: data["authorId"].(uint),
		PlayUrl:  data["playUrl"].(string),
		CoverUrl: data["coverUrl"].(string),
		Title:    data["title"].(string),
	}
	if err := DB.Create(&video).Error; err != nil {
		return err
	}
	return nil
}

func GetPublishedVideos(userID int) ([]*Video, error) {
	var Videos []*Video
	err := DB.Where("author_id = ?", userID).Find(&Videos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return Videos, nil
}
