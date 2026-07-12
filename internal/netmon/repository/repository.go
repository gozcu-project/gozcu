package repository

import "github.com/gozcu-project/gozcu/internal/netmon/policy"

// PolicyRepository — Spring Boot control plane'den policy çeken sözleşme.
// Tek sorumluluğu: Spring Boot REST API → Policy domain modeli dönüşümü.
// BPF map güncellemesi bu interface'in sorumluluğunda değildir.
type PolicyRepository interface {
	Load() (policy.Policy, error)
}
