package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Registration
func TestRegistration_Success(t *testing.T) {
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Mock the Server struct
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	server := &Server{
		Repository: mockRepo,
	}

	// Sample Registration request data
	body := generated.RegistrationParam{
		FullName:    "testFullName",
		Password:    "test@Password1",
		PhoneNumber: "+6282222222",
	}

	jsonBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/registration", bytes.NewReader(jsonBytes))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set up the expected behavior of the mock
	mockRepo.EXPECT().CreateNewUser(gomock.Any(), gomock.Any()).Return(
		repository.GetRegistrationOutput{Id: 1},
		nil,
	)

	// Call the Registration function
	err := server.Registration(c)

	assert.NoError(t, err)
}

func TestRegistration_Error_CreateNewUser(t *testing.T) {
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Mock the Server struct
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	server := &Server{
		Repository: mockRepo,
	}

	// Sample Registration request data
	body := generated.RegistrationParam{
		FullName:    "testFullName",
		Password:    "test@Password1",
		PhoneNumber: "+6282222222",
	}
	jsonBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/registration", bytes.NewReader(jsonBytes))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set up the expected behavior of the mock
	mockRepo.EXPECT().CreateNewUser(gomock.Any(), gomock.Any()).Return(
		repository.GetRegistrationOutput{},
		errors.New("err"),
	)

	// Call the Registration function
	err := server.Registration(c)

	assert.Nil(t, err)
}

func TestRegistration_Err_EmptyBody(t *testing.T) {
	e := echo.New()
	// Mock the Server struct
	server := &Server{}

	// Sample Registration request data
	body := generated.RegistrationParam{}
	jsonBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/registration", bytes.NewReader(jsonBytes))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the Registration function
	err := server.Registration(c)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	// Assert the response body (assuming JSON response)
	// You may need to update this based on the actual structure of RegistrationResponse
	assert.NotNil(t, rec.Body.String()) // TODO make compare return value
}

// Login
func TestLogin_Success(t *testing.T) {
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Mock the Server struct
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	server := &Server{
		Repository: mockRepo,
	}

	// Sample login request data
	body := generated.LoginParam{
		Password:    "my@Password1",
		PhoneNumber: "+6282222222",
	}

	jsonBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBytes))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set up the expected behavior of the mock
	mockRepo.EXPECT().GetUserByPhoneNumber(gomock.Any(), gomock.Any()).Return(
		repository.GetLoginOutput{Id: 1, Password: "$2a$04$BN7qD4ROTQKOoagz6Ez5xucaSFNkKWYhT9UJF7pd4jgKvaRsLBKFW"},
		nil,
	)
	mockRepo.EXPECT().UpdateUserSuccesLogin(gomock.Any(), gomock.Any()).Return(
		nil,
	)

	// Call the Registration function
	err := server.Login(c)

	assert.NoError(t, err)
}

func TestLogin_Error_GetUserByPhoneNumber(t *testing.T) {
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Mock the Server struct
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	server := &Server{
		Repository: mockRepo,
	}

	// Sample login request data
	body := generated.LoginParam{
		Password:    "my@Password1",
		PhoneNumber: "+6282222222",
	}

	jsonBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBytes))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set up the expected behavior of the mock
	mockRepo.EXPECT().GetUserByPhoneNumber(gomock.Any(), gomock.Any()).Return(
		repository.GetLoginOutput{},
		errors.New("err"),
	)

	// Call the login function
	err := server.Login(c)

	assert.NoError(t, err)
}

func TestLogin_Error_UpdateUserSuccesLogin(t *testing.T) {
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Mock the Server struct
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	server := &Server{
		Repository: mockRepo,
	}

	// Sample login request data
	body := generated.LoginParam{
		Password:    "my@Password1",
		PhoneNumber: "+6282222222",
	}

	jsonBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBytes))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set up the expected behavior of the mock
	mockRepo.EXPECT().GetUserByPhoneNumber(gomock.Any(), gomock.Any()).Return(
		repository.GetLoginOutput{Id: 1, Password: "$2a$04$BN7qD4ROTQKOoagz6Ez5xucaSFNkKWYhT9UJF7pd4jgKvaRsLBKFW"},
		nil,
	)
	mockRepo.EXPECT().UpdateUserSuccesLogin(gomock.Any(), gomock.Any()).Return(
		errors.New("Err"),
	)

	// Call the login function
	err := server.Login(c)

	assert.NoError(t, err)
}

func TestLogin_Error_comparePasswords(t *testing.T) {
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Mock the Server struct
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	server := &Server{
		Repository: mockRepo,
	}

	// Sample login request data
	body := generated.LoginParam{
		Password:    "my@Password11",
		PhoneNumber: "+6282222222",
	}

	jsonBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBytes))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set up the expected behavior of the mock
	mockRepo.EXPECT().GetUserByPhoneNumber(gomock.Any(), gomock.Any()).Return(
		repository.GetLoginOutput{Id: 1, Password: "$2a$04$BN7qD4ROTQKOoagz6Ez5xucaSFNkKWYhT9UJF7pd4jgKvaRsLBKFW"},
		nil,
	)

	// Call the login function
	err := server.Login(c)

	assert.NoError(t, err)
}

// my profile
func TestMyProfile_Success(t *testing.T) {
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Mock the Server struct
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	server := &Server{
		Repository: mockRepo,
	}

	req := httptest.NewRequest(http.MethodGet, "/my-profile", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCIsInBob25lTnVtYmVyIjoiKzYyODIyMjIyMjIiLCJleHAiOjE3MDcxNTEzNTJ9.89iLa6pq7tUKUWsWlagLi7uMkRcTDSc0XIP8KOw2BvU"))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the Registration function
	err := server.MyProfile(c)

	assert.NoError(t, err)
}

func TestMyProfile_malformed(t *testing.T) {
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Mock the Server struct
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	server := &Server{
		Repository: mockRepo,
	}

	req := httptest.NewRequest(http.MethodGet, "/my-profile", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", "test"))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Call the Registration function
	err := server.MyProfile(c)

	assert.NoError(t, err)
}

// Update Profile
func TestUpdateProfile_Success(t *testing.T) {
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Mock the Server struct
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	server := &Server{
		Repository: mockRepo,
	}

	// Sample login request data
	var mockFullName = "MyName"
	var mockPhoneNumber = "+6282222222"
	body := generated.UpdateProfileParam{
		FullName:    &mockFullName,
		PhoneNumber: &mockPhoneNumber,
	}

	jsonBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/update-profile", bytes.NewReader(jsonBytes))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCIsInBob25lTnVtYmVyIjoiKzYyODIyMjIyMjIiLCJleHAiOjE3MDcxNTEzNTJ9.89iLa6pq7tUKUWsWlagLi7uMkRcTDSc0XIP8KOw2BvU"))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set up the expected behavior of the mock
	mockRepo.EXPECT().UpdateUserByPhoneNumber(gomock.Any(), gomock.Any()).Return(
		repository.UpdateUserOutput{Id: 1},
		nil,
	)

	// Call the Registration function
	err := server.UpdateProfile(c)

	assert.NoError(t, err)
}

func TestUpdateProfile_Error_Conflict(t *testing.T) {
	e := echo.New()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Mock the Server struct
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	server := &Server{
		Repository: mockRepo,
	}

	// Sample login request data
	var mockFullName = "MyName"
	var mockPhoneNumber = "+6282222222"
	body := generated.UpdateProfileParam{
		FullName:    &mockFullName,
		PhoneNumber: &mockPhoneNumber,
	}

	jsonBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/update-profile", bytes.NewReader(jsonBytes))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCIsInBob25lTnVtYmVyIjoiKzYyODIyMjIyMjIiLCJleHAiOjE3MDcxNTEzNTJ9.89iLa6pq7tUKUWsWlagLi7uMkRcTDSc0XIP8KOw2BvU"))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Set up the expected behavior of the mock
	mockRepo.EXPECT().UpdateUserByPhoneNumber(gomock.Any(), gomock.Any()).Return(
		repository.UpdateUserOutput{Id: 1},
		errors.New("err"),
	)

	// Call the Registration function
	err := server.UpdateProfile(c)

	assert.NoError(t, err)
}
