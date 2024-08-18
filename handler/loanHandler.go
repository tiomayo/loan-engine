package handler

import (
	"loan-engine/entity"
	"loan-engine/repository"
	"loan-engine/states"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/looplab/fsm"
)

type handlerInterface interface {
	Propose(ctx *fiber.Ctx) error
	Approve(ctx *fiber.Ctx) error
	Invest(ctx *fiber.Ctx) error
	Disburse(ctx *fiber.Ctx) error
	List(ctx *fiber.Ctx) error
}
type Handler struct {
	repo  repository.RepoInterface
	state *fsm.FSM
}

func New(repo repository.RepoInterface, state *fsm.FSM) handlerInterface {
	return Handler{
		repo:  repo,
		state: state,
	}
}

// Approve implements handlerInterface.
func (h Handler) Approve(ctx *fiber.Ctx) error {
	reqBody := entity.ApproveRequestDto{}
	if form, err := ctx.MultipartForm(); err == nil {

		if borrowerID := form.Value["borrower_id"]; len(borrowerID) > 0 {
			reqBody.BorrowerID, _ = strconv.Atoi(borrowerID[0])
		}

		if validatorID := form.Value["validator_id"]; len(validatorID) > 0 {
			reqBody.ValidatorID, _ = strconv.Atoi(validatorID[0])
		}

		if files := form.File["image"]; len(files) > 0 {
			reqBody.Image = files[0]
		}
	}
	if reqBody.BorrowerID == 0 || reqBody.ValidatorID == 0 || reqBody.Image == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid field")
	}

	loanData := h.repo.Get(reqBody.BorrowerID)
	if loanData.Status == "" {
		return fiber.ErrNotFound
	}

	if err := loanData.FSM.Event(ctx.UserContext(), states.STATE_APPROVED); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	loanData.ValidatorID = reqBody.ValidatorID
	loanData.ApprovalDate = time.Now()
	loanData.Status = loanData.FSM.Current()

	h.repo.Set(loanData)

	return ctx.JSON(loanData)
}

// Disburse implements handlerInterface.
func (h Handler) Disburse(ctx *fiber.Ctx) error {
	panic("unimplemented")
}

// Invest implements handlerInterface.
func (h Handler) Invest(ctx *fiber.Ctx) error {
	reqBody := entity.InvestRequestDto{}
	if err := ctx.BodyParser(&reqBody); err != nil {
		return err
	}
	if reqBody.BorrowerID == 0 || reqBody.InvestedValue == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid field")
	}

	loanData := h.repo.Get(reqBody.BorrowerID)
	if loanData.Status == "" {
		return fiber.ErrNotFound
	}

	totalInvested := loanData.InvestedAmount + reqBody.InvestedValue
	if totalInvested == loanData.PrincipalAmount {
		if err := loanData.FSM.Event(ctx.UserContext(), states.STATE_INVESTED); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		// TODO: Blast email for each investor
	} else if totalInvested > loanData.PrincipalAmount {
		return fiber.NewError(fiber.StatusBadRequest, "Principal amount exceeded")
	}

	loanData.InvestedAmount = totalInvested

	loanData.Status = loanData.FSM.Current()

	h.repo.Set(loanData)

	return ctx.JSON(loanData)
}

// Propose implements handlerInterface.
func (h Handler) Propose(ctx *fiber.Ctx) error {
	newLoan := entity.Loan{
		FSM: h.state,
	}
	if err := ctx.BodyParser(&newLoan); err != nil {
		return err
	}
	if newLoan.BorrowerID == 0 || newLoan.PrincipalAmount == 0 || newLoan.Rate == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid field")
	}
	newLoan.Status = newLoan.FSM.Current()

	newLoan.FSM.Event(ctx.UserContext(), states.STATE_PROPOSED)

	h.repo.Set(newLoan)

	return ctx.JSON(newLoan)
}

// List implements handlerInterface.
func (h Handler) List(ctx *fiber.Ctx) error {
	loanList := h.repo.List()

	return ctx.JSON(loanList)
}
