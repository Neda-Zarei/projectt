package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/app"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/plan/domain"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/plan/port"
)

type PlanHandler struct {
	app     app.App
	service port.Service
}

func NewPlanHandler(a app.App, service port.Service) *PlanHandler {
	return &PlanHandler{app: a, service: service}
}

func (h *PlanHandler) CreatePlan(c echo.Context) error {
	var plan domain.Plan
	if err := c.Bind(&plan); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request"})
	}

	if err := h.service.CreatePlan(c.Request().Context(), &plan); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, plan)
}

func (h *PlanHandler) ListPlans(c echo.Context) error {
	page := parseQueryParamInt(c, "page", 1)
	limit := parseQueryParamInt(c, "limit", 10)
	offset := (page - 1) * limit

	plans, err := h.service.ListPlans(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to fetch plans"})
	}

	return c.JSON(http.StatusOK, plans)
}

func (h *PlanHandler) GetPlan(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid plan ID"})
	}

	plan, err := h.service.GetPlanByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Plan not found"})
	}

	return c.JSON(http.StatusOK, plan)
}

func (h *PlanHandler) UpdatePlan(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid plan ID"})
	}

	plan, err := h.service.GetPlanByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "Plan not found"})
	}

	if err := c.Bind(plan); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request"})
	}

	if err := h.service.UpdatePlan(c.Request().Context(), plan); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to update plan"})
	}

	return c.JSON(http.StatusOK, plan)
}

func (h *PlanHandler) TogglePlanActive(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid plan ID"})
	}

	if err := h.service.TogglePlanActive(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to toggle plan status"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Plan status updated successfully"})
}

func (h *PlanHandler) DeletePlan(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid plan ID"})
	}

	if err := h.service.DeletePlan(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to delete plan"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Plan deleted successfully"})
}
