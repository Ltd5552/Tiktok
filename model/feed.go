package model

import (
	"Tiktok/pkg/log"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Video struct {
	AuthorId      uint   `gorm:"column:author_id" json:"author_id"`
	PlayUrl       string `gorm:"column:play_url" json:"play_url"`
	CoverUrl      string `gorm:"column:cover_url" json:"cover_url"`
	FavoriteCount int64  `gorm:"column:favourite_count" json:"favorite_count"`
	CommentCount  int64  `gorm:"column:comment_count" json:"comment_count"`
	Title         string `gorm:"column:title" json:"title"`
	gorm.Model
}

type Favorite struct {
	UserId  uint `gorm:"column:user_id"`
	VideoId uint `gorm:"column:video_id"`
}

// 判断用户是否给该视频点赞
func JudgeFavorite(userId uint, videoId uint) bool {
	var tmp Favorite
	if err := DB.Where("user_id = ?", userId).Where("video_id = ?", videoId).Find(&tmp).Error; err == gorm.ErrRecordNotFound {
		return false
	} else if err != nil {
		log.Error("select Favorite error" + err.Error())
		return false
	}
	return true
}

// 从数据库中获取video信息

func GetVideo(id uint, lastestTime int) ([]Video, error){
	var videos []Video
	if lastestTime != 0 {
		time := time.Unix(int64(lastestTime), 0)
		err := DB.Where("create_at < ?", time).Order("create_at desc").Limit(30).Find(&videos).Error
		if err != nil{
			log.Error("select sql failed", zap.Error(err))
			return nil, err
		}
		return videos, nil
	}
	err := DB.Order("create_at desc").Limit(30).Find(&videos).Error
	if err != nil{
		log.Error("select sql failed", zap.Error(err))
		return nil, err
	}
	return videos, nil
}
