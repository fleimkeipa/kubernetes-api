package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fleimkeipa/kubernetes-api/controller"
	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/repositories"
	"github.com/fleimkeipa/kubernetes-api/uc"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name       string
		wantJSON   string
		request    model.UserRequest
		wantStatus int
	}{
		{
			name: "successful creation",
			request: model.UserRequest{
				Username: "test-user-1",
				Email:    "test@example.com",
				Password: "password",
				RoleID:   7,
			},
			wantStatus: http.StatusCreated,
			wantJSON:   `{"data":"test-user-1","message":"User created successfully."}`,
		},
		{
			name: "invalid request data",
			request: model.UserRequest{
				Username: "test-user-1",
				Email:    "test@example.com",
				Password: "password",
				RoleID:   7,
			},
			wantStatus: http.StatusBadRequest,
			wantJSON:   `{"error":"Failed to bind request: ...","message":"Invalid request data. Please check your input and try again."}`,
		},
		{
			name: "user creation failure",
			request: model.UserRequest{
				Username: "test-user-1",
				Email:    "test@example.com",
				Password: "password",
				RoleID:   7,
			},
			wantStatus: http.StatusInternalServerError,
			wantJSON:   `{"error":"Failed to create user: user creation failed","message":"User creation failed. Please verify the details and try again."}`,
		},
		{
			name: "nil echo.Context object",
			request: model.UserRequest{
				Username: "test-user-1",
				Email:    "test@example.com",
				Password: "password",
				RoleID:   7,
			},
			wantStatus: http.StatusInternalServerError,
			wantJSON:   `{"error":"...","message":"..."}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create a new echo context
			e := echo.New()
			reqBody, err := json.Marshal(tt.request)
			assert.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			eventRepo := repositories.NewEventRepository(test_db)
			eventUC := uc.NewEventUC(eventRepo)

			// create a new UserHandler
			userRepo := repositories.NewUserRepository(test_db)
			userUC := uc.NewUserUC(userRepo, eventUC)
			userHandler := controller.NewUserHandlers(userUC)

			// call the Create function
			err = userHandler.CreateUser(c)
			assert.NoError(t, err)

			// check the status code
			assert.Equal(t, tt.wantStatus, rec.Code)

			// check the response JSON
			var respJSON string
			err = json.Unmarshal(rec.Body.Bytes(), &respJSON)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, respJSON)
		})
	}
}
