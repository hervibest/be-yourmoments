package route

func (c *RouteConfig) SetupUserRoute() {

	userRoutes := c.App.Group("/api/users", c.AuthMiddleware)
	userRoutes.Get("/current", c.AuthController.Current)
	userRoutes.Delete("/logout", c.AuthController.Logout)

}
