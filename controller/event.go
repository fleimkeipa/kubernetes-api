package controller

import (
	"fmt"
	"net/http"

	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/uc"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type EventHandler struct {
	eventsUC *uc.EventUC
	logger   *zap.SugaredLogger
}

func NewEventHandler(eventsUC *uc.EventUC, logger *zap.SugaredLogger) *EventHandler {
	return &EventHandler{
		eventsUC: eventsUC,
		logger:   logger,
	}
}

// List godoc
//
//	@Summary		List events
//	@Description	Retrieves a list of events from the database.
//	@Tags			events
//	@Accept			json
//	@Produce		json
//	@Param			kind			query		string			false	"kind to filter events by"
//	@Param			event_kind		query		string			false	"event kind to filter events by"
//	@Param			creation_time	query		string			false	"creation time to filter events by"
//	@Param			owner_id		query		string			false	"owner id to filter events by"
//	@Param			owner_username	query		string			false	"owner username to filter events by"
//	@Success		200				{object}	SuccessResponse	"List of events"
//	@Failure		500				{object}	FailureResponse "Interval error"
//	@Router			/events [get]
func (rc *EventHandler) List(c echo.Context) error {
	// Extract filtering options from the query parameters
	var opts = rc.getEventsFindOpts(c)

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

func (rc *EventHandler) getEventsFindOpts(c echo.Context) model.EventFindOpts {
	return model.EventFindOpts{
		PaginationOpts: getPagination(c),
		Kind:           getFilter(c, "kind"),
		EventKind:      getFilter(c, "event_kind"),
		CreationTime:   getFilter(c, "creation_time"),
		OwnerID:        getFilter(c, "owner_id"),
		OwnerUsername:  getFilter(c, "owner_username"),
	}
}
