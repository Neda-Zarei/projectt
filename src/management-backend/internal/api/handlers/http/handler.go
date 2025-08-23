package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/app"
	mw "hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/api/middleware"
)

type Handler struct {
	app  app.App
	echo *echo.Echo
	auth *AuthHandler
	user *UserHandler
	plan *PlanHandler
}

func NewHandler(a app.App) *Handler {
	return &Handler{
		app:  a,
		echo: echo.New(),
		auth: NewAuthHandler(a.UserService(), a.Config().JWT),
		user: NewUserHandler(a.UserService()),
		plan: NewPlanHandler(a.PlanService()),
	}
}

func (h *Handler) SetupRoutes() *echo.Echo {
	e := h.echo
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())
	e.Use(mw.Logger(h.app.Logger()))
	e.Use(middleware.Recover())

	//public routes
	e.POST("/api/auth/login", h.auth.Login)

	//protected routes
	api := e.Group("/api")
	api.Use(mw.NewAuthMiddleware(h.app).ValidateJWT())

	//user routes
	api.GET("/users", h.user.ListUsers)
	api.POST("/users", h.user.CreateUser)
	api.GET("/users/:id", h.user.GetUser)
	api.PUT("/users/:id", h.user.UpdateUser)
	api.PATCH("/users/:id/toggle-active", h.user.ToggleUserActive)
	api.DELETE("/users/:id", h.user.DeleteUser)

	//plan routes
	api.GET("/plans", h.plan.ListPlans)
	api.POST("/plans", h.plan.CreatePlan)
	api.GET("/plans/:id", h.plan.GetPlan)
	api.PUT("/plans/:id", h.plan.UpdatePlan)
	api.PATCH("/plans/:id/toggle-active", h.plan.TogglePlanActive)
	api.DELETE("/plans/:id", h.plan.DeletePlan)

	return e
}
