package strategy

import (
	"fmt"

	"github.com/gozcu-project/gozcu/internal/approval"
)

type RejectedStrategy struct{}

func (s *RejectedStrategy) Execute(resp *approval.Response, args []string) error {
	return fmt.Errorf("istek reddedildi (onaylayan: %s)", resp.ResolvedBy)
}
