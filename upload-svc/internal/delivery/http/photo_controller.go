package http

import (
	"be-yourmoments/upload-svc/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PhotoController interface {
	UploadPhoto(ctx *fiber.Ctx) error
	PhotoRoute(app *fiber.App)
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

	priceStr := ctx.FormValue("price", "0")
	price, _ := strconv.Atoi(priceStr)

	err = c.photoUsecase.UploadPhoto(ctx.UserContext(), file, priceStr, price)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
	})

}
