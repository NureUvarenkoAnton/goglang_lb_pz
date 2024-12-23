package transport

//
// import (
// 	"bytes"
// 	"database/sql"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"
//
// 	"NureUvarenkoAnton/apzkr-pzpi-21-7-uvarenko-anton/Task2/apzkr-pzpi-21-7-uvarenko-anton-task2/internal/core"
// 	"NureUvarenkoAnton/apzkr-pzpi-21-7-uvarenko-anton/Task2/apzkr-pzpi-21-7-uvarenko-anton-task2/internal/pkg"
// 	"NureUvarenkoAnton/apzkr-pzpi-21-7-uvarenko-anton/Task2/apzkr-pzpi-21-7-uvarenko-anton-task2/internal/pkg/jwt"
// 	"NureUvarenkoAnton/apzkr-pzpi-21-7-uvarenko-anton/Task2/apzkr-pzpi-21-7-uvarenko-anton-task2/internal/pkg/middleware"
//
// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )
//
// func TestWalkHalder_CreateWalkRequest(t *testing.T) {
// 	mockWalkService := new(mockIWalkService)
// 	handler := NewWalkHandler(mockWalkService)
// 	jwtHandler := jwt.NewJWT("abc")
//
// 	token, _ := jwtHandler.GenUserToken(1, core.UsersUserTypeAdmin, time.Now().Add(time.Hour))
//
// 	gin.SetMode(gin.TestMode)
// 	r := gin.Default()
// 	r.Use(middleware.TokenVerifier(*jwtHandler, []core.UsersUserType{core.UsersUserTypeAdmin}))
// 	r.POST("/walks", handler.CreateWalkRequest)
//
// 	tests := []struct {
// 		name         string
// 		userId       int
// 		authToken    string
// 		requestBody  interface{}
// 		mockResponse error
// 		expectedCode int
// 	}{
// 		{"Unauthorized", 0, "", nil, nil, http.StatusUnauthorized},
// 		{"BadRequest", 1, token, "invalid json", nil, http.StatusBadRequest},
// 		{"InternalServerError", 1, token, map[string]interface{}{"walkedId": 1, "petId": 1, "startTime": time.Now()}, errors.New("some error"), http.StatusInternalServerError},
// 		{"Success", 1, token, map[string]interface{}{"walkedId": 1, "petId": 1, "startTime": time.Now()}, nil, http.StatusOK},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var reqBody *bytes.Buffer
// 			if tt.requestBody != nil {
// 				if str, ok := tt.requestBody.(string); ok {
// 					reqBody = bytes.NewBuffer([]byte(str))
// 				} else {
// 					body, _ := json.Marshal(tt.requestBody)
// 					reqBody = bytes.NewBuffer(body)
// 				}
// 			} else {
// 				reqBody = new(bytes.Buffer)
// 			}
//
// 			req, _ := http.NewRequest(http.MethodPost, "/walks", reqBody)
// 			w := httptest.NewRecorder()
// 			req.Header.Set("Content-Type", "application/json")
// 			req.Header.Set("Authorization", "Bearer "+tt.authToken)
//
// 			if tt.userId != 0 {
// 				req.Header.Set("user_id", string(rune(tt.userId)))
// 			}
//
// 			if tt.name != "Unauthorized" && tt.name != "BadRequest" {
// 				mockWalkService.On("CreateWalk", mock.Anything, mock.Anything).Return(tt.mockResponse).Once()
// 			}
//
// 			r.ServeHTTP(w, req)
//
// 			assert.Equal(t, tt.expectedCode, w.Code)
// 		})
// 	}
// }
//
// func TestWalkHalder_GetWalksByWalkerId(t *testing.T) {
// 	mockWalkService := new(mockIWalkService)
// 	handler := NewWalkHandler(mockWalkService)
//
// 	gin.SetMode(gin.TestMode)
// 	r := gin.Default()
// 	r.GET("/walks/walker/:walkerId", handler.GetWalksByParams)
//
// 	jwtHandler := jwt.NewJWT("abc")
// 	token, _ := jwtHandler.GenUserToken(1, core.UsersUserTypeAdmin, time.Now().Add(time.Hour))
//
// 	tests := []struct {
// 		name         string
// 		uri          string
// 		mockResponse interface{}
// 		mockError    error
// 		expectedCode int
// 	}{
// 		{"BadRequest", "/walks/walker/invalid", nil, nil, http.StatusBadRequest},
// 		{"NotFound", "/walks/walker/1", nil, pkg.ErrNotFound, http.StatusNotFound},
// 		{"InternalServerError", "/walks/walker/1", nil, errors.New("some error"), http.StatusInternalServerError},
// 		{"Success", "/walks/walker/1", []core.Walk{}, nil, http.StatusOK},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			req, _ := http.NewRequest(http.MethodGet, tt.uri, nil)
// 			req.Header.Set("Authorization", "Bearer "+token)
//
// 			w := httptest.NewRecorder()
//
// 			if tt.name != "BadRequest" {
// 				mockWalkService.On("GetWalksByWalkerId", mock.Anything, sql.NullInt64{Int64: 1, Valid: true}).Return(tt.mockResponse, tt.mockError).Once()
// 			}
//
// 			r.ServeHTTP(w, req)
//
// 			assert.Equal(t, tt.expectedCode, w.Code)
// 		})
// 	}
// }
//
// func TestWalkHalder_GetWalksByOwnerId(t *testing.T) {
// 	mockWalkService := new(mockIWalkService)
// 	handler := NewWalkHandler(mockWalkService)
//
// 	gin.SetMode(gin.TestMode)
// 	r := gin.Default()
// 	r.GET("/walks/owner/:ownerId", handler.GetWalksByOwnerId)
//
// 	tests := []struct {
// 		name         string
// 		uri          string
// 		mockResponse interface{}
// 		mockError    error
// 		expectedCode int
// 	}{
// 		{"BadRequest", "/walks/owner/invalid", nil, nil, http.StatusBadRequest},
// 		{"NotFound", "/walks/owner/1", nil, pkg.ErrNotFound, http.StatusNotFound},
// 		{"InternalServerError", "/walks/owner/1", nil, errors.New("some error"), http.StatusInternalServerError},
// 		{"Success", "/walks/owner/1", []core.Walk{}, nil, http.StatusOK},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			req, _ := http.NewRequest(http.MethodGet, tt.uri, nil)
// 			w := httptest.NewRecorder()
//
// 			if tt.name != "BadRequest" {
// 				mockWalkService.On("GetWalksByOwnerId", mock.Anything, sql.NullInt64{Int64: 1, Valid: true}).Return(tt.mockResponse, tt.mockError).Once()
// 			}
//
// 			r.ServeHTTP(w, req)
//
// 			assert.Equal(t, tt.expectedCode, w.Code)
// 		})
// 	}
// }
//
// func TestWalkHalder_UpdateWalkState(t *testing.T) {
// 	mockWalkService := new(mockIWalkService)
// 	handler := NewWalkHandler(mockWalkService)
//
// 	gin.SetMode(gin.TestMode)
// 	r := gin.Default()
// 	r.PUT("/walks/state", handler.UpdateWalkState)
//
// 	tests := []struct {
// 		name         string
// 		requestBody  interface{}
// 		mockResponse error
// 		expectedCode int
// 	}{
// 		{"BadRequest_InvalidJson", "invalid json", nil, http.StatusBadRequest},
// 		{"BadRequest_InvalidState", map[string]interface{}{"walkId": 1, "state": "invalid"}, nil, http.StatusBadRequest},
// 		{"NotFound", map[string]interface{}{"walkId": 1, "state": string(core.WalksStateAccepted)}, pkg.ErrNotFound, http.StatusNotFound},
// 		{"InternalServerError", map[string]interface{}{"walkId": 1, "state": string(core.WalksStateAccepted)}, errors.New("some error"), http.StatusInternalServerError},
// 		{"Success", map[string]interface{}{"walkId": 1, "state": string(core.WalksStateAccepted)}, nil, http.StatusOK},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var reqBody *bytes.Buffer
// 			if tt.requestBody != nil {
// 				if str, ok := tt.requestBody.(string); ok {
// 					reqBody = bytes.NewBuffer([]byte(str))
// 				} else {
// 					body, _ := json.Marshal(tt.requestBody)
// 					reqBody = bytes.NewBuffer(body)
// 				}
// 			} else {
// 				reqBody = new(bytes.Buffer)
// 			}
//
// 			req, _ := http.NewRequest(http.MethodPut, "/walks/state", reqBody)
// 			w := httptest.NewRecorder()
// 			req.Header.Set("Content-Type", "application/json")
//
// 			if tt.name != "BadRequest_InvalidJson" && tt.name != "BadRequest_InvalidState" {
// 				mockWalkService.On("UpdateWalkState", mock.Anything, mock.Anything).Return(tt.mockResponse).Once()
// 			}
//
// 			r.ServeHTTP(w, req)
//
// 			assert.Equal(t, tt.expectedCode, w.Code)
// 		})
// 	}
// }
