package policy

import (
	"net"
	"sync"

	"github.com/gozcu-project/gozcu/internal/netmon/model"
)

type Rule struct {
	CIDR  *net.IPNet
	Ports portSet
}

type portSet map[uint16]struct{}

func (ps portSet) contains(port uint16) bool {
	_, ok := ps[port]
	return ok
}

func (ps portSet) isEmpty() bool {
	return len(ps) == 0
}

type WhitelistPolicy struct {
	mu    sync.RWMutex
	rules []Rule
}

func NewWhitelistPolicy(rules []Rule) *WhitelistPolicy {
	return &WhitelistPolicy{rules: rules}
}

func (w *WhitelistPolicy) Update(rules []Rule) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.rules = rules
}

func (w *WhitelistPolicy) Evaluate(conn model.Connection) Decision {
	w.mu.RLock()
	defer w.mu.RUnlock()
	for _, rule := range w.rules {
		if rule.matches(conn) {
			return Allow
		}
	}
	return Block
}

func (r Rule) matches(conn model.Connection) bool {
	return r.CIDR.Contains(conn.DestIP) && r.portMatches(conn.DstPort)
}

func (r Rule) portMatches(port uint16) bool {
	return r.Ports.isEmpty() || r.Ports.contains(port)
}
