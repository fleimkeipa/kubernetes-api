package controller

import (
	"net/http"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/uc"

	"github.com/labstack/echo/v4"
	v1 "k8s.io/api/apps/v1"
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
	var input model.PodsRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	var opts = metav1.CreateOptions{
		TypeMeta: metav1.TypeMeta{
			Kind:       "",
			APIVersion: "",
		},
		DryRun:          []string{},
		FieldManager:    "",
		FieldValidation: "",
	}
	var deployment = v1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       input.TypeMeta.Kind,
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      input.ObjectMeta.Name,
			Namespace: input.ObjectMeta.NameSpace,
		},
	}

	_, err := rc.deploymentUC.Create(c.Request().Context(), &deployment, opts)
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
