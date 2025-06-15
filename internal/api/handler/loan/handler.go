package loan

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	domain "github.com/hinha/los-technical/internal/domain/loan"
	"github.com/hinha/los-technical/internal/domain/response"
	"github.com/hinha/los-technical/internal/pkg/utils"
	"github.com/labstack/echo/v4"
)

// @title Loan API
// @version 1.0
// @description API for managing loans
// @BasePath /loan

// Handler handles HTTP requests for loan operations
type Handler struct {
	service   domain.Service
	validator *validator.Validate
}

// NewHandler creates a new loan handler
func NewHandler(service domain.Service) *Handler {
	validate := validator.New()

	// Register custom validation for loan state
	_ = validate.RegisterValidation("validLoanState", utils.ValidateLoanState)

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
	e.GET("/loans", h.GetLoans)
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
	BorrowerID      string  `json:"borrower_id" validate:"required" example:"amr-001"`
	PrincipalAmount float64 `json:"principal_amount" validate:"required,gt=0" example:"1000000"`
	Rate            float64 `json:"rate" validate:"required,gte=0" example:"12.5"`
	ROI             float64 `json:"roi" validate:"required,gte=0" example:"10"`
}

// CreateLoan handles the creation of a new loan
// @Summary Create a new loan
// @Description Creates a new loan with the given borrower and loan details
// @Tags loans
// @Accept json
// @Produce json
// @Param request body CreateLoanRequest true "Loan creation request"
// @Success 201 {object} response.Response "Loan created successfully"
// @Failure 400 {object} response.Response "Invalid request or validation error"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /loans [post]
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
// @Summary Get loan by ID
// @Description Retrieves a loan by its ID
// @Tags loans
// @Accept json
// @Produce json
// @Param id path string true "Loan ID"
// @Success 200 {object} response.Response "Loan details retrieved successfully"
// @Failure 404 {object} response.Response "Loan not found"
// @Router /loans/{id} [get]
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
	ValidatorID string `json:"validator_id" validate:"required" example:"LOS-123"`
	ProofURL    string `json:"proof_url" validate:"required" example:"https://storage.your.com/loan-proof/visit123.jpeg"`
}

// ApproveLoan handles the approval of a loan
// @Summary Approve a loan
// @Description Approves a loan with validator details
// @Tags loans
// @Accept json
// @Produce json
// @Param id path string true "Loan ID"
// @Param request body ApproveLoanRequest true "Loan approval request"
// @Success 200 {object} response.Response "Loan approved successfully"
// @Failure 400 {object} response.Response "Invalid request or state validation error"
// @Router /loans/{id}/approve [post]
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
	InvestorID string  `json:"investor_id" validate:"required" example:"investor-001"`
	Email      string  `json:"email" validate:"required,email" example:"client@mail.com"`
	Amount     float64 `json:"amount" validate:"required,gt=0" example:"100000"`
}

// AddInvestment handles adding an investment to a loan
// @Summary Add investment to loan
// @Description Adds an investment to an existing loan
// @Tags loans
// @Accept json
// @Produce json
// @Param id path string true "Loan ID"
// @Param request body AddInvestmentRequest true "Investment details"
// @Success 200 {string} string "Investment added successfully"
// @Failure 400 {string} string "Invalid request or state validation error"
// @Router /loans/{id}/invest [post]
func (h *Handler) AddInvestment(c echo.Context) error {
	id := c.Param("id")

	var req AddInvestmentRequest
	if err := c.Bind(&req); err != nil {
		return response.DefaultResponse(c, "Invalid request body", nil, nil, http.StatusBadRequest)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.DefaultResponse(c, "Validation error", nil, err.Error(), http.StatusBadRequest)
	}

	// Validate loan state - must be APPROVED or INVESTED to add investment
	if err := h.validateLoanStateForAction(id, domain.StateApproved, domain.StateInvested); err != nil {
		return response.DefaultResponse(c, "State validation error", nil, err.Error(), http.StatusBadRequest)
	}

	if err := h.service.AddInvestment(id, req.InvestorID, req.Email, req.Amount); err != nil {
		return response.DefaultResponse(c, "Failed to add investment", nil, err.Error(), http.StatusBadRequest)
	}

	return response.DefaultResponse(c, "OK", req, nil, http.StatusOK)
}

// DisburseLoanRequest represents the request body for disbursing a loan
type DisburseLoanRequest struct {
	FieldOfficerID  string `json:"field_officer_id" validate:"required" example:"OFC-001"`
	SignedAgreement string `json:"signed_agreement" validate:"required" example:"https://storage.your.com/loan-agreement/signed123.pdf"`
}

// DisburseLoan handles the disbursement of a loan
// @Description Disburses an approved and invested loan
// @Tags loans
// @Accept json
// @Produce json
// @Param id path string true "Loan ID"
// @Param request body DisburseLoanRequest true "Disbursement details"
// @Success 200 {string} string "Loan disbursed successfully"
// @Failure 400 {string} string "Invalid request or state validation error"
// @Router /loans/{id}/disburse [post]
func (h *Handler) DisburseLoan(c echo.Context) error {

	id := c.Param("id")

	var req DisburseLoanRequest
	if err := c.Bind(&req); err != nil {
		return response.DefaultResponse(c, "Invalid request body", nil, nil, http.StatusBadRequest)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.DefaultResponse(c, "Validation error", nil, err.Error(), http.StatusBadRequest)
	}

	// Validate loan state - must be INVESTED to be disbursed
	if err := h.validateLoanStateForAction(id, domain.StateInvested); err != nil {
		return response.DefaultResponse(c, "State validation error", nil, err.Error(), http.StatusBadRequest)
	}

	if err := h.service.DisburseLoan(id, req.FieldOfficerID, req.SignedAgreement); err != nil {
		return response.DefaultResponse(c, "Failed to disburse loan", nil, err.Error(), http.StatusBadRequest)
	}

	return response.DefaultResponse(c, "OK", req, nil, http.StatusOK)
}

// GenerateAgreementLetterRequest represents the request body for generating an agreement letter
type GenerateAgreementLetterRequest struct {
	LetterURL string `json:"letter_url" validate:"required"`
}

// GenerateAgreementLetter handles generating an agreement letter for a loan
// @Summary Generate agreement letter
// @Description Generates an agreement letter for a loan
// @Tags loans
// @Accept json
// @Produce json
// @Param id path string true "Loan ID"
// @Param request body GenerateAgreementLetterRequest true "Agreement letter details"
// @Success 200 {string} string "Agreement letter generated successfully"
// @Failure 400 {string} string "Invalid request"
// @Router /loans/{id}/agreement [post]
func (h *Handler) GenerateAgreementLetter(c echo.Context) error {
	id := c.Param("id")

	var req GenerateAgreementLetterRequest
	if err := c.Bind(&req); err != nil {
		return response.DefaultResponse(c, "Invalid request body", nil, nil, http.StatusBadRequest)
	}

	if err := h.validator.Struct(req); err != nil {
		return response.DefaultResponse(c, "Validation error", nil, err.Error(), http.StatusBadRequest)
	}

	if err := h.service.GenerateAgreementLetter(id, req.LetterURL); err != nil {
		return response.DefaultResponse(c, "Failed to generate agreement letter", nil, err.Error(), http.StatusBadRequest)
	}
	return response.DefaultResponse(c, "OK", nil, nil, http.StatusOK)
}

// GetLoansByBorrower handles retrieving all loans for a borrower
// @Summary Get loans by borrower
// @Description Retrieves all loans associated with a borrower
// @Tags loans
// @Accept json
// @Produce json
// @Param borrowerId path string true "Borrower ID"
// @Success 200 {array} response.Response "List of loans" "List of loans"
// @Failure 500 {string} string "Internal server error"
// @Router /loans/borrower/{borrowerId} [get]
func (h *Handler) GetLoansByBorrower(c echo.Context) error {

	borrowerId := c.Param("borrowerId")

	loans, err := h.service.GetLoansByBorrower(borrowerId)
	if err != nil {
		return response.DefaultResponse(c, "Failed to retrieve loans", nil, err.Error(), http.StatusInternalServerError)
	}

	return response.DefaultResponse(c, "OK", loans, nil, http.StatusOK)
}

// GetLoansByStateRequest represents the request for retrieving loans by state
type GetLoansByStateRequest struct {
	State domain.LoanState `validate:"required,validLoanState"`
}

// GetLoansByState handles retrieving all loans in a specific state
// @Summary Get loans by state
// @Description Retrieves all loans in a specific state
// @Tags loans
// @Accept json
// @Produce json
// @Param state path string true "Loan state" Enums(PROPOSED, APPROVED, INVESTED, DISBURSED)
// @Success 200 {array} response.Response "List of loans"
// @Failure 400 {string} string "Invalid state"
// @Failure 500 {string} string "Internal server error"
// @Router /loans/state/{state} [get]
func (h *Handler) GetLoansByState(c echo.Context) error {
	stateStr := c.Param("state")

	req := GetLoansByStateRequest{
		State: domain.LoanState(stateStr),
	}

	if err := h.validator.Struct(req); err != nil {
		return response.DefaultResponse(c, "Validation error", nil, err.Error(), http.StatusBadRequest)
	}

	loans, err := h.service.GetLoansByState(req.State)
	if err != nil {
		return response.DefaultResponse(c, "Failed to retrieve loans", nil, err.Error(), http.StatusInternalServerError)
	}

	return response.DefaultResponse(c, "OK", loans, nil, http.StatusOK)
}

// GetLoans handles retrieving all loans with pagination
// @Summary Get all loans
// @Description Retrieves all loans with pagination
// @Tags loans
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Success 200 {array} response.Response "List of loans"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /loans [get]
func (h *Handler) GetLoans(c echo.Context) error {
	// Parse page parameter
	pageStr := c.QueryParam("page")
	page := 1 // Default page
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
	}

	// Parse limit parameter
	limitStr := c.QueryParam("limit")
	limit := 10 // Default limit
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10
		}
	}

	loans, err := h.service.GetLoans(page, limit)
	if err != nil {
		return response.DefaultResponse(c, "Failed to retrieve loans", nil, err.Error(), http.StatusInternalServerError)
	}

	return response.DefaultResponse(c, "OK", loans, nil, http.StatusOK)
}
