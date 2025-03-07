package http

import (
	"be-yourmoments/photo-svc/internal/config"

	"github.com/gofiber/fiber/v2"
)

func (c *photoController) Route(app *fiber.App) {
	api := app.Group(config.EndpointPrefix)
	api.Post("/:id", c.UploadPhoto)

}
