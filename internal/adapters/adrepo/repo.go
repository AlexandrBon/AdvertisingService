package adrepo

import (
	"advertisingService/internal/ads"
	"sync"
	"time"
)

const layout = "01/02/06"

type Repository interface {
	AddAd(ad ads.Ad) int64
	GetAdByID(adID int64) (ads.Ad, bool)
	Update(adID int64, ad ads.Ad) ads.Ad
	GetFilteredAds(filter []func(ad ads.Ad) bool) []ads.Ad
	DeleteAd(adID int64) bool
}

type AdRepository struct {
	dataBase map[int64]ads.Ad
	mtx      sync.Mutex
}

func New() Repository {
	return &AdRepository{dataBase: map[int64]ads.Ad{}}
}

func (ap *AdRepository) AddAd(ad ads.Ad) int64 {
	ap.mtx.Lock()
	defer ap.mtx.Unlock()
	adID := int64(len(ap.dataBase))
	ad.CreationDate = time.Now().UTC().Format(layout)
	ad.LastUpdate = time.Now().UTC().Format(layout)
	ap.dataBase[adID] = ad
	return adID
}

func (ap *AdRepository) GetAdByID(adID int64) (ads.Ad, bool) {
	ap.mtx.Lock()
	defer ap.mtx.Unlock()
	ad, ok := ap.dataBase[adID]
	return ad, ok
}

func (ap *AdRepository) Update(adID int64, ad ads.Ad) ads.Ad {
	ap.mtx.Lock()
	defer ap.mtx.Unlock()
	if prevAdState, ok := ap.dataBase[adID]; ok {
		ad.CreationDate = prevAdState.CreationDate
		ad.LastUpdate = time.Now().UTC().Format(layout)
		ap.dataBase[adID] = ad
		return ad
	}
	return ads.Ad{}
}

func checkAd(filter []func(ad ads.Ad) bool, ad ads.Ad) bool {
	for _, f := range filter {
		if !f(ad) {
			return false
		}
	}
	return true
}

func (ap *AdRepository) GetFilteredAds(filter []func(ad ads.Ad) bool) []ads.Ad {
	ap.mtx.Lock()
	defer ap.mtx.Unlock()
	var filteredAds []ads.Ad
	for _, ad := range ap.dataBase {
		if checkAd(filter, ad) {
			filteredAds = append(filteredAds, ad)
		}
	}
	return filteredAds
}

func (ap *AdRepository) DeleteAd(adID int64) bool {
	_, ok := ap.GetAdByID(adID)
	if !ok {
		return false
	}
	delete(ap.dataBase, adID)
	return true
}
