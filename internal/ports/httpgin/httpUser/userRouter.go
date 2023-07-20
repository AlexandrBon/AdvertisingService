package httpUser

import (
	"advertisingService/internal/userApp"
	"github.com/gin-gonic/gin"
)

func AppRouter(r *gin.RouterGroup, ua userApp.App) {
	r.POST("/users", createUser(ua))                // Метод для создания пользователя
	r.PUT("/users/:user_id/status", updateUser(ua)) // Метод для изменения данных пользователя
	r.DELETE("/users/:user_id", deleteUser(ua))     // Метод для удаления пользователя по id
	r.GET("/users/:user_id", getUser(ua))           // Метод для получения пользователя по id
}
