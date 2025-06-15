package loan

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	domain "github.com/hinha/los-technical/internal/domain/loan"
	mock "github.com/hinha/los-technical/internal/domain/loan/mock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateLoan(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func(*mock.MockService)
		expectedStatus int
		expectedMsg    string
	}{
		{
			name: "Success",
			requestBody: map[string]interface{}{
				"borrower_id":      "borrower-123",
				"principal_amount": 1000.0,
				"rate":             5.0,
				"roi":              10.0,
			},
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().CreateLoan("borrower-123", 1000.0, 5.0, 10.0).Return(&domain.Loan{
					ID:              "loan-123",
					BorrowerID:      "borrower-123",
					PrincipalAmount: 1000.0,
					Rate:            5.0,
					ROI:             10.0,
					State:           domain.StateProposed,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedMsg:    "Loan created successfully",
		},
		{
			name: "Invalid Request - Missing Required Field",
			requestBody: map[string]interface{}{
				"borrower_id": "borrower-123",
				"rate":        5.0,
				"roi":         10.0,
				// Missing principal_amount
			},
			mockSetup:      func(mockService *mock.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Validation error",
		},
		{
			name: "Service Error",
			requestBody: map[string]interface{}{
				"borrower_id":      "borrower-123",
				"principal_amount": 1000.0,
				"rate":             5.0,
				"roi":              10.0,
			},
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().CreateLoan("borrower-123", 1000.0, 5.0, 10.0).Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Failed to create loan",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock.NewMockService(ctrl)
			tc.mockSetup(mockService)

			handler := NewHandler(mockService)

			// Create request
			reqBody, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/loans", bytes.NewBuffer(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e := echo.New()
			c := e.NewContext(req, rec)

			// Execute
			err := handler.CreateLoan(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, tc.expectedMsg, response["message"])
		})
	}
}

func TestGetLoan(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		loanID         string
		mockSetup      func(*mock.MockService)
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:   "Success",
			loanID: "loan-123",
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GetLoan("loan-123").Return(&domain.Loan{
					ID:              "loan-123",
					BorrowerID:      "borrower-123",
					PrincipalAmount: 1000.0,
					Rate:            5.0,
					ROI:             10.0,
					State:           domain.StateProposed,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "OK",
		},
		{
			name:   "Loan Not Found",
			loanID: "non-existent",
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GetLoan("non-existent").Return(nil, errors.New("loan not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedMsg:    "Loan not found",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock.NewMockService(ctrl)
			tc.mockSetup(mockService)

			handler := NewHandler(mockService)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/loans/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.loanID)

			// Execute
			err := handler.GetLoan(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, tc.expectedMsg, response["message"])
		})
	}
}

func TestApproveLoan(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		loanID         string
		requestBody    map[string]interface{}
		mockSetup      func(*mock.MockService)
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:   "Success",
			loanID: "loan-123",
			requestBody: map[string]interface{}{
				"validator_id": "validator-123",
				"proof_url":    "http://example.com/proof",
			},
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GetLoan("loan-123").Return(&domain.Loan{
					ID:    "loan-123",
					State: domain.StateProposed,
				}, nil)
				mockService.EXPECT().ApproveLoan("loan-123", "validator-123", "http://example.com/proof").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "OK",
		},
		{
			name:   "Invalid Request - Missing Required Field",
			loanID: "loan-123",
			requestBody: map[string]interface{}{
				"validator_id": "validator-123",
				// Missing proof_url
			},
			mockSetup:      func(mockService *mock.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Validation error",
		},
		{
			name:   "Invalid State",
			loanID: "loan-123",
			requestBody: map[string]interface{}{
				"validator_id": "validator-123",
				"proof_url":    "http://example.com/proof",
			},
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GetLoan("loan-123").Return(&domain.Loan{
					ID:    "loan-123",
					State: domain.StateApproved, // Already approved
				}, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "State validation error",
		},
		{
			name:   "Service Error",
			loanID: "loan-123",
			requestBody: map[string]interface{}{
				"validator_id": "validator-123",
				"proof_url":    "http://example.com/proof",
			},
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GetLoan("loan-123").Return(&domain.Loan{
					ID:    "loan-123",
					State: domain.StateProposed,
				}, nil)
				mockService.EXPECT().ApproveLoan("loan-123", "validator-123", "http://example.com/proof").Return(errors.New("service error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Failed to approve loan",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock.NewMockService(ctrl)
			tc.mockSetup(mockService)

			handler := NewHandler(mockService)

			// Create request
			reqBody, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/loans/:id/approve")
			c.SetParamNames("id")
			c.SetParamValues(tc.loanID)

			// Execute
			err := handler.ApproveLoan(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, tc.expectedMsg, response["message"])
		})
	}
}

func TestAddInvestment(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		loanID         string
		requestBody    map[string]interface{}
		mockSetup      func(*mock.MockService)
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:   "Success",
			loanID: "loan-123",
			requestBody: map[string]interface{}{
				"investor_id": "investor-123",
				"email":       "investor@example.com",
				"amount":      500.0,
			},
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GetLoan("loan-123").Return(&domain.Loan{
					ID:    "loan-123",
					State: domain.StateApproved,
				}, nil)
				mockService.EXPECT().AddInvestment("loan-123", "investor-123", "investor@example.com", 500.0).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "OK",
		},
		{
			name:   "Invalid Request - Missing Required Field",
			loanID: "loan-123",
			requestBody: map[string]interface{}{
				"investor_id": "investor-123",
				"email":       "investor@example.com",
				// Missing amount
			},
			mockSetup:      func(mockService *mock.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Validation error",
		},
		{
			name:   "Invalid Email Format",
			loanID: "loan-123",
			requestBody: map[string]interface{}{
				"investor_id": "investor-123",
				"email":       "invalid-email",
				"amount":      500.0,
			},
			mockSetup:      func(mockService *mock.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Validation error",
		},
		{
			name:   "Invalid State",
			loanID: "loan-123",
			requestBody: map[string]interface{}{
				"investor_id": "investor-123",
				"email":       "investor@example.com",
				"amount":      500.0,
			},
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GetLoan("loan-123").Return(&domain.Loan{
					ID:    "loan-123",
					State: domain.StateProposed, // Not approved yet
				}, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "State validation error",
		},
		{
			name:   "Service Error",
			loanID: "loan-123",
			requestBody: map[string]interface{}{
				"investor_id": "investor-123",
				"email":       "investor@example.com",
				"amount":      500.0,
			},
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GetLoan("loan-123").Return(&domain.Loan{
					ID:    "loan-123",
					State: domain.StateApproved,
				}, nil)
				mockService.EXPECT().AddInvestment("loan-123", "investor-123", "investor@example.com", 500.0).Return(errors.New("service error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Failed to add investment",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock.NewMockService(ctrl)
			tc.mockSetup(mockService)

			handler := NewHandler(mockService)

			// Create request
			reqBody, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/loans/:id/invest")
			c.SetParamNames("id")
			c.SetParamValues(tc.loanID)

			// Execute
			err := handler.AddInvestment(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, tc.expectedMsg, response["message"])
		})
	}
}

func TestDisburseLoan(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		loanID         string
		requestBody    map[string]interface{}
		mockSetup      func(*mock.MockService)
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:   "Success",
			loanID: "loan-123",
			requestBody: map[string]interface{}{
				"field_officer_id": "officer-123",
				"signed_agreement": "http://example.com/signed",
			},
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GetLoan("loan-123").Return(&domain.Loan{
					ID:    "loan-123",
					State: domain.StateInvested,
				}, nil)
				mockService.EXPECT().DisburseLoan("loan-123", "officer-123", "http://example.com/signed").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "OK",
		},
		{
			name:   "Invalid Request - Missing Required Field",
			loanID: "loan-123",
			requestBody: map[string]interface{}{
				"field_officer_id": "officer-123",
				// Missing signed_agreement
			},
			mockSetup:      func(mockService *mock.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Validation error",
		},
		{
			name:   "Invalid State",
			loanID: "loan-123",
			requestBody: map[string]interface{}{
				"field_officer_id": "officer-123",
				"signed_agreement": "http://example.com/signed",
			},
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GetLoan("loan-123").Return(&domain.Loan{
					ID:    "loan-123",
					State: domain.StateApproved, // Not invested yet
				}, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "State validation error",
		},
		{
			name:   "Service Error",
			loanID: "loan-123",
			requestBody: map[string]interface{}{
				"field_officer_id": "officer-123",
				"signed_agreement": "http://example.com/signed",
			},
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GetLoan("loan-123").Return(&domain.Loan{
					ID:    "loan-123",
					State: domain.StateInvested,
				}, nil)
				mockService.EXPECT().DisburseLoan("loan-123", "officer-123", "http://example.com/signed").Return(errors.New("service error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Failed to disburse loan",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock.NewMockService(ctrl)
			tc.mockSetup(mockService)

			handler := NewHandler(mockService)

			// Create request
			reqBody, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/loans/:id/disburse")
			c.SetParamNames("id")
			c.SetParamValues(tc.loanID)

			// Execute
			err := handler.DisburseLoan(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, tc.expectedMsg, response["message"])
		})
	}
}

func TestGenerateAgreementLetter(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		loanID         string
		requestBody    map[string]interface{}
		mockSetup      func(*mock.MockService)
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:   "Success",
			loanID: "loan-123",
			requestBody: map[string]interface{}{
				"letter_url": "http://example.com/letter",
			},
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GenerateAgreementLetter("loan-123", "http://example.com/letter").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "OK",
		},
		{
			name:        "Invalid Request - Missing Required Field",
			loanID:      "loan-123",
			requestBody: map[string]interface{}{
				// Missing letter_url
			},
			mockSetup:      func(mockService *mock.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Validation error",
		},
		{
			name:   "Service Error",
			loanID: "loan-123",
			requestBody: map[string]interface{}{
				"letter_url": "http://example.com/letter",
			},
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GenerateAgreementLetter("loan-123", "http://example.com/letter").Return(errors.New("service error"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Failed to generate agreement letter",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock.NewMockService(ctrl)
			tc.mockSetup(mockService)

			handler := NewHandler(mockService)

			// Create request
			reqBody, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/loans/:id/agreement")
			c.SetParamNames("id")
			c.SetParamValues(tc.loanID)

			// Execute
			err := handler.GenerateAgreementLetter(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, tc.expectedMsg, response["message"])
		})
	}
}

func TestGetLoansByBorrower(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		borrowerID     string
		mockSetup      func(*mock.MockService)
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:       "Success",
			borrowerID: "borrower-123",
			mockSetup: func(mockService *mock.MockService) {
				loans := []*domain.Loan{
					{
						ID:         "loan-1",
						BorrowerID: "borrower-123",
						State:      domain.StateProposed,
					},
					{
						ID:         "loan-2",
						BorrowerID: "borrower-123",
						State:      domain.StateApproved,
					},
				}
				mockService.EXPECT().GetLoansByBorrower("borrower-123").Return(loans, nil)
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "OK",
		},
		{
			name:       "No Loans Found",
			borrowerID: "borrower-456",
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GetLoansByBorrower("borrower-456").Return([]*domain.Loan{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "OK",
		},
		{
			name:       "Service Error",
			borrowerID: "borrower-123",
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GetLoansByBorrower("borrower-123").Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Failed to retrieve loans",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock.NewMockService(ctrl)
			tc.mockSetup(mockService)

			handler := NewHandler(mockService)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/loans/borrower/:borrowerId")
			c.SetParamNames("borrowerId")
			c.SetParamValues(tc.borrowerID)

			// Execute
			err := handler.GetLoansByBorrower(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, tc.expectedMsg, response["message"])
		})
	}
}

func TestGetLoansByState(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		state          string
		mockSetup      func(*mock.MockService)
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:  "Success",
			state: "APPROVED",
			mockSetup: func(mockService *mock.MockService) {
				loans := []*domain.Loan{
					{
						ID:    "loan-1",
						State: domain.StateApproved,
					},
					{
						ID:    "loan-2",
						State: domain.StateApproved,
					},
				}
				mockService.EXPECT().GetLoansByState(domain.StateApproved).Return(loans, nil)
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "OK",
		},
		{
			name:  "No Loans Found",
			state: "DISBURSED",
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GetLoansByState(domain.StateDisbursed).Return([]*domain.Loan{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "OK",
		},
		{
			name:           "Invalid State",
			state:          "INVALID",
			mockSetup:      func(mockService *mock.MockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Validation error",
		},
		{
			name:  "Service Error",
			state: "PROPOSED",
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GetLoansByState(domain.StateProposed).Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Failed to retrieve loans",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock.NewMockService(ctrl)
			tc.mockSetup(mockService)

			handler := NewHandler(mockService)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			e := echo.New()
			c := e.NewContext(req, rec)
			c.SetPath("/loans/state/:state")
			c.SetParamNames("state")
			c.SetParamValues(tc.state)

			// Execute
			err := handler.GetLoansByState(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, tc.expectedMsg, response["message"])
		})
	}
}

func TestGetLoans(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name           string
		page           string
		limit          string
		mockSetup      func(*mock.MockService)
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:  "Success",
			page:  "1",
			limit: "10",
			mockSetup: func(mockService *mock.MockService) {
				loans := []*domain.Loan{
					{
						ID:    "loan-1",
						State: domain.StateProposed,
					},
					{
						ID:    "loan-2",
						State: domain.StateApproved,
					},
				}
				mockService.EXPECT().GetLoans(1, 10).Return(loans, nil)
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "OK",
		},
		{
			name:  "Default Page and Limit",
			page:  "",
			limit: "",
			mockSetup: func(mockService *mock.MockService) {
				loans := []*domain.Loan{
					{
						ID:    "loan-1",
						State: domain.StateProposed,
					},
					{
						ID:    "loan-2",
						State: domain.StateApproved,
					},
				}
				mockService.EXPECT().GetLoans(1, 10).Return(loans, nil)
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "OK",
		},
		{
			name:  "Invalid Page",
			page:  "invalid",
			limit: "10",
			mockSetup: func(mockService *mock.MockService) {
				loans := []*domain.Loan{
					{
						ID:    "loan-1",
						State: domain.StateProposed,
					},
					{
						ID:    "loan-2",
						State: domain.StateApproved,
					},
				}
				mockService.EXPECT().GetLoans(1, 10).Return(loans, nil)
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "OK",
		},
		{
			name:  "Invalid Limit",
			page:  "1",
			limit: "invalid",
			mockSetup: func(mockService *mock.MockService) {
				loans := []*domain.Loan{
					{
						ID:    "loan-1",
						State: domain.StateProposed,
					},
					{
						ID:    "loan-2",
						State: domain.StateApproved,
					},
				}
				mockService.EXPECT().GetLoans(1, 10).Return(loans, nil)
			},
			expectedStatus: http.StatusOK,
			expectedMsg:    "OK",
		},
		{
			name:  "Service Error",
			page:  "1",
			limit: "10",
			mockSetup: func(mockService *mock.MockService) {
				mockService.EXPECT().GetLoans(1, 10).Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Failed to retrieve loans",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock.NewMockService(ctrl)
			tc.mockSetup(mockService)

			handler := NewHandler(mockService)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			e := echo.New()
			c := e.NewContext(req, rec)

			// Add query parameters if provided
			q := req.URL.Query()
			if tc.page != "" {
				q.Add("page", tc.page)
			}
			if tc.limit != "" {
				q.Add("limit", tc.limit)
			}
			req.URL.RawQuery = q.Encode()

			// Execute
			err := handler.GetLoans(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			// Parse response
			var response map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &response)
			assert.Equal(t, tc.expectedMsg, response["message"])
		})
	}
}
