package store_test

import (
	"demo-store/common"
	"demo-store/store"
	"fmt"
	"testing"
)

var key1 = "key1"
var value1 = "value1"
var owner1 = "user1"

var key2 = "key2"
var value2 = "value2"
var owner2 = "user2"

func TestLruListReadEntry(t *testing.T) {

	list := store.NewLruEntryList(CreateMockTracer(), 0)

	list.AddEntry(key1, value1, owner1)

	value, err := list.ReadEntry(key1)

	if err != nil {
		t.Errorf("Returned unexpected error: got %v want %v", err, "nil")
	}
	if value != value1 {
		t.Errorf("Returned unexpected value: got %v want %v", value, value1)
	}
}

func TestLruListUpdateEntryIfKeyExists(t *testing.T) {

	list := store.NewLruEntryList(CreateMockTracer(), 0)

	list.AddEntry(key1, value1, owner1)

	err := list.UpdateEntry(key1, value2)

	if err != nil {
		t.Errorf("Returned unexpected error: got %v want %v", err, "nil")
	}

	value, err := list.ReadEntry(key1)
	if value != value2 {
		t.Errorf("Returned unexpected value: got %v want %v", value, value1)
	}
}

func TestLruListUpdateEntryFailsIfKeyExists(t *testing.T) {

	list := store.NewLruEntryList(CreateMockTracer(), 0)

	err := list.UpdateEntry(key1, value2)

	expected := common.ErrorKeyNotFound
	if err != expected {
		t.Errorf("Returned unexpected error: got %v want %v", err, expected)
	}
}

func TestLruListUpdateEntryUpdatesWritesCount(t *testing.T) {

	list := store.NewLruEntryList(CreateMockTracer(), 0)

	list.AddEntry(key1, value1, owner1)
	numbOfWrites := 10
	for i := 0; i < numbOfWrites; i++ {
		list.UpdateEntry(key1, value2)
	}

	expected := numbOfWrites + 1 // include the initial added write
	entry, _ := list.FindEntry(key1)

	if entry.Writes != expected {
		t.Errorf("Returned unexpected writes count: got %v want %v", entry.Writes, expected)
	}
}

func TestLruListReadEntryFailsIfKeyExists(t *testing.T) {

	list := store.NewLruEntryList(CreateMockTracer(), 0)

	_, err := list.ReadEntry(key1)

	expected := common.ErrorKeyNotFound
	if err != expected {
		t.Errorf("Returned unexpected error: got %v want %v", err, expected)
	}
}

func TestLruListReadEntryUpdatesReadCount(t *testing.T) {

	list := store.NewLruEntryList(CreateMockTracer(), 0)

	list.AddEntry(key1, value1, owner1)
	numbOfReads := 10
	for i := 0; i < numbOfReads; i++ {
		list.ReadEntry(key1)
	}

	expected := numbOfReads
	entry, _ := list.FindEntry(key1)
	if entry.Reads != expected {
		t.Errorf("Returned unexpected reads count: got %v want %v", entry.Reads, expected)
	}
}

func TestLruListFindEntryFailsIfKeyExists(t *testing.T) {

	list := store.NewLruEntryList(CreateMockTracer(), 0)

	_, err := list.FindEntry(key1)

	expected := common.ErrorKeyNotFound
	if err != expected {
		t.Errorf("Returned unexpected error: got %v want %v", err, expected)
	}
}

func TestLruListFindEntryReturnsEntry(t *testing.T) {

	list := store.NewLruEntryList(CreateMockTracer(), 0)
	list.AddEntry(key1, value1, owner1)
	entry, _ := list.FindEntry(key1)

	expected := key1
	if entry.Key != expected {
		t.Errorf("Returned unexpected error: got %v want %v", entry.Key, expected)
	}

	expectedEntry := store.NewEntry(key1, value1, owner1).String()
	if entry.String() != expectedEntry {
		t.Errorf("Returned unexpected error: got %v want %v", entry.String(), expectedEntry)
	}
}

func TestLruDeleteEntryFailsIfKeyExists(t *testing.T) {

	list := store.NewLruEntryList(CreateMockTracer(), 0)

	err := list.DeleteEntry(key1)

	expected := common.ErrorKeyNotFound
	if err != expected {
		t.Errorf("Returned unexpected error: got %v want %v", err, expected)
	}
}

func TestLruDeleteEntryRemovesEntry(t *testing.T) {

	list := store.NewLruEntryList(CreateMockTracer(), 0)
	list.AddEntry(key1, value1, owner1)
	list.DeleteEntry(key1)

	_, err := list.FindEntry(key1)
	expected := common.ErrorKeyNotFound
	if err != expected {
		t.Errorf("Returned unexpected error: got %v want %v", err, expected)
	}
}

func TestLruDepthRemoveTheFirstItemsAdded(t *testing.T) {
	depth := 10
	list := store.NewLruEntryList(CreateMockTracer(), depth)

	numberOfInserts := 15
	for i := 0; i < numberOfInserts; i++ {
		list.AddEntry(fmt.Sprint("key", i), value1, owner1)
	}

	numberOfRemoved := 5
	for i := 0; i < numberOfRemoved; i++ {

		_, err := list.FindEntry(fmt.Sprint("key", i))
		expected := common.ErrorKeyNotFound
		if err != expected {
			t.Errorf("Returned unexpected error: got %v want %v", err, expected)
		}
	}
}

func TestLruDepthPushesFirsKeyToTopAfterARead(t *testing.T) {
	depth := 5
	list := store.NewLruEntryList(CreateMockTracer(), depth)

	numberOfInserts := 5
	for i := 0; i < numberOfInserts; i++ {
		list.AddEntry(fmt.Sprint("key", i), value1, owner1)
	}

	expected := value2
	list.UpdateEntry("key0", expected)

	numberOfNewInserts := 2
	for i := numberOfInserts; i < numberOfNewInserts; i++ {
		list.AddEntry(fmt.Sprint("key", i), value1, owner1)
	}

	entry, _ := list.FindEntry("key0")

	if entry.Value != expected {
		t.Errorf("Returned unexpected value: got %v want %v", entry.Value, expected)
	}
}

func TestListAllReturnsAllEntrie(t *testing.T) {
	list := store.NewLruEntryList(CreateMockTracer(), 0)

	numberOfInserts := 15
	for i := 0; i < numberOfInserts; i++ {
		list.AddEntry(fmt.Sprint("key", i), fmt.Sprint("value", i), fmt.Sprint("owner", i))
	}

	entries := list.ListAll()

	for i, entry := range entries {

		expectedEntry := store.NewEntry(fmt.Sprint("key", i), fmt.Sprint("value", i), fmt.Sprint("owner", i))
		if entry.String() != expectedEntry.String() {
			t.Errorf("Returned unexpected key: got %v want %v", entry.Key, expectedEntry.String())
		}
	}
}
