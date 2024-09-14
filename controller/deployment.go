package controller

import (
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
//	@Param			deployment	body		model.DeploymentRequest	true	"Deployment request body"
//	@Success		201			{object}	map[string]string		"Successfully created deployment"
//	@Failure		400			{object}	map[string]string		"Bad request or error message"
//	@Router			/deployments [post]
func (rc *DeploymentHandler) Create(c echo.Context) error {
	var request model.DeploymentRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	deployment, err := rc.deploymentUC.Create(c.Request().Context(), &request.Deployment, request.Opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"deployment": deployment.Name})
}

// List godoc
//
//	@Summary		List deployments
//	@Description	Retrieves a list of deployments from the Kubernetes cluster, optionally filtered by namespace.
//	@Tags			deployments
//	@Accept			json
//	@Produce		json
//	@Param			namespace	query		string					false	"Namespace to filter deployments by"
//	@Success		200			{object}	map[string]interface{}	"List of deployments"
//	@Failure		400			{object}	map[string]string		"Bad request or error message"
//	@Router			/deployments [get]
func (rc *DeploymentHandler) List(c echo.Context) error {
	var namespace = c.QueryParam("namespace")

	var opts = metav1.ListOptions{}

	list, err := rc.deploymentUC.Get(c.Request().Context(), namespace, opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"data": list})
}

// GetByNameOrUID godoc
//
//	@Summary		Get a deployment by name or UID
//	@Description	Retrieves a deployment from the Kubernetes cluster by its name or UID, optionally filtered by namespace.
//	@Tags			deployments
//	@Accept			json
//	@Produce		json
//	@Param			namespace	query		string					false	"Namespace to filter the deployment by"
//	@Param			id			path		string					true	"Name or UID of the deployment"
//	@Success		200			{object}	map[string]interface{}	"Details of the requested deployment"
//	@Failure		400			{object}	map[string]string		"Bad request or error message"
//	@Router			/deployments/{id} [get]
func (rc *DeploymentHandler) GetByNameOrUID(c echo.Context) error {
	var namespace = c.QueryParam("namespace")
	var nameOrUID = c.Param("id")

	var opts = metav1.ListOptions{}

	list, err := rc.deploymentUC.GetByNameOrUID(c.Request().Context(), namespace, nameOrUID, opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"data": list})
}

// Delete godoc
//
//	@Summary		Delete a deployment by name or UID
//	@Description	Deletes a deployment from the Kubernetes cluster by its name or UID, optionally filtered by namespace.
//	@Tags			deployments
//	@Accept			json
//	@Produce		json
//	@Param			namespace	query		string				false	"Namespace to filter the deployment by"
//	@Param			id			path		string				true	"Name or UID of the deployment"
//	@Success		200			{string}	string				"Success message"
//	@Failure		400			{object}	map[string]string	"Bad request or error message"
//	@Router			/deployments/{id} [delete]
func (rc *DeploymentHandler) Delete(c echo.Context) error {
	var namespace = c.QueryParam("namespace")
	var nameOrUID = c.Param("id")

	var opts = metav1.DeleteOptions{}

	if err := rc.deploymentUC.Delete(c.Request().Context(), namespace, nameOrUID, opts); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, "deleted succesfully")
}

// Update godoc
//
//	@Summary		Update an existing deployment
//	@Description	Updates an existing deployment in the Kubernetes cluster.
//	@Tags			deployments
//	@Accept			json
//	@Produce		json
//	@Param			deployment	body		model.DeploymentRequest	true	"Deployment request body"
//	@Success		200			{object}	map[string]string		"Successfully updated the deployment"
//	@Failure		400			{object}	map[string]string		"Bad request or invalid data"
//	@Router			/deployments [put]
func (rc *DeploymentHandler) Update(c echo.Context) error {
	var request model.DeploymentRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	deployment, err := rc.deploymentUC.Update(c.Request().Context(), &request.Deployment, metav1.UpdateOptions(request.Opts))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"deployment": deployment.Name})
}
