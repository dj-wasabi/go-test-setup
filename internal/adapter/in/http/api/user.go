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

// cs.uc --> domain/services/organisation

var (
	HttpUserRequestsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "adapter_in_http_request_users",
			Help: "Number of total User related requests.",
		},
	)
)

func (cs *ApiHandler) UserCreate(c *gin.Context) {
	log_id := GetXAppLogId(c)
	ctx := utils.NewContextWrapper(c, log_id).Build()

	var u model.User
	if err := c.ShouldBindJSON(&u); err != nil {
		error := model.NewError(err.Error())
		c.JSON(http.StatusBadRequest, error)
		return
	}
	userObject, userError := model.NewUser(u.GetUsername(), u.GetPassword(), u.GetRole(), u.GetEnabled(), u.GetOrgId())
	if userError != nil {
		error := model.NewError(userError.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, error)
		return
	}
	HttpUserRequestsTotal.Inc()
	timeStart := time.Now()
	createMessage, createError := cs.uc.UserCreate(ctx, userObject)
	timeEnd := float64(time.Since(timeStart).Seconds())

	if createError != nil {
		user_create_requests.WithLabelValues("failure").Observe(timeEnd)
		c.AbortWithStatusJSON(http.StatusConflict, createError)
		return
	} else {
		user_create_requests.WithLabelValues("successful").Observe(timeEnd)
		c.JSON(http.StatusOK, createMessage)
	}
}

func (cs *ApiHandler) GetUserByID(c *gin.Context, user string) {
	// log_id := GetXAppLogId(c)
	// ctx := utils.NewContextWrapper(c, log_id).Build()

	HttpUserRequestsTotal.Inc()

}
