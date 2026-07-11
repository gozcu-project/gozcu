package repository

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gozcu-project/gozcu/internal/netmon/policy"
)

type whitelistEntry struct {
	CIDR  string   `json:"cidr"`
	Ports []uint16 `json:"ports"`
}

type WhitelistRepository struct {
	client  *http.Client
	baseURL string
}

func NewWhitelistRepository(client *http.Client, baseURL string) *WhitelistRepository {
	return &WhitelistRepository{client: client, baseURL: baseURL}
}

func (r *WhitelistRepository) Load() ([]policy.Rule, error) {
	resp, err := r.client.Get(fmt.Sprintf("%s/api/netmon/whitelist", r.baseURL))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("backend %d döndü", resp.StatusCode)
	}

	var entries []whitelistEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return nil, err
	}

	return toRules(entries)
}

func (r *WhitelistRepository) StartSync(p *policy.WhitelistPolicy, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			rules, err := r.Load()
			if err != nil {
				continue // sessizce atla, mevcut policy'yi koru
			}
			p.Update(rules)
		}
	}()
}

func toRules(entries []whitelistEntry) ([]policy.Rule, error) {
	rules := make([]policy.Rule, 0, len(entries))
	for _, e := range entries {
		_, cidr, err := net.ParseCIDR(e.CIDR)
		if err != nil {
			return nil, fmt.Errorf("geçersiz CIDR %s: %w", e.CIDR, err)
		}
		ports := make(map[uint16]struct{}, len(e.Ports))
		for _, p := range e.Ports {
			ports[p] = struct{}{}
		}
		rules = append(rules, policy.Rule{CIDR: cidr, Ports: ports})
	}
	return rules, nil
}
