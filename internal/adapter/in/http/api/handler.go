package api

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	middleware "github.com/oapi-codegen/gin-middleware"
	"github.com/prometheus/client_golang/prometheus"
	"werner-dijkerman.nl/test-setup/internal/core/port/in"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	intmid "werner-dijkerman.nl/test-setup/internal/middleware"
	"werner-dijkerman.nl/test-setup/pkg/config"
	"werner-dijkerman.nl/test-setup/pkg/logging"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

type envelope map[string]any

var logid string = "X-APP-LOG-ID"

var (
	authentication_requests_per_state = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "adapter_in_http_request_authentication",
		Help: "A histogram of authentications request durations with in seconds.",
	}, []string{"state"})
	user_create_requests = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "adapter_in_http_request_create_user",
		Help: "A histogram for creation of user related request with durations in seconds.",
	}, []string{"state"})
	organisation_create_requests = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "adapter_in_http_request_create_organisation",
		Help: "A histogram for creation of organisation related request with durations in seconds.",
	}, []string{"state"})
)

type ApiHandler struct {
	uc  in.ApiUseCases
	log *slog.Logger
}

func NewApiService(as in.ApiUseCases) *ApiHandler {
	return &ApiHandler{
		uc:  as,
		log: logging.Initialize(),
	}
}

func GetXAppLogId(c *gin.Context) string {
	log_id := c.Request.Header.Get(logid)

	if log_id == "" {
		logId_string := uuid.New()
		log_id = logId_string.String()
	}

	return log_id
}

func registerMetrics() {
	_ = prometheus.Register(authentication_requests_per_state)
	_ = prometheus.Register(user_create_requests)
	_ = prometheus.Register(organisation_create_requests)
}

func NewAuthenticator(po out.PortUser, h *ApiHandler, l *slog.Logger) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return intmid.ValidateSecurityScheme(po, l, input)
	}
}

func NewGinServer(po out.PortUser, h *ApiHandler, c *config.Config, l *slog.Logger) *http.Server {
	swagger, err := GetSwagger()

	if err != nil {
		h.log.Error(fmt.Sprintf("Error loading swagger spec\n: %s", err.Error()))
		os.Exit(1)
	}

	swagger.Servers = nil

	dir := filepath.Dir(c.Http.Logfile)
	ok, err := utils.IsWritable(dir)
	if !ok {
		h.log.Error(fmt.Sprintf("Error while create a file in directory: '%v'", err))
		os.Exit(1)
	}

	f, _ := os.Create(c.Http.Logfile)
	gin.DefaultWriter = io.MultiWriter(f)
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(intmid.JsonLoggerMiddleware())

	r.Use(middleware.OapiRequestValidatorWithOptions(swagger,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: NewAuthenticator(po, h, l),
			},
		},
	))
	// Register metrics
	registerMetrics()

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	RegisterHandlers(r, h)

	s := &http.Server{
		Handler:      r,
		Addr:         net.JoinHostPort("0.0.0.0", c.Http.Listen),
		IdleTimeout:  time.Duration(c.Http.Timeout.Idle) * time.Second,
		ReadTimeout:  time.Duration(c.Http.Timeout.Read) * time.Second,
		WriteTimeout: time.Duration(c.Http.Timeout.Write) * time.Second,
		ErrorLog:     slog.NewLogLogger(h.log.Handler(), slog.LevelError),
	}
	return s
}
