package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
)

// cs.uc --> domain/services/organisation

var (
	HttpOrganisationRequestsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "organisation_http_requests",
			Help: "Number of total Organisation related requests.",
		},
	)
)

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

	HttpOrganisationRequestsTotal.Inc()
	timeStart := time.Now()
	createOutput, createError := cs.uc.CreateOrganisation(context.Background(), organisationObject)
	timeEnd := float64(time.Since(timeStart).Seconds())
	if createError != nil {
		organisation_create_requests.WithLabelValues("failure").Observe(timeEnd)
		c.JSON(http.StatusInternalServerError, createError)
		return
	}
	organisation_create_requests.WithLabelValues("successful").Observe(timeEnd)
	c.JSON(http.StatusOK, createOutput)
}

func (cs *ApiHandler) GetAllOrganisations(c *gin.Context, params GetAllOrganisationsParams) {
	cs.log.Info("Get all organisations")
	HttpOrganisationRequestsTotal.Inc()
	data, _ := cs.uc.GetAllOrganisations(context.Background())

	c.JSON(http.StatusOK, data)
}

// dummy endpoint to just validate the authentication part
func (cs *ApiHandler) ListTags(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"error": "Ik heb hier ook geen data!"})
}
