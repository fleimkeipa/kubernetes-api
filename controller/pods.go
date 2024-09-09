package controller

import (
	"net/http"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/uc"

	"github.com/labstack/echo/v4"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodsHandler struct {
	podsUC *uc.PodsUC
}

func NewPodsHandler(podsUC *uc.PodsUC) *PodsHandler {
	return &PodsHandler{
		podsUC: podsUC,
	}
}

func (rc *PodsHandler) Create(c echo.Context) error {
	var request model.PodsRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	pod, err := rc.podsUC.Create(c.Request().Context(), &request.Pod, request.Opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"pod": pod.Name})
}

func (rc *PodsHandler) Update(c echo.Context) error {
	var request model.PodsRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	pod, err := rc.podsUC.Update(c.Request().Context(), &request.Pod, metav1.UpdateOptions(request.Opts))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"pod": pod.Name})
}

func (rc *PodsHandler) List(c echo.Context) error {
	var namespace = c.QueryParam("namespace")

	var opts = metav1.ListOptions{}

	list, err := rc.podsUC.Get(c.Request().Context(), namespace, opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"data": list})
}

func (rc *PodsHandler) GetByNameOrUID(c echo.Context) error {
	var namespace = c.QueryParam("namespace")
	var nameOrUID = c.Param("id")

	var opts = metav1.ListOptions{}

	list, err := rc.podsUC.GetByNameOrUID(c.Request().Context(), namespace, nameOrUID, opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"data": list})
}

func (rc *PodsHandler) Delete(c echo.Context) error {
	var namespace = c.QueryParam("namespace")
	var nameOrUID = c.Param("id")

	var opts = metav1.DeleteOptions{}

	if err := rc.podsUC.Delete(c.Request().Context(), namespace, nameOrUID, opts); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, "deleted succesfully")
}
