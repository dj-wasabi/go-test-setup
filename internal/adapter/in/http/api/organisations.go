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
	result, err := cs.uc.CreateOrganisation(context.Background(), model.NewOrganization(e.GetName(), e.GetDescription(), e.GetFqdn(), e.GetEnabled(), e.GetAdmins()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (cs *ApiHandler) GetAllOrganisations(c *gin.Context, params GetAllOrganisationsParams) {
	cs.log.Info("Get all organisations")
	data, _ := cs.uc.GetAllOrganisations(context.Background())

	c.JSON(http.StatusOK, data)
}

// dummy endpoint to just validate the authentication part
func (cs *ApiHandler) ListTags(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"error": "Ik heb hier ook geen data!"})
}
