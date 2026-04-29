package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"golang-train/internal/identity/application/usecase"
)

type Handler struct {
	getUser *usecase.GetUserUsecase
}

func RegisterRoutes(r *gin.RouterGroup, getUser *usecase.GetUserUsecase) {
	h := &Handler{getUser: getUser}
	r.GET("/users/:id", h.get)
}

func (h *Handler) get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	u, err := h.getUser.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    u.ID,
		"email": u.Email,
		"name":  u.Name,
	})
}
