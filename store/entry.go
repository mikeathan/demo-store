package store

import (
	"fmt"
	"time"
)

type Entry struct {
	Key    string `json:"key"`
	Value  string `json:"-"`
	Owner  string `json:"owner"`
	Reads  int    `json:"reads"`
	Writes int    `json:"writes"`
	Age    int64  `json:"age"`

	Timestamp time.Time `json:"-"`
}

func NewEntry(key string, value string, owner string) *Entry {
	entry := Entry{
		Key:   key,
		Value: value,
		Owner: owner,
		Reads: 0,
	}

	entry.WriteValue(value)
	return &entry
}

func (e *Entry) Clone() *Entry {
	newEntry := Entry{
		Key:       e.Key,
		Value:     e.Value,
		Owner:     e.Owner,
		Reads:     e.Reads,
		Writes:    e.Writes,
		Timestamp: e.Timestamp,
	}

	newEntry.Age = time.Since(e.Timestamp).Milliseconds()
	return &newEntry
}

func (e *Entry) ReadValue() string {
	e.Reads++
	e.Timestamp = time.Now()
	return e.Value
}

func (e *Entry) WriteValue(value string) {
	e.Value = value
	e.Writes++
	e.Timestamp = time.Now()
}

func (e *Entry) String() string {
	return fmt.Sprintf("key: %s value: %s owner: %s", e.Key, e.Value, e.Owner)
}
