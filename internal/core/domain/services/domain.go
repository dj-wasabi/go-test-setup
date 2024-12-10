package services

import (
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	"werner-dijkerman.nl/test-setup/internal/core/port/in"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	"werner-dijkerman.nl/test-setup/pkg/logging"
)

type domainServices struct {
	log *slog.Logger
	org out.PortOrganisation
	usr out.PortUser
}

var (
	model_authentication_requests = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "domain_service_http_request_authentication",
		Help: "A histogram of authentications request durations with in seconds.",
	}, []string{"state"})
)

func registerMetrics() {
	prometheus.Register(model_authentication_requests)
}

func NewdomainServices(org out.PortOrganisation, usr out.PortUser) in.ApiUseCases {
	registerMetrics()
	return &domainServices{
		log: logging.Initialize(),
		org: org,
		usr: usr,
	}
}
