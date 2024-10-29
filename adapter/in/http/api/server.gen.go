// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
	strictgin "github.com/oapi-codegen/runtime/strictmiddleware/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Authenticate
	// (POST /auth/login)
	AuthenticateLogin(c *gin.Context)
	// Returns all organisations
	// (GET /organisations)
	GetAllOrganisations(c *gin.Context, params GetAllOrganisationsParams)
	// Create an organisation
	// (POST /organisations)
	CreateOrganisation(c *gin.Context)
	// Listing dummy stuff
	// (GET /test)
	ListTags(c *gin.Context)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// AuthenticateLogin operation middleware
func (siw *ServerInterfaceWrapper) AuthenticateLogin(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.AuthenticateLogin(c)
}

// GetAllOrganisations operation middleware
func (siw *ServerInterfaceWrapper) GetAllOrganisations(c *gin.Context) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAllOrganisationsParams

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", c.Request.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter limit: %w", err), http.StatusBadRequest)
		return
	}

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", c.Request.URL.Query(), &params.Page)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter page: %w", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetAllOrganisations(c, params)
}

// CreateOrganisation operation middleware
func (siw *ServerInterfaceWrapper) CreateOrganisation(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CreateOrganisation(c)
}

// ListTags operation middleware
func (siw *ServerInterfaceWrapper) ListTags(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.ListTags(c)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.POST(options.BaseURL+"/auth/login", wrapper.AuthenticateLogin)
	router.GET(options.BaseURL+"/organisations", wrapper.GetAllOrganisations)
	router.POST(options.BaseURL+"/organisations", wrapper.CreateOrganisation)
	router.GET(options.BaseURL+"/test", wrapper.ListTags)
}

type AuthenticateLoginRequestObject struct {
	Body *AuthenticateLoginJSONRequestBody
}

type AuthenticateLoginResponseObject interface {
	VisitAuthenticateLoginResponse(w http.ResponseWriter) error
}

type AuthenticateLogin200JSONResponse string

func (response AuthenticateLogin200JSONResponse) VisitAuthenticateLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type AuthenticateLogindefaultResponseHeaders struct {
	XCORRELATIONID string
}

type AuthenticateLogindefaultJSONResponse struct {
	Body       Error
	Headers    AuthenticateLogindefaultResponseHeaders
	StatusCode int
}

func (response AuthenticateLogindefaultJSONResponse) VisitAuthenticateLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-CORRELATION-ID", fmt.Sprint(response.Headers.XCORRELATIONID))
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetAllOrganisationsRequestObject struct {
	Params GetAllOrganisationsParams
}

type GetAllOrganisationsResponseObject interface {
	VisitGetAllOrganisationsResponse(w http.ResponseWriter) error
}

type GetAllOrganisations200ResponseHeaders struct {
	XNEXT     string
	XPREVIOUS string
}

type GetAllOrganisations200JSONResponse struct {
	Body    GetAllOrganisations
	Headers GetAllOrganisations200ResponseHeaders
}

func (response GetAllOrganisations200JSONResponse) VisitGetAllOrganisationsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-NEXT", fmt.Sprint(response.Headers.XNEXT))
	w.Header().Set("X-PREVIOUS", fmt.Sprint(response.Headers.XPREVIOUS))
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetAllOrganisationsdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response GetAllOrganisationsdefaultJSONResponse) VisitGetAllOrganisationsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type CreateOrganisationRequestObject struct {
	Body *CreateOrganisationJSONRequestBody
}

type CreateOrganisationResponseObject interface {
	VisitCreateOrganisationResponse(w http.ResponseWriter) error
}

type CreateOrganisation201ResponseHeaders struct {
	ContentLocation string
	XCORRELATIONID  string
	XOBJECTID       string
}

type CreateOrganisation201Response struct {
	Headers CreateOrganisation201ResponseHeaders
}

func (response CreateOrganisation201Response) VisitCreateOrganisationResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Location", fmt.Sprint(response.Headers.ContentLocation))
	w.Header().Set("X-CORRELATION-ID", fmt.Sprint(response.Headers.XCORRELATIONID))
	w.Header().Set("X-OBJECT-ID", fmt.Sprint(response.Headers.XOBJECTID))
	w.WriteHeader(201)
	return nil
}

type CreateOrganisationdefaultResponseHeaders struct {
	XCORRELATIONID string
}

type CreateOrganisationdefaultJSONResponse struct {
	Body       Error
	Headers    CreateOrganisationdefaultResponseHeaders
	StatusCode int
}

func (response CreateOrganisationdefaultJSONResponse) VisitCreateOrganisationResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-CORRELATION-ID", fmt.Sprint(response.Headers.XCORRELATIONID))
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type ListTagsRequestObject struct {
}

type ListTagsResponseObject interface {
	VisitListTagsResponse(w http.ResponseWriter) error
}

type ListTags200ResponseHeaders struct {
	XCORRELATIONID string
}

type ListTags200JSONResponse struct {
	Body    Dummy
	Headers ListTags200ResponseHeaders
}

func (response ListTags200JSONResponse) VisitListTagsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-CORRELATION-ID", fmt.Sprint(response.Headers.XCORRELATIONID))
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response.Body)
}

type ListTagsdefaultResponseHeaders struct {
	XCORRELATIONID string
}

type ListTagsdefaultJSONResponse struct {
	Body       Error
	Headers    ListTagsdefaultResponseHeaders
	StatusCode int
}

func (response ListTagsdefaultJSONResponse) VisitListTagsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-CORRELATION-ID", fmt.Sprint(response.Headers.XCORRELATIONID))
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Authenticate
	// (POST /auth/login)
	AuthenticateLogin(ctx context.Context, request AuthenticateLoginRequestObject) (AuthenticateLoginResponseObject, error)
	// Returns all organisations
	// (GET /organisations)
	GetAllOrganisations(ctx context.Context, request GetAllOrganisationsRequestObject) (GetAllOrganisationsResponseObject, error)
	// Create an organisation
	// (POST /organisations)
	CreateOrganisation(ctx context.Context, request CreateOrganisationRequestObject) (CreateOrganisationResponseObject, error)
	// Listing dummy stuff
	// (GET /test)
	ListTags(ctx context.Context, request ListTagsRequestObject) (ListTagsResponseObject, error)
}

type StrictHandlerFunc = strictgin.StrictGinHandlerFunc
type StrictMiddlewareFunc = strictgin.StrictGinMiddlewareFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// AuthenticateLogin operation middleware
func (sh *strictHandler) AuthenticateLogin(ctx *gin.Context) {
	var request AuthenticateLoginRequestObject

	var body AuthenticateLoginJSONRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.Error(err)
		return
	}
	request.Body = &body

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.AuthenticateLogin(ctx, request.(AuthenticateLoginRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "AuthenticateLogin")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(AuthenticateLoginResponseObject); ok {
		if err := validResponse.VisitAuthenticateLoginResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetAllOrganisations operation middleware
func (sh *strictHandler) GetAllOrganisations(ctx *gin.Context, params GetAllOrganisationsParams) {
	var request GetAllOrganisationsRequestObject

	request.Params = params

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetAllOrganisations(ctx, request.(GetAllOrganisationsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetAllOrganisations")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(GetAllOrganisationsResponseObject); ok {
		if err := validResponse.VisitGetAllOrganisationsResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// CreateOrganisation operation middleware
func (sh *strictHandler) CreateOrganisation(ctx *gin.Context) {
	var request CreateOrganisationRequestObject

	var body CreateOrganisationJSONRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.Error(err)
		return
	}
	request.Body = &body

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.CreateOrganisation(ctx, request.(CreateOrganisationRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreateOrganisation")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(CreateOrganisationResponseObject); ok {
		if err := validResponse.VisitCreateOrganisationResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// ListTags operation middleware
func (sh *strictHandler) ListTags(ctx *gin.Context) {
	var request ListTagsRequestObject

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.ListTags(ctx, request.(ListTagsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ListTags")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(ListTagsResponseObject); ok {
		if err := validResponse.VisitListTagsResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
