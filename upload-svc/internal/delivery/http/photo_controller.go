package http

import (
	"be-yourmoments/upload-svc/internal/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type PhotoController interface {
	UploadPhoto(ctx *fiber.Ctx) error
	Route(app *fiber.App)
}

type photoController struct {
	photoUsecase usecase.PhotoUsecase
}

func NewPhotoController(photoUsecase usecase.PhotoUsecase) PhotoController {
	return &photoController{
		photoUsecase: photoUsecase,
	}
}

func (c *photoController) UploadPhoto(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("photo")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid photo")
	}

	err = c.photoUsecase.UploadPhoto(ctx.UserContext(), file)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
	})

}
