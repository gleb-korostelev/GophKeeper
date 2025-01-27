package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gleb-korostelev/GophKeeper/middleware"
	MockService "github.com/gleb-korostelev/GophKeeper/mocks"
	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gleb-korostelev/GophKeeper/models/profile"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestGetUserCards(t *testing.T) {
	mc := minimock.NewController(t)

	mockAuthSvc := MockService.NewAuthSvcMock(mc)
	mockProfileSvc := MockService.NewProfileSvcMock(mc)

	tests := []struct {
		name           string
		setupMocks     func()
		contextIssuer  string
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "Successful retrieval",
			setupMocks: func() {
				mockAuthSvc.GetAccountByUserNameMock.Expect(
					minimock.AnyContext, "test_user",
				).Return(models.Account{
					Username:    "test_user",
					AccountType: models.AccountAuthorizedUser,
				}, nil)

				mockProfileSvc.GetUserCardsMock.Expect(
					minimock.AnyContext, "test_user",
				).Return([]profile.CardInfo{
					{
						CardNumber:     "1234567812345678",
						CardHolder:     "John Doe",
						ExpirationDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						Cvv:            "123",
						Metadata:       "Primary card",
					},
					{
						CardNumber:     "8765432187654321",
						CardHolder:     "Jane Doe",
						ExpirationDate: time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
						Cvv:            "456",
						Metadata:       "Backup card",
					},
				}, nil)
			},
			contextIssuer:  "test_user",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": map[string]interface{}{
					"username": "test_user",
					"cards": []interface{}{
						map[string]interface{}{
							"card_holder":     "John Doe",
							"card_number":     "1234567812345678",
							"cvv":             "123",
							"expiration_date": "2025-01-01T00:00:00Z",
							"metadata":        "Primary card",
						},
						map[string]interface{}{
							"card_holder":     "Jane Doe",
							"card_number":     "8765432187654321",
							"cvv":             "456",
							"expiration_date": "2024-12-31T00:00:00Z",
							"metadata":        "Backup card",
						},
					},
				},
				"message": "Success",
				"success": true,
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
			contextIssuer:  "test_user",
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"success": false,
				"message": "not enough rights",
			},
		},
		{
			name: "Error retrieving user account",
			setupMocks: func() {
				mockAuthSvc.GetAccountByUserNameMock.Expect(
					minimock.AnyContext, "test_user",
				).Return(models.Account{}, errors.New("user not found"))
			},
			contextIssuer:  "test_user",
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"success": false,
				"message": "user not found",
			},
		},
		{
			name: "Error retrieving cards",
			setupMocks: func() {
				mockAuthSvc.GetAccountByUserNameMock.Expect(
					minimock.AnyContext, "test_user",
				).Return(models.Account{
					Username:    "test_user",
					AccountType: models.AccountAuthorizedUser,
				}, nil)

				mockProfileSvc.GetUserCardsMock.Expect(
					minimock.AnyContext, "test_user",
				).Return(nil, errors.New("database error"))
			},
			contextIssuer:  "test_user",
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"success": false,
				"message": "database error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			h := &Implementation{
				AuthSvc:    mockAuthSvc,
				ProfileSvc: mockProfileSvc,
			}

			req := httptest.NewRequest("GET", "/api/v1/cards", nil)

			ctx := context.WithValue(req.Context(), middleware.CtxKeyUserID, tt.contextIssuer)
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()

			h.GetUserCards(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			var expectedBody []byte
			if body, ok := tt.expectedBody.(map[string]interface{}); ok {
				expectedBody, _ = json.Marshal(body)
			} else {
				expectedBody, _ = json.Marshal(tt.expectedBody)
			}

			assert.JSONEq(t, string(expectedBody), rec.Body.String())
		})
	}
}
