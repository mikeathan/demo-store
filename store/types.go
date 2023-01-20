package store

import (
	"demo-store/users"
	"demo-store/utils"
)

type Store interface {
	RegisterShutdownListener(listener *ShutdownListener)
	MakePutRequest(key string, value string, owner string) error
	MakeGetRequest(key string) (string, error)
	MakeListAllRequest() []*Entry
	MakeListRequest(key string) (*Entry, error)
	MakeDeleteRequest(key string, owner string) error
	MakeShutdownRequest()
	UserDatabase() users.UserDatabase
}

type KvStore struct {
	Tracer           utils.Tracer
	userDatabase     users.UserDatabase
	lruData          LruEntryList
	putChannel       chan PutRequest
	getChannel       chan GetRequest
	listAllChannel   chan ListAllRequest
	listChannel      chan ListRequest
	deleteChannel    chan DeleteRequest
	shutdownChannel  chan ShutdownRequest
	shutdownListener *ShutdownListener
}
