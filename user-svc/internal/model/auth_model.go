package model

import (
	"time"
)

type RegisterByPhoneRequest struct {
	Username    string     `json:"username" validate:"required,max=100"`
	Password    string     `json:"password" validate:"required"`
	PhoneNumber string     `json:"phone_number" validate:"required,min=10,max=15"`
	BirthDate   *time.Time `json:"birth_date" validate:"required"`
	//TODO
}

type RegisterByGoogleRequest struct {
	Token string `json:"token" validate:"required"`
}

type GoogleSignInClaim struct {
	Email             string `json:"email" validate:"required,max=255"`
	Username          string `json:"username" validate:"required,max=255"`
	ProfilePictureUrl string `json:"picture" validate:"required"`
	GoogleId          string `json:"sub" validate:"required,max=15"`
}

type RegisterByEmailRequest struct {
	Username  string     `json:"username" validate:"required,max=100"`
	Email     string     `json:"email" validate:"required,max=255"`
	Password  string     `json:"password" validate:"required"`
	BirthDate *time.Time `json:"birth_date" validate:"required"`
}

type ResendEmailUserRequest struct {
	Email string `json:"email" validate:"required,email,max=100"`
}

type VerifyEmailUserRequest struct {
	Email string `json:"email" validate:"required,email,max=100"`
	Token string `validate:"required"`
}

type SendResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email,max=100"`
}

type ValidateResetTokenRequest struct {
	Email string `json:"email" validate:"required,email,max=100"`
	Token string `json:"token" validate:"required"`
}

type ResetPasswordUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,max=100"`
	Token    string `validate:"required"`
}

type LoginUserRequest struct {
	MultipleParam string `json:"multiple_param" validate:"required,max=100"`
	Password      string `json:"password" validate:"required,max=100"`
}

type VerifyUserRequest struct {
	Token string `validate:"required"`
}

type AuthResponse struct {
	UserId      string
	Username    string
	Email       string
	PhoneNumber string
	Token       string
	ExpiresAt   time.Time
}

type LogoutUserRequest struct {
	UserId       string
	AccessToken  string
	ExpiresAt    time.Time
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type AccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// type UpdateUserRequest struct {
// 	Id         string   `json:"id"`
// 	Name       string   `json:"name"`
// 	Categories []string `json:"categories"`
// }

// type GetUserRequest struct {
// 	Id string `json:"id"`
// }

// type SearchUserRequest struct {
// 	Name  string `json:"name" validate:"max=100"`
// 	Email string `json:"email" validate:"max=100"`
// 	Page  int    `json:"page" validate:"min=1"`
// 	Size  int    `json:"size" validate:"min=1,max=100"`
// }

type UserResponse struct {
	Id                    string     `json:"id"`
	Username              string     `json:"username"`
	Email                 string     `json:"email"`
	EmailVerifiedAt       *time.Time `json:"email_verified_at,omitempty"`
	PhoneNumber           string     `json:"phone_number"`
	PhoneNumberVerifiedAt *time.Time `json:"phone_number_verified_at,omitempty"`
	GoogleId              string     `json:"google_id"`
	CreatedAt             *time.Time `json:"created_at,omitempty"`
	UpdatedAt             *time.Time `json:"updated_at,omitempty"`
}

// type UserResponse struct {
// 	Id         string                  `json:"id"`
// 	Name       string                  `json:"name"`
// 	Email      string                  `json:"email,omitempty"`
// 	Categories *[]*UserCatUserResponse `json:"categories,omitempty"`
// 	// UserDetail *UserDetailResponse     `json:"user_detail,omitempty"`
// 	CreatedAt *time.Time `json:"created_at,omitempty"`
// 	UpdatedAt *time.Time `json:"updated_at"`
// }

// type UserCatUserResponse struct {
// 	UserId         string     `json:"user_id"`
// 	UserCategoryId string     `json:"user_category_id"`
// 	Name           string     `json:"name"`
// 	Description    string     `json:"description"`
// 	CreatedAt      *time.Time `json:"created_at,omitempty"`
// 	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
// }
