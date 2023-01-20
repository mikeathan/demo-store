package endpoints

import (
	"demo-store/utils"
	"net/http"
)

type PingHandler struct {
	Tracer     utils.Tracer
	httpMethod string
}

func (p *PingHandler) HttpMethod() string {
	return p.httpMethod
}
func (p *PingHandler) Handle(args *HttpMethodHandlerParams, resp http.ResponseWriter, req *http.Request) HttpResult {
	return p.handleRequest(args, req, resp)
}

func (p *PingHandler) handleRequest(args *HttpMethodHandlerParams, req *http.Request, resp http.ResponseWriter) HttpResult {

	resp.Header().Set("Content-Type", "text/plain")
	resp.Write([]byte("pong"))
	return CreateHttpResponse("Ok", http.StatusOK)
}
