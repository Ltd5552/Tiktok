package handlers

import (
	"Tiktok/model"
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
	userID, exit := c.Get("ID")
	if exit == false {
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
	//上传视频到minio
	if err := minio.UploadFile("video", videoByte, videoName, "video/mp4"); err != nil {
		log.Infos(c, "video upload err", zap.Error(err))
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error()})
		return
	}
	//获取视频链接
	videoUrl, err := minio.GetFile("video", videoName)
	if err != nil {
		log.Infos(c, "videoUrl get err", zap.Error(err))
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error()})
		return
	}
	//获取视频封面
	coverByte, err := getCover(videoUrl.String(), 1)
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
	coverUrl, err := minio.GetFile("cover", videoName)
	if err != nil {
		log.Infos(c, "coverUrl get err", zap.Error(err))
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error()})
		return
	}
	//创建video并将视频和封面链接存入数据库
	video := map[string]interface{}{
		"authorId": userID,
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
	userID, exit := c.Get("ID")
	if exit == false {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "please login first"})
		return
	}

	modelVideos, err := model.GetPublishedVideos(userID.(int))
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
