package handlers

import (
	"net/http"

	"github.com/KishiEdward/back-skenaid/config"
	"github.com/KishiEdward/back-skenaid/models"
	"github.com/gin-gonic/gin"
)

type CartHandler struct{}

func NewCartHandler() *CartHandler {
	return &CartHandler{}
}

func getUserID(c *gin.Context) uint {
	return 1
}

func (h *CartHandler) GetCart(c *gin.Context) {
	userID := getUserID(c)
	var cart models.Cart

	err := config.DB.Preload("Items.Product").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		cart = models.Cart{UserID: userID}
		config.DB.Create(&cart)
	}

	var total float64
	for _, item := range cart.Items {
		total += item.Subtotal
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil keranjang",
		"data": gin.H{
			"items":      cart.Items,
			"total":      total,
			"item_count": len(cart.Items),
		},
	})
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	var input struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	userID := getUserID(c)
	var cart models.Cart
	if err := config.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		cart = models.Cart{UserID: userID}
		config.DB.Create(&cart)
	}

	var product models.Product
	if err := config.DB.First(&product, input.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Produk tidak ditemukan"})
		return
	}

	var cartItem models.CartItem
	err := config.DB.Where("cart_id = ? AND product_id = ?", cart.ID, input.ProductID).First(&cartItem).Error

	if err == nil {
		// Update jika produk sudah ada
		cartItem.Quantity += input.Quantity
		cartItem.Subtotal = float64(cartItem.Quantity) * product.Price
		config.DB.Save(&cartItem)
	} else {
		// Buat baru jika belum ada
		newItem := models.CartItem{
			CartID:    cart.ID,
			ProductID: input.ProductID,
			Quantity:  input.Quantity,
			Subtotal:  float64(input.Quantity) * product.Price,
		}
		config.DB.Create(&newItem)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Produk berhasil ditambahkan ke keranjang"})
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	cartItemID := c.Param("id")
	var input struct {
		Quantity int `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	var cartItem models.CartItem
	if err := config.DB.Preload("Product").First(&cartItem, cartItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Item keranjang tidak ditemukan"})
		return
	}

	cartItem.Quantity = input.Quantity
	cartItem.Subtotal = float64(input.Quantity) * cartItem.Product.Price
	config.DB.Save(&cartItem)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Jumlah produk diperbarui"})
}

func (h *CartHandler) RemoveCartItem(c *gin.Context) {
	cartItemID := c.Param("id")
	config.DB.Delete(&models.CartItem{}, cartItemID)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Produk dihapus dari keranjang"})
}

func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := getUserID(c)
	var cart models.Cart
	if err := config.DB.Where("user_id = ?", userID).First(&cart).Error; err == nil {
		// Hapus semua item yang memiliki cart_id milik user ini
		config.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Keranjang berhasil dikosongkan"})
}