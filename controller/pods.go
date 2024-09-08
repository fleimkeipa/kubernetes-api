package controller

import (
	"kub/model"
	"kub/uc"
	"net/http"

	"github.com/labstack/echo/v4"
	corev1 "k8s.io/api/core/v1"
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
	var input model.PodsRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	var containers = []corev1.Container{}
	for _, v := range input.Spec.Containers {
		containers = append(containers, corev1.Container{
			Name:  v.Name,
			Image: v.Image,
		})
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
	var pod = corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       input.TypeMeta.Kind,
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      input.ObjectMeta.Name,
			Namespace: input.ObjectMeta.NameSpace,
		},
		Spec: corev1.PodSpec{
			Containers: containers,
		},
	}

	_, err := rc.podsUC.Create(c.Request().Context(), &pod, opts)
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
