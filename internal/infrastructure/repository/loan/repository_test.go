package loan

import (
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	domain "github.com/hinha/los-technical/internal/domain/loan"
)

func TestInMemoryRepository_Save(t *testing.T) {
	// Setup test logger
	logger := logrus.New()
	logger.SetOutput(logrus.StandardLogger().Out)

	// Test cases
	tests := []struct {
		name    string
		loan    *domain.Loan
		wantErr bool
	}{
		{
			name: "Save new loan successfully",
			loan: &domain.Loan{
				ID:              "loan-123",
				BorrowerID:      "borrower-123",
				PrincipalAmount: 10000,
				Rate:            0.05,
				ROI:             0.1,
				State:           domain.StateProposed,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Error when saving loan with duplicate borrower ID",
			loan: &domain.Loan{
				ID:              "loan-456",
				BorrowerID:      "borrower-123", // Same as previous test
				PrincipalAmount: 20000,
				Rate:            0.06,
				ROI:             0.12,
				State:           domain.StateProposed,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new repository for each test
			repo := NewInMemoryRepository(logger)

			// For the duplicate test, first save the initial loan
			if tt.name == "Error when saving loan with duplicate borrower ID" {
				initialLoan := &domain.Loan{
					ID:              "loan-123",
					BorrowerID:      "borrower-123",
					PrincipalAmount: 10000,
					Rate:            0.05,
					ROI:             0.1,
					State:           domain.StateProposed,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}
				_ = repo.Save(initialLoan)
			}

			// Execute the test
			err := repo.Save(tt.loan)

			// Assert the result
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify the loan was saved
				savedLoan, findErr := repo.FindByID(tt.loan.ID)
				assert.NoError(t, findErr)
				assert.Equal(t, tt.loan.ID, savedLoan.ID)
				assert.Equal(t, tt.loan.BorrowerID, savedLoan.BorrowerID)
				assert.Equal(t, tt.loan.PrincipalAmount, savedLoan.PrincipalAmount)
				assert.Equal(t, tt.loan.State, savedLoan.State)
			}
		})
	}
}

func TestInMemoryRepository_FindByID(t *testing.T) {
	// Setup test logger
	logger := logrus.New()
	logger.SetOutput(logrus.StandardLogger().Out)

	// Create a repository with some test data
	repo := NewInMemoryRepository(logger)
	testLoan := &domain.Loan{
		ID:              "loan-123",
		BorrowerID:      "borrower-123",
		PrincipalAmount: 10000,
		Rate:            0.05,
		ROI:             0.1,
		State:           domain.StateProposed,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	_ = repo.Save(testLoan)

	// Test cases
	tests := []struct {
		name    string
		id      string
		want    *domain.Loan
		wantErr bool
	}{
		{
			name:    "Find existing loan",
			id:      "loan-123",
			want:    testLoan,
			wantErr: false,
		},
		{
			name:    "Error when loan not found",
			id:      "non-existent-id",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the test
			got, err := repo.FindByID(tt.id)

			// Assert the result
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.BorrowerID, got.BorrowerID)
				assert.Equal(t, tt.want.PrincipalAmount, got.PrincipalAmount)
				assert.Equal(t, tt.want.State, got.State)
			}
		})
	}
}

func TestInMemoryRepository_Update(t *testing.T) {
	// Setup test logger
	logger := logrus.New()
	logger.SetOutput(logrus.StandardLogger().Out)

	// Test cases
	tests := []struct {
		name    string
		setup   func(*InMemoryRepository) *domain.Loan
		update  func(*domain.Loan) *domain.Loan
		wantErr bool
	}{
		{
			name: "Update existing loan successfully",
			setup: func(repo *InMemoryRepository) *domain.Loan {
				loan := &domain.Loan{
					ID:              "loan-123",
					BorrowerID:      "borrower-123",
					PrincipalAmount: 10000,
					Rate:            0.05,
					ROI:             0.1,
					State:           domain.StateProposed,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}
				_ = repo.Save(loan)
				return loan
			},
			update: func(loan *domain.Loan) *domain.Loan {
				loan.State = domain.StateApproved
				loan.ApprovedInfo = &domain.Approval{
					ValidatorID: "validator-123",
					ProofURL:    "http://example.com/proof",
					Date:        time.Now(),
				}
				return loan
			},
			wantErr: false,
		},
		{
			name: "Error when updating non-existent loan",
			setup: func(repo *InMemoryRepository) *domain.Loan {
				return &domain.Loan{
					ID:              "non-existent-id",
					BorrowerID:      "borrower-456",
					PrincipalAmount: 20000,
					Rate:            0.06,
					ROI:             0.12,
					State:           domain.StateProposed,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}
			},
			update: func(loan *domain.Loan) *domain.Loan {
				return loan // No changes needed for this test
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new repository for each test
			repo := NewInMemoryRepository(logger)

			// Setup the test data
			loan := tt.setup(repo)

			// Update the loan
			updatedLoan := tt.update(loan)

			// Execute the test
			err := repo.Update(updatedLoan)

			// Assert the result
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify the loan was updated
				savedLoan, findErr := repo.FindByID(updatedLoan.ID)
				assert.NoError(t, findErr)
				assert.Equal(t, updatedLoan.State, savedLoan.State)
				if updatedLoan.ApprovedInfo != nil {
					assert.NotNil(t, savedLoan.ApprovedInfo)
					assert.Equal(t, updatedLoan.ApprovedInfo.ValidatorID, savedLoan.ApprovedInfo.ValidatorID)
				}
			}
		})
	}
}

func TestInMemoryRepository_FindByBorrowerID(t *testing.T) {
	// Setup test logger
	logger := logrus.New()
	logger.SetOutput(logrus.StandardLogger().Out)

	// Create a repository with some test data
	repo := NewInMemoryRepository(logger)

	// Add loan for borrower-123
	loan1 := &domain.Loan{
		ID:              "loan-123",
		BorrowerID:      "borrower-123",
		PrincipalAmount: 10000,
		Rate:            0.05,
		ROI:             0.1,
		State:           domain.StateProposed,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	_ = repo.Save(loan1)

	// Add loan for borrower-456
	loan3 := &domain.Loan{
		ID:              "loan-789",
		BorrowerID:      "borrower-456",
		PrincipalAmount: 30000,
		Rate:            0.07,
		ROI:             0.14,
		State:           domain.StateProposed,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	_ = repo.Save(loan3)

	// Test cases
	tests := []struct {
		name       string
		borrowerID string
		wantCount  int
		wantErr    bool
	}{
		{
			name:       "Find loans for borrower-123",
			borrowerID: "borrower-123",
			wantCount:  1,
			wantErr:    false,
		},
		{
			name:       "Find loans for borrower-456",
			borrowerID: "borrower-456",
			wantCount:  1,
			wantErr:    false,
		},
		{
			name:       "Find loans for non-existent borrower",
			borrowerID: "non-existent-borrower",
			wantCount:  0,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the test
			got, err := repo.FindByBorrowerID(tt.borrowerID)

			// Assert the result
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCount, len(got))

				// Verify all returned loans have the correct borrower ID
				for _, loan := range got {
					assert.Equal(t, tt.borrowerID, loan.BorrowerID)
				}
			}
		})
	}
}

func TestInMemoryRepository_FindByState(t *testing.T) {
	// Setup test logger
	logger := logrus.New()
	logger.SetOutput(logrus.StandardLogger().Out)

	// Create a repository with some test data
	repo := NewInMemoryRepository(logger)

	// Add loans with different states
	loan1 := &domain.Loan{
		ID:              "loan-123",
		BorrowerID:      "borrower-123",
		PrincipalAmount: 10000,
		Rate:            0.05,
		ROI:             0.1,
		State:           domain.StateProposed,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	_ = repo.Save(loan1)

	loan2 := &domain.Loan{
		ID:              "loan-456",
		BorrowerID:      "borrower-456",
		PrincipalAmount: 20000,
		Rate:            0.06,
		ROI:             0.12,
		State:           domain.StateApproved,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	_ = repo.Save(loan2)

	loan3 := &domain.Loan{
		ID:              "loan-789",
		BorrowerID:      "borrower-789",
		PrincipalAmount: 30000,
		Rate:            0.07,
		ROI:             0.14,
		State:           domain.StateProposed,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	_ = repo.Save(loan3)

	// Test cases
	tests := []struct {
		name      string
		state     domain.LoanState
		wantCount int
		wantErr   bool
	}{
		{
			name:      "Find PROPOSED loans",
			state:     domain.StateProposed,
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:      "Find APPROVED loans",
			state:     domain.StateApproved,
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:      "Find INVESTED loans (none exist)",
			state:     domain.StateInvested,
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the test
			got, err := repo.FindByState(tt.state)

			// Assert the result
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCount, len(got))

				// Verify all returned loans have the correct state
				for _, loan := range got {
					assert.Equal(t, tt.state, loan.State)
				}
			}
		})
	}
}

func TestInMemoryRepository_FindAll(t *testing.T) {
	// Setup test logger
	logger := logrus.New()
	logger.SetOutput(logrus.StandardLogger().Out)

	// Create a repository with some test data
	repo := NewInMemoryRepository(logger)

	// Add multiple loans
	for i := 1; i <= 10; i++ {
		loan := &domain.Loan{
			ID:              fmt.Sprintf("loan-%d", i),
			BorrowerID:      fmt.Sprintf("borrower-%d", i),
			PrincipalAmount: float64(i * 10000),
			Rate:            0.05,
			ROI:             0.1,
			State:           domain.StateProposed,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		_ = repo.Save(loan)
	}

	// Test cases
	tests := []struct {
		name      string
		page      int
		limit     int
		wantCount int
		wantErr   bool
	}{
		{
			name:      "First page with 5 items",
			page:      1,
			limit:     5,
			wantCount: 5,
			wantErr:   false,
		},
		{
			name:      "Second page with 5 items",
			page:      2,
			limit:     5,
			wantCount: 5,
			wantErr:   false,
		},
		{
			name:      "Page beyond available data",
			page:      3,
			limit:     5,
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "Get all items with large limit",
			page:      1,
			limit:     20,
			wantCount: 10,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the test
			got, err := repo.FindAll(tt.page, tt.limit)

			// Assert the result
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCount, len(got))
			}
		})
	}
}
