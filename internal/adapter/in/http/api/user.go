package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
	//nolint:ineffassign,staticcheck
	ctx := utils.NewContextWrapper(c, log_id).Build()

	ctx, span := tracer.Start(c.Request.Context(), "InUserCreate")
	defer span.End()
	span.SetAttributes(
		attribute.String("http.api.file", "user"),
		attribute.String("http.api.function", "UserCreate"),
		attribute.String("code.type", "adapter.in"),
	)

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
	span.AddEvent("Creating user user", trace.WithAttributes(
		attribute.String("username", u.GetUsername()),
		attribute.String("code.type", "adapter.in"),
	))
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

func (cs *ApiHandler) GetUserByID(c *gin.Context, userId string) {
	log_id := GetXAppLogId(c)
	//nolint:ineffassign,staticcheck
	ctx := utils.NewContextWrapper(c, log_id).Build()

	ctx, span := tracer.Start(c.Request.Context(), "InUserGetById")
	defer span.End()
	span.SetAttributes(
		attribute.String("http.api.file", "user"),
		attribute.String("http.api.function", "GetUserByID"),
		attribute.String("code.type", "adapter.in"),
	)

	span.AddEvent("Getting information from token", trace.WithAttributes(
		attribute.String("userid", userId),
	))
	auth := utils.NewAuthentication()
	ad, _, adError := auth.GetAuthenticationDetails(c.Request, log_id)
	if adError != nil {
		c.AbortWithStatusJSON(http.StatusConflict, adError)
		return
	}

	if ad.GetUserId() == userId || ad.GetRole() == "admin" {
		span.AddEvent("Allowed to get data from this user", trace.WithAttributes(
			attribute.String("userid", userId),
		))
		user, userErr := cs.uc.UserGet(ctx, userId, log_id)
		if userErr != nil {
			span.RecordError(errors.New(userErr.Error))
			span.SetStatus(codes.Error, userErr.Error)
			c.AbortWithStatusJSON(http.StatusConflict, userErr)
			return
		} else {
			c.JSON(http.StatusOK, user)
		}
	} else {
		myError := model.GetError("USR0006", log_id)
		c.AbortWithStatusJSON(http.StatusForbidden, myError)
		return
	}
}
