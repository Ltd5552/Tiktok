package model

import (
	"Tiktok/pkg/log"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Video struct {
	AuthorId      uint   `gorm:"column:author_id"`
	PlayUrl       string `gorm:"column:play_url"`
	CoverUrl      string `gorm:"column:cover_url"`
	FavoriteCount int64  `gorm:"column:favorite_count"`
	CommentCount  int64  `gorm:"column:comment_count"`
	Title         string `gorm:"column:title"`
	gorm.Model
}

type Favorite struct {
	UserId  int `gorm:"column:user_id"`
	VideoId int `gorm:"column:video_id"`
}

// JudgeFavorite 判断用户是否给该视频点赞
func JudgeFavorite(userId int, videoId uint) bool {
	var favorite Favorite
	err := DB.Where("user_id = ? AND video_id = ?", userId, videoId).Find(&favorite).Error
	if err != nil {
		log.Error("select Favorite error" + err.Error())
		return false
	}
	if favorite.VideoId == 0 {
		return false
	}
	return true
}

// GetVideoByTime 从数据库中获取video信息
func GetVideoByTime(latestTime int64) ([]Video, error) {
	var videos []Video
	Time := time.Unix(latestTime, 0)
	if latestTime != 0 && Time.Before(time.Now()) {
		err := DB.Where("created_at < ?", Time).Order("created_at desc").Limit(30).Find(&videos).Error
		if err != nil {
			log.Error("select sql failed", zap.Error(err))
			return nil, err
		}
		return videos, nil
	}
	err := DB.Order("created_at desc").Limit(30).Find(&videos).Error
	if err != nil {
		log.Error("select sql failed", zap.Error(err))
		return nil, err
	}
	return videos, nil
}
