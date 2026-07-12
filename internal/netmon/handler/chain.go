package handler

import (
	"log"

	"github.com/gozcu-project/gozcu/internal/netmon/model"
)

// Chain — birden fazla handler'ı sırayla çalıştırır.
// Bir handler hata verirse zincir devam eder — hata loglanır.
// Bu, alert gönderme başarısızlığının diğer handler'ları
// (audit, metric) engellemesini önler.
type Chain struct {
	handlers []EventHandler
}

func NewChain(handlers ...EventHandler) *Chain {
	return &Chain{handlers: handlers}
}

func (c *Chain) Handle(event model.BlockEvent) error {
	for _, h := range c.handlers {
		if err := h.Handle(event); err != nil {
			log.Printf("[handler] hata (zincir devam ediyor): %v", err)
		}
	}
	return nil
}
