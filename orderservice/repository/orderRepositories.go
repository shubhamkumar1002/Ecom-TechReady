package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"orderService/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (ord *OrderRepository) Create(orderCreateDTO *model.OrderCreateDTO) (*model.Order, error) {

	if len(orderCreateDTO.Items) == 0 {
		return nil, errors.New("cannot create an order with zero items")
	}
	var totalPrice int
	var productIDsAsString []string
	for _, item := range orderCreateDTO.Items {
		var itemTotalprice int
		productIDsAsString = append(productIDsAsString, item.ProductID.String())
		if item.Quantity > 0 {
			itemTotalprice = item.Price * item.Quantity
		}
		totalPrice = totalPrice + itemTotalprice
	}
	type ProductDetailsRequest struct {
		ProductIDs []string `json:"product_ids"`
	}
	requestPayload := ProductDetailsRequest{ProductIDs: productIDsAsString}
	jsonBody, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request body for products service: %w", err)
	}

	productServiceURL := "http://product-service.default.svc.cluster.local/product/details"
	resp, err := http.Post(productServiceURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to contact products service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("products service returned an error: %s", resp.Status)
	}

	type ProductDetailsResponse struct {
		ID       string  `json:"id"`
		Price    float64 `json:"price"`
		Quantity int     `json:"quantity"`
	}
	var productDetails []ProductDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&productDetails); err != nil {
		return nil, fmt.Errorf("failed to decode response from products service: %w", err)
	}

	productMap := make(map[string]ProductDetailsResponse)
	for _, p := range productDetails {
		productMap[p.ID] = p
	}

	var finalOrderItems []model.OrderItem
	var totalAmount float64

	for _, requestedItem := range orderCreateDTO.Items {
		product, found := productMap[requestedItem.ProductID.String()]
		if !found {
			return nil, fmt.Errorf("product with ID '%s' was not found", requestedItem.ProductID.String())
		}

		if product.Quantity < requestedItem.Quantity {
			return nil, fmt.Errorf("not enough stock for product ID '%s'", requestedItem.ProductID.String())
		}

		finalOrderItems = append(finalOrderItems, model.OrderItem{
			ProductID: requestedItem.ProductID.String(),
			Quantity:  requestedItem.Quantity,
			UnitPrice: product.Price,
		})

		totalAmount += product.Price * float64(requestedItem.Quantity)
	}

	var newOrder *model.Order
	err = ord.DB.Transaction(func(tx *gorm.DB) error {
		orderID := uuid.New()
		newOrder = &model.Order{
			ID:          orderID,
			UserID:      orderCreateDTO.UserID,
			TotalAmount: totalAmount,
			OrderStatus: model.OrderPlaced,
		}

		if err := tx.Create(newOrder).Error; err != nil {
			return err
		}

		for i := range finalOrderItems {
			finalOrderItems[i].OrderID = orderID
		}
		if err := tx.Create(&finalOrderItems).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to save order to database: %w", err)
	}

	newOrder.OrderItems = finalOrderItems
	return newOrder, nil
}

func (ord *OrderRepository) GetOrderByID(id uuid.UUID) (*model.Order, error) {
	var order model.Order
	if err := ord.DB.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (ord *OrderRepository) GetOrders() ([]model.Order, error) {
	var orders []model.Order
	err := ord.DB.Preload("OrderItems").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (ord *OrderRepository) UpdateStatus(id uuid.UUID, orderstatus string) (*model.Order, error) {
	var updateOrder model.Order
	if orderstatus != "" {
		err := ord.DB.Model(&updateOrder).Where("id = ?", id).Update("order_status", orderstatus).Update("updated_at", time.Now()).Error
		if err != nil {
			return nil, err
		}
	}

	return &updateOrder, nil
}
