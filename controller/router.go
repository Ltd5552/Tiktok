package controller

import (
	"Tiktok/controller/handlers"
	"Tiktok/pkg/jwt"
	"Tiktok/pkg/metric"
	"Tiktok/pkg/trace"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	r := gin.New()

	// 初始化trace
	trace.Set(r)

	// 初始化metric
	metric.Set(r)

	// feed组，视频流，无需验证登录状态
	feed := r.Group("/douyin/feed")
	feed.GET("/", handlers.GetFeed)

	// 创建路由组，jwt验证
	douyin := r.Group("/douyin")
	douyin.Use(jwt.VerifyMiddleware())

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

	return r
}
