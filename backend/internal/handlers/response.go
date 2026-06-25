package handlers

import "github.com/labstack/echo/v4"

func readJSON(c echo.Context, target any) error {
	return c.Bind(target)
}

func writeJSON(c echo.Context, status int, payload any) error {
	return c.JSON(status, payload)
}

func writeError(c echo.Context, status int, message string) error {
	return c.JSON(status, map[string]string{"error": message})
}
