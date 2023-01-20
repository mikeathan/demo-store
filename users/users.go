package users

import (
	"errors"
)

type User struct {
	UserName     string
	HashPassword string
}

func CreateUser(username string, password string) *User {
	hashPwd, _ := HashPassword(password)
	return &User{UserName: username, HashPassword: hashPwd}
}

var ErrorUserNotFound = errors.New("User not found")
var ErrorUserExists = errors.New("User exists")
var ErrorUserAuthentication = errors.New("User authentication failed. Username or Password is invalid")

type UserDatabase interface {
	AddUser(username string, password string) error
	Authenticate(username string, password string) error
	IsAdmin(username string) bool
}

type UserStorage struct {
	data map[string]*User
}

func (u *UserStorage) IsAdmin(username string) bool {
	user, err := u.FindUser(username)
	if err != nil {
		return false
	}

	return user.UserName == "admin"
}
func (user *UserStorage) FindUser(username string) (*User, error) {

	if user, ok := user.data[username]; ok {
		return user, nil
	}

	return nil, ErrorUserNotFound
}

func (user *UserStorage) AddUser(username string, password string) error {

	if _, ok := user.data[username]; ok {
		return ErrorUserExists
	}

	user.data[username] = CreateUser(username, password)
	return nil
}

func (u *UserStorage) Authenticate(username string, password string) error {

	user, err := u.FindUser(username)
	if err != nil {
		return err
	}
	if ok := PasswordHashMatches(password, user.HashPassword); !ok {
		return ErrorUserAuthentication
	}

	return nil
}
