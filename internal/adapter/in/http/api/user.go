package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"werner-dijkerman.nl/test-setup/internal/core/domain/model"
)

// cs.uc --> domain/services/organisation

func (cs *ApiHandler) UserCreate(c *gin.Context) {
	cs.log.Debug("Ceate a user")
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
	createMessage, createError := cs.uc.UserCreate(context.Background(), userObject)
	if createError != nil {
		c.AbortWithStatusJSON(http.StatusConflict, createError)
		return
	} else {
		c.JSON(http.StatusOK, createMessage)
	}
}

func (cs *ApiHandler) GetUserByID(c *gin.Context, user string) {

}
