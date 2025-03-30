package http

import (
	"be-yourmoments/user-svc/internal/delivery/http/middleware"
	"be-yourmoments/user-svc/internal/helper"
	"be-yourmoments/user-svc/internal/model"
	"be-yourmoments/user-svc/internal/usecase"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	Current(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	Logout(ctx *fiber.Ctx) error
	RegisterByEmail(ctx *fiber.Ctx) error
	RegisterByGoogleSignIn(ctx *fiber.Ctx) error
	RegisterByPhoneNumber(ctx *fiber.Ctx) error
	RequestAccessToken(ctx *fiber.Ctx) error
	RequestResetPassword(ctx *fiber.Ctx) error
	ResendEmailVerification(ctx *fiber.Ctx) error
	ResetPassword(ctx *fiber.Ctx) error
	ValidateResetPassword(ctx *fiber.Ctx) error
	VerifyEmail(ctx *fiber.Ctx) error
}

type authController struct {
	authUseCase     usecase.AuthUseCase
	customValidator helper.CustomValidator
}

func NewAuthController(authUseCase usecase.AuthUseCase, customValidator helper.CustomValidator) AuthController {
	return &authController{
		authUseCase:     authUseCase,
		customValidator: customValidator,
	}
}

func (c *authController) RegisterByPhoneNumber(ctx *fiber.Ctx) error {
	request := new(model.RegisterByPhoneRequest)
	if err := ctx.BodyParser(request); err != nil {
		fiber.NewError(http.StatusBadRequest, "bad request")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validatonErrs.GetValidationErrors(),
			Message: "validation error",
		})
	}

	response, err := c.authUseCase.RegisterByPhoneNumber(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.UserResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *authController) RegisterByGoogleSignIn(ctx *fiber.Ctx) error {
	request := new(model.RegisterByGoogleRequest)
	if err := ctx.BodyParser(request); err != nil {
		fiber.NewError(http.StatusBadRequest, "bad request")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validatonErrs.GetValidationErrors(),
			Message: "validation error",
		})
	}

	response, token, err := c.authUseCase.RegisterByGoogleSignIn(ctx.Context(), request)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.UserResponse]{
		Success: true,
		Data:    response,
		Token:   token,
	})
}

func (c *authController) RegisterByEmail(ctx *fiber.Ctx) error {
	request := new(model.RegisterByEmailRequest)
	if err := ctx.BodyParser(request); err != nil {
		fiber.NewError(http.StatusBadRequest, "bad request")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validatonErrs.GetValidationErrors(),
			Message: "validation error",
		})
	}

	response, err := c.authUseCase.RegisterByEmail(ctx.Context(), request)
	if err != nil {
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			// fiberErr sekarang berisi error dari tipe *fiber.Error
			fmt.Println("Ini fiber error dengan code:", fiberErr.Code)
		} else {
			// error tidak sesuai tipe *fiber.Error
			fmt.Println("Ini bukan fiber error")
		}
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.UserResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *authController) ResendEmailVerification(ctx *fiber.Ctx) error {
	request := new(model.ResendEmailUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		fiber.NewError(http.StatusBadRequest, "bad request")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validatonErrs.GetValidationErrors(),
			Message: "validation error",
		})
	}

	if err := c.authUseCase.ResendEmailVerification(ctx.Context(), request.Email); err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}

func (c *authController) VerifyEmail(ctx *fiber.Ctx) error {
	request := new(model.VerifyEmailUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		fiber.NewError(http.StatusBadRequest, "bad request")
	}

	request.Token = ctx.Params("token")

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validatonErrs.GetValidationErrors(),
			Message: "validation error",
		})
	}

	if err := c.authUseCase.VerifyEmail(ctx.Context(), request); err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}

func (c *authController) RequestResetPassword(ctx *fiber.Ctx) error {
	request := new(model.SendResetPasswordRequest)
	if err := ctx.BodyParser(request); err != nil {
		fiber.NewError(http.StatusBadRequest, "bad request")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validatonErrs.GetValidationErrors(),
			Message: "validation error",
		})
	}

	if err := c.authUseCase.RequestResetPassword(ctx.Context(), request.Email); err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})

}

func (c *authController) ValidateResetPassword(ctx *fiber.Ctx) error {
	request := new(model.ValidateResetTokenRequest)
	if err := ctx.BodyParser(request); err != nil {
		fiber.NewError(http.StatusBadRequest, "bad request")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validatonErrs.GetValidationErrors(),
			Message: "validation error",
		})
	}

	valid, err := c.authUseCase.ValidateResetPassword(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
		Data: map[string]interface{}{
			"valid": valid,
		},
	})
}

func (c *authController) ResetPassword(ctx *fiber.Ctx) error {
	request := new(model.ResetPasswordUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		fiber.NewError(http.StatusBadRequest, "bad request")
	}
	request.Token = ctx.Params("token")

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validatonErrs.GetValidationErrors(),
			Message: "validation error",
		})
	}

	if err := c.authUseCase.ResetPassword(ctx.Context(), request); err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
	})
}

func (c *authController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		fiber.NewError(http.StatusBadRequest, "bad request")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validatonErrs.GetValidationErrors(),
			Message: "validation error",
		})
	}

	userResponse, tokenResponse, err := c.authUseCase.Login(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*model.UserResponse]{
		Success: true,
		Data:    userResponse,
		Token:   tokenResponse,
	})
}

func (c *authController) Current(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	userResponse, err := c.authUseCase.Current(ctx.Context(), auth.Email)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*model.UserResponse]{
		Success: true,
		Data:    userResponse,
	})
}

func (c *authController) RequestAccessToken(ctx *fiber.Ctx) error {
	request := new(model.AccessTokenRequest)
	if err := ctx.BodyParser(request); err != nil {
		fiber.NewError(http.StatusBadRequest, "bad request")
	}

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validatonErrs.GetValidationErrors(),
			Message: "validation error",
		})
	}
	userResponse, tokenResponse, err := c.authUseCase.AccessTokenRequest(ctx.Context(), request.RefreshToken)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*model.UserResponse]{
		Success: true,
		Data:    userResponse,
		Token:   tokenResponse,
	})
}

func (c *authController) Logout(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	request := new(model.LogoutUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		return fiber.ErrBadRequest
	}

	request.UserId = auth.UserId
	request.AccessToken = auth.Token

	valid, err := c.authUseCase.Logout(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
		Data: map[string]interface{}{
			"valid": valid,
		},
	})
}
