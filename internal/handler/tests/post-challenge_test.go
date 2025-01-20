package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gleb-korostelev/GophKeeper/internal/handler"
	MockService "github.com/gleb-korostelev/GophKeeper/mocks"
	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestPostChallenge(t *testing.T) {
	mc := minimock.NewController(t)

	mockAuthSvc := MockService.NewAuthSvcMock(mc)

	tests := []struct {
		name           string
		setupMocks     func()
		requestBody    interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Successful challenge generation",
			setupMocks: func() {
				mockAuthSvc.GetChallengeMock.Expect(
					minimock.AnyContext,
					models.Profile{Username: "test_user", Password: ""},
				).Return("random_challenge", nil)
			},
			requestBody: map[string]string{
				"username": "test_user",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": map[string]interface{}{
					"challenge": "random_challenge",
				},
				"message": "Success",
				"success": true,
			},
		},
		{
			name:           "Invalid JSON in request body",
			setupMocks:     func() {},
			requestBody:    "invalid_json",
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"success": false,
				"message": "invalid request body",
			},
		},
		{
			name: "Error generating challenge",
			setupMocks: func() {
				mockAuthSvc.GetChallengeMock.Expect(
					minimock.AnyContext,
					models.Profile{Username: "test_user", Password: ""},
				).Return("", errors.New("challenge generation error"))
			},
			requestBody: map[string]string{
				"username": "test_user",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"success": false,
				"message": "challenge generation error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			h := &handler.Implementation{
				AuthSvc: mockAuthSvc,
			}

			var reqBody []byte
			if body, ok := tt.requestBody.(map[string]string); ok {
				reqBody, _ = json.Marshal(body)
			} else {
				reqBody = []byte(tt.requestBody.(string))
			}

			req := httptest.NewRequest("POST", "/api/v1/challenge", bytes.NewBuffer(reqBody))
			rec := httptest.NewRecorder()

			h.PostChallenge(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			expectedJSON, _ := json.Marshal(tt.expectedBody)
			assert.JSONEq(t, string(expectedJSON), rec.Body.String())
		})
	}
}
