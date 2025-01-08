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
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"werner-dijkerman.nl/test-setup/internal/core/port/in"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	intmid "werner-dijkerman.nl/test-setup/internal/middleware"
	"werner-dijkerman.nl/test-setup/pkg/config"
	"werner-dijkerman.nl/test-setup/pkg/logging"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

type envelope map[string]any

const logid string = "X-APP-LOG-ID"

var (
	tracer                            = otel.Tracer("werner-dijkerman.nl/test-setup/internal/adapter/in/http/api")
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
	uc  in.ApiUseCasesInterface
	log *slog.Logger
}

func NewApiService(as in.ApiUseCasesInterface) *ApiHandler {
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

func NewAuthenticator(po out.PortUserInterface, ts out.PortStoreInterface, h *ApiHandler, l *slog.Logger) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return intmid.ValidateSecurityScheme(po, ts, l, input)
	}
}

func NewGinServer(po out.PortUserInterface, ts out.PortStoreInterface, h *ApiHandler, c *config.Config, l *slog.Logger) *http.Server {
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

	if c.Tracing.Enabled {
		l.Info("Enable Tracing via the otelgin middleware package")
		r.Use(otelgin.Middleware(c.Tracing.Appname))
	}

	r.Use(middleware.OapiRequestValidatorWithOptions(swagger,
		&middleware.Options{
			Options: openapi3filter.Options{
				AuthenticationFunc: NewAuthenticator(po, ts, h, l),
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

func InitTracer(c *config.Config, l *slog.Logger) func(context.Context) error {
	headers := map[string]string{
		"content-type": "application/json",
	}

	l.Debug(fmt.Sprintf("Configuring Tracing app '%v' using endpoint %v", c.Tracing.Appname, c.Tracing.Endpoint))
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint(c.Tracing.Endpoint),
			otlptracehttp.WithInsecure(),
			otlptracehttp.WithHeaders(headers),
		),
	)

	if err != nil {
		l.Error(fmt.Sprintf("Error while initialising OTEL endpoint: %v", err))
	}

	tracerprovider := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(c.Tracing.Appname),
			),
		),
	)
	otel.SetTracerProvider(tracerprovider)
	return exporter.Shutdown
}
