package users_test

import (
	"demo-store/common"
	"demo-store/users"
	"fmt"
	"testing"
)

func TestAuthFailsWithWrongPassword(t *testing.T) {

	storage := users.CreateUserDatabase()
	storage.AddUser("user1", "11111")

	err := storage.Authenticate("user1", "22222")

	expected := users.ErrorUserAuthentication
	if err != expected {
		t.Errorf("Unexpected error: got %v want %v,", err, expected)
	}
}

func TestAuthFailsWithWrongUsername(t *testing.T) {

	storage := users.CreateUserDatabase()
	storage.AddUser("user1", "11111")

	err := storage.Authenticate("user2", "11111")

	expected := users.ErrorUserNotFound
	if err != expected {
		t.Errorf("Unexpected error: got %v want %v,", err, expected)
	}
}

func TestAuthSucceeds(t *testing.T) {

	storage := users.CreateUserDatabase()
	storage.AddUser("user1", "11111")

	err := storage.Authenticate("user1", "11111")

	var expected error = nil
	if err != expected {
		t.Errorf("Unexpected error: got %v want %v,", err, expected)
	}
}

func TestAuthFailsIfAddingExistingUser(t *testing.T) {

	storage := users.CreateUserDatabase()
	storage.AddUser("user1", "11111")

	err := storage.AddUser("user1", "11111")

	expected := users.ErrorUserExists
	if err != err {
		t.Errorf("Unexpected error: got %v want %v,", err, expected)
	}
}

func TestGenerateUsers(t *testing.T) {

	t.Skip("component test")
	storage := users.CreateUserDatabase()
	storage.AddUser("user_a", "passwordA")
	storage.AddUser("user_b", "passwordB")
	storage.AddUser("user_c", "passwordC")
	storage.AddUser("admin", "Password1")
	st := storage.(*users.UserStorage)
	users.SaveToCache("../cache", st)

	storage2, err := users.LoadFromCache("../cache")
	if err != nil {
		t.Error(err)
	}

	fmt.Println(storage2)
}

func TestIsAdmin(t *testing.T) {
	storage := users.CreateUserDatabase()
	storage.AddUser("user_a", "passwordA")
	storage.AddUser("user_b", "passwordB")
	storage.AddUser("user_c", "passwordC")
	storage.AddUser("admin", "Password1")

	if ok := storage.IsAdmin("admin"); !ok {
		t.Errorf("Unexpected error: user admin: got %v want %v,", ok, "true")
	}

	if ok := storage.IsAdmin("user_a"); ok {
		t.Errorf("Unexpected error: user user_a: got %v want %v,", ok, "false")
	}

	if ok := storage.IsAdmin("new_user"); ok {
		t.Errorf("Unexpected error: user new_user: got %v want %v,", ok, "true")
	}
}

func TestToObject(t *testing.T) {

	testUsers := []users.User{
		*users.CreateUser("user1", "111"),
		*users.CreateUser("user2", "222"),
	}
	jsonStr, _ := common.ToJson(testUsers)
	deserializedUsers, _ := users.ToObject(jsonStr)
	for i, user := range deserializedUsers {
		if user.UserName != testUsers[i].UserName {
			t.Errorf("Unexpected user : got %v want %v,", user.UserName, testUsers[i].UserName)
		}
	}
}

func TestUserDatabaseLoadWhenUserCacheDontExists(t *testing.T) {

	db := users.Load("somedir")
	if db == nil {
		t.Errorf("Unexpected error : got %v want %v,", "nil", "init db")
	}
}

func TestUserDatabaseLoadWhenUserCacheExists(t *testing.T) {

	db := users.Load("../cache")
	if db == nil {
		t.Errorf("Unexpected error : got %v want %v,", db, "init db")
	}

	if err := db.Authenticate("user_a", "passwordA"); err != nil {
		t.Errorf("Unexpected error : got %v want %v,", err, "nil")
	}
}
