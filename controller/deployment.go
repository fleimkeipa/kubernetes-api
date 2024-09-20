package controller

import (
	"fmt"
	"net/http"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/uc"

	"github.com/labstack/echo/v4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentHandler struct {
	deploymentUC *uc.DeploymentUC
}

func NewDeploymentHandler(deploymentUC *uc.DeploymentUC) *DeploymentHandler {
	return &DeploymentHandler{
		deploymentUC: deploymentUC,
	}
}

// Create godoc
//
//	@Summary		Create a new deployment
//	@Description	Creates a new deployment in the Kubernetes cluster.
//	@Tags			deployments
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			deployment		body		model.DeploymentCreateRequest	true	"Deployment request body"
//	@Success		201				{object}	map[string]string				"Suxccessfully created deployment"
//	@Failure		400				{object}	FailureResponse					"Bad request or error message"
//	@Failure		500				{object}	FailureResponse					"Interval error"
//	@Router			/deployments [post]
func (rc *DeploymentHandler) Create(c echo.Context) error {
	var request model.DeploymentCreateRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to parse request body: %v", err),
			Message: "Invalid request format. Please ensure your data is correctly formatted.",
		})
	}

	deployment, err := rc.deploymentUC.Create(c.Request().Context(), &request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to create deployment: %v", err),
			Message: "There was an error creating the deployment. Please check your data and try again.",
		})
	}

	return c.JSON(http.StatusCreated, SuccessResponse{
		Data:    deployment.Name,
		Message: "Deployment created successfully.",
	})
}

// Update godoc
//
//	@Summary		Update an existing deployment
//	@Description	Updates an existing deployment in the Kubernetes cluster.
//	@Tags			deployments
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			deployment		body		model.DeploymentUpdateRequest	true	"Deployment request body"
//	@Success		200				{object}	SuccessResponse					"Successfully updated the deployment"
//	@Failure		400				{object}	FailureResponse					"Bad request or invalid data"
//	@Failure		500				{object}	FailureResponse					"Interval error"
//	@Router			/deployments [put]
func (rc *DeploymentHandler) Update(c echo.Context) error {
	var id = c.Param("id")

	var request model.DeploymentUpdateRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to parse request body: %v", err),
			Message: "Invalid request format. Please ensure your data is correctly formatted.",
		})
	}

	deployment, err := rc.deploymentUC.Update(c.Request().Context(), "", id, &request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to update deployment: %v", err),
			Message: "There was an error updating the deployment. Please check your data and try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    deployment.Name,
		Message: "Deployment updated successfully.",
	})
}

// List godoc
//
//	@Summary		List deployments
//	@Description	Retrieves a list of deployments from the Kubernetes cluster, optionally filtered by namespace.
//	@Tags			deployments
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			namespace		query		string			false	"Namespace to filter deployments by"
//	@Success		200				{object}	SuccessResponse	"List of deployments"
//	@Failure		500				{object}	FailureResponse	"Interval error"
//	@Router			/deployments [get]
func (rc *DeploymentHandler) List(c echo.Context) error {
	var namespace = c.QueryParam("namespace")

	var opts = metav1.ListOptions{}

	list, err := rc.deploymentUC.List(c.Request().Context(), namespace, opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to list deployments: %v", err),
			Message: "There was an issue retrieving deployments. Please try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    list,
		Message: "Deployments retrieved successfully.",
	})
}

// GetByNameOrUID godoc
//
//	@Summary		Get a deployment by name or UID
//	@Description	Retrieves a deployment from the Kubernetes cluster by its name or UID, optionally filtered by namespace.
//	@Tags			deployments
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			namespace		query		string			false	"Namespace to filter the deployment by"
//	@Param			id				path		string			true	"Name or UID of the deployment"
//	@Success		200				{object}	SuccessResponse	"Details of the requested deployment"
//	@Failure		500				{object}	FailureResponse	"Interval error"
//	@Router			/deployments/{id} [get]
func (rc *DeploymentHandler) GetByNameOrUID(c echo.Context) error {
	var namespace = c.QueryParam("namespace")
	var nameOrUID = c.Param("id")

	var opts = metav1.ListOptions{}

	list, err := rc.deploymentUC.GetByNameOrUID(c.Request().Context(), namespace, nameOrUID, opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to retrieve deployment: %v", err),
			Message: "Could not find the requested deployment. Please verify the name or UID and try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    list,
		Message: "Deployment retrieved successfully.",
	})
}

// Delete godoc
//
//	@Summary		Delete a deployment by name or UID
//	@Description	Deletes a deployment from the Kubernetes cluster by its name or UID, optionally filtered by namespace.
//	@Tags			deployments
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			namespace		query		string			false	"Namespace to filter the deployment by"
//	@Param			id				path		string			true	"Name or UID of the deployment"
//	@Success		200				{string}	SuccessResponse	"Success message"
//	@Failure		500				{object}	FailureResponse	"Bad request or error message"
//	@Router			/deployments/{id} [delete]
func (rc *DeploymentHandler) Delete(c echo.Context) error {
	var namespace = c.QueryParam("namespace")
	var nameOrUID = c.Param("id")

	var opts = metav1.DeleteOptions{}

	if err := rc.deploymentUC.Delete(c.Request().Context(), namespace, nameOrUID, opts); err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to delete deployment: %v", err),
			Message: "There was an error deleting the deployment. Please check the name or UID and try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Deployment deleted successfully.",
	})
}
