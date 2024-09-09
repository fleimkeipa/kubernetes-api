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

// Create godoc
//
//	@Summary		Create a new pod
//	@Description	Creates a new pod in the Kubernetes cluster.
//	@Tags			pods
//	@Accept			json
//	@Produce		json
//	@Param			pod	body		model.PodsRequest	true	"Pod request body"
//	@Success		201	{object}	map[string]string	"Successfully created the pod"
//	@Failure		400	{object}	map[string]string	"Bad request or invalid data"
//	@Router			/pods [post]
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

// Update godoc
//
//	@Summary		Update an existing pod
//	@Description	Updates an existing pod in the Kubernetes cluster.
//	@Tags			pods
//	@Accept			json
//	@Produce		json
//	@Param			pod	body		model.PodsRequest	true	"Pod request body"
//	@Success		200	{object}	map[string]string	"Successfully updated the pod"
//	@Failure		400	{object}	map[string]string	"Bad request or invalid data"
//	@Router			/pods [put]
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

// List godoc
//
//	@Summary		List pods
//	@Description	Retrieves a list of pods from the Kubernetes cluster, optionally filtered by namespace.
//	@Tags			pods
//	@Accept			json
//	@Produce		json
//	@Param			namespace	query		string					false	"Namespace to filter pods by"
//	@Success		200			{object}	map[string]interface{}	"List of pods"
//	@Failure		400			{object}	map[string]string		"Bad request or invalid data"
//	@Router			/pods [get]
func (rc *PodsHandler) List(c echo.Context) error {
	var namespace = c.QueryParam("namespace")

	var opts = metav1.ListOptions{}

	list, err := rc.podsUC.Get(c.Request().Context(), namespace, opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"data": list})
}

// GetByNameOrUID godoc
//
//	@Summary		Get a pod by name or UID
//	@Description	Retrieves a pod from the Kubernetes cluster by its name or UID, optionally filtered by namespace.
//	@Tags			pods
//	@Accept			json
//	@Produce		json
//	@Param			namespace	query		string					false	"Namespace to filter the pod by"
//	@Param			id			path		string					true	"Name or UID of the pod"
//	@Success		200			{object}	map[string]interface{}	"Details of the requested pod"
//	@Failure		400			{object}	map[string]string		"Bad request or invalid data"
//	@Router			/pods/{id} [get]
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

// Delete godoc
//
//	@Summary		Delete a pod by name or UID
//	@Description	Deletes a pod from the Kubernetes cluster by its name or UID, optionally filtered by namespace.
//	@Tags			pods
//	@Accept			json
//	@Produce		json
//	@Param			namespace	query		string				false	"Namespace to filter the pod by"
//	@Param			id			path		string				true	"Name or UID of the pod"
//	@Success		200			{string}	string				"Success message"
//	@Failure		400			{object}	map[string]string	"Bad request or error message"
//	@Router			/pods/{id} [delete]
func (rc *PodsHandler) Delete(c echo.Context) error {
	var namespace = c.QueryParam("namespace")
	var nameOrUID = c.Param("id")

	var opts = metav1.DeleteOptions{}

	if err := rc.podsUC.Delete(c.Request().Context(), namespace, nameOrUID, opts); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, "deleted succesfully")
}
