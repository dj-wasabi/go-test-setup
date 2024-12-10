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

var (
	HttpAuthenticationRequestsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "authentication_requests",
			Help: "Number of total Authentication related requests.",
		},
	)
)

// cs.uc --> domain/services/authentication.go

func (cs *ApiHandler) AuthenticateLogin(c *gin.Context) {
	var e model.AuthenticateRequest
	if err := c.ShouldBindJSON(&e); err != nil {
		error := model.NewError(err.Error())
		c.JSON(http.StatusBadRequest, error)
		return
	}

	errCheck := model.ValidateAuthenticationData(&e)
	if errCheck != nil {
		error := model.NewError(errCheck.Error())
		c.JSON(http.StatusBadRequest, error)
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
