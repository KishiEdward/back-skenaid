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

func (h *UserHandler) GetProfile(c *gin.Context) {
	firebaseUID, exists := c.Get("firebase_uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Sesi tidak valid",
		})
		return
	}

	var user models.User
	if err := config.DB.Where("firebase_uid = ?", firebaseUID).First(&user).Error; err != nil {
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