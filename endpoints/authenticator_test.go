package endpoints_test

import (
	"demo-store/common"
	"demo-store/endpoints"
	"errors"
	"testing"
)

func TestCtor(t *testing.T) {

	auth := endpoints.NewRouteAuthenticator(CreateMockTracer())
	if auth == nil {
		t.Errorf("Error returned unexpected object: got %v want %v", auth, "not nil")
	}
}

func TestTokenReturnsErrorAuthorizationHeaderMissing(t *testing.T) {
	auth := endpoints.NewRouteAuthenticator(CreateMockTracer())

	_, err := auth.GetUsername("")

	if err != common.ErrorAuthorizationHeaderMissing {
		t.Errorf("Returned unexpected error: got %v want %v", err, common.ErrorAuthorizationHeaderMissing)
	}
}

func TestTokenReturnsUsername(t *testing.T) {

	excpected := "someusername"
	auth := endpoints.NewRouteAuthenticator(CreateMockTracer())
	val, _ := auth.GetUsername(excpected)

	if val != excpected {
		t.Errorf("Returned unexpected username: got %v want %v", val, excpected)
	}
}
func TestAuthenticatorUserGetsExtractedFromBearerToken(t *testing.T) {

	token := "Bearer 1234567890"
	expected := "user"

	auth := endpoints.NewRouteAuthenticatorWithTokenizer(CreateMockTracer(), NewMockTokenizer(expected, nil))
	val, _ := auth.GetUsername(token)

	if val != expected {
		t.Errorf("Returned unexpected username: got %v want %v", val, expected)
	}
}

func TestAuthenticatorUserIsReturnedIfHeaderIsNotBearerToken(t *testing.T) {

	token := "user"
	expected := "user"

	auth := endpoints.NewRouteAuthenticatorWithTokenizer(CreateMockTracer(), NewMockTokenizer(expected, nil))
	val, _ := auth.GetUsername(token)

	if val != expected {
		t.Errorf("Returned unexpected username: got %v want %v", val, expected)
	}
}

func TestAuthenticatorTokenizerErrorIsHandled(t *testing.T) {

	token := "Bearer 123456"
	expected := errors.New("some error")

	auth := endpoints.NewRouteAuthenticatorWithTokenizer(CreateMockTracer(), NewMockTokenizer(token, errors.New("some error")))
	_, err := auth.GetUsername(token)

	if errors.Is(err, expected) {
		t.Errorf("Returned unexpected error: got %v want %v", err, expected)
	}
}
