package loan

import (
	"errors"
	"github.com/hinha/los-technical/internal/domain/response"
	"net/http"

	"github.com/go-playground/validator/v10"
	domain "github.com/hinha/los-technical/internal/domain/loan"
	"github.com/hinha/los-technical/internal/pkg/utils"
	"github.com/hinha/los-technical/internal/usecase/loan"
	"github.com/labstack/echo/v4"
)

// Handler handles HTTP requests for loan operations
type Handler struct {
	service   *loan.LoanService
	validator *validator.Validate
}

// NewHandler creates a new loan handler
func NewHandler(service *loan.LoanService) *Handler {
	validate := validator.New()

	// Register custom validation for loan state
	validate.RegisterValidation("validLoanState", utils.ValidateLoanState)

	return &Handler{
		service:   service,
		validator: validate,
	}
}

// validateLoanStateForAction validates if a loan is in the correct state for a specific action
func (h *Handler) validateLoanStateForAction(loanID string, requiredStates ...domain.LoanState) error {
	loanData, err := h.service.GetLoan(loanID)
	if err != nil {
		return err
	}

	for _, state := range requiredStates {
		if loanData.State == state {
			return nil
		}
	}

	return errors.New("loanData is not in the required state for this action")
}

// RegisterRoutes registers the loan API routes
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.POST("/loans", h.CreateLoan)
	e.GET("/loans/:id", h.GetLoan)
	e.POST("/loans/:id/approve", h.ApproveLoan)
	e.POST("/loans/:id/invest", h.AddInvestment)
	e.POST("/loans/:id/disburse", h.DisburseLoan)
	e.POST("/loans/:id/agreement", h.GenerateAgreementLetter)
	e.GET("/loans/borrower/:borrowerId", h.GetLoansByBorrower)
	e.GET("/loans/state/:state", h.GetLoansByState)
}

// CreateLoanRequest represents the request body for creating a loan
type CreateLoanRequest struct {
	BorrowerID      string  `json:"borrowerId" validate:"required"`
	PrincipalAmount float64 `json:"principalAmount" validate:"required,gt=0"`
	Rate            float64 `json:"rate" validate:"required,gte=0"`
	ROI             float64 `json:"roi" validate:"required,gte=0"`
}

// CreateLoan handles the creation of a new loan
func (h *Handler) CreateLoan(c echo.Context) error {
	var req CreateLoanRequest
	if err := c.Bind(&req); err != nil {
		return response.DefaultResponse(c, "Invalid request body", nil, nil, http.StatusBadRequest)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.DefaultResponse(c, "Validation error", nil, err.Error(), http.StatusBadRequest)
	}

	loan, err := h.service.CreateLoan(req.BorrowerID, req.PrincipalAmount, req.Rate, req.ROI)
	if err != nil {
		return response.DefaultResponse(c, "Failed to create loan", nil, err.Error(), http.StatusInternalServerError)
	}

	return response.DefaultResponse(c, "Loan created successfully", loan, nil, http.StatusCreated)
}

// GetLoan handles retrieving a loan by ID
func (h *Handler) GetLoan(c echo.Context) error {
	id := c.Param("id")

	loan, err := h.service.GetLoan(id)
	if err != nil {
		return response.DefaultResponse(c, "Loan not found", nil, err.Error(), http.StatusNotFound)
	}

	return response.DefaultResponse(c, "OK", loan, nil, http.StatusOK)
}

// ApproveLoanRequest represents the request body for approving a loan
type ApproveLoanRequest struct {
	ValidatorID string `json:"validatorId" validate:"required"`
	ProofURL    string `json:"proofUrl" validate:"required"`
}

// ApproveLoan handles the approval of a loan
func (h *Handler) ApproveLoan(c echo.Context) error {
	id := c.Param("id")

	var req ApproveLoanRequest
	if err := c.Bind(&req); err != nil {
		return response.DefaultResponse(c, "Invalid request body", nil, nil, http.StatusBadRequest)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.DefaultResponse(c, "Validation error", nil, err.Error(), http.StatusBadRequest)
	}

	// Validate loan state - must be PROPOSED to be approved
	if err := h.validateLoanStateForAction(id, domain.StateProposed); err != nil {
		return response.DefaultResponse(c, "State validation error", nil, err.Error(), http.StatusBadRequest)
	}

	if err := h.service.ApproveLoan(id, req.ValidatorID, req.ProofURL); err != nil {
		return response.DefaultResponse(c, "Failed to approve loan", nil, err.Error(), http.StatusBadRequest)
	}

	return response.DefaultResponse(c, "OK", req, nil, http.StatusOK)
}

// AddInvestmentRequest represents the request body for adding an investment
type AddInvestmentRequest struct {
	InvestorID string  `json:"investorId" validate:"required"`
	Email      string  `json:"email" validate:"required,email"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
}

// AddInvestment handles adding an investment to a loan
func (h *Handler) AddInvestment(c echo.Context) error {
	id := c.Param("id")

	var req AddInvestmentRequest
	if err := c.Bind(&req); err != nil {
		return response.DefaultResponse(c, "Invalid request body", nil, nil, http.StatusBadRequest)
	}

	if err := h.validator.Struct(req); err != nil {
		return c.String(http.StatusBadRequest, "Validation error: "+err.Error())
	}

	// Validate loan state - must be APPROVED or INVESTED to add investment
	if err := h.validateLoanStateForAction(id, domain.StateApproved, domain.StateInvested); err != nil {
		return c.String(http.StatusBadRequest, "State validation error: "+err.Error())
	}

	if err := h.service.AddInvestment(id, req.InvestorID, req.Email, req.Amount); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// DisburseLoanRequest represents the request body for disbursing a loan
type DisburseLoanRequest struct {
	FieldOfficerID  string `json:"fieldOfficerId" validate:"required"`
	SignedAgreement string `json:"signedAgreement" validate:"required"`
}

// DisburseLoan handles the disbursement of a loan
func (h *Handler) DisburseLoan(c echo.Context) error {

	id := c.Param("id")

	var req DisburseLoanRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.validator.Struct(req); err != nil {
		return c.String(http.StatusBadRequest, "Validation error: "+err.Error())
	}

	// Validate loan state - must be INVESTED to be disbursed
	if err := h.validateLoanStateForAction(id, domain.StateInvested); err != nil {
		return c.String(http.StatusBadRequest, "State validation error: "+err.Error())
	}

	if err := h.service.DisburseLoan(id, req.FieldOfficerID, req.SignedAgreement); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// GenerateAgreementLetterRequest represents the request body for generating an agreement letter
type GenerateAgreementLetterRequest struct {
	LetterURL string `json:"letterUrl" validate:"required"`
}

// GenerateAgreementLetter handles generating an agreement letter for a loan
func (h *Handler) GenerateAgreementLetter(c echo.Context) error {
	id := c.Param("id")

	var req GenerateAgreementLetterRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.validator.Struct(req); err != nil {
		return c.String(http.StatusBadRequest, "Validation error: "+err.Error())
	}

	if err := h.service.GenerateAgreementLetter(id, req.LetterURL); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

// GetLoansByBorrower handles retrieving all loans for a borrower
func (h *Handler) GetLoansByBorrower(c echo.Context) error {

	borrowerId := c.Param("borrowerId")

	loans, err := h.service.GetLoansByBorrower(borrowerId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, loans)
}

// GetLoansByStateRequest represents the request for retrieving loans by state
type GetLoansByStateRequest struct {
	State domain.LoanState `validate:"required,validLoanState"`
}

// GetLoansByState handles retrieving all loans in a specific state
func (h *Handler) GetLoansByState(c echo.Context) error {
	stateStr := c.Param("state")

	req := GetLoansByStateRequest{
		State: domain.LoanState(stateStr),
	}

	if err := h.validator.Struct(req); err != nil {
		return c.String(http.StatusBadRequest, "Validation error: "+err.Error())
	}

	loans, err := h.service.GetLoansByState(req.State)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, loans)
}
