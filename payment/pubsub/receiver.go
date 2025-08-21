package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"log"
	"os"
	"paymentService/model"
	"paymentService/repository"
)

func CheckForPublishedPayments(repo *repository.PaymentRepository) {
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")
	stockUpdatedEventSub := os.Getenv("STOCK_UPDATED_EVENT_SUB")
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("PubSub Client error: %v", err)
	}
	sub := client.Subscription(stockUpdatedEventSub)

	go func() {
		err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
			var evt model.PaymentEvent
			if err := json.Unmarshal(m.Data, &evt); err != nil {
				log.Printf("Invalid message: %v", err)
				m.Nack()
				return
			}

			log.Printf("Received CheckForPublishedPayments: %+v", evt)

			orderExist, err := repo.CheckOrder(evt.OrderID)
			if orderExist {
				err = repo.UpdateStatus(evt.OrderID, evt.PaymentStatus)
				if err != nil {
					log.Printf("Payment status update failed: %v", err)
					m.Nack()
					return
				}
			} else {
				newPayment := &model.PaymentCreateDTO{
					OrderID:       evt.OrderID,
					TotalAmount:   evt.TotalAmount,
					PaymentStatus: model.PaymentStatus(evt.PaymentStatus),
				}
				_, err = repo.CreatePayment(newPayment)
				if err != nil {
					log.Printf("Payment creation failed: %v", err)
					m.Nack()
					return
				}
			}

			m.Ack()
		})
		if err != nil {
			log.Fatalf("PubSub subscription error: %v", err)
		}
	}()
}

func CheckForSuccessOrderEvent(repo *repository.PaymentRepository) {
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")
	orderUpdateEventSub := os.Getenv("ORDER_UPDATE_EVENT_SUB")
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("PubSub Client error: %v", err)
	}
	sub := client.Subscription(orderUpdateEventSub)

	go func() {
		err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
			var evt model.PaymentEvent
			if err := json.Unmarshal(m.Data, &evt); err != nil {
				log.Printf("Invalid message: %v", err)
				m.Nack()
				return
			}

			log.Printf("Received CheckForSuccessOrderEvent : %+v", evt)

			orderExist, err := repo.CheckOrder(evt.OrderID)
			if orderExist {
				err = repo.UpdateStatus(evt.OrderID, evt.PaymentStatus)
				if err != nil {
					log.Printf("Payment status update failed: %v", err)
					m.Nack()
					return
				}
			} else {
				newPayment := &model.PaymentCreateDTO{
					OrderID:       evt.OrderID,
					TotalAmount:   evt.TotalAmount,
					PaymentStatus: model.PaymentStatus(evt.PaymentStatus),
				}
				_, err = repo.CreatePayment(newPayment)
				if err != nil {
					log.Printf("Payment creation failed: %v", err)
					m.Nack()
					return
				}
			}
			m.Ack()
		})
		if err != nil {
			log.Fatalf("PubSub subscription error: %v", err)
		}
	}()
}
