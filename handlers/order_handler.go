package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}


func (h *OrderHandler) Checkout(c *gin.Context) {
	var input struct {
		ShippingAddress string `json:"shipping_address" binding:"required"`
		Notes           string `json:"notes"`
		PaymentMethod   string `json:"payment_method" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	
	mockOrder := gin.H{
		"id":               1,
		"total_amount":     150000,
		"status":           "pending",
		"shipping_address": input.ShippingAddress,
		"notes":            input.Notes,
		"payment_method":   input.PaymentMethod,
		"items":            []gin.H{}, // Kosong sementara
		"created_at":       "2026-06-14T10:00:00Z",
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Checkout berhasil",
		"data":    mockOrder,
	})
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil riwayat pesanan",
		"data":    []gin.H{},
	})
}

func (h *OrderHandler) GetOrderDetail(c *gin.Context) {
	orderID := c.Param("id")

	mockOrder := gin.H{
		"id":     orderID,
		"status": "pending",
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil detail pesanan",
		"data":    mockOrder,
	})
}