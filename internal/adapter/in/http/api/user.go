package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
	"werner-dijkerman.nl/test-setup/internal/core/port/in"
)

// cs.uc --> domain/services/organisation

func (cs *ApiHandler) UserCreate(c *gin.Context) {
	cs.log.Debug("Ceate an organisation")
	var u in.UserIn
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input, _ := model.NewUser(u.GetUsername(), u.GetPassword(), u.GetEnabled(), u.GetRoles())
	message, err := cs.uc.UserCreate(context.Background(), input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, err)
	} else {
		output := envelope{"id": message}
		c.JSON(http.StatusOK, output)
	}
}
