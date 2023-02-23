package handlers

import (
	"Tiktok/model"
	"Tiktok/pkg/jwt"
	"Tiktok/pkg/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list"`
	NextTime  int64   `json:"next_time"`
}

// VideosConv 将video数据从数据库结构体转成response结构体
func VideosConv(userID int, videos []model.Video) ([]Video, int64) {
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
				Name: user.Name}
		}
		convVideo.PlayUrl = video.PlayUrl
		convVideo.CoverUrl = video.CoverUrl
		convVideo.CommentCount = video.CommentCount
		if userID == 0 {
			convVideo.IsFavorite = false
		} else {
			convVideo.IsFavorite = model.JudgeFavorite(userID, video.Model.ID)
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
	var time int64
	if latestTime != "" {
		var err error
		time, err = strconv.ParseInt(latestTime, 10, 64)
		if err != nil {
			log.Errors(c, "time to int error"+err.Error())
			c.JSON(http.StatusOK, FeedResponse{
				Response: Response{StatusCode: 1, StatusMsg: "time to int error"},
			})
			return
		}
	}
	var videoList []Video
	var userID int
	// 获取登录者的ID信息
	token := c.Query("token")

	if token == "" {
		userID = 0
	}
	userID, _ = jwt.ParseToken(token)
	modelVideoList, err := model.GetVideoByTime(time)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: err.Error()},
		})
	}
	videoList, newTime := VideosConv(userID, modelVideoList)
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		NextTime:  newTime,
		VideoList: videoList,
	})
}
