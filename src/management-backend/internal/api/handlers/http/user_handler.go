package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/admin/domain"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/admin/port"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/api/dto"
)

type UserHandler struct {
	service port.Service
}

func NewUserHandler(s port.Service) *UserHandler {
	return &UserHandler{service: s}
}

// @Summary      Create new user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user  body  dto.CreateUserRequest  true  "User object"
// @Success      201  {object}  dto.UserResponse
// @Failure      default  {object}  dto.Error
// @Router       /users [post]
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

// @Summary      Get list of users (paginated + filter)
// @Tags         user
// @Produce      json
// @Param        page   query  int    false  "Page number"
// @Param        size   query  int    false  "Page size"
// @Param        name   query  string false  "Filter by name"
// @Param        email  query  string false  "Filter by email"
// @Param        phone  query  string false  "Filter by phone"
// @Success      200  {object}  dto.ListUsersResponse
// @Failure      default  {object}  dto.Error
// @Router       /users [get]
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

// @Summary      Get user by ID
// @Tags         user
// @Produce      json
// @Param        id  path  string  true  "User ID"
// @Success      200  {object}  dto.UserResponse
// @Failure      default  {object}  dto.Error
// @Router       /users/{id} [get]
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

// @Summary      Update user info
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id    path  string      true  "User ID"
// @Param        user  body  dto.UserResponse true  "User object"
// @Success      200  {object}  dto.UserResponse
// @Failure      default  {object}  dto.Error
// @Router       /users/{id} [put]
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

// @Summary      Activate/Deactivate user account
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "User ID"
// @Param        active body dto.ToggleUserActiveRequest  true  "Activation status"
// @Success      200  {string}  string  "Status updated"
// @Failure      default  {object}  dto.Error
// @Router       /users/{id} [patch]
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
