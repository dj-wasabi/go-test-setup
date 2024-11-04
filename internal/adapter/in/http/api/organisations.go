package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/internal/core/port/in"
)

// cs.uc --> domain/services/organisation

func (cs *ApiHandler) CreateOrganisation(c *gin.Context) {
	cs.log.Debug("Ceate an organisation")
	var e in.OrganisationIn
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cs.uc.CreateOrganisation(context.Background(), model.NewOrganization(e.GetName(), e.GetDescription(), e.GetFqdn(), e.GetEnabled(), e.GetAdmins()))
	c.JSON(http.StatusOK, e)
}

func (cs *ApiHandler) GetAllOrganisations(c *gin.Context, params GetAllOrganisationsParams) {
	cs.log.Info("Get all organisations")
	// AllOrganisations := &in.ListOrganisationsIn{}
	// AllOrganisations := &model.ListOrganisations{}
	data, _ := cs.uc.GetAllOrganisations(context.Background())

	// for i := range data {
	// 	AllOrganisations = append(AllOrganisations.Organisations, i)
	// }

	c.JSON(http.StatusOK, data)
}

func (cs *ApiHandler) ListTags(c *gin.Context) {
	// cs.log.Debug("About to Create an Experience")
	var e model.Organization
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cs.uc.CreateOrganisation(context.Background(), model.NewOrganization(e.GetName(), e.GetDescription(), e.GetFqdn(), e.GetEnabled(), e.GetAdmins()))
}
