package dispatcher

import (
	"github.com/gozcu-project/gozcu/internal/netmon/handler"
	"github.com/gozcu-project/gozcu/internal/netmon/model"
)

// Dispatcher — kernel'den gelen BlockEvent'ı handler zincirine dağıtır.
// Observer Pattern: event gelince zincir çalışır.
// Dispatcher policy kararı vermez — sadece event iletir.
type Dispatcher struct {
	chain *handler.Chain
}

func New(chain *handler.Chain) *Dispatcher {
	return &Dispatcher{chain: chain}
}

// Dispatch — BlockEvent'ı handler zincirine iletir.
func (d *Dispatcher) Dispatch(event model.BlockEvent) {
	if err := d.chain.Handle(event); err != nil {
		// Chain kendi içinde hataları loglar, buraya ulaşmamalı
		_ = err
	}
}
