package endpoints

import (
	"demo-store/common"
	"demo-store/store"
	"demo-store/utils"
	"net/http"
	"strings"
)

type GetHandler struct {
	Tracer     utils.Tracer
	httpMethod string
	store      store.Store
}

func (p *GetHandler) HttpMethod() string {
	return p.httpMethod
}

func (p *GetHandler) Handle(args *HttpMethodHandlerParams, resp http.ResponseWriter, req *http.Request) HttpResult {
	return p.handleRequest(args, req, resp)
}

func (p *GetHandler) handleRequest(args *HttpMethodHandlerParams, req *http.Request, resp http.ResponseWriter) HttpResult {

	path := args.Get(PathParameter)
	key := strings.Trim(req.URL.Path, path)

	if key == "" {
		return CreateHttpResponseFromError(common.ErrorKeyNotSet)
	}

	username := args.Get(UsernameParameter)
	if username == "" {
		return CreateHttpResponseFromError(common.ErrorAuthorizationHeaderMissing)
	}

	value, err := p.store.MakeGetRequest(key)
	if err != nil {
		return CreateHttpResponseFromError(err)
	}

	resp.Header().Set("Content-Type", "text/plain")
	resp.Write([]byte(value))

	return CreateHttpResponse("Ok", http.StatusOK)
}
