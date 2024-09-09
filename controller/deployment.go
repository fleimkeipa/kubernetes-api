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

func NewDeploymentHandler(podsUC *uc.DeploymentUC) *DeploymentHandler {
	return &DeploymentHandler{
		deploymentUC: podsUC,
	}
}

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

func (rc *DeploymentHandler) List(c echo.Context) error {
	var namespace = c.QueryParam("namespace")

	var opts = metav1.ListOptions{}

	list, err := rc.deploymentUC.Get(c.Request().Context(), namespace, opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"data": list})
}

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
