package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
)

// cs.uc --> domain/services/organisation

func (cs *ApiHandler) CreateOrganisation(c *gin.Context) {
	cs.log.Debug("Ceate an organisation")
	var e model.Organisation
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	organisationObject, organisationError := model.NewOrganization(e.GetName(), e.GetDescription(), e.GetFqdn(), e.GetEnabled(), e.GetAdmins())
	if organisationError != nil {
		error := model.NewError(organisationError.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, error)
		return
	}

	createOutput, createError := cs.uc.CreateOrganisation(context.Background(), organisationObject)
	if createError != nil {
		c.JSON(http.StatusInternalServerError, createError)
		return
	}
	c.JSON(http.StatusOK, createOutput)
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
