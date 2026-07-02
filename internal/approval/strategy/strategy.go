package strategy

import "github.com/gozcu-project/gozcu/internal/approval"

type ApprovalStrategy interface {
	Execute(resp *approval.Response, args []string) error
}
