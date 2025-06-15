//go:generate mockgen -source=service.go -destination=mock/service_mock.go -package provider github.com/hinha/los-technical
package loan

// EmailSender defines the interface for sending emails
type EmailSender interface {
	SendAgreementEmail(email, loanID, agreementURL string) error
}

// Service defines the interface for loan operations
type Service interface {
	CreateLoan(borrowerID string, principal, rate, roi float64) (*Loan, error)
	ApproveLoan(id, validatorID, proofURL string) error
	AddInvestment(id, investorID, email string, amount float64) error
	DisburseLoan(id, fieldOfficerID, signedAgreement string) error
	GenerateAgreementLetter(id string, letterURL string) error
	GetLoan(id string) (*Loan, error)
	GetLoansByBorrower(borrowerID string) ([]*Loan, error)
	GetLoansByState(state LoanState) ([]*Loan, error)
	GetLoans(page, limit int) ([]*Loan, error)
}
