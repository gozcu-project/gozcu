package strategy

import (
	"fmt"

	"github.com/gozcu-project/gozcu/internal/approval"
)

type TimeoutStrategy struct{}

func (s *TimeoutStrategy) Execute(resp *approval.Response, args []string) error {
	return fmt.Errorf("onay zaman aşımına uğradı (ID: %d) — lütfen tekrar deneyin", resp.ID)
}
