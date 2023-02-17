package handlers

import (
	"Tiktok/controller/common"
	"Tiktok/dal/db"
	"Tiktok/pkg/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct{
	common.Response
	Video_list []common.Video
	Next_time int64
}

func GetFeed(c *gin.Context) {
	lastestTime := c.Query("latest_time")
	login, _ := c.Get("Login")
	time, err := strconv.Atoi(lastestTime)
	if err !=nil {
		log.Error("time to int error" + err.Error())
		c.JSON(http.StatusOK, FeedResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: "time to int error"},
		})
	}
	var video_list []common.Video
	var id int64
	if login == false {
		id = -1
	} else{
		tmp, _ := c.Get("ID")
		var ok bool
		if id, ok = tmp.(int64); !ok{
			log.Error("id to int error")
			c.JSON(http.StatusOK, FeedResponse{
				Response: common.Response{StatusCode: 1, StatusMsg: "id to int error"},
			})
			id = -1
		}
	}
	video_list, newtime := db.GetVedio(id, time)
	c.JSON(http.StatusOK, FeedResponse{
		Response: common.Response{StatusCode: 1},
		Video_list: video_list,
		Next_time: newtime,
	})
}
