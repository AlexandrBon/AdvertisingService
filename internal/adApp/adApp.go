package adApp

import (
	"advertisingService/internal/adapters/adrepo"
	"advertisingService/internal/adapters/userrepo"
	"advertisingService/internal/ads"
	"fmt"
	"github.com/AlexandrBon/Validator"
)

var (
	ErrBadRequest = fmt.Errorf("bad request")
	ErrForbidden  = fmt.Errorf("forbidden")
)

type App interface {
	CreateAd(title string, text string, usedID int64) (ads.Ad, error)
	ChangeAdStatus(adID int64, userID int64, published bool) (ads.Ad, error)
	UpdateAd(adID int64, userID int64, Title string, Text string) (ads.Ad, error)
	ListAds(filter []func(ad ads.Ad) bool) ([]ads.Ad, error)
	GetAdByID(adID int64) (ads.Ad, error)
	DeleteAd(adID int64, userID int64) error
}

type AdApp struct {
	repository adrepo.Repository
	userRepo   userrepo.Repository
}

func checkIfUserExists(repo userrepo.Repository, userID int64) bool {
	_, ok := repo.GetUserByID(userID)
	return ok
}

func (aa *AdApp) CreateAd(title string, text string, userID int64) (ads.Ad, error) {
	if !checkIfUserExists(aa.userRepo, userID) {
		return ads.Ad{}, ErrBadRequest
	}

	ad := ads.Ad{Title: title, Text: text, AuthorID: userID, Published: false}
	err := Validator.Validate(ad)
	if err != nil {
		return ads.Ad{}, ErrBadRequest
	}
	id := aa.repository.AddAd(ad)

	ad.ID = id
	ad = aa.repository.Update(ad.ID, ad)

	return ad, nil
}

func (aa *AdApp) ChangeAdStatus(adID int64, userID int64, published bool) (ads.Ad, error) {
	if !checkIfUserExists(aa.userRepo, userID) {
		return ads.Ad{}, ErrBadRequest
	}

	ad, ok := aa.repository.GetAdByID(adID)
	if !ok {
		return ads.Ad{}, ErrBadRequest
	}
	if ad.AuthorID != userID {
		return ads.Ad{}, ErrForbidden
	}

	ad.Published = published
	aa.repository.Update(ad.ID, ad)
	return ad, nil
}

func (aa *AdApp) UpdateAd(adID int64, userID int64, Title string, Text string) (ads.Ad, error) {
	if !checkIfUserExists(aa.userRepo, userID) {
		return ads.Ad{}, ErrBadRequest
	}

	ad, ok := aa.repository.GetAdByID(adID)
	if !ok {
		return ads.Ad{}, ErrBadRequest
	}
	if ad.AuthorID != userID {
		return ads.Ad{}, ErrForbidden
	}

	tmpAd := ads.Ad{Title: Title, Text: Text}
	err := Validator.Validate(tmpAd)
	if err != nil {
		return ads.Ad{}, ErrBadRequest
	}

	ad.Title = Title
	ad.Text = Text
	aa.repository.Update(ad.ID, ad)
	return ad, nil
}

func (aa *AdApp) ListAds(filter []func(ad ads.Ad) bool) ([]ads.Ad, error) {
	if len(filter) == 0 {
		filter = append(filter, func(ad ads.Ad) bool {
			return ad.Published
		})
	}
	return aa.repository.GetFilteredAds(filter), nil
}

func (aa *AdApp) GetAdByID(adID int64) (ads.Ad, error) {
	ad, ok := aa.repository.GetAdByID(adID)
	if !ok {
		return ads.Ad{}, ErrBadRequest
	}
	return ad, nil
}

func (aa *AdApp) DeleteAd(adID int64, userID int64) error {
	if !checkIfUserExists(aa.userRepo, userID) {
		return ErrBadRequest
	}

	ad, ok := aa.repository.GetAdByID(adID)
	if !ok {
		return ErrBadRequest
	}

	if ad.AuthorID != userID {
		return ErrForbidden
	}

	_ = aa.repository.DeleteAd(adID)

	return nil
}

func NewApp(repo adrepo.Repository, userRepo userrepo.Repository) App {
	return &AdApp{repository: repo, userRepo: userRepo}
}
