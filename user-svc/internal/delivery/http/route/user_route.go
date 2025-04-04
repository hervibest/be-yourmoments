package route

func (c *RouteConfig) SetupUserRoute() {

	userRoutes := c.App.Group("/api/users", c.AuthMiddleware)
	userRoutes.Get("/current", c.AuthController.Current)
	userRoutes.Delete("/logout", c.AuthController.Logout)

	userRoutes.Get("/profile", c.UserController.GetUserProfile)
	userRoutes.Put("/profile", c.UserController.UpdateUserProfile)
	userRoutes.Patch("/profile/:userProfId", c.UserController.UpdateUserProfileImage)
	userRoutes.Patch("/profile/cover/:userProfId", c.UserController.UpdateUserCoverImage)
}
