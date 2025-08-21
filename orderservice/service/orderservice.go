package service

import (
	"github.com/google/uuid"
	"log"
	"orderService/model"
	"orderService/pubsub"
	"orderService/repository"
)

type OrderService struct {
	Repo *repository.OrderRepository
}

func (os *OrderService) Create(ocd *model.OrderCreateDTO, authHeader string) (*model.Order, error) {
	newOrder, err := os.Repo.Create(ocd, authHeader)
	if err != nil {
		return nil, err
	}

	log.Printf("Order %s created successfully. Publishing event.", newOrder.ID)
	eventPayload := model.OrderCreatedEvent{
		OrderID: newOrder.ID,
		Items:   ocd.Items,
		Total:   newOrder.TotalAmount,
	}
	pubsub.SubmitCreateMessage(eventPayload)
	return newOrder, nil

}
func (os *OrderService) GetOrderById(id uuid.UUID) (*model.Order, error) {
	return os.Repo.GetOrderByID(id)
}

func (os *OrderService) GetOrders() ([]model.Order, error) {
	return os.Repo.GetOrders()
}

func (os *OrderService) UpdateOrderStatus(id uuid.UUID, orderstatus string, paymentstatus string) (*model.Order, error) {
	updatedOrder, err := os.Repo.UpdateStatus(id, orderstatus)
	if err != nil {
		return nil, err
	}

	if paymentstatus != "" {
		paymentEvent := model.PaymentEvent{
			OrderID:       id,
			PaymentStatus: paymentstatus,
		}
		pubsub.SubmitUpdateMessage(paymentEvent)
	}

	return updatedOrder, nil
}
