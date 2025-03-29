package http

import (
	"be-yourmoments/upload-svc/internal/config"

	"github.com/gofiber/fiber/v2"
)

func (c *photoController) PhotoRoute(app *fiber.App) {
	// api := app.Group(config.EndpointPrefix)
	// // api.Post("/single", c.UploadPhoto)
}

func (c *facecamController) FacecamRoute(app *fiber.App) {
	api := app.Group(config.EndpointPrefix)
	api.Post("/facecam/single", c.UploadFacecam)
}
