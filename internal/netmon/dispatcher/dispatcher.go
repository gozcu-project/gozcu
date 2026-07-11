package dispatcher

import (
	"github.com/gozcu-project/gozcu/internal/netmon/handler"
	"github.com/gozcu-project/gozcu/internal/netmon/model"
	"github.com/gozcu-project/gozcu/internal/netmon/policy"
)

type Dispatcher struct {
	chains map[policy.Decision]*handler.Chain
}

func New() *Dispatcher {
	return &Dispatcher{
		chains: make(map[policy.Decision]*handler.Chain),
	}
}

func (d *Dispatcher) WithHandler(decision policy.Decision, chain *handler.Chain) *Dispatcher {
	d.chains[decision] = chain
	return d
}

func (d *Dispatcher) Dispatch(decision policy.Decision, conn model.Connection) error {
	chain, ok := d.chains[decision]
	if !ok {
		return nil
	}
	return chain.Handle(conn)
}
