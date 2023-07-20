package httpUser

import (
	"advertisingService/internal/users"
	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type updateUserRequest struct {
	ID       string `json:"user_id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type userResponse struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

func UserSuccessResponse(user *users.User) *gin.H {
	return &gin.H{
		"data": userResponse{
			ID:       user.ID,
			Nickname: user.Nickname,
			Email:    user.Email,
		},
		"error": nil,
	}
}

func UserErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
