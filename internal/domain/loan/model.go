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
	ID              string  `json:"id"`
	BorrowerID      string  `json:"borrower_id"`
	PrincipalAmount float64 `json:"principal_amount"`
	Rate            float64 `json:"rate"`
	ROI             float64 `json:"roi"`
	AgreementLetter string  `json:"agreement_letter"`

	State         LoanState     `json:"state"`
	ApprovedInfo  *Approval     `json:"approved_info"`
	Investors     []Investor    `json:"investors"`
	DisbursedInfo *Disbursement `json:"disbursed_info"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type Approval struct {
	ValidatorID string    `json:"validator_id"`
	ProofURL    string    `json:"proof_url"`
	Date        time.Time `json:"date"`
}

type Investor struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
	Email  string  `json:"email"`
}

type Disbursement struct {
	SignedAgreement string    `json:"signed_agreement"`
	FieldOfficerID  string    `json:"field_officer_id"`
	Date            time.Time `json:"date"`
}
