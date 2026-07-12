package repository

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gozcu-project/gozcu/internal/netmon/policy"
)

type whitelistEntry struct {
	CIDR    string   `json:"cidr"`
	Ports   []uint16 `json:"ports"`
	Version uint32   `json:"version"`
}

type whitelistResponse struct {
	Version uint32           `json:"version"`
	Rules   []whitelistEntry `json:"rules"`
}

// WhitelistRepository — Spring Boot /api/netmon/whitelist endpoint'inden
// policy çeker ve Policy domain modeline dönüştürür.
type WhitelistRepository struct {
	client  *http.Client
	baseURL string
}

func NewWhitelistRepository(client *http.Client, baseURL string) *WhitelistRepository {
	return &WhitelistRepository{client: client, baseURL: baseURL}
}

func (r *WhitelistRepository) Load() (policy.Policy, error) {
	resp, err := r.client.Get(fmt.Sprintf("%s/api/netmon/whitelist", r.baseURL))
	if err != nil {
		return policy.Policy{}, fmt.Errorf("policy çekilemedi: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return policy.Policy{}, fmt.Errorf("backend %d döndü", resp.StatusCode)
	}

	var response whitelistResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return policy.Policy{}, fmt.Errorf("policy parse edilemedi: %w", err)
	}

	return toPolicy(response)
}

func toPolicy(response whitelistResponse) (policy.Policy, error) {
	rules := make([]policy.Rule, 0, len(response.Rules))
	for _, e := range response.Rules {
		_, cidr, err := net.ParseCIDR(e.CIDR)
		if err != nil {
			return policy.Policy{}, fmt.Errorf("geçersiz CIDR %s: %w", e.CIDR, err)
		}
		rules = append(rules, policy.Rule{
			CIDR:  cidr,
			Ports: e.Ports,
		})
	}
	return policy.Policy{
		Version: response.Version,
		Rules:   rules,
	}, nil
}
