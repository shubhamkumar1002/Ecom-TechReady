package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"log"
	"productService/models"
	"productService/repository"
)

func CheckForCreateOrder(repo *repository.ProductRepository) {

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "steam-way-468010-t0")
	if err != nil {
		log.Fatalf("PubSub Client error: %v", err)
	}
	sub := client.Subscription("orderservicecreate-sub")
	log.Println("Product service listening for OrderCreated events...")
	go func() {
		err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
			var evt models.OrderCreatedEvent
			if err := json.Unmarshal(m.Data, &evt); err != nil {
				log.Printf("Invalid message: %v", err)
				m.Nack()
				return
			}
			log.Printf("Received CheckForCreateOrder: %+v", evt)
			totalAmount, err := repo.ReduceStockForOrder(evt.Items)
			if err != nil {
				log.Printf("Failed to update stock for OrderID %s: %v", evt.OrderID, err)
				failureEvent := models.StockUpdateFailedEvent{
					OrderID: evt.OrderID,
					Reason:  err.Error(),
				}
				SubmitFailureMessage(failureEvent)
				log.Printf("Published StockUpdateFailed event for OrderID: %s", evt.OrderID)

				m.Ack()
				return
			}

			log.Printf("Successfully updated stock for OrderID %s. Total amount: %.2f", evt.OrderID, totalAmount)

			SubmitSuccessMessage(evt)
			log.Printf("Published StockUpdatedEvent for OrderID: %s", evt.OrderID)

			m.Ack()
		})
		if err != nil {
			log.Fatalf("PubSub subscription error: %v", err)
		}
	}()
}
