package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"golang-train/internal/auth/application/usecase"
	// sharedmw "golang-train/internal/shared/middleware"
)

type Handler struct {
	login *usecase.LoginUsecase
}

func RegisterRoutes(r *gin.Engine, login *usecase.LoginUsecase) {
	h := &Handler{login: login}
	r.POST("/auth/login", h.loginHandler)
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) loginHandler(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// tenantID, ok := sharedmw.TenantFromGin(c)
	// if !ok {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "tenant missing"})
	// 	return
	// }
	tenantID := "tenant1"

	token, err := h.login.Execute(c.Request.Context(), tenantID, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": token})
}
