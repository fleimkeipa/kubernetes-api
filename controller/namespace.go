package controller

import (
	"kub/model"
	"kub/uc"
	"net/http"

	"github.com/labstack/echo/v4"
	corev1 "k8s.io/api/core/v1"
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
	var namespace = corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       input.TypeMeta.Kind,
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      input.ObjectMeta.Name,
			Namespace: input.ObjectMeta.NameSpace,
		},
	}

	_, err := rc.namespaceUC.Create(c.Request().Context(), &namespace, opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"pod": namespace.Name})
}
