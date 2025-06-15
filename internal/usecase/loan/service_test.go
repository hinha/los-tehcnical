package loan

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	domain "github.com/hinha/los-technical/internal/domain/loan"
	mock "github.com/hinha/los-technical/internal/domain/loan/mock"
)

func TestCreateLoan(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		borrowerID  string
		principal   float64
		rate        float64
		roi         float64
		mockSetup   func(*mock.MockLoanRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:       "Success",
			borrowerID: "borrower-123",
			principal:  1000.0,
			rate:       0.05,
			roi:        0.1,
			mockSetup: func(repo *mock.MockLoanRepository) {
				repo.EXPECT().Save(gomock.Any()).Return(nil)
			},
			expectError: false,
		},
		{
			name:       "Repository Error",
			borrowerID: "borrower-123",
			principal:  1000.0,
			rate:       0.05,
			roi:        0.1,
			mockSetup: func(repo *mock.MockLoanRepository) {
				repo.EXPECT().Save(gomock.Any()).Return(errors.New("database error"))
			},
			expectError: true,
			errorMsg:    "database error",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockLoanRepository(ctrl)
			mockEmailSender := mock.NewMockEmailSender(ctrl)
			logger := logrus.New()

			// Configure mocks
			tc.mockSetup(mockRepo)

			// Create service
			service := NewLoanService(mockRepo, mockEmailSender, logger)

			// Execute
			loan, err := service.CreateLoan(tc.borrowerID, tc.principal, tc.rate, tc.roi)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, loan)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, loan)
				assert.Equal(t, tc.borrowerID, loan.BorrowerID)
				assert.Equal(t, tc.principal, loan.PrincipalAmount)
				assert.Equal(t, tc.rate, loan.Rate)
				assert.Equal(t, tc.roi, loan.ROI)
				assert.Equal(t, domain.StateProposed, loan.State)
			}
		})
	}
}

func TestApproveLoan(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		loanID      string
		validatorID string
		proofURL    string
		mockSetup   func(*mock.MockLoanRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Success",
			loanID:      "loan-123",
			validatorID: "validator-123",
			proofURL:    "http://example.com/proof",
			mockSetup: func(repo *mock.MockLoanRepository) {
				loan := &domain.Loan{
					ID:    "loan-123",
					State: domain.StateProposed,
				}
				repo.EXPECT().FindByID("loan-123").Return(loan, nil)
				repo.EXPECT().Update(gomock.Any()).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "Loan Not Found",
			loanID:      "loan-123",
			validatorID: "validator-123",
			proofURL:    "http://example.com/proof",
			mockSetup: func(repo *mock.MockLoanRepository) {
				repo.EXPECT().FindByID("loan-123").Return(nil, errors.New("loan not found"))
			},
			expectError: true,
			errorMsg:    "failed to find loan",
		},
		{
			name:        "Invalid State",
			loanID:      "loan-123",
			validatorID: "validator-123",
			proofURL:    "http://example.com/proof",
			mockSetup: func(repo *mock.MockLoanRepository) {
				loan := &domain.Loan{
					ID:    "loan-123",
					State: domain.StateApproved,
				}
				repo.EXPECT().FindByID("loan-123").Return(loan, nil)
			},
			expectError: true,
			errorMsg:    "loan must be in PROPOSED state",
		},
		{
			name:        "Update Error",
			loanID:      "loan-123",
			validatorID: "validator-123",
			proofURL:    "http://example.com/proof",
			mockSetup: func(repo *mock.MockLoanRepository) {
				loan := &domain.Loan{
					ID:    "loan-123",
					State: domain.StateProposed,
				}
				repo.EXPECT().FindByID("loan-123").Return(loan, nil)
				repo.EXPECT().Update(gomock.Any()).Return(errors.New("update error"))
			},
			expectError: true,
			errorMsg:    "update error",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockLoanRepository(ctrl)
			mockEmailSender := mock.NewMockEmailSender(ctrl)
			logger := logrus.New()

			// Configure mocks
			tc.mockSetup(mockRepo)

			// Create service
			service := NewLoanService(mockRepo, mockEmailSender, logger)

			// Execute
			err := service.ApproveLoan(tc.loanID, tc.validatorID, tc.proofURL)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAddInvestment(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		loanID      string
		investorID  string
		email       string
		amount      float64
		mockSetup   func(*mock.MockLoanRepository, *mock.MockEmailSender)
		expectError bool
		errorMsg    string
	}{
		{
			name:       "Success - Partial Investment",
			loanID:     "loan-123",
			investorID: "investor-123",
			email:      "investor@example.com",
			amount:     500.0,
			mockSetup: func(repo *mock.MockLoanRepository, emailSender *mock.MockEmailSender) {
				loan := &domain.Loan{
					ID:              "loan-123",
					State:           domain.StateApproved,
					PrincipalAmount: 1000.0,
				}
				repo.EXPECT().FindByID("loan-123").Return(loan, nil)
				repo.EXPECT().Update(gomock.Any()).Return(nil)
			},
			expectError: false,
		},
		{
			name:       "Success - Full Investment",
			loanID:     "loan-123",
			investorID: "investor-123",
			email:      "investor@example.com",
			amount:     1000.0,
			mockSetup: func(repo *mock.MockLoanRepository, emailSender *mock.MockEmailSender) {
				loan := &domain.Loan{
					ID:              "loan-123",
					State:           domain.StateApproved,
					PrincipalAmount: 1000.0,
					AgreementLetter: "http://example.com/agreement",
				}
				repo.EXPECT().FindByID("loan-123").Return(loan, nil)
				emailSender.EXPECT().SendAgreementEmail("investor@example.com", "loan-123", "http://example.com/agreement").Return(nil)
				repo.EXPECT().Update(gomock.Any()).Return(nil)
			},
			expectError: false,
		},
		{
			name:       "Loan Not Found",
			loanID:     "loan-123",
			investorID: "investor-123",
			email:      "investor@example.com",
			amount:     500.0,
			mockSetup: func(repo *mock.MockLoanRepository, emailSender *mock.MockEmailSender) {
				repo.EXPECT().FindByID("loan-123").Return(nil, errors.New("loan not found"))
			},
			expectError: true,
			errorMsg:    "failed to find loan",
		},
		{
			name:       "Invalid State",
			loanID:     "loan-123",
			investorID: "investor-123",
			email:      "investor@example.com",
			amount:     500.0,
			mockSetup: func(repo *mock.MockLoanRepository, emailSender *mock.MockEmailSender) {
				loan := &domain.Loan{
					ID:    "loan-123",
					State: domain.StateProposed,
				}
				repo.EXPECT().FindByID("loan-123").Return(loan, nil)
			},
			expectError: true,
			errorMsg:    "loan must be in APPROVED or INVESTED state",
		},
		{
			name:       "Investment Exceeds Principal",
			loanID:     "loan-123",
			investorID: "investor-123",
			email:      "investor@example.com",
			amount:     1500.0,
			mockSetup: func(repo *mock.MockLoanRepository, emailSender *mock.MockEmailSender) {
				loan := &domain.Loan{
					ID:              "loan-123",
					State:           domain.StateApproved,
					PrincipalAmount: 1000.0,
				}
				repo.EXPECT().FindByID("loan-123").Return(loan, nil)
			},
			expectError: true,
			errorMsg:    "investment exceeds principal",
		},
		{
			name:       "Email Sending Error",
			loanID:     "loan-123",
			investorID: "investor-123",
			email:      "investor@example.com",
			amount:     1000.0,
			mockSetup: func(repo *mock.MockLoanRepository, emailSender *mock.MockEmailSender) {
				loan := &domain.Loan{
					ID:              "loan-123",
					State:           domain.StateApproved,
					PrincipalAmount: 1000.0,
					AgreementLetter: "http://example.com/agreement",
				}
				repo.EXPECT().FindByID("loan-123").Return(loan, nil)
				emailSender.EXPECT().SendAgreementEmail("investor@example.com", "loan-123", "http://example.com/agreement").Return(errors.New("email error"))
			},
			expectError: true,
			errorMsg:    "failed to send agreement to investor",
		},
		{
			name:       "Update Error",
			loanID:     "loan-123",
			investorID: "investor-123",
			email:      "investor@example.com",
			amount:     500.0,
			mockSetup: func(repo *mock.MockLoanRepository, emailSender *mock.MockEmailSender) {
				loan := &domain.Loan{
					ID:              "loan-123",
					State:           domain.StateApproved,
					PrincipalAmount: 1000.0,
				}
				repo.EXPECT().FindByID("loan-123").Return(loan, nil)
				repo.EXPECT().Update(gomock.Any()).Return(errors.New("update error"))
			},
			expectError: true,
			errorMsg:    "update error",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockLoanRepository(ctrl)
			mockEmailSender := mock.NewMockEmailSender(ctrl)
			logger := logrus.New()

			// Configure mocks
			tc.mockSetup(mockRepo, mockEmailSender)

			// Create service
			service := NewLoanService(mockRepo, mockEmailSender, logger)

			// Execute
			err := service.AddInvestment(tc.loanID, tc.investorID, tc.email, tc.amount)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDisburseLoan(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name            string
		loanID          string
		fieldOfficerID  string
		signedAgreement string
		mockSetup       func(*mock.MockLoanRepository)
		expectError     bool
		errorMsg        string
	}{
		{
			name:            "Success",
			loanID:          "loan-123",
			fieldOfficerID:  "officer-123",
			signedAgreement: "http://example.com/signed",
			mockSetup: func(repo *mock.MockLoanRepository) {
				loan := &domain.Loan{
					ID:    "loan-123",
					State: domain.StateInvested,
				}
				repo.EXPECT().FindByID("loan-123").Return(loan, nil)
				repo.EXPECT().Update(gomock.Any()).Return(nil)
			},
			expectError: false,
		},
		{
			name:            "Empty Loan ID",
			loanID:          "",
			fieldOfficerID:  "officer-123",
			signedAgreement: "http://example.com/signed",
			mockSetup:       func(repo *mock.MockLoanRepository) {},
			expectError:     true,
			errorMsg:        "loan ID cannot be empty",
		},
		{
			name:            "Empty Field Officer ID",
			loanID:          "loan-123",
			fieldOfficerID:  "",
			signedAgreement: "http://example.com/signed",
			mockSetup:       func(repo *mock.MockLoanRepository) {},
			expectError:     true,
			errorMsg:        "field officer ID cannot be empty",
		},
		{
			name:            "Empty Signed Agreement",
			loanID:          "loan-123",
			fieldOfficerID:  "officer-123",
			signedAgreement: "",
			mockSetup:       func(repo *mock.MockLoanRepository) {},
			expectError:     true,
			errorMsg:        "signed agreement cannot be empty",
		},
		{
			name:            "Loan Not Found",
			loanID:          "loan-123",
			fieldOfficerID:  "officer-123",
			signedAgreement: "http://example.com/signed",
			mockSetup: func(repo *mock.MockLoanRepository) {
				repo.EXPECT().FindByID("loan-123").Return(nil, errors.New("loan not found"))
			},
			expectError: true,
			errorMsg:    "failed to find loan",
		},
		{
			name:            "Invalid State",
			loanID:          "loan-123",
			fieldOfficerID:  "officer-123",
			signedAgreement: "http://example.com/signed",
			mockSetup: func(repo *mock.MockLoanRepository) {
				loan := &domain.Loan{
					ID:    "loan-123",
					State: domain.StateApproved,
				}
				repo.EXPECT().FindByID("loan-123").Return(loan, nil)
			},
			expectError: true,
			errorMsg:    "loan must be in INVESTED state",
		},
		{
			name:            "Update Error",
			loanID:          "loan-123",
			fieldOfficerID:  "officer-123",
			signedAgreement: "http://example.com/signed",
			mockSetup: func(repo *mock.MockLoanRepository) {
				loan := &domain.Loan{
					ID:    "loan-123",
					State: domain.StateInvested,
				}
				repo.EXPECT().FindByID("loan-123").Return(loan, nil)
				repo.EXPECT().Update(gomock.Any()).Return(errors.New("update error"))
			},
			expectError: true,
			errorMsg:    "update error",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockLoanRepository(ctrl)
			mockEmailSender := mock.NewMockEmailSender(ctrl)
			logger := logrus.New()

			// Configure mocks
			tc.mockSetup(mockRepo)

			// Create service
			service := NewLoanService(mockRepo, mockEmailSender, logger)

			// Execute
			err := service.DisburseLoan(tc.loanID, tc.fieldOfficerID, tc.signedAgreement)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGenerateAgreementLetter(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		loanID      string
		letterURL   string
		mockSetup   func(*mock.MockLoanRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:      "Success",
			loanID:    "loan-123",
			letterURL: "http://example.com/letter",
			mockSetup: func(repo *mock.MockLoanRepository) {
				loan := &domain.Loan{
					ID: "loan-123",
				}
				repo.EXPECT().FindByID("loan-123").Return(loan, nil)
				repo.EXPECT().Update(gomock.Any()).Return(nil)
			},
			expectError: false,
		},
		{
			name:        "Empty Loan ID",
			loanID:      "",
			letterURL:   "http://example.com/letter",
			mockSetup:   func(repo *mock.MockLoanRepository) {},
			expectError: true,
			errorMsg:    "loan ID cannot be empty",
		},
		{
			name:        "Empty Letter URL",
			loanID:      "loan-123",
			letterURL:   "",
			mockSetup:   func(repo *mock.MockLoanRepository) {},
			expectError: true,
			errorMsg:    "letter URL cannot be empty",
		},
		{
			name:      "Loan Not Found",
			loanID:    "loan-123",
			letterURL: "http://example.com/letter",
			mockSetup: func(repo *mock.MockLoanRepository) {
				repo.EXPECT().FindByID("loan-123").Return(nil, errors.New("loan not found"))
			},
			expectError: true,
			errorMsg:    "failed to find loan",
		},
		{
			name:      "Update Error",
			loanID:    "loan-123",
			letterURL: "http://example.com/letter",
			mockSetup: func(repo *mock.MockLoanRepository) {
				loan := &domain.Loan{
					ID: "loan-123",
				}
				repo.EXPECT().FindByID("loan-123").Return(loan, nil)
				repo.EXPECT().Update(gomock.Any()).Return(errors.New("update error"))
			},
			expectError: true,
			errorMsg:    "update error",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockLoanRepository(ctrl)
			mockEmailSender := mock.NewMockEmailSender(ctrl)
			logger := logrus.New()

			// Configure mocks
			tc.mockSetup(mockRepo)

			// Create service
			service := NewLoanService(mockRepo, mockEmailSender, logger)

			// Execute
			err := service.GenerateAgreementLetter(tc.loanID, tc.letterURL)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetLoan(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		loanID      string
		mockSetup   func(*mock.MockLoanRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:   "Success",
			loanID: "loan-123",
			mockSetup: func(repo *mock.MockLoanRepository) {
				loan := &domain.Loan{
					ID: "loan-123",
				}
				repo.EXPECT().FindByID("loan-123").Return(loan, nil)
			},
			expectError: false,
		},
		{
			name:        "Empty Loan ID",
			loanID:      "",
			mockSetup:   func(repo *mock.MockLoanRepository) {},
			expectError: true,
			errorMsg:    "loan ID cannot be empty",
		},
		{
			name:   "Loan Not Found",
			loanID: "loan-123",
			mockSetup: func(repo *mock.MockLoanRepository) {
				repo.EXPECT().FindByID("loan-123").Return(nil, errors.New("loan not found"))
			},
			expectError: true,
			errorMsg:    "loan not found",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockLoanRepository(ctrl)
			mockEmailSender := mock.NewMockEmailSender(ctrl)
			logger := logrus.New()

			// Configure mocks
			tc.mockSetup(mockRepo)

			// Create service
			service := NewLoanService(mockRepo, mockEmailSender, logger)

			// Execute
			loan, err := service.GetLoan(tc.loanID)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, loan)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, loan)
				assert.Equal(t, tc.loanID, loan.ID)
			}
		})
	}
}

func TestGetLoansByBorrower(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		borrowerID  string
		mockSetup   func(*mock.MockLoanRepository)
		expectError bool
		errorMsg    string
		expectedLen int
	}{
		{
			name:       "Success",
			borrowerID: "borrower-123",
			mockSetup: func(repo *mock.MockLoanRepository) {
				loans := []*domain.Loan{
					{ID: "loan-1", BorrowerID: "borrower-123"},
					{ID: "loan-2", BorrowerID: "borrower-123"},
				}
				repo.EXPECT().FindByBorrowerID("borrower-123").Return(loans, nil)
			},
			expectError: false,
			expectedLen: 2,
		},
		{
			name:        "Empty Borrower ID",
			borrowerID:  "",
			mockSetup:   func(repo *mock.MockLoanRepository) {},
			expectError: true,
			errorMsg:    "borrower ID cannot be empty",
			expectedLen: 0,
		},
		{
			name:       "Repository Error",
			borrowerID: "borrower-123",
			mockSetup: func(repo *mock.MockLoanRepository) {
				repo.EXPECT().FindByBorrowerID("borrower-123").Return(nil, errors.New("database error"))
			},
			expectError: true,
			errorMsg:    "database error",
			expectedLen: 0,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockLoanRepository(ctrl)
			mockEmailSender := mock.NewMockEmailSender(ctrl)
			logger := logrus.New()

			// Configure mocks
			tc.mockSetup(mockRepo)

			// Create service
			service := NewLoanService(mockRepo, mockEmailSender, logger)

			// Execute
			loans, err := service.GetLoansByBorrower(tc.borrowerID)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, loans)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, loans)
				assert.Len(t, loans, tc.expectedLen)
				if tc.expectedLen > 0 {
					assert.Equal(t, tc.borrowerID, loans[0].BorrowerID)
				}
			}
		})
	}
}

func TestGetLoansByState(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		state       domain.LoanState
		mockSetup   func(*mock.MockLoanRepository)
		expectError bool
		errorMsg    string
		expectedLen int
	}{
		{
			name:  "Success",
			state: domain.StateApproved,
			mockSetup: func(repo *mock.MockLoanRepository) {
				loans := []*domain.Loan{
					{ID: "loan-1", State: domain.StateApproved},
					{ID: "loan-2", State: domain.StateApproved},
				}
				repo.EXPECT().FindByState(domain.StateApproved).Return(loans, nil)
			},
			expectError: false,
			expectedLen: 2,
		},
		{
			name:  "Repository Error",
			state: domain.StateApproved,
			mockSetup: func(repo *mock.MockLoanRepository) {
				repo.EXPECT().FindByState(domain.StateApproved).Return(nil, errors.New("database error"))
			},
			expectError: true,
			errorMsg:    "database error",
			expectedLen: 0,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockLoanRepository(ctrl)
			mockEmailSender := mock.NewMockEmailSender(ctrl)
			logger := logrus.New()

			// Configure mocks
			tc.mockSetup(mockRepo)

			// Create service
			service := NewLoanService(mockRepo, mockEmailSender, logger)

			// Execute
			loans, err := service.GetLoansByState(tc.state)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, loans)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, loans)
				assert.Len(t, loans, tc.expectedLen)
				if tc.expectedLen > 0 {
					assert.Equal(t, tc.state, loans[0].State)
				}
			}
		})
	}
}

func TestGetLoans(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name        string
		page        int
		limit       int
		mockSetup   func(*mock.MockLoanRepository)
		expectError bool
		errorMsg    string
		expectedLen int
	}{
		{
			name:  "Success",
			page:  1,
			limit: 10,
			mockSetup: func(repo *mock.MockLoanRepository) {
				loans := []*domain.Loan{
					{ID: "loan-1"},
					{ID: "loan-2"},
				}
				repo.EXPECT().FindAll(1, 10).Return(loans, nil)
			},
			expectError: false,
			expectedLen: 2,
		},
		{
			name:  "Default Page and Limit",
			page:  0,
			limit: 0,
			mockSetup: func(repo *mock.MockLoanRepository) {
				loans := []*domain.Loan{
					{ID: "loan-1"},
					{ID: "loan-2"},
				}
				repo.EXPECT().FindAll(1, 10).Return(loans, nil)
			},
			expectError: false,
			expectedLen: 2,
		},
		{
			name:  "Repository Error",
			page:  1,
			limit: 10,
			mockSetup: func(repo *mock.MockLoanRepository) {
				repo.EXPECT().FindAll(1, 10).Return(nil, errors.New("database error"))
			},
			expectError: true,
			errorMsg:    "database error",
			expectedLen: 0,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock.NewMockLoanRepository(ctrl)
			mockEmailSender := mock.NewMockEmailSender(ctrl)
			logger := logrus.New()

			// Configure mocks
			tc.mockSetup(mockRepo)

			// Create service
			service := NewLoanService(mockRepo, mockEmailSender, logger)

			// Execute
			loans, err := service.GetLoans(tc.page, tc.limit)

			// Assert
			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
				assert.Nil(t, loans)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, loans)
				assert.Len(t, loans, tc.expectedLen)
			}
		})
	}
}

func TestNewLoanService(t *testing.T) {
	type args struct {
		repo        domain.LoanRepository
		emailSender domain.EmailSender
		logger      *logrus.Logger
	}
	tests := []struct {
		name string
		args args
		want domain.Service
	}{
		{
			name: "Success",
			args: args{
				repo:        nil,
				emailSender: nil,
				logger:      nil,
			},
			want: NewLoanService(nil, nil, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, NewLoanService(tt.args.repo, tt.args.emailSender, tt.args.logger), "NewLoanService(%v, %v, %v)", tt.args.repo, tt.args.emailSender, tt.args.logger)
		})
	}
}
