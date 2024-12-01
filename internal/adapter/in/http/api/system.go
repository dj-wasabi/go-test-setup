package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (cs *ApiHandler) GetHealth(c *gin.Context) {
	output := envelope{"health": true}
	c.JSON(http.StatusOK, output)

}

func (cs *ApiHandler) GetMetrics(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}
