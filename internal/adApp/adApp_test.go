package adApp

import (
	"advertisingService/internal/adapters/adrepo"
	"advertisingService/internal/adapters/userrepo"
	"advertisingService/internal/ads"
	"advertisingService/internal/userApp"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AdAppSuiteTest struct {
	suite.Suite
	adRepo      adrepo.Repository
	userRepo    userrepo.Repository
	adAppTest   App
	userAppTest userApp.App
}

func (suite *AdAppSuiteTest) SetupTest() {
	suite.adRepo = adrepo.New()
	suite.userRepo = userrepo.New()
	suite.adAppTest = NewApp(suite.adRepo, suite.userRepo)
	suite.userAppTest = userApp.NewApp(suite.userRepo)
}

func (suite *AdAppSuiteTest) TearDownTest() {
	// close DB
}

func (suite *AdAppSuiteTest) TestAddAd0() {
	_, err := suite.userAppTest.CreateUser("a", "a@mail.ru")
	suite.Nil(err)
	ad, err := suite.adAppTest.CreateAd("dog", "new", 0)
	suite.Nil(err)
	suite.Equal(int64(0), ad.ID)
	suite.Equal(false, ad.Published)
	suite.Equal("dog", ad.Title)
	suite.Equal("new", ad.Text)
	suite.Equal(int64(0), ad.AuthorID)
}

func (suite *AdAppSuiteTest) TestAddAd1() {
	_, err := suite.adAppTest.CreateAd("dog", "new", 0)
	suite.NotNil(err)
	suite.ErrorIs(err, ErrBadRequest)
}

func (suite *AdAppSuiteTest) TestAddAd2() {
	_, err := suite.userAppTest.CreateUser("a", "a@mail.ru")
	suite.Nil(err)

	_, err = suite.adAppTest.CreateAd("", "", 0)
	suite.NotNil(err)
	suite.ErrorIs(err, ErrBadRequest)
}

func (suite *AdAppSuiteTest) TestUpdateAd0() {
	_, err := suite.userAppTest.CreateUser("a", "a@mail.ru")
	suite.Nil(err)

	_, err = suite.adAppTest.CreateAd("dog", "new", 0)
	suite.Nil(err)

	ad, err := suite.adAppTest.UpdateAd(0, 0, "cat", "old")
	suite.Nil(err)
	suite.Equal(int64(0), ad.ID)
	suite.Equal(false, ad.Published)
	suite.Equal("cat", ad.Title)
	suite.Equal("old", ad.Text)
	suite.Equal(int64(0), ad.AuthorID)
}

func (suite *AdAppSuiteTest) TestUpdateAd1() {
	_, err := suite.userAppTest.CreateUser("a", "a@mail.ru")
	suite.Nil(err)

	_, err = suite.adAppTest.CreateAd("dog", "new", 0)
	suite.Nil(err)

	_, err = suite.adAppTest.UpdateAd(0, 0, "", "old")
	suite.NotNil(err)
	suite.ErrorIs(err, ErrBadRequest)
}

func (suite *AdAppSuiteTest) TestChangeAdStatus0() {
	_, err := suite.userAppTest.CreateUser("a", "a@mail.ru")
	suite.Nil(err)

	_, err = suite.adAppTest.CreateAd("dog", "new", 0)
	suite.Nil(err)

	ad, err := suite.adAppTest.ChangeAdStatus(0, 0, true)
	suite.Nil(err)
	suite.Equal(true, ad.Published)
}

func (suite *AdAppSuiteTest) TestChangeAdStatus1() {
	_, err := suite.userAppTest.CreateUser("a", "a@mail.ru")
	suite.Nil(err)

	_, err = suite.userAppTest.CreateUser("b", "b@mail.ru")
	suite.Nil(err)

	_, err = suite.adAppTest.CreateAd("dog", "new", 0)
	suite.Nil(err)

	_, err = suite.adAppTest.ChangeAdStatus(0, 1, true)
	suite.NotNil(err)
	suite.ErrorIs(err, ErrForbidden)
}

func (suite *AdAppSuiteTest) TestListAds() {
	_, err := suite.userAppTest.CreateUser("a", "a@mail.ru")
	suite.Nil(err)

	_, err = suite.adAppTest.CreateAd("dog", "new", 0)
	suite.Nil(err)

	var filter []func(ads.Ad) bool
	filter = append(filter, func(ad ads.Ad) bool {
		return ad.Published == false
	})
	adList, err := suite.adAppTest.ListAds(filter)
	suite.Nil(err)
	suite.Equal(1, len(adList))
	suite.Equal("dog", adList[0].Title)
}

func (suite *AdAppSuiteTest) TestGetAdByID0() {
	_, err := suite.userAppTest.CreateUser("a", "a@mail.ru")
	suite.Nil(err)

	ad, err := suite.adAppTest.CreateAd("dog", "new", 0)
	suite.Nil(err)

	suite.Nil(err)
	suite.Equal(int64(0), ad.ID)
	suite.Equal(false, ad.Published)
	suite.Equal("dog", ad.Title)
	suite.Equal("new", ad.Text)
	suite.Equal(int64(0), ad.AuthorID)
}

func (suite *AdAppSuiteTest) TestGetAdByID1() {
	_, err := suite.userAppTest.CreateUser("a", "a@mail.ru")
	suite.Nil(err)

	_, err = suite.adAppTest.CreateAd("dog", "new", 0)
	suite.Nil(err)

	_, err = suite.adAppTest.GetAdByID(1)
	suite.NotNil(err)
	suite.ErrorIs(err, ErrBadRequest)
}

func (suite *AdAppSuiteTest) TestDeleteAd0() {
	_, err := suite.userAppTest.CreateUser("a", "a@mail.ru")
	suite.Nil(err)

	_, err = suite.adAppTest.CreateAd("dog", "new", 0)
	suite.Nil(err)

	err = suite.adAppTest.DeleteAd(0, 0)
	suite.Nil(err)
}

func (suite *AdAppSuiteTest) TestDeleteAd1() {
	_, err := suite.userAppTest.CreateUser("a", "a@mail.ru")
	suite.Nil(err)

	_, err = suite.userAppTest.CreateUser("b", "b@mail.ru")
	suite.Nil(err)

	_, err = suite.adAppTest.CreateAd("dog", "new", 0)
	suite.Nil(err)

	err = suite.adAppTest.DeleteAd(0, 1)
	suite.NotNil(err)
	suite.ErrorIs(err, ErrForbidden)
}

func TestAdAppSuite(t *testing.T) {
	suite.Run(t, new(AdAppSuiteTest))
}
