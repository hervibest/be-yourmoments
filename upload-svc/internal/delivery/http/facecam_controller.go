package http

import (
	"be-yourmoments/upload-svc/internal/usecase"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type FacecamController interface {
	UploadFacecam(ctx *fiber.Ctx) error
	FacecamRoute(app *fiber.App)
}

type facecamController struct {
	facecamUseCase usecase.FacecamUseCase
}

func NewFacecamController(facecamUseCase usecase.FacecamUseCase) FacecamController {
	return &facecamController{
		facecamUseCase: facecamUseCase,
	}
}

func (c *facecamController) UploadFacecam(ctx *fiber.Ctx) error {
	log.Println("Upload facecam via http")
	file, err := ctx.FormFile("facecam")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid facecam")
	}

	userId := ctx.FormValue("userId", "xx-yy-zz")
	fmt.Println(userId)

	err = c.facecamUseCase.UploadFacecam(ctx.UserContext(), file, userId)
	if err != nil {
		return err
	}

	return ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"success": true,
	})

}
