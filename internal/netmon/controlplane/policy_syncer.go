package controlplane

import (
	"fmt"
	"log"
	"time"

	"github.com/gozcu-project/gozcu/internal/netmon/policy"
	"github.com/gozcu-project/gozcu/internal/netmon/repository"
)

// PolicyStore — BPF map güncelleme sözleşmesi.
// store.BPFPolicyStore bu interface'i implement eder.
type PolicyStore interface {
	Replace(p policy.Policy) error
}

// PolicySyncer — Spring Boot control plane'den policy çekip
// BPFPolicyStore üzerinden kernel'e yükler.
type PolicySyncer struct {
	repo      repository.PolicyRepository
	store     PolicyStore
	validator policy.Validator
	tracker   *VersionTracker
}

func NewPolicySyncer(
	repo repository.PolicyRepository,
	store PolicyStore,
	validator policy.Validator,
	tracker *VersionTracker,
) *PolicySyncer {
	return &PolicySyncer{
		repo:      repo,
		store:     store,
		validator: validator,
		tracker:   tracker,
	}
}

// Sync — tek bir policy sync döngüsü çalıştırır.
// Hata durumunda aktif BPF policy değişmez.
func (s *PolicySyncer) Sync() error {
	p, err := s.repo.Load()
	if err != nil {
		return fmt.Errorf("policy yüklenemedi: %w", err)
	}

	if s.tracker.IsCurrent(p.Version) {
		return nil
	}

	if err := s.validator.Validate(p); err != nil {
		return fmt.Errorf("policy doğrulama başarısız: %w", err)
	}

	if err := s.store.Replace(p); err != nil {
		return fmt.Errorf("policy kernel'e yüklenemedi: %w", err)
	}

	s.tracker.Update(p.Version)
	log.Printf("[controlplane] policy güncellendi: version=%d rules=%d",
		p.Version, len(p.Rules))

	return nil
}

// StartSync — periyodik policy sync başlatır.
// Hata olursa aktif policy korunur, log yazılır, devam edilir.
func (s *PolicySyncer) StartSync(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			if err := s.Sync(); err != nil {
				log.Printf("[controlplane] sync hatası (aktif policy korunuyor): %v", err)
			}
		}
	}()
}
