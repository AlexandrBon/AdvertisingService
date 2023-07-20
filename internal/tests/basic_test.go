package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("testUser", "testUser@mail.ru")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(0))
	assert.False(t, response.Data.Published)
}

func TestChangeAdStatus(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("testUser", "testUser@mail.ru")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response, err = client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)
	assert.True(t, response.Data.Published)

	response, err = client.changeAdStatus(0, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)

	response, err = client.changeAdStatus(0, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("testUser", "testUser@mail.ru")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response, err = client.updateAd(0, response.Data.ID, "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
}

func TestGetUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("testUser", "testUser@mail.ru")
	assert.NoError(t, err)

	response, err := client.getUser(0)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), response.Data.ID)
	assert.Equal(t, "testUser", response.Data.Nickname)
	assert.Equal(t, "testUser@mail.ru", response.Data.Email)
}

func TestDeleteUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("testUser", "testUser@mail.ru")
	assert.NoError(t, err)

	err = client.deleteUser(0)
	assert.NoError(t, err)

	_, err = client.getUser(0)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestListAds(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("testUser", "testUser@mail.ru")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(0, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAds(map[string]string{})
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}

func TestGetAdByID(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("testUser", "testUser@mail.ru")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(0, "best cat", "not for sale")
	assert.NoError(t, err)

	response, err = client.getAdByID(0)
	assert.NoError(t, err)
	assert.Equal(t, publishedAd.Data.Title, response.Data.Title)
	assert.Equal(t, publishedAd.Data.Text, response.Data.Text)
	assert.True(t, response.Data.Published)
	assert.Zero(t, response.Data.AuthorID)
}

func TestGetAdByTitle(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("testUser", "testUser@mail.ru")
	assert.NoError(t, err)

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(0, "best cat", "not for sale")
	assert.NoError(t, err)

	responseList, err := client.getAdsByTitle("hello")
	assert.NoError(t, err)
	assert.Equal(t, publishedAd.Data.Title, responseList.Data[0].Title)
	assert.Equal(t, publishedAd.Data.Text, responseList.Data[0].Text)
	assert.True(t, responseList.Data[0].Published)
	assert.Zero(t, responseList.Data[0].AuthorID)
}
