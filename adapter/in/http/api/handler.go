package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"
	"werner-dijkerman.nl/test-setup/pkg/config"
	"werner-dijkerman.nl/test-setup/pkg/logging"
	"werner-dijkerman.nl/test-setup/port/in"
)

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

func NewGinServer(h *ApiHandler, c *config.Config) *http.Server {
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
	r.Use(jsonLoggerMiddleware())

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	r.Use(middleware.OapiRequestValidator(swagger))
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

func jsonLoggerMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})

			log["time"] = params.TimeStamp.Format("2006-01-02 15:04:05")
			log["status_code"] = params.StatusCode
			log["path"] = params.Path
			log["method"] = params.Method
			log["remote_addr"] = params.ClientIP
			log["user_agent"] = params.Request.UserAgent()
			log["referer"] = params.Request.Referer()
			// time endresult is logged in milliseconds: 1000943 == 1 second.
			log["duration"] = time.Duration(params.Latency) / 1000

			s, _ := json.Marshal(log)
			return string(s) + "\n"
		},
	)
}
