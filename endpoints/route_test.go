package endpoints_test

import (
	"demo-store/endpoints"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInsecureRouteMethodNotAllowed(t *testing.T) {

	route := CreateMockRouteWithPing("somepath", &MockTracer{})
	req, _ := http.NewRequest(http.MethodPost, "somepath", nil)

	rr := httptest.NewRecorder()
	route.ServeHTTP(rr, req)

	expected := http.StatusMethodNotAllowed
	if rr.Code != expected {
		t.Errorf("handler returned unexpected code: got %v want %v", rr.Code, expected)
	}
}

func TestHttpResult(t *testing.T) {

	err := errors.New("Mock Error")
	resp := endpoints.CreateHttpResponseFromError(err)

	expected := http.StatusInternalServerError
	if resp.Code != expected {
		t.Errorf("handler returned unexpected code: got %v want %v", resp.Code, expected)
	}

}
func TestInsecureRoutePath(t *testing.T) {

	expected := "somepath"
	route := CreateMockRouteWithPing(expected, &MockTracer{})
	if route.RootPath() != expected {
		t.Errorf("handler returned unexpected path: got %v want %v", route.RootPath(), expected)
	}
}

func TestSecureRouteMethodNotAllowed(t *testing.T) {

	mockStore := NewMockStore()
	route := CreateMockRouteWithList("somepath", &MockTracer{}, mockStore, NewMockAuthenticator(""))
	req, _ := http.NewRequest(http.MethodPost, "somepath", nil)

	rr := httptest.NewRecorder()
	route.ServeHTTP(rr, req)

	expected := http.StatusMethodNotAllowed
	if rr.Code != expected {
		t.Errorf("handler returned unexpected code: got %v want %v", rr.Code, expected)
	}
}

func TestSecureRoutePath(t *testing.T) {
	expected := "somepath"
	mockStore := NewMockStore()
	route := CreateMockRouteWithList(expected, &MockTracer{}, mockStore, NewMockAuthenticator(""))
	if route.RootPath() != expected {
		t.Errorf("handler returned unexpected path: got %v want %v", route.RootPath(), expected)
	}
}

func TestCreateStoreRouteMethods(t *testing.T) {
	mockStore := NewMockStore()
	auth := NewMockAuthenticator("")
	route := endpoints.CreateStoreRoute(mockStore.Tracer, mockStore, auth)

	secureRoute := route.(*endpoints.SecureRoute)
	expectedPath := "/store/"
	if secureRoute.RootPath() != expectedPath {
		t.Errorf("handler returned unexpected path: got %v want %v", secureRoute.RootPath(), expectedPath)
	}
	expected := []string{
		"PUT",
		"GET",
		"DELETE",
	}
	for i, handlder := range secureRoute.MethodHandlers {
		if handlder.HttpMethod() != expected[i] {
			t.Errorf("handler returned unexpected method: got %v want %v", handlder.HttpMethod(), expected[i])
		}
	}
}

func TestCreatePingRouteMethods(t *testing.T) {
	mockStore := NewMockStore()
	route := endpoints.CreatePingRoute(mockStore.Tracer, mockStore)

	insecureRoute := route.(*endpoints.InsecureRoute)
	expectedPath := "/ping/"
	if insecureRoute.RootPath() != expectedPath {
		t.Errorf("handler returned unexpected path: got %v want %v", insecureRoute.RootPath(), expectedPath)
	}

	expected := []string{
		"GET",
	}
	for i, handlder := range insecureRoute.MethodHandlers {
		if handlder.HttpMethod() != expected[i] {
			t.Errorf("handler returned unexpected method: got %v want %v", handlder.HttpMethod(), expected[i])
		}
	}
}

func TestCreateListRouteMethods(t *testing.T) {
	mockStore := NewMockStore()
	auth := NewMockAuthenticator("")
	route := endpoints.CreateListRoute(mockStore.Tracer, mockStore, auth)

	secureRoute := route.(*endpoints.SecureRoute)
	expectedPath := "/list/"
	if secureRoute.RootPath() != expectedPath {
		t.Errorf("handler returned unexpected path: got %v want %v", secureRoute.RootPath(), expectedPath)
	}

	expected := []string{
		"GET",
	}
	for i, handlder := range secureRoute.MethodHandlers {

		if handlder.HttpMethod() != expected[i] {
			t.Errorf("handler returned unexpected method: got %v want %v", handlder.HttpMethod(), expected[i])
		}
	}
}

func TestCreateShutdownRouteMethods(t *testing.T) {
	mockStore := NewMockStore()
	auth := NewMockAuthenticator("")
	route := endpoints.CreateShutdownRoute(mockStore.Tracer, mockStore, auth)

	secureRoute := route.(*endpoints.SecureRoute)
	expectedPath := "/shutdown/"
	if secureRoute.RootPath() != expectedPath {
		t.Errorf("handler returned unexpected path: got %v want %v", secureRoute.RootPath(), expectedPath)
	}

	expected := []string{
		"GET",
	}
	for i, handlder := range secureRoute.MethodHandlers {

		if handlder.HttpMethod() != expected[i] {
			t.Errorf("handler returned unexpected method: got %v want %v", handlder.HttpMethod(), expected[i])
		}
	}
}

func TestCreateLoginRouteMethods(t *testing.T) {
	mockStore := NewMockStore()
	route := endpoints.CreateLoginRoute(mockStore.Tracer, mockStore.UserDatabase())

	insecureRoute := route.(*endpoints.InsecureRoute)
	expectedPath := "/login/"
	if insecureRoute.RootPath() != expectedPath {
		t.Errorf("handler returned unexpected path: got %v want %v", insecureRoute.RootPath(), expectedPath)
	}

	expected := []string{
		"GET",
	}
	for i, handlder := range insecureRoute.MethodHandlers {

		if handlder.HttpMethod() != expected[i] {
			t.Errorf("handler returned unexpected method: got %v want %v", handlder.HttpMethod(), expected[i])
		}
	}
}

func TestCreateInsecureApiRoutes(t *testing.T) {
	mockStore := NewMockStore()
	routes := endpoints.APIRoutes(mockStore.Tracer, mockStore)

	expectedPaths := []string{
		"/ping/",
		"/login/",
	}
	for i, route := range routes.Insecure {
		if route.RootPath() != expectedPaths[i] {
			t.Errorf("handler returned unexpected path: got %v want %v", route.RootPath(), expectedPaths[i])
		}
	}
}

func TestCreateSecureApiRoutes(t *testing.T) {
	mockStore := NewMockStore()
	routes := endpoints.APIRoutes(mockStore.Tracer, mockStore)

	expectedPaths := []string{
		"/store/",
		"/list/",
		"/shutdown/",
	}
	for i, route := range routes.Secure {
		if route.RootPath() != expectedPaths[i] {
			t.Errorf("handler returned unexpected path: got %v want %v", route.RootPath(), expectedPaths[i])
		}
	}
}
