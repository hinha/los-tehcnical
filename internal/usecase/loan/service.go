package loan

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	domain "github.com/hinha/los-technical/internal/domain/loan"
	"github.com/hinha/los-technical/internal/pkg/utils"
)

// EmailSender defines the interface for sending emails
type EmailSender interface {
	SendAgreementEmail(email, loanID, agreementURL string) error
}

// LoanService handles loan business logic
type LoanService struct {
	repo        domain.LoanRepository
	emailSender EmailSender
	logger      *logrus.Logger
}

// NewLoanService creates a new loan service
func NewLoanService(repo domain.LoanRepository, emailSender EmailSender, logger *logrus.Logger) *LoanService {
	return &LoanService{
		repo:        repo,
		emailSender: emailSender,
		logger:      logger,
	}
}

// CreateLoan creates a new loan in the PROPOSED state
func (s *LoanService) CreateLoan(borrowerID string, principal, rate, roi float64) (*domain.Loan, error) {
	s.logger.WithFields(logrus.Fields{
		"layer":            "service",
		"function":         "CreateLoan",
		"borrower_id":      borrowerID,
		"principal_amount": principal,
		"rate":             rate,
		"roi":              roi,
	}).Info("Creating new loan")

	loan := &domain.Loan{
		ID:              utils.GenerateUUID(),
		BorrowerID:      borrowerID,
		PrincipalAmount: principal,
		Rate:            rate,
		ROI:             roi,
		State:           domain.StateProposed,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err := s.repo.Save(loan)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "CreateLoan",
			"error":    err.Error(),
		}).Error("Failed to create loan")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"layer":    "service",
		"function": "CreateLoan",
		"loan_id":  loan.ID,
	}).Info("Loan created successfully")
	return loan, nil
}

// ApproveLoan transitions a loan from PROPOSED to APPROVED state
func (s *LoanService) ApproveLoan(id, validatorID, proofURL string) error {
	s.logger.WithFields(logrus.Fields{
		"layer":        "service",
		"function":     "ApproveLoan",
		"loan_id":      id,
		"validator_id": validatorID,
	}).Info("Approving loan")

	loan, err := s.repo.FindByID(id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "ApproveLoan",
			"loan_id":  id,
			"error":    err.Error(),
		}).Error("Failed to find loan")
		return fmt.Errorf("failed to find loan: %w", err)
	}

	if loan.State != domain.StateProposed {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "ApproveLoan",
			"loan_id":  id,
			"state":    loan.State,
		}).Error("Loan not in PROPOSED state")
		return errors.New("loan must be in PROPOSED state to be approved")
	}

	loan.ApprovedInfo = &domain.Approval{
		ValidatorID: validatorID,
		ProofURL:    proofURL,
		Date:        time.Now(),
	}
	loan.State = domain.StateApproved
	loan.UpdatedAt = time.Now()

	err = s.repo.Update(loan)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "ApproveLoan",
			"loan_id":  id,
			"error":    err.Error(),
		}).Error("Failed to update loan")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"layer":    "service",
		"function": "ApproveLoan",
		"loan_id":  id,
	}).Info("Loan approved successfully")
	return nil
}

// AddInvestment adds an investment to a loan
// If the total invested amount equals the principal, the loan transitions to INVESTED state
func (s *LoanService) AddInvestment(id, investorID, email string, amount float64) error {
	s.logger.WithFields(logrus.Fields{
		"layer":       "service",
		"function":    "AddInvestment",
		"loan_id":     id,
		"investor_id": investorID,
		"email":       email,
		"amount":      amount,
	}).Info("Adding investment to loan")

	loan, err := s.repo.FindByID(id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "AddInvestment",
			"loan_id":  id,
			"error":    err.Error(),
		}).Error("Failed to find loan")
		return fmt.Errorf("failed to find loan: %w", err)
	}

	if loan.State != domain.StateApproved && loan.State != domain.StateInvested {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "AddInvestment",
			"loan_id":  id,
			"state":    loan.State,
		}).Error("Loan not in APPROVED or INVESTED state")
		return errors.New("loan must be in APPROVED or INVESTED state to add investment")
	}

	// Calculate total invested amount including the new investment
	total := amount
	for _, inv := range loan.Investors {
		total += inv.Amount
	}

	if total > loan.PrincipalAmount {
		s.logger.WithFields(logrus.Fields{
			"layer":     "service",
			"function":  "AddInvestment",
			"loan_id":   id,
			"total":     total,
			"principal": loan.PrincipalAmount,
		}).Error("Investment exceeds principal")
		return fmt.Errorf("investment exceeds principal: %.2f > %.2f", total, loan.PrincipalAmount)
	}

	// Add the new investor
	loan.Investors = append(loan.Investors, domain.Investor{
		ID:     investorID,
		Amount: amount,
		Email:  email,
	})

	// If total equals principal, transition to INVESTED state
	if total == loan.PrincipalAmount {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "AddInvestment",
			"loan_id":  id,
		}).Info("Loan fully funded, transitioning to INVESTED state")

		loan.State = domain.StateInvested

		// Send agreement emails to all investors
		for _, inv := range loan.Investors {
			s.logger.WithFields(logrus.Fields{
				"layer":       "service",
				"function":    "AddInvestment",
				"loan_id":     id,
				"investor_id": inv.ID,
				"email":       inv.Email,
			}).Info("Sending agreement email to investor")

			if err := s.emailSender.SendAgreementEmail(inv.Email, loan.ID, loan.AgreementLetter); err != nil {
				s.logger.WithFields(logrus.Fields{
					"layer":       "service",
					"function":    "AddInvestment",
					"loan_id":     id,
					"investor_id": inv.ID,
					"error":       err.Error(),
				}).Error("Failed to send agreement email")
				return fmt.Errorf("failed to send agreement to investor %s: %w", inv.ID, err)
			}
		}
	}

	loan.UpdatedAt = time.Now()
	err = s.repo.Update(loan)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "AddInvestment",
			"loan_id":  id,
			"error":    err.Error(),
		}).Error("Failed to update loan")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"layer":    "service",
		"function": "AddInvestment",
		"loan_id":  id,
	}).Info("Investment added successfully")
	return nil
}

// DisburseLoan transitions a loan from INVESTED to DISBURSED state
func (s *LoanService) DisburseLoan(id, fieldOfficerID, signedAgreement string) error {
	s.logger.WithFields(logrus.Fields{
		"layer":            "service",
		"function":         "DisburseLoan",
		"loan_id":          id,
		"field_officer_id": fieldOfficerID,
	}).Info("Disbursing loan")

	// Validate input
	if id == "" {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "DisburseLoan",
		}).Error("Loan ID cannot be empty")
		return errors.New("loan ID cannot be empty")
	}
	if fieldOfficerID == "" {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "DisburseLoan",
			"loan_id":  id,
		}).Error("Field officer ID cannot be empty")
		return errors.New("field officer ID cannot be empty")
	}
	if signedAgreement == "" {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "DisburseLoan",
			"loan_id":  id,
		}).Error("Signed agreement cannot be empty")
		return errors.New("signed agreement cannot be empty")
	}

	loan, err := s.repo.FindByID(id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "DisburseLoan",
			"loan_id":  id,
			"error":    err.Error(),
		}).Error("Failed to find loan")
		return fmt.Errorf("failed to find loan: %w", err)
	}

	if loan.State != domain.StateInvested {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "DisburseLoan",
			"loan_id":  id,
			"state":    loan.State,
		}).Error("Loan not in INVESTED state")
		return errors.New("loan must be in INVESTED state to be disbursed")
	}

	loan.DisbursedInfo = &domain.Disbursement{
		SignedAgreement: signedAgreement,
		FieldOfficerID:  fieldOfficerID,
		Date:            time.Now(),
	}
	loan.State = domain.StateDisbursed
	loan.UpdatedAt = time.Now()

	err = s.repo.Update(loan)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "DisburseLoan",
			"loan_id":  id,
			"error":    err.Error(),
		}).Error("Failed to update loan")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"layer":    "service",
		"function": "DisburseLoan",
		"loan_id":  id,
	}).Info("Loan disbursed successfully")
	return nil
}

// GenerateAgreementLetter generates an agreement letter for a loan
func (s *LoanService) GenerateAgreementLetter(id string, letterURL string) error {
	s.logger.WithFields(logrus.Fields{
		"layer":      "service",
		"function":   "GenerateAgreementLetter",
		"loan_id":    id,
		"letter_url": letterURL,
	}).Info("Generating agreement letter")

	if id == "" {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "GenerateAgreementLetter",
		}).Error("Loan ID cannot be empty")
		return errors.New("loan ID cannot be empty")
	}
	if letterURL == "" {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "GenerateAgreementLetter",
			"loan_id":  id,
		}).Error("Letter URL cannot be empty")
		return errors.New("letter URL cannot be empty")
	}

	loan, err := s.repo.FindByID(id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "GenerateAgreementLetter",
			"loan_id":  id,
			"error":    err.Error(),
		}).Error("Failed to find loan")
		return fmt.Errorf("failed to find loan: %w", err)
	}

	loan.AgreementLetter = letterURL
	loan.UpdatedAt = time.Now()

	err = s.repo.Update(loan)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "GenerateAgreementLetter",
			"loan_id":  id,
			"error":    err.Error(),
		}).Error("Failed to update loan")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"layer":    "service",
		"function": "GenerateAgreementLetter",
		"loan_id":  id,
	}).Info("Agreement letter generated successfully")
	return nil
}

// GetLoan retrieves a loan by ID
func (s *LoanService) GetLoan(id string) (*domain.Loan, error) {
	s.logger.WithFields(logrus.Fields{
		"layer":    "service",
		"function": "GetLoan",
		"loan_id":  id,
	}).Info("Retrieving loan by ID")

	if id == "" {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "GetLoan",
		}).Error("Loan ID cannot be empty")
		return nil, errors.New("loan ID cannot be empty")
	}

	loan, err := s.repo.FindByID(id)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "GetLoan",
			"loan_id":  id,
			"error":    err.Error(),
		}).Error("Failed to find loan")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"layer":    "service",
		"function": "GetLoan",
		"loan_id":  id,
	}).Info("Loan retrieved successfully")
	return loan, nil
}

// GetLoansByBorrower retrieves all loans for a borrower
func (s *LoanService) GetLoansByBorrower(borrowerID string) ([]*domain.Loan, error) {
	s.logger.WithFields(logrus.Fields{
		"layer":       "service",
		"function":    "GetLoansByBorrower",
		"borrower_id": borrowerID,
	}).Info("Retrieving loans by borrower ID")

	if borrowerID == "" {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "GetLoansByBorrower",
		}).Error("Borrower ID cannot be empty")
		return nil, errors.New("borrower ID cannot be empty")
	}

	loans, err := s.repo.FindByBorrowerID(borrowerID)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"layer":       "service",
			"function":    "GetLoansByBorrower",
			"borrower_id": borrowerID,
			"error":       err.Error(),
		}).Error("Failed to find loans by borrower ID")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"layer":       "service",
		"function":    "GetLoansByBorrower",
		"borrower_id": borrowerID,
		"count":       len(loans),
	}).Info("Loans retrieved successfully")
	return loans, nil
}

// GetLoansByState retrieves all loans in a specific state
func (s *LoanService) GetLoansByState(state domain.LoanState) ([]*domain.Loan, error) {
	s.logger.WithFields(logrus.Fields{
		"layer":    "service",
		"function": "GetLoansByState",
		"state":    state,
	}).Info("Retrieving loans by state")

	loans, err := s.repo.FindByState(state)
	if err != nil {
		s.logger.WithFields(logrus.Fields{
			"layer":    "service",
			"function": "GetLoansByState",
			"state":    state,
			"error":    err.Error(),
		}).Error("Failed to find loans by state")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"layer":    "service",
		"function": "GetLoansByState",
		"state":    state,
		"count":    len(loans),
	}).Info("Loans retrieved successfully")
	return loans, nil
}
