package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gleb-korostelev/GophKeeper/internal/handler"
	"github.com/gleb-korostelev/GophKeeper/middleware"
	MockService "github.com/gleb-korostelev/GophKeeper/mocks"
	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/models/profile"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestPostUploadInfo(t *testing.T) {
	mc := minimock.NewController(t)

	mockAuthSvc := MockService.NewAuthSvcMock(mc)
	mockProfileSvc := MockService.NewProfileSvcMock(mc)

	tests := []struct {
		name           string
		setupMocks     func()
		contextIssuer  string
		requestBody    interface{}
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Successful upload",
			setupMocks: func() {
				mockAuthSvc.GetAccountByUserNameMock.Expect(
					minimock.AnyContext, "test_user",
				).Return(models.Account{
					Username:    "test_user",
					AccountType: models.AccountAuthorizedUser,
				}, nil)

				mockProfileSvc.UploadInfoMock.Expect(
					minimock.AnyContext,
					profile.CardInfo{
						Username:       "test_user",
						CardNumber:     "1234567812345678",
						CardHolder:     "John Doe",
						ExpirationDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						Cvv:            "123",
						Metadata:       "Test metadata",
					},
				).Return(nil)
			},
			contextIssuer: "test_user",
			requestBody: map[string]interface{}{
				"card_number":     "1234567812345678",
				"card_holder":     "John Doe",
				"expiration_date": "2025-01-01T00:00:00Z",
				"cvv":             "123",
				"metadata":        "Test metadata",
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"success": true,
				"message": "Success",
			},
		},
		{
			name:           "Invalid JSON request body",
			setupMocks:     func() {},
			contextIssuer:  "test_user",
			requestBody:    "invalid_json",
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"success": false,
				"message": "invalid request body",
			},
		},
		{
			name: "Unauthorized user",
			setupMocks: func() {
				mockAuthSvc.GetAccountByUserNameMock.Expect(
					minimock.AnyContext, "test_user",
				).Return(models.Account{
					Username:    "test_user",
					AccountType: models.AccountUnauthorizedUser,
				}, nil)
			},
			contextIssuer: "test_user",
			requestBody: map[string]interface{}{
				"card_number":     "1234567812345678",
				"card_holder":     "John Doe",
				"expiration_date": "2025-01-01T00:00:00Z",
				"cvv":             "123",
				"metadata":        "Test metadata",
			},
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"success": false,
				"message": "not enough rights",
			},
		},
		{
			name: "Error in UploadInfo",
			setupMocks: func() {
				mockAuthSvc.GetAccountByUserNameMock.Expect(
					minimock.AnyContext, "test_user",
				).Return(models.Account{
					Username:    "test_user",
					AccountType: models.AccountAuthorizedUser,
				}, nil)

				mockProfileSvc.UploadInfoMock.Expect(
					minimock.AnyContext,
					profile.CardInfo{
						Username:       "test_user",
						CardNumber:     "1234567812345678",
						CardHolder:     "John Doe",
						ExpirationDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						Cvv:            "123",
						Metadata:       "Test metadata",
					},
				).Return(errors.New("upload error"))
			},
			contextIssuer: "test_user",
			requestBody: map[string]interface{}{
				"card_number":     "1234567812345678",
				"card_holder":     "John Doe",
				"expiration_date": "2025-01-01T00:00:00Z",
				"cvv":             "123",
				"metadata":        "Test metadata",
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"success": false,
				"message": "upload error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			h := &handler.Implementation{
				AuthSvc:    mockAuthSvc,
				ProfileSvc: mockProfileSvc,
			}

			var reqBody []byte
			if body, ok := tt.requestBody.(map[string]interface{}); ok {
				reqBody, _ = json.Marshal(body)
			} else {
				reqBody = []byte(tt.requestBody.(string))
			}

			req := httptest.NewRequest("POST", "/api/v1/upload-card-info", bytes.NewBuffer(reqBody))
			ctx := context.WithValue(req.Context(), middleware.CtxKeyAddress, tt.contextIssuer)
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()

			h.PostUploadInfo(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			expectedJSON, _ := json.Marshal(tt.expectedBody)
			assert.JSONEq(t, string(expectedJSON), rec.Body.String())
		})
	}
}
