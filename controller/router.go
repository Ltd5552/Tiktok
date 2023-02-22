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

	// 设置文件最大容量
	r.MaxMultipartMemory = 2 << 26

	// 初始化trace
	trace.Set(r)

	// 初始化metric
	metric.Set(r)

	// feed组，视频流，无需验证登录状态
	feed := r.Group("/douyin/feed")
	feed.GET("/", handlers.GetFeed)

	// user组，用户，只有信息需要验证token
	user := r.Group("/douyin/user")
	user.GET("/", handlers.GetUserInfo)
	user.POST("/login/", handlers.UserLogin)
	user.POST("/register/", handlers.UserRegister)

	douyin := r.Group("/douyin")

	// publish组，投稿，action的token在参数中
	publish := douyin.Group("/publish")
	publish.POST("/action/", handlers.PublishAction)
	publish.GET("/list/", handlers.GetPublishList)

	// favorite组，喜欢
	favorite := douyin.Group("/favorite")
	favorite.Use(jwt.VerifyMiddleware())
	favorite.POST("/action/", handlers.FavoriteAction)
	favorite.GET("/list/", handlers.GetFavoriteList)

	// comment组，评论
	comment := douyin.Group("/comment")
	comment.POST("/action/", handlers.CommentAction)
	comment.GET("/list/", handlers.GetCommentList)

	return r
}
