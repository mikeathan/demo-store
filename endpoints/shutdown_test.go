package endpoints_test

import (
	"demo-store/common"
	"demo-store/endpoints"
	"demo-store/store"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createMockShutdownequestWithAuthenticator(store store.Store, authenticator endpoints.Authenticator) *httptest.ResponseRecorder {

	path := "shutdown"

	route := CreateMockRouteWithShutdown(path, &MockTracer{}, store, authenticator)
	req, err := http.NewRequest(http.MethodGet, path, nil)

	if err != nil {
		return nil
	}

	rr := httptest.NewRecorder()
	route.ServeHTTP(rr, req)

	return rr
}

func createMockShutdownRequestWithUsername(store store.Store, username string) *httptest.ResponseRecorder {

	return createMockShutdownequestWithAuthenticator(store, NewMockAuthenticatorWithValue(username))
}

func TestShutdownReturnsForbiddenForNonAdminUser(t *testing.T) {
	err := common.ErrorUnauthorisedOwner

	mockStore := NewMockStore()
	rr := createMockShutdownRequestWithUsername(mockStore, input1.Owner)

	AssertErrorHttpCode(err, rr.Code, t)
}

func TestShutdownReturnsErrorIfOwnerIsMissing(t *testing.T) {
	err := common.ErrorAuthorizationHeaderMissing
	mockStore := NewMockStore()

	rr := createMockShutdownequestWithAuthenticator(mockStore, NewMockAuthenticatorWithValue(""))

	AssertErrorHttpCode(err, rr.Code, t)
}

func TestShutdownReturnsSuccessIfOwnerIsAdmin(t *testing.T) {

	mockStore := NewMockStore()
	mockStore.UserDatabase().AddUser("admin", "123")
	rr := createMockShutdownRequestWithUsername(mockStore, "admin")

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("handler returned unexpected code: got %v want %v", rr.Code, expected)
	}
}
