package route

func (c *RouteConfig) SetupUserRoute() {

	userRoutes := c.App.Group("/api/users", c.AuthMiddleware)
	userRoutes.Get("/current", c.AuthController.Current)
	userRoutes.Delete("/logout", c.AuthController.Logout)

	userRoutes.Get("/profile", c.UserController.GetUserProfile)
	userRoutes.Put("/profile", c.UserController.UpdateUserProfile)
	userRoutes.Patch("/profile/:userProfId", c.UserController.UpdateUserProfileImage)
	userRoutes.Patch("/profile/cover/:userProfId", c.UserController.UpdateUserCoverImage)

	userRoutes.Post("/message", c.UserController.GetUserProfile)

	// userRoutes.Get("/users", listUsers)
	userRoutes.Post("/room", c.ChatController.GetOrCreateRoom)
	userRoutes.Get("/token/:uid", c.ChatController.GetCustomToken)
	userRoutes.Post("/send-message", c.ChatController.SendMessage)
	// userRoutes.Post("/send-notification", sendNotification) // route baru untuk FCM

}
