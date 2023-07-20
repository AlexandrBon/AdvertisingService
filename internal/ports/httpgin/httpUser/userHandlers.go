package httpUser

import (
	"advertisingService/internal/userApp"
	"advertisingService/internal/users"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func getHTTPStatus(err error) int {
	if errors.Is(err, userApp.ErrBadRequest) {
		return http.StatusBadRequest
	} else if errors.Is(err, userApp.ErrEmailConflict) {
		return http.StatusConflict
	}
	return -1
}

func createUser(ua userApp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createUserRequest
		err := c.BindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		response, err := ua.CreateUser(reqBody.Nickname, reqBody.Email)
		if err != nil {
			c.JSON(getHTTPStatus(err), UserErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(&response))
	}
}

func updateUser(ua userApp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateUserRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		response, err := ua.UpdateUser(userID, reqBody.Nickname, reqBody.Email)

		if err != nil {
			c.JSON(getHTTPStatus(err), UserErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(&response))
	}
}

func deleteUser(ua userApp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		err = ua.DeleteUser(userID)

		if err != nil {
			c.JSON(getHTTPStatus(err), UserErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(&users.User{}))
	}
}

func getUser(ua userApp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, UserErrorResponse(err))
			return
		}

		response, err := ua.GetUser(userID)

		if err != nil {
			c.JSON(getHTTPStatus(err), UserErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, UserSuccessResponse(&response))
	}
}
