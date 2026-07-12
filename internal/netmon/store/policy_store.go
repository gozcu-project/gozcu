package store

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
	"unsafe"

	"github.com/cilium/ebpf"

	"github.com/gozcu-project/gozcu/internal/netmon/bpf"
	"github.com/gozcu-project/gozcu/internal/netmon/policy"
)

/*
 * BPFPolicyStore — whitelist policy'sini A/B double-buffer map
 * üzerinden kernel'e atomik olarak yükler.
 *
 * Atomic swap sırası:
 *   1. Pasif map'i temizle
 *   2. Yeni policy'yi pasif map'e yükle
 *   3. Policy doğrula
 *   4. active_map pointer'ını atomik olarak pasif map'e çevir
 *   5. Aktif/pasif rolleri değiştir
 *
 * Kernel hiçbir anda yarı güncellenmiş policy görmez.
 */
type BPFPolicyStore struct {
	mu     sync.Mutex
	maps   bpf.Maps
	active uint32 // 0 = whitelist_a, 1 = whitelist_b
}

// whitelistKey — kernel struct whitelist_key ile birebir eşleşmeli.
type whitelistKey struct {
	PrefixLen uint32
	IP        uint32
}

// whitelistValue — kernel struct whitelist_value ile birebir eşleşmeli.
type whitelistValue struct {
	Ports     [16]uint16
	PortCount uint8
	_         [1]byte // padding
}

func NewBPFPolicyStore(maps bpf.Maps) *BPFPolicyStore {
	return &BPFPolicyStore{
		maps:   maps,
		active: 0,
	}
}

// Replace — yeni policy'yi atomik olarak kernel'e yükler.
// Hata durumunda aktif policy değişmez.
func (s *BPFPolicyStore) Replace(p policy.Policy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	passive := s.passiveMap()

	if err := s.clearMap(passive); err != nil {
		return fmt.Errorf("pasif map temizlenemedi: %w", err)
	}

	if err := s.loadRules(passive, p.Rules); err != nil {
		return fmt.Errorf("kurallar yüklenemedi: %w", err)
	}

	if err := s.swapActive(); err != nil {
		return fmt.Errorf("policy swap başarısız: %w", err)
	}

	if err := s.updateVersion(p.Version); err != nil {
		return fmt.Errorf("version güncellenemedi: %w", err)
	}

	s.active = 1 - s.active
	return nil
}

func (s *BPFPolicyStore) passiveMap() *ebpf.Map {
	if s.active == 0 {
		return s.maps.WhitelistB
	}
	return s.maps.WhitelistA
}

func (s *BPFPolicyStore) clearMap(m *ebpf.Map) error {
	var keys []whitelistKey
	var vals []whitelistValue

	iter := m.Iterate()
	var k whitelistKey
	var v whitelistValue
	for iter.Next(&k, &v) {
		keys = append(keys, k)
		vals = append(vals, v)
	}

	for _, key := range keys {
		if err := m.Delete(key); err != nil {
			return err
		}
	}
	return nil
}

func (s *BPFPolicyStore) loadRules(m *ebpf.Map, rules []policy.Rule) error {
	for _, rule := range rules {
		key := toKernelKey(rule.CIDR)
		val := toKernelValue(rule.Ports)
		if err := m.Put(key, val); err != nil {
			return fmt.Errorf("kural yüklenemedi %s: %w", rule.CIDR, err)
		}
	}
	return nil
}

func (s *BPFPolicyStore) swapActive() error {
	key := uint32(0)
	newActive := 1 - s.active
	return s.maps.ActiveMap.Put(
		unsafe.Pointer(&key),
		unsafe.Pointer(&newActive),
	)
}

func (s *BPFPolicyStore) updateVersion(version uint32) error {
	key := uint32(0)
	return s.maps.PolicyVersion.Put(
		unsafe.Pointer(&key),
		unsafe.Pointer(&version),
	)
}

func toKernelKey(cidr *net.IPNet) whitelistKey {
	ones, _ := cidr.Mask.Size()
	ip := cidr.IP.To4()
	return whitelistKey{
		PrefixLen: uint32(ones),
		IP:        binary.LittleEndian.Uint32(ip),
	}
}

func toKernelValue(ports []uint16) whitelistValue {
	var val whitelistValue
	count := len(ports)
	if count > 16 {
		count = 16
	}
	for i := 0; i < count; i++ {
		val.Ports[i] = ports[i]
	}
	val.PortCount = uint8(count)
	return val
}
