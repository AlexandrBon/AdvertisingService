package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("testUser", "testUser@mail.ru")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), response.Data.ID)
	assert.Equal(t, "testUser", response.Data.Nickname)
	assert.Equal(t, "testUser@mail.ru", response.Data.Email)
}

func TestUpdateUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("testUser", "testUser@mail.ru")
	assert.NoError(t, err)

	response, err := client.updateUser(0, "cat", "cat@mail.ru")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), response.Data.ID)
	assert.Equal(t, "cat", response.Data.Nickname)
	assert.Equal(t, "cat@mail.ru", response.Data.Email)
}

func TestChangeStatusAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("testUser", "testUser@mail.ru")
	assert.NoError(t, err)

	_, err = client.createUser("testUser", "testUser0@mail.ru")
	assert.NoError(t, err)

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(1, resp.Data.ID, true)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestUpdateAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("testUser", "testUser@mail.ru")
	assert.NoError(t, err)

	_, err = client.createUser("testUser", "testUser0@mail.ru")
	assert.NoError(t, err)

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(1, resp.Data.ID, "title", "text")
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestCreateAd_ID(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("testUser", "testUser@mail.ru")
	assert.NoError(t, err)

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(0))

	resp, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(1))

	resp, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(2))
}

func TestCreateAdWithoutUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createAd(0, "hello", "world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateUserWithSameEmail(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("testUser", "testUser@mail.ru")
	assert.NoError(t, err)

	_, err = client.createUser("testUser", "testUser@mail.ru")
	assert.ErrorIs(t, err, ErrEmailConflict)
}

func TestListAdsFilterByAuthor(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("testUser", "testUser@mail.ru")
	assert.NoError(t, err)

	_, err = client.createUser("testUser", "testUser0@mail.ru")
	assert.NoError(t, err)

	response0, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.createAd(1, "my", "boots")
	assert.NoError(t, err)

	ads, err := client.listAds(map[string]string{
		"authorID": "0",
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, len(ads.Data))
	assert.Equal(t, response0.Data.AuthorID, ads.Data[0].AuthorID)
}

func TestListAdsFilterByTitle(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("testUser", "testUser@mail.ru")
	assert.NoError(t, err)

	response0, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.createAd(0, "my", "boots")
	assert.NoError(t, err)

	ads, err := client.listAds(map[string]string{
		"title": "hel",
	})

	assert.NoError(t, err)
	assert.Equal(t, 1, len(ads.Data))
	assert.Equal(t, response0.Data.Title, ads.Data[0].Title)
}
