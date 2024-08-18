package repository

import (
	"loan-engine/config"
	"loan-engine/entity"
)

type routesRepository struct {
	cfg *config.Cfg
}

type RepoInterface interface {
	Get(int) entity.Loan
	Set(entity.Loan) error
	List() []entity.Loan
}

func New(cfg *config.Cfg) RepoInterface {
	return routesRepository{
		cfg: cfg,
	}
}

func (r routesRepository) Get(borrowerID int) entity.Loan {
	if v, found := r.cfg.Loans[borrowerID]; found {
		return v
	}
	return entity.Loan{}
}

func (r routesRepository) Set(req entity.Loan) error {
	r.cfg.Loans[req.BorrowerID] = req
	return nil
}

// List implements RepoInterface.
func (r routesRepository) List() []entity.Loan {
	result := []entity.Loan{}
	for _, v := range r.cfg.Loans {
		result = append(result, entity.Loan{
			BorrowerID:      v.BorrowerID,
			PrincipalAmount: v.PrincipalAmount,
			Rate:            v.Rate,
			ROI:             v.ROI,
			AgreementLetter: v.AgreementLetter,
			InvestedAmount:  v.InvestedAmount,
			Status:          v.FSM.Current(),
		})
	}
	return result
}
