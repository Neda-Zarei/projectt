package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/app"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/api/dto"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/user/domain"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/user/port"
)

type UserHandler struct {
	app     app.App
	service port.Service
}

func NewUserHandler(a app.App, service port.Service) *UserHandler {
	return &UserHandler{app: a, service: service}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request"})
	}

	//generating a temp password (but i think we should send via email)
	tempPassword := "TempPassword123!" // todo: generate secure random password

	user := &domain.AdminUser{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := h.service.CreateUser(c.Request().Context(), user, tempPassword); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	page := parseQueryParamInt(c, "page", 1)
	limit := parseQueryParamInt(c, "limit", 10)
	offset := (page - 1) * limit

	filters := map[string]string{
		"name":  c.QueryParam("name"),
		"email": c.QueryParam("email"),
		"phone": c.QueryParam("phone"),
	}

	users, err := h.service.ListUsers(c.Request().Context(), limit, offset, filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to fetch users"})
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		}
	}

	//todo: Get total count for pagination
	total := len(userResponses) //should be actual total count

	response := dto.ListUsersResponse{
		Users: userResponses,
		Pagination: dto.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}

	return c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid user ID"})
	}

	user, err := h.service.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "User not found"})
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid user ID"})
	}

	user, err := h.service.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{"error": "User not found"})
	}

	var updateData struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request"})
	}

	if updateData.FirstName != "" {
		user.FirstName = updateData.FirstName
	}
	if updateData.LastName != "" {
		user.LastName = updateData.LastName
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}

	if err := h.service.UpdateUser(c.Request().Context(), user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to update user"})
	}

	response := dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *UserHandler) ToggleUserActive(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid user ID"})
	}

	if err := h.service.ToggleUserActive(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to toggle user status"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "User status updated successfully"})
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid user ID"})
	}

	if err := h.service.DeleteUser(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to delete user"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "User deleted successfully"})
}
