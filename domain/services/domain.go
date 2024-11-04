package services

import (
	"log/slog"

	"werner-dijkerman.nl/test-setup/pkg/logging"
	"werner-dijkerman.nl/test-setup/port/in"
	"werner-dijkerman.nl/test-setup/port/out"
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
