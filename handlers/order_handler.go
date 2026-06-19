package handlers

import (
	"net/http"

	"github.com/KishiEdward/back-skenaid/config"
	"github.com/KishiEdward/back-skenaid/models"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct{}

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

	userID := uint(1)

	tx := config.DB.Begin()

	var cart models.Cart
	if err := tx.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Keranjang kosong atau tidak ditemukan"})
		return
	}

	if len(cart.Items) == 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Keranjang belanja kosong"})
		return
	}

	var totalAmount float64
	order := models.Order{
		UserID:          userID,
		Status:          "pending",
		ShippingAddress: input.ShippingAddress,
		Notes:           input.Notes,
		PaymentMethod:   input.PaymentMethod,
	}

	for _, item := range cart.Items {
		totalAmount += item.Subtotal

		if item.Product.Stock < item.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Stok tidak cukup untuk produk: " + item.Product.Name})
			return
		}

		newStock := item.Product.Stock - item.Quantity
		if err := tx.Model(&item.Product).Update("stock", newStock).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengurangi stok"})
			return
		}

		orderItem := models.OrderItem{
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			Price:       item.Product.Price,
			Quantity:    item.Quantity,
			Subtotal:    item.Subtotal,
		}
		order.Items = append(order.Items, orderItem)
	}

	order.TotalAmount = totalAmount

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat pesanan"})
		return
	}

	if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengosongkan keranjang"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Checkout berhasil",
		"data":    order,
	})
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	firebaseUID, exists := c.Get("firebase_uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Sesi tidak valid"})
		return
	}

	var user models.User
	if err := config.DB.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User tidak ditemukan"})
		return
	}

	var orders []models.Order
	if err := config.DB.Preload("Items").Preload("Items.Product").
		Where("user_id = ?", user.ID).
		Order("created_at desc").
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil riwayat pesanan",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil riwayat pesanan",
		"data":    orders,
	})
}

func (h *OrderHandler) GetOrderDetail(c *gin.Context) {
	orderID := c.Param("id")

	var order models.Order
	if err := config.DB.Preload("Items").First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Pesanan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil detail pesanan",
		"data":    order,
	})
}