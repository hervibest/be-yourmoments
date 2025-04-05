package route

import (
	http "be-yourmoments/user-svc/internal/delivery/http/controller"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	AuthController http.AuthController
	UserController http.UserController
	ChatController http.ChatController
	AuthMiddleware fiber.Handler
}

func (r *RouteConfig) Setup() {
	r.SetupAuthRoute()
	r.SetupUserRoute()
}
