package endpoints

import (
	"demo-store/common"
	"demo-store/store"
	"demo-store/utils"
	"net/http"
)

type ShutdownHandler struct {
	Tracer     utils.Tracer
	httpMethod string
	store      store.Store
}

func (p *ShutdownHandler) HttpMethod() string {
	return p.httpMethod
}

func (p *ShutdownHandler) Handle(args *HttpMethodHandlerParams, resp http.ResponseWriter, req *http.Request) HttpResult {
	return p.handleRequest(args, req, resp)
}

func (p *ShutdownHandler) handleRequest(args *HttpMethodHandlerParams, req *http.Request, resp http.ResponseWriter) HttpResult {
	username := args.Get(UsernameParameter)
	if username == "" {
		return CreateHttpResponseFromError(common.ErrorAuthorizationHeaderMissing)
	}

	if !p.store.UserDatabase().IsAdmin(username) {
		return CreateHttpResponseFromError(common.ErrorUnauthorisedOwner)
	}

	p.store.MakeShutdownRequest()
	return CreateHttpResponse("Ok", http.StatusOK)
}
