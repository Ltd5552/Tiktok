package db

import (
	"Tiktok/pkg/log"
	"time"
)

type video struct {
	Id             int    `gorm:"column:id;primaryKey"`
	Author_id      int    `gorm:"column:author_id"`
	Play_url       string `gorm:"column:play_url"`
	Cover_url      string `gorm:"column:cover_url"`
	Favorite_count int    `gorm:"column:favourite_count"`
	Comment_count  int    `gorm:"column:comment_count"`
	Title          string `gorm:"column:title"`
	Create_at      *time.Time `gorm:"column:create_at"`
	Update_at	   *time.Time `gorm:"column:update_at"`
	Delete_at	   *time.Time `gorm:"column:delete_at"`
}

func GetVedio(lastestTime int) []video{
	var videos []video
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