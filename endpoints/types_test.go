package endpoints_test

import (
	"demo-store/common"
	"demo-store/endpoints"
	"demo-store/store"
	"demo-store/users"
	"demo-store/utils"
	"errors"
	"net/http"
	"strings"
	"testing"
)

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

type TestInput struct {
	Key   string
	Value string
	Owner string
}

type MockAuthenticator struct {
	MockValue string
	MockError error
	Tokenizer utils.Tokenizer
}

func NewMockAuthenticatorWithError(err error) endpoints.Authenticator {
	return &MockAuthenticator{MockError: err}
}
func NewMockAuthenticatorWithValue(value string) endpoints.Authenticator {
	return &MockAuthenticator{MockValue: value}
}

func NewMockAuthenticator(value string) endpoints.Authenticator {
	var err error
	if value == "" {
		err = common.ErrorAuthorizationHeaderMissing
	}
	return &MockAuthenticator{MockValue: value, MockError: err}
}

func NewMockAuthenticatorWithToken(value string, tokenizer utils.Tokenizer) endpoints.Authenticator {
	var err error
	if value == "" {
		err = common.ErrorAuthorizationHeaderMissing
	}
	return &MockAuthenticator{MockValue: value, MockError: err, Tokenizer: tokenizer}
}

func (a *MockAuthenticator) GetUsername(bearerToken string) (string, error) {
	if a.MockError != nil {
		return "", a.MockError
	}

	return a.MockValue, nil
}

func CreateTestInput(key string, value string, owner string) TestInput {
	return TestInput{Key: key, Value: value, Owner: owner}
}

var input1 TestInput = CreateTestInput("key1", "some value 1", "user1")
var input2 TestInput = CreateTestInput("key2", "some value 2", "user2")

func CreateMockRouteWithDelete(path string, tracer utils.Tracer, kvStore store.Store, authenticator endpoints.Authenticator) endpoints.Route {
	var methods []endpoints.HttpMethodHandler
	methods = append(methods, endpoints.CreateDelete(tracer, kvStore))

	return &endpoints.SecureRoute{Path: path, Tracer: tracer, MethodHandlers: methods, Authenticator: authenticator}
}

func CreateMockRouteWithGet(path string, tracer utils.Tracer, kvStore store.Store, authenticator endpoints.Authenticator) endpoints.Route {
	var methods []endpoints.HttpMethodHandler
	methods = append(methods, endpoints.CreateGet(tracer, kvStore))

	return &endpoints.SecureRoute{Path: path, Tracer: tracer, MethodHandlers: methods, Authenticator: authenticator}
}

func CreateMockRouteWithPut(path string, tracer utils.Tracer, kvStore store.Store, authenticator endpoints.Authenticator) endpoints.Route {
	var methods []endpoints.HttpMethodHandler
	methods = append(methods, endpoints.CreatePut(tracer, kvStore))

	return &endpoints.SecureRoute{Path: path, Tracer: tracer, MethodHandlers: methods, Authenticator: authenticator}
}

func CreateMockRouteWithList(path string, tracer utils.Tracer, kvStore store.Store, authenticator endpoints.Authenticator) endpoints.Route {
	var methods []endpoints.HttpMethodHandler
	methods = append(methods, endpoints.CreateList(tracer, kvStore))

	return &endpoints.SecureRoute{Path: path, Tracer: tracer, MethodHandlers: methods, Authenticator: authenticator}
}

func CreateMockRouteWithPing(path string, tracer utils.Tracer) endpoints.Route {
	var methods []endpoints.HttpMethodHandler
	methods = append(methods, endpoints.CreatePing(tracer))

	return &endpoints.InsecureRoute{Path: path, Tracer: tracer, MethodHandlers: methods}
}

func CreateMockRouteWithLogin(path string, tracer utils.Tracer, userDatabase users.UserDatabase, tokenizer utils.Tokenizer) endpoints.Route {
	var methods []endpoints.HttpMethodHandler
	methods = append(methods, endpoints.CreateLoginWithTokenizer(tracer, userDatabase, tokenizer))

	return &endpoints.InsecureRoute{Path: path, Tracer: tracer, MethodHandlers: methods}
}

func CreateMockRouteWithShutdown(path string, tracer utils.Tracer, kvStore store.Store, authenticator endpoints.Authenticator) endpoints.Route {
	var methods []endpoints.HttpMethodHandler
	methods = append(methods, endpoints.CreateShutdown(tracer, kvStore))

	return &endpoints.SecureRoute{Path: path, Tracer: tracer, MethodHandlers: methods, Authenticator: authenticator}
}

func TestGetBodySuccess(t *testing.T) {

	expected := "somedata"
	req, _ := http.NewRequest(http.MethodGet, "", strings.NewReader(expected))
	res := endpoints.GetBody(req)
	if res != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res, expected)
	}
}

func TestGetBodyNullReader(t *testing.T) {

	expected := ""
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	res := endpoints.GetBody(req)
	if res != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res, expected)
	}
}

func TestGetBodyClosedReader(t *testing.T) {

	r := MockReader{}

	req, _ := http.NewRequest(http.MethodGet, "", &r)
	res := endpoints.GetBody(req)
	expected := ""
	if res != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", res, expected)
	}
}

func AssertErrorHttpCode(err error, got int, t *testing.T) {
	want := endpoints.CreateHttpResponseFromError(err)
	if got != want.Code {
		t.Errorf("handler returned unexpected code: got %v want %v", got, want)
		t.Errorf("				   unexpected message: got %v want %v", err.Error(), want.Message)
	}
}

func NewMockStore() *store.KvStore {
	return store.CreateKvStore(CreateMockTracer(), users.CreateUserDatabase(), 0)
}

type MockReader struct {
}

func (e *MockReader) Read(p []byte) (int, error) {
	return 0, errors.New("test errr")
}

func TestHttpMethodHandlerParamsPath(t *testing.T) {

	param := endpoints.CreatePathParameter(input1.Value)

	if param.Get(endpoints.PathParameter) != input1.Value {
		t.Errorf("handler returned unexpected value: got %v want %v", param.Get(endpoints.PathParameter), input1.Value)
	}
}

func TestHttpMethodHandlerParamsPathAndUsername(t *testing.T) {

	param := endpoints.CreatePathAndUsernameParameter(input1.Value, input1.Owner)

	if param.Get(endpoints.PathParameter) != input1.Value {
		t.Errorf("handler returned unexpected value: got %v want %v", param.Get(endpoints.PathParameter), input1.Value)
	}

	if param.Get(endpoints.UsernameParameter) != input1.Owner {
		t.Errorf("handler returned unexpected value: got %v want %v", param.Get(endpoints.UsernameParameter), input1.Owner)
	}
}

func TestHttpMethodHandlerParamsAdd(t *testing.T) {

	param := endpoints.CreatePathAndUsernameParameter("path", "username")

	param.Add(input1.Key, input1.Value)
	if param.Get(input1.Key) != input1.Value {
		t.Errorf("handler returned unexpected value: got %v want %v", param.Get(input1.Key), input1.Value)
	}
}

func TestHttpMethodHandlerParamsGet(t *testing.T) {

	param := endpoints.CreatePathAndUsernameParameter("path", "username")

	expected := ""
	if param.Get(input1.Key) != expected {
		t.Errorf("handler returned unexpected value: got %v want %v", param.Get(input1.Key), expected)
	}
}

func TestCreateHttpResponseFromErrorErrorValidatingJwtToken(t *testing.T) {

	resp := endpoints.CreateHttpResponseFromError(common.ErrorValidatingJwtToken)

	expected := http.StatusUnauthorized
	if resp.Code != expected {
		t.Errorf("handler returned unexpected value: got %v want %v", resp.Code, expected)
	}
}
