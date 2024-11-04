package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
)

// cs.uc --> domain/services/authentication.go

func (cs *ApiHandler) AuthenticateLogin(c *gin.Context) {
	// cs.log.Debug("About to Create an Experience")
	var e model.AuthenticationRequest
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cs.uc.AuthenticateLogin(context.Background(), e.GetUsername(), e.GetPassword())
	c.JSON(http.StatusOK, e)
}
