package store_test

import (
	"demo-store/common"
	"demo-store/store"
	"demo-store/users"
	"demo-store/utils"
	"testing"
	"time"
)

func TestPutRequestInsertsToStore(t *testing.T) {

	key := "key1"
	value1 := "data1"
	owner := "testUser1"
	mockStore := NewMockStore()

	err := mockStore.MakePutRequest(key, value1, owner)
	if err != nil {
		t.Errorf("Put unexpected value got %s want %s", err, "nil")
	}

	value, _ := mockStore.MakeGetRequest(key)
	if value != value1 {
		t.Errorf("Get unexpected value got %s want %s", value1, value)
	}
}

func TestGetReturnsErrorIfItemNotInStore(t *testing.T) {

	key := "key1"
	mockStore := NewMockStore()

	_, err := mockStore.MakeGetRequest(key)
	if err != common.ErrorKeyNotFound {
		t.Errorf("Get unexpected error got %v want %v", err, common.ErrorKeyNotFound)
	}
}

func TestListAllReturnsAllItems(t *testing.T) {

	key := "key1"
	value := "data1"
	owner := "testUser1"
	mockStore := NewMockStore()
	mockStore.Put(key, value, owner)

	entries := mockStore.MakeListAllRequest()
	newEntry := entries[0].String()
	expected := store.NewEntry(key, value, owner).String()

	if newEntry != expected {
		t.Errorf("Get unexpected error got %v want %v", newEntry, expected)
	}
}

func TestListReturnsItem(t *testing.T) {

	key := "key1"
	value := "data1"
	owner := "testUser1"
	mockStore := NewMockStore()
	mockStore.Put(key, value, owner)

	entry, _ := mockStore.MakeListRequest(key)
	newEntry := entry.String()
	expected := store.NewEntry(key, value, owner).String()

	if newEntry != expected {
		t.Errorf("Get unexpected error got %v want %v", newEntry, expected)
	}
}

func TestDeleteRemovesItem(t *testing.T) {

	key := "key1"
	value := "data1"
	owner := "testUser1"
	mockStore := NewMockStore()
	mockStore.Put(key, value, owner)

	mockStore.MakeDeleteRequest(key, owner)

	_, err := mockStore.MakeListRequest(key)
	if err != common.ErrorKeyNotFound {
		t.Errorf("Get unexpected error got %v want %v", common.ErrorKeyNotFound, err)
	}
}

func TestShutdown(t *testing.T) {

	mockStore := NewMockStore()
	sl := store.CreateShutdownListener()
	mockStore.RegisterShutdownListener(sl)
	mockStore.MakeShutdownRequest()

	go func() {
		resp := <-sl.Listener
		if resp != true {
			t.Errorf("Get unexpected error got %v want %v", false, true)
		}
	}()

	time.Sleep(1 * time.Second)
}

type MockTracer struct {
}

func (h *MockTracer) LogInfo(message ...any) {

}

func (h *MockTracer) LogError(message ...any) {

}

func (h *MockTracer) LogWarning(message ...any) {

}

func (h *MockTracer) Close() {
}

func CreateMockTracer() utils.Tracer {
	return &MockTracer{}
}

func NewMockStore() *store.KvStore {
	return store.CreateKvStore(CreateMockTracer(), users.CreateUserDatabase(), 0)
}
