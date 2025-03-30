package middleware

import (
	"be-yourmoments/user-svc/internal/helper"
	"be-yourmoments/user-svc/internal/model"
	"be-yourmoments/user-svc/internal/usecase"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// TODO SHOULD TOKEN VALIDATED ?
func NewUserAuth(authUseCase usecase.AuthUseCase, validator helper.CustomValidator) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := strings.TrimPrefix(ctx.Get("Authorization", ""), "Bearer ")
		if token == "" || token == "NOT_FOUND" {
			return fiber.NewError(fiber.ErrUnauthorized.Code, "Unauthorized access")
		}

		request := &model.VerifyUserRequest{Token: token}
		if validationErr := validator.ValidateUseCase(request); validationErr != nil {
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
				Success: false,
				Errors:  validationErr,
				Message: "validation error",
			})
		}

		auth, err := authUseCase.Verify(ctx.UserContext(), request)
		if err != nil {
			return err
		}

		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.AuthResponse {
	return ctx.Locals("auth").(*model.AuthResponse)
}
