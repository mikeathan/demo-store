package endpoints_test

import (
	"demo-store/common"
	"demo-store/endpoints"
	"demo-store/store"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createMockListRequestWithAuthenticator(store store.Store, url string, authenticator endpoints.Authenticator) *httptest.ResponseRecorder {
	path := "store"

	route := CreateMockRouteWithList(path, &MockTracer{}, store, authenticator)
	req, err := http.NewRequest(http.MethodGet, path+url, nil)

	if err != nil {
		return nil
	}

	rr := httptest.NewRecorder()
	route.ServeHTTP(rr, req)

	return rr
}

func createMockListRequestWithUsername(store store.Store, url string, username string) *httptest.ResponseRecorder {

	return createMockListRequestWithAuthenticator(store, url, NewMockAuthenticator(username))
}

func TestListAllReturnsSuccess(t *testing.T) {

	mockStore := NewMockStore()
	mockStore.MakePutRequest(input1.Key, input1.Value, input1.Owner)
	rr := createMockListRequestWithUsername(mockStore, "", input1.Owner)

	expectedStatus := http.StatusOK
	if rr.Code != expectedStatus {
		t.Errorf("handler returned unexpected code: got %v want %v", rr.Code, expectedStatus)
	}

	expectedBody, _ := common.ToJson(mockStore.MakeListAllRequest())
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
	}
}

func TestListReturnsFailed(t *testing.T) {

	mockStore := NewMockStore()
	mockStore.MakePutRequest(input1.Key, input1.Value, input1.Owner)
	rr := createMockListRequestWithUsername(mockStore, input2.Value, input1.Owner)

	AssertErrorHttpCode(common.ErrorKeyNotFound, rr.Code, t)
}

func TestListAllReturnsFailed(t *testing.T) {

	mockStore := NewMockStore()
	mockStore.MakePutRequest(input1.Key, input1.Value, input1.Owner)
	rr := createMockListRequestWithUsername(mockStore, "", "")

	AssertErrorHttpCode(common.ErrorUnauthorisedOwner, rr.Code, t)
}

func TestListReturnsSuccess(t *testing.T) {

	mockStore := NewMockStore()
	mockStore.MakePutRequest(input1.Key, input1.Value, input1.Owner)
	rr := createMockListRequestWithUsername(mockStore, input1.Key, input1.Owner)

	expectedStatus := http.StatusOK
	if rr.Code != expectedStatus {
		t.Errorf("handler returned unexpected code: got %v want %v", rr.Code, expectedStatus)
	}

	expectedEntry, _ := mockStore.MakeListRequest(input1.Key)
	expectedBody, _ := common.ToJson(expectedEntry)
	if rr.Body.String() != expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
	}

	expectedContentType := "application/json"
	if rr.Header().Get("Content-Type") != expectedContentType {
		t.Errorf("handler returned unexpected code: got %v want %v", rr.Header().Get("Content-Type"), expectedContentType)
	}
}

func TestListReturnsStatusNotFoundIfKeyNotFound(t *testing.T) {

	err := common.ErrorKeyNotFound

	mockStore := NewMockStore()
	mockStore.MakePutRequest(input1.Key, input1.Value, input1.Owner)

	rr := createMockListRequestWithAuthenticator(mockStore, input2.Key, NewMockAuthenticatorWithError(err))

	AssertErrorHttpCode(err, rr.Code, t)
}

func TestListReturnsErrorIfAuthenticationMissing(t *testing.T) {

	err := common.ErrorAuthorizationHeaderMissing

	mockStore := NewMockStore()
	mockStore.MakePutRequest(input1.Key, input1.Value, input1.Owner)
	rr := createMockListRequestWithAuthenticator(mockStore, input1.Key, NewMockAuthenticatorWithValue(""))

	AssertErrorHttpCode(err, rr.Code, t)
}
