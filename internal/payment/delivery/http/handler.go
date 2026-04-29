package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"golang-train/internal/payment/application/usecase"
)

type Handler struct {
	create *usecase.CreatePaymentUsecase
	list   *usecase.ListPaymentsUsecase
	export *usecase.ExportPaymentsUsecase
}

func RegisterRoutes(r *gin.RouterGroup, create *usecase.CreatePaymentUsecase, list *usecase.ListPaymentsUsecase, export *usecase.ExportPaymentsUsecase) {
	h := &Handler{create: create, list: list, export: export}
	r.POST("/payments", h.createPayment)
	r.GET("/payments", h.listPayments)
	r.POST("/payments/export", h.exportPayments)
}

type createReq struct {
	UserID   uint64 `json:"user_id"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

func (h *Handler) createPayment(c *gin.Context) {
	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p, err := h.create.Execute(c.Request.Context(), req.UserID, req.Amount, req.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": p.ID})
}

func (h *Handler) listPayments(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))
	rows, err := h.list.Execute(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *Handler) exportPayments(c *gin.Context) {
	url, err := h.export.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}
