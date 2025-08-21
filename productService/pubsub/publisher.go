package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"log"
	"os"
	"productService/models"
)

func SubmitFailureMessage(paymentEvent models.StockUpdateFailedEvent) {
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")
	stockUpdateFailedEvent := os.Getenv("STOCK_UPDATE_FAILED_EVENT")
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("PubSub Client error: %v", err)
	}

	topic := client.Topic(stockUpdateFailedEvent)
	defer topic.Stop()
	messageData, err := json.Marshal(paymentEvent)
	if err != nil {
		log.Fatalf("Failed to marshal order event: %v", err)
	}
	result := topic.Publish(ctx, &pubsub.Message{
		Data: messageData,
	})

	id, err := result.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	log.Printf("Message published with ID: %s", id)
}

func SubmitSuccessMessage(orderEvent models.OrderCreatedEvent) {
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")
	stockUpdateSuccessEvent := os.Getenv("STOCK_UPDATED_EVENT_SUB")
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("PubSub Client error: %v", err)
	}

	topic := client.Topic(stockUpdateSuccessEvent)
	defer topic.Stop()
	messageData, err := json.Marshal(orderEvent)
	if err != nil {
		log.Fatalf("Failed to marshal order event: %v", err)
	}
	result := topic.Publish(ctx, &pubsub.Message{
		Data: messageData,
	})

	id, err := result.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	log.Printf("Message published with ID: %s", id)
}
