package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"log"
	"orderService/model"
	"orderService/repository"
)

func CheckForFailedOrderEvent(repo *repository.OrderRepository) {

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "steam-way-468010-t0")
	if err != nil {
		log.Fatalf("PubSub Client error: %v", err)
	}
	sub := client.Subscription("StockUpdateFailedEvent-sub")

	go func() {
		err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
			var evt model.PaymentEvent
			if err := json.Unmarshal(m.Data, &evt); err != nil {
				log.Printf("Invalid message: %v", err)
				m.Nack()
				return
			}

			log.Printf("Received CheckForFailedOrderEvent: %+v", evt)

			orderExist, err := repo.GetOrderByID(evt.OrderID)
			if orderExist != nil {
				_, err = repo.UpdateStatus(evt.OrderID, "FAILED")
				if err != nil {
					log.Printf("Payment status update failed: %v", err)
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
