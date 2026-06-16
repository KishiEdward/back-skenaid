package handlers

import (
	"net/http"

	"github.com/KishiEdward/back-skenaid/config"
	"github.com/KishiEdward/back-skenaid/models"
	"github.com/gin-gonic/gin"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GET /v1/user/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	// ==========================================
	// Nanti buka komentar ini jika AuthMiddleware 
	// sudah menyimpan "user_id" dari Firebase:
	//
	// userIDRaw, exists := c.Get("user_id")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Tidak ada akses"})
	// 	return
	// }
	// userID := userIDRaw.(uint)
	// ==========================================

	// Hardcode ID 1 untuk testing UI Profil & Checkout
	userID := uint(1) 

	var user models.User
	// GORM otomatis mencari berdasarkan Primary Key (ID) bawaan gorm.Model
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false, 
			"message": "Data user tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Berhasil mengambil profil user",
		"data":    user,
	})
}