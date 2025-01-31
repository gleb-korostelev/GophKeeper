package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gleb-korostelev/GophKeeper/middleware"
	MockService "github.com/gleb-korostelev/GophKeeper/mocks"
	"github.com/gleb-korostelev/GophKeeper/models"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestDeleteCardInfo(t *testing.T) {
	mc := minimock.NewController(t)

	mockAuthSvc := MockService.NewAuthSvcMock(mc)
	mockProfileSvc := MockService.NewProfileSvcMock(mc)

	tests := []struct {
		name           string
		setupMocks     func()
		requestBody    interface{}
		contextIssuer  string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Successful deletion",
			setupMocks: func() {
				mockAuthSvc.GetAccountByUserNameMock.Expect(
					minimock.AnyContext, "test_user",
				).Return(models.Account{
					Username:    "test_user",
					AccountType: models.AccountAuthorizedUser,
				}, nil)

				mockProfileSvc.DeleteCardMock.Expect(
					minimock.AnyContext, "test_user", "1234567812345678",
				).Return(nil)
			},
			requestBody: map[string]string{
				"card_number": "1234567812345678",
			},
			contextIssuer:  "test_user",
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"success": true,
				"message": "Success",
			},
		},
		{
			name:       "Error retrieving issuer",
			setupMocks: func() {},
			requestBody: map[string]string{
				"card_number": "1234567812345678",
			},
			contextIssuer:  "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"success": false,
				"message": "bearer token is not correct",
			},
		},
		{
			name: "User not found",
			setupMocks: func() {
				mockAuthSvc.GetAccountByUserNameMock.Expect(
					minimock.AnyContext, "test_user",
				).Return(models.Account{}, errors.New("user not found"))
			},
			requestBody: map[string]string{
				"card_number": "1234567812345678",
			},
			contextIssuer:  "test_user",
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"success": false,
				"message": "user not found",
			},
		},
		{
			name: "Not enough rights",
			setupMocks: func() {
				mockAuthSvc.GetAccountByUserNameMock.Expect(
					minimock.AnyContext, "test_user",
				).Return(models.Account{
					Username:    "test_user",
					AccountType: models.AccountUnauthorizedUser,
				}, nil)
			},
			requestBody: map[string]string{
				"card_number": "1234567812345678",
			},
			contextIssuer:  "test_user",
			expectedStatus: http.StatusForbidden,
			expectedBody: map[string]interface{}{
				"success": false,
				"message": "not enough rights",
			},
		},
		{
			name: "Invalid request body",
			setupMocks: func() {
				mockAuthSvc.GetAccountByUserNameMock.Expect(
					minimock.AnyContext, "test_user",
				).Return(models.Account{
					Username:    "test_user",
					AccountType: models.AccountAuthorizedUser,
				}, nil)
			},
			requestBody:    "invalid_json",
			contextIssuer:  "test_user",
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"success": false,
				"message": "json: cannot unmarshal string into Go value of type handler.DeleteCardInfoReq",
			},
		},
		{
			name: "Error deleting card",
			setupMocks: func() {
				mockAuthSvc.GetAccountByUserNameMock.Expect(
					minimock.AnyContext, "test_user",
				).Return(models.Account{
					Username:    "test_user",
					AccountType: models.AccountAuthorizedUser,
				}, nil)

				mockProfileSvc.DeleteCardMock.Expect(
					minimock.AnyContext, "test_user", "1234567812345678",
				).Return(errors.New("card deletion error"))
			},
			requestBody: map[string]string{
				"card_number": "1234567812345678",
			},
			contextIssuer:  "test_user",
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"success": false,
				"message": "card deletion error",
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

			reqBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("DELETE", "/api/v1/cards", bytes.NewBuffer(reqBody))

			ctx := context.WithValue(req.Context(), middleware.CtxKeyUserID, tt.contextIssuer)
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()

			h.DeleteCardInfo(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			expectedJSON, _ := json.Marshal(tt.expectedBody)
			assert.JSONEq(t, string(expectedJSON), rec.Body.String())
		})
	}
}
