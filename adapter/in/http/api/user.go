package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"werner-dijkerman.nl/test-setup/domain/model"
	"werner-dijkerman.nl/test-setup/port/in"
)

// cs.uc --> domain/services/organisation

func (cs *ApiHandler) UserCreate(c *gin.Context) {
	cs.log.Debug("Ceate an organisation")
	var u in.UserIn
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	message := cs.uc.UserCreate(context.Background(), model.NewUser(u.GetUsername(), u.GetPassword(), u.GetEnabled(), u.GetRoles()))
	if message == "duplicate" {
		output := envelope{"error": "Duplicate account"}
		c.AbortWithStatusJSON(http.StatusConflict, output)
	} else {
		output := envelope{"id": message}
		c.JSON(http.StatusOK, output)
	}
}
