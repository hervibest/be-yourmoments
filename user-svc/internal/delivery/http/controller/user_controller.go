package http

import (
	"be-yourmoments/user-svc/internal/delivery/http/middleware"
	"be-yourmoments/user-svc/internal/helper"
	"be-yourmoments/user-svc/internal/model"
	"be-yourmoments/user-svc/internal/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	GetUserProfile(ctx *fiber.Ctx) error
	UpdateUserProfile(ctx *fiber.Ctx) error
	UpdateUserProfileImage(ctx *fiber.Ctx) error
	UpdateUserCoverImage(ctx *fiber.Ctx) error
}

type userController struct {
	userUseCase     usecase.UserUseCase
	customValidator helper.CustomValidator
}

func NewUserController(userUseCase usecase.UserUseCase, customValidator helper.CustomValidator) UserController {
	return &userController{userUseCase: userUseCase, customValidator: customValidator}
}

func (c *userController) GetUserProfile(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)

	userProfileResponse, err := c.userUseCase.GetUserProfile(ctx.Context(), auth.UserId)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*model.UserProfileResponse]{
		Success: true,
		Data:    userProfileResponse,
	})
}

func (c *userController) UpdateUserProfile(ctx *fiber.Ctx) error {
	request := new(model.RequestUpdateUserProfile)
	if err := ctx.BodyParser(request); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	auth := middleware.GetUser(ctx)
	request.UserId = auth.UserId

	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validatonErrs.GetValidationErrors(),
			Message: "validation error",
		})
	}

	response, err := c.userUseCase.UpdateUserProfile(ctx.Context(), request)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[*model.UserProfileResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *userController) UpdateUserProfileImage(ctx *fiber.Ctx) error {

	userProfId := ctx.Params("userProfId", "")
	if userProfId == "" {
		return fiber.NewError(http.StatusBadRequest, "userProfId is required")
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "missing file: "+err.Error())
	}

	const maxFileSize = 5 * 1024 * 1024 // 5MB
	if file.Size > maxFileSize {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "File size exceeds the 2MB limit")
	}

	success, err := c.userUseCase.UpdateUserProfileImage(ctx.Context(), file, userProfId)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
		Data: map[string]interface{}{
			"success": success,
		},
	})
}

func (c *userController) UpdateUserCoverImage(ctx *fiber.Ctx) error {

	userProfId := ctx.Params("userProfId", "")
	if userProfId == "" {
		return fiber.NewError(http.StatusBadRequest, "userProfId is required")
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "missing file: "+err.Error())
	}

	const maxFileSize = 4 * 1024 * 1024 // 5MB
	if file.Size > maxFileSize {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "File size exceeds the 2MB limit")
	}

	success, err := c.userUseCase.UpdateUserCoverImage(ctx.Context(), file, userProfId)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusOK).JSON(model.WebResponse[any]{
		Success: true,
		Data: map[string]interface{}{
			"success": success,
		},
	})
}
