package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gozcu-project/gozcu/internal/approval"
	"github.com/gozcu-project/gozcu/internal/config"
	"github.com/gozcu-project/gozcu/internal/executor"
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

	if resp.Status == "APPROVED" && resp.RiskLevel == "LOW" {
		fmt.Println("Düşük riskli komut, otomatik onaylandı.")
	} else {
		fmt.Printf("Onay isteği gönderildi (ID: %d). Bekleniyor...\n", resp.ID)
	}

	resp, err = client.WaitForResolution(resp.ID, 2)
	if err != nil {
		fail("onay beklenirken hata: %v", err)
	}

	switch resp.Status {
	case "APPROVED":
		if resp.RiskLevel != "LOW" {
			fmt.Println("Onaylandı. Komut çalıştırılıyor...")
		}
		if err := executor.Run(args); err != nil {
			fail("komut çalıştırılamadı: %v", err)
		}
	case "REJECTED":
		fail("İstek reddedildi (onaylayan: %s).", resp.ResolvedBy)
	default:
		fail("Onay zaman aşımına uğradı (durum: %s). Lütfen tekrar deneyin.", resp.Status)
	}
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}