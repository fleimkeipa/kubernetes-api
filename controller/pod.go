package controller

import (
	"fmt"
	"net/http"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/uc"

	"github.com/labstack/echo/v4"
)

type PodHandler struct {
	podsUC *uc.PodUC
}

func NewPodHandler(podsUC *uc.PodUC) *PodHandler {
	return &PodHandler{
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
//	@Param			Authorization	header		string					true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			pod				body		model.PodsCreateRequest	true	"Pod request body"
//	@Success		201				{object}	SuccessResponse			"Successfully created the pod"
//	@Failure		400				{object}	FailureResponse			"Bad request or invalid data"
//	@Failure		500				{object}	FailureResponse			"Interval error"
//	@Router			/pods [post]
func (rc *PodHandler) Create(c echo.Context) error {
	var request model.PodsCreateRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to bind request: %v", err),
			Message: "Invalid request data. Please check your input and try again.",
		})
	}

	pod, err := rc.podsUC.Create(c.Request().Context(), request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to create pod: %v", err),
			Message: "Pod creation failed. Please verify the details and try again.",
		})
	}

	return c.JSON(http.StatusCreated, SuccessResponse{
		Data:    pod.Name,
		Message: "Pod created successfully.",
	})
}

// UpdatePod godoc
//
//	@Summary		Update an existing pod
//	@Description	Update specific fields of an existing pod in the Kubernetes cluster. The following fields are changeable:
//	@Description	- containers.image
//	@Description	- initContainers.image
//	@Description	- tolerations (only additions)
//	@Description	- activeDeadlineSeconds
//	@Description	- terminationGracePeriodSeconds
//	@Tags			pods
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			pod				body		model.PodsUpdateRequest	true	"Pod update request body"
//	@Param			namespace		query		string					false	"Namespace to filter the pod by"
//	@Param			id				path		string					true	"Name or UID of the pod"
//	@Success		200				{object}	SuccessResponse			"Pod successfully updated"
//	@Failure		400				{object}	FailureResponse			"Bad request or invalid input data"
//	@Failure		500				{object}	FailureResponse			"Interval error"
//	@Router			/pods/{id} [put]
func (rc *PodHandler) Update(c echo.Context) error {
	id := c.Param("id")
	namespace := c.QueryParam("namespace")

	var request model.PodsUpdateRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to bind request: %v", err),
			Message: "Invalid request data. Please check your input and try again.",
		})
	}

	pod, err := rc.podsUC.Update(c.Request().Context(), namespace, id, &request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to update pod: %v", err),
			Message: "Pod update failed. Please verify the details and try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    pod.Name,
		Message: "Pod updated successfully.",
	})
}

// List godoc
//
//	@Summary		List pods
//	@Description	Retrieves a list of pods from the Kubernetes cluster. You can filter results by namespace or paginate the response using the limit and continue parameters.
//	@Tags			pods
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			limit			query		string			false	"Maximum number of pods to retrieve"
//	@Param			continue		query		string			false	"Pagination token for fetching more pods"
//	@Param			namespace		query		string			false	"Namespace to filter pods by"
//	@Success		200				{object}	SuccessResponse	"List of pods"
//	@Failure		500				{object}	FailureResponse	"Interval error"
//	@Router			/pods [get]
func (rc *PodHandler) List(c echo.Context) error {
	namespace := c.QueryParam("namespace")

	opts := getKubeListOpts(c)

	list, err := rc.podsUC.List(c.Request().Context(), namespace, opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to retrieve pods: %v", err),
			Message: "Error fetching the list of pods. Please try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    list.ConvertMini(),
		Message: "Pods retrieved successfully.",
	})
}

// GetByNameOrUID godoc
//
//	@Summary		Get a pod by name or UID
//	@Description	Retrieves a pod from the Kubernetes cluster by its name or UID, optionally filtered by namespace.
//	@Tags			pods
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			namespace		query		string			false	"Namespace to filter the pod by"
//	@Param			id				path		string			true	"Name or UID of the pod"
//	@Success		200				{object}	SuccessResponse	"Details of the requested pod"
//	@Failure		500				{object}	FailureResponse	"Interval error"
//	@Router			/pods/{id} [get]
func (rc *PodHandler) GetByNameOrUID(c echo.Context) error {
	namespace := c.QueryParam("namespace")
	nameOrUID := c.Param("id")

	opts := model.ListOptions{}

	list, err := rc.podsUC.GetByNameOrUID(c.Request().Context(), namespace, nameOrUID, opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to retrieve pod: %v", err),
			Message: "Error fetching the pod details. Please verify the pod name or UID and try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    list,
		Message: "Pod retrieved successfully.",
	})
}

// Delete godoc
//
//	@Summary		Delete a pod by name or UID
//	@Description	Deletes a pod from the Kubernetes cluster by its name or UID, optionally filtered by namespace.
//	@Tags			pods
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			namespace		query		string			false	"Namespace to filter the pod by"
//	@Param			id				path		string			true	"Name or UID of the pod"
//	@Success		200				{string}	SuccessResponse	"Success message"
//	@Failure		500				{object}	FailureResponse	"Interval error"
//	@Router			/pods/{id} [delete]
func (rc *PodHandler) Delete(c echo.Context) error {
	namespace := c.QueryParam("namespace")
	nameOrUID := c.Param("id")

	opts := model.DeleteOptions{}

	if err := rc.podsUC.Delete(c.Request().Context(), namespace, nameOrUID, opts); err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to delete pod: %v", err),
			Message: "Error deleting the pod. Please verify the pod name or UID and try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Message: "Pod deleted successfully.",
	})
}
