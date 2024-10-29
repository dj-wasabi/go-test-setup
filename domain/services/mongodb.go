package services

import (
	"werner-dijkerman.nl/test-setup/port/in"
	"werner-dijkerman.nl/test-setup/port/out"
)

type DBServices struct {
	db out.OrganisationsDBPort
}

func NewDBService(db out.OrganisationsDBPort) in.ApiUseCases {
	return &DBServices{
		db: db,
	}
}
