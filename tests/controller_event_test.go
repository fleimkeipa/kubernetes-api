package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fleimkeipa/kubernetes-api/controller"
	"github.com/fleimkeipa/kubernetes-api/pkg"
	"github.com/fleimkeipa/kubernetes-api/repositories"
	"github.com/fleimkeipa/kubernetes-api/uc"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestEventHandler_List(t *testing.T) {
	// Mock EventUC
	test_db, terminateDB = pkg.GetTestInstance(context.TODO())
	defer terminateDB()

	eventRepo := repositories.NewEventRepository(test_db)
	eventUC := uc.NewEventUC(eventRepo)
	handler := controller.NewEventHandler(eventUC)

	// Test successful event retrieval with no filters
	t.Run("successful retrieval with no filters", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/events", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err = handler.List(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp controller.SuccessResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp.Data)
	})

	// Test successful event retrieval with filters
	t.Run("successful retrieval with filters", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/events?kind=pod&event_kind=create", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err = handler.List(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		var resp controller.SuccessResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp.Data)
	})

	// Test error handling for event retrieval failure
	t.Run("error handling for event retrieval failure", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/events", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err = handler.List(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		var resp controller.FailureResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp.Error)
	})

	// Test error handling for invalid query parameters
	t.Run("error handling for invalid query parameters", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/events?invalid_param=invalid_value", nil)
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err = handler.List(c)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		var resp controller.FailureResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.NotNil(t, resp.Error)
	})
}
