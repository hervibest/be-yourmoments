package adapter

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

type FirestoreClientAdapter interface {
	Collection(path string) *firestore.CollectionRef
	BulkWriter(ctx context.Context) *firestore.BulkWriter
	Close() error
}

func NewFirestoreClientAdapter(app *firebase.App) FirestoreClientAdapter {
	ctx := context.Background()
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing firebase firestore: %v", err)
	}

	return firestoreClient
}
