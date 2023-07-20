package userrepo

import (
	"advertisingService/internal/users"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserRepoSuiteTest struct {
	suite.Suite
	userRepo Repository
}

func (suite *UserRepoSuiteTest) SetupTest() {
	suite.userRepo = New()
}

func (suite *UserRepoSuiteTest) TearDownTest() {
	// close DB
}

func (suite *UserRepoSuiteTest) TestSimple() {
	user := users.User{ID: 0, Nickname: "a", Email: "a@mail.ru"}
	userID := suite.userRepo.AddUser(user)

	user0, ok := suite.userRepo.GetUserByID(userID)
	suite.True(ok)
	suite.Equal(user0, user)

	suite.userRepo.Update(0, users.User{ID: 0, Nickname: "b", Email: "b@mail.ru"})
	user0, ok = suite.userRepo.GetUserByID(0)
	suite.True(ok)
	suite.Equal("b", user0.Nickname)
	suite.Equal("b@mail.ru", user0.Email)

	suite.True(suite.userRepo.CheckIfEmailAlreadyInUse(user0.Email))

	suite.True(suite.userRepo.DeleteUser(0))
	suite.False(suite.userRepo.DeleteUser(0))
}

func TestUserRepoSuite(t *testing.T) {
	suite.Run(t, new(UserRepoSuiteTest))
}
