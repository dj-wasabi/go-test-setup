package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/pkg/utils"
)

var (
	HttpAuthenticationRequestsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "adapter_in_http_request_authentications",
			Help: "Number of total Authentication related requests.",
		},
	)
)

// cs.uc --> domain/services/authentication.go

func (cs *ApiHandler) AuthenticateLogin(c *gin.Context) {
	log_id := GetXAppLogId(c)
	//nolint:ineffassign,staticcheck
	ctx := utils.NewContextWrapper(c, log_id).Build()

	ctx, span := tracer.Start(c.Request.Context(), "InAuthenticate")
	defer span.End()

	var e model.AuthenticateRequest
	if err := c.ShouldBindJSON(&e); err != nil {
		utils.HandleHTTPError(c, http.StatusBadRequest, err)
		return
	}

	if errCheck := model.ValidateAuthenticationData(&e); errCheck != nil {
		utils.HandleHTTPError(c, http.StatusBadRequest, errCheck)
		return
	}

	HttpAuthenticationRequestsTotal.Inc()
	timeStart := time.Now()

	token, tokenErr := cs.uc.AuthenticateLoginService(ctx, e.GetUsername(), e.GetPassword(), log_id)
	timeEnd := float64(time.Since(timeStart).Seconds())

	if tokenErr != nil {
		authentication_requests_per_state.WithLabelValues("failure").Observe(timeEnd)
		c.JSON(http.StatusUnauthorized, tokenErr)
		return
	}
	authentication_requests_per_state.WithLabelValues("successful").Observe(timeEnd)
	c.JSON(http.StatusOK, token)
}
