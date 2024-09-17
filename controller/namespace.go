package controller

import (
	"fmt"
	"net/http"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/uc"

	"github.com/labstack/echo/v4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceHandler struct {
	namespaceUC *uc.NamespaceUC
}

func NewNamespaceHandler(namespaceUC *uc.NamespaceUC) *NamespaceHandler {
	return &NamespaceHandler{
		namespaceUC: namespaceUC,
	}
}

// List godoc
//
//	@Summary		List namespaces
//	@Description	Retrieves a list of namespaces from the Kubernetes cluster.
//	@Tags			namespaces
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	SuccessResponse	"List of namespaces"
//	@Failure		500				{object}	FailureResponse	"Bad request or error message"
//	@Router			/namespaces [get]
func (rc *NamespaceHandler) List(c echo.Context) error {
	var opts = metav1.ListOptions{}
	list, err := rc.namespaceUC.List(c.Request().Context(), opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to retrieve namespaces: %v", err),
			Message: "There was an error retrieving the list of namespaces. Please try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    list,
		Message: "Namespaces retrieved successfully.",
	})
}

// Create godoc
//
//	@Summary		Create a new namespace
//	@Description	Creates a new namespace in the Kubernetes cluster.
//	@Tags			namespaces
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			namespace		body		model.NamespaceRequest	true	"Namespace request body"
//	@Success		201				{object}	SuccessResponse			"Successfully created namespace"
//	@Failure		400				{object}	FailureResponse			"Bad request or error message"
//	@Failure		500				{object}	FailureResponse			"Interval error"
//	@Router			/namespaces [post]
func (rc *NamespaceHandler) Create(c echo.Context) error {
	var input model.NamespaceRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to bind namespace request: %v", err),
			Message: "Invalid input. Please verify the data and try again.",
		})
	}

	namespace, err := rc.namespaceUC.Create(c.Request().Context(), &input.Namespace, input.Opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to create namespace: %v", err),
			Message: "There was an error creating the namespace. Please try again.",
		})
	}

	return c.JSON(http.StatusCreated, SuccessResponse{
		Data:    namespace.Name,
		Message: "Namespace created successfully.",
	})
}

// GetByNameOrUID godoc
//
//	@Summary		Get a namespace by name or UID
//	@Description	Retrieves a namespace from the Kubernetes cluster by its name or UID.
//	@Tags			namespaces
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			id				path		string			true	"Name or UID of the namespace"
//	@Success		200				{object}	SuccessResponse	"Details of the requested namespace"
//	@Failure		500				{object}	FailureResponse	"Interval error"
//	@Router			/namespaces/{id} [get]
func (rc *NamespaceHandler) GetByNameOrUID(c echo.Context) error {
	var nameOrUID = c.Param("id")

	var opts = metav1.ListOptions{}

	list, err := rc.namespaceUC.GetByNameOrUID(c.Request().Context(), nameOrUID, opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to retrieve namespace: %v", err),
			Message: "Error retrieving namespace. Please check the name or UID and try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    list,
		Message: "Namespace retrieved successfully.",
	})
}

// Delete godoc
//
//	@Summary		Delete a namespace by name or UID
//	@Description	Deletes a namespace from the Kubernetes cluster by its name or UID.
//	@Tags			namespaces
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			id				path		string			true	"Name or UID of the namespace"
//	@Success		200				{string}	SuccessResponse	"Success message"
//	@Failure		500				{object}	FailureResponse	"Interval error"
//	@Router			/namespaces/{id} [delete]
func (rc *NamespaceHandler) Delete(c echo.Context) error {
	var nameOrUID = c.Param("id")

	var opts = metav1.DeleteOptions{}

	if err := rc.namespaceUC.Delete(c.Request().Context(), nameOrUID, opts); err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to delete namespace: %v", err),
			Message: "Error deleting namespace. Please check the name or UID and try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Namespace deleted successfully.",
	})
}

// Update godoc
//
//	@Summary		Update an existing namespace
//	@Description	Updates an existing namespace in the Kubernetes cluster.
//	@Tags			namespaces
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			namespace		body		model.NamespaceRequest	true	"Namespace request body"
//	@Success		200				{object}	SuccessResponse			"Successfully updated the namespace"
//	@Failure		400				{object}	FailureResponse			"Bad request or invalid data"
//	@Failure		500				{object}	FailureResponse			"Interval error"
//	@Router			/namespaces [put]
func (rc *NamespaceHandler) Update(c echo.Context) error {
	var request model.NamespaceRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to bind namespace request: %v", err),
			Message: "Invalid input. Please verify the data and try again.",
		})
	}

	namespace, err := rc.namespaceUC.Update(c.Request().Context(), &request.Namespace, metav1.UpdateOptions(request.Opts))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to update namespace: %v", err),
			Message: "Error updating namespace. Please try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    namespace.Name,
		Message: "Namespace updated successfully.",
	})
}
