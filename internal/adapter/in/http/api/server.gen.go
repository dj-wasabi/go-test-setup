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
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"
	"github.com/oapi-codegen/runtime"
	strictgin "github.com/oapi-codegen/runtime/strictmiddleware/gin"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// AuthenticatePostRequest Authenticate
type AuthenticatePostRequest struct {
	Password string `json:"password" validate:"required"`
	Username string `json:"username" validate:"required,min=6,max=64"`
}

// AuthenticatePostResponse Authenticate
type AuthenticatePostResponse struct {
	Token string `json:"token" validate:"required"`
}

// Error Default error response
type Error struct {
	Message string `json:"message"`
}

// MetricsGetResponse Metrics for Prometheus response
type MetricsGetResponse = string

// Organisation Return of the Organisation object
type Organisation struct {
	Admins      []string  `json:"admins"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description" validate:"required,min=1,max=256"`
	Enabled     bool      `json:"enabled"`
	Fqdn        string    `json:"fqdn"`
	Id          string    `json:"id" validate:"required,mongodb"`
	Name        string    `json:"name" validate:"required,min=6,max=64"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OrganisationIn Creating the Organisation object
type OrganisationIn struct {
	Admins      []string `json:"admins"`
	Description string   `json:"description" validate:"required,min=1,max=256"`
	Enabled     bool     `json:"enabled"`
	Fqdn        string   `json:"fqdn" validate:"required,fqdn"`
	Name        string   `json:"name" validate:"required,min=6,max=64"`
}

// Organisations Overview of all organisations
type Organisations = []Organisation

// UserIn Object to create the User
type UserIn struct {
	Enabled  *bool   `json:"enabled,omitempty"`
	OrgId    string  `json:"org_id"`
	Password string  `json:"password"`
	Role     *string `json:"role,omitempty"`
	Username string  `json:"username"`
}

// UserNoPassword Return of the User object
type UserNoPassword struct {
	CreatedAt time.Time `json:"created_at"`
	Enabled   bool      `json:"enabled"`
	Id        string    `json:"id"`
	OrgId     string    `json:"org_id"`
	Role      string    `json:"role"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
}

// GetAllOrganisationsParams defines parameters for GetAllOrganisations.
type GetAllOrganisationsParams struct {
	// Limit Size of the page, maximum is 100, default is 25
	Limit *int32 `form:"limit,omitempty" json:"limit,omitempty"`

	// Page The page to return
	Page *int32 `form:"page,omitempty" json:"page,omitempty"`
}

// AuthenticateLoginJSONRequestBody defines body for AuthenticateLogin for application/json ContentType.
type AuthenticateLoginJSONRequestBody = AuthenticatePostRequest

// CreateOrganisationJSONRequestBody defines body for CreateOrganisation for application/json ContentType.
type CreateOrganisationJSONRequestBody = OrganisationIn

// UserCreateJSONRequestBody defines body for UserCreate for application/json ContentType.
type UserCreateJSONRequestBody = UserIn

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Authenticate
	// (POST /v1/auth/login)
	AuthenticateLogin(c *gin.Context)
	// Health
	// (GET /v1/health)
	GetHealth(c *gin.Context)
	// Metrics
	// (GET /v1/metrics)
	GetMetrics(c *gin.Context)
	// Create an organisation
	// (POST /v1/organisation)
	CreateOrganisation(c *gin.Context)
	// Returns all organisations
	// (GET /v1/organisations)
	GetAllOrganisations(c *gin.Context, params GetAllOrganisationsParams)
	// Create User Account
	// (POST /v1/user)
	UserCreate(c *gin.Context)
	// Get all users
	// (GET /v1/user/{user})
	GetUserByID(c *gin.Context, user string)
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

// GetHealth operation middleware
func (siw *ServerInterfaceWrapper) GetHealth(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetHealth(c)
}

// GetMetrics operation middleware
func (siw *ServerInterfaceWrapper) GetMetrics(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetMetrics(c)
}

// CreateOrganisation operation middleware
func (siw *ServerInterfaceWrapper) CreateOrganisation(c *gin.Context) {

	c.Set(BearerAuthScopes, []string{"admin"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CreateOrganisation(c)
}

// GetAllOrganisations operation middleware
func (siw *ServerInterfaceWrapper) GetAllOrganisations(c *gin.Context) {

	var err error

	c.Set(BearerAuthScopes, []string{"admin"})

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

// UserCreate operation middleware
func (siw *ServerInterfaceWrapper) UserCreate(c *gin.Context) {

	c.Set(BearerAuthScopes, []string{"admin", "write", "readonly"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.UserCreate(c)
}

// GetUserByID operation middleware
func (siw *ServerInterfaceWrapper) GetUserByID(c *gin.Context) {

	var err error

	// ------------- Path parameter "user" -------------
	var user string

	err = runtime.BindStyledParameterWithOptions("simple", "user", c.Param("user"), &user, runtime.BindStyledParameterOptions{Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter user: %w", err), http.StatusBadRequest)
		return
	}

	c.Set(BearerAuthScopes, []string{"admin", "write", "readonly"})

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetUserByID(c, user)
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

	router.POST(options.BaseURL+"/v1/auth/login", wrapper.AuthenticateLogin)
	router.GET(options.BaseURL+"/v1/health", wrapper.GetHealth)
	router.GET(options.BaseURL+"/v1/metrics", wrapper.GetMetrics)
	router.POST(options.BaseURL+"/v1/organisation", wrapper.CreateOrganisation)
	router.GET(options.BaseURL+"/v1/organisations", wrapper.GetAllOrganisations)
	router.POST(options.BaseURL+"/v1/user", wrapper.UserCreate)
	router.GET(options.BaseURL+"/v1/user/:user", wrapper.GetUserByID)
}

type AuthenticateLoginRequestObject struct {
	Body *AuthenticateLoginJSONRequestBody
}

type AuthenticateLoginResponseObject interface {
	VisitAuthenticateLoginResponse(w http.ResponseWriter) error
}

type AuthenticateLogin200JSONResponse AuthenticatePostResponse

func (response AuthenticateLogin200JSONResponse) VisitAuthenticateLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type AuthenticateLogin403JSONResponse Error

func (response AuthenticateLogin403JSONResponse) VisitAuthenticateLoginResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(403)

	return json.NewEncoder(w).Encode(response)
}

type GetHealthRequestObject struct {
}

type GetHealthResponseObject interface {
	VisitGetHealthResponse(w http.ResponseWriter) error
}

type GetHealth200JSONResponse MetricsGetResponse

func (response GetHealth200JSONResponse) VisitGetHealthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetMetricsRequestObject struct {
}

type GetMetricsResponseObject interface {
	VisitGetMetricsResponse(w http.ResponseWriter) error
}

type GetMetrics200Response struct {
}

func (response GetMetrics200Response) VisitGetMetricsResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type CreateOrganisationRequestObject struct {
	Body *CreateOrganisationJSONRequestBody
}

type CreateOrganisationResponseObject interface {
	VisitCreateOrganisationResponse(w http.ResponseWriter) error
}

type CreateOrganisation201JSONResponse Organisation

func (response CreateOrganisation201JSONResponse) VisitCreateOrganisationResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response)
}

type CreateOrganisationdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response CreateOrganisationdefaultJSONResponse) VisitCreateOrganisationResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
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
	Body    Organisations
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

type UserCreateRequestObject struct {
	Body *UserCreateJSONRequestBody
}

type UserCreateResponseObject interface {
	VisitUserCreateResponse(w http.ResponseWriter) error
}

type UserCreate201JSONResponse struct {
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Enabled   *bool      `json:"enabled,omitempty"`
	Id        *string    `json:"id,omitempty"`
	OrgId     *string    `json:"org_id,omitempty"`
	Role      *string    `json:"role,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	Username  *string    `json:"username,omitempty"`
}

func (response UserCreate201JSONResponse) VisitUserCreateResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response)
}

type UserCreate409JSONResponse string

func (response UserCreate409JSONResponse) VisitUserCreateResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)

	return json.NewEncoder(w).Encode(response)
}

type GetUserByIDRequestObject struct {
	User string `json:"user"`
}

type GetUserByIDResponseObject interface {
	VisitGetUserByIDResponse(w http.ResponseWriter) error
}

type GetUserByID200JSONResponse UserNoPassword

func (response GetUserByID200JSONResponse) VisitGetUserByIDResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetUserByID404JSONResponse string

func (response GetUserByID404JSONResponse) VisitGetUserByIDResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Authenticate
	// (POST /v1/auth/login)
	AuthenticateLogin(ctx context.Context, request AuthenticateLoginRequestObject) (AuthenticateLoginResponseObject, error)
	// Health
	// (GET /v1/health)
	GetHealth(ctx context.Context, request GetHealthRequestObject) (GetHealthResponseObject, error)
	// Metrics
	// (GET /v1/metrics)
	GetMetrics(ctx context.Context, request GetMetricsRequestObject) (GetMetricsResponseObject, error)
	// Create an organisation
	// (POST /v1/organisation)
	CreateOrganisation(ctx context.Context, request CreateOrganisationRequestObject) (CreateOrganisationResponseObject, error)
	// Returns all organisations
	// (GET /v1/organisations)
	GetAllOrganisations(ctx context.Context, request GetAllOrganisationsRequestObject) (GetAllOrganisationsResponseObject, error)
	// Create User Account
	// (POST /v1/user)
	UserCreate(ctx context.Context, request UserCreateRequestObject) (UserCreateResponseObject, error)
	// Get all users
	// (GET /v1/user/{user})
	GetUserByID(ctx context.Context, request GetUserByIDRequestObject) (GetUserByIDResponseObject, error)
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

// GetHealth operation middleware
func (sh *strictHandler) GetHealth(ctx *gin.Context) {
	var request GetHealthRequestObject

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetHealth(ctx, request.(GetHealthRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetHealth")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(GetHealthResponseObject); ok {
		if err := validResponse.VisitGetHealthResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetMetrics operation middleware
func (sh *strictHandler) GetMetrics(ctx *gin.Context) {
	var request GetMetricsRequestObject

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetMetrics(ctx, request.(GetMetricsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetMetrics")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(GetMetricsResponseObject); ok {
		if err := validResponse.VisitGetMetricsResponse(ctx.Writer); err != nil {
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

// UserCreate operation middleware
func (sh *strictHandler) UserCreate(ctx *gin.Context) {
	var request UserCreateRequestObject

	var body UserCreateJSONRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.Status(http.StatusBadRequest)
		ctx.Error(err)
		return
	}
	request.Body = &body

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.UserCreate(ctx, request.(UserCreateRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "UserCreate")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(UserCreateResponseObject); ok {
		if err := validResponse.VisitUserCreateResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetUserByID operation middleware
func (sh *strictHandler) GetUserByID(ctx *gin.Context, user string) {
	var request GetUserByIDRequestObject

	request.User = user

	handler := func(ctx *gin.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetUserByID(ctx, request.(GetUserByIDRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetUserByID")
	}

	response, err := handler(ctx, request)

	if err != nil {
		ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
	} else if validResponse, ok := response.(GetUserByIDResponseObject); ok {
		if err := validResponse.VisitGetUserByIDResponse(ctx.Writer); err != nil {
			ctx.Error(err)
		}
	} else if response != nil {
		ctx.Error(fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"/DsAAP//N/VRypIlAAA=",
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