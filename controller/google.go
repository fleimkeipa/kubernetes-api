package controller

import (
	"context"
	"io"
	"net/http"

	"github.com/fleimkeipa/kubernetes-api/config"

	"github.com/labstack/echo/v4"
)

// Google Login godoc
//
// @Summary Redirect to Google login page
// @Description This endpoint initiates the Google OAuth2 login process by redirecting the user to Google's login page.
// @Tags OAuth
// @Success 303 {string} string "Redirects to Google login page"
// @Router /auth/google/login [get]
func GoogleLogin(c echo.Context) error {
	var url = config.AppConfig.GoogleLoginConfig.AuthCodeURL("randomstate")

	c.Redirect(303, url)
	return c.JSON(303, url)
}

// Google Callback godoc
//
// @Summary Google OAuth2 callback
// @Description This endpoint handles the callback from Google after a user authorizes the app. It exchanges the authorization code for an access token and retrieves the user's profile information.
// @Tags OAuth
// @Param state query string true "State for CSRF protection"
// @Param code query string true "Authorization code returned by Google"
// @Success 200 {string} string "User's Google profile data"
// @Failure 400 {string} string "Error message"
// @Router /auth/google/callback [get]
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
