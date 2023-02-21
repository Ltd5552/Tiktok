package handlers

import (
	"Tiktok/model"
	"Tiktok/pkg/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommentResponse struct {
	Response
	Comment Comment
}

type CommentListResponse struct {
	Response
	Comment []Comment
}

func CommentConv(comment model.Comment) (Comment, error) {
	var newComment Comment
	newComment.Id = comment.Model.ID
	tmpUser, err := model.ReadUser(strconv.Itoa(int(comment.CommenterId)))
	if err != nil {
		return Comment{}, err
	}
	newComment.User = User{
		Id:   tmpUser.ID,
		Name: tmpUser.Name,
	}
	newComment.Content = comment.Text
	newComment.CreateDate = comment.CreatedAt.Format("01-02")
	return newComment, nil
}

func CommentAction(c *gin.Context) {
	login, _ := c.Get("Login")
	tmpVideoId := c.Query("video_id")
	actionType := c.Query("action_type")
	var id uint
	if login == false {
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "comment need login first",
			},
			Comment: Comment{},
		})
		return
	} else {
		tmp, _ := c.Get("ID")
		var ok bool
		if id, ok = tmp.(uint); !ok {
			log.Errors(c, "id to int error")
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: "id to int error"},
			})
			return
		}
	}
	if actionType == "1" {
		commentText := c.Query("comment_text")
		videoId, err := strconv.Atoi(tmpVideoId)
		if err != nil {
			log.Errors(c, "video_id conv int failed", zap.Error(err))
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: "video_id conv int failed"},
			})
			return
		}
		comment, err := model.CreateComment(id, uint(videoId), commentText)
		if err != nil {
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
		}
		newComment, err := CommentConv(comment)
		if err != nil {
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
		}
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{StatusCode: 0},
			Comment:  newComment,
		})
	} else if actionType == "2" {
		commentId := c.Query("comment_id")
		tmpId, err := strconv.Atoi(commentId)
		if err != nil {
			log.Errors(c, "comment_id conv int failed", zap.Error(err))
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: "comment_id conv int failed"},
			})
			return
		}
		comment, err := model.DeleteComment(uint(tmpId))
		if err != nil {
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
		}
		newComment, err := CommentConv(comment)
		if err != nil {
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
		}
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{StatusCode: 0},
			Comment:  newComment,
		})
	}
}

func GetCommentList(c *gin.Context) {
	videoId := c.Query("vedio_id")
	login, _ := c.Get("Login")
	if login == false {
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "get comment need login first",
			},
		})
		return
	} else {
		tmp, _ := c.Get("ID")
		var ok bool
		if _, ok = tmp.(uint); !ok {
			log.Errors(c, "id to int error")
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: "id to int error"},
			})
			return
		}
	}
	commentList, err := model.GetComment(videoId)
	if err != nil {
		log.Errors(c, "Get comment_list error", zap.Error(err))
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Get comment_list error"},
		})
	}
	var comments []Comment
	for _, comment := range commentList {
		newComment, err := CommentConv(comment)
		if err != nil {
			c.JSON(http.StatusOK, CommentListResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
		comments = append(comments, newComment)
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{StatusCode: 0},
		Comment:  comments,
	})
}
