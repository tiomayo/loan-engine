package entity

import (
	"mime/multipart"
	"time"

	"github.com/looplab/fsm"
)

type Loan struct {
	BorrowerID      int       `json:"borrower_id"`
	ValidatorID     int       `json:"validator_id"`
	PrincipalAmount float64   `json:"principal_amount"`
	Rate            int       `json:"rate"`
	ROI             float64   `json:"roi"`
	AgreementLetter string    `json:"agreement_letter"`
	ApprovalDate    time.Time `json:"-"`
	InvestedAmount  float64
	Status          string

	FSM *fsm.FSM `json:"-"`
}

type ApproveRequestDto struct {
	BorrowerID  int                   `json:"borrower_id" form:"borrower_id"`
	ValidatorID int                   `json:"validator_id" form:"validator_id"`
	Image       *multipart.FileHeader `json:"image" form:"image"`
}

type InvestRequestDto struct {
	BorrowerID    int     `json:"borrower_id"`
	InvestedValue float64 `json:"invested_value"`
}
