package service

import (
	"github.com/google/uuid"
	"productService/models"
	"productService/repository"
)

type ProductService struct {
	Repo *repository.ProductRepository
}

func (ps *ProductService) Create(pcd *models.Product) (*models.Product, error) {
	return ps.Repo.Create(pcd)
}
func (ps *ProductService) GetProductById(id uuid.UUID) (*models.Product, error) {
	return ps.Repo.GetProductByID(id)
}

func (ps *ProductService) GetProducts() ([]models.Product, error) {
	return ps.Repo.GetProducts()
}

func (ps *ProductService) UpdateProduct(id uuid.UUID, product models.Product) (*models.Product, error) {
	return ps.Repo.UpdateProduct(id, product)
}

func (ps *ProductService) DeleteProduct(id uuid.UUID) (*models.Product, error) {
	return ps.Repo.DeleteProduct(id)
}

func (ps *ProductService) GetProductDetailsByIDs(productIDs []uuid.UUID) ([]models.ProductDetailsResponse, error) {
	products, err := ps.Repo.GetProductDetailsByIDs(productIDs)
	if err != nil {
		return nil, err
	}

	var responseDTOs []models.ProductDetailsResponse
	for _, p := range products {
		responseDTOs = append(responseDTOs, models.ProductDetailsResponse{
			ID:       p.ID,
			Name:     p.Name,
			Price:    p.Price,
			Quantity: p.Quantity,
		})
	}

	return responseDTOs, nil
}
