package route

func (c *RouteConfig) SetupAuthRoute() {

	userRoutes := c.App.Group("/api/user")
	userRoutes.Post("/register/email", c.AuthController.RegisterByEmail)
	userRoutes.Post("/register/google", c.AuthController.RegisterByGoogleSignIn)
	userRoutes.Post("/register/phone", c.AuthController.RegisterByPhoneNumber)
	userRoutes.Post("/request-resend-email", c.AuthController.ResendEmailVerification)
	userRoutes.Post("/verify/:token", c.AuthController.VerifyEmail)

	userRoutes.Post("/login", c.AuthController.Login)
	userRoutes.Post("/request-access-token", c.AuthController.RequestAccessToken)

	userRoutes.Post("/reset-password/request", c.AuthController.RequestResetPassword)
	userRoutes.Post("/reset-password/validate", c.AuthController.ValidateResetPassword)
	userRoutes.Post("/reset-password/reset/:token", c.AuthController.ResetPassword)
}
