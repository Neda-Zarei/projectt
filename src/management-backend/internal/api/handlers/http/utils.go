package http

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func parseUintParam(c echo.Context, paramName string) (uint, error) {
	param := c.Param(paramName)
	id, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func parseQueryParamInt(c echo.Context, paramName string, defaultValue int) int {
	param := c.QueryParam(paramName)
	if param == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(param)
	if err != nil {
		return defaultValue
	}
	return value
}
