package endpoints

import (
	"demo-store/common"
	"demo-store/store"
	"demo-store/users"
	"demo-store/utils"
	"errors"
	"io/ioutil"
	"net/http"
)

const (
	PathParameter     = "path"
	UsernameParameter = "username"
)

type HttpMethodHandler interface {
	HttpMethod() string
	Handle(args *HttpMethodHandlerParams, resp http.ResponseWriter, req *http.Request) HttpResult
}

type HttpMethodHandlerParams struct {
	data map[string]string
}

func CreatePathParameter(path string) *HttpMethodHandlerParams {
	return &HttpMethodHandlerParams{data: map[string]string{PathParameter: path}}
}

func CreatePathAndUsernameParameter(path string, username string) *HttpMethodHandlerParams {
	return &HttpMethodHandlerParams{data: map[string]string{PathParameter: path, UsernameParameter: username}}
}

func (p *HttpMethodHandlerParams) Add(key string, value string) {
	p.data[key] = value
}

func (p *HttpMethodHandlerParams) Get(key string) string {
	value, ok := p.data[key]
	if !ok {
		return ""
	}
	return value
}

type HttpResult struct {
	Message string
	Code    int
}

func CreateHttpResponseFromError(err error) HttpResult {
	switch {
	case errors.Is(err, common.ErrorValidatingJwtToken):
		return CreateHttpResponse(err.Error(), http.StatusUnauthorized)

	case errors.Is(err, common.ErrorKeyNotFound):
		return CreateHttpResponse("Key not found", http.StatusNotFound)

	case errors.Is(err, common.ErrorUnauthorisedOwner):
		return CreateHttpResponse("Forbiden", http.StatusForbidden)

	case errors.Is(err, common.ErrorAuthorizationFailed):
		return CreateHttpResponse("Unauthorized", http.StatusUnauthorized)

	case errors.Is(err, common.ErrorInvalidAuthorizationHeader):
		return CreateHttpResponse(err.Error(), http.StatusUnprocessableEntity)

	case errors.Is(err, common.ErrorAuthorizationHeaderMissing):
		return CreateHttpResponse("Forbidden", http.StatusForbidden)

	case errors.Is(err, common.ErrorKeyNotSet):
		return CreateHttpResponse(err.Error(), http.StatusUnprocessableEntity)

	case errors.Is(err, common.ErrorStoreValueNotSet):
		return CreateHttpResponse(err.Error(), http.StatusUnprocessableEntity)

	default:
		return CreateHttpResponse(err.Error(), http.StatusInternalServerError)
	}
}

type Routes struct {
	Insecure []Route
	Secure   []Route
}

func APIRoutes(tracer utils.Tracer, kvStore store.Store) *Routes {

	routes := Routes{Insecure: []Route{}, Secure: []Route{}}
	authenticator := NewRouteAuthenticator(tracer)

	routes.Secure = append(routes.Secure, CreateStoreRoute(tracer, kvStore, authenticator))
	routes.Secure = append(routes.Secure, CreateListRoute(tracer, kvStore, authenticator))
	routes.Secure = append(routes.Secure, CreateShutdownRoute(tracer, kvStore, authenticator))

	routes.Insecure = append(routes.Insecure, CreatePingRoute(tracer, kvStore))
	routes.Insecure = append(routes.Insecure, CreateLoginRoute(tracer, kvStore.UserDatabase()))

	return &routes
}

func CreateHttpResponse(message string, code int) HttpResult {
	return HttpResult{Message: message, Code: code}
}

func CreateStoreRoute(tracer utils.Tracer, kvStore store.Store, authenticator Authenticator) Route {
	var methods []HttpMethodHandler
	methods = append(methods, CreatePut(tracer, kvStore))
	methods = append(methods, CreateGet(tracer, kvStore))
	methods = append(methods, CreateDelete(tracer, kvStore))

	return &SecureRoute{Path: "/store/", Tracer: tracer, MethodHandlers: methods, Authenticator: authenticator}
}

func CreatePingRoute(tracer utils.Tracer, kvStore store.Store) Route {
	var methods []HttpMethodHandler
	methods = append(methods, CreatePing(tracer))

	return &InsecureRoute{Path: "/ping/", Tracer: tracer, MethodHandlers: methods}
}

func CreateListRoute(tracer utils.Tracer, kvStore store.Store, authenticator Authenticator) Route {
	var methods []HttpMethodHandler
	methods = append(methods, CreateList(tracer, kvStore))

	return &SecureRoute{Path: "/list/", Tracer: tracer, MethodHandlers: methods, Authenticator: authenticator}
}

func CreateShutdownRoute(tracer utils.Tracer, kvStore store.Store, authenticator Authenticator) Route {
	var methods []HttpMethodHandler
	methods = append(methods, CreateShutdown(tracer, kvStore))

	return &SecureRoute{Path: "/shutdown/", Tracer: tracer, MethodHandlers: methods, Authenticator: authenticator}
}

func CreateLoginRoute(tracer utils.Tracer, users users.UserDatabase) Route {
	var methods []HttpMethodHandler
	methods = append(methods, CreateLogin(tracer, users))

	return &InsecureRoute{Path: "/login/", Tracer: tracer, MethodHandlers: methods}
}

func CreatePing(tracer utils.Tracer) *PingHandler {
	return &PingHandler{Tracer: tracer, httpMethod: http.MethodGet}
}

func CreatePut(tracer utils.Tracer, kvStore store.Store) *PutHandler {
	return &PutHandler{Tracer: tracer, httpMethod: http.MethodPut, store: kvStore}
}

func CreateGet(tracer utils.Tracer, kvStore store.Store) *GetHandler {
	return &GetHandler{Tracer: tracer, httpMethod: http.MethodGet, store: kvStore}
}

func CreateList(tracer utils.Tracer, kvStore store.Store) *ListHandler {
	return &ListHandler{Tracer: tracer, httpMethod: http.MethodGet, store: kvStore}
}

func CreateDelete(tracer utils.Tracer, kvStore store.Store) *DeleteHandler {
	return &DeleteHandler{Tracer: tracer, httpMethod: http.MethodDelete, store: kvStore}
}

func CreateShutdown(tracer utils.Tracer, kvStore store.Store) *ShutdownHandler {
	return &ShutdownHandler{Tracer: tracer, httpMethod: http.MethodGet, store: kvStore}
}

func CreateLogin(tracer utils.Tracer, users users.UserDatabase) *LoginHandler {
	return &LoginHandler{Tracer: tracer, httpMethod: http.MethodGet, Users: users, Tokenizer: utils.NewJwtTokenizer(tracer)}
}

func CreateLoginWithTokenizer(tracer utils.Tracer, users users.UserDatabase, tokenizer utils.Tokenizer) *LoginHandler {
	return &LoginHandler{Tracer: tracer, httpMethod: http.MethodGet, Users: users, Tokenizer: tokenizer}
}

func GetBody(req *http.Request) string {
	if req.Body == nil {
		return ""
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {

		return ""
	}

	return string(body)
}
