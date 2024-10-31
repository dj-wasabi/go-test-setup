package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"werner-dijkerman.nl/test-setup/domain/model"
)

// cs.uc --> domain/services/organisation

func (cs *ApiHandler) UserCreate(c *gin.Context) {
	cs.log.Debug("Ceate an organisation")
	var u model.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cs.uc.UserCreate(context.Background(), model.NewUser(u.GetUsername(), u.GetPassword(), u.GetEnabled(), u.GetRoles()))
	c.JSON(http.StatusOK, u)
}
