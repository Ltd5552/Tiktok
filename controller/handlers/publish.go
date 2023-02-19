package handlers

import (
	"Tiktok/model"
	"Tiktok/pkg/jwt"
	"Tiktok/pkg/log"
	"Tiktok/pkg/minio"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
	"strconv"
)

type VideoListResponse struct {
	Response
	VideoList []*Video `json:"video_list"`
}

func PublishAction(c *gin.Context) {
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	token := c.Query("token")
	userID, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	finalName := fmt.Sprintf("%d_%s", userID, filename)
	//上传到minio
	saveFile := filepath.Join("./public/", finalName)
	if err := minio.UploadFile("", finalName, "data", saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 数据库创建
	// 创建video
	video := map[string]interface{}{
		"authorId": userID,
		"playUrl":  "",
		"coverUrl": "",
		"title":    "",
	}
	if err := model.CreatVideo(video); err != nil {
		log.Infos(c, "video upload err", zap.Error(err))
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "User register err"})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

func GetPublishList(c *gin.Context) {
	token := c.Query("token")
	userID, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	modelVideos, err := model.GetPublishedVideos(userID)
	var videos []*Video
	for i := 0; i < len(modelVideos); i++ {
		modelUser, _ := model.ReadUser(strconv.Itoa(int(modelVideos[i].AuthorId)))

		user := User{
			modelUser.ID,
			modelUser.Name,
		}

		videos[i] = &Video{modelVideos[i].ID,
			user,
			modelVideos[i].PlayUrl,
			modelVideos[i].CoverUrl,
			modelVideos[i].FavoriteCount,
			modelVideos[i].CommentCount,
			true,
			modelVideos[i].Title}
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: videos,
	})
}
