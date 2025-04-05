package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func NewFirebaseConfig() *firebase.App {
	ctx := context.Background()
	firebaseApp, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile("../../serviceAccountKey.json"))
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}

	return firebaseApp
}
