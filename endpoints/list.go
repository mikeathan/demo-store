package endpoints

import (
	"demo-store/common"
	"demo-store/store"
	"demo-store/utils"
	"net/http"
	"strings"
)

type ListHandler struct {
	Tracer     utils.Tracer
	httpMethod string
	store      store.Store
}

func (p *ListHandler) HttpMethod() string {
	return p.httpMethod
}

func (p *ListHandler) Handle(args *HttpMethodHandlerParams, resp http.ResponseWriter, req *http.Request) HttpResult {
	return p.handleRequest(args, req, resp)
}

func (p *ListHandler) handleRequest(args *HttpMethodHandlerParams, req *http.Request, resp http.ResponseWriter) HttpResult {

	username := args.Get(UsernameParameter)
	if username == "" {
		return CreateHttpResponseFromError(common.ErrorAuthorizationHeaderMissing)
	}

	path := args.Get(PathParameter)
	key := strings.Trim(req.URL.Path, path)
	if key != "" {
		return p.handleFindRequest(key, resp)
	}

	return p.handleFindAllRequest(resp)
}

func (p *ListHandler) handleFindAllRequest(resp http.ResponseWriter) HttpResult {
	entries := p.store.MakeListAllRequest()
	err := writeResponse(entries, resp)
	if err != nil {
		return CreateHttpResponseFromError(err)
	}

	return CreateHttpResponse("Ok", http.StatusOK)
}

func (p *ListHandler) handleFindRequest(key string, resp http.ResponseWriter) HttpResult {
	entry, err := p.store.MakeListRequest(key)
	if err != nil {
		return CreateHttpResponseFromError(err)
	}

	err = writeResponse(entry, resp)
	if err != nil {
		return CreateHttpResponseFromError(err)
	}

	return CreateHttpResponse("Ok", http.StatusOK)
}

func writeResponse(entries any, resp http.ResponseWriter) error {
	json, err := common.ToJson(entries)
	if err != nil {
		return err
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(json))

	return nil
}
