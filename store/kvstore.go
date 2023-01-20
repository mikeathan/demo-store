package store

import (
	"demo-store/common"
	"demo-store/users"
	"demo-store/utils"
	"time"
)

func (s *KvStore) MakePutRequest(key string, value string, owner string) error {
	req := CreatePutRequest(key, value, owner)

	s.putChannel <- req
	return <-req.Response
}

func (s *KvStore) MakeGetRequest(key string) (string, error) {
	req := CreateGetRequest(key)
	s.getChannel <- req

	resp := <-req.Response
	return resp.Value, resp.Error
}

func (s *KvStore) MakeListAllRequest() []*Entry {
	req := CreateListAllRequest()
	s.listAllChannel <- req

	return <-req.Response
}

func (s *KvStore) MakeListRequest(key string) (*Entry, error) {
	req := CreateListRequest(key)
	s.listChannel <- req

	resp := <-req.Response
	return resp.Entry, resp.Error
}

func (s *KvStore) MakeDeleteRequest(key string, owner string) error {
	req := CreateDeleteRequest(key, owner)
	s.deleteChannel <- req

	return <-req.Response
}

func (s *KvStore) MakeShutdownRequest() {

	req := CreateShutdownRequest()
	s.shutdownChannel <- req
}

func (s *KvStore) RegisterShutdownListener(listener *ShutdownListener) {
	s.shutdownListener = listener
}

func CreateKvStore(tracer utils.Tracer, users users.UserDatabase, depth int) *KvStore {

	kvStore := &KvStore{
		Tracer:           tracer,
		lruData:          *NewLruEntryList(tracer, depth),
		putChannel:       make(chan PutRequest),
		getChannel:       make(chan GetRequest),
		listAllChannel:   make(chan ListAllRequest),
		listChannel:      make(chan ListRequest),
		deleteChannel:    make(chan DeleteRequest),
		shutdownChannel:  make(chan ShutdownRequest),
		shutdownListener: nil,
	}

	kvStore.userDatabase = users
	kvStore.monitor()

	return kvStore
}
func (s *KvStore) UserDatabase() users.UserDatabase {
	return s.userDatabase
}

func (s *KvStore) monitor() {
	go func() {
		shutdown := false
		for !shutdown {
			select {
			case req := <-s.putChannel:
				err := s.Put(req.Key, req.Value, req.Owner)
				req.Response <- err

			case req := <-s.getChannel:
				value, err := s.Get(req.Key)
				req.Response <- CreateGetResponse(value, err)

			case req := <-s.listAllChannel:
				value := s.ListAll()
				req.Response <- value

			case req := <-s.listChannel:
				value, err := s.List(req.Key)
				req.Response <- CreateListResponse(value, err)

			case req := <-s.deleteChannel:
				err := s.Delete(req.Key, req.Owner)
				req.Response <- err

			case <-s.shutdownChannel:
				shutdown = true
				if s.shutdownListener != nil {

					time.Sleep(500 * time.Millisecond)
					s.shutdownListener.Listener <- shutdown
				}
			}
		}
	}()
}

func (s *KvStore) Put(key string, value string, owner string) error {

	entry, err := s.lruData.FindEntry(key)
	if err != nil {
		s.lruData.AddEntry(key, value, owner)
		return nil

	} else if entry.Owner == owner || s.userDatabase.IsAdmin(owner) {
		s.lruData.UpdateEntry(key, value)
		return nil
	}
	s.Tracer.LogError("User", owner, " cannot update key.")
	return common.ErrorUnauthorisedOwner
}

func (s *KvStore) Get(key string) (string, error) {

	value, err := s.lruData.ReadEntry(key)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (s *KvStore) ListAll() []*Entry {

	return s.lruData.ListAll()
}

func (s *KvStore) List(key string) (*Entry, error) {

	entry, err := s.lruData.FindEntry(key)
	if err != nil {
		return nil, err
	}

	return entry.Clone(), nil
}

func (s *KvStore) Delete(key string, owner string) error {

	entry, err := s.lruData.FindEntry(key)
	if err != nil {
		return err
	}

	if entry.Owner == owner || s.userDatabase.IsAdmin(owner) {
		s.lruData.DeleteEntry(key)
		return nil
	}

	s.Tracer.LogError("User", owner, " cannot delete key.")
	return common.ErrorUnauthorisedOwner
}
