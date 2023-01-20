package store

import (
	"container/list"
	"demo-store/common"
	"demo-store/utils"
)

// type OrderedList interface {

// }

type LruEntryList struct {
	data        map[string]*list.Element
	orderedData *list.List
	tracer      utils.Tracer
	depth       int
}

func NewLruEntryList(tracer utils.Tracer, depth int) *LruEntryList {

	return &LruEntryList{data: make(map[string]*list.Element), orderedData: list.New(), tracer: tracer, depth: depth}
}

func (s *LruEntryList) AddEntry(key string, value string, owner string) {

	entry := NewEntry(key, value, owner)
	elem := s.orderedData.PushFront(entry)
	s.data[entry.Key] = elem

	s.tracer.LogInfo("Key", entry.Key, "added")

	if s.depth == 0 {
		return
	}

	// if list is full remove last accessed key
	for len(s.data) > s.depth {
		last := s.orderedData.Back()
		remove, _ := last.Value.(*Entry)

		s.tracer.LogInfo("Key", entry.Key, "dropped")

		s.orderedData.Remove(last)
		delete(s.data, remove.Key)
	}
}

func (s *LruEntryList) UpdateEntry(key string, value string) error {
	entry, err := s.FindEntry(key)
	if err != nil {
		return err
	}

	entry.WriteValue(value)

	// push to top as its been written
	elem := s.data[entry.Key]
	s.orderedData.MoveToFront(elem)

	s.tracer.LogInfo("Key", entry.Key, "updated")
	return nil
}

func (s *LruEntryList) ReadEntry(key string) (string, error) {

	entry, err := s.FindEntry(key)
	if err != nil {
		return "", err
	}

	valeue := entry.ReadValue()

	// push to top as its been read
	elem := s.data[entry.Key]
	s.orderedData.MoveToFront(elem)

	s.tracer.LogInfo("Key", entry.Key, "accessed")
	return valeue, nil
}

func (s *LruEntryList) DeleteEntry(key string) error {

	entry, err := s.FindEntry(key)
	if err != nil {
		return err
	}

	//remove from ordered list
	elem := s.data[key]
	s.orderedData.Remove(elem)

	// remove from dictionary
	delete(s.data, entry.Key)

	s.tracer.LogInfo("Key", entry.Key, " deleted")

	return nil
}

func (s *LruEntryList) FindEntry(key string) (*Entry, error) {

	elem, ok := s.data[key]
	if !ok {
		s.tracer.LogError("Key", key, " not found")
		return nil, common.ErrorKeyNotFound
	}

	entry := elem.Value.(*Entry)
	return entry, nil
}

func (s *LruEntryList) ListAll() []*Entry {

	// copy them in correct order as the dictionary can be accessed in random order
	entryList := make([]*Entry, 0, len(s.data))
	if len(s.data) > 0 {
		front := s.orderedData.Back()
		entry := front.Value.(*Entry)
		entryList = append(entryList, entry.Clone())

		for {

			front = front.Prev()
			if front == nil {
				break
			}
			entry := front.Value.(*Entry)
			entryList = append(entryList, entry.Clone())
		}
	}

	return entryList
}
