package endpoints_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func createMockPingRequest() *httptest.ResponseRecorder {
	path := "ping"
	route := CreateMockRouteWithPing(path, &MockTracer{})
	req, err := http.NewRequest(http.MethodGet, path, nil)

	if err != nil {
		return nil
	}

	rr := httptest.NewRecorder()
	route.ServeHTTP(rr, req)

	return rr
}

func TestPingReturnsPong(t *testing.T) {

	rr := createMockPingRequest()

	expected := "pong"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestPingReturnsSuccess(t *testing.T) {
	rr := createMockPingRequest()

	expected := http.StatusOK
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
