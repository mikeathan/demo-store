package endpoints

import (
	"demo-store/common"
	"demo-store/users"
	"demo-store/utils"
	"fmt"
	"net/http"
)

type LoginHandler struct {
	Tracer     utils.Tracer
	httpMethod string
	Users      users.UserDatabase
	Tokenizer  utils.Tokenizer
}

func (p *LoginHandler) HttpMethod() string {
	return p.httpMethod
}

func (p *LoginHandler) Handle(args *HttpMethodHandlerParams, resp http.ResponseWriter, req *http.Request) HttpResult {
	return p.handleRequest(args, req, resp)
}

func (p *LoginHandler) handleRequest(args *HttpMethodHandlerParams, req *http.Request, resp http.ResponseWriter) HttpResult {
	username, password, ok := req.BasicAuth()
	if !ok {
		return CreateHttpResponseFromError(common.ErrorInvalidAuthorizationHeader)
	}
	if username == "" || password == "" {
		return CreateHttpResponseFromError(common.ErrorAuthorizationHeaderMissing)
	}

	err := p.Users.Authenticate(username, password)
	if err != nil {
		return CreateHttpResponseFromError(common.ErrorAuthorizationFailed)
	}

	tokenString, err := p.Tokenizer.CreateToken(username)
	if err != nil {
		return CreateHttpResponseFromError(common.ErrorAuthorizationFailed)
	}

	resp.Header().Set("Content-Type", "text/plain")
	resp.Write([]byte(fmt.Sprintf("Bearer %s", tokenString)))

	return CreateHttpResponse("Ok", http.StatusOK)
}
