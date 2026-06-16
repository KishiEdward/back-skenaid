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

	userID := uint(1) // Sementara hardcode user ID 1 untuk testing

	// Mulai Database Transaction (jika satu gagal, semua dibatalkan)
	tx := config.DB.Begin()

	// 1. Ambil keranjang user
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

	// 2. Buat Order Baru
	var totalAmount float64
	order := models.Order{
		UserID:          userID,
		Status:          "pending",
		ShippingAddress: input.ShippingAddress,
		Notes:           input.Notes,
		PaymentMethod:   input.PaymentMethod,
	}

	// 3. Pindahkan CartItem ke OrderItem DAN Kurangi Stok
	for _, item := range cart.Items {
		totalAmount += item.Subtotal

		// CEK STOK
		if item.Product.Stock < item.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Stok tidak cukup untuk produk: " + item.Product.Name})
			return
		}

		// KURANGI STOK PRODUK
		newStock := item.Product.Stock - item.Quantity
		if err := tx.Model(&item.Product).Update("stock", newStock).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengurangi stok"})
			return
		}

		// Buat item pesanan
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

	// Simpan pesanan ke database
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat pesanan"})
		return
	}

	// 4. Kosongkan Keranjang
	if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengosongkan keranjang"})
		return
	}

	// Selesaikan transaksi
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Checkout berhasil",
		"data":    order,
	})
}

func (h *OrderHandler) GetMyOrders(c *gin.Context) {
	userID := uint(1) // Hardcode sementara

	var orders []models.Order
	if err := config.DB.Preload("Items").Where("user_id = ?", userID).Order("created_at desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal mengambil pesanan"})
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