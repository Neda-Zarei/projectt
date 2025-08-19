package middleware

import (
	"net/http"
	"strings"

	"hamgit.ir/arcaptcha/arcaptcha-dumbledore/management-backend/app"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	app app.App
}

func NewAuthMiddleware(a app.App) *AuthMiddleware {
	return &AuthMiddleware{app: a}
}

func (m *AuthMiddleware) ValidateJWT() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Authorization header required"})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Bearer token required"})
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(m.app.Config().JWT.Secret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Invalid token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Invalid token claims"})
			}

			c.Set("userID", claims["userID"])
			c.Set("email", claims["email"])
			c.Set("role", claims["role"])

			return next(c)
		}
	}
}
