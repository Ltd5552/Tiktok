package main

import (
	"Tiktok/controller/handlers"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter() {

	r := gin.New()

	// 创建路由组
	douyin := r.Group("/douyin")

	// feed组，视频流
	feed := douyin.Group("/feed")
	feed.GET("/", handlers.GetFeed)

	// user组，用户
	user := douyin.Group("/user")
	user.GET("/", handlers.GetUserInfo)
	user.POST("/register", handlers.UserRegister)
	user.POST("/login", handlers.UserLogin)

	// publish组，投稿
	publish := douyin.Group("/publish")
	publish.POST("/action", handlers.PublishAction)
	publish.GET("/list", handlers.GetPublishList)

	// favorite组，喜欢
	favorite := douyin.Group("/favorite")
	favorite.POST("/action", handlers.FavoriteAction)
	favorite.GET("/list", handlers.GetFavoriteList)

	// comment组，评论
	comment := douyin.Group("/comment")
	comment.POST("/action", handlers.CommentAction)
	comment.GET("/list", handlers.GetCommentList)

	// relation组，关注
	relation := douyin.Group("/relation")
	relation.POST("/action", handlers.RelationAction)
	relation.GET("/follow/list", handlers.GetRelationFollowList)
	relation.GET("/follower/list", handlers.GetRelationFollowerList)

	// 为什么使用原生的http而不直接r.Run()?
	err := http.ListenAndServe(":8088", r)
	if err != nil {
		fmt.Println("log...")
	}
}
