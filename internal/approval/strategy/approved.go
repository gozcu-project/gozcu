package strategy

import (
	"fmt"

	"github.com/gozcu-project/gozcu/internal/approval"
	"github.com/gozcu-project/gozcu/internal/executor"
)

type ApprovedStrategy struct{}

func (s *ApprovedStrategy) Execute(resp *approval.Response, args []string) error {
	if resp.RiskLevel != "LOW" {
		fmt.Println("Onaylandı. Komut çalıştırılıyor...")
	}
	return executor.Run(args)
}
