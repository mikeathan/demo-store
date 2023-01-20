package endpoints

import (
	"demo-store/common"
	"demo-store/store"
	"demo-store/utils"
	"net/http"
	"strings"
)

type DeleteHandler struct {
	Tracer     utils.Tracer
	httpMethod string
	store      store.Store
}

func (p *DeleteHandler) HttpMethod() string {
	return p.httpMethod
}

func (p *DeleteHandler) Handle(args *HttpMethodHandlerParams, resp http.ResponseWriter, req *http.Request) HttpResult {
	return p.handleRequest(args, req, resp)
}

func (p *DeleteHandler) handleRequest(args *HttpMethodHandlerParams, req *http.Request, resp http.ResponseWriter) HttpResult {

	path := args.Get(PathParameter)
	key := strings.Trim(req.URL.Path, path)

	if key == "" {
		return CreateHttpResponseFromError(common.ErrorKeyNotSet)
	}

	username := args.Get(UsernameParameter)
	if username == "" {
		return CreateHttpResponseFromError(common.ErrorAuthorizationHeaderMissing)
	}

	err := p.store.MakeDeleteRequest(key, username)
	if err != nil {
		return CreateHttpResponseFromError(err)
	}

	return CreateHttpResponse("Ok", http.StatusOK)
}
