package httpAd

import (
	"advertisingService/internal/adApp"
	"advertisingService/internal/ads"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func getHTTPStatus(err error) int {
	if errors.Is(err, adApp.ErrForbidden) {
		return http.StatusForbidden
	} else if errors.Is(err, adApp.ErrBadRequest) {
		return http.StatusBadRequest
	}
	return -1
}

// Метод для создания объявления (httpAd)
func createAd(a adApp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.BindJSON(&reqBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		response, err := a.CreateAd(reqBody.Title, reqBody.Text, reqBody.UserID)

		if err != nil {
			c.JSON(getHTTPStatus(err), AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&response))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a adApp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adID, err := strconv.ParseInt(c.Param("ad_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(adApp.ErrBadRequest))
			return
		}

		response, err := a.ChangeAdStatus(adID, reqBody.UserID, reqBody.Published)

		if err != nil {
			c.JSON(getHTTPStatus(err), AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&response))
	}
}

// Метод для получения объявления по его id
func getAdByID(a adApp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		adID, err := strconv.ParseInt(c.Param("ad_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		response, err := a.GetAdByID(adID)

		if err != nil {
			c.JSON(getHTTPStatus(err), AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&response))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a adApp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adID, err := strconv.ParseInt(c.Param("ad_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		response, err := a.UpdateAd(adID, reqBody.UserID, reqBody.Title, reqBody.Text)

		if err != nil {
			c.JSON(getHTTPStatus(err), AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&response))
	}
}

// Метод для получения списка объявлений с фильтрами(по умолчанию фильтр published=true)
func listAds(a adApp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter []func(ads.Ad) bool

		userID, err := strconv.ParseInt(c.Query("authorID"), 10, 64)
		if err == nil {
			filter = append(filter, func(ad ads.Ad) bool {
				return ad.AuthorID == userID
			})
		}

		published, err := strconv.ParseBool(c.Query("published"))
		if err == nil {
			filter = append(filter, func(ad ads.Ad) bool {
				return ad.Published == published
			})
		}

		creationTime := c.Query("creationTime")
		if creationTime != "" {
			filter = append(filter, func(ad ads.Ad) bool {
				return ad.CreationDate >= creationTime
			})
		}

		title := c.Query("title")
		if title != "" {
			filter = append(filter, func(ad ads.Ad) bool {
				return strings.Contains(ad.Title, title)
			})
		}

		response, _ := a.ListAds(filter)

		c.JSON(http.StatusOK, AdListSuccessResponse(&response))
	}
}

// Метод для удаления объявления по id
func deleteAd(a adApp.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody deleteAdRequest
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adID, err := strconv.ParseInt(c.Param("ad_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		err = a.DeleteAd(adID, reqBody.UserID)

		if err != nil {
			c.JSON(getHTTPStatus(err), AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(&ads.Ad{}))
	}
}
