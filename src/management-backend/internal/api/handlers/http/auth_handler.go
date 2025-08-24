package http

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/arcaptcha/arcaptcha-go"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/config"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/admin/port"
	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/internal/api/dto"
)

var ErrCaptchaFailed = errors.New("captcha failed")

type AuthHandler struct {
	service   port.Service
	cfg       config.JWTConfig
	arcaptcha config.ArcaptchaConfig
}

func NewAuthHandler(s port.Service, c config.JWTConfig, a config.ArcaptchaConfig) *AuthHandler {
	return &AuthHandler{service: s, cfg: c, arcaptcha: a}
}

// @Summary      User login with captcha
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        loginRequest  body  dto.LoginRequest true "Login credentials"
// @Success      200  {object}  dto.LoginResponse
// @Failure      default  {object}  dto.Error
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Invalid request"})
	}

	if err := h.arcaptchaVerify(req.CaptchaToken); err != nil {
		code := http.StatusUnauthorized
		return c.JSON(code, &dto.Error{Code: code, Message: ErrCaptchaFailed.Error()})
	}

	user, err := h.service.Authenticate(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Invalid credentials"})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = user.ID
	claims["email"] = user.Email
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(h.cfg.Expiration)).Unix()

	tokenString, err := token.SignedString([]byte(h.cfg.Secret))
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

func (h *AuthHandler) arcaptchaVerify(token string) error {
	website := arcaptcha.NewWebsite(h.arcaptcha.SiteKey, h.arcaptcha.SecretKey)
	res, err := website.Verify(token)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCaptchaFailed, err)
	}
	if !res.Success {
		return ErrCaptchaFailed
	}
	return nil
}
