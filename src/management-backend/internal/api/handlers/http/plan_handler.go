package http

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/plan/domain"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/plan/port"
)

var Validate = validator.New()

type PlanHandler struct {
	service port.Service
}

func NewPlanHandler(s port.Service) *PlanHandler {
	return &PlanHandler{service: s}
}

func (h *PlanHandler) CreatePlan(c echo.Context) error {
	var plan domain.Plan
	if err := c.Bind(&plan); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request"})
	}
	if err := Validate.Struct(plan); err != nil {
		validationErrors := err.(validator.ValidationErrors)

		errors := make(map[string]string)
		for _, e := range validationErrors {
			errors[e.Field()] = fmt.Sprintf("failed on '%s' validation", e.Tag())
		}

		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":  "validation failed",
			"fields": errors,
		})
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

// @Summary      Get active plans of a user
// @Tags         plan
// @Produce      json
// @Param        id  path  string  true  "User ID"
// @Success      200  {array}  dto.PlanResponse
// @Failure      default  {object}  dto.Error
// @Router       /users/{id}/plans [get]
func (h *PlanHandler) UserPlan(c echo.Context) error { return nil }

// @Summary      Assign new plan to user
// @Tags         plan
// @Accept       json
// @Produce      json
// @Param        id    path  string     true  "User ID"
// @Param        plan  body  dto.PlanResponse true  "Plan object"
// @Success      201  {object}  dto.PlanResponse
// @Failure      default  {object}  dto.Error
// @Router       /users/{id}/plans [post]
func (h *PlanHandler) AssignPlan(c echo.Context) error { return nil }

// @Summary      Renew a user's plan
// @Tags         plan
// @Produce      json
// @Param        id      path  string  true  "User ID"
// @Param        planId  path  string  true  "Plan ID"
// @Success      200  {string}  string  "Plan renewed"
// @Failure      default  {object}  dto.Error
// @Router       /users/{id}/plans/{planId}/renew [post]
func (h *PlanHandler) RenewPlan(c echo.Context) error { return nil }
