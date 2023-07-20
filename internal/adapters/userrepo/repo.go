package userrepo

import (
	"advertisingService/internal/users"
	"sync"
)

type Repository interface {
	AddUser(user users.User) int64
	Update(userID int64, user users.User)
	GetUserByID(userID int64) (users.User, bool)
	CheckIfEmailAlreadyInUse(email string) bool
	DeleteUser(userID int64) bool
}

type UserRepository struct {
	dataBase map[int64]users.User
	mtx      sync.Mutex
}

func New() Repository {
	return &UserRepository{dataBase: map[int64]users.User{}}
}

func (up *UserRepository) CheckIfEmailAlreadyInUse(email string) bool {
	up.mtx.Lock()
	defer up.mtx.Unlock()
	for _, user := range up.dataBase {
		if user.Email == email {
			return true
		}
	}
	return false
}

func (up *UserRepository) AddUser(user users.User) int64 {
	up.mtx.Lock()
	defer up.mtx.Unlock()
	userID := int64(len(up.dataBase))
	up.dataBase[userID] = user
	return userID
}

func (up *UserRepository) Update(userID int64, user users.User) {
	up.mtx.Lock()
	defer up.mtx.Unlock()
	if _, ok := up.dataBase[userID]; ok {
		up.dataBase[userID] = user
	}
}

func (up *UserRepository) GetUserByID(userID int64) (users.User, bool) {
	up.mtx.Lock()
	defer up.mtx.Unlock()
	ad, ok := up.dataBase[userID]
	return ad, ok
}

func (up *UserRepository) DeleteUser(userID int64) bool {
	up.mtx.Lock()
	defer up.mtx.Unlock()
	_, ok := up.dataBase[userID]
	if !ok {
		return false
	}
	delete(up.dataBase, userID)
	return true
}
