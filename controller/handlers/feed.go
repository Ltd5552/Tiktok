package handlers

import (
	"Tiktok/model"
	"Tiktok/pkg/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct{
	Response
	Video_list []Video
	Next_time int64
}

// 将video数据从数据库结构体转成response结构体
func VideosConv(id uint, videos []model.Video) ([]Video, int64){
	var convVideos []Video
	var lastestTime int64
	for _, video := range(videos){
		var convVideo Video
		convVideo.Id = video.Model.ID
		user, err := model.ReadUser(strconv.Itoa(int(video.Author_id)))
		if err != nil{
			convVideo.Author = User{}
		} else{
			convVideo.Author = User{
				Id: user.Model.ID,
				Name: user.Name,
			}
		}
		convVideo.PlayUrl = video.Play_url
		convVideo.CoverUrl = video.Cover_url
		convVideo.CommentCount = video.Comment_count
		if id == 0{
			convVideo.IsFavorite = false
		} else{
			convVideo.IsFavorite = model.JudgeFavorite(id, video.Model.ID)
		}
		convVideos = append(convVideos, convVideo)
		// 记录创建最晚时间
		lastestTime = video.Model.CreatedAt.Unix()
	}
	return convVideos, lastestTime
}

func GetFeed(c *gin.Context) {
	// 获取是否登录和上次最晚时间
	lastestTime := c.Query("latest_time")
	login, _ := c.Get("Login")
	time, err := strconv.Atoi(lastestTime)
	if err !=nil {
		log.Error("time to int error" + err.Error())
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "time to int error"},
		})
	}
	var video_list []Video
	var id uint
	if login == false {
		id = 0
	} else{
		// 获取登录者的ID信息
		tmp, _ := c.Get("ID")
		var ok bool
		if id, ok = tmp.(uint); !ok{
			log.Error("id to int error")
			c.JSON(http.StatusOK, FeedResponse{
				Response: Response{StatusCode: 1, StatusMsg: "id to int error"},
			})
			id = 0
		}
	}
	model_video_list, err := model.GetVideo(id, time)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	}
	video_list, new_time := VideosConv(id, model_video_list)
	c.JSON(http.StatusOK, FeedResponse{
		Response: Response{StatusCode: 0},
		Video_list: video_list,
		Next_time: new_time,
	})
}
