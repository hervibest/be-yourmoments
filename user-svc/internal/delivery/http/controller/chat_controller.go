package http

import (
	"be-yourmoments/user-svc/internal/helper"
	"be-yourmoments/user-svc/internal/model"
	"be-yourmoments/user-svc/internal/usecase"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type ChatController interface {
	GetCustomToken(ctx *fiber.Ctx) error
	GetOrCreateRoom(ctx *fiber.Ctx) error
	SendMessage(ctx *fiber.Ctx) error
}

type chatController struct {
	chatUseCase     usecase.ChatUseCase
	customValidator helper.CustomValidator
}

func NewChatController(chatUseCase usecase.ChatUseCase, customValidator helper.CustomValidator) ChatController {
	return &chatController{chatUseCase: chatUseCase, customValidator: customValidator}
}

func (c *chatController) GetCustomToken(ctx *fiber.Ctx) error {
	request := new(model.RequestCustomToken)
	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validatonErrs.GetValidationErrors(),
			Message: "validation error",
		})
	}

	response, err := c.chatUseCase.GetCustomToken(ctx.Context(), request)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.CustomTokenResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *chatController) GetOrCreateRoom(ctx *fiber.Ctx) error {
	request := new(model.RequestGetOrCreateRoom)
	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validatonErrs.GetValidationErrors(),
			Message: "validation error",
		})
	}

	response, err := c.chatUseCase.GetOrCreateRoom(ctx.Context(), request)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.GetOrCreateRoomResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *chatController) SendMessage(ctx *fiber.Ctx) error {
	request := new(model.RequestSendMessage)
	if validatonErrs := c.customValidator.ValidateUseCase(request); validatonErrs != nil {
		return ctx.Status(http.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validatonErrs.GetValidationErrors(),
			Message: "validation error",
		})
	}

	if err := c.chatUseCase.SendMessage(ctx.Context(), request); err != nil {
		fmt.Println(err)
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(model.WebResponse[*model.GetOrCreateRoomResponse]{
		Success: true,
	})
}
