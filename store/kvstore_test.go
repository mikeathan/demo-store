package store_test

import (
	"demo-store/common"
	"demo-store/store"
	"testing"
)

func TestUsers(t *testing.T) {

	mockStore := NewMockStore()

	mockStore.UserDatabase().AddUser("testUser1", "1111")

	err := mockStore.UserDatabase().Authenticate("testUser1", "1111")

	if err != nil {
		t.Errorf("Second Put unexpected value got %s want %s", err, "nil")
	}
}

func TestPutOwnerCanChangeTheValue(t *testing.T) {

	key := "key1"
	value1 := "data1"
	value2 := "data2"
	owner := "testUser1"
	mockStore := NewMockStore()
	mockStore.Put(key, value1, owner)
	mockStore.Put(key, value2, owner)
	newEntry, _ := mockStore.List(key)
	if newEntry.Value != value2 {
		t.Errorf("Second Put unexpected value got %s want %s", value2, newEntry.Value)
	}
}

func TestPutDifferntOwnerCannotChangeExistingValue(t *testing.T) {

	key := "key1"
	value1 := "data1"
	value2 := "data2"
	owner := "tesUser1"
	owner2 := "testUser2"
	mockStore := NewMockStore()
	mockStore.Put(key, value1, owner)

	err := mockStore.Put(key, value2, owner2)
	if err != common.ErrorUnauthorisedOwner {
		t.Errorf("Unexpected error got %v want %v", err, common.ErrorUnauthorisedOwner)
	}
}

func TestPutReturnsErrorIfKeyNotExist(t *testing.T) {

	key := "key1"
	value := "data1"
	owner := "testUser1"
	mockStore := NewMockStore()
	mockStore.Put(key, value, owner)

	newEntry, _ := mockStore.List(key)
	if newEntry == nil {
		t.Errorf("Expected key %s not found", key)
	}

	if newEntry.Value != value {
		t.Errorf("Expected value %s not found", value)
	}
	if newEntry.Owner != owner {
		t.Errorf("Expected owner %s not found", owner)
	}
}

func TestListAllReturnsItemAdded(t *testing.T) {

	key := "key1"
	value := "data1"
	owner := "testUser1"
	mockStore := NewMockStore()
	mockStore.Put(key, value, owner)

	entries := mockStore.ListAll()

	newEntry := entries[0].String()
	expected := store.NewEntry(key, value, owner).String()

	if newEntry != expected {
		t.Errorf("Get unexpected error got %v want %v", newEntry, expected)
	}
}

func TestListReturnsEmtpyList(t *testing.T) {

	mockStore := NewMockStore()

	entries := mockStore.ListAll()
	if entries == nil {
		t.Errorf("Unexpected error: List is nil,")
	}
	if len(entries) != 0 {
		t.Errorf("Unexpected list size: got %v want %v,", len(entries), 0)
	}

}

func TestListKeyReturnsTheKeyAndEntry(t *testing.T) {

	key := "key1"
	value := "data1"
	owner := "testUser1"
	store := NewMockStore()
	store.Put(key, value, owner)

	entry, error := store.List(key)
	if error != nil {
		t.Errorf("Expected key %s not found", key)
	}

	if entry.Key != key {
		t.Errorf("Expected value %s not found", value)
	}

	if entry.Value != value {
		t.Errorf("Expected value %s not found", value)
	}

	if entry.Owner != owner {
		t.Errorf("Expected owner %s not found", owner)
	}
}

func TestListKeyReturnsErrorIfKeyNotExists(t *testing.T) {

	key := "key1"
	mockStore := NewMockStore()

	_, error := mockStore.List(key)
	if error != common.ErrorKeyNotFound {
		t.Errorf("Unexpected error: got %v want %v,", error, common.ErrorKeyNotFound)
	}
}

func TestGetReturnsValueFromMatchingKey(t *testing.T) {

	key := "key1"
	value := "data1"
	owner := "testUser1"
	mockStore := NewMockStore()
	mockStore.Put(key, value, owner)
	entryValue, _ := mockStore.Get(key)

	if entryValue != value {
		t.Errorf("Unexpected error: got %v want %v,", entryValue, value)
	}
}

func TestGetReturnsErrorIfKeyNotExists(t *testing.T) {

	key := "key1"
	mockStore := NewMockStore()
	_, err := mockStore.Get(key)
	if err != common.ErrorKeyNotFound {
		t.Errorf("Unexpected error: got %v want %v,", err, common.ErrorKeyNotFound)
	}
}

func TestDeleteReturnsErrorIfKeyNotExists(t *testing.T) {

	key := "key1"
	owner1 := "testUser1"
	mockStore := NewMockStore()
	err := mockStore.Delete(key, owner1)
	if err != common.ErrorKeyNotFound {
		t.Errorf("Unexpected error: got %v want %v,", err, common.ErrorKeyNotFound)
	}
}

func TestDeleteRemovesKey(t *testing.T) {

	key := "key1"
	value := "data1"
	owner1 := "testUser1"
	mockStore := NewMockStore()
	mockStore.Put(key, value, owner1)
	err := mockStore.Delete(key, owner1)
	if err != nil {
		t.Errorf("Unexpected error: got %v want %v,", err, "nil")
	}
	_, err = mockStore.List(key)
	if err != common.ErrorKeyNotFound {
		t.Errorf("Unexpected error: got %v want %v,", err, common.ErrorKeyNotFound)
	}
}

func TestDeleteOnlyAuthorisedOwnerCanRemoveKey(t *testing.T) {

	key := "key1"
	value := "data1"
	owner1 := "testUser1"
	owner2 := "testUser2"
	mockStore := NewMockStore()
	mockStore.Put(key, value, owner1)
	err := mockStore.Delete(key, owner2)
	if err != common.ErrorUnauthorisedOwner {
		t.Errorf("Unexpected error: got %v want %v,", err, common.ErrorKeyNotFound)
	}
}

func TestDeleteAdminCanOverrideUsery(t *testing.T) {

	key := "key1"
	value := "data1"
	owner1 := "testUser1"
	owner2 := "admin"
	mockStore := NewMockStore()
	mockStore.UserDatabase().AddUser("admin", "111")
	mockStore.Put(key, value, owner1)
	err := mockStore.Delete(key, owner2)
	if err != nil {
		t.Errorf("Unexpected error: got %v want %v,", err, "nil")
	}
}

func TestDeleteOnlyOwnerCanDleeteKey(t *testing.T) {

	key := "key1"
	value := "data1"
	owner1 := "testUser1"
	owner2 := "testUser2"
	mockStore := NewMockStore()
	mockStore.Put(key, value, owner1)
	err := mockStore.Delete(key, owner2)
	if err != common.ErrorUnauthorisedOwner {
		t.Errorf("Unexpected error: got %v want %v,", err, "nil")
	}
}

func TestPutOriginalOwnerAfterAdminOverridwsValueOnlyAuthorisedOwnerCanRemoveKey(t *testing.T) {

	key := "key1"
	value := "data1"
	owner1 := "testUser1"

	mockStore := NewMockStore()
	mockStore.Put(key, value, owner1)
	mockStore.Put(key, "updated value", "admin")
	err := mockStore.Put(key, "override updated value", owner1)
	if err != nil {
		t.Errorf("Unexpected error: got %v want %v,", err, common.ErrorKeyNotFound)
	}
}
