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

type GithubAuthHandler struct {
	userUC *uc.UserUC
}

func NewGithubAuthHandler(userUC *uc.UserUC) *GithubAuthHandler {
	return &GithubAuthHandler{
		userUC: userUC,
	}
}

// Github Login godoc
//
//	@Summary		Redirect to Github login page
//	@Description	This endpoint initiates the Github OAuth2 login process by redirecting the user to Githubs login page.
//	@Tags			oAuth
//	@Success		303	{object}	map[string]string	"Redirects to Github login page"
//	@Failure		400	{object}	FailureResponse		"Error message"
//	@Router			/auth/github_login [get]
func (rc *GithubAuthHandler) GithubLogin(c echo.Context) error {
	// Load Github OAuth2 configuration
	config.GithubConfig()

	// Generate the Github OAuth2 login URL
	var url = config.AppConfig.GitHubLoginConfig.AuthCodeURL("randomstate")

	// Attempt to redirect the user to the Github login page
	if err := c.Redirect(http.StatusSeeOther, url); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to redirect to Github login page: %v", err),
			Message: "There was an issue starting the Github login process. Please try again later.",
		})
	}

	// Return the URL in case the redirect works but the client needs the login URL
	return c.JSON(http.StatusSeeOther, echo.Map{
		"url":     url,
		"message": "Redirecting to Github for login...",
	})
}

// Github Callback godoc
//
//	@Summary		Github OAuth2 callback
//	@Description	This endpoint handles the callback from Github after a user authorizes the app. It exchanges the authorization code for an access token and retrieves the users profile information.
//	@Tags			oAuth
//	@Param			state	query		string			true	"State for CSRF protection"
//	@Param			code	query		string			true	"Authorization code returned by Github"
//	@Success		200		{object}	AuthResponse	"User's Github profile data"
//	@Failure		400		{object}	FailureResponse	"Error message"
//	@Failure		500		{object}	FailureResponse	"Interval error"
//	@Router			/auth/github_callback [get]
func (rc *GithubAuthHandler) GithubCallback(c echo.Context) error {
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

	var Githubcon = config.GithubConfig()

	token, err := Githubcon.Exchange(context.Background(), code)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to exchange authorization code for access token: %v", err),
			Message: "There was an issue communicating with Github. Please try again.",
		})
	}

	var client = &http.Client{}

	var url = "https://api.github.com/user"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to create request for Github: %v", err),
			Message: "Unable to create request for Github.",
		})
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to fetch user data from Github: %v", err),
			Message: "Unable to retrieve your profile information from Github.",
		})
	}

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to read response from Github: %v", err),
			Message: "There was an issue processing the response from Github.",
		})
	}

	var GithubUser = new(model.GithubUser)
	if err := json.Unmarshal(userData, GithubUser); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to parse Github user data: %v", err),
			Message: "There was an issue parsing your profile information.",
		})
	}

	user, err := rc.userUC.GetByUsernameOrEmail(c.Request().Context(), GithubUser.Email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("User retrieval failed: %v", err),
			Message: "We could not find a user associated with your Github account.",
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
		Message:  "Successfully logged in with Github.",
	})
}
