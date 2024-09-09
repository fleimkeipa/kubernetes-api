package controller

import (
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

func (rc *NamespaceHandler) Get(c echo.Context) error {
	var opts = metav1.ListOptions{}
	list, err := rc.namespaceUC.Get(c.Request().Context(), opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"data": list})
}

func (rc *NamespaceHandler) Create(c echo.Context) error {
	var input model.NamespaceRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	namespace, err := rc.namespaceUC.Create(c.Request().Context(), &input.Namespace, input.Opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"pod": namespace.Name})
}

func (rc *NamespaceHandler) GetByNameOrUID(c echo.Context) error {
	var nameOrUID = c.Param("id")

	var opts = metav1.ListOptions{}

	list, err := rc.namespaceUC.GetByNameOrUID(c.Request().Context(), nameOrUID, opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"data": list})
}

func (rc *NamespaceHandler) Delete(c echo.Context) error {
	var nameOrUID = c.Param("id")

	var opts = metav1.DeleteOptions{}

	if err := rc.namespaceUC.Delete(c.Request().Context(), nameOrUID, opts); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, "deleted succesfully")
}
