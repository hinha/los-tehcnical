package loan

import "time"

type LoanState string

const (
	StateProposed  LoanState = "PROPOSED"
	StateApproved  LoanState = "APPROVED"
	StateInvested  LoanState = "INVESTED"
	StateDisbursed LoanState = "DISBURSED"
)

type Loan struct {
	ID              string
	BorrowerID      string
	PrincipalAmount float64
	Rate            float64
	ROI             float64
	AgreementLetter string

	State         LoanState
	ApprovedInfo  *Approval
	Investors     []Investor
	DisbursedInfo *Disbursement
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Approval struct {
	ValidatorID string
	ProofURL    string
	Date        time.Time
}

type Investor struct {
	ID     string
	Amount float64
	Email  string
}

type Disbursement struct {
	SignedAgreement string
	FieldOfficerID  string
	Date            time.Time
}
