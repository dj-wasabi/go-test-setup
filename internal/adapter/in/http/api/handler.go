package api

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
	"github.com/prometheus/client_golang/prometheus"
	"werner-dijkerman.nl/test-setup/internal/core/port/in"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	intmid "werner-dijkerman.nl/test-setup/internal/middleware"
	"werner-dijkerman.nl/test-setup/pkg/config"
	"werner-dijkerman.nl/test-setup/pkg/logging"
)

type envelope map[string]any

var (
	authentication_requests_per_state = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "authentication_requests_per_state",
		Help: "A histogram of authentications request durations in seconds per state.",
	},
		[]string{"state"})
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

func RegisterMetrics() {
	prometheus.Register(authentication_requests_per_state)

}

func NewAuthenticator(po out.PortUser, h *ApiHandler, l *slog.Logger) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return intmid.ValidateSecurityScheme(po, l, ctx, input)
	}
}

func NewGinServer(po out.PortUser, h *ApiHandler, c *config.Config, l *slog.Logger) *http.Server {
	swagger, err := GetSwagger()

	if err != nil {
		h.log.Error(fmt.Sprintf("Error loading swagger spec\n: %s", err.Error()))
		os.Exit(1)
	}

	swagger.Servers = nil

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
	RegisterMetrics()

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
