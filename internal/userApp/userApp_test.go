package userApp

import (
	"advertisingService/internal/adapters/userrepo"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
)

type UserAppSuiteTest struct {
	suite.Suite
	userAppTest App
}

func (suite *UserAppSuiteTest) SetupTest() {
	userRepo := userrepo.New()
	suite.userAppTest = NewApp(userRepo)
}

func (suite *UserAppSuiteTest) TearDownTest() {
	// close DB
}

func (suite *UserAppSuiteTest) TestAddUser() {
	user, err := suite.userAppTest.CreateUser("a", "a@mail.ru")
	suite.Nil(err)
	suite.Equal(int64(0), user.ID)

	_, err = suite.userAppTest.CreateUser("b", "a@mail.ru")
	suite.NotNil(err)
	suite.ErrorIs(err, ErrEmailConflict)
}

func (suite *UserAppSuiteTest) TestDeleteUser() {
	user, err := suite.userAppTest.CreateUser("a", "a@mail.ru")
	suite.Nil(err)

	err = suite.userAppTest.DeleteUser(user.ID)
	suite.Nil(err)

	err = suite.userAppTest.DeleteUser(user.ID)
	suite.NotNil(err)
	suite.ErrorIs(err, ErrBadRequest)
}

func (suite *UserAppSuiteTest) TestUpdateUser0() {
	user, err := suite.userAppTest.CreateUser("a", "a@mail.ru")
	suite.Nil(err)

	updateUser, err := suite.userAppTest.UpdateUser(user.ID, "b", "a@mail.ru")
	suite.Nil(err)
	suite.Equal(int64(0), updateUser.ID)
	suite.Equal("b", updateUser.Nickname)
	suite.Equal("a@mail.ru", updateUser.Email)

	err = suite.userAppTest.DeleteUser(user.ID)
	suite.Nil(err)

	_, err = suite.userAppTest.UpdateUser(user.ID, "b", "a@mail.ru")
	suite.NotNil(err)
	suite.ErrorIs(err, ErrBadRequest)
}

func (suite *UserAppSuiteTest) TestUpdateUser1() {
	user, err := suite.userAppTest.CreateUser("a", "a@mail.ru")
	suite.Nil(err)

	_, err = suite.userAppTest.UpdateUser(user.ID, "", "a@mail.ru")
	suite.NotNil(err)
	suite.ErrorIs(err, ErrBadRequest)
}

func (suite *UserAppSuiteTest) TestGetUser() {
	user, err := suite.userAppTest.CreateUser("a", "a@mail.ru")
	suite.Nil(err)

	user, err = suite.userAppTest.GetUser(user.ID)
	suite.Nil(err)
	suite.Equal(int64(0), user.ID)

	err = suite.userAppTest.DeleteUser(user.ID)
	suite.Nil(err)

	_, err = suite.userAppTest.GetUser(user.ID)
	suite.NotNil(err)
	suite.ErrorIs(err, ErrBadRequest)
}

func (suite *UserAppSuiteTest) TestTableTest() {
	type Test struct {
		In     int64
		Expect string
	}

	tests := []Test{
		{0, "a"},
		{1, "b"},
		{2, "c"},
		{3, "d"},
		{4, "e"},
	}

	// setup
	_, err := suite.userAppTest.CreateUser("a", "a@mail.ru")
	suite.Nil(err)
	_, err = suite.userAppTest.CreateUser("b", "b@mail.ru")
	suite.Nil(err)
	_, err = suite.userAppTest.CreateUser("c", "c@mail.ru")
	suite.Nil(err)
	_, err = suite.userAppTest.CreateUser("d", "d@mail.ru")
	suite.Nil(err)
	_, err = suite.userAppTest.CreateUser("e", "e@mail.ru")
	suite.Nil(err)
	//

	for _, test := range tests {
		user, err := suite.userAppTest.GetUser(test.In)
		suite.Nil(err)
		if user.Nickname != test.Expect {
			suite.Failf("test %d: expect %s got %v", strconv.FormatInt(test.In, 10), test.Expect, user.Nickname)
		}
	}
}

func TestUserAppSuite(t *testing.T) {
	suite.Run(t, new(UserAppSuiteTest))
}
