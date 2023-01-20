package endpoints

import (
	"demo-store/utils"
	"fmt"
	"net/http"
)

type Route interface {
	ServeHTTP(resp http.ResponseWriter, req *http.Request)
	RootPath() string
}

type SecureRoute struct {
	Tracer         utils.Tracer
	Path           string
	MethodHandlers []HttpMethodHandler
	Authenticator  Authenticator
}

type InsecureRoute struct {
	Tracer         utils.Tracer
	Path           string
	MethodHandlers []HttpMethodHandler
}

func (p *InsecureRoute) RootPath() string {
	return p.Path
}

func (p *InsecureRoute) log(r *http.Request) {
	p.Tracer.LogInfo(fmt.Sprintf("source: %v method: %s URL: %s", r.RemoteAddr, r.Method, r.URL.Path))
}

func (p *InsecureRoute) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	for _, methodHandler := range p.MethodHandlers {
		switch req.Method {
		case methodHandler.HttpMethod():

			p.log(req)

			httpResp := methodHandler.Handle(CreatePathParameter(p.Path), resp, req)
			if httpResp.Code != http.StatusOK {

				p.Tracer.LogError(fmt.Sprintf("Message: %s Code: %d", httpResp.Message, httpResp.Code))
				http.Error(resp, httpResp.Message, httpResp.Code)
			} else {
				p.Tracer.LogInfo(fmt.Sprintf("Message: %s Code: %d", httpResp.Message, httpResp.Code))
			}
			return
		}
	}

	http.Error(resp, "Method not allowed", http.StatusMethodNotAllowed)
}

func (p *SecureRoute) RootPath() string {
	return p.Path
}

func (p *SecureRoute) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	for _, methodHandler := range p.MethodHandlers {
		switch req.Method {

		case methodHandler.HttpMethod():
			p.log(req)

			var httpResp HttpResult
			username, err := p.Authenticator.GetUsername(req.Header.Get("Authorization"))
			if err != nil {
				httpResp = CreateHttpResponseFromError(err)
			} else {
				httpResp = methodHandler.Handle(CreatePathAndUsernameParameter(p.Path, username), resp, req)
			}

			if httpResp.Code != http.StatusOK {

				p.Tracer.LogError(fmt.Sprintf("Message: %s Code: %d", httpResp.Message, httpResp.Code))
				http.Error(resp, httpResp.Message, httpResp.Code)
			} else {
				p.Tracer.LogInfo(fmt.Sprintf("Message: %s Code: %d", httpResp.Message, httpResp.Code))
			}
			return
		}
	}

	http.Error(resp, "Method not allowed", http.StatusMethodNotAllowed)
}

func (p *SecureRoute) log(r *http.Request) {
	p.Tracer.LogInfo(fmt.Sprintf("source: %v method: %s URL: %s", r.RemoteAddr, r.Method, r.URL.Path))
}
