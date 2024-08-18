package states

import (
	"github.com/looplab/fsm"
)

const (
	STATE_PROPOSED  = "proposed"
	STATE_APPROVED  = "approved"
	STATE_INVESTED  = "invested"
	STATE_DISBURSED = "disbursed"
)

var LoanState *fsm.FSM

func StateDefinition() *fsm.FSM {
	if LoanState == nil {
		LoanState = fsm.NewFSM(
			STATE_PROPOSED,
			fsm.Events{
				{Name: STATE_PROPOSED, Src: []string{STATE_PROPOSED}, Dst: STATE_PROPOSED},
				{Name: STATE_APPROVED, Src: []string{STATE_PROPOSED}, Dst: STATE_APPROVED},
				{Name: STATE_INVESTED, Src: []string{STATE_APPROVED}, Dst: STATE_INVESTED},
				{Name: STATE_DISBURSED, Src: []string{STATE_INVESTED}, Dst: STATE_DISBURSED},
			},
			fsm.Callbacks{},
		)
	}
	return LoanState
}
