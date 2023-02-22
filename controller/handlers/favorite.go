package handlers

import (
	"Tiktok/model"
	"Tiktok/pkg/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type FavoriteListResponse struct {
	Response
	VideoList []*Video `json:"video_list"`
}

func FavoriteAction(c *gin.Context) {
	videoId := c.Query("video_id")
	actionType := c.Query("action_type")
	userID, exit := c.Get("ID")
	if !exit {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "please login first"})
		return
	}
	videoID, err := strconv.Atoi(videoId)
	if err != nil {
		log.Errors(c, "video_id conv int failed", zap.Error(err))
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "video_id conv int failed"})
		return
	}
	if actionType == "1" {
		//点赞
		if err := model.Like(userID.(int), videoID); err != nil {
			log.Errors(c, "like video err:", zap.Error(err))
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "like video err"})
			return
		}

	} else if actionType == "2" {
		//取消点赞
		if err := model.Dislike(userID.(int), videoID); err != nil {
			log.Errors(c, "dislike video err:", zap.Error(err))
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "dislike video err"})
			return
		}
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "success",
	})
}

func GetFavoriteList(c *gin.Context) {
	userID, exit := c.Get("ID")
	if !exit {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "please login first"})
		return
	}
	// 如何权衡model.Video和handler.Video?
	modelVideos, err := model.GetVideoByLike(userID.(int))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error()})
		return
	}
	var videos []*Video
	for i := 0; i < len(modelVideos); i++ {
		modelUser, _ := model.ReadUser(strconv.Itoa(int(modelVideos[i].AuthorId)))
		user := User{modelUser.ID, modelUser.Name}
		videos[i] = &Video{modelVideos[i].ID,
			user,
			modelVideos[i].PlayUrl,
			modelVideos[i].CoverUrl,
			modelVideos[i].FavoriteCount,
			modelVideos[i].CommentCount,
			true,
			modelVideos[i].Title}
	}
	log.Infos(c, "get favorite lists")
	c.JSON(http.StatusOK, FavoriteListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: videos,
	})
}
