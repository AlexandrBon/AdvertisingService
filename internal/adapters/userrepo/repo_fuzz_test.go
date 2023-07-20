package userrepo

import (
	"advertisingService/internal/users"
	"testing"
)

func FuzzUserRepository_DeleteUser(f *testing.F) {
	userRepo := New()

	f.Fuzz(func(t *testing.T, userID int64) {
		ID := userRepo.AddUser(users.User{ID: userID})
		got := userRepo.DeleteUser(ID)
		if got != true {
			t.Errorf("For (%d) Expect: true, but got: %t", userID, got)
		}
	})
}
