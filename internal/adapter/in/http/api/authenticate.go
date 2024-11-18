package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
)

// cs.uc --> domain/services/authentication.go

func (cs *ApiHandler) AuthenticateLogin(c *gin.Context) {
	var e model.AuthenticationRequest
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := cs.uc.AuthenticateLoginService(context.Background(), e.GetUsername(), e.GetPassword())
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	c.JSON(http.StatusOK, token)
}
