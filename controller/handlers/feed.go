package handlers

import (
	"Tiktok/model"
	"Tiktok/pkg/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video
	NextTime  int64
}

// VideosConv 将video数据从数据库结构体转成response结构体
func VideosConv(id uint, videos []model.Video) ([]Video, int64) {
	var convVideos []Video
	var latestTime int64
	for _, video := range videos {
		var convVideo Video
		convVideo.Id = video.Model.ID
		user, err := model.ReadUser(strconv.Itoa(int(video.AuthorId)))
		if err != nil {
			convVideo.Author = User{}
		} else {
			convVideo.Author = User{
				Id:   user.Model.ID,
				Name: user.Name,
			}
		}
		convVideo.PlayUrl = video.PlayUrl
		convVideo.CoverUrl = video.CoverUrl
		convVideo.CommentCount = video.CommentCount
		if id == 0 {
			convVideo.IsFavorite = false
		} else {
			convVideo.IsFavorite = model.JudgeFavorite(id, video.Model.ID)
		}
		convVideos = append(convVideos, convVideo)
		// 记录创建最晚时间
		latestTime = video.Model.CreatedAt.Unix()
	}
	return convVideos, latestTime
}

func GetFeed(c *gin.Context) {
	// 获取是否登录和上次最晚时间
	latestTime := c.Query("latest_time")
	login, _ := c.Get("Login")
	time, err := strconv.Atoi(latestTime)
	if err != nil {
		log.Errors(c, "time to int error"+err.Error())
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "time to int error"},
		})
	}
	var videoList []Video
	var id uint
	if login == false {
		id = 0
	} else {
		// 获取登录者的ID信息
		tmp, _ := c.Get("ID")
		var ok bool
		if id, ok = tmp.(uint); !ok {
			log.Errors(c, "id to int error")
			c.JSON(http.StatusOK, FeedResponse{
				Response: Response{StatusCode: 1, StatusMsg: "id to int error"},
			})
			id = 0
		}
	}
	modelVideoList, err := model.GetVideo(time)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	}
	videoList, newTime := VideosConv(id, modelVideoList)
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  newTime,
	})
}
