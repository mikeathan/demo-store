package endpoints_test

import (
	"demo-store/common"
	"demo-store/endpoints"
	"demo-store/store"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createMockDeleteRequestWithAuthenticator(store store.Store, url string, authenticator endpoints.Authenticator) *httptest.ResponseRecorder {
	path := "store"

	route := CreateMockRouteWithDelete(path, &MockTracer{}, store, authenticator)
	req, err := http.NewRequest(http.MethodDelete, path+url, nil)

	if err != nil {
		return nil
	}

	rr := httptest.NewRecorder()
	route.ServeHTTP(rr, req)

	return rr
}

func createMockDeleteRequestWithUsername(store store.Store, url string, username string) *httptest.ResponseRecorder {

	return createMockDeleteRequestWithAuthenticator(store, url, NewMockAuthenticator(username))
}

func TestDeleteRemovesKeyFromStore(t *testing.T) {
	mockStore := NewMockStore()
	// add data
	mockStore.MakePutRequest(input1.Key, input1.Value, input1.Owner)

	rr := createMockDeleteRequestWithUsername(mockStore, input1.Key, input1.Owner)
	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("handler returned unexpected code: got %v want %v", rr.Code, expected)
	}

	_, err := mockStore.MakeListRequest(input1.Key)
	expectedError := common.ErrorKeyNotFound
	if err != expectedError {
		t.Errorf("handler returned unexpected code: got %v want %v", err.Error(), expectedError)
	}
}

func TestDeleteReturnsErrorIfKeyIsMissing(t *testing.T) {
	mockStore := NewMockStore()
	err := common.ErrorAuthorizationHeaderMissing
	rr := createMockDeleteRequestWithAuthenticator(mockStore, "", NewMockAuthenticatorWithError(err))

	AssertErrorHttpCode(err, rr.Code, t)
}

func TestDeleteAtuhenticatorReturnsError(t *testing.T) {
	mockStore := NewMockStore()

	err := common.ErrorAuthorizationHeaderMissing
	rr := createMockDeleteRequestWithAuthenticator(mockStore, input1.Key, NewMockAuthenticatorWithError(err))

	AssertErrorHttpCode(err, rr.Code, t)
}

func TestDeleteReturnsErrorIfOwnerIsMissing(t *testing.T) {
	mockStore := NewMockStore()

	err := common.ErrorAuthorizationHeaderMissing
	rr := createMockDeleteRequestWithAuthenticator(mockStore, input1.Key, NewMockAuthenticatorWithValue(""))

	AssertErrorHttpCode(err, rr.Code, t)
}

func TestDeleteReturnsErrorIfKeyNotFound(t *testing.T) {
	mockStore := NewMockStore()

	rr := createMockDeleteRequestWithUsername(mockStore, "", input1.Owner)

	AssertErrorHttpCode(common.ErrorKeyNotSet, rr.Code, t)
}

func TestDeleteReturnsForbiddenIfOwnerIsNotTheCreator(t *testing.T) {

	mockStore := NewMockStore()
	// add data
	mockStore.MakePutRequest(input1.Key, input1.Value, input1.Owner)

	rr := createMockDeleteRequestWithUsername(mockStore, input1.Key, input2.Owner)
	AssertErrorHttpCode(common.ErrorUnauthorisedOwner, rr.Code, t)
}
