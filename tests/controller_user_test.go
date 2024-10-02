package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fleimkeipa/kubernetes-api/controller"
	"github.com/fleimkeipa/kubernetes-api/model"
	"github.com/fleimkeipa/kubernetes-api/pkg"
	"github.com/fleimkeipa/kubernetes-api/repositories"
	"github.com/fleimkeipa/kubernetes-api/uc"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	// Setup
	test_db, terminateDB := pkg.GetTestInstance(context.TODO())
	defer terminateDB()

	userRepo := repositories.NewUserRepository(test_db)
	userUC := uc.NewUserUC(userRepo)

	rc := controller.NewUserHandlers(userUC)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte(`{"username":"test","email":"test@example.com","password":"password","role_id":1}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Test with valid input
	if userUC == nil {
		t.Fatal("userUC is nil")
	}
	userUC.On(
		"Create",
		mock.Anything,
		model.User{
			Username: "test",
			Email:    "test@example.com",
			Password: "password",
			RoleID:   1,
		},
	).Return(nil)
	err := rc.CreateUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	var resp controller.SuccessResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "test", resp.Data)

	// Test with invalid input (binding error)
	req = httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte(`{"username":"test","email":"test@example.com","password":""}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	err = rc.CreateUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	var failureResp controller.FailureResponse
	json.Unmarshal(rec.Body.Bytes(), &failureResp)
	assert.NotNil(t, failureResp.Error)

	// Test with internal server error (userUC.Create error)
	if userUC == nil {
		t.Fatal("userUC is nil")
	}
	userUC.On("Create", mock.Anything, model.User{Username: "test", Email: "test@example.com", Password: "password", RoleID: 1}).Return(errors.New("internal error"))
	req = httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer([]byte(`{"username":"test","email":"test@example.com","password":"password","role_id":1}`)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	err = rc.CreateUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &failureResp)
	assert.NotNil(t, failureResp.Error)
}
