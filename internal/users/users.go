package users

type User struct {
	ID       int64
	Nickname string `validate:"minmax:1,20"`
	Email    string `validate:"minmax:3,50"`
}
