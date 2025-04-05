package adapter

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

type MessagingClientAdapter interface {
	Send(ctx context.Context, message *messaging.Message) (string, error)
}

func NewMessagingClientAdapter(app *firebase.App) MessagingClientAdapter {
	ctx := context.Background()
	messagingClientAdapter, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error initializing firebase messaging: %v", err)
	}

	return messagingClientAdapter
}
