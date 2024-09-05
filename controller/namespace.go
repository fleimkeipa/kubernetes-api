package controller

import (
	"kub/uc"
	"net/http"

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
