package endpoints

import (
	"demo-store/common"
	"demo-store/users"
	"demo-store/utils"
	"strings"
)

type Authenticator interface {
	GetUsername(bearerToken string) (string, error)
}
type RouteAuthenticator struct {
	UserDatabase users.UserDatabase
	Tokenizer    utils.Tokenizer
}

func NewRouteAuthenticator(tracer utils.Tracer) Authenticator {
	return &RouteAuthenticator{Tokenizer: utils.NewJwtTokenizer(tracer)}
}

func NewRouteAuthenticatorWithTokenizer(tracer utils.Tracer, tokenizer utils.Tokenizer) Authenticator {
	return &RouteAuthenticator{Tokenizer: tokenizer}
}

func (p *RouteAuthenticator) GetUsername(bearerToken string) (string, error) {

	if bearerToken == "" {
		return "", common.ErrorAuthorizationHeaderMissing
	}

	if strings.HasPrefix(bearerToken, utils.BearerTokenHeader) {
		bearerToken = strings.ReplaceAll(bearerToken, utils.BearerTokenHeader, "")
		username, err := p.Tokenizer.GetUsernameFromToken(bearerToken)
		if err != nil {
			return "", err
		}

		return username, nil
	}

	return bearerToken, nil
}
