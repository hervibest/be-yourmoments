package http

import (
	"be-yourmoments/internal/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type PhotoController interface{}
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

	err1, _ := c.photoUsecase.UploadPhoto(ctx.UserContext(), file)
	if err1 != nil {
		return err1
	}

	return nil

}
