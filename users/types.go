package users

import (
	"demo-store/common"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

func CreateUserDatabase() UserDatabase {

	userStorage := UserStorage{}
	userStorage.data = make(map[string]*User)

	return &userStorage
}

func Load(path string) UserDatabase {
	usersDatabase, err := LoadFromCache(path)
	if err != nil {
		usersDatabase = CreateUserDatabase()
	}

	return usersDatabase
}

func LoadFromCache(path string) (UserDatabase, error) {

	if !common.DirExists(path) {
		return nil, os.ErrNotExist
	}

	file, err := os.Open(filepath.Join(path, "users.dat"))
	defer file.Close()
	if err != nil {
		//log.Fatal(fmt.Printf("Error opening user.dat file %s", err))
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		//log.Fatal(fmt.Printf("Error reading user.dat file %s", err))
		return nil, err
	}

	users, serr := ToObject(string(data))
	if serr != nil {
		//log.Fatal(fmt.Printf("Error deserialising user.dat file %s", serr))
		return nil, serr
	}

	userStorage := UserStorage{}
	userStorage.data = make(map[string]*User)

	for _, user := range users {
		userStorage.data[user.UserName] = &User{UserName: user.UserName, HashPassword: user.HashPassword}
	}

	return &userStorage, nil
}

func SaveToCache(path string, user *UserStorage) error {
	common.CreateDirIfNotExists(path)

	file, err := os.Create(filepath.Join(path, "users.dat"))
	defer file.Close()

	if err != nil {
		//log.Fatal(fmt.Printf("Error initialising user.dat file %s", err))
		return err
	}

	values := make([]*User, 0, len(user.data))
	for k := range user.data {
		values = append(values, user.data[k])
	}

	json, _ := common.ToJson(values)
	file.Write([]byte(json))
	return nil
}

func ToObject(jsonString string) ([]User, error) {
	var users []User
	err := json.Unmarshal([]byte(jsonString), &users)

	if err != nil {
		return nil, err
	}
	return users, nil
}
