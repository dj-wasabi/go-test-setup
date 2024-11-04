package services

import (
	"log/slog"

	"werner-dijkerman.nl/test-setup/internal/core/port/in"
	"werner-dijkerman.nl/test-setup/internal/core/port/out"
	"werner-dijkerman.nl/test-setup/pkg/logging"
)

type domainServices struct {
	log *slog.Logger
	org out.PortOrganisation
	usr out.PortUser
}

func NewdomainServices(org out.PortOrganisation, usr out.PortUser) in.ApiUseCases {
	return &domainServices{
		log: logging.Initialize(),
		org: org,
		usr: usr,
	}
}
