package policy

import "github.com/gozcu-project/gozcu/internal/netmon/model"

type Decision int

const (
	Allow Decision = iota
	Block
)

func (d Decision) String() string {
	return [...]string{"ALLOW", "BLOCK"}[d]
}

type Policy interface {
	Evaluate(conn model.Connection) Decision
}
