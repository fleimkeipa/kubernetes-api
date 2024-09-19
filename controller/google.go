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
//	@Success		303	{object}	map[string]string	"Redirects to Google login page"
//	@Failure		400	{object}	FailureResponse		"Error message"
//	@Router			/auth/google_login [get]
func (rc *GoogleAuthHandler) GoogleLogin(c echo.Context) error {
	// Load Google OAuth2 configuration
	config.GoogleConfig()

	// Generate the Google OAuth2 login URL
	var url = config.AppConfig.GoogleLoginConfig.AuthCodeURL("randomstate")

	// Attempt to redirect the user to the Google login page
	if err := c.Redirect(http.StatusSeeOther, url); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to redirect to Google login page: %v", err),
			Message: "There was an issue starting the Google login process. Please try again later.",
		})
	}

	// Return the URL in case the redirect works but the client needs the login URL
	return c.JSON(http.StatusSeeOther, echo.Map{
		"url":     url,
		"message": "Redirecting to Google for login...",
	})
}

// Google Callback godoc
//
//	@Summary		Google OAuth2 callback
//	@Description	This endpoint handles the callback from Google after a user authorizes the app. It exchanges the authorization code for an access token and retrieves the users profile information.
//	@Tags			oAuth
//	@Param			state	query		string			true	"State for CSRF protection"
//	@Param			code	query		string			true	"Authorization code returned by Google"
//	@Success		200		{object}	AuthResponse	"User's Google profile data"
//	@Failure		400		{object}	FailureResponse	"Error message"
//	@Failure		500		{object}	FailureResponse	"Interval error"
//	@Router			/auth/google_callback [get]
func (rc *GoogleAuthHandler) GoogleCallback(c echo.Context) error {
	var state = c.QueryParam("state")
	if state != "randomstate" {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Message: "State parameter mismatch! Please restart the login process.",
		})
	}

	var code = c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Message: "Authorization code missing! Unable to proceed with login.",
		})
	}

	var googlecon = config.GoogleConfig()

	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to exchange authorization code for access token: %v", err),
			Message: "There was an issue communicating with Google. Please try again.",
		})
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to fetch user data from Google: %v", err),
			Message: "Unable to retrieve your profile information from Google.",
		})
	}

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to read response from Google: %v", err),
			Message: "There was an issue processing the response from Google.",
		})
	}

	var googleUser = new(model.GoogleUser)
	if err := json.Unmarshal(userData, googleUser); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to parse Google user data: %v", err),
			Message: "There was an issue parsing your profile information.",
		})
	}

	user, err := rc.userUC.GetByUsernameOrEmail(c.Request().Context(), googleUser.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("User retrieval failed: %v", err),
			Message: "We could not find a user associated with your Google account.",
		})
	}

	jwt, err := util.GenerateJWT(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("JWT generation failed: %v", err),
			Message: "There was an issue generating your authentication token.",
		})
	}

	return c.JSON(http.StatusOK, AuthResponse{
		Token:    jwt,
		Type:     "oauth2",
		Username: user.Username,
		Message:  "Successfully logged in with Google.",
	})
}
