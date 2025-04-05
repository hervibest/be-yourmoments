package model

type RequestGetOrCreateRoom struct {
	UserA string `json:"user_a" validate:"required"`
	UserB string `json:"user_B" validate:"required"`
}

type GetOrCreateRoomResponse struct {
	RoomId  string `json:"room_id"`
	Created bool   `json:"created"`
}

type RequestCustomToken struct {
	UserId string `json:"user_id" validate:"required"`
}

// TODO memastikan snake_case bukan pascalCase
type RequestSendMessage struct {
	RoomID   string `json:"room_id" validate:"required"`
	SenderID string `json:"sender_id" validate:"required"`
	Message  string `json:"message" validate:"required"`
}

type CustomTokenResponse struct {
	Token string `json:"token"`
}
