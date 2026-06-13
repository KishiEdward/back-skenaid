package main

import (
	"log"

	"github.com/KishiEdward/back-skenaid/config"
	"github.com/KishiEdward/back-skenaid/models"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.InitDatabase()

	products := []models.Product{
		{
			Name:        "Kalung Rantai Anime Berserk",
			Description: "Kalung rantai anti karat dengan statement, cocok untuk OOTD.",
			Price:       75000,
			Stock:       50,
			Category:    "Kalung",
			ImageURL:    "https://i.ibb.co.com/DgGB856X/berserk.jpg",
			Material:    "Titanium murni",
		},
		{
			Name:        "Kalung Rantai Liontin Bintang",
			Description: "Kalung rantai anti karat dengan liontin bintang, cocok untuk OOTD.",
			Price:       75000,
			Stock:       50,
			Category:    "Kalung",
			ImageURL:    "https://i.ibb.co.com/KjkZkZxh/bintang.jpg",
			Material:    "Titanium murni",
		},
	}

	for _, p := range products {
		config.DB.Create(&p)
	}

	log.Printf("Seed berhasil: %d produk ditambahkan", len(products))
}