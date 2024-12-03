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
	input, myerr := model.NewUser(u.GetUsername(), u.GetPassword(), u.GetRole(), u.GetEnabled(), u.GetOrgId())
	if myerr != nil {
		error := &model.Error{
			Message: myerr.Error(),
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, error)
		return
	}
	message, err := cs.uc.UserCreate(context.Background(), input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusConflict, err)
		return
	} else {
		output := envelope{"id": message}
		c.JSON(http.StatusOK, output)
	}
}

func (cs *ApiHandler) GetUserByID(c *gin.Context, user string) {

}
