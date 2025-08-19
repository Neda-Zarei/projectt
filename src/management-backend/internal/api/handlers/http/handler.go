package http

import (
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/app"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/adapter/repository"
	mw "hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/api/middleware"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/plan"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/plan/port"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/user"
	userport "hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/user/port"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	app         app.App
	echo        *echo.Echo
	auth        *AuthHandler
	user        *UserHandler
	plan        *PlanHandler
	userService userport.Service
	planService port.Service
}

func NewHandler(a app.App) *Handler {
	e := echo.New()

	userRepo := repository.NewUserRepository(a.Config().DB) // You'll need to implement this
	userService := user.NewService(userRepo)

	planRepo := repository.NewPlanRepository(a.Config().DB) // You'll need to implement this
	planService := plan.NewService(planRepo)

	return &Handler{
		app:         a,
		echo:        e,
		userService: userService,
		planService: planService,
		auth:        NewAuthHandler(a, userService),
		user:        NewUserHandler(a, userService),
		plan:        NewPlanHandler(a, planService),
	}
}

func (h *Handler) SetupRoutes() *echo.Echo {
	e := h.echo
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
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
