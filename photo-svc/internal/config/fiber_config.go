package config

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		Prefork:      false,
		AppName:      "photo-svc",
		ErrorHandler: CustomError(),
		BodyLimit:    100 * 1024 * 1024,
	})

	app.Use(cors.ConfigDefault)
	return app
}

func CustomError() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := http.StatusInternalServerError
		if err, ok := err.(*fiber.Error); ok {
			code = err.Code
		}

		message := &Message{
			Success: false,
			Error:   err.Error(),
		}

		return ctx.Status(code).JSON(message)
	}
}

type Message struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
