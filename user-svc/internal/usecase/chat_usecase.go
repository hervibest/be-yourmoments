package usecase

import (
	"be-yourmoments/user-svc/internal/adapter"
	"be-yourmoments/user-svc/internal/model"
	"context"
	"html"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

type ChatUseCase interface {
	GetCustomToken(ctx context.Context, req *model.RequestCustomToken) (*model.CustomTokenResponse, error)
	GetOrCreateRoom(ctx context.Context, req *model.RequestGetOrCreateRoom) (*model.GetOrCreateRoomResponse, error)
	SendMessage(ctx context.Context, req *model.RequestSendMessage) error
}

type chatUseCase struct {
	firestoreClientAdapter adapter.FirestoreClientAdapter
	authClientAdapter      adapter.AuthClientAdapter
	messagingClientAdapter adapter.MessagingClientAdapter
	perspectiveAdapter     adapter.PerspectiveAdapter
}

func NewChatUseCase(firestoreClientAdapter adapter.FirestoreClientAdapter, authClientAdapter adapter.AuthClientAdapter,
	messagingClientAdapter adapter.MessagingClientAdapter, perspectiveAdapter adapter.PerspectiveAdapter) ChatUseCase {
	return &chatUseCase{
		firestoreClientAdapter: firestoreClientAdapter,
		authClientAdapter:      authClientAdapter,
		messagingClientAdapter: messagingClientAdapter,
		perspectiveAdapter:     perspectiveAdapter,
	}
}

func (u *chatUseCase) GetOrCreateRoom(ctx context.Context, req *model.RequestGetOrCreateRoom) (*model.GetOrCreateRoomResponse, error) {
	roomUserId := generateRoomId(req.UserA, req.UserB)
	roomRef := u.firestoreClientAdapter.Collection("rooms").Doc(roomUserId)
	created := false
	var roomId string

	docSnap, err := roomRef.Get(ctx)
	if err != nil {
		roomId = ulid.Make().String()
		participants := []string{req.UserA, req.UserB}
		_, err := roomRef.Set(ctx, map[string]interface{}{
			"roomUserId":   roomUserId,
			"roomId":       roomId,
			"participants": participants,
			"createdAt":    firestore.ServerTimestamp,
		})
		if err != nil {
			return nil, fiber.ErrInternalServerError
		}
		created = true
	} else {
		data := docSnap.Data()
		if rn, ok := data["roomId"].(string); ok {
			roomId = rn
		} else {
			roomId = ulid.Make().String()
		}
	}

	response := &model.GetOrCreateRoomResponse{
		RoomId:  roomId,
		Created: created,
	}

	return response, nil
}

func generateRoomId(userA, userB string) string {
	if userA < userB {
		return userA + "_" + userB
	}
	return userB + "_" + userA
}

func (u *chatUseCase) GetCustomToken(ctx context.Context, req *model.RequestCustomToken) (*model.CustomTokenResponse, error) {
	token, err := u.authClientAdapter.CustomToken(ctx, req.UserId)
	if err != nil {
		return nil, fiber.NewError(http.StatusBadRequest)
	}

	response := &model.CustomTokenResponse{
		Token: token,
	}

	return response, nil
}

func (u *chatUseCase) SendMessage(ctx context.Context, req *model.RequestSendMessage) error {
	trimmed := strings.TrimSpace(req.Message)
	if trimmed == "" {
		return fiber.NewError(http.StatusUnprocessableEntity, "Empty message not allowed")
	}

	if len(trimmed) > 500 {
		return fiber.NewError(http.StatusUnprocessableEntity, "Message too long")
	}

	isToxic, err := u.perspectiveAdapter.IsToxicMessage(trimmed)
	if err != nil {
		log.Printf("error checking toxicity: %v", err)
	}

	if isToxic {
		return fiber.NewError(http.StatusUnprocessableEntity, "Toxic content detected")
	}

	safeMessage := html.EscapeString(trimmed)
	_, _, err = u.firestoreClientAdapter.
		Collection("rooms").
		Doc(req.RoomID).
		Collection("messages").
		Add(ctx, map[string]interface{}{
			"senderId":  req.SenderID,
			"message":   safeMessage,
			"timestamp": firestore.ServerTimestamp,
		})

	if err != nil {
		return fiber.ErrInternalServerError
	}

	return nil
}
