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

// AuthenticateRequest Authenticate
type AuthenticateRequest struct {
	Password string `json:"password" validate:"required,min=6,max=256"`
	Username string `json:"username" validate:"required,min=6,max=64,alphanum"`
}

// AuthenticateToken Authenticate
type AuthenticateToken struct {
	Token string `json:"token" validate:"required"`
}

// Error Default error response
type Error struct {
	Error string `json:"error"`
}

// MetricsGetResponse Metrics for Prometheus response
type MetricsGetResponse = string

// Organisation The model Organisation object
type Organisation struct {
	Admins      []string  `json:"admins"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description" validate:"required,min=1,max=256"`
	Enabled     bool      `json:"enabled"`
	Fqdn        string    `json:"fqdn"`
	Id          string    `json:"id,omitempty"`
	Name        string    `json:"name" validate:"required,min=6,max=64,alphanum"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// OrganisationIn Input to create the Organisation object
type OrganisationIn struct {
	Admins      []string `json:"admins"`
	Description string   `json:"description" validate:"required,min=1,max=256"`
	Enabled     bool     `json:"enabled"`
	Fqdn        string   `json:"fqdn" validate:"required,fqdn"`
	Name        string   `json:"name" validate:"required,min=6,max=64,alphanum"`
}

// Organisations Overview of all organisations
type Organisations = []Organisation

// UserIn Input to create the User object
type UserIn struct {
	Enabled  bool   `json:"enabled"`
	OrgId    string `json:"org_id" validate:"required"`
	Password string `json:"password" validate:"required,min=6,max=256"`
	Role     string `json:"role,omitempty" validate:"oneof=admin write readonly"`
	Username string `json:"username" validate:"required,min=6,max=64,alphanum"`
}

// UserNoPassword Return of the User object
type UserNoPassword struct {
	CreatedAt time.Time `json:"created_at"`
	Enabled   bool      `json:"enabled"`
	Id        string    `json:"id" validate:"required"`
	OrgId     string    `json:"org_id" validate:"required"`
	Role      string    `json:"role" validate:"oneof=admin write readonly"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username" validate:"required,min=6,max=64,alphanum"`
}

// GetAllOrganisationsParams defines parameters for GetAllOrganisations.
type GetAllOrganisationsParams struct {
	// Limit Size of the page, maximum is 100, default is 25
	Limit *int32 `form:"limit,omitempty" json:"limit,omitempty"`

	// Page The page to return
	Page *int32 `form:"page,omitempty" json:"page,omitempty"`
}

// AuthenticateLoginJSONRequestBody defines body for AuthenticateLogin for application/json ContentType.
type AuthenticateLoginJSONRequestBody = AuthenticateRequest

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

type AuthenticateLogin200JSONResponse AuthenticateToken

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

type GetHealth503JSONResponse Error

func (response GetHealth503JSONResponse) VisitGetHealthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(503)

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

type CreateOrganisation403JSONResponse Error

func (response CreateOrganisation403JSONResponse) VisitCreateOrganisationResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(403)

	return json.NewEncoder(w).Encode(response)
}

type CreateOrganisation409JSONResponse string

func (response CreateOrganisation409JSONResponse) VisitCreateOrganisationResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)

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

type GetAllOrganisations403JSONResponse Error

func (response GetAllOrganisations403JSONResponse) VisitGetAllOrganisationsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(403)

	return json.NewEncoder(w).Encode(response)
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

type UserCreate403JSONResponse Error

func (response UserCreate403JSONResponse) VisitUserCreateResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(403)

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

type GetUserByID403JSONResponse Error

func (response GetUserByID403JSONResponse) VisitGetUserByIDResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(403)

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
