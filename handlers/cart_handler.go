package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
}

func NewCartHandler() *CartHandler {
	return &CartHandler{}
}

func (h *CartHandler) GetCart(c *gin.Context) {

	mockData := gin.H{
		"items":      []gin.H{}, 
		"total":      0,
		"item_count": 0,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil keranjang",
		"data":    mockData,
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Produk berhasil ditambahkan ke keranjang",
	})
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	cartItemID := c.Param("id")
	_ = cartItemID 
	var input struct {
		Quantity int `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Jumlah produk diperbarui",
	})
}

func (h *CartHandler) RemoveCartItem(c *gin.Context) {
	cartItemID := c.Param("id")
	_ = cartItemID

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Produk dihapus dari keranjang",
	})
}

func (h *CartHandler) ClearCart(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Keranjang berhasil dikosongkan",
	})
}