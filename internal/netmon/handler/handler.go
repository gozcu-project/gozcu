package handler

import "github.com/gozcu-project/gozcu/internal/netmon/model"

type Handler interface {
	Handle(conn model.Connection) error
}

type Chain struct {
	handlers []Handler
}

func NewChain(handlers ...Handler) *Chain {
	return &Chain{handlers: handlers}
}

func (c *Chain) Handle(conn model.Connection) error {
	for _, h := range c.handlers {
		if err := h.Handle(conn); err != nil {
			return err
		}
	}
	return nil
}
