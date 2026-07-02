package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gozcu-project/gozcu/internal/approval"
	"github.com/gozcu-project/gozcu/internal/approval/strategy"
	"github.com/gozcu-project/gozcu/internal/config"
	"github.com/gozcu-project/gozcu/internal/identity"
	"github.com/gozcu-project/gozcu/internal/tlsutil"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fail("kullanım: gozcu-gate <komut> [argümanlar]")
	}

	cfg := config.Load()

	httpClient, err := tlsutil.NewMTLSClient(cfg)
	if err != nil {
		fail("mTLS client kurulamadı: %v", err)
	}

	client := approval.NewClient(httpClient, cfg.BackendBaseURL)

	resp, err := client.Create(approval.CreateRequest{
		RequestedBy: identity.ResolveUser(),
		HostName:    identity.Hostname(),
		Command:     strings.Join(args, " "),
	})
	if err != nil {
		fail("onay isteği oluşturulamadı: %v", err)
	}

	if resp.Status == approval.StatusApproved && resp.RiskLevel == "LOW" {
		fmt.Println("Düşük riskli komut, otomatik onaylandı.")
	} else {
		fmt.Printf("Onay isteği gönderildi (ID: %d). Bekleniyor...\n", resp.ID)
	}

	resp, err = client.WaitForResolution(resp.ID, 2)
	if err != nil {
		fail("onay beklenirken hata: %v", err)
	}

	strat, err := strategy.Resolve(resp.Status)
	if err != nil {
		fail("%v", err)
	}

	if err := strat.Execute(resp, args); err != nil {
		fail("%v", err)
	}
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
