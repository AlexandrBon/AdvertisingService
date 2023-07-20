package userApp

import (
	"advertisingService/internal/adapters/userrepo"
	"advertisingService/internal/users"
	"fmt"
	"github.com/AlexandrBon/Validator"
)

var ErrBadRequest = fmt.Errorf("bad request")
var ErrEmailConflict = fmt.Errorf("email already in use")

type App interface {
	CreateUser(nickname string, email string) (users.User, error)
	UpdateUser(userID int64, nickname string, email string) (users.User, error)
	GetUser(usedID int64) (users.User, error)
	DeleteUser(userID int64) error
}

type UserApp struct {
	repository userrepo.Repository
}

func (ua *UserApp) GetUser(userID int64) (users.User, error) {
	user, ok := ua.repository.GetUserByID(userID)
	if !ok {
		return users.User{}, ErrBadRequest
	}
	return user, nil
}

func (ua *UserApp) CreateUser(nickname string, email string) (users.User, error) {
	if ua.repository.CheckIfEmailAlreadyInUse(email) {
		return users.User{}, ErrEmailConflict
	}

	user := users.User{Nickname: nickname, Email: email}
	err := Validator.Validate(user)
	if err != nil {
		return users.User{}, ErrBadRequest
	}

	id := ua.repository.AddUser(user)
	user.ID = id
	ua.repository.Update(user.ID, user)
	return user, nil
}

func (ua *UserApp) UpdateUser(userID int64, nickname string, email string) (users.User, error) {
	user, ok := ua.repository.GetUserByID(userID)

	if ua.repository.CheckIfEmailAlreadyInUse(email) && !ok {
		return users.User{}, ErrEmailConflict
	}

	if !ok {
		return users.User{}, ErrBadRequest
	}

	tmpUser := users.User{Nickname: nickname, Email: email}
	err := Validator.Validate(tmpUser)
	if err != nil {
		return users.User{}, ErrBadRequest
	}

	user.Nickname = nickname
	user.Email = email
	ua.repository.Update(user.ID, user)
	return user, nil
}

func (ua *UserApp) DeleteUser(userID int64) error {
	ok := ua.repository.DeleteUser(userID)
	if !ok {
		return ErrBadRequest
	}
	return nil
}

func NewApp(repo userrepo.Repository) App {
	return &UserApp{repository: repo}
}
