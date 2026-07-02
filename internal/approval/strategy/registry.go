package strategy

import (
	"fmt"

	"github.com/gozcu-project/gozcu/internal/approval"
)

var registry = map[approval.Status]ApprovalStrategy{
	approval.StatusApproved: &ApprovedStrategy{},
	approval.StatusRejected: &RejectedStrategy{},
	approval.StatusTimeout:  &TimeoutStrategy{},
}

func Resolve(status approval.Status) (ApprovalStrategy, error) {
	s, ok := registry[status]
	if !ok {
		return nil, fmt.Errorf("bilinmeyen onay durumu: %s", status)
	}
	return s, nil
}
