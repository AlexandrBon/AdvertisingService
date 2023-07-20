package httpAd

import (
	"advertisingService/internal/adApp"
	"github.com/gin-gonic/gin"
)

func AppRouter(r *gin.RouterGroup, a adApp.App) {
	r.POST("/ads", createAd(a))                    // Метод для создания объявления (httpAd)
	r.PUT("/ads/:ad_id/status", changeAdStatus(a)) // Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
	r.PUT("/ads/:ad_id", updateAd(a))              // Метод для обновления текста(Text) или заголовка(Title) объявления
	r.GET("/ads/:ad_id", getAdByID(a))             // Метод для получения объявления по ID
	r.GET("/ads", listAds(a))                      // Метод для получения списка объявлений по фильтрам (по умолчанию выводятся объявления с Published == true)
	r.DELETE("/ads/:ad_id", deleteAd(a))           // Метод для удаления объявления по id
}
