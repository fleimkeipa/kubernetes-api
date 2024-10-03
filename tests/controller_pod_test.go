package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fleimkeipa/kubernetes-api/controller"
	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/uc"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPodHandlerCreate(t *testing.T) {
	tests := []struct {
		name       string
		request    model.PodsCreateRequest
		podsUC     *uc.PodUC
		wantStatus int
		wantJSON   string
	}{
		{
			name: "successful creation",
			request: model.PodsCreateRequest{
				Pod: model.Pod{
					ObjectMeta: model.ObjectMeta{
						Name: "test-pod",
					},
				},
			},
			podsUC:     &uc.PodUC{},
			wantStatus: http.StatusCreated,
			wantJSON:   `{"data":"test-pod","message":"Pod created successfully."}`,
		},
		{
			name: "invalid request data",
			request: model.PodsCreateRequest{
				Pod: model.Pod{},
			},
			podsUC:     &uc.PodUC{},
			wantStatus: http.StatusBadRequest,
			wantJSON:   `{"error":"Failed to bind request: ...","message":"Invalid request data. Please check your input and try again."}`,
		},
		{
			name: "pod creation failure",
			request: model.PodsCreateRequest{
				Pod: model.Pod{
					ObjectMeta: model.ObjectMeta{
						Name: "test-pod",
					},
				},
			},
			podsUC:     &uc.PodUC{},
			wantStatus: http.StatusInternalServerError,
			wantJSON:   `{"error":"Failed to create pod: pod creation failed","message":"Pod creation failed. Please verify the details and try again."}`,
		},
		{
			name: "nil podsUC field",
			request: model.PodsCreateRequest{
				Pod: model.Pod{
					ObjectMeta: model.ObjectMeta{
						Name: "test-pod",
					},
				},
			},
			podsUC:     nil,
			wantStatus: http.StatusInternalServerError,
			wantJSON:   `{"error":"...","message":"..."}`,
		},
		{
			name: "nil echo.Context object",
			request: model.PodsCreateRequest{
				Pod: model.Pod{
					ObjectMeta: model.ObjectMeta{
						Name: "test-pod",
					},
				},
			},
			podsUC:     &uc.PodUC{},
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
			req := httptest.NewRequest(http.MethodPost, "/pods", bytes.NewBuffer(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// create a new PodHandler
			podHandler := controller.NewPodHandler(tt.podsUC)

			// call the Create function
			err = podHandler.Create(c)
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
