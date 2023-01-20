package endpoints_test

import (
	"demo-store/common"
	"demo-store/users"
	"demo-store/utils"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func CreateMockUserDatabase() users.UserDatabase {
	return users.CreateUserDatabase()
}

func createMockLoginRequest(userDatabase users.UserDatabase, authenticationHeader string, tokenizer utils.Tokenizer) *httptest.ResponseRecorder {
	path := "login"

	route := CreateMockRouteWithLogin(path, &MockTracer{}, userDatabase, tokenizer)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if authenticationHeader != "" {
		req.Header.Set("Authorization", authenticationHeader)
	}

	if err != nil {
		return nil
	}

	rr := httptest.NewRecorder()
	route.ServeHTTP(rr, req)

	return rr
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func TestLoginSuccess(t *testing.T) {

	userDb := CreateMockUserDatabase()
	userDb.AddUser(input1.Owner, "abc")
	basicAuthHeader := "Basic " + basicAuth(input1.Owner, "abc")

	rr := createMockLoginRequest(userDb, basicAuthHeader, utils.NewJwtTokenizer(&MockTracer{}))

	expectedStatus := http.StatusOK
	if rr.Code != expectedStatus {
		t.Errorf("handler returned unexpected code: got %v want %v", rr.Code, expectedStatus)
	}
}

func TestLoginFailedWithMissingHeader(t *testing.T) {

	userDb := CreateMockUserDatabase()
	userDb.AddUser(input1.Owner, "abc")
	rr := createMockLoginRequest(userDb, "", utils.NewJwtTokenizer(&MockTracer{}))

	AssertErrorHttpCode(common.ErrorInvalidAuthorizationHeader, rr.Code, t)
}

func TestLoginFailedWithWrongUsername(t *testing.T) {

	userDb := CreateMockUserDatabase()
	userDb.AddUser(input1.Owner, "abc")
	basicAuthHeader := "Basic " + basicAuth("wrong", "abc")

	rr := createMockLoginRequest(userDb, basicAuthHeader, utils.NewJwtTokenizer(&MockTracer{}))

	AssertErrorHttpCode(common.ErrorAuthorizationFailed, rr.Code, t)
}

func TestLoginFailedWithWrongPassword(t *testing.T) {

	userDb := CreateMockUserDatabase()
	userDb.AddUser(input1.Owner, "abc")
	basicAuthHeader := "Basic " + basicAuth(input1.Owner, "cbd")

	rr := createMockLoginRequest(userDb, basicAuthHeader, utils.NewJwtTokenizer(&MockTracer{}))

	AssertErrorHttpCode(common.ErrorAuthorizationFailed, rr.Code, t)
}

func TestLoginFailedWithMissingUsername(t *testing.T) {

	userDb := CreateMockUserDatabase()
	userDb.AddUser(input1.Owner, "abc")
	basicAuthHeader := "Basic " + basicAuth("", "abc")

	rr := createMockLoginRequest(userDb, basicAuthHeader, utils.NewJwtTokenizer(&MockTracer{}))

	AssertErrorHttpCode(common.ErrorAuthorizationHeaderMissing, rr.Code, t)
}

func TestLoginFailedWithMissingPassword(t *testing.T) {

	userDb := CreateMockUserDatabase()
	userDb.AddUser(input1.Owner, "abc")
	basicAuthHeader := "Basic " + basicAuth(input1.Owner, "")

	rr := createMockLoginRequest(userDb, basicAuthHeader, utils.NewJwtTokenizer(&MockTracer{}))

	AssertErrorHttpCode(common.ErrorAuthorizationHeaderMissing, rr.Code, t)
}

func TestLoginFailedWithInvalidToken(t *testing.T) {

	userDb := CreateMockUserDatabase()
	userDb.AddUser(input1.Owner, "abc")
	basicAuthHeader := "Basic " + basicAuth(input1.Owner, "abc")

	tokenizer := NewMockTokenizer("", errors.New("invalid token"))
	rr := createMockLoginRequest(userDb, basicAuthHeader, tokenizer)

	AssertErrorHttpCode(common.ErrorAuthorizationFailed, rr.Code, t)
}

type MockTokenizer struct {
	MockValue string
	MockError error
}

func (t *MockTokenizer) GetUsernameFromToken(tokenString string) (string, error) {

	return t.MockValue, t.MockError
}

func (t *MockTokenizer) CreateToken(username string) (string, error) {
	return t.MockValue, t.MockError
}

func NewMockTokenizer(value string, err error) utils.Tokenizer {
	return &MockTokenizer{MockValue: value, MockError: err}
}
