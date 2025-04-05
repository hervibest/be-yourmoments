package adapter

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

type AuthClientAdapter interface {
	CreateOIDCProviderConfig(ctx context.Context, config *auth.OIDCProviderConfigToCreate) (*auth.OIDCProviderConfig, error)
	CustomToken(ctx context.Context, uid string) (string, error)
}

func NewAuthClientAdapter(app *firebase.App) AuthClientAdapter {
	ctx := context.Background()
	authClientAdapter, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error initializing firebase auth: %v", err)
	}

	return authClientAdapter
}
