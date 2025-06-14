package loan

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"

	domain "github.com/hinha/los-technical/internal/domain/loan"
)

// InMemoryRepository is a simple in-memory implementation of the LoanRepository interface
type InMemoryRepository struct {
	loans  map[string]*domain.Loan
	mutex  sync.RWMutex
	logger *logrus.Logger
}

// NewInMemoryRepository creates a new in-memory loan repository
func NewInMemoryRepository(logger *logrus.Logger) *InMemoryRepository {
	return &InMemoryRepository{
		loans:  make(map[string]*domain.Loan),
		logger: logger,
	}
}

// Save persists a new loan to the repository
func (r *InMemoryRepository) Save(loan *domain.Loan) error {
	r.logger.WithFields(logrus.Fields{
		"layer":    "repository",
		"function": "Save",
		"loan_id":  loan.ID,
	}).Info("Saving loan to repository")

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.loans[loan.ID]; exists {
		r.logger.WithFields(logrus.Fields{
			"layer":    "repository",
			"function": "Save",
			"loan_id":  loan.ID,
		}).Error("Loan already exists")
		return fmt.Errorf("loan with ID %s already exists", loan.ID)
	}

	r.loans[loan.ID] = loan
	r.logger.WithFields(logrus.Fields{
		"layer":    "repository",
		"function": "Save",
		"loan_id":  loan.ID,
	}).Info("Loan saved successfully")
	return nil
}

// FindByID retrieves a loan by its ID
func (r *InMemoryRepository) FindByID(id string) (*domain.Loan, error) {
	r.logger.WithFields(logrus.Fields{
		"layer":    "repository",
		"function": "FindByID",
		"loan_id":  id,
	}).Info("Finding loan by ID")

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	loan, exists := r.loans[id]
	if !exists {
		r.logger.WithFields(logrus.Fields{
			"layer":    "repository",
			"function": "FindByID",
			"loan_id":  id,
		}).Error("Loan not found")
		return nil, fmt.Errorf("loan with ID %s not found", id)
	}

	r.logger.WithFields(logrus.Fields{
		"layer":    "repository",
		"function": "FindByID",
		"loan_id":  id,
	}).Info("Loan found successfully")
	return loan, nil
}

// Update updates an existing loan in the repository
func (r *InMemoryRepository) Update(loan *domain.Loan) error {
	r.logger.WithFields(logrus.Fields{
		"layer":    "repository",
		"function": "Update",
		"loan_id":  loan.ID,
	}).Info("Updating loan in repository")

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.loans[loan.ID]; !exists {
		r.logger.WithFields(logrus.Fields{
			"layer":    "repository",
			"function": "Update",
			"loan_id":  loan.ID,
		}).Error("Loan not found for update")
		return fmt.Errorf("loan with ID %s not found", loan.ID)
	}

	r.loans[loan.ID] = loan
	r.logger.WithFields(logrus.Fields{
		"layer":    "repository",
		"function": "Update",
		"loan_id":  loan.ID,
	}).Info("Loan updated successfully")
	return nil
}

// FindByBorrowerID retrieves all loans for a specific borrower
func (r *InMemoryRepository) FindByBorrowerID(borrowerID string) ([]*domain.Loan, error) {
	r.logger.WithFields(logrus.Fields{
		"layer":       "repository",
		"function":    "FindByBorrowerID",
		"borrower_id": borrowerID,
	}).Info("Finding loans by borrower ID")

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*domain.Loan
	for _, loan := range r.loans {
		if loan.BorrowerID == borrowerID {
			result = append(result, loan)
		}
	}

	r.logger.WithFields(logrus.Fields{
		"layer":       "repository",
		"function":    "FindByBorrowerID",
		"borrower_id": borrowerID,
		"count":       len(result),
	}).Info("Found loans for borrower")
	return result, nil
}

// FindByState retrieves all loans in a specific state
func (r *InMemoryRepository) FindByState(state domain.LoanState) ([]*domain.Loan, error) {
	r.logger.WithFields(logrus.Fields{
		"layer":    "repository",
		"function": "FindByState",
		"state":    state,
	}).Info("Finding loans by state")

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*domain.Loan
	for _, loan := range r.loans {
		if loan.State == state {
			result = append(result, loan)
		}
	}

	r.logger.WithFields(logrus.Fields{
		"layer":    "repository",
		"function": "FindByState",
		"state":    state,
		"count":    len(result),
	}).Info("Found loans by state")
	return result, nil
}
