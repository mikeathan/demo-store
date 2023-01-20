package endpoints_test

import (
	"demo-store/common"
	"demo-store/endpoints"
	"demo-store/store"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func createMockPutRequestWithAuthenticator(store store.Store, url string, authenticator endpoints.Authenticator, body string) *httptest.ResponseRecorder {
	path := "store"

	route := CreateMockRouteWithPut(path, CreateMockTracer(), store, authenticator)
	req, err := http.NewRequest(http.MethodPut, path+url, strings.NewReader(body))

	if err != nil {
		return nil
	}

	rr := httptest.NewRecorder()
	route.ServeHTTP(rr, req)

	return rr
}

func createMockPutRequestWithUsername(store store.Store, url string, username string, body string) *httptest.ResponseRecorder {

	return createMockPutRequestWithAuthenticator(store, url, NewMockAuthenticator(username), body)
}

func TestPutReturnsSuccess(t *testing.T) {

	mockStore := NewMockStore()
	rr := createMockPutRequestWithUsername(mockStore, input1.Key, input1.Owner, input1.Value)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("handler returned unexpected code: got %v want %v", rr.Code, expected)
	}
}

func TestPutStatusUnprocessableEntityIfNotKeyPassed(t *testing.T) {

	mockStore := NewMockStore()
	rr := createMockPutRequestWithUsername(mockStore, "", input1.Owner, input1.Value)

	AssertErrorHttpCode(common.ErrorKeyNotSet, rr.Code, t)
}

func TestPutStatusUnprocessableEntityIfNoAuthenticationIsPassed(t *testing.T) {

	err := common.ErrorAuthorizationHeaderMissing

	mockStore := NewMockStore()
	mockStore.MakePutRequest(input1.Key, input1.Value, input1.Owner)

	rr := createMockPutRequestWithAuthenticator(mockStore, input1.Key, NewMockAuthenticatorWithValue(""), "")

	AssertErrorHttpCode(err, rr.Code, t)
}

func TestPutStatusUnprocessableEntityIfNoBodyIsPassed(t *testing.T) {

	err := common.ErrorStoreValueNotSet

	mockStore := NewMockStore()
	mockStore.MakePutRequest(input1.Key, input1.Value, input1.Owner)

	rr := createMockPutRequestWithUsername(mockStore, input1.Key, input1.Owner, "")

	AssertErrorHttpCode(err, rr.Code, t)
}

func TestPutReturnsFobiddenIfNotOwnerAttemptsToUpdateKey(t *testing.T) {

	mockStore := NewMockStore()
	mockStore.MakePutRequest(input1.Key, input1.Value, input1.Owner)
	rr := createMockPutRequestWithUsername(mockStore, input1.Key, input2.Owner, input1.Value)

	AssertErrorHttpCode(common.ErrorUnauthorisedOwner, rr.Code, t)
}
