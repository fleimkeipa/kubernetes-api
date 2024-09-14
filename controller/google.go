package controller

import (
	"context"
	"io"
	"net/http"

	"github.com/fleimkeipa/kubernetes-api/config"

	"github.com/labstack/echo/v4"
)

func GoogleLogin(c echo.Context) error {
	var url = config.AppConfig.GoogleLoginConfig.AuthCodeURL("randomstate")

	c.Redirect(303, url)
	return c.JSON(303, url)
}

func GoogleCallback(c echo.Context) error {
	var state = c.QueryParam("state")
	if state != "randomstate" {
		return c.String(400, "States don't Match!!")
	}

	var code = c.QueryParam("code")

	var googlecon = config.GoogleConfig()

	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		return c.String(400, "Code-Token Exchange Failed")
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.String(400, "User Data Fetch Failed")
	}

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.String(400, "JSON Parsing Failed")
	}

	return c.String(200, string(userData))
}
