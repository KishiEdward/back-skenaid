package repositories

import (
	"github.com/KishiEdward/back-skenaid/config"
	"github.com/KishiEdward/back-skenaid/models"
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

func (r *ProductRepository) FindAll(page, limit int, category string) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := config.DB.Model(&models.Product{}).Where("is_active = ?", true)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	query.Count(&total)
	offset := (page - 1) * limit
	result := query.Offset(offset).Limit(limit).Find(&products)

	return products, total, result.Error
}

func (r *ProductRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	result := config.DB.First(&product, id)
	return &product, result.Error
}

func (r *ProductRepository) Create(product *models.Product) error {
	return config.DB.Create(product).Error
}

func (r *ProductRepository) Update(product *models.Product) error {
	return config.DB.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return config.DB.Delete(&models.Product{}, id).Error
}