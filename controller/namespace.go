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

// Get godoc
//
//	@Summary		List namespaces
//	@Description	Retrieves a list of namespaces from the Kubernetes cluster.
//	@Tags			namespaces
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}	"List of namespaces"
//	@Failure		400	{object}	map[string]string		"Bad request or error message"
//	@Router			/namespaces [get]
func (rc *NamespaceHandler) Get(c echo.Context) error {
	var opts = metav1.ListOptions{}
	list, err := rc.namespaceUC.Get(c.Request().Context(), opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"data": list})
}

// Create godoc
//
//	@Summary		Create a new namespace
//	@Description	Creates a new namespace in the Kubernetes cluster.
//	@Tags			namespaces
//	@Accept			json
//	@Produce		json
//	@Param			namespace	body		model.NamespaceRequest	true	"Namespace request body"
//	@Success		201			{object}	map[string]string		"Successfully created namespace"
//	@Failure		400			{object}	map[string]string		"Bad request or error message"
//	@Router			/namespaces [post]
func (rc *NamespaceHandler) Create(c echo.Context) error {
	var input model.NamespaceRequest

	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	namespace, err := rc.namespaceUC.Create(c.Request().Context(), &input.Namespace, input.Opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"namespace": namespace.Name})
}

// GetByNameOrUID godoc
//
//	@Summary		Get a namespace by name or UID
//	@Description	Retrieves a namespace from the Kubernetes cluster by its name or UID.
//	@Tags			namespaces
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string					true	"Name or UID of the namespace"
//	@Success		200	{object}	map[string]interface{}	"Details of the requested namespace"
//	@Failure		400	{object}	map[string]string		"Bad request or error message"
//	@Router			/namespaces/{id} [get]
func (rc *NamespaceHandler) GetByNameOrUID(c echo.Context) error {
	var nameOrUID = c.Param("id")

	var opts = metav1.ListOptions{}

	list, err := rc.namespaceUC.GetByNameOrUID(c.Request().Context(), nameOrUID, opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"data": list})
}

// Delete godoc
//
//	@Summary		Delete a namespace by name or UID
//	@Description	Deletes a namespace from the Kubernetes cluster by its name or UID.
//	@Tags			namespaces
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string				true	"Name or UID of the namespace"
//	@Success		200	{string}	string				"Success message"
//	@Failure		400	{object}	map[string]string	"Bad request or error message"
//	@Router			/namespaces/{id} [delete]
func (rc *NamespaceHandler) Delete(c echo.Context) error {
	var nameOrUID = c.Param("id")

	var opts = metav1.DeleteOptions{}

	if err := rc.namespaceUC.Delete(c.Request().Context(), nameOrUID, opts); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, "deleted succesfully")
}

// Update godoc
//
//	@Summary		Update an existing namespace
//	@Description	Updates an existing namespace in the Kubernetes cluster.
//	@Tags			namespaces
//	@Accept			json
//	@Produce		json
//	@Param			namespace	body		model.NamespaceRequest	true	"Namespace request body"
//	@Success		200			{object}	map[string]string		"Successfully updated the namespace"
//	@Failure		400			{object}	map[string]string		"Bad request or invalid data"
//	@Router			/namespaces [put]
func (rc *NamespaceHandler) Update(c echo.Context) error {
	var request model.NamespaceRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	namespace, err := rc.namespaceUC.Update(c.Request().Context(), &request.Namespace, metav1.UpdateOptions(request.Opts))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"namespace": namespace.Name})
}
