package adapter

import (
	"be-yourmoments/user-svc/internal/helper/utils"
	"be-yourmoments/user-svc/internal/model"
	"context"
	"errors"

	"google.golang.org/api/idtoken"
)

type GoogleTokenAdapter interface {
	ValidateGoogleToken(ctx context.Context, token string) (*model.GoogleSignInClaim, error)
	GetClientId() string
}

type googleTokenAdapter struct {
	clientId string
}

func NewGoogleTokenAdapter() GoogleTokenAdapter {
	clientId := utils.GetEnv("GOOGLE_CLIENT_ID")
	return &googleTokenAdapter{
		clientId: clientId,
	}
}

func (a *googleTokenAdapter) ValidateGoogleToken(ctx context.Context, token string) (*model.GoogleSignInClaim, error) {
	payload, err := idtoken.Validate(context.Background(), token, a.clientId)
	if err != nil {
		return nil, errors.New("invalid google token")
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		return nil, errors.New("invalid email claims")

	}

	profilePictureUrl, ok := payload.Claims["picture"].(string)
	if !ok {
		return nil, errors.New("invalid url picture claims")

	}

	givenName, ok := payload.Claims["given_name"].(string)
	if !ok {
		return nil, errors.New("invalid url picture claims")

	}

	googleId, ok := payload.Claims["sub"].(string)
	if !ok {
		return nil, errors.New("invalid google id claims")

	}

	googleClaims := &model.GoogleSignInClaim{
		Email:             email,
		Username:          givenName,
		ProfilePictureUrl: profilePictureUrl,
		GoogleId:          googleId,
	}

	return googleClaims, nil
}

func (a *googleTokenAdapter) GetClientId() string {
	return a.clientId
}
