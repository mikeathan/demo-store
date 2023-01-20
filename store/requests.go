package store

type PutRequest struct {
	Key      string
	Value    string
	Owner    string
	Response chan error
}

type GetRequest struct {
	Key      string
	Response chan GetResponse
}

type ListAllRequest struct {
	Response chan []*Entry
}

type ListRequest struct {
	Key      string
	Response chan ListResponse
}

type DeleteRequest struct {
	Key      string
	Owner    string
	Response chan error
}

type ShutdownRequest struct {
}

type ShutdownListener struct {
	Listener chan bool
}

type GetResponse struct {
	Value string
	Error error
}
type ListResponse struct {
	Entry *Entry
	Error error
}

func CreatePutRequest(key string, value string, owner string) PutRequest {
	return PutRequest{Key: key, Value: value, Owner: owner, Response: make(chan error)}
}

func CreateGetRequest(key string) GetRequest {
	return GetRequest{Key: key, Response: make(chan GetResponse)}
}

func CreateListAllRequest() ListAllRequest {
	return ListAllRequest{Response: make(chan []*Entry)}
}

func CreateListRequest(key string) ListRequest {
	return ListRequest{Key: key, Response: make(chan ListResponse)}
}

func CreateDeleteRequest(key string, owner string) DeleteRequest {
	return DeleteRequest{Key: key, Owner: owner, Response: make(chan error)}
}

func CreateShutdownRequest() ShutdownRequest {
	return ShutdownRequest{}
}

func CreateGetResponse(value string, err error) GetResponse {
	return GetResponse{Value: value, Error: err}
}

func CreateListResponse(entry *Entry, err error) ListResponse {
	return ListResponse{Entry: entry, Error: err}
}

func CreateShutdownListener() *ShutdownListener {
	return &ShutdownListener{Listener: make(chan bool)}
}
