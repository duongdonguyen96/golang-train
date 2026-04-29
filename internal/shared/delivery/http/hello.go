package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterHello @Summary Hello API
// @Description Say hello
// @Tags test
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /hello [get]
func RegisterHello(r *gin.Engine) {
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Nguyen"})
	})
}
