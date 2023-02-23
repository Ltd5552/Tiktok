package handlers

import (
	"Tiktok/model"
	"Tiktok/pkg/jwt"
	"Tiktok/pkg/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommentResponse struct {
	Response
	Comment Comment `json:"comment"`
}

type CommentListResponse struct {
	Response
	Comment []Comment `json:"comment_list"`
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
	token := c.Query("token")
	tmpVideoId := c.Query("video_id")
	actionType := c.Query("action_type")

	userID, err := jwt.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "please login first"}})
		return
	}

	if actionType == "1" {
		commentText := c.Query("comment_text")
		videoID, err := strconv.Atoi(tmpVideoId)
		if err != nil {
			log.Errors(c, "video_id conv int failed", zap.Error(err))
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "video_id conv int failed"}})
			return
		}
		comment, err := model.CreateComment(userID, videoID, commentText)
		if err != nil {
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
		newComment, err := CommentConv(comment)
		if err != nil {
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
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
			return
		}
		newComment, err := CommentConv(comment)
		if err != nil {
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{StatusCode: 0},
			Comment:  newComment,
		})
	}
}

func GetCommentList(c *gin.Context) {
	videoId := c.Query("video_id")

	commentList, err := model.GetComment(videoId)
	if err != nil {
		log.Errors(c, "Get comment_list error", zap.Error(err))
		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "Get comment_list error"},
		})
		return
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
