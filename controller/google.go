package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fleimkeipa/kubernetes-api/config"
	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/uc"
	"github.com/fleimkeipa/kubernetes-api/util"

	"github.com/labstack/echo/v4"
)

type GoogleAuthHandler struct {
	userUC *uc.UserUC
}

func NewGoogleAuthHandler(userUC *uc.UserUC) *GoogleAuthHandler {
	return &GoogleAuthHandler{
		userUC: userUC,
	}
}

// Google Login godoc
//
//	@Summary		Redirect to Google login page
//	@Description	This endpoint initiates the Google OAuth2 login process by redirecting the user to Googles login page.
//	@Tags			oAuth
//	@Success		303	{string}	string	"Redirects to Google login page"
//	@Router			/auth/google_login [get]
func (rc *GoogleAuthHandler) GoogleLogin(c echo.Context) error {
	var url = config.AppConfig.GoogleLoginConfig.AuthCodeURL("randomstate")

	c.Redirect(303, url)
	return c.JSON(303, url)
}

// Google Callback godoc
//
//	@Summary		Google OAuth2 callback
//	@Description	This endpoint handles the callback from Google after a user authorizes the app. It exchanges the authorization code for an access token and retrieves the users profile information.
//	@Tags			oAuth
//	@Param			state	query		string	true	"State for CSRF protection"
//	@Param			code	query		string	true	"Authorization code returned by Google"
//	@Success		200		{string}	string	"User's Google profile data"
//	@Failure		400		{string}	string	"Error message"
//	@Router			/auth/google_callback [get]
func (rc *GoogleAuthHandler) GoogleCallback(c echo.Context) error {
	var state = c.QueryParam("state")
	if state != "randomstate" {
		return c.String(http.StatusBadRequest, "States don't Match!!")
	}

	var code = c.QueryParam("code")

	var googlecon = config.GoogleConfig()

	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Errorf("Code-Token Exchange Failed : %w", err).Error()})
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Errorf("user Data Fetch Failed %w", err).Error()})
	}

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Errorf("JSON Parsing Failed %w", err).Error()})
	}

	var googleUser = new(model.GoogleUser)
	if err := json.Unmarshal(userData, googleUser); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	user, err := rc.userUC.GetByUsernameOrEmail(c.Request().Context(), googleUser.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	jwt, err := util.GenerateJWT(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": jwt, "username": user.Username, "message": "Successfully logged in"})
}
