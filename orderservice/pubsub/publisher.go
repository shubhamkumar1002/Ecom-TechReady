package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"log"
	"orderService/model"
	"os"
)

func SubmitCreateMessage(paymentEvent model.OrderCreatedEvent) {
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")
	orderCreateEvent := os.Getenv("ORDER_CREATE_EVENT")
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("PubSub Client error: %v", err)
	}
	defer client.Close()

	topic := client.Topic(orderCreateEvent)
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

func SubmitUpdateMessage(paymentEvent model.PaymentEvent) {
	ctx := context.Background()
	projectID := os.Getenv("PROJECT_ID")
	orderUpdateEvent := os.Getenv("ORDER_UPDATE_EVENT")
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("PubSub Client error: %v", err)
	}
	defer client.Close()

	topic := client.Topic(orderUpdateEvent)
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
