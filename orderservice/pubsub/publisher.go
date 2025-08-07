package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"log"
	"orderService/model"
)

func SubmitCreateMessage(paymentEvent model.OrderCreatedEvent) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "steam-way-468010-t0")
	if err != nil {
		log.Fatalf("PubSub Client error: %v", err)
	}
	defer client.Close()

	// Specify the existing topic name.
	topic := client.Topic("orderservicecreate")
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
	client, err := pubsub.NewClient(ctx, "steam-way-468010-t0")
	if err != nil {
		log.Fatalf("PubSub Client error: %v", err)
	}
	defer client.Close()

	// Specify the existing topic name.
	topic := client.Topic("orderserviceUpdate")
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
