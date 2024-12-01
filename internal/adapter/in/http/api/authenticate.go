package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
)

var (
	HttpAuthenticationRequestsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "authentication_requests",
			Help: "Number of Authentication requests in total.",
		},
	)
)

// cs.uc --> domain/services/authentication.go

func (cs *ApiHandler) AuthenticateLogin(c *gin.Context) {
	var e model.AuthenticatePostRequest
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate body, need proper validation error output.
	validate := validator.New(validator.WithRequiredStructEnabled())
	errCheck := validate.Struct(e)
	if errCheck != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errCheck.Error()})
		return
	}

	HttpAuthenticationRequestsTotal.Inc()
	timeStart := time.Now()
	token, err := cs.uc.AuthenticateLoginService(context.Background(), e.GetUsername(), e.GetPassword())
	timeEnd := float64(time.Since(timeStart).Seconds())

	if err != nil {
		authentication_requests_per_state.WithLabelValues("failure").Observe(timeEnd)
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	authentication_requests_per_state.WithLabelValues("successful").Observe(timeEnd)
	c.JSON(http.StatusOK, token)
}
