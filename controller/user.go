package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/uc"
	"github.com/fleimkeipa/kubernetes-api/util"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserHandlers struct {
	userUC *uc.UserUC
}

func NewUserHandlers(uc *uc.UserUC) *UserHandlers {
	return &UserHandlers{
		userUC: uc,
	}
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	This endpoint registers a new user by providing username, email, password, and role ID.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			body	body		model.Register		true	"User registration input"
//	@Success		201		{object}	map[string]string	"User created"
//	@Failure		400		{object}	map[string]string	"Error message"
//	@Router			/auth/register [post]
func (rc *UserHandlers) Register(c echo.Context) error {
	var input model.Register

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	var user = model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		RoleID:   input.RoleID,
	}

	_, err := rc.userUC.Create(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"user": user.Username})
}

// Login godoc
//
//	@Summary		User login
//	@Description	This endpoint allows a user to log in by providing a valid username and password.
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			body	body		model.Login				true	"User login input"
//	@Success		200		{object}	map[string]interface{}	"Successfully logged in with JWT token"
//	@Failure		400		{object}	map[string]string		"Error message"
//	@Router			/auth/login [post]
func (rc *UserHandlers) Login(c echo.Context) error {
	var input model.Login

	if err := c.Bind(&input); err != nil {
		var errorMessage string
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			validationError := validationErrors[0]
			if validationError.Tag() == "required" {
				errorMessage = fmt.Sprintf("%s not provided", validationError.Field())
			}
		}

		return c.JSON(http.StatusBadRequest, echo.Map{"error": errorMessage})
	}

	user, err := rc.userUC.GetUserByUsername(c.Request().Context(), input.Username)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := model.ValidateUserPassword(user.Password, input.Password); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	jwt, err := util.GenerateJWT(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": jwt, "username": input.Username, "message": "Successfully logged in"})
}
