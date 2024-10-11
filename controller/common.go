package controller

import (
	"strconv"

	"github.com/fleimkeipa/kubernetes-api/model"

	"github.com/labstack/echo/v4"
)

type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type FailureResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type AuthResponse struct {
	Type     string `json:"type" example:"basic,oauth2"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func getPagination(c echo.Context) model.PaginationOpts {
	limitQuery := c.QueryParam("limit")
	skipQuery := c.QueryParam("skip")

	limit, _ := strconv.Atoi(limitQuery)

	if limit == 0 {
		limit = 30
	}

	skip, _ := strconv.Atoi(skipQuery)

	return model.PaginationOpts{
		Skip:  skip,
		Limit: limit,
	}
}

func getFilter(c echo.Context, query string) model.Filter {
	param := c.QueryParam(query)
	if param == "" {
		return model.Filter{}
	}

	return model.Filter{
		IsSended: true,
		Value:    param,
	}
}

func getKubeListOpts(c echo.Context) model.ListOptions {
	limitQuery := c.QueryParam("limit")

	continueQuery := c.QueryParam("continue")

	limit, _ := strconv.Atoi(limitQuery)

	return model.ListOptions{
		Continue: continueQuery,
		Limit:    int64(limit),
	}
}
