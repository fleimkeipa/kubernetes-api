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

// CreateUser godoc
//
//	@Summary		CreateUser creates a new user
//	@Description	This endpoint creates a new user by providing username, email, password, and role ID.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		model.UserRequest	true	"User creation input"
//	@Success		201				{object}	SuccessResponse		"user username"
//	@Failure		400				{object}	FailureResponse		"Error message including details on failure"
//	@Failure		500				{object}	FailureResponse		"Interval error"
//	@Router			/users [post]
func (rc *UserHandlers) CreateUser(c echo.Context) error {
	var input model.UserRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to bind request: %v", err),
			Message: "Invalid request format. Please check the input data and try again.",
		})
	}

	var user = model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		RoleID:   input.RoleID,
	}

	_, err := rc.userUC.Create(c.Request().Context(), user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to create user: %v", err),
			Message: "User creation failed. Please check the provided details and try again.",
		})
	}

	return c.JSON(http.StatusCreated, SuccessResponse{
		Data:    user.Username,
		Message: "User created successfully.",
	})
}

// UpdateUser godoc
//
//	@Summary		UpdateUser updates an existing user
//	@Description	This endpoint updates a user by providing username, email, password, and role ID.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		model.UserRequest	true	"User update input"
//	@Success		200				{object}	SuccessResponse		"user username"
//	@Failure		400				{object}	FailureResponse		"Error message including details on failure"
//	@Failure		500				{object}	FailureResponse		"Interval error"
//	@Router			/users/{id} [put]
func (rc *UserHandlers) UpdateUser(c echo.Context) error {
	var id = c.Param("id")
	var input model.UserRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to bind request: %v", err),
			Message: "Invalid request format. Please check the input data and try again.",
		})
	}

	var user = model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
		RoleID:   input.RoleID,
	}

	_, err := rc.userUC.Update(c.Request().Context(), id, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to update user: %v", err),
			Message: "User update failed. Please check the provided details and try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    user.Username,
		Message: "User updated successfully.",
	})
}

// DeleteUser godoc
//
//	@Summary		DeleteUser deletes an existing user
//	@Description	This endpoint deletes a user by providing user id.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	SuccessResponse	"user username"
//	@Failure		500				{object}	FailureResponse	"Interval error"
//	@Router			/users/{id} [delete]
func (rc *UserHandlers) DeleteUser(c echo.Context) error {
	var id = c.Param("id")

	if err := rc.userUC.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to delete user: %v", err),
			Message: "User delete failed. Please check the provided details and try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "User deleted successfully.",
	})
}

// Login godoc
//
//	@Summary		User login
//	@Description	This endpoint allows a user to log in by providing a valid username and password.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		model.Login		true	"User login input"
//	@Success		200		{object}	SuccessResponse	"Successfully logged in with JWT token"
//	@Failure		400		{object}	FailureResponse	"Error message including details on failure"
//	@Failure		500		{object}	FailureResponse	"Interval error"
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

		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to bind request: %v", errorMessage),
			Message: "Invalid login details. Please ensure all required fields are provided and try again.",
		})
	}

	user, err := rc.userUC.GetByUsernameOrEmail(c.Request().Context(), input.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("User not found: %v", err),
			Message: "User not found. Please check the username and try again.",
		})
	}

	if err := model.ValidateUserPassword(user.Password, input.Password); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Invalid password: %v", err),
			Message: "Invalid password. Please check the password and try again.",
		})
	}

	jwt, err := util.GenerateJWT(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to generate JWT: %v", err),
			Message: "Login failed. Please try again later.",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token":    jwt,
		"username": input.Username,
		"message":  "Successfully logged in",
	})
}

// List godoc
//
//	@Summary		List all users
//	@Description	Retrieves a filtered and paginated list of users from the database based on query parameters.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			limit			query		string			false	"Limit the number of users returned"
//	@Param			skip			query		string			false	"Number of users to skip for pagination"
//	@Param			username		query		string			false	"Filter users by username"
//	@Param			email			query		string			false	"Filter users by email"
//	@Param			role_id			query		string			false	"Filter users by role ID"
//	@Success		200				{object}	SuccessResponse	"Successful response containing the list of users"
//	@Failure		500				{object}	FailureResponse	"Interval error"
//	@Router			/users [get]
func (rc *UserHandlers) List(c echo.Context) error {
	var opts = rc.getUsersFindOpts(c, model.ZeroCreds)

	list, err := rc.userUC.List(c.Request().Context(), &opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to retrieve user list: %v", err),
			Message: "Unable to retrieve the list of users. Please check the query parameters and try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    list,
		Message: "Users retrieved successfully.",
	})
}

func (rc *UserHandlers) getUsersFindOpts(c echo.Context, fields ...string) model.UserFindOpts {
	return model.UserFindOpts{
		PaginationOpts: getPagination(c),
		FieldsOpts: model.FieldsOpts{
			Fields: fields,
		},
		Username: getFilter(c, "username"),
		Email:    getFilter(c, "email"),
		RoleID:   getFilter(c, "role_id"),
	}
}
