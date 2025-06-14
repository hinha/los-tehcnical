//go:generate mockgen -source=repository.go -destination=mock/repository_mock.go -package provider github.com/hinha/los-technical
package loan

// LoanRepository defines the interface for loan data persistence
type LoanRepository interface {
	// Save persists a new loan to the repository
	Save(loan *Loan) error

	// FindByID retrieves a loan by its ID
	FindByID(id string) (*Loan, error)

	// Update updates an existing loan in the repository
	Update(loan *Loan) error

	// FindByBorrowerID retrieves all loans for a specific borrower
	FindByBorrowerID(borrowerID string) ([]*Loan, error)

	// FindByState retrieves all loans in a specific state
	FindByState(state LoanState) ([]*Loan, error)

	// FindAll retrieves all loans with pagination
	FindAll(page, limit int) ([]*Loan, error)
}
