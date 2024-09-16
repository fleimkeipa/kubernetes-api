package controller

import (
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
//	@Param			kind			query		string					false	"kind to filter events by"
//	@Param			event_kind		query		string					false	"event kind to filter events by"
//	@Param			creation_time	query		string					false	"creation time to filter events by"
//	@Param			owner_id		query		string					false	"owner id to filter events by"
//	@Param			owner_username	query		string					false	"owner username to filter events by"
//	@Success		200				{object}	map[string]interface{}	"List of events"
//	@Failure		400				{object}	map[string]string		"Bad request or invalid data"
//	@Router			/events [get]
func (rc *EventHandler) List(c echo.Context) error {
	var opts = rc.getEventsFindOpts(c)

	list, err := rc.eventsUC.List(c.Request().Context(), &opts)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"data": list})
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
