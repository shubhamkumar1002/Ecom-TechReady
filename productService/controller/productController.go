package controller

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"productService/models"
	"productService/service"
)

type ProductController struct {
	Service service.ProductService
}

func (pc *ProductController) CreateProduct(ctx iris.Context) {
	var product models.Product
	if err := ctx.ReadJSON(&product); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid input"})
		return
	}

	newProduct, err := pc.Service.Create(&product)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to create product"})
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(newProduct)
}

func (pc *ProductController) GetProductByID(ctx iris.Context) {
	idParam := ctx.Params().Get("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid product ID"})
		return
	}

	product, err := pc.Service.GetProductById(id)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"message": fmt.Sprintf("Product not found with ProductID: %s", id)})
		return
	}

	ctx.JSON(product)
}

func (pc *ProductController) GetProducts(ctx iris.Context) {
	result, err := pc.Service.GetProducts()
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Internal server error"})
		return
	}

	if result == nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(iris.Map{"message": "No products found"})
		return
	}
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(iris.Map{"message": "Products retrieved successfully", "products": result})
}

func (pc *ProductController) UpdateProduct(ctx iris.Context) {
	var productUpdate models.Product
	idParam := ctx.Params().Get("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid product ID"})
		return
	}

	if err := ctx.ReadJSON(&productUpdate); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid status"})
		return
	}

	if _, err := pc.Service.UpdateProduct(id, productUpdate); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to update status"})
		return
	}

	ctx.JSON(iris.Map{"message": "Product updated"})
}

func (pc *ProductController) GetProductDetails(ctx iris.Context) {
	var requestBody struct {
		ProductIDs []uuid.UUID `json:"product_ids"`
	}

	if err := ctx.ReadJSON(&requestBody); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid JSON request body"})
		return
	}

	if len(requestBody.ProductIDs) == 0 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "product_ids array cannot be empty"})
		return
	}

	productDetails, err := pc.Service.GetProductDetailsByIDs(requestBody.ProductIDs)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Could not retrieve product details"})
		fmt.Println("Error from service:", err)
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(productDetails)
}

func (pc *ProductController) DeleteProduct(ctx iris.Context) {
	idParam := ctx.Params().Get("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(iris.Map{"error": "Invalid product ID"})
		return
	}

	if _, err := pc.Service.DeleteProduct(id); err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(iris.Map{"error": "Failed to delete product"})
		return
	}

	ctx.JSON(iris.Map{"message": "Product Deleted"})
}
