package controller

import (
	"fmt"
	"net/http"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/uc"

	"github.com/labstack/echo/v4"
)

type EventHandler struct {
	eventsUC *uc.EventUC
}

func NewEventHandler(eventsUC *uc.EventUC) *EventHandler {
	return &EventHandler{
		eventsUC: eventsUC,
	}
}

// List godoc
//
//	@Summary		List events
//	@Description	Retrieves a list of events from the database.
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			kind			query		string			false	"kind to filter events by"
//	@Param			event_kind		query		string			false	"event kind to filter events by"
//	@Param			creation_time	query		string			false	"creation time to filter events by"
//	@Param			owner_id		query		string			false	"owner id to filter events by"
//	@Param			owner_username	query		string			false	"owner username to filter events by"
//	@Success		200				{object}	SuccessResponse	"List of events"
//	@Failure		500				{object}	FailureResponse	"Interval error"
//	@Router			/events [get]
func (rc *EventHandler) List(c echo.Context) error {
	// Extract filtering options from the query parameters
	opts := rc.getEventsFindOpts(c)

	// Attempt to retrieve the list of events
	list, err := rc.eventsUC.List(c.Request().Context(), &opts)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, FailureResponse{
			Error:   fmt.Sprintf("Failed to retrieve events: %v", err),
			Message: "There was an error fetching the events. Please verify the filters and try again.",
		})
	}

	// Return the list of events if successful
	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    list,
		Message: "Events retrieved successfully.",
	})
}

// GetByID godoc
//
//	@Summary		Get a event by ID
//	@Description	Retrieves a event from Database by its ID, optionally filtered by namespace.
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string			true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			id				path		string			true	"ID of the event"
//	@Success		200				{object}	SuccessResponse	"Details of the requested event"
//	@Failure		500				{object}	FailureResponse	"Interval error"
//	@Router			/events/{id} [get]
func (rc *EventHandler) GetByID(c echo.Context) error {
	nameOrUID := c.Param("id")

	event, err := rc.eventsUC.GetByID(c.Request().Context(), nameOrUID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, FailureResponse{
			Error:   fmt.Sprintf("Failed to retrieve event: %v", err),
			Message: "Error fetching the event details. Please verify the event name or UID and try again.",
		})
	}

	return c.JSON(http.StatusOK, SuccessResponse{
		Data:    event,
		Message: "Event retrieved successfully.",
	})
}

func (rc *EventHandler) getEventsFindOpts(c echo.Context) model.EventFindOpts {
	return model.EventFindOpts{
		PaginationOpts: getPagination(c),
		Category:       getFilter(c, "kind"),
		Type:           getFilter(c, "event_kind"),
		CreatedAt:      getFilter(c, "created_at"),
		OwnerID:        getFilter(c, "owner_id"),
		OwnerUsername:  getFilter(c, "owner_username"),
	}
}
