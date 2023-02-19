package model

import (
	"Tiktok/pkg/log"
	"time"

	"gorm.io/gorm"
)

type Video struct {
	Author_id      uint  `gorm:"column:author_id"`
	Play_url       string `gorm:"column:play_url"`
	Cover_url      string `gorm:"column:cover_url"`
	Favorite_count int64  `gorm:"column:favourite_count"`
	Comment_count  int64  `gorm:"column:comment_count"`
	Title          string `gorm:"column:title"`
	gorm.Model
}

type favorite struct{
	User_id int64 `gorm:"column:user_id"`
	Video_id int64 `gorm:"column:vidoe_id"`
}

// 判断用户是否给该视频点赞
func JudgeFavorite(user_id uint, video_id uint) bool{
	var tmp favorite
	if err := DB.Where("user_id = ?", user_id).Where("video_id = ?", video_id).Find(&tmp).Error; err == gorm.ErrRecordNotFound {
		return false
	} else if err !=nil {
		log.Error("select favorite error"+err.Error())
		return false
	}
	return true
}

// 从数据库中获取video信息
func GetVideo(id uint, lastestTime int) ([]Video){
	var videos []Video
	if lastestTime !=0 {
		time := time.Unix(int64(lastestTime), 0)
		err := DB.Where("create_at < ?", time).Order("create_at desc").Limit(30).Find(&videos).Error
		if err != nil{
			log.Error("select sql failed")
			return nil
		}
		return videos
	}
	err := DB.Order("create_at desc").Limit(30).Find(&videos).Error
	if err != nil{
		log.Error("select sql failed")
		return nil
	}
	return videos
}