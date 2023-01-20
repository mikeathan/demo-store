package endpoints_test

import (
	"demo-store/common"
	"demo-store/endpoints"
	"demo-store/store"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createMockGetRequestWithAuthenticator(store store.Store, url string, authenticator endpoints.Authenticator) *httptest.ResponseRecorder {
	path := "store"

	route := CreateMockRouteWithGet(path, &MockTracer{}, store, authenticator)
	req, err := http.NewRequest(http.MethodGet, path+url, nil)

	if err != nil {
		return nil
	}

	rr := httptest.NewRecorder()
	route.ServeHTTP(rr, req)

	return rr
}

func createMockGetRequestWithUsername(store store.Store, url string, username string) *httptest.ResponseRecorder {

	return createMockGetRequestWithAuthenticator(store, url, NewMockAuthenticator(username))
}

func TestGetReturnsErrorIfKeyIsMissing(t *testing.T) {
	mockStore := NewMockStore()
	rr := createMockGetRequestWithUsername(mockStore, "", input1.Owner)

	AssertErrorHttpCode(common.ErrorKeyNotSet, rr.Code, t)
}

func TestGetReturnsErrorIfOwnerIsMissing(t *testing.T) {
	mockStore := NewMockStore()
	rr := createMockGetRequestWithAuthenticator(mockStore, input1.Key, NewMockAuthenticatorWithValue(""))

	AssertErrorHttpCode(common.ErrorAuthorizationHeaderMissing, rr.Code, t)
}

func TestGetReturnsErrorIfKeyNotFound(t *testing.T) {
	mockStore := NewMockStore()
	rr := createMockGetRequestWithUsername(mockStore, input1.Key, input1.Owner)

	AssertErrorHttpCode(common.ErrorKeyNotFound, rr.Code, t)
}

func TestGetReturnsSuccess(t *testing.T) {

	mockStore := NewMockStore()
	mockStore.MakePutRequest(input1.Key, input1.Value, input1.Owner)
	rr := createMockGetRequestWithUsername(mockStore, input1.Key, input1.Owner)

	expectedCode := http.StatusOK
	if rr.Code != expectedCode {
		t.Errorf("handler returned unexpected code: got %v want %v", rr.Code, expectedCode)

		expectedBody := input1.Value
		if rr.Body.String() != expectedBody {
			t.Errorf("handler returned unexpected code: got %v want %v", rr.Body.String(), expectedBody)
		}

		expectedContentType := "text/plain"
		if rr.Header().Get("Content-Type") != expectedContentType {
			t.Errorf("handler returned unexpected code: got %v want %v", rr.Header().Get("Content-Type"), expectedContentType)
		}
	}
}
