package util

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// check for valid admin token
func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := ValidateJWT(c); err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Authentication required"})
		}

		if err := ValidateAdminRoleJWT(c); err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Only Administrator is allowed to perform this action"})
		}

		return next(c)
	}
}

// check for valid viewer token
func JWTAuthViewer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := ValidateJWT(c); err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Authentication required"})
		}

		if err := ValidateViewerRoleJWT(c); err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Only registered Viewers are allowed to perform this action"})
		}

		return next(c)
	}
}
