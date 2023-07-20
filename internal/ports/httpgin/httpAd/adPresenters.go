package httpAd

import (
	"advertisingService/internal/ads"
	"github.com/gin-gonic/gin"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type AdResponse struct {
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Text         string `json:"text"`
	AuthorID     int64  `json:"author_id"`
	Published    bool   `json:"published"`
	CreationTime string `json:"creationTime"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type deleteAdRequest struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`
}

func AdListSuccessResponse(adList *[]ads.Ad) *gin.H {
	var adsResponse []AdResponse
	for _, ad := range *adList {
		adsResponse = append(adsResponse, AdResponse{
			ID:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorID:     ad.AuthorID,
			Published:    ad.Published,
			CreationTime: ad.CreationDate,
		})
	}
	return &gin.H{
		"data":  adsResponse,
		"error": nil,
	}
}

func AdSuccessResponse(ad *ads.Ad) *gin.H {
	return &gin.H{
		"data": AdResponse{
			ID:           ad.ID,
			Title:        ad.Title,
			Text:         ad.Text,
			AuthorID:     ad.AuthorID,
			Published:    ad.Published,
			CreationTime: ad.CreationDate,
		},
		"error": nil,
	}
}

func AdErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
