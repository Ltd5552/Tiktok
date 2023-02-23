package handlers

import (
	"Tiktok/config"
	"Tiktok/model"
	"Tiktok/pkg/jwt"
	"Tiktok/pkg/log"
	"Tiktok/pkg/minio"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"go.uber.org/zap"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
)

type VideoListResponse struct {
	Response
	VideoList []*Video `json:"video_list"`
}

func PublishAction(c *gin.Context) {
	token := c.PostForm("token")
	userID, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "please login first"})
		return
	}

	//读取fileHeader
	fileHeader, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error()})
		return
	}
	//从header打开关联的file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error()})
		return
	}

	videoName := fmt.Sprintf("%d_%s", userID, filepath.Base(fileHeader.Filename))

	//将file类型读取为byte
	videoByte, err := io.ReadAll(file)
	if err != nil {
		log.Error("Read file error")
	}
	//上传视频到minio
	if err := minio.UploadFile("video", videoByte, videoName, "video/mp4"); err != nil {
		log.Infos(c, "video upload err", zap.Error(err))
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error()})
		return
	}
	//获取视频链接
	videoUrl := fmt.Sprintf("http://%s:%s/video/%s", config.MinioSetting.Host, config.MinioSetting.Port, videoName)
	//获取视频封面
	coverByte, err := getCover(videoUrl, 1)
	if err != nil {
		log.Infos(c, "cover get err", zap.Error(err))
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error()})
		return
	}

	//上传图片到minio
	if err := minio.UploadFile("cover", coverByte, videoName, "image/jpeg"); err != nil {
		log.Infos(c, "cover upload err", zap.Error(err))
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error()})
		return
	}
	//获取封面链接
	coverUrl := fmt.Sprintf("http://%s:%s/cover/%s.jpg", config.MinioSetting.Host, config.MinioSetting.Port, videoName)
	//coverUrl := "https://images.ltd7.ltd/img/2022-summary/summer.jpg"
	//创建video并将视频和封面链接存入数据库
	video := map[string]interface{}{
		"authorId": uint(userID),
		"playUrl":  videoUrl,
		"coverUrl": coverUrl,
		"title":    videoName,
	}
	if err := model.CreatVideo(video); err != nil {
		log.Infos(c, "video model save err", zap.Error(err))
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error()})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  videoName + " uploaded successfully",
	})
}

func GetPublishList(c *gin.Context) {
	token := c.Query("token")
	userID, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "please login first"})
		return
	}

	modelVideos, err := model.GetPublishedVideos(userID)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error()})
		return
	}
	var videos []*Video
	for _, modelVideo := range modelVideos {
		modelUser, _ := model.ReadUser(strconv.Itoa(int(modelVideo.AuthorId)))
		user := User{modelUser.ID, modelUser.Name}

		video := &Video{modelVideo.ID,
			user,
			modelVideo.PlayUrl,
			modelVideo.CoverUrl,
			modelVideo.FavoriteCount,
			modelVideo.CommentCount,
			model.JudgeFavorite(userID, modelVideo.ID),
			modelVideo.Title}
		videos = append(videos, video)
	}
	log.Infos(c, "get published lists")
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: videos,
	})
}

func getCover(videoPath string, frameNum int) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg_go.Input(videoPath).Filter("select", ffmpeg_go.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg_go.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf).
		Run()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
