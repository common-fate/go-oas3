// This file is generated by github.com/common-fate/go-oas3. DO NOT EDIT.

package output

import (
	"context"
	"encoding/json"
	chi "github.com/go-chi/chi/v5"
	"net/http"
)

type Hooks struct {
	RequestSecurityParseFailed    func(*http.Request, string, RequestProcessingResult)
	RequestSecurityParseCompleted func(*http.Request, string)
	RequestSecurityCheckFailed    func(*http.Request, string, string, RequestProcessingResult)
	RequestSecurityCheckCompleted func(*http.Request, string, string)
	RequestBodyUnmarshalFailed    func(*http.Request, string, RequestProcessingResult)
	RequestHeaderParseFailed      func(*http.Request, string, string, RequestProcessingResult)
	RequestPathParseFailed        func(*http.Request, string, string, RequestProcessingResult)
	RequestQueryParseFailed       func(*http.Request, string, string, RequestProcessingResult)
	RequestBodyValidationFailed   func(*http.Request, string, RequestProcessingResult)
	RequestHeaderValidationFailed func(*http.Request, string, RequestProcessingResult)
	RequestPathValidationFailed   func(*http.Request, string, RequestProcessingResult)
	RequestQueryValidationFailed  func(*http.Request, string, RequestProcessingResult)
	RequestBodyUnmarshalCompleted func(*http.Request, string)
	RequestHeaderParseCompleted   func(*http.Request, string)
	RequestPathParseCompleted     func(*http.Request, string)
	RequestQueryParseCompleted    func(*http.Request, string)
	RequestParseCompleted         func(*http.Request, string)
	RequestProcessingCompleted    func(*http.Request, string)
	RequestRedirectStarted        func(*http.Request, string, string)
	ResponseBodyMarshalCompleted  func(*http.Request, string)
	ResponseBodyWriteCompleted    func(*http.Request, string, int)
	ResponseBodyMarshalFailed     func(http.ResponseWriter, *http.Request, string, error)
	ResponseBodyWriteFailed       func(*http.Request, string, int, error)
	ServiceCompleted              func(*http.Request, string)
}

type requestProcessingResultType uint8

const (
	BodyUnmarshalFailed requestProcessingResultType = iota + 1
	BodyValidationFailed
	HeaderParseFailed
	HeaderValidationFailed
	QueryParseFailed
	QueryValidationFailed
	PathParseFailed
	PathValidationFailed
	SecurityParseFailed
	SecurityCheckFailed
	ParseSucceed
)

type RequestProcessingResult struct {
	error error
	typee requestProcessingResultType
}

func NewRequestProcessingResult(t requestProcessingResultType, err error) RequestProcessingResult {
	return RequestProcessingResult{
		error: err,
		typee: t,
	}
}

func (r RequestProcessingResult) Type() requestProcessingResultType {
	return r.typee
}

func (r RequestProcessingResult) Err() error {
	return r.error
}

func CarsHandler(impl CarsService, r chi.Router, hooks *Hooks) http.Handler {
	router := &carsRouter{router: r, service: impl, hooks: hooks}

	router.mount()

	return router.router
}

type carsRouter struct {
	router  chi.Router
	service CarsService
	hooks   *Hooks
}

func (router *carsRouter) mount() {
	router.router.Post("/cars", router.PostCars)
}

func (router *carsRouter) parsePostCarsRequest(r *http.Request) (request PostCarsRequest) {
	request.ProcessingResult = RequestProcessingResult{typee: ParseSucceed}

	var (
		body      PostCarsRequestBody
		decodeErr error
	)
	decodeErr = json.NewDecoder(r.Body).Decode(&body)
	if decodeErr != nil {
		request.ProcessingResult = RequestProcessingResult{error: decodeErr, typee: BodyUnmarshalFailed}
		if router.hooks.RequestBodyUnmarshalFailed != nil {
			router.hooks.RequestBodyUnmarshalFailed(r, "PostCars", request.ProcessingResult)

			return
		}

		return
	}

	request.Body = body

	if err := request.Body.Validate(); err != nil {
		request.ProcessingResult = RequestProcessingResult{error: err, typee: BodyValidationFailed}
		if router.hooks.RequestBodyValidationFailed != nil {
			router.hooks.RequestBodyValidationFailed(r, "PostCars", request.ProcessingResult)
		}

		return
	}

	if router.hooks.RequestBodyUnmarshalCompleted != nil {
		router.hooks.RequestBodyUnmarshalCompleted(r, "PostCars")
	}

	if router.hooks.RequestParseCompleted != nil {
		router.hooks.RequestParseCompleted(r, "PostCars")
	}

	return
}

func (router *carsRouter) PostCars(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	response := router.service.PostCars(r.Context(), router.parsePostCarsRequest(r))

	if response.statusCode() == 302 && response.redirectURL() != "" {
		if router.hooks.RequestRedirectStarted != nil {
			router.hooks.RequestRedirectStarted(r, "PostCars", response.redirectURL())
		}

		http.Redirect(w, r, response.redirectURL(), 302)

		if router.hooks.ServiceCompleted != nil {
			router.hooks.ServiceCompleted(r, "PostCars")
		}

		return
	}

	for header, value := range response.headers() {
		w.Header().Set(header, value)
	}

	for _, c := range response.cookies() {
		cookie := c
		http.SetCookie(w, &cookie)
	}

	if router.hooks.RequestProcessingCompleted != nil {
		router.hooks.RequestProcessingCompleted(r, "PostCars")
	}

	w.WriteHeader(response.statusCode())

	if router.hooks.ServiceCompleted != nil {
		router.hooks.ServiceCompleted(r, "PostCars")
	}
}

type response struct {
	statusCode  int
	body        interface{}
	contentType string
	redirectURL string
	headers     map[string]string
	cookies     []http.Cookie
}

type responseInterface interface {
	statusCode() int
	body() interface{}
	contentType() string
	redirectURL() string
	cookies() []http.Cookie
	headers() map[string]string
}

type PostCarsResponse interface {
	responseInterface
	postCarsResponse()
}

type postCarsResponse struct {
	response
}

func (postCarsResponse) postCarsResponse() {}

func (response postCarsResponse) statusCode() int {
	return response.response.statusCode
}

func (response postCarsResponse) body() interface{} {
	return response.response.body
}

func (response postCarsResponse) contentType() string {
	return response.response.contentType
}

func (response postCarsResponse) redirectURL() string {
	return response.response.redirectURL
}

func (response postCarsResponse) headers() map[string]string {
	return response.response.headers
}

func (response postCarsResponse) cookies() []http.Cookie {
	return response.response.cookies
}

type postCarsStatusCodeResponseBuilder struct {
	response
}

func PostCarsResponseBuilder() *postCarsStatusCodeResponseBuilder {
	return new(postCarsStatusCodeResponseBuilder)
}

func (builder *postCarsStatusCodeResponseBuilder) StatusCode200() *PostCars200ResponseBuilder {
	builder.response.statusCode = 200

	return &PostCars200ResponseBuilder{response: builder.response}
}

type PostCars200ResponseBuilder struct {
	response
}

func (builder *PostCars200ResponseBuilder) Build() PostCarsResponse {
	return postCarsResponse{response: builder.response}
}

type CarsService interface {
	PostCars(context.Context, PostCarsRequest) PostCarsResponse
}

type PostCarsRequest struct {
	Body             PostCarsRequestBody
	ProcessingResult RequestProcessingResult
}

type SecurityScheme string

const ()

type securityProcessor struct {
	scheme  SecurityScheme
	extract func(r *http.Request) (string, string, bool)
	handle  func(r *http.Request, scheme SecurityScheme, name string, value string) error
}

var securityExtractorsFuncs = map[SecurityScheme]func(r *http.Request) (string, string, bool){}

type SecuritySchemas interface{}

type SecurityCheckResult struct {
	Scheme SecurityScheme
	Value  string
}
