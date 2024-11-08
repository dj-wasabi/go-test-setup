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
	"werner-dijkerman.nl/test-setup/internal/adapter/out/mongodb"
	"werner-dijkerman.nl/test-setup/internal/core/port/in"
	intmid "werner-dijkerman.nl/test-setup/internal/middleware"
	"werner-dijkerman.nl/test-setup/pkg/config"
	"werner-dijkerman.nl/test-setup/pkg/logging"
)

type envelope map[string]any

type ApiHandler struct {
	uc  in.ApiUseCases
	log *slog.Logger
}

func NewApiService(s in.ApiUseCases) *ApiHandler {
	return &ApiHandler{
		uc:  s,
		log: logging.Initialize(),
	}
}

func NewAuthenticator(mc *mongodb.MongodbConnection, h *ApiHandler) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return intmid.ValidateSecurityScheme(mc, ctx, input)
	}
}

func NewGinServer(mc *mongodb.MongodbConnection, h *ApiHandler, c *config.Config) *http.Server {
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
				AuthenticationFunc: NewAuthenticator(mc, h),
			},
		},
	))

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	// r.Use(middleware.OapiRequestValidator(swagger))
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
