package http

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/app"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/api/dto"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/user/port"
)

type AuthHandler struct {
	app     app.App
	service port.Service
}

func NewAuthHandler(a app.App, service port.Service) *AuthHandler {
	return &AuthHandler{app: a, service: service}
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request"})
	}

	// todo: add Arcaptcha verification here
	// if !verifyArcaptcha(req.CaptchaToken) {
	//     return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid captcha"})
	// }

	user, err := h.service.Authenticate(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Invalid credentials"})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = user.ID
	claims["email"] = user.Email
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(h.app.Config().JWT.Expiration)).Unix()

	tokenString, err := token.SignedString([]byte(h.app.Config().JWT.Secret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to generate token"})
	}

	response := dto.LoginResponse{
		Token: tokenString,
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
		},
	}

	return c.JSON(http.StatusOK, response)
}
