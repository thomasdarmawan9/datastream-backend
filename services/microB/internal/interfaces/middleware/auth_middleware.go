package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/thomasdarmawan9/datastream-backend/services/microB/internal/infrastructure/auth"
)

func JWTAuth(jwtManager *auth.JWTManager, allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing token"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid auth header"})
			}

			claims, err := jwtManager.Verify(parts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}

			// check role
			if len(allowedRoles) > 0 {
				allowed := false
				for _, r := range allowedRoles {
					if claims.Role == r {
						allowed = true
						break
					}
				}
				if !allowed {
					return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
				}
			}

			// simpan user di context
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}
