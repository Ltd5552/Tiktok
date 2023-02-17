package db

import (
	"Tiktok/controller/common"
	"Tiktok/pkg/log"
	"time"

	"gorm.io/gorm"
)

type video struct {
	Id             int64  `gorm:"column:id;primaryKey"`
	Author_id      int64  `gorm:"column:author_id"`
	Play_url       string `gorm:"column:play_url"`
	Cover_url      string `gorm:"column:cover_url"`
	Favorite_count int64  `gorm:"column:favourite_count"`
	Comment_count  int64  `gorm:"column:comment_count"`
	Title          string `gorm:"column:title"`
	Create_at      *time.Time `gorm:"column:create_at"`
	Update_at	   *time.Time `gorm:"column:update_at"`
	Delete_at	   *time.Time `gorm:"column:delete_at"`
}

type favorite struct{
	User_id int64 `gorm:"column:user_id"`
	Video_id int64 `gorm:"column:vidoe_id"`
}

func JudgeFavorite(user_id int64, video_id int64) bool{
	var tmp favorite
	if err := DB.Where("user_id = ?", user_id).Where("video_id = ?", video_id).Find(&tmp).Error; err == gorm.ErrRecordNotFound {
		return false
	} else if err !=nil {
		log.Error("select favorite error"+err.Error())
		return false
	}
	return true
}

func VideosConv(id int64, videos []video) ([]common.Video, int64){
	var convVedios []common.Video
	var lastestTime int64
	for _, video := range(videos){
		var convVedio common.Video
		convVedio.Id = video.Id
		convVedio.Author = GetUser(video.Author_id)
		convVedio.PlayUrl = video.Play_url
		convVedio.CoverUrl = video.Cover_url
		convVedio.CommentCount = video.Comment_count
		if id == -1{
			convVedio.IsFavorite = false
		} else{
			convVedio.IsFavorite = JudgeFavorite(id, video.Id)
		}
		convVedios = append(convVedios, convVedio)
		lastestTime = video.Create_at.Unix()
	}
	return convVedios, lastestTime
}

func GetVedio(id int64, lastestTime int) ([]common.Video, int64){
	var videos []video
	if lastestTime !=0 {
		time := time.Unix(int64(lastestTime), 0)
		err := DB.Where("create_at < ?", time).Order("create_at desc").Limit(30).Find(&videos).Error
		if err != nil{
			log.Error("select sql failed")
			return nil, 0
		}
		return VideosConv(id, videos)
	}
	err := DB.Order("create_at desc").Limit(30).Find(&videos).Error
	if err != nil{
		log.Error("select sql failed")
		return nil, 0
	}
	return VideosConv(id, videos)
}