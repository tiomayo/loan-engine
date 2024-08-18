package config

import (
	"loan-engine/entity"
	"loan-engine/states"
	"sync"

	"github.com/looplab/fsm"
)

var once sync.Once

type Cfg struct {
	Loans map[int]entity.Loan
	State *fsm.FSM
}

var Config *Cfg

func Get() *Cfg {
	if Config == nil {
		once.Do(
			func() {
				stateDef := states.StateDefinition()
				Config = &Cfg{
					Loans: map[int]entity.Loan{},
					State: stateDef,
				}
			})
	}

	return Config
}
