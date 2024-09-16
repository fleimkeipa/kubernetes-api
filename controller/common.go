package controller

import (
	"strconv"

	"github.com/fleimkeipa/kubernetes-api/model"

	"github.com/labstack/echo/v4"
)

func getPagination(c echo.Context) model.PaginationOpts {
	var limitQuery = c.QueryParam("limit")
	var skipQuery = c.QueryParam("skip")

	limit, _ := strconv.Atoi(limitQuery)

	skip, _ := strconv.Atoi(skipQuery)

	return model.PaginationOpts{
		Skip:  uint(skip),
		Limit: uint(limit),
	}
}

func getFilter(c echo.Context, query string) model.Filter {
	var param = c.QueryParam(query)
	if param == "" {
		return model.Filter{}
	}

	return model.Filter{
		IsSended: true,
		Value:    param,
	}
}
