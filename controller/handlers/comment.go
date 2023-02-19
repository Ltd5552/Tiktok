package handlers

import (
	"Tiktok/model"
	"Tiktok/pkg/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommentResponse struct{
	Response
	Comment Comment
}

func CommentConv(comment model.Comment) (Comment, error){
	var new_comment Comment
	new_comment.Id = comment.Model.ID
	tmp_user, err := model.ReadUser(strconv.Itoa(int(comment.Commenter_id)))
	if err != nil {
		return Comment{}, err
	}
	new_comment.User = User{
		Id: tmp_user.ID,
		Name: tmp_user.Name,
	}
	new_comment.Content = comment.Text
	new_comment.CreateDate = comment.CreatedAt.Format("01-02")
	return new_comment, nil
}

func CommentAction(c *gin.Context) {
	login, _ := c.Get("Login")
	tmp_video_id := c.Query("video_id")
	action_type := c.Query("action_type")
	var id uint
	if login == false {
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg: "comment need login first",
			},
			Comment: Comment{},
		})
		return
	} else {
		tmp, _ := c.Get("ID")
		var ok bool
		if id, ok = tmp.(uint); !ok{
			log.Error("id to int error")
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: "id to int error"},
			})
			return
		} 
	}
	if action_type == "1" {
		comment_text := c.Query("comment_text")
		video_id, err := strconv.Atoi(tmp_video_id)
		if err != nil{
			log.Error("video_id conv int failed", zap.Error(err))
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: "video_id conv int failed"},
			})
			return
		}
		comment, err := model.CreateComment(id, uint(video_id), comment_text)
		if err !=nil {
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
		}
		new_comment, err := CommentConv(comment)
		if err !=nil {
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
		}
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{StatusCode: 0},
			Comment: new_comment,
		})
	} else if action_type =="2" {
		comment_id := c.Query("comment_id")
		tmp_id, err := strconv.Atoi(comment_id)
		if err != nil{
			log.Error("comment_id conv int failed", zap.Error(err))
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: "comment_id conv int failed"},
			})
			return
		}
		comment, err := model.DeleteComment(uint(tmp_id))
		if err !=nil {
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
		}
		new_comment, err := CommentConv(comment)
		if err !=nil {
			c.JSON(http.StatusOK, CommentResponse{
				Response: Response{StatusCode: 1, StatusMsg: err.Error()},
			})
		}
		c.JSON(http.StatusOK, CommentResponse{
			Response: Response{StatusCode: 0},
			Comment: new_comment,
		})
	}
}

func GetCommentList(c *gin.Context) {

}
