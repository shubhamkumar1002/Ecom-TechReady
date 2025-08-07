package repository

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"productService/models"
	"time"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (prd *ProductRepository) Create(productCreateDTO *models.Product) (*models.Product, error) {
	productId := uuid.New()

	newProduct := &models.Product{
		ID:          productId,
		Name:        productCreateDTO.Name,
		Description: productCreateDTO.Description,
		Price:       productCreateDTO.Price,
		Quantity:    productCreateDTO.Quantity,
		CreatedAt:   time.Now(),
	}

	err := prd.DB.Create(&newProduct).Error
	if err != nil {
		return nil, err
	}

	return newProduct, nil
}

func (prd *ProductRepository) GetProductByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := prd.DB.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (prd *ProductRepository) GetProducts() ([]models.Product, error) {
	var products []models.Product
	err := prd.DB.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (prd *ProductRepository) UpdateProduct(id uuid.UUID, productInput models.Product) (*models.Product, error) {
	var updateProduct models.Product
	err := prd.DB.Model(&updateProduct).
		Where("id = ?", id).
		Omit("ID", "CreatedAt").
		Updates(productInput).Error
	if err != nil {
		return nil, err
	}
	return &updateProduct, nil
}

func (prd *ProductRepository) GetProductDetailsByIDs(productIDs []uuid.UUID) ([]models.Product, error) {
	var products []models.Product
	err := prd.DB.Where("id IN ?", productIDs).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (repo *ProductRepository) ReduceStockForOrder(items []models.ItemRequest) (float64, error) {
	var totalAmount float64

	err := repo.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			var product models.Product

			if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?",
				item.ProductID).First(&product).Error; err != nil {
				log.Println("product with ID %s not found", item.ProductID)
				return fmt.Errorf("product with ID %s not found", item.ProductID)
			}

			if product.Quantity < item.Quantity {
				log.Println("not enough quantity for product ID %s. Available: %d, Requested: %d",
					item.ProductID, product.Quantity, item.Quantity)
				return fmt.Errorf("not enough quantity for product ID %s. Available: %d, Requested: %d",
					item.ProductID, product.Quantity, item.Quantity)
			}

			quantity := product.Quantity - item.Quantity
			if err := tx.Model(&product).Update("quantity", quantity).Error; err != nil {
				return err
			}

			totalAmount += product.Price * float64(item.Quantity)
		}

		return nil
	})

	return totalAmount, err
}

func (prd *ProductRepository) DeleteProduct(id uuid.UUID) (*models.Product, error) {
	var productToDelete models.Product
	err := prd.DB.First(&productToDelete, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	err = prd.DB.Delete(&productToDelete).Error
	if err != nil {
		return nil, err
	}

	return &productToDelete, nil
}
