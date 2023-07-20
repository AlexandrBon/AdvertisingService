package adrepo

import (
	"advertisingService/internal/ads"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AdRepoSuiteTest struct {
	suite.Suite
	adRepo Repository
}

func (suite *AdRepoSuiteTest) SetupTest() {
	suite.adRepo = New()
}

func (suite *AdRepoSuiteTest) TearDownTest() {
	// close DB
}

func (suite *AdRepoSuiteTest) TestSimple() {
	ad := ads.Ad{ID: 0, Title: "cat", Text: "new", AuthorID: 0}
	adID := suite.adRepo.AddAd(ad)

	ad0, ok := suite.adRepo.GetAdByID(adID)
	suite.True(ok)
	suite.Equal(ad.Title, ad0.Title)
	suite.Equal(ad.Text, ad0.Text)
	suite.Equal(ad.Published, ad0.Published)

	newAd := suite.adRepo.Update(adID, ads.Ad{ID: 0, Title: "dog", Text: "old", AuthorID: 0})
	suite.Equal("dog", newAd.Title)
	suite.Equal("old", newAd.Text)
	suite.LessOrEqual(ad0.LastUpdate, newAd.LastUpdate)
	suite.Equal(ad0.CreationDate, newAd.CreationDate)

	ads0 := suite.adRepo.GetFilteredAds([]func(ads.Ad) bool{
		func(ad ads.Ad) bool {
			return ad.Published == false
		},
	})
	suite.Len(ads0, 1)
	suite.Equal("dog", ads0[0].Title)
	suite.Equal("old", ads0[0].Text)

	deleted := suite.adRepo.DeleteAd(0)
	suite.True(deleted)
}

func TestAdRepoSuite(t *testing.T) {
	suite.Run(t, new(AdRepoSuiteTest))
}
